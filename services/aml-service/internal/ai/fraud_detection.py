#!/usr/bin/env python3

import sys
import json
import numpy as np
from sklearn.ensemble import RandomForestClassifier
from sklearn.preprocessing import StandardScaler
from datetime import datetime
import joblib
import os

class FraudDetectionModel:
    def __init__(self, model_path='fraud_model.joblib'):
        self.model_path = model_path
        self.model = None
        self.scaler = StandardScaler()
        self.load_model()

    def load_model(self):
        if os.path.exists(self.model_path):
            self.model = joblib.load(self.model_path)
        else:
            self.model = RandomForestClassifier(
                n_estimators=100,
                max_depth=10,
                random_state=42
            )

    def save_model(self):
        joblib.dump(self.model, self.model_path)

    def preprocess_features(self, features):
        # Convert datetime to timestamp
        if isinstance(features['transaction_time'], str):
            features['transaction_time'] = datetime.fromisoformat(features['transaction_time'])
        
        # Extract features for model
        X = np.array([
            features['amount'],
            features['transaction_count'],
            features['average_amount'],
            features['amount_deviation'],
            features['user_age'],
            features['account_age'],
            features['country_risk'],
            features['ip_risk'],
            features['login_attempts'],
            features['failed_logins'],
            features['device_changes'],
            float(features['pep_status']),
            float(features['sanction_list']),
            float(features['watch_list'])
        ]).reshape(1, -1)

        # Scale features
        X = self.scaler.fit_transform(X)
        return X

    def predict(self, features):
        X = self.preprocess_features(features)
        
        # Get prediction and probabilities
        prediction = self.model.predict(X)[0]
        probabilities = self.model.predict_proba(X)[0]
        
        # Calculate risk level
        if prediction == 1:
            if probabilities[1] > 0.8:
                risk_level = "HIGH"
            elif probabilities[1] > 0.5:
                risk_level = "MEDIUM"
            else:
                risk_level = "LOW"
        else:
            risk_level = "LOW"

        # Get feature importance for explanation
        feature_importance = dict(zip([
            'amount', 'transaction_count', 'average_amount', 'amount_deviation',
            'user_age', 'account_age', 'country_risk', 'ip_risk',
            'login_attempts', 'failed_logins', 'device_changes',
            'pep_status', 'sanction_list', 'watch_list'
        ], self.model.feature_importances_))

        # Generate explanation
        risk_factors = []
        for feature, importance in feature_importance.items():
            if importance > 0.1:  # Only include significant factors
                risk_factors.append(f"{feature}: {importance:.2f}")

        return {
            "risk_level": risk_level,
            "risk_score": float(probabilities[1]),
            "explanation": f"Risk factors: {', '.join(risk_factors)}",
            "risk_factors": risk_factors,
            "confidence": float(max(probabilities))
        }

    def update(self, features_list, labels):
        X = np.array([self.preprocess_features(f)[0] for f in features_list])
        y = np.array(labels)
        
        # Fit scaler on all features
        self.scaler.fit(X)
        X = self.scaler.transform(X)
        
        # Update model
        self.model.fit(X, y)
        self.save_model()

def main():
    if len(sys.argv) < 2:
        print("Usage: python fraud_detection.py [predict|update|info] [data]")
        sys.exit(1)

    model = FraudDetectionModel()
    command = sys.argv[1]

    if command == "predict":
        if len(sys.argv) < 3:
            print("Error: Missing features data")
            sys.exit(1)
        features = json.loads(sys.argv[2])
        result = model.predict(features)
        print(json.dumps(result))

    elif command == "update":
        if len(sys.argv) < 3:
            print("Error: Missing training data")
            sys.exit(1)
        data = json.loads(sys.argv[2])
        model.update(data['features'], data['labels'])
        print(json.dumps({"status": "success"}))

    elif command == "info":
        info = {
            "model_type": "RandomForestClassifier",
            "n_estimators": model.model.n_estimators,
            "max_depth": model.model.max_depth,
            "feature_importance": dict(zip([
                'amount', 'transaction_count', 'average_amount', 'amount_deviation',
                'user_age', 'account_age', 'country_risk', 'ip_risk',
                'login_attempts', 'failed_logins', 'device_changes',
                'pep_status', 'sanction_list', 'watch_list'
            ], model.model.feature_importances_.tolist()))
        }
        print(json.dumps(info))

    else:
        print(f"Unknown command: {command}")
        sys.exit(1)

if __name__ == "__main__":
    main() 