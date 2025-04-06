from pydantic import BaseModel, Field
from typing import Dict, List, Optional
from datetime import datetime
import uuid

class FaceMatchRequest(BaseModel):
    document_id: str = Field(..., description="ID of the document containing the face")
    selfie_id: str = Field(..., description="ID of the selfie")
    verification_id: str = Field(..., description="ID of the verification")
    document_image: str = Field(..., description="Base64 encoded document image")
    selfie_image: str = Field(..., description="Base64 encoded selfie image")

class FaceMatchResponse(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid.uuid4()), description="ID of the face match result")
    verification_id: str = Field(..., description="ID of the verification")
    document_id: str = Field(..., description="ID of the document")
    selfie_id: str = Field(..., description="ID of the selfie")
    is_match: bool = Field(..., description="Whether the faces match")
    confidence: float = Field(..., description="Confidence score for the match")
    created_at: datetime = Field(default_factory=datetime.now, description="When the match was created")
