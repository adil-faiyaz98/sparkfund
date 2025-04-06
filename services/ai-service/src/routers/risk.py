import logging
from fastapi import APIRouter, HTTPException
from typing import List, Optional
from datetime import datetime

from ..services.risk_service import RiskService
from ..models.risk import RiskAnalysisRequest, RiskAnalysisResponse

router = APIRouter()
logger = logging.getLogger(__name__)
risk_service = RiskService()

@router.post("/analyze", response_model=RiskAnalysisResponse)
async def analyze_risk(request: RiskAnalysisRequest):
    """
    Analyze risk based on user data and device information
    """
    try:
        # Analyze risk
        result = await risk_service.analyze_risk(
            user_id=request.user_id,
            verification_id=request.verification_id,
            device_info=request.device_info
        )
        
        return result
    except Exception as e:
        logger.error(f"Error analyzing risk: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error analyzing risk: {str(e)}")

@router.get("/factors")
async def get_risk_factors():
    """
    Get a list of risk factors
    """
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

@router.get("/levels")
async def get_risk_levels():
    """
    Get risk level thresholds
    """
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
