import os
import logging
import random
import uuid
from datetime import datetime
from fastapi import FastAPI, Depends, HTTPException, status, File, UploadFile, Form
from fastapi.middleware.cors import CORSMiddleware
from fastapi.security import APIKeyHeader
from typing import List, Optional, Dict, Any
from pydantic import BaseModel, Field

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)
logger = logging.getLogger(__name__)

# Configuration
class Settings(BaseModel):
    API_KEY: str = os.getenv("API_KEY", "your-api-key")
    HOST: str = os.getenv("HOST", "0.0.0.0")
    PORT: int = int(os.getenv("PORT", "8000"))
    DEBUG: bool = os.getenv("DEBUG", "False").lower() == "true"
    CORS_ORIGINS: List[str] = [
        "http://localhost",
        "http://localhost:8080",
        "http://localhost:8081",
        "http://kyc-service:8080",
    ]

settings = Settings()

# Create FastAPI app
app = FastAPI(
    title="AI Service for KYC",
    description="AI-powered service for document verification, facial recognition, risk analysis, and anomaly detection",
    version="1.0.0",
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# API Key security
API_KEY_NAME = "X-API-Key"
api_key_header = APIKeyHeader(name=API_KEY_NAME, auto_error=False)

async def get_api_key(api_key_header: str = Depends(api_key_header)):
    if api_key_header == settings.API_KEY:
        return api_key_header
    raise HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="Invalid API Key",
    )

# Models
class DeviceInfo(BaseModel):
    ip_address: str
    user_agent: str
    device_type: Optional[str] = None
    os: Optional[str] = None
    browser: Optional[str] = None
    location: Optional[str] = None
    captured_time: Optional[datetime] = Field(default_factory=datetime.now)

class DocumentAnalysisRequest(BaseModel):
    document_id: str
    verification_id: str
    document_image: Optional[str] = None

class ExtractedData(BaseModel):
    full_name: Optional[str] = None
    document_number: Optional[str] = None
    date_of_birth: Optional[str] = None
    expiry_date: Optional[str] = None
    issuing_country: Optional[str] = None

class DocumentAnalysisResponse(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    verification_id: str
    document_id: str
    document_type: str
    is_authentic: bool
    confidence: float
    extracted_data: ExtractedData
    issues: List[str] = Field(default_factory=list)
    created_at: datetime = Field(default_factory=datetime.now)

class FaceMatchRequest(BaseModel):
    document_id: str
    selfie_id: str
    verification_id: str
    document_image: Optional[str] = None
    selfie_image: Optional[str] = None

class FaceMatchResponse(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    verification_id: str
    document_id: str
    selfie_id: str
    is_match: bool
    confidence: float
    created_at: datetime = Field(default_factory=datetime.now)

class RiskAnalysisRequest(BaseModel):
    user_id: str
    verification_id: str
    device_info: DeviceInfo

class RiskAnalysisResponse(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    verification_id: str
    user_id: str
    risk_score: float
    risk_level: str
    risk_factors: List[str] = Field(default_factory=list)
    device_info: DeviceInfo
    ip_address: str
    location: Optional[str] = None
    created_at: datetime = Field(default_factory=datetime.now)

class AnomalyDetectionRequest(BaseModel):
    user_id: str
    verification_id: str
    device_info: DeviceInfo

class AnomalyDetectionResponse(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    verification_id: str
    user_id: str
    is_anomaly: bool
    anomaly_score: float
    anomaly_type: Optional[str] = None
    reasons: List[str] = Field(default_factory=list)
    device_info: DeviceInfo
    created_at: datetime = Field(default_factory=datetime.now)

class AIModelInfo(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    name: str
    version: str
    type: str
    accuracy: float
    last_trained_at: datetime = Field(default_factory=datetime.now)
    created_at: datetime = Field(default_factory=datetime.now)
    updated_at: datetime = Field(default_factory=datetime.now)

class AIModelList(BaseModel):
    models: List[AIModelInfo]

# Document Verification Endpoints
@app.post("/api/v1/document/analyze", response_model=DocumentAnalysisResponse, dependencies=[Depends(get_api_key)])
async def analyze_document(document_id: str = Form(...), verification_id: str = Form(...), document_file: Optional[UploadFile] = File(None)):
    """Analyze a document for authenticity and extract information"""
    logger.info(f"Analyzing document: {document_id}")

    # Simulate document analysis
    document_type = random.choice(["PASSPORT", "DRIVERS_LICENSE", "ID_CARD", "RESIDENCE_PERMIT"])
    is_authentic = random.random() > 0.1  # 90% chance of being authentic
    confidence = 70.0 + random.random() * 25.0

    # Create extracted data
    extracted_data = ExtractedData(
        full_name="John Smith",
        document_number="X123456789",
        date_of_birth="1990-01-01",
        expiry_date="2030-01-01",
        issuing_country="United States"
    )

    # Create issues if not authentic
    issues = []
    if not is_authentic:
        issues.append("Document appears to be manipulated")
        issues.append("Security features missing")

    # Create response
    return DocumentAnalysisResponse(
        verification_id=verification_id,
        document_id=document_id,
        document_type=document_type,
        is_authentic=is_authentic,
        confidence=confidence,
        extracted_data=extracted_data,
        issues=issues
    )

@app.post("/api/v1/document/analyze-base64", response_model=DocumentAnalysisResponse, dependencies=[Depends(get_api_key)])
async def analyze_document_base64(request: DocumentAnalysisRequest):
    """Analyze a document from base64 encoded image"""
    logger.info(f"Analyzing document from base64: {request.document_id}")

    # Simulate document analysis
    document_type = random.choice(["PASSPORT", "DRIVERS_LICENSE", "ID_CARD", "RESIDENCE_PERMIT"])
    is_authentic = random.random() > 0.1  # 90% chance of being authentic
    confidence = 70.0 + random.random() * 25.0

    # Create extracted data
    extracted_data = ExtractedData(
        full_name="John Smith",
        document_number="X123456789",
        date_of_birth="1990-01-01",
        expiry_date="2030-01-01",
        issuing_country="United States"
    )

    # Create issues if not authentic
    issues = []
    if not is_authentic:
        issues.append("Document appears to be manipulated")
        issues.append("Security features missing")

    # Create response
    return DocumentAnalysisResponse(
        verification_id=request.verification_id,
        document_id=request.document_id,
        document_type=document_type,
        is_authentic=is_authentic,
        confidence=confidence,
        extracted_data=extracted_data,
        issues=issues
    )

@app.get("/api/v1/document/types", dependencies=[Depends(get_api_key)])
async def get_document_types():
    """Get a list of supported document types"""
    return {
        "document_types": [
            "PASSPORT",
            "DRIVERS_LICENSE",
            "ID_CARD",
            "RESIDENCE_PERMIT",
            "UTILITY_BILL",
            "BANK_STATEMENT"
        ]
    }

# Face Recognition Endpoints
@app.post("/api/v1/face/match", response_model=FaceMatchResponse, dependencies=[Depends(get_api_key)])
async def match_faces(
    document_id: str = Form(...),
    selfie_id: str = Form(...),
    verification_id: str = Form(...),
    document_file: Optional[UploadFile] = File(None),
    selfie_file: Optional[UploadFile] = File(None)
):
    """Match a selfie with a document photo"""
    logger.info(f"Matching faces: document_id={document_id}, selfie_id={selfie_id}")

    # Simulate face matching
    is_match = random.random() > 0.15  # 85% chance of matching

    if is_match:
        confidence = 75.0 + random.random() * 20.0  # 75-95% confidence for matches
    else:
        confidence = 30.0 + random.random() * 40.0  # 30-70% confidence for non-matches

    # Create response
    return FaceMatchResponse(
        verification_id=verification_id,
        document_id=document_id,
        selfie_id=selfie_id,
        is_match=is_match,
        confidence=confidence
    )

@app.post("/api/v1/face/match-base64", response_model=FaceMatchResponse, dependencies=[Depends(get_api_key)])
async def match_faces_base64(request: FaceMatchRequest):
    """Match faces from base64 encoded images"""
    logger.info(f"Matching faces from base64: document_id={request.document_id}, selfie_id={request.selfie_id}")

    # Simulate face matching
    is_match = random.random() > 0.15  # 85% chance of matching

    if is_match:
        confidence = 75.0 + random.random() * 20.0  # 75-95% confidence for matches
    else:
        confidence = 30.0 + random.random() * 40.0  # 30-70% confidence for non-matches

    # Create response
    return FaceMatchResponse(
        verification_id=request.verification_id,
        document_id=request.document_id,
        selfie_id=request.selfie_id,
        is_match=is_match,
        confidence=confidence
    )

@app.get("/api/v1/face/thresholds", dependencies=[Depends(get_api_key)])
async def get_face_match_thresholds():
    """Get face matching thresholds"""
    return {
        "thresholds": {
            "high_confidence": 90.0,
            "medium_confidence": 75.0,
            "low_confidence": 60.0
        }
    }

# Risk Analysis Endpoints
@app.post("/api/v1/risk/analyze", response_model=RiskAnalysisResponse, dependencies=[Depends(get_api_key)])
async def analyze_risk(request: RiskAnalysisRequest):
    """Analyze risk based on user data and device information"""
    logger.info(f"Analyzing risk: user_id={request.user_id}")

    # Simulate risk analysis
    risk_score = 5.0 + random.random() * 20.0  # 5-25% risk score

    # Determine risk level based on score
    if risk_score <= 15.0:
        risk_level = "LOW"
        risk_factors = []
    elif risk_score <= 50.0:
        risk_level = "MEDIUM"
        risk_factors = ["UNUSUAL_LOCATION"]
    else:
        risk_level = "HIGH"
        risk_factors = ["UNUSUAL_LOCATION", "MULTIPLE_ATTEMPTS", "IP_FRAUD_ASSOCIATION"]

    # Add some randomness to risk factors
    if random.random() < 0.2:
        risk_factors.append("DEVICE_CHANGE")

    if random.random() < 0.1:
        risk_factors.append("TIME_ANOMALY")

    # Create response
    return RiskAnalysisResponse(
        verification_id=request.verification_id,
        user_id=request.user_id,
        risk_score=risk_score,
        risk_level=risk_level,
        risk_factors=risk_factors,
        device_info=request.device_info,
        ip_address=request.device_info.ip_address,
        location=request.device_info.location
    )

@app.get("/api/v1/risk/factors", dependencies=[Depends(get_api_key)])
async def get_risk_factors():
    """Get a list of risk factors"""
    return {
        "risk_factors": [
            "UNUSUAL_LOCATION",
            "MULTIPLE_ATTEMPTS",
            "IP_FRAUD_ASSOCIATION",
            "DEVICE_CHANGE",
            "TIME_ANOMALY",
            "VPN_DETECTED",
            "PROXY_DETECTED",
            "TOR_DETECTED",
            "SUSPICIOUS_BEHAVIOR",
            "SANCTIONED_COUNTRY"
        ]
    }

@app.get("/api/v1/risk/levels", dependencies=[Depends(get_api_key)])
async def get_risk_levels():
    """Get risk level thresholds"""
    return {
        "risk_levels": {
            "LOW": {
                "min": 0.0,
                "max": 15.0,
                "description": "Low risk, proceed with verification"
            },
            "MEDIUM": {
                "min": 15.1,
                "max": 50.0,
                "description": "Medium risk, additional verification may be required"
            },
            "HIGH": {
                "min": 50.1,
                "max": 100.0,
                "description": "High risk, manual verification required"
            }
        }
    }

# Anomaly Detection Endpoints
@app.post("/api/v1/anomaly/detect", response_model=AnomalyDetectionResponse, dependencies=[Depends(get_api_key)])
async def detect_anomalies(request: AnomalyDetectionRequest):
    """Detect anomalies in user behavior"""
    logger.info(f"Detecting anomalies: user_id={request.user_id}")

    # Simulate anomaly detection
    is_anomaly = random.random() < 0.1  # 10% chance of anomaly

    if is_anomaly:
        anomaly_score = 70.0 + random.random() * 30.0  # 70-100% anomaly score
        anomaly_type = random.choice([
            "MULTIPLE_VERIFICATION_ATTEMPTS",
            "DIFFERENT_DEVICE",
            "UNUSUAL_TIME",
            "LOCATION_CHANGE",
            "SUSPICIOUS_ACTIVITY_PATTERN"
        ])
        reasons = [
            f"{anomaly_type} detected",
            "Multiple verification attempts in short time",
            "Different device than usual"
        ]
    else:
        anomaly_score = random.random() * 30.0  # 0-30% anomaly score
        anomaly_type = None
        reasons = []

    # Create response
    return AnomalyDetectionResponse(
        verification_id=request.verification_id,
        user_id=request.user_id,
        is_anomaly=is_anomaly,
        anomaly_score=anomaly_score,
        anomaly_type=anomaly_type,
        reasons=reasons,
        device_info=request.device_info
    )

@app.get("/api/v1/anomaly/types", dependencies=[Depends(get_api_key)])
async def get_anomaly_types():
    """Get a list of anomaly types"""
    return {
        "anomaly_types": [
            "MULTIPLE_VERIFICATION_ATTEMPTS",
            "DIFFERENT_DEVICE",
            "UNUSUAL_TIME",
            "LOCATION_CHANGE",
            "BROWSER_CHANGE",
            "RAPID_LOCATION_CHANGE",
            "SUSPICIOUS_ACTIVITY_PATTERN",
            "MULTIPLE_FAILED_ATTEMPTS"
        ]
    }

# AI Models Endpoints
@app.get("/api/v1/models", response_model=AIModelList, dependencies=[Depends(get_api_key)])
async def list_models():
    """List all AI models"""
    # Create default models
    models = [
        AIModelInfo(
            name="Document Verification Model",
            version="1.0.0",
            type="DOCUMENT",
            accuracy=0.98,
            last_trained_at=datetime.now()
        ),
        AIModelInfo(
            name="Face Recognition Model",
            version="1.0.0",
            type="FACE",
            accuracy=0.95,
            last_trained_at=datetime.now()
        ),
        AIModelInfo(
            name="Risk Analysis Model",
            version="1.0.0",
            type="RISK",
            accuracy=0.92,
            last_trained_at=datetime.now()
        ),
        AIModelInfo(
            name="Anomaly Detection Model",
            version="1.0.0",
            type="ANOMALY",
            accuracy=0.90,
            last_trained_at=datetime.now()
        )
    ]

    return AIModelList(models=models)

@app.get("/api/v1/models/{model_id}", response_model=AIModelInfo, dependencies=[Depends(get_api_key)])
async def get_model(model_id: str):
    """Get AI model information"""
    # Create a mock model
    model = AIModelInfo(
        id=model_id,
        name="Document Verification Model",
        version="1.0.0",
        type="DOCUMENT",
        accuracy=0.98,
        last_trained_at=datetime.now()
    )

    return model

@app.get("/api/v1/models/type/{model_type}", response_model=AIModelInfo, dependencies=[Depends(get_api_key)])
async def get_model_by_type(model_type: str):
    """Get latest AI model by type"""
    # Create a mock model based on type
    if model_type == "DOCUMENT":
        model = AIModelInfo(
            name="Document Verification Model",
            version="1.0.0",
            type="DOCUMENT",
            accuracy=0.98,
            last_trained_at=datetime.now()
        )
    elif model_type == "FACE":
        model = AIModelInfo(
            name="Face Recognition Model",
            version="1.0.0",
            type="FACE",
            accuracy=0.95,
            last_trained_at=datetime.now()
        )
    elif model_type == "RISK":
        model = AIModelInfo(
            name="Risk Analysis Model",
            version="1.0.0",
            type="RISK",
            accuracy=0.92,
            last_trained_at=datetime.now()
        )
    elif model_type == "ANOMALY":
        model = AIModelInfo(
            name="Anomaly Detection Model",
            version="1.0.0",
            type="ANOMALY",
            accuracy=0.90,
            last_trained_at=datetime.now()
        )
    else:
        raise HTTPException(status_code=404, detail=f"No model found for type: {model_type}")

    return model

@app.get("/health", tags=["Health"])
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "service": "ai-service",
        "version": "1.0.0",
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "src.main:app",
        host=settings.HOST,
        port=settings.PORT,
        reload=settings.DEBUG,
    )
