import logging
from fastapi import APIRouter, HTTPException
from typing import List, Optional
from datetime import datetime

from ..services.anomaly_service import AnomalyService
from ..models.anomaly import AnomalyDetectionRequest, AnomalyDetectionResponse

router = APIRouter()
logger = logging.getLogger(__name__)
anomaly_service = AnomalyService()

@router.post("/detect", response_model=AnomalyDetectionResponse)
async def detect_anomalies(request: AnomalyDetectionRequest):
    """
    Detect anomalies in user behavior
    """
    try:
        # Detect anomalies
        result = await anomaly_service.detect_anomalies(
            user_id=request.user_id,
            verification_id=request.verification_id,
            device_info=request.device_info
        )
        
        return result
    except Exception as e:
        logger.error(f"Error detecting anomalies: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error detecting anomalies: {str(e)}")

@router.get("/types")
async def get_anomaly_types():
    """
    Get a list of anomaly types
    """
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
