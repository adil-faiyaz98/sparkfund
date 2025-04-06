import os
import logging
import uuid
import json
import shutil
from datetime import datetime
from typing import Dict, List, Optional
from fastapi import UploadFile

from ..models.ai_model import AIModelInfo
from ..config import settings

logger = logging.getLogger(__name__)

class ModelService:
    def __init__(self):
        self.model_path = settings.MODEL_PATH
        logger.info(f"Initializing ModelService with model path: {self.model_path}")
        
        # Ensure model directories exist
        os.makedirs(settings.DOCUMENT_MODEL_PATH, exist_ok=True)
        os.makedirs(settings.FACE_MODEL_PATH, exist_ok=True)
        os.makedirs(settings.RISK_MODEL_PATH, exist_ok=True)
        os.makedirs(settings.ANOMALY_MODEL_PATH, exist_ok=True)
        
        # Create models.json if it doesn't exist
        self.models_file = os.path.join(self.model_path, "models.json")
        if not os.path.exists(self.models_file):
            with open(self.models_file, "w") as f:
                json.dump({"models": []}, f)
    
    async def list_models(self) -> List[AIModelInfo]:
        """
        List all AI models
        """
        logger.info("Listing AI models")
        
        try:
            # Load models from JSON file
            with open(self.models_file, "r") as f:
                data = json.load(f)
            
            # Convert to AIModelInfo objects
            models = []
            for model_data in data.get("models", []):
                models.append(AIModelInfo(**model_data))
            
            # If no models found, create default models
            if not models:
                models = await self._create_default_models()
            
            return models
        except Exception as e:
            logger.error(f"Error listing models: {str(e)}")
            raise
    
    async def get_model(self, model_id: str) -> Optional[AIModelInfo]:
        """
        Get AI model information by ID
        """
        logger.info(f"Getting AI model: {model_id}")
        
        try:
            # Load models from JSON file
            with open(self.models_file, "r") as f:
                data = json.load(f)
            
            # Find model by ID
            for model_data in data.get("models", []):
                if model_data.get("id") == model_id:
                    return AIModelInfo(**model_data)
            
            return None
        except Exception as e:
            logger.error(f"Error getting model: {str(e)}")
            raise
    
    async def get_model_by_type(self, model_type: str) -> Optional[AIModelInfo]:
        """
        Get latest AI model by type
        """
        logger.info(f"Getting AI model by type: {model_type}")
        
        try:
            # Load models from JSON file
            with open(self.models_file, "r") as f:
                data = json.load(f)
            
            # Find models by type
            models_of_type = []
            for model_data in data.get("models", []):
                if model_data.get("type") == model_type:
                    models_of_type.append(AIModelInfo(**model_data))
            
            # Sort by version (assuming semantic versioning)
            if models_of_type:
                return sorted(models_of_type, key=lambda m: m.version, reverse=True)[0]
            
            return None
        except Exception as e:
            logger.error(f"Error getting model by type: {str(e)}")
            raise
    
    async def upload_model(self, model_type: str, model_name: str, model_version: str, model_file: UploadFile) -> str:
        """
        Upload a new AI model
        """
        logger.info(f"Uploading AI model: type={model_type}, name={model_name}, version={model_version}")
        
        try:
            # Generate model ID
            model_id = str(uuid.uuid4())
            
            # Determine model directory
            if model_type == "DOCUMENT":
                model_dir = settings.DOCUMENT_MODEL_PATH
            elif model_type == "FACE":
                model_dir = settings.FACE_MODEL_PATH
            elif model_type == "RISK":
                model_dir = settings.RISK_MODEL_PATH
            elif model_type == "ANOMALY":
                model_dir = settings.ANOMALY_MODEL_PATH
            else:
                raise ValueError(f"Invalid model type: {model_type}")
            
            # Create model directory if it doesn't exist
            os.makedirs(model_dir, exist_ok=True)
            
            # Save model file
            file_extension = os.path.splitext(model_file.filename)[1]
            model_file_path = os.path.join(model_dir, f"{model_id}{file_extension}")
            
            with open(model_file_path, "wb") as f:
                shutil.copyfileobj(model_file.file, f)
            
            # Create model info
            now = datetime.now()
            model_info = AIModelInfo(
                id=model_id,
                name=model_name,
                version=model_version,
                type=model_type,
                accuracy=0.95,  # Default accuracy
                last_trained_at=now,
                created_at=now,
                updated_at=now
            )
            
            # Add model to models.json
            with open(self.models_file, "r") as f:
                data = json.load(f)
            
            data["models"].append(model_info.dict())
            
            with open(self.models_file, "w") as f:
                json.dump(data, f, default=str)
            
            return model_id
        except Exception as e:
            logger.error(f"Error uploading model: {str(e)}")
            raise
    
    async def _create_default_models(self) -> List[AIModelInfo]:
        """
        Create default AI models
        """
        logger.info("Creating default AI models")
        
        try:
            now = datetime.now()
            
            # Create default models
            models = [
                AIModelInfo(
                    id=str(uuid.uuid4()),
                    name="Document Verification Model",
                    version="1.0.0",
                    type="DOCUMENT",
                    accuracy=0.98,
                    last_trained_at=now,
                    created_at=now,
                    updated_at=now
                ),
                AIModelInfo(
                    id=str(uuid.uuid4()),
                    name="Face Recognition Model",
                    version="1.0.0",
                    type="FACE",
                    accuracy=0.95,
                    last_trained_at=now,
                    created_at=now,
                    updated_at=now
                ),
                AIModelInfo(
                    id=str(uuid.uuid4()),
                    name="Risk Analysis Model",
                    version="1.0.0",
                    type="RISK",
                    accuracy=0.92,
                    last_trained_at=now,
                    created_at=now,
                    updated_at=now
                ),
                AIModelInfo(
                    id=str(uuid.uuid4()),
                    name="Anomaly Detection Model",
                    version="1.0.0",
                    type="ANOMALY",
                    accuracy=0.90,
                    last_trained_at=now,
                    created_at=now,
                    updated_at=now
                )
            ]
            
            # Save models to JSON file
            with open(self.models_file, "w") as f:
                json.dump({"models": [model.dict() for model in models]}, f, default=str)
            
            return models
        except Exception as e:
            logger.error(f"Error creating default models: {str(e)}")
            raise
