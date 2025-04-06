import os
import logging
import uuid
from fastapi import APIRouter, UploadFile, File, Form, HTTPException, BackgroundTasks
from fastapi.responses import JSONResponse
from typing import List, Optional
import aiofiles
from datetime import datetime
import json

from ..services.document_service import DocumentService
from ..models.document import DocumentAnalysisRequest, DocumentAnalysisResponse

router = APIRouter()
logger = logging.getLogger(__name__)
document_service = DocumentService()

@router.post("/analyze", response_model=DocumentAnalysisResponse)
async def analyze_document(
    document_id: str = Form(...),
    verification_id: str = Form(...),
    document_file: UploadFile = File(...),
):
    """
    Analyze a document for authenticity and extract information
    """
    try:
        # Save the uploaded file
        file_extension = os.path.splitext(document_file.filename)[1]
        file_path = f"uploads/{document_id}{file_extension}"
        
        async with aiofiles.open(file_path, 'wb') as out_file:
            content = await document_file.read()
            await out_file.write(content)
        
        # Analyze the document
        result = await document_service.analyze_document(
            document_id=document_id,
            verification_id=verification_id,
            file_path=file_path
        )
        
        return result
    except Exception as e:
        logger.error(f"Error analyzing document: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error analyzing document: {str(e)}")

@router.post("/analyze-base64")
async def analyze_document_base64(request: DocumentAnalysisRequest):
    """
    Analyze a document from base64 encoded image
    """
    try:
        # Analyze the document
        result = await document_service.analyze_document_base64(
            document_id=request.document_id,
            verification_id=request.verification_id,
            base64_image=request.document_image
        )
        
        return result
    except Exception as e:
        logger.error(f"Error analyzing document: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error analyzing document: {str(e)}")

@router.get("/types")
async def get_document_types():
    """
    Get a list of supported document types
    """
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
