import os
from typing import List
from pydantic import BaseSettings

class Settings(BaseSettings):
    # Application settings
    APP_NAME: str = "ai-service"
    DEBUG: bool = os.getenv("DEBUG", "False").lower() == "true"
    HOST: str = os.getenv("HOST", "0.0.0.0")
    PORT: int = int(os.getenv("PORT", "8000"))
    
    # Security settings
    API_KEY: str = os.getenv("API_KEY", "your-api-key")
    CORS_ORIGINS: List[str] = [
        "http://localhost",
        "http://localhost:8080",
        "http://localhost:3000",
        "http://kyc-service:8080",
    ]
    
    # Model settings
    MODEL_PATH: str = os.getenv("MODEL_PATH", "./models")
    DOCUMENT_MODEL_PATH: str = os.path.join(MODEL_PATH, "document")
    FACE_MODEL_PATH: str = os.path.join(MODEL_PATH, "face")
    RISK_MODEL_PATH: str = os.path.join(MODEL_PATH, "risk")
    ANOMALY_MODEL_PATH: str = os.path.join(MODEL_PATH, "anomaly")
    
    # Storage settings
    UPLOAD_DIR: str = os.getenv("UPLOAD_DIR", "./uploads")
    
    # Logging settings
    LOG_LEVEL: str = os.getenv("LOG_LEVEL", "INFO")
    
    class Config:
        env_file = ".env"

settings = Settings()
