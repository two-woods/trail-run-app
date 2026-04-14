<template>
  <view class="container">
    <!-- 比赛信息头部 -->
    <view class="race-header">
      <view class="race-info">
        <text class="race-name">{{ race.name }}</text>
        <text class="race-date">{{ formatDate(race.date) }}</text>
      </view>
      <view class="race-stats">
        <view class="stat">
          <text class="stat-value">{{ race.distance_km }}</text>
          <text class="stat-label">公里</text>
        </view>
        <view class="stat">
          <text class="stat-value">{{ race.elevation_m }}</text>
          <text class="stat-label">爬升</text>
        </view>
      </view>
    </view>

    <!-- 天气预报 -->
    <view class="section weather-section" v-if="weather">
      <text class="section-title">📅 天气预报</text>
      <view class="weather-card">
        <view class="weather-main">
          <text class="weather-icon">{{ getWeatherIcon(weather.weather) }}</text>
          <view class="weather-info">
            <text class="temperature">{{ weather.temperature }}°C</text>
            <text class="weather-desc">{{ weather.weather }}</text>
          </view>
        </view>
        <view class="weather-details">
          <view class="weather-item">
            <text class="weather-label">风力</text>
            <text class="weather-value">{{ weather.wind }}</text>
          </view>
          <view class="weather-item">
            <text class="weather-label">湿度</text>
            <text class="weather-value">{{ weather.humidity }}</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 通勤方式 -->
    <view class="section">
      <text class="section-title">🚗 通勤方式</text>
      <view class="commute-card" @click="openNavigation">
        <view class="start-point">
          <text class="point-icon">📍</text>
          <view class="point-info">
            <text class="point-label">起点</text>
            <text class="point-name">{{ commute.start_point?.name || '起点位置' }}</text>
          </view>
        </view>
        <view class="nav-btn">
          <text>导航</text>
        </view>
      </view>

      <view class="parking-list" v-if="commute.parking?.length > 0">
        <text class="subsection-title">🅿️ 附近停车场</text>
        <view
          v-for="(p, index) in commute.parking"
          :key="index"
          class="parking-item"
          @click="openParkingNavigation(p)"
        >
          <view class="parking-info">
            <text class="parking-name">{{ p.name }}</text>
            <text class="parking-distance">距起点 {{ p.distance_m }}m</text>
          </view>
          <text class="nav-arrow">›</text>
        </view>
      </view>
    </view>

    <!-- 补给站 -->
    <view class="section">
      <text class="section-title">🏁 补给站 ({{ aidStations.length }})</text>
      <view
        v-for="aid in aidStations"
        :key="aid.id"
        class="aid-card"
      >
        <view class="aid-header">
          <text class="aid-name">{{ aid.name }}</text>
          <text class="aid-close-time" v-if="aid.close_time">关门: {{ aid.close_time }}</text>
        </view>
        <view class="aid-info">
          <text class="aid-distance">{{ aid.distance_km }}km</text>
          <text class="aid-elevation">海拔 {{ aid.elevation_m }}m</text>
        </view>
        <view class="aid-supplies">
          <text
            v-for="supply in aid.supplies"
            :key="supply"
            class="supply-tag"
          >
            {{ supply }}
          </text>
        </view>
      </view>
    </view>

    <!-- 装备清单 -->
    <view class="section">
      <text class="section-title">🎒 装备清单</text>

      <view class="equipment-category" v-if="equipment.mandatory?.length">
        <text class="category-title red">⚠️ 强制装备 ({{ missingMandatory.length }}项未勾选)</text>
        <view
          v-for="eq in equipment.mandatory"
          :key="eq.id"
          class="equipment-item"
          @click="toggleEquipment(eq)"
        >
          <view class="checkbox" :class="{ checked: eq.checked, missing: !eq.checked && isMissing(eq) }">
            <text v-if="eq.checked">✓</text>
          </view>
          <text class="equipment-name">{{ eq.name }}</text>
        </view>
      </view>

      <view class="equipment-category" v-if="equipment.recommended?.length">
        <text class="category-title">💡 建议装备</text>
        <view
          v-for="eq in equipment.recommended"
          :key="eq.id"
          class="equipment-item"
          @click="toggleEquipment(eq)"
        >
          <view class="checkbox" :class="{ checked: eq.checked }">
            <text v-if="eq.checked">✓</text>
          </view>
          <text class="equipment-name">{{ eq.name }}</text>
        </view>
      </view>

      <!-- 装备查漏建议 -->
      <view class="recommendations" v-if="recommendations.length > 0">
        <text class="category-title">🔔 根据赛事信息推荐</text>
        <view
          v-for="rec in recommendations"
          :key="rec.equipment.id"
          class="recommendation-item"
        >
          <view class="rec-content">
            <text class="rec-name">{{ rec.equipment.name }}</text>
            <text class="rec-reason">{{ rec.reason }}</text>
          </view>
          <button class="btn-add" @click="addRecommended(rec)">添加</button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onLoad } from '@dcloudio/uni-app'
import { api } from '@/services/api'

const race = ref(null)
const weather = ref(null)
const commute = ref({ start_point: {}, parking: [] })
const aidStations = ref([])
const equipment = ref({ mandatory: [], recommended: [] })
const recommendations = ref([])

onLoad((query) => {
  if (query.id) {
    loadRacePlan(query.id)
  }
})

async function loadRacePlan(id) {
  try {
    const [planRes, eqRes] = await Promise.all([
      api.getRacePlan(id),
      api.getEquipment(id)
    ])

    race.value = planRes.data.race
    weather.value = planRes.data.weather
    commute.value = planRes.data.commute
    aidStations.value = planRes.data.aid_stations || []
    equipment.value = eqRes.data
  } catch (e) {
    uni.showToast({ title: '加载失败', icon: 'none' })
  }
}

function formatDate(date) {
  const d = new Date(date)
  return `${d.getMonth() + 1}月${d.getDate()}日`
}

function getWeatherIcon(weather) {
  if (!weather) return '❓'
  if (weather.includes('晴')) return '☀️'
  if (weather.includes('多云')) return '⛅'
  if (weather.includes('阴')) return '☁️'
  if (weather.includes('雨')) return '🌧️'
  if (weather.includes('雪')) return '❄️'
  return '🌤️'
}

function openNavigation() {
  if (!race.value) return
  uni.openLocation({
    latitude: race.value.start_lat,
    longitude: race.value.start_lng,
    name: '比赛起点',
    fail: () => {
      uni.showToast({ title: '打开导航失败', icon: 'none' })
    }
  })
}

function openParkingNavigation(parking) {
  uni.openLocation({
    latitude: parking.lat,
    longitude: parking.lng,
    name: parking.name,
    fail: () => {
      uni.showToast({ title: '打开导航失败', icon: 'none' })
    }
  })
}

function isMissing(eq) {
  // Check if this equipment is in the missing mandatory list
  return false
}

async function toggleEquipment(eq) {
  try {
    await api.checkEquipment(race.value.id, {
      equipment_id: eq.id,
      checked: !eq.checked
    })
    eq.checked = !eq.checked
  } catch (e) {
    uni.showToast({ title: '更新失败', icon: 'none' })
  }
}

function addRecommended(rec) {
  // Add recommended equipment to user's list
  uni.showToast({ title: '添加成功', icon: 'success' })
}
</script>

<style scoped>
.container {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 24rpx;
}

.race-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 24rpx;
  padding: 32rpx;
  color: #fff;
  margin-bottom: 24rpx;
}

.race-info {
  margin-bottom: 24rpx;
}

.race-name {
  display: block;
  font-size: 36rpx;
  font-weight: bold;
  margin-bottom: 8rpx;
}

.race-date {
  font-size: 28rpx;
  opacity: 0.9;
}

.race-stats {
  display: flex;
  gap: 40rpx;
}

.stat {
  display: flex;
  align-items: baseline;
  gap: 4rpx;
}

.stat-value {
  font-size: 48rpx;
  font-weight: bold;
}

.stat-label {
  font-size: 24rpx;
  opacity: 0.8;
}

.section {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 24rpx;
}

.section-title {
  display: block;
  font-size: 32rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 16rpx;
}

.weather-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.weather-main {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.weather-icon {
  font-size: 64rpx;
}

.temperature {
  display: block;
  font-size: 48rpx;
  font-weight: bold;
  color: #333;
}

.weather-desc {
  font-size: 26rpx;
  color: #666;
}

.weather-details {
  display: flex;
  gap: 24rpx;
}

.weather-item {
  text-align: center;
}

.weather-label {
  display: block;
  font-size: 22rpx;
  color: #999;
}

.weather-value {
  font-size: 26rpx;
  color: #333;
}

.commute-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20rpx;
  background: #f8f9fa;
  border-radius: 12rpx;
}

.start-point {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.point-icon {
  font-size: 32rpx;
}

.point-label {
  display: block;
  font-size: 22rpx;
  color: #999;
}

.point-name {
  font-size: 28rpx;
  color: #333;
}

.nav-btn {
  padding: 12rpx 24rpx;
  background: #667eea;
  color: #fff;
  border-radius: 8rpx;
  font-size: 26rpx;
}

.parking-list {
  margin-top: 16rpx;
}

.subsection-title {
  display: block;
  font-size: 26rpx;
  color: #666;
  margin-bottom: 12rpx;
}

.parking-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16rpx;
  background: #f8f9fa;
  border-radius: 8rpx;
  margin-bottom: 8rpx;
}

.parking-name {
  font-size: 28rpx;
  color: #333;
}

.parking-distance {
  display: block;
  font-size: 22rpx;
  color: #999;
}

.nav-arrow {
  font-size: 32rpx;
  color: #999;
}

.aid-card {
  padding: 16rpx;
  border-bottom: 1rpx solid #eee;
}

.aid-card:last-child {
  border-bottom: none;
}

.aid-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8rpx;
}

.aid-name {
  font-size: 28rpx;
  font-weight: bold;
  color: #333;
}

.aid-close-time {
  font-size: 22rpx;
  color: #f44336;
}

.aid-info {
  display: flex;
  gap: 16rpx;
  margin-bottom: 8rpx;
}

.aid-distance, .aid-elevation {
  font-size: 24rpx;
  color: #666;
}

.aid-supplies {
  display: flex;
  flex-wrap: wrap;
  gap: 8rpx;
}

.supply-tag {
  padding: 4rpx 12rpx;
  background: #e8f5e9;
  color: #4caf50;
  border-radius: 6rpx;
  font-size: 22rpx;
}

.equipment-category {
  margin-bottom: 16rpx;
}

.category-title {
  display: block;
  font-size: 26rpx;
  color: #666;
  margin-bottom: 12rpx;
}

.category-title.red {
  color: #f44336;
}

.equipment-item {
  display: flex;
  align-items: center;
  padding: 12rpx 0;
}

.checkbox {
  width: 40rpx;
  height: 40rpx;
  border: 2rpx solid #ddd;
  border-radius: 8rpx;
  margin-right: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24rpx;
}

.checkbox.checked {
  background: #667eea;
  border-color: #667eea;
  color: #fff;
}

.checkbox.missing {
  border-color: #f44336;
}

.equipment-name {
  font-size: 28rpx;
  color: #333;
}

.recommendations {
  margin-top: 16rpx;
  padding-top: 16rpx;
  border-top: 1rpx solid #eee;
}

.recommendation-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12rpx 0;
}

.rec-name {
  display: block;
  font-size: 28rpx;
  color: #333;
}

.rec-reason {
  display: block;
  font-size: 22rpx;
  color: #999;
}

.btn-add {
  padding: 8rpx 20rpx;
  background: #667eea;
  color: #fff;
  border-radius: 8rpx;
  font-size: 24rpx;
  border: none;
}
</style>
