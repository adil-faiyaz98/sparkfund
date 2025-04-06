import os
import logging
import uuid
from fastapi import APIRouter, UploadFile, File, Form, HTTPException, BackgroundTasks
from fastapi.responses import JSONResponse
from typing import List, Optional
import aiofiles
from datetime import datetime
import json

from ..services.face_service import FaceService
from ..models.face import FaceMatchRequest, FaceMatchResponse

router = APIRouter()
logger = logging.getLogger(__name__)
face_service = FaceService()

@router.post("/match", response_model=FaceMatchResponse)
async def match_faces(
    document_id: str = Form(...),
    selfie_id: str = Form(...),
    verification_id: str = Form(...),
    document_file: UploadFile = File(...),
    selfie_file: UploadFile = File(...),
):
    """
    Match a selfie with a document photo
    """
    try:
        # Save the uploaded files
        doc_extension = os.path.splitext(document_file.filename)[1]
        selfie_extension = os.path.splitext(selfie_file.filename)[1]
        
        doc_path = f"uploads/{document_id}{doc_extension}"
        selfie_path = f"uploads/{selfie_id}{selfie_extension}"
        
        async with aiofiles.open(doc_path, 'wb') as out_file:
            content = await document_file.read()
            await out_file.write(content)
        
        async with aiofiles.open(selfie_path, 'wb') as out_file:
            content = await selfie_file.read()
            await out_file.write(content)
        
        # Match the faces
        result = await face_service.match_faces(
            document_id=document_id,
            selfie_id=selfie_id,
            verification_id=verification_id,
            document_path=doc_path,
            selfie_path=selfie_path
        )
        
        return result
    except Exception as e:
        logger.error(f"Error matching faces: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error matching faces: {str(e)}")

@router.post("/match-base64")
async def match_faces_base64(request: FaceMatchRequest):
    """
    Match faces from base64 encoded images
    """
    try:
        # Match the faces
        result = await face_service.match_faces_base64(
            document_id=request.document_id,
            selfie_id=request.selfie_id,
            verification_id=request.verification_id,
            document_image=request.document_image,
            selfie_image=request.selfie_image
        )
        
        return result
    except Exception as e:
        logger.error(f"Error matching faces: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error matching faces: {str(e)}")

@router.get("/thresholds")
async def get_face_match_thresholds():
    """
    Get face matching thresholds
    """
    return {
        "thresholds": {
            "high_confidence": 90.0,
            "medium_confidence": 75.0,
            "low_confidence": 60.0
        }
    }
