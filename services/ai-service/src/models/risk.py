from pydantic import BaseModel, Field
from typing import Dict, List, Optional
from datetime import datetime
import uuid

class DeviceInfo(BaseModel):
    ip_address: str = Field(..., description="IP address of the device")
    user_agent: str = Field(..., description="User agent of the device")
    device_type: Optional[str] = Field(None, description="Type of the device")
    os: Optional[str] = Field(None, description="Operating system of the device")
    browser: Optional[str] = Field(None, description="Browser of the device")
    mac_address: Optional[str] = Field(None, description="MAC address of the device")
    location: Optional[str] = Field(None, description="Location of the device")
    coordinates: Optional[str] = Field(None, description="Coordinates of the device")
    isp: Optional[str] = Field(None, description="ISP of the device")
    country_code: Optional[str] = Field(None, description="Country code of the device")
    captured_time: Optional[datetime] = Field(default_factory=datetime.now, description="When the device info was captured")

class RiskAnalysisRequest(BaseModel):
    user_id: str = Field(..., description="ID of the user")
    verification_id: str = Field(..., description="ID of the verification")
    device_info: DeviceInfo = Field(..., description="Device information")

class RiskAnalysisResponse(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid.uuid4()), description="ID of the risk analysis result")
    verification_id: str = Field(..., description="ID of the verification")
    user_id: str = Field(..., description="ID of the user")
    risk_score: float = Field(..., description="Risk score")
    risk_level: str = Field(..., description="Risk level")
    risk_factors: List[str] = Field(default_factory=list, description="Risk factors")
    device_info: DeviceInfo = Field(..., description="Device information")
    ip_address: str = Field(..., description="IP address")
    location: Optional[str] = Field(None, description="Location")
    created_at: datetime = Field(default_factory=datetime.now, description="When the analysis was created")
