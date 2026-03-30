from fastapi import FastAPI
from pydantic import BaseModel
from typing import Optional
import base64

app = FastAPI(title="Trail Run Image Service")


class ImageGenerationRequest(BaseModel):
    race_name: str
    date: str
    gpx_data: str  # Base64 encoded GPX
    distance_km: float
    elevation_m: int
    finish_time: str
    user_nickname: Optional[str] = None


class ImageServiceResponse(BaseModel):
    image_url: str = ""
    image_data: Optional[str] = None  # Base64 for small images


@app.post("/generate", response_model=ImageServiceResponse)
async def generate_image(req: ImageGenerationRequest):
    """
    Generate trail race image with:
    - Race name and date
    - Trail map with GPX track
    - Key stats (distance, elevation, finish time)
    """
    try:
        gpx_bytes = base64.b64decode(req.gpx_data)

        # TODO: Implement actual image generation
        # For now, return a placeholder response
        return ImageServiceResponse(
            image_url="",
            image_data=None
        )
    except Exception as e:
        return ImageServiceResponse(image_url="", image_data=None)


@app.get("/health")
async def health_check():
    return {"status": "ok"}
