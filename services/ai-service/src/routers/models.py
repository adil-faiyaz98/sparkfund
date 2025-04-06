import logging
from fastapi import APIRouter, HTTPException, UploadFile, File, Form
from typing import List, Optional
from datetime import datetime

from ..services.model_service import ModelService
from ..models.ai_model import AIModelInfo, AIModelList

router = APIRouter()
logger = logging.getLogger(__name__)
model_service = ModelService()

@router.get("/", response_model=AIModelList)
async def list_models():
    """
    List all AI models
    """
    try:
        models = await model_service.list_models()
        return {"models": models}
    except Exception as e:
        logger.error(f"Error listing models: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error listing models: {str(e)}")

@router.get("/{model_id}", response_model=AIModelInfo)
async def get_model(model_id: str):
    """
    Get AI model information
    """
    try:
        model = await model_service.get_model(model_id)
        if not model:
            raise HTTPException(status_code=404, detail="Model not found")
        return model
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error getting model: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error getting model: {str(e)}")

@router.get("/type/{model_type}", response_model=AIModelInfo)
async def get_model_by_type(model_type: str):
    """
    Get latest AI model by type
    """
    try:
        model = await model_service.get_model_by_type(model_type)
        if not model:
            raise HTTPException(status_code=404, detail=f"No model found for type: {model_type}")
        return model
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error getting model by type: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error getting model by type: {str(e)}")

@router.post("/upload/{model_type}")
async def upload_model(
    model_type: str,
    model_name: str = Form(...),
    model_version: str = Form(...),
    model_file: UploadFile = File(...),
):
    """
    Upload a new AI model
    """
    try:
        model_id = await model_service.upload_model(
            model_type=model_type,
            model_name=model_name,
            model_version=model_version,
            model_file=model_file
        )
        return {"model_id": model_id, "message": "Model uploaded successfully"}
    except Exception as e:
        logger.error(f"Error uploading model: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error uploading model: {str(e)}")
