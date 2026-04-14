<template>
  <view class="container">
    <view class="header">
      <text class="title">发现赛事</text>
    </view>

    <view class="filters">
      <view class="filter-item" @click="showFilter = true">
        <text>{{ selectedProvince || '全部地区' }}</text>
        <text class="arrow">▼</text>
      </view>
      <view class="filter-item" @click="showDifficultyFilter = true">
        <text>{{ selectedDifficulty || '全部难度' }}</text>
        <text class="arrow">▼</text>
      </view>
    </view>

    <view class="race-list">
      <view
        v-for="race in races"
        :key="race.id"
        class="race-card"
        @click="goToRaceDetail(race.id)"
      >
        <view class="race-header">
          <text class="race-name">{{ race.name }}</text>
          <view class="difficulty-tag" :class="race.difficulty.toLowerCase()">
            {{ race.difficulty }}
          </view>
        </view>
        <view class="race-info">
          <text class="race-date">{{ formatDate(race.date) }}</text>
          <text class="race-location">{{ race.location }}</text>
        </view>
        <view class="race-stats">
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
      </view>

      <view v-if="races.length === 0" class="empty-state">
        <text>暂无赛事</text>
      </view>
    </view>

    <view class="load-more" v-if="races.length > 0" @click="loadMore">
      <text v-if="loading">加载中...</text>
      <text v-else-if="hasMore">加载更多</text>
      <text v-else>没有更多了</text>
    </view>
  </view>
</template>

<script setup>
import { ref, onMount } from 'vue'
import { api } from '@/services/api'

const races = ref([])
const loading = ref(false)
const hasMore = ref(true)
const page = ref(1)
const pageSize = 20

const selectedProvince = ref('')
const selectedDifficulty = ref('')
const showFilter = ref(false)
const showDifficultyFilter = ref(false)

onMount(() => {
  loadRaces()
})

async function loadRaces(refresh = false) {
  if (loading.value) return

  if (refresh) {
    page.value = 1
    races.value = []
    hasMore.value = true
  }

  loading.value = true

  try {
    const params = {
      page: page.value,
      page_size: pageSize
    }

    if (selectedProvince.value) {
      params.province = selectedProvince.value
    }
    if (selectedDifficulty.value) {
      params.difficulty = selectedDifficulty.value
    }

    const res = await api.getRaces(params)

    if (refresh) {
      races.value = res.data.races || []
    } else {
      races.value = [...races.value, ...(res.data.races || [])]
    }

    hasMore.value = races.value.length < res.data.total
    page.value++
  } catch (e) {
    uni.showToast({ title: '加载失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

function loadMore() {
  if (hasMore.value && !loading.value) {
    loadRaces()
  }
}

function formatDate(date) {
  const d = new Date(date)
  return `${d.getMonth() + 1}月${d.getDate()}日`
}

function goToRaceDetail(id) {
  uni.navigateTo({ url: `/pages/race-detail/race-detail?id=${id}` })
}
</script>

<style scoped>
.container {
  min-height: 100vh;
  background: #f5f5f5;
}

.header {
  padding: 32rpx;
  background: #fff;
}

.title {
  font-size: 40rpx;
  font-weight: bold;
  color: #333;
}

.filters {
  display: flex;
  padding: 24rpx 32rpx;
  background: #fff;
  border-bottom: 1rpx solid #eee;
}

.filter-item {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28rpx;
  color: #666;
}

.filter-item:first-child {
  border-right: 1rpx solid #eee;
}

.arrow {
  margin-left: 8rpx;
  font-size: 20rpx;
}

.race-list {
  padding: 24rpx;
}

.race-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 24rpx;
}

.race-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
}

.race-name {
  font-size: 32rpx;
  font-weight: bold;
  color: #333;
}

.difficulty-tag {
  padding: 6rpx 16rpx;
  border-radius: 8rpx;
  font-size: 22rpx;
}

.difficulty-tag.easy {
  background: #e8f5e9;
  color: #4caf50;
}

.difficulty-tag.medium {
  background: #fff3e0;
  color: #ff9800;
}

.difficulty-tag.hard {
  background: #ffebee;
  color: #f44336;
}

.difficulty-tag.extrem {
  background: #f3e5f5;
  color: #9c27b0;
}

.race-info {
  display: flex;
  gap: 24rpx;
  margin-bottom: 16rpx;
}

.race-date, .race-location {
  font-size: 26rpx;
  color: #666;
}

.race-stats {
  display: flex;
  gap: 32rpx;
}

.stat {
  display: flex;
  align-items: baseline;
  gap: 4rpx;
}

.stat-value {
  font-size: 32rpx;
  font-weight: bold;
  color: #333;
}

.stat-label {
  font-size: 22rpx;
  color: #999;
}

.empty-state {
  text-align: center;
  padding: 80rpx 0;
  color: #999;
}

.load-more {
  text-align: center;
  padding: 32rpx;
  color: #666;
  font-size: 26rpx;
}
</style>
