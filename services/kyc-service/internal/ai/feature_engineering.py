import pandas as pd
import numpy as np
from datetime import datetime
from typing import Dict, Any, Optional
import logging

logger = logging.getLogger(__name__)

class FeatureEngineer:
    def __init__(self):
        self.categorical_features = ['email_domain', 'Status', 'DocumentType']
        self.numerical_features = ['age', 'RiskScore', 'notes_finbert']

    def calculate_age(self, dob: str, created_at: str) -> Optional[float]:
        """Calculate age from date of birth and creation date"""
        try:
            if not dob or not created_at:
                return None
            dob_date = pd.to_datetime(dob)
            created_date = pd.to_datetime(created_at)
            return (created_date - dob_date).days / 365.25
        except Exception as e:
            logger.warning(f"Error calculating age: {e}")
            return None

    def extract_email_domain(self, email: str) -> str:
        """Extract domain from email address"""
        try:
            if not email or '@' not in email:
                return 'unknown'
            return email.split('@')[1].lower()
        except Exception as e:
            logger.warning(f"Error extracting email domain: {e}")
            return 'unknown'

    def combine_location(self, address: Dict[str, str]) -> str:
        """Combine location fields into a single string"""
        fields = ['address', 'city', 'country', 'postal_code']
        return ', '.join(str(address.get(f, '')) for f in fields if address.get(f))

    def engineer_features(self, data: Dict[str, Any]) -> Dict[str, Any]:
        """Generate all features from raw data"""
        features = {}
        
        # Basic features
        features['age'] = self.calculate_age(
            data.get('date_of_birth'),
            data.get('created_at', datetime.now().isoformat())
        )
        features['email_domain'] = self.extract_email_domain(data.get('email'))
        features['location'] = self.combine_location(data.get('address', {}))
        
        # Risk score (invert if needed)
        risk_score = data.get('risk_score', 0)
        features['risk_score'] = 100 - risk_score if risk_score else 0
        
        # Document features
        features['document_type'] = data.get('document_type', 'unknown')
        features['document_number'] = data.get('document_number')
        
        # Personal information
        features['first_name'] = data.get('first_name')
        features['last_name'] = data.get('last_name')
        
        return features

    def get_feature_names(self) -> Dict[str, list]:
        """Get lists of feature names by type"""
        return {
            'categorical': self.categorical_features,
            'numerical': self.numerical_features,
            'text': ['location', 'document_number', 'first_name', 'last_name']
        }