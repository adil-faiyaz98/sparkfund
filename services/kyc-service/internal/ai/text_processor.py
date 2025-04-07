from transformers import AutoTokenizer, AutoModelForSequenceClassification
from huggingface_hub import login
import torch
import numpy as np
import logging

logger = logging.getLogger(__name__)

class TextProcessor:
    def __init__(self, api_token: str = None):
        """Initialize FinBERT text processor"""
        if api_token:
            try:
                login(token=api_token)
            except Exception as e:
                logger.error(f"Error logging in to Hugging Face: {e}")
                raise

        self.tokenizer = AutoTokenizer.from_pretrained('ProsusAI/finbert')
        self.model = AutoModelForSequenceClassification.from_pretrained('ProsusAI/finbert')
        self.model.eval()

    def get_text_features(self, text: str) -> np.ndarray:
        """Extract features from text using FinBERT"""
        if not text or pd.isna(text):
            return np.zeros(1)

        inputs = self.tokenizer(
            text,
            padding=True,
            truncation=True,
            return_tensors='pt',
            max_length=512
        )

        with torch.no_grad():
            outputs = self.model(**inputs)
            # Get probability for fraudulent class (index 1)
            probs = outputs.logits.softmax(dim=1)[:, 1].numpy()
            return probs