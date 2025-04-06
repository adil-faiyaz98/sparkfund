from pydantic import BaseModel, Field
from typing import Dict, List, Optional
from datetime import datetime
import uuid

class DocumentAnalysisRequest(BaseModel):
    document_id: str = Field(..., description="ID of the document to analyze")
    verification_id: str = Field(..., description="ID of the verification")
    document_image: str = Field(..., description="Base64 encoded document image")

class ExtractedData(BaseModel):
    full_name: Optional[str] = None
    document_number: Optional[str] = None
    date_of_birth: Optional[str] = None
    expiry_date: Optional[str] = None
    issuing_country: Optional[str] = None
    issuing_authority: Optional[str] = None
    address: Optional[str] = None
    nationality: Optional[str] = None
    gender: Optional[str] = None
    mrz: Optional[str] = None
    additional_fields: Optional[Dict[str, str]] = None

class DocumentAnalysisResponse(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid.uuid4()), description="ID of the document analysis result")
    verification_id: str = Field(..., description="ID of the verification")
    document_id: str = Field(..., description="ID of the document")
    document_type: str = Field(..., description="Type of the document")
    is_authentic: bool = Field(..., description="Whether the document is authentic")
    confidence: float = Field(..., description="Confidence score for the authenticity")
    extracted_data: ExtractedData = Field(..., description="Data extracted from the document")
    issues: List[str] = Field(default_factory=list, description="Issues found with the document")
    created_at: datetime = Field(default_factory=datetime.now, description="When the analysis was created")
