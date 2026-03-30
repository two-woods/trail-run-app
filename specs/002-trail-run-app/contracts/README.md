# Interface Contracts: 越野跑App MVP

**Date**: 2026-03-30

---

## 1. REST API Contract

**Base URL**: `/api/v1`
**Authentication**: Bearer JWT Token

### Response Envelope

```json
{
  "success": true,
  "data": { ... },
  "error": null
}
```

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message"
  }
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| UNAUTHORIZED | 401 | Invalid or missing JWT |
| FORBIDDEN | 403 | No permission |
| NOT_FOUND | 404 | Resource not found |
| VALIDATION_ERROR | 400 | Invalid input |
| INTERNAL_ERROR | 500 | Server error |

---

## 2. API Endpoints

### 2.1 Authentication

#### POST /auth/register

**Request**:
```json
{
  "email": "user@example.com",
  "password": "securePassword123",
  "nickname": "越野达人"
}
```

**Response** (201):
```json
{
  "success": true,
  "data": {
    "user_id": "uuid",
    "token": "jwt_token"
  }
}
```

#### POST /auth/login

**Request**:
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Response** (200):
```json
{
  "success": true,
  "data": {
    "user_id": "uuid",
    "nickname": "越野达人",
    "avatar_url": "https://...",
    "token": "jwt_token"
  }
}
```

---

### 2.2 User

#### GET /users/me

**Headers**: `Authorization: Bearer <token>`

**Response** (200):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "nickname": "越野达人",
    "avatar_url": "https://...",
    "phone": "13800138000",
    "itra_account": null,
    "created_at": "2026-03-30T10:00:00Z"
  }
}
```

---

### 2.3 Races

#### GET /races

**Query Params**:
- `province` (optional): 省份筛选
- `city` (optional): 城市筛选
- `date_from` (optional): 开始日期
- `date_to` (optional): 结束日期
- `difficulty` (optional): Easy/Medium/Hard/Extrem
- `page` (default: 1): 页码
- `page_size` (default: 20): 每页数量

**Response** (200):
```json
{
  "success": true,
  "data": {
    "races": [
      {
        "id": "uuid",
        "name": "2026柴古唐斯100",
        "date": "2026-04-12",
        "location": "浙江临海",
        "distance_km": 102.5,
        "elevation_m": 5200,
        "difficulty": "Hard",
        "itra_points": 6,
        "status": "published"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

#### GET /races/:id

**Response** (200):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "2026柴古唐斯100",
    "date": "2026-04-12",
    "location": "浙江临海",
    "province": "浙江",
    "city": "临海市",
    "distance_km": 102.5,
    "elevation_m": 5200,
    "difficulty": "Hard",
    "itra_points": 6,
    "start_lat": 28.6562,
    "start_lng": 121.1234,
    "end_lat": 28.6562,
    "end_lng": 121.1234,
    "route_gpx_url": "https://cdn.example.com/race.gpx",
    "weather_city_code": "331001",
    "aid_stations": [
      {
        "id": "uuid",
        "name": "CP1-瓦窑",
        "distance_km": 15.5,
        "elevation_m": 890,
        "lat": 28.6821,
        "lng": 121.1567,
        "supplies": ["水", "运动饮料", "能量胶"],
        "close_time": "08:30:00"
      }
    ],
    "equipment": [
      {
        "id": "uuid",
        "name": "号码簿",
        "is_mandatory": true,
        "category": "穿着"
      }
    ]
  }
}
```

---

### 2.4 Race Day Plan

#### GET /races/:id/plan

**Response** (200):
```json
{
  "success": true,
  "data": {
    "race": {
      "id": "uuid",
      "name": "2026柴古唐斯100",
      "date": "2026-04-12",
      "distance_km": 102.5,
      "elevation_m": 5200
    },
    "commute": {
      "start_point": {
        "lat": 28.6562,
        "lng": 121.1234,
        "name": "起点"
      },
      "parking": [
        {
          "name": "市民广场停车场",
          "distance_m": 300,
          "lat": 28.6550,
          "lng": 121.1240
        }
      ]
    },
    "aid_stations": [...],
    "weather": {
      "date": "2026-04-12",
      "weather": "多云",
      "temperature": "15-22",
      "wind": "东南风3级",
      "humidity": "65%",
      "rain_probability": "30%"
    }
  }
}
```

#### GET /races/:id/equipment

**Headers**: `Authorization: Bearer <token>`

**Response** (200):
```json
{
  "success": true,
  "data": {
    "mandatory": [
      {
        "id": "uuid",
        "name": "号码簿",
        "category": "穿着",
        "checked": false
      }
    ],
    "recommended": [
      {
        "id": "uuid",
        "name": "登山杖",
        "category": "装备",
        "checked": true
      }
    ]
  }
}
```

#### POST /races/:id/equipment/check

**Headers**: `Authorization: Bearer <token>`

**Request**:
```json
{
  "equipment_id": "uuid",
  "checked": true
}
```

**Response** (200):
```json
{
  "success": true,
  "data": {
    "equipment_id": "uuid",
    "checked": true
  }
}
```

---

### 2.5 Results

#### GET /results

**Headers**: `Authorization: Bearer <token>`

**Response** (200):
```json
{
  "success": true,
  "data": {
    "results": [
      {
        "id": "uuid",
        "race_name": "2026柴古唐斯100",
        "race_date": "2026-04-12",
        "distance_km": 102.5,
        "elevation_m": 5200,
        "finish_time": "22:45:30",
        "ranking": 156,
        "ranking_age_group": 23,
        "status": "finished",
        "generated_image_url": "https://cdn.example.com/images/xxx.png"
      }
    ]
  }
}
```

#### POST /results

**Headers**: `Authorization: Bearer <token>`

**Request**:
```json
{
  "race_id": "uuid",
  "finish_time": "22:45:30",
  "ranking": 156,
  "ranking_age_group": 23,
  "gpx_url": "https://cdn.example.com/uploaded.gpx",
  "photos": ["https://cdn.example.com/photo1.jpg"]
}
```

**Response** (201):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "finished"
  }
}
```

#### GET /results/:id

**Response** (200):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "race": {
      "id": "uuid",
      "name": "2026柴古唐斯100",
      "date": "2026-04-12"
    },
    "finish_time": "22:45:30",
    "ranking": 156,
    "ranking_age_group": 23,
    "gpx_url": "https://...",
    "geojson": {
      "type": "LineString",
      "coordinates": [[121.1234, 28.6562], ...]
    },
    "photos": ["https://..."],
    "generated_image_url": "https://...",
    "created_at": "2026-04-12T22:45:30Z"
  }
}
```

#### POST /results/:id/image

**Headers**: `Authorization: Bearer <token>`

**Request**:
```json
{
  "race_name": "2026柴古唐斯100",
  "date": "2026-04-12",
  "gpx_data": "base64_encoded_gpx_xml",
  "stats": {
    "distance_km": 102.5,
    "elevation_m": 5200,
    "finish_time": "22:45:30"
  }
}
```

**Response** (202):
```json
{
  "success": true,
  "data": {
    "image_url": "https://cdn.example.com/results/xxx.png",
    "status": "processing"
  }
}
```

---

## 3. Python Image Service Contract

**Base URL**: `http://python-image-service:5000`

### POST /generate

**Request**:
```json
{
  "race_name": "2026柴古唐斯100",
  "date": "2026-04-12",
  "gpx_data": "base64_encoded_gpx_content",
  "distance_km": 102.5,
  "elevation_m": 5200,
  "finish_time": "22:45:30",
  "user_nickname": "越野达人"
}
```

**Response** (200):
```json
{
  "image_url": "https://oss.example.com/trail-images/xxx.png",
  "image_data": null
}
```

**Response** (500):
```json
{
  "error": "GPX parse error"
}
```
