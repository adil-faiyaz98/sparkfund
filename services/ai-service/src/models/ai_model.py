from pydantic import BaseModel, Field
from typing import Dict, List, Optional
from datetime import datetime
import uuid

class AIModelInfo(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid.uuid4()), description="ID of the AI model")
    name: str = Field(..., description="Name of the AI model")
    version: str = Field(..., description="Version of the AI model")
    type: str = Field(..., description="Type of the AI model")
    accuracy: float = Field(..., description="Accuracy of the AI model")
    last_trained_at: datetime = Field(default_factory=datetime.now, description="When the AI model was last trained")
    created_at: datetime = Field(default_factory=datetime.now, description="When the AI model was created")
    updated_at: datetime = Field(default_factory=datetime.now, description="When the AI model was last updated")

class AIModelList(BaseModel):
    models: List[AIModelInfo] = Field(..., description="List of AI models")
