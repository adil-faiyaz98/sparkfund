#!/usr/bin/env python3

import sys
import json
import numpy as np
import pandas as pd
from sklearn.preprocessing import StandardScaler
from sklearn.ensemble import RandomForestClassifier
import joblib
from datetime import datetime
import time

def load_model(model_path):
    """Load the trained model and scaler from disk."""
    try:
        model, scaler = joblib.load(model_path)
        return model, scaler
    except Exception as e:
        print(f"Error loading model: {str(e)}")
        sys.exit(1)

def preprocess_features(features):
    """Convert the input features into a format suitable for the model."""
    # Convert the features dictionary to a DataFrame
    df = pd.DataFrame([features])
    
    # Extract time-based features
    if 'transaction_time' in df.columns:
        df['hour_of_day'] = pd.to_datetime(df['transaction_time']).dt.hour
        df['day_of_week'] = pd.to_datetime(df['transaction_time']).dt.dayofweek
    
    # Create derived features
    if 'transaction_amount' in df.columns:
        df['log_amount'] = np.log1p(df['transaction_amount'])
    
    # Handle categorical features
    if 'document_type' in df.columns:
        df = pd.get_dummies(df, columns=['document_type'], prefix='doc_type')
    
    # Select features for the model
    feature_columns = [
        'transaction_amount', 'log_amount', 'hour_of_day', 'day_of_week',
        'transaction_frequency', 'user_age', 'account_age', 'previous_fraud_reports',
        'document_quality', 'document_age', 'country_risk_score', 'ip_risk_score',
        'login_attempts', 'failed_logins'
    ]
    
    # Add document type dummy columns
    doc_type_columns = [col for col in df.columns if col.startswith('doc_type_')]
    feature_columns.extend(doc_type_columns)
    
    # Ensure all required features are present
    missing_features = [col for col in feature_columns if col not in df.columns]
    if missing_features:
        print(f"Missing required features: {missing_features}")
        sys.exit(1)
    
    return df[feature_columns]

def calculate_risk_level(risk_score):
    """Convert numerical risk score to risk level."""
    if risk_score < 33:
        return "low"
    elif risk_score < 67:
        return "medium"
    else:
        return "high"

def generate_explanation(prediction, features):
    """Generate a human-readable explanation of the prediction."""
    risk_factors = []
    
    # Check transaction amount
    if features['transaction_amount'].iloc[0] > 10000:
        risk_factors.append("High transaction amount")
    
    # Check document quality
    if features['document_quality'].iloc[0] < 0.7:
        risk_factors.append("Low document quality")
    
    # Check country risk
    if features['country_risk_score'].iloc[0] > 0.7:
        risk_factors.append("High-risk country")
    
    # Check login attempts
    if features['failed_logins'].iloc[0] > 3:
        risk_factors.append("Multiple failed login attempts")
    
    if not risk_factors:
        risk_factors.append("No significant risk factors identified")
    
    return " | ".join(risk_factors)

def predict(model_path, features_json):
    """Make a prediction using the loaded model."""
    # Load model and scaler
    model, scaler = load_model(model_path)
    
    # Parse input features
    try:
        features = json.loads(features_json)
    except json.JSONDecodeError:
        print("Error: Invalid JSON input")
        sys.exit(1)
    
    # Preprocess features
    X = preprocess_features(features)
    
    # Scale features
    X_scaled = scaler.transform(X)
    
    # Make prediction
    prediction_proba = model.predict_proba(X_scaled)[0]
    prediction = model.predict(X_scaled)[0]
    
    # Calculate risk score (0-100)
    risk_score = prediction_proba[1] * 100
    
    # Generate prediction result
    result = {
        "risk_level": calculate_risk_level(risk_score),
        "risk_score": float(risk_score),
        "confidence": float(prediction_proba[1]),
        "explanation": generate_explanation(prediction, X),
        "risk_factors": generate_explanation(prediction, X).split(" | "),
        "recommended_action": "Enhanced verification required" if risk_score > 66 else "Standard verification"
    }
    
    print(json.dumps(result))
    sys.exit(0)

def update(model_path, training_data_json):
    """Update the model with new training data."""
    try:
        data = json.loads(training_data_json)
        features = data['features']
        labels = data['labels']
    except json.JSONDecodeError:
        print("Error: Invalid JSON input")
        sys.exit(1)
    
    # Load existing model and scaler
    model, scaler = load_model(model_path)
    
    # Preprocess features
    X = pd.DataFrame(features)
    X = preprocess_features(X)
    
    # Scale features
    X_scaled = scaler.transform(X)
    
    # Update model
    model.fit(X_scaled, labels)
    
    # Save updated model
    joblib.dump((model, scaler), model_path)
    
    print(json.dumps({"status": "success", "message": "Model updated successfully"}))
    sys.exit(0)

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python fraud_detection.py [predict|update] <model_path> [features_json|training_data_json]")
        sys.exit(1)
    
    command = sys.argv[1]
    model_path = sys.argv[2]
    
    if command == "predict":
        if len(sys.argv) != 4:
            print("Usage: python fraud_detection.py predict <model_path> <features_json>")
            sys.exit(1)
        predict(model_path, sys.argv[3])
    elif command == "update":
        if len(sys.argv) != 4:
            print("Usage: python fraud_detection.py update <model_path> <training_data_json>")
            sys.exit(1)
        update(model_path, sys.argv[3])
    else:
        print("Invalid command. Use 'predict' or 'update'")
        sys.exit(1) 