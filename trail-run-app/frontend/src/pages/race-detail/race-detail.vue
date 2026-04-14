<template>
  <view class="container">
    <view v-if="loading" class="loading">加载中...</view>

    <view v-else-if="race" class="race-detail">
      <view class="map-section">
        <map
          id="raceMap"
          class="map"
          :latitude="race.start_lat"
          :longitude="race.start_lng"
          :scale="12"
          :markers="markers"
          :polyline="polyline"
          show-location
        />
      </view>

      <view class="content">
        <view class="header">
          <text class="race-name">{{ race.name }}</text>
          <view class="difficulty-tag" :class="race.difficulty.toLowerCase()">
            {{ race.difficulty }}
          </view>
        </view>

        <view class="info-row">
          <text class="info-item">📅 {{ formatDate(race.date) }}</text>
          <text class="info-item">📍 {{ race.location }}</text>
        </view>

        <view class="stats-row">
          <view class="stat">
            <text class="stat-value">{{ race.distance_km }}</text>
            <text class="stat-label">公里</text>
          </view>
          <view class="stat">
            <text class="stat-value">{{ race.elevation_m }}</text>
            <text class="stat-label">爬升(m)</text>
          </view>
          <view class="stat" v-if="race.itra_points">
            <text class="stat-value">{{ race.itra_points }}</text>
            <text class="stat-label">ITRA积分</text>
          </view>
        </view>

        <view class="section">
          <text class="section-title">补给站 ({{ race.aid_stations?.length || 0 }})</text>
          <view
            v-for="aid in race.aid_stations"
            :key="aid.id"
            class="aid-item"
          >
            <view class="aid-info">
              <text class="aid-name">{{ aid.name }}</text>
              <text class="aid-distance">{{ aid.distance_km }}km | 海拔{{ aid.elevation_m }}m</text>
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

        <view class="section">
          <text class="section-title">装备清单</text>
          <view v-if="equipment.mandatory?.length" class="equipment-category">
            <text class="category-title">强制装备</text>
            <view
              v-for="eq in equipment.mandatory"
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
        </view>
      </view>

      <view class="bottom-actions">
        <button class="btn-secondary" @click="goToNavigation">导航到起点</button>
        <button class="btn-primary" @click="registerRace">报名该赛事</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onLoad } from '@dcloudio/uni-app'
import { api } from '@/services/api'

const race = ref(null)
const loading = ref(true)
const equipment = ref({})

function formatDate(date) {
  const d = new Date(date)
  return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日`
}

const markers = computed(() => {
  if (!race.value?.aid_stations) return []

  const result = []

  result.push({
    id: 0,
    latitude: race.value.start_lat,
    longitude: race.value.start_lng,
    iconPath: '/static/marker-start.png',
    width: 32,
    height: 32,
    callout: {
      content: '起点',
      color: '#333',
      fontSize: 12,
      borderRadius: 8,
      padding: 8,
      display: 'ALWAYS'
    }
  })

  race.value.aid_stations.forEach((aid, index) => {
    result.push({
      id: index + 1,
      latitude: aid.lat,
      longitude: aid.lng,
      iconPath: '/static/marker-aid.png',
      width: 28,
      height: 28,
      callout: {
        content: `${aid.name}\n${aid.distance_km}km`,
        color: '#333',
        fontSize: 12,
        borderRadius: 8,
        padding: 8,
        display: 'BYCLICK'
      }
    })
  })

  return result
})

const polyline = computed(() => {
  return []
})

async function loadRaceDetail(id) {
  loading.value = true
  try {
    const res = await api.getRace(id)
    race.value = res.data

    const eqRes = await api.getEquipment(id)
    equipment.value = eqRes.data
  } catch (e) {
    uni.showToast({ title: '加载失败', icon: 'none' })
  } finally {
    loading.value = false
  }
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

function goToNavigation() {
  if (!race.value) return

  uni.openLocation({
    latitude: race.value.start_lat,
    longitude: race.value.start_lng,
    name: race.value.name,
    fail: () => {
      uni.showToast({ title: '打开导航失败', icon: 'none' })
    }
  })
}

function registerRace() {
  uni.showToast({ title: '报名功能待开发', icon: 'none' })
}

onLoad((query) => {
  if (query.id) {
    loadRaceDetail(query.id)
  }
})
</script>

<style scoped>
.container {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 120rpx;
}

.loading {
  text-align: center;
  padding: 100rpx;
  color: #666;
}

.map-section {
  height: 400rpx;
}

.map {
  width: 100%;
  height: 100%;
}

.content {
  padding: 32rpx;
  background: #fff;
  margin-top: -24rpx;
  border-radius: 24rpx 24rpx 0 0;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24rpx;
}

.race-name {
  font-size: 40rpx;
  font-weight: bold;
  color: #333;
}

.difficulty-tag {
  padding: 8rpx 20rpx;
  border-radius: 8rpx;
  font-size: 24rpx;
}

.difficulty-tag.easy { background: #e8f5e9; color: #4caf50; }
.difficulty-tag.medium { background: #fff3e0; color: #ff9800; }
.difficulty-tag.hard { background: #ffebee; color: #f44336; }
.difficulty-tag.extrem { background: #f3e5f5; color: #9c27b0; }

.info-row {
  display: flex;
  gap: 24rpx;
  margin-bottom: 24rpx;
}

.info-item {
  font-size: 26rpx;
  color: #666;
}

.stats-row {
  display: flex;
  gap: 40rpx;
  padding: 24rpx 0;
  border-top: 1rpx solid #eee;
  border-bottom: 1rpx solid #eee;
  margin-bottom: 24rpx;
}

.stat {
  display: flex;
  align-items: baseline;
  gap: 4rpx;
}

.stat-value {
  font-size: 40rpx;
  font-weight: bold;
  color: #333;
}

.stat-label {
  font-size: 24rpx;
  color: #999;
}

.section {
  margin-bottom: 32rpx;
}

.section-title {
  display: block;
  font-size: 32rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 16rpx;
}

.aid-item {
  padding: 16rpx 0;
  border-bottom: 1rpx solid #eee;
}

.aid-item:last-child {
  border-bottom: none;
}

.aid-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8rpx;
}

.aid-name {
  font-size: 28rpx;
  color: #333;
}

.aid-distance {
  font-size: 24rpx;
  color: #999;
}

.aid-supplies {
  display: flex;
  flex-wrap: wrap;
  gap: 8rpx;
}

.supply-tag {
  padding: 4rpx 12rpx;
  background: #f5f5f5;
  border-radius: 6rpx;
  font-size: 22rpx;
  color: #666;
}

.equipment-category {
  margin-bottom: 16rpx;
}

.category-title {
  font-size: 26rpx;
  color: #666;
  margin-bottom: 12rpx;
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

.equipment-name {
  font-size: 28rpx;
  color: #333;
}

.bottom-actions {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  gap: 24rpx;
  padding: 24rpx 32rpx;
  background: #fff;
  box-shadow: 0 -4rpx 16rpx rgba(0, 0, 0, 0.05);
}

.btn-secondary {
  flex: 1;
  height: 88rpx;
  line-height: 88rpx;
  background: #fff;
  color: #667eea;
  border: 2rpx solid #667eea;
  border-radius: 16rpx;
  font-size: 28rpx;
}

.btn-primary {
  flex: 2;
  height: 88rpx;
  line-height: 88rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  border-radius: 16rpx;
  font-size: 28rpx;
  font-weight: bold;
  border: none;
}
</style>
