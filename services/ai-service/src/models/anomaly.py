from pydantic import BaseModel, Field
from typing import Dict, List, Optional
from datetime import datetime
import uuid

from .risk import DeviceInfo

class AnomalyDetectionRequest(BaseModel):
    user_id: str = Field(..., description="ID of the user")
    verification_id: str = Field(..., description="ID of the verification")
    device_info: DeviceInfo = Field(..., description="Device information")

class AnomalyDetectionResponse(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid.uuid4()), description="ID of the anomaly detection result")
    verification_id: str = Field(..., description="ID of the verification")
    user_id: str = Field(..., description="ID of the user")
    is_anomaly: bool = Field(..., description="Whether an anomaly was detected")
    anomaly_score: float = Field(..., description="Anomaly score")
    anomaly_type: Optional[str] = Field(None, description="Type of anomaly")
    reasons: List[str] = Field(default_factory=list, description="Reasons for the anomaly")
    device_info: DeviceInfo = Field(..., description="Device information")
    created_at: datetime = Field(default_factory=datetime.now, description="When the detection was created")
