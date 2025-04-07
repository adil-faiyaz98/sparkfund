from fastapi import FastAPI, UploadFile, File, HTTPException, Depends, Header
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import Optional, List
import os
import uvicorn
import uuid
from datetime import datetime

app = FastAPI(
    title="SparkFund AI Service",
    description="AI-powered document verification and facial recognition for KYC",
    version="1.0.0"
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Models
class HealthResponse(BaseModel):
    status: str
    service: str
    version: str

class DocumentVerificationRequest(BaseModel):
    document_type: str
    user_id: str

class DocumentVerificationResponse(BaseModel):
    id: str
    user_id: str
    document_type: str
    status: str
    confidence: float
    created_at: str
    updated_at: str

class FacialRecognitionRequest(BaseModel):
    user_id: str
    document_id: Optional[str] = None

class FacialRecognitionResponse(BaseModel):
    id: str
    user_id: str
    document_id: Optional[str] = None
    status: str
    confidence: float
    created_at: str
    updated_at: str

# Routes
@app.get("/health", response_model=HealthResponse)
async def health():
    return {
        "status": "UP",
        "service": "ai-service",
        "version": "1.0.0"
    }

@app.post("/api/v1/document/verify", response_model=DocumentVerificationResponse)
async def verify_document(
    request: DocumentVerificationRequest,
    file: UploadFile = File(...)
):
    # Mock implementation
    return {
        "id": f"doc-{uuid.uuid4()}",
        "user_id": request.user_id,
        "document_type": request.document_type,
        "status": "VERIFIED",
        "confidence": 0.95,
        "created_at": datetime.now().isoformat(),
        "updated_at": datetime.now().isoformat()
    }

@app.post("/api/v1/facial/verify", response_model=FacialRecognitionResponse)
async def verify_facial(
    request: FacialRecognitionRequest,
    file: UploadFile = File(...)
):
    # Mock implementation
    return {
        "id": f"face-{uuid.uuid4()}",
        "user_id": request.user_id,
        "document_id": request.document_id,
        "status": "VERIFIED",
        "confidence": 0.92,
        "created_at": datetime.now().isoformat(),
        "updated_at": datetime.now().isoformat()
    }

if __name__ == "__main__":
    port = int(os.getenv("PORT", "8000"))
    host = os.getenv("HOST", "0.0.0.0")
    uvicorn.run("app:app", host=host, port=port, reload=False)
