import pandas as pd
import numpy as np
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler, OneHotEncoder
from sklearn.compose import ColumnTransformer
from sklearn.pipeline import Pipeline
from transformers import AutoTokenizer, AutoModelForSequenceClassification  # For ProsusAI/finbert and embeddings
from huggingface_hub import login

# --- 0. Login to Hugging Face ---
# Replace 'your_token_here' with your actual Hugging Face API token.
# You can also set it as an environment variable HF_API_TOKEN.
# If the token is not valid, the models will not be downloaded successfully.
try:
    login(token="your_token_here")  # or use environment variable HF_API_TOKEN
except Exception as e:
    print(f"Error logging in to Hugging Face: {e}")
    print("Please make sure you have a valid API token and it is correctly set.")
    exit()

# --- 1. Data Loading and Initial Setup ---
# Replace 'your_data.csv' with the actual path to your data file.  It should be a CSV where each row represents a KYCVerification record.
try:
    data = pd.read_csv('your_data.csv')
except FileNotFoundError:
    print("Error: 'your_data.csv' not found. Please replace with the correct file path.")
    exit()


# Assuming RiskScore needs to be inverted (lower score = lower risk)
data['RiskScore'] = 100 - data['RiskScore']

# Create the target variable
data['is_fraudulent'] = np.where(data['Status'] == 'REJECTED', 1, 0)

# --- 2. Feature Engineering Functions ---

def calculate_age(row):
    if pd.isna(row['DateOfBirth']) or pd.isna(row['CreatedAt']):
        return np.nan  # Handle missing dates
    return (pd.to_datetime(row['CreatedAt']) - pd.to_datetime(row['DateOfBirth'])).days // 365

def extract_email_domain(email):
    if pd.isna(email) or '@' not in email:
        return 'unknown'
    return email.split('@')[1]

# --- 3. Apply Feature Engineering ---

data['age'] = data.apply(calculate_age, axis=1)
data['email_domain'] = data['Email'].apply(extract_email_domain)

# Combine location fields
data['location'] = data['Address'].fillna('') + ', ' + data['City'].fillna('') + ', ' + data['Country'].fillna('') + ', ' + data['PostalCode'].fillna('')

# --- 4. Text Preprocessing and Feature Extraction (with ProsusAI/finbert) ---

# Load the tokenizer and model
finbert_tokenizer = AutoTokenizer.from_pretrained('ProsusAI/finbert')
finbert_model = AutoModelForSequenceClassification.from_pretrained('ProsusAI/finbert')

def get_finbert_features(text):
    if pd.isna(text):
        return np.zeros(1)  # Return a zero vector for missing text
    inputs = finbert_tokenizer(text, padding=True, truncation=True, return_tensors='pt', max_length=512)
    with torch.no_grad():
        outputs = finbert_model(**inputs)
        # Assuming the relevant class (fraudulent) is at index 1
        return outputs.logits.softmax(dim=1)[:, 1].numpy()

# Apply text feature extraction
data['notes_finbert'] = data['Notes'].apply(lambda x: get_finbert_features(x)[0])

# --- 5. Select Features and Target ---

# Define features (excluding original text and ID fields, and the target variable)
features = [
    'age', 'email_domain', 'location', 'notes_finbert',  # Engineered
    'Status', 'DocumentType', 'DocumentNumber', 'FirstName', 'LastName',  # Structured
    'RiskScore',  # Structured
    # Add other relevant structured features (excluding 'Notes', 'Email', 'Address', 'City', 'Country', 'PostalCode' as we have engineered features from them)
]

# Include these only if they are relevant and preprocessed appropriately:
# 'UserID', 'DocumentURL',  ... other fields ...

X = data[features]
y = data['is_fraudulent']

# --- 6. Preprocessing Pipeline ---
# Define which columns need which preprocessing

categorical_features = ['email_domain', 'Status', 'DocumentType'] # Add other categorical features here
numerical_features = ['age', 'RiskScore', 'notes_finbert']

# Create transformers for preprocessing
preprocessor = ColumnTransformer(
    transformers=[
        ('num', StandardScaler(), numerical_features),  # Scale numerical
        ('cat', OneHotEncoder(handle_unknown='ignore'), categorical_features), # One-Hot encode
        ('text', 'passthrough', ['location', 'DocumentNumber', 'FirstName', 'LastName']), #Passthrough
    ],
    remainder='passthrough'  # Keep other columns as is (for now)
)

# --- 7. Data Splitting ---
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42, stratify=y)

# --- 8. Apply Preprocessing ---
# Create a pipeline that first preprocesses the data, then apply logistic regression
pipeline = Pipeline(steps=[('preprocessor', preprocessor)])

# Fit the pipeline on the training data
pipeline.fit(X_train, y_train)

# Transform training and testing data
X_train_processed = pipeline.transform(X_train)
X_test_processed = pipeline.transform(X_test)

# --- 9. Output ---
# Now you have X_train_processed and X_test_processed, which are ready to be used for training your fraud detection model.
print("Processed Training Data Shape:", X_train_processed.shape)
print("Processed Testing Data Shape:", X_test_processed.shape)
print("\nFirst 5 rows of processed training data:\n", X_train_processed[:5])

# ---  (Optional) To Get Feature Names After One-Hot Encoding ---
# Get the feature names after one-hot encoding
# categorical_features_names = pipeline.named_steps['preprocessor'].named_transformers_['cat'].get_feature_names_out(categorical_features)
# all_feature_names = list(numerical_features) + list(categorical_features_names)  + ['location', 'DocumentNumber', 'FirstName', 'LastName'] + ['notes_finbert']
# print("\nFeature names after preprocessing:\n", all_feature_names)
