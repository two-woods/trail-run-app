<template>
  <view class="container">
    <!-- Map Section -->
    <view class="map-section">
      <map
        id="resultMap"
        class="map"
        :latitude="mapCenter.lat"
        :longitude="mapCenter.lng"
        :scale="mapScale"
        :polyline="polyline"
        :markers="markers"
        :show-location="true"
        @tap="onMapTap"
        @zoomchange="onZoomChange"
        @end="onMapEnd"
      ></map>
      <view class="map-controls">
        <button class="map-btn" @click="resetMapView">重置视图</button>
      </view>
    </view>

    <!-- Stats Section -->
    <view class="stats-section" v-if="stats">
      <view class="stats-grid">
        <view class="stat-item">
          <text class="stat-value">{{ stats.distance_km?.toFixed(2) || '-' }}</text>
          <text class="stat-label">公里</text>
        </view>
        <view class="stat-item">
          <text class="stat-value">{{ stats.elevation_gain_m || '-' }}</text>
          <text class="stat-label">爬升(m)</text>
        </view>
        <view class="stat-item">
          <text class="stat-value">{{ stats.elevation_loss_m || '-' }}</text>
          <text class="stat-label">下降(m)</text>
        </view>
        <view class="stat-item">
          <text class="stat-value">{{ formatDuration(stats.duration_secs) }}</text>
          <text class="stat-label">用时</text>
        </view>
      </view>
      <view class="stats-secondary">
        <view class="stat-row">
          <text class="stat-label">最高海拔: {{ stats.max_elevation_m || '-' }}m</text>
          <text class="stat-label">最低海拔: {{ stats.min_elevation_m || '-' }}m</text>
        </view>
        <view class="stat-row" v-if="stats.avg_speed_kmh">
          <text class="stat-label">平均速度: {{ stats.avg_speed_kmh?.toFixed(2) }} km/h</text>
        </view>
      </view>
    </view>

    <!-- Race Info Section -->
    <view class="info-section" v-if="result">
      <view class="info-header">
        <text class="race-name">{{ result.race?.name || '比赛' }}</text>
        <view class="status-tag" :class="result.status">
          {{ statusText(result.status) }}
        </view>
      </view>

      <view class="info-grid">
        <view class="info-item">
          <text class="info-label">比赛日期</text>
          <text class="info-value">{{ formatDate(result.race?.date) }}</text>
        </view>
        <view class="info-item">
          <text class="info-label">比赛距离</text>
          <text class="info-value">{{ result.race?.distance_km || '-' }} km</text>
        </view>
        <view class="info-item">
          <text class="info-label">爬升</text>
          <text class="info-value">{{ result.race?.elevation_m || '-' }} m</text>
        </view>
        <view class="info-item">
          <text class="info-label">关门时间</text>
          <text class="info-value">{{ result.race?.cutoff_time || '-' }}</text>
        </view>
      </view>
    </view>

    <!-- Result Details -->
    <view class="result-section" v-if="result">
      <view class="section-title">完赛成绩</view>
      <view class="result-grid">
        <view class="result-item" v-if="result.finish_time">
          <text class="result-label">完赛时间</text>
          <text class="result-value">{{ result.finish_time }}</text>
        </view>
        <view class="result-item" v-if="result.ranking">
          <text class="result-label">总排名</text>
          <text class="result-value">#{{ result.ranking }}</text>
        </view>
        <view class="result-item" v-if="result.ranking_age_group">
          <text class="result-label">年龄组排名</text>
          <text class="result-value">#{{ result.ranking_age_group }}</text>
        </view>
      </view>
    </view>

    <!-- Photos Section -->
    <view class="photos-section" v-if="result?.photos?.length">
      <view class="section-title">照片</view>
      <scroll-view class="photos-scroll" scroll-x="true">
        <image
          v-for="(photo, idx) in result.photos"
          :key="idx"
          :src="photo"
          class="photo-item"
          mode="aspectFill"
          @click="previewPhoto(photo)"
        />
      </scroll-view>
    </view>

    <!-- Generated Image -->
    <view class="generated-section" v-if="result">
      <view class="section-title">完赛证书</view>

      <!-- Generated image display -->
      <view v-if="result.generated_image_url" class="certificate-container">
        <image
          :src="result.generated_image_url"
          class="generated-image"
          mode="widthFix"
          @click="previewPhoto(result.generated_image_url)"
        />
        <button class="btn-save" @click="saveImage">保存到相册</button>
      </view>

      <!-- Generate button -->
      <view v-else class="generate-prompt">
        <text class="prompt-text">生成您的专属完赛证书</text>
        <button class="btn-generate" @click="generateCertificate" :disabled="generating">
          {{ generating ? '生成中...' : '生成分享图' }}
        </button>
      </view>
    </view>

    <!-- Loading State -->
    <view v-if="loading" class="loading-state">
      <text>加载中...</text>
    </view>
  </view>
</template>

<script setup>
import { ref, onLoad } from 'vue'
import { api } from '@/services/api'

const result = ref(null)
const stats = ref(null)
const loading = ref(true)
const generating = ref(false)

// Map data
const mapCenter = ref({ lat: 30.5728, lng: 114.2525 }) // Default: Wuhan
const mapScale = ref(12)
const polyline = ref([])
const markers = ref([])

onLoad((query) => {
  if (query.id) {
    loadResult(query.id)
  }
})

async function loadResult(id) {
  try {
    loading.value = true
    const res = await api.getResult(id)
    result.value = res.data

    if (res.data.stats) {
      stats.value = res.data.stats
    }

    // Setup map if GeoJSON available
    if (res.data.geojson) {
      setupMap(res.data.geojson, res.data.race)
    }
  } catch (e) {
    console.error('Failed to load result:', e)
    uni.showToast({ title: '加载失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

function setupMap(geojson, race) {
  const features = geojson.features || []

  // Process LineString track
  const trackFeature = features.find(f =>
    f.geometry?.type === 'LineString' && !f.properties?.marker
  )

  if (trackFeature?.geometry?.coordinates?.length) {
    const coords = trackFeature.geometry.coordinates
    const points = coords.map(c => ({ longitude: c[0], latitude: c[1] }))

    polyline.value = [{
      points: points,
      color: '#667eea',
      width: 4,
      dottedLine: false
    }]

    // Fit map to track bounds
    if (points.length > 0) {
      const lats = points.map(p => p.latitude)
      const lngs = points.map(p => p.longitude)
      const minLat = Math.min(...lats)
      const maxLat = Math.max(...lats)
      const minLng = Math.min(...lngs)
      const maxLng = Math.max(...lngs)

      mapCenter.value = {
        lat: (minLat + maxLat) / 2,
        lng: (minLng + maxLng) / 2
      }

      // Calculate scale based on bounds
      const latDiff = maxLat - minLat
      const lngDiff = maxLng - minLng
      const maxDiff = Math.max(latDiff, lngDiff)
      if (maxDiff > 0.1) mapScale.value = 10
      else if (maxDiff > 0.05) mapScale.value = 11
      else if (maxDiff > 0.01) mapScale.value = 12
      else mapScale.value = 13
    }
  }

  // Process markers (start/end points)
  const startFeature = features.find(f => f.properties?.marker === 'start')
  const endFeature = features.find(f => f.properties?.marker === 'end')

  if (startFeature?.geometry?.coordinates) {
    markers.value.push({
      id: 'start',
      longitude: startFeature.geometry.coordinates[0],
      latitude: startFeature.geometry.coordinates[1],
      width: 24,
      height: 24,
      iconPath: '/static/marker-start.png',
      title: '起点'
    })
  }

  if (endFeature?.geometry?.coordinates) {
    markers.value.push({
      id: 'end',
      longitude: endFeature.geometry.coordinates[0],
      latitude: endFeature.geometry.coordinates[1],
      width: 24,
      height: 24,
      iconPath: '/static/marker-end.png',
      title: '终点'
    })
  }

  // Add aid station markers if available from race
  if (race?.aid_stations?.length) {
    race.aid_stations.forEach((station, idx) => {
      if (station.latitude && station.longitude) {
        markers.value.push({
          id: `aid-${idx}`,
          longitude: station.longitude,
          latitude: station.latitude,
          width: 20,
          height: 20,
          iconPath: '/static/marker-aid.png',
          title: station.name,
          callout: {
            content: station.name,
            color: '#333',
            fontSize: 12,
            borderRadius: 4,
            padding: 6,
            display: 'BYCLICK'
          }
        })
      }
    })
  }
}

function resetMapView() {
  if (result.value?.geojson) {
    setupMap(result.value.geojson, result.value.race)
  }
}

function onMapTap(e) {
  console.log('Map tapped:', e.detail)
}

function onZoomChange(e) {
  mapScale.value = e.detail.scale
}

function onMapEnd(e) {
  console.log('Map gesture end:', e.detail)
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日`
}

function formatDuration(secs) {
  if (!secs) return '-'
  const h = Math.floor(secs / 3600)
  const m = Math.floor((secs % 3600) / 60)
  const s = secs % 60
  if (h > 0) {
    return `${h}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
  }
  return `${m}:${s.toString().padStart(2, '0')}`
}

function statusText(status) {
  const map = {
    registered: '已报名',
    ongoing: '进行中',
    finished: '已完赛'
  }
  return map[status] || status
}

function previewPhoto(url) {
  uni.previewImage({
    urls: [url],
    current: url
  })
}

async function generateCertificate() {
  if (!result.value?.race) {
    uni.showToast({ title: '缺少比赛数据', icon: 'none' })
    return
  }

  try {
    generating.value = true
    const race = result.value.race

    // Fetch GPX data if we have a GPX URL
    let gpxData = ''
    if (result.value.gpx_url) {
      try {
        const gpxResp = await uni.request({
          url: result.value.gpx_url,
          method: 'GET'
        })
        gpxData = base64Encode(gpxResp.data)
      } catch (e) {
        console.error('Failed to fetch GPX:', e)
      }
    }

    const res = await api.generateImage(result.value.id, {
      race_name: race.name,
      date: formatDate(race.date),
      gpx_data: gpxData,
      distance_km: stats.value?.distance_km || race.distance_km,
      elevation_m: stats.value?.elevation_gain_m || race.elevation_m,
      finish_time: result.value.finish_time || '00:00:00'
    })

    if (res.data.image_data) {
      result.value.generated_image_url = `data:image/jpeg;base64,${res.data.image_data}`
      uni.showToast({ title: '生成成功', icon: 'success' })
    }
  } catch (e) {
    console.error('Failed to generate image:', e)
    uni.showToast({ title: '生成失败', icon: 'none' })
  } finally {
    generating.value = false
  }
}

function base64Encode(str) {
  // Simple base64 encoding for Latin-1 strings
  return btoa(unescape(encodeURIComponent(str)))
}

async function saveImage() {
  if (!result.value?.generated_image_url) return

  try {
    // For base64 data URLs, save directly
    const base64Data = result.value.generated_image_url.split(',')[1]
    const filePath = `${wx.env.USER_DATA_PATH}/certificate_${result.value.id}.jpg`

    const fs = uni.getFileSystemManager()
    fs.writeFile({
      filePath: filePath,
      data: base64Data,
      encoding: 'base64',
      success: () => {
        uni.saveImageToPhotosAlbum({
          filePath: filePath,
          success: () => {
            uni.showToast({ title: '已保存到相册', icon: 'success' })
          },
          fail: () => {
            uni.showToast({ title: '保存失败', icon: 'none' })
          }
        })
      },
      fail: () => {
        uni.showToast({ title: '保存失败', icon: 'none' })
      }
    })
  } catch (e) {
    console.error('Failed to save image:', e)
    uni.showToast({ title: '保存失败', icon: 'none' })
  }
}
</script>

<style scoped>
.container {
  min-height: 100vh;
  background: #f5f5f5;
}

.map-section {
  position: relative;
  height: 400rpx;
  background: #ddd;
}

.map {
  width: 100%;
  height: 100%;
}

.map-controls {
  position: absolute;
  right: 24rpx;
  bottom: 24rpx;
}

.map-btn {
  padding: 12rpx 24rpx;
  background: #fff;
  border-radius: 8rpx;
  font-size: 24rpx;
  box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.15);
}

.stats-section {
  background: #fff;
  padding: 24rpx;
  margin-bottom: 16rpx;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16rpx;
  margin-bottom: 16rpx;
}

.stat-item {
  text-align: center;
}

.stat-value {
  display: block;
  font-size: 36rpx;
  font-weight: bold;
  color: #667eea;
}

.stat-label {
  font-size: 22rpx;
  color: #999;
}

.stats-secondary {
  border-top: 1rpx solid #eee;
  padding-top: 16rpx;
}

.stat-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8rpx;
}

.info-section,
.result-section,
.photos-section,
.generated-section {
  background: #fff;
  padding: 24rpx;
  margin-bottom: 16rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 16rpx;
}

.info-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
}

.race-name {
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
}

.status-tag {
  padding: 6rpx 16rpx;
  border-radius: 8rpx;
  font-size: 22rpx;
}

.status-tag.registered {
  background: #e3f2fd;
  color: #2196f3;
}

.status-tag.ongoing {
  background: #fff3e0;
  color: #ff9800;
}

.status-tag.finished {
  background: #e8f5e9;
  color: #4caf50;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24rpx;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.info-label {
  font-size: 24rpx;
  color: #999;
}

.info-value {
  font-size: 28rpx;
  color: #333;
}

.result-grid {
  display: flex;
  gap: 48rpx;
}

.result-item {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.result-label {
  font-size: 24rpx;
  color: #999;
}

.result-value {
  font-size: 32rpx;
  font-weight: bold;
  color: #667eea;
}

.photos-scroll {
  white-space: nowrap;
}

.photo-item {
  width: 200rpx;
  height: 200rpx;
  border-radius: 8rpx;
  margin-right: 12rpx;
  display: inline-block;
}

.generated-image {
  width: 100%;
  border-radius: 12rpx;
}

.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 100rpx;
  color: #999;
}

.certificate-container {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.btn-save {
  padding: 16rpx 32rpx;
  background: #fff;
  border: 2rpx solid #667eea;
  color: #667eea;
  border-radius: 32rpx;
  font-size: 28rpx;
}

.generate-prompt {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 48rpx 0;
  gap: 24rpx;
}

.prompt-text {
  font-size: 28rpx;
  color: #666;
}

.btn-generate {
  padding: 20rpx 64rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  border-radius: 32rpx;
  font-size: 32rpx;
  border: none;
}

.btn-generate[disabled] {
  background: #ccc;
}
</style>
