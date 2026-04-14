<template>
  <view class="container">
    <view class="header">
      <text class="title">我的比赛</text>
    </view>

    <view class="race-list">
      <view
        v-for="result in results"
        :key="result.id"
        class="result-card"
        @click="goToResultDetail(result.id)"
      >
        <view class="result-header">
          <text class="race-name">{{ result.race_name }}</text>
          <view class="status-tag" :class="result.status">
            {{ statusText(result.status) }}
          </view>
        </view>

        <view class="result-info">
          <text class="race-date">{{ formatDate(result.race_date) }}</text>
        </view>

        <view class="result-stats">
          <view class="stat">
            <text class="stat-value">{{ result.distance_km }}</text>
            <text class="stat-label">公里</text>
          </view>
          <view class="stat">
            <text class="stat-value">{{ result.elevation_m }}</text>
            <text class="stat-label">爬升(m)</text>
          </view>
          <view class="stat" v-if="result.finish_time">
            <text class="stat-value">{{ result.finish_time }}</text>
            <text class="stat-label">完赛时间</text>
          </view>
          <view class="stat" v-if="result.ranking">
            <text class="stat-value">#{{ result.ranking }}</text>
            <text class="stat-label">总排名</text>
          </view>
        </view>

        <view class="result-image" v-if="result.generated_image_url">
          <image :src="result.generated_image_url" mode="aspectFill" />
        </view>
      </view>

      <view v-if="results.length === 0" class="empty-state">
        <text class="empty-icon">🏅</text>
        <text class="empty-text">暂无完赛记录</text>
        <button class="btn-discover" @click="goToRaces">去发现比赛</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onShow } from 'vue'
import { api } from '@/services/api'

const results = ref([])

onShow(() => {
  loadResults()
})

async function loadResults() {
  try {
    const res = await api.getResults()
    results.value = res.data.results || []
  } catch (e) {
    console.error('Failed to load results:', e)
  }
}

function formatDate(date) {
  const d = new Date(date)
  return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日`
}

function statusText(status) {
  const map = {
    registered: '已报名',
    ongoing: '进行中',
    finished: '已完赛'
  }
  return map[status] || status
}

function goToResultDetail(id) {
  uni.navigateTo({ url: `/pages/result-detail/result-detail?id=${id}` })
}

function goToRaces() {
  uni.switchTab({ url: '/pages/race-list/race-list' })
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

.race-list {
  padding: 24rpx;
}

.result-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 24rpx;
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12rpx;
}

.race-name {
  font-size: 32rpx;
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

.result-info {
  margin-bottom: 16rpx;
}

.race-date {
  font-size: 26rpx;
  color: #666;
}

.result-stats {
  display: flex;
  gap: 32rpx;
  flex-wrap: wrap;
}

.stat {
  display: flex;
  align-items: baseline;
  gap: 4rpx;
}

.stat-value {
  font-size: 28rpx;
  font-weight: bold;
  color: #333;
}

.stat-label {
  font-size: 22rpx;
  color: #999;
}

.result-image {
  margin-top: 16rpx;
  border-radius: 12rpx;
  overflow: hidden;
}

.result-image image {
  width: 100%;
  height: 300rpx;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 100rpx 0;
}

.empty-icon {
  font-size: 80rpx;
  margin-bottom: 24rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #999;
  margin-bottom: 32rpx;
}

.btn-discover {
  padding: 16rpx 48rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  border-radius: 32rpx;
  font-size: 28rpx;
  border: none;
}
</style>
