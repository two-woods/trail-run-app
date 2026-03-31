from fastapi import FastAPI, HTTPException
from fastapi.responses import StreamingResponse
from pydantic import BaseModel
from typing import Optional
import base64
import io
import math
import xml.etree.ElementTree as ET
from datetime import datetime

from PIL import Image, ImageDraw, ImageFont, ImageFilter
import matplotlib.pyplot as plt
import matplotlib.patches as mpatches
from matplotlib.patches import FancyBboxPatch
import numpy as np

app = FastAPI(title="Trail Run Image Service")


class ImageGenerationRequest(BaseModel):
    race_name: str
    date: str
    gpx_data: str  # Base64 encoded GPX content
    distance_km: float
    elevation_m: int
    finish_time: str
    user_nickname: Optional[str] = None


class ImageServiceResponse(BaseModel):
    image_url: str = ""
    image_data: Optional[str] = None  # Base64 for small images


def parse_gpx(gpx_content: str) -> tuple[list, list, float, int, int]:
    """
    Parse GPX content and extract:
    - latitudes, longitudes
    - elevations
    Returns: (lats, lons, total_distance_km, total_elevation_gain, total_elevation_loss)
    """
    try:
        root = ET.fromstring(gpx_content)
    except Exception:
        raise ValueError("Invalid GPX format")

    ns = {'gpx': 'http://www.topografix.com/GPX/1/1'}
    lats, lons, elevations = [], [], []

    for trk in root.findall('.//gpx:trk', ns):
        for seg in trk.findall('gpx:trkseg', ns):
            for pt in seg.findall('gpx:trkpt', ns):
                lat = float(pt.get('lat'))
                lon = float(pt.get('lon'))
                ele_elem = pt.find('gpx:ele', ns)
                ele = float(ele_elem.text) if ele_elem is not None else 0

                lats.append(lat)
                lons.append(lon)
                elevations.append(ele)

    if len(lats) < 2:
        return [0], [0], [0], 0, 0, 0

    # Calculate distance using haversine
    total_dist = 0
    for i in range(1, len(lats)):
        total_dist += haversine(lats[i-1], lons[i-1], lats[i], lons[i])

    # Calculate elevation gain/loss
    elev_gain, elev_loss = 0, 0
    for i in range(1, len(elevations)):
        diff = elevations[i] - elevations[i-1]
        if diff > 0:
            elev_gain += diff
        else:
            elev_loss += abs(diff)

    return lats, lons, total_dist, int(elev_gain), int(elev_loss)


def haversine(lat1, lon1, lat2, lon2):
    """Calculate distance between two points in km"""
    R = 6371.0
    lat1_rad = math.radians(lat1)
    lat2_rad = math.radians(lat2)
    delta_lat = math.radians(lat2 - lat1)
    delta_lon = math.radians(lon2 - lon1)

    a = math.sin(delta_lat/2)**2 + math.cos(lat1_rad) * math.cos(lat2_rad) * math.sin(delta_lon/2)**2
    c = 2 * math.atan2(math.sqrt(a), math.sqrt(1-a))

    return R * c


def create_gradient_color(value, min_val, max_val):
    """Create gradient color from blue to red based on value"""
    if max_val == min_val:
        return (0.4, 0.6, 1.0)
    ratio = (value - min_val) / (max_val - min_val)
    # Blue -> Cyan -> Green -> Yellow -> Red
    if ratio < 0.25:
        return (0.4, 0.4 + ratio*4*0.2, 1.0)
    elif ratio < 0.5:
        return (0.4, 0.6 - (ratio-0.25)*4*0.2, 1.0 - (ratio-0.25)*4*0.4)
    elif ratio < 0.75:
        return (0.4 + (ratio-0.5)*4*0.6, 0.4 - (ratio-0.5)*4*0.4, 0.6 - (ratio-0.5)*4*0.6)
    else:
        return (1.0, 0.0, 0.0)


def generate_vertical_gradient(height, width, color1, color2):
    """Generate vertical gradient image"""
    gradient = np.zeros((height, width, 3), dtype=np.uint8)
    for y in range(height):
        ratio = y / height
        r = int(color1[0] * (1 - ratio) + color2[0] * ratio)
        g = int(color1[1] * (1 - ratio) + color2[1] * ratio)
        b = int(color1[2] * (1 - ratio) + color2[2] * ratio)
        gradient[y, :] = [r, g, b]
    return Image.fromarray(gradient)


async def generate_image(req: ImageGenerationRequest) -> bytes:
    """
    Generate trail race image with:
    - Race name and date
    - Trail map with GPX track (colored by elevation)
    - Key stats (distance, elevation, finish time)
    """
    # Parse GPX
    try:
        gpx_bytes = base64.b64decode(req.gpx_data)
        gpx_content = gpx_bytes.decode('utf-8')
    except Exception:
        gpx_content = req.gpx_data

    try:
        lats, lons, distance, elev_gain, elev_loss = parse_gpx(gpx_content)
    except Exception:
        lats, lons = [30.5], [114.3]  # fallback
        distance, elev_gain, elev_loss = 0, 0, 0

    # Image dimensions
    img_width = 800
    img_height = 1000
    map_height = 500

    # Create base gradient background
    bg = generate_vertical_gradient(img_height, img_width, (25, 35, 70), (15, 20, 40))
    draw = ImageDraw.Draw(bg)

    # Try to load fonts
    try:
        title_font = ImageFont.truetype("/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf", 42)
        subtitle_font = ImageFont.truetype("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 28)
        stat_font = ImageFont.truetype("/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf", 32)
        small_font = ImageFont.truetype("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 20)
        tiny_font = ImageFont.truetype("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 16)
    except:
        title_font = subtitle_font = stat_font = small_font = tiny_font = ImageFont.load_default()

    # === Title Section ===
    title_y = 60
    draw.text((img_width//2, title_y), "完赛证书", font=title_font, fill=(255, 255, 255), anchor="mm")
    draw.text((img_width//2, title_y + 55), "FINISH CERTIFICATE", font=subtitle_font, fill=(180, 180, 200), anchor="mm")

    # === Race Name ===
    draw.text((img_width//2, 160), req.race_name, font=stat_font, fill=(255, 255, 255), anchor="mm")

    # === Date ===
    draw.text((img_width//2, 200), req.date, font=small_font, fill=(150, 150, 180), anchor="mm")

    # === Stats Card ===
    card_padding = 20
    card_height = 120
    card_y = 240
    card = Image.new('RGBA', (img_width - 80, card_height), (255, 255, 255, 30))
    card_draw = ImageDraw.Draw(card)
    card_draw.rounded_rectangle([0, 0, img_width - 80, card_height], radius=15, fill=(255, 255, 255, 20))
    bg.paste(card, (40, card_y), card)

    # Stats
    stats = [
        (f"{req.distance_km:.2f}", "公里"),
        (f"{req.elevation_m}", "爬升m"),
        (req.finish_time, "用时"),
    ]

    stat_width = (img_width - 80) // 3
    for i, (value, label) in enumerate(stats):
        x = 40 + stat_width * i + stat_width // 2
        y = card_y + 35
        draw.text((x, y), value, font=stat_font, fill=(255, 255, 255), anchor="mm")
        draw.text((x, y + 40), label, font=tiny_font, fill=(150, 150, 180), anchor="mm")

    # === Map Section ===
    map_y = 380
    map_margin = 40
    map_width = img_width - map_margin * 2

    # Map background
    map_bg = Image.new('RGB', (map_width, map_height), (40, 50, 70))
    map_draw = ImageDraw.Draw(map_bg)

    # Draw track on map
    if len(lats) > 1 and distance > 0:
        # Normalize coordinates to map bounds
        min_lat, max_lat = min(lats), max(lats)
        min_lon, max_lon = min(lons), max(lons)

        # Add padding
        lat_pad = (max_lat - min_lat) * 0.1 or 0.01
        lon_pad = (max_lon - min_lon) * 0.1 or 0.01
        min_lat -= lat_pad
        max_lat += lat_pad
        min_lon -= lon_pad
        max_lon += lon_pad

        # Get elevations for coloring
        elevations = []
        try:
            root = ET.fromstring(gpx_content)
            ns = {'gpx': 'http://www.topografix.com/GPX/1/1'}
            for trk in root.findall('.//gpx:trk', ns):
                for seg in trk.findall('gpx:trkseg', ns):
                    for pt in seg.findall('gpx:trkpt', ns):
                        ele_elem = pt.find('gpx:ele', ns)
                        elevations.append(float(ele_elem.text) if ele_elem is not None else 0)
        except:
            elevations = [0] * len(lats)

        min_ele = min(elevations) if elevations else 0
        max_ele = max(elevations) if elevations else 1

        # Draw track with gradient coloring
        point_size = 3
        for i in range(len(lats) - 1):
            x1 = int((lons[i] - min_lon) / (max_lon - min_lon) * (map_width - 20) + 10)
            y1 = int((max_lat - lats[i]) / (max_lat - min_lat) * (map_height - 20) + 10)
            x2 = int((lons[i+1] - min_lon) / (max_lon - min_lon) * (map_width - 20) + 10)
            y2 = int((max_lat - lats[i+1]) / (max_lat - min_lat) * (map_height - 20) + 10)

            # Get color based on elevation
            if i < len(elevations):
                color = create_gradient_color(elevations[i], min_ele, max_ele)
                color_int = (int(color[0]*255), int(color[1]*255), int(color[2]*255))
            else:
                color_int = (100, 200, 255)

            map_draw.line([(x1, y1), (x2, y2)], fill=color_int, width=3)

        # Draw start marker (green)
        start_x = int((lons[0] - min_lon) / (max_lon - min_lon) * (map_width - 20) + 10)
        start_y = int((max_lat - lats[0]) / (max_lat - min_lat) * (map_height - 20) + 10)
        map_draw.ellipse([start_x-6, start_y-6, start_x+6, start_y+6], fill=(50, 200, 100))

        # Draw end marker (red)
        end_x = int((lons[-1] - min_lon) / (max_lon - min_lon) * (map_width - 20) + 10)
        end_y = int((max_lat - lats[-1]) / (max_lat - min_lat) * (map_height - 20) + 10)
        map_draw.ellipse([end_x-6, end_y-6, end_x+6, end_y+6], fill=(255, 80, 80))

    # Add subtle grid
    for i in range(5):
        x = map_width * i // 4
        map_draw.line([(x, 0), (x, map_height)], fill=(60, 70, 100), width=1)
    for i in range(5):
        y = map_height * i // 4
        map_draw.line([(0, y), (map_width, y)], fill=(60, 70, 100), width=1)

    # Apply blur for depth effect
    map_bg = map_bg.filter(ImageFilter.GaussianBlur(radius=2))

    # Composite map onto background
    bg.paste(map_bg, (map_margin, map_y))

    # Map border
    map_draw = ImageDraw.Draw(bg)
    map_draw.rectangle([map_margin, map_y, map_margin + map_width, map_y + map_height],
                       outline=(100, 120, 180), width=2)

    # === Elevation Profile ===
    profile_y = map_y + map_height + 20
    profile_height = 100
    profile_width = map_width

    profile_bg = Image.new('RGB', (profile_width, profile_height), (30, 35, 55))
    profile_draw = ImageDraw.Draw(profile_bg)

    if len(lats) > 1 and elevations:
        min_ele = min(elevations)
        max_ele = max(elevations)
        ele_range = max_ele - min_ele or 1

        # Build elevation polyline
        polyline_points = []
        for i, ele in enumerate(elevations):
            x = int(i / (len(elevations) - 1) * (profile_width - 20)) + 10
            y = int((max_ele - ele) / ele_range * (profile_height - 20)) + 10
            polyline_points.append((x, y))

        # Draw filled area
        fill_points = polyline_points + [(profile_width - 10, profile_height - 10), (10, profile_height - 10)]
        profile_draw.polygon(fill_points, fill=(80, 100, 180, 100))

        # Draw line
        for i in range(len(polyline_points) - 1):
            profile_draw.line([polyline_points[i], polyline_points[i+1]], fill=(100, 180, 255), width=2)

        # Elevation labels
        profile_draw.text((15, 10), f"{int(max_ele)}m", font=tiny_font, fill=(150, 150, 180))
        profile_draw.text((15, profile_height - 25), f"{int(min_ele)}m", font=tiny_font, fill=(150, 150, 180))

    bg.paste(profile_bg, (map_margin, profile_y))

    # === Footer ===
    footer_y = profile_y + profile_height + 20
    if req.user_nickname:
        draw.text((img_width//2, footer_y), f"@{req.user_nickname}", font=small_font, fill=(120, 120, 160), anchor="mm")

    # Add app branding
    draw.text((img_width//2, img_height - 30), "越野跑 App", font=tiny_font, fill=(80, 80, 100), anchor="mm")

    # Convert to bytes
    img_bytes = io.BytesIO()
    bg.convert('RGB').save(img_bytes, format='JPEG', quality=95)
    img_bytes.seek(0)
    return img_bytes.getvalue()


@app.post("/generate", response_model=ImageServiceResponse)
async def generate_image_endpoint(req: ImageGenerationRequest):
    """
    Generate trail race image with:
    - Race name and date
    - Trail map with GPX track
    - Key stats (distance, elevation, finish time)
    """
    try:
        image_bytes = await generate_image(req)
        image_b64 = base64.b64encode(image_bytes).decode('utf-8')

        return ImageServiceResponse(
            image_url="",
            image_data=image_b64
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/health")
async def health_check():
    return {"status": "ok"}
