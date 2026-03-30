# Research Report: 越野跑App MVP

**Date**: 2026-03-30
**Status**: Complete

---

## 1. Uni-app + Go Backend Integration

**Decision**: HTTP REST API with JWT authentication

### Architecture

```
Uni-app Client (uni.request)
    ↓ HTTPS
Go Backend (Gorilla Mux)
    ├── JWT Authentication
    └── RESTful Endpoints /api/v1/*
```

### Uni-app HTTP Client Pattern

```javascript
// api/request.js
const BASE_URL = 'https://api.example.com'

const request = (options) => {
  return new Promise((resolve, reject) => {
    uni.request({
      url: BASE_URL + options.url,
      method: options.method || 'GET',
      data: options.data || {},
      header: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + uni.getStorageSync('token'),
      },
      success: (res) => {
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve(res.data)
        } else if (res.statusCode === 401) {
          uni.reLaunch({ url: '/pages/login/login' })
          reject(new Error('Unauthorized'))
        } else {
          reject(res.data)
        }
      }
    })
  })
}
```

### Go JWT Middleware

```go
// gorilla mux with JWT auth
r := mux.NewRouter()
api := r.PathPrefix("/api/v1").Subrouter()

// Protected routes with auth middleware
protected := api.PathPrefix("").Subrouter()
protected.Use(authMiddleware.AuthRequired)
protected.HandleFunc("/users/me", handlers.GetProfile).Methods("GET")
```

**Alternatives Considered**:
- WebSocket: Not needed for MVP (REST sufficient)
- GraphQL: Overkill for this scope

---

## 2. 高德地图 (Amap) SDK Integration

**Decision**: 高德小程序SDK + Web服务API

### Key申请
- 应用类型: 微信小程序
- 需要: Web服务Key (HTTP API) + 小程序Key (map组件)

### Map Display Patterns

```xml
<map
  latitude="{{center.lat}}"
  longitude="{{center.lng}}"
  polyline="{{polyline}}"
  markers="{{markers}}"
  scale="14"
  style="width:100%;height:100vh;"
/>
```

### GPX轨迹渲染

```javascript
// GPX → 高德坐标格式
function parseGPX(gpxData) {
  const points = []
  const xml = new DOMParser().parseFromString(gpxData, 'text/xml')
  xml.querySelectorAll('trkpt').forEach(pt => {
    points.push({
      latitude: parseFloat(pt.getAttribute('lat')),
      longitude: parseFloat(pt.getAttribute('lon'))
    })
  })
  return points
}

// 绘制赛道路线
const polyline = [{
  points: parseGPX(gpxData),
  color: '#FF6B35',
  width: 4,
  borderWidth: 2
}]
```

### Markers (补给站/起点/终点)

```javascript
const markers = [{
  id: station.id,
  latitude: station.lat,
  longitude: station.lng,
  iconPath: '/assets/aid.png',
  callout: {
    content: `${station.name}\n${station.distance_km}km`,
    display: 'BYCLICK'
  }
}]
```

**Note**: 高德使用 GCJ-02 坐标系，与WGS-84不同，GPX文件需转换

**Alternatives Considered**:
- Google Maps: 国内不稳，放弃
- Mapbox: 国内支持不如高德

---

## 3. Go + Python 微服务通信

**Decision**: HTTP REST (JSON) + Separate Docker Containers

### Why HTTP REST over gRPC

| Factor | HTTP REST | gRPC |
|--------|-----------|------|
| 性能需求 | 30s SLA足够 | 过度设计 |
| 开发速度 | 快 | 需要写proto |
| 调试 | curl/browser | grpcurl |
| 团队熟悉度 | 高 | 中 |

### Go → Python 调用

```go
// internal/service/image.go
type ImageService struct {
    client   *http.Client
    baseURL  string
}

func (s *ImageService) Generate(ctx context.Context, req ImageRequest) (string, error) {
    payload, _ := json.Marshal(req)

    resp, err := s.client.Post(s.baseURL+"/generate", "application/json", bytes.NewReader(payload))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var result ImageResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return result.ImageURL, nil
}
```

### Python FastAPI Endpoint

```python
# image-service/app/main.py
from fastapi import FastAPI
app = FastAPI()

@app.post("/generate")
async def generate_image(req: ImageRequest):
    gpx_data = base64.b64decode(req.gpx_data)
    image = create_trail_image(gpx_data, req.race_name, req.stats)
    return {"image_url": upload_to_oss(image)}
```

### Docker Compose

```yaml
services:
  backend:
    build: ./backend
    environment:
      - IMAGE_SERVICE_URL=http://python-image-service:5000

  python-image-service:
    build: ./image-service
    ports:
      - "5000:5000"
```

**Alternatives Considered**:
- gRPC: 性能优势不明显，复杂度高
- Sidecar: Kubernetes才需要
- Process spawn: 无隔离，不生产级

---

## 4. GPX轨迹解析

**Decision**: `github.com/tkrajina/gpxgo` 库

### 安装

```bash
go get github.com/tkrajina/gpxgo/gpx
```

### 解析和统计

```go
import "github.com/tkrajina/gpxgo/gpx"

func parseGPX(gpxData []byte) (*gpx.GPX, error) {
    return gpx.Parse(gpxData)
}

func calculateStats(segments []gpx.GPXSegment) (distance, elevationGain float64) {
    for i := 0; i < len(segments[0].Points)-1; i++ {
        p1 := segments[0].Points[i]
        p2 := segments[0].Points[i+1]

        distance += haversine(p1.Latitude, p1.Longitude, p2.Latitude, p2.Longitude)
        if p2.Elevation.Value > p1.Elevation.Value {
            elevationGain += p2.Elevation.Value - p1.Elevation.Value
        }
    }
    return
}
```

### 转换为GeoJSON (for map display)

```go
type GeoJSON struct {
    Type        string        `json:"type"`
    Coordinates [][2]float64 `json:"coordinates"` // [lon, lat]
}

func toGeoJSON(segments []gpx.GPXSegment) GeoJSON {
    var coords [][2]float64
    for _, seg := range segments {
        for _, pt := range seg.Points {
            coords = append(coords, [2]float64{pt.Longitude, pt.Latitude})
        }
    }
    return GeoJSON{Type: "LineString", Coordinates: coords}
}
```

**Alternatives Considered**:
- 自解析GPX XML: 重复造轮子
- 其他库: gpxgo最成熟，Star最多

---

## 5. 天气API

**Decision**: 高德天气API (与地图同一服务商)

### API调用

```
GET https://restapi.amap.com/v3/weather/weatherInfo
? key=YOUR_KEY
& city=城市编码
& extensions=all (实时+预报)
```

### Uni-app 封装

```javascript
// services/weather.js
export async function getWeather(cityCode) {
  const res = await uni.request({
    url: 'https://restapi.amap.com/v3/weather/weatherInfo',
    data: {
      key: process.env.AMAP_KEY,
      city: cityCode,
      extensions: 'all'
    }
  })
  return res.lives[0] // 实时天气
}
```

### 返回数据

```json
{
  "province": "浙江",
  "city": "临海市",
  "weather": "晴",
  "temperature": "18",
  "wind": "北风",
  "windpower": "3",
  "humidity": "45"
}
```

**Alternatives Considered**:
- OpenWeatherMap: 国内不稳
- 和风天气: 需要单独注册

---

## 6. 离线缓存策略

**Decision**: Uni-app Storage + 赛事数据预缓存

### 分层缓存

```
┌─────────────────────────┐
│  Uni Storage            │  关键赛事数据
│  - 当前比赛信息           │  - 装备清单
│  - 路线数据              │  - 补给站信息
└─────────────────────────┘
```

### 实现方案

```javascript
// 比赛前自动缓存
async function cacheRaceForOffline(raceId) {
  const [race, plan, equipment] = await Promise.all([
    api.getRace(raceId),
    api.getRacePlan(raceId),
    api.getEquipment(raceId)
  ])

  uni.setStorage({
    key: `race_${raceId}`,
    data: { race, plan, equipment }
  })
}

// 离线时读取缓存
function getOfflineRaceData(raceId) {
  const cached = uni.getStorageSync(`race_${raceId}`)
  if (cached) {
    return cached.data
  }
  return null
}
```

**Offline触发**: 赛事日前一天自动缓存

---

## Summary

| Decision | Choice | Rationale |
|----------|--------|-----------|
| 前端框架 | Uni-app (Vue) | 一次开发多端 |
| 后端框架 | Go + Gorilla Mux | 高并发，成熟稳定 |
| 认证 | JWT | 标准方案 |
| 地图 | 高德SDK | 国内最优 |
| GPX解析 | gpxgo库 | 最成熟 |
| 图片服务 | Python FastAPI | 图像处理优势 |
| 通信协议 | HTTP REST | 简单够用 |
| 天气API | 高德天气 | 统一服务商 |
| 缓存 | Uni Storage | 端侧缓存 |

All NEEDS CLARIFICATION items resolved.
