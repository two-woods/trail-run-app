<template>
  <view class="container">
    <!-- 未登录状态 -->
    <view v-if="!isLoggedIn" class="guest-view">
      <view class="hero">
        <text class="hero-title">越野跑</text>
        <text class="hero-subtitle">记录每一次突破</text>
      </view>
      <button class="btn-login" @click="goToLogin">登录/注册</button>
    </view>

    <!-- 已登录状态 -->
    <view v-else class="user-view">
      <view class="user-header">
        <image
          class="avatar"
          :src="userInfo.avatar_url || '/static/default-avatar.png'"
          mode="aspectFill"
        />
        <view class="user-info">
          <text class="nickname">{{ userInfo.nickname }}</text>
          <text class="welcome">开始今天的训练</text>
        </view>
      </view>

      <view class="quick-actions">
        <view class="action-card" @click="goToRaces">
          <text class="action-icon">🏃</text>
          <text class="action-title">赛事</text>
          <text class="action-desc">发现更多比赛</text>
        </view>

        <view class="action-card" @click="goToMyRaces">
          <text class="action-icon">🏅</text>
          <text class="action-title">我的比赛</text>
          <text class="action-desc">查看完赛记录</text>
        </view>
      </view>

      <view class="section">
        <text class="section-title">即将到来的比赛</text>
        <view class="empty-state" v-if="upcomingRaces.length === 0">
          <text>暂无报名的比赛</text>
          <button class="btn-discover" @click="goToRaces">去发现</button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onShow } from 'vue'
import { api } from '@/services/api'

const isLoggedIn = ref(false)
const userInfo = ref({})
const upcomingRaces = ref([])

onShow(() => {
  checkAuth()
})

function checkAuth() {
  const token = uni.getStorageSync('token')
  const user = uni.getStorageSync('user')

  if (token && user) {
    isLoggedIn.value = true
    userInfo.value = user
    loadUpcomingRaces()
  } else {
    isLoggedIn.value = false
    userInfo.value = {}
  }
}

async function loadUpcomingRaces() {
  try {
    // TODO: 获取用户报名的比赛
  } catch (e) {
    console.error('Failed to load races:', e)
  }
}

function goToLogin() {
  uni.navigateTo({ url: '/pages/login/login' })
}

function goToRaces() {
  uni.switchTab({ url: '/pages/race-list/race-list' })
}

function goToMyRaces() {
  uni.navigateTo({ url: '/pages/my-races/my-races' })
}
</script>

<style scoped>
.container {
  min-height: 100vh;
  background: #f5f5f5;
}

.guest-view {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 40rpx;
}

.hero {
  text-align: center;
  margin-bottom: 80rpx;
}

.hero-title {
  display: block;
  font-size: 72rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 16rpx;
}

.hero-subtitle {
  display: block;
  font-size: 32rpx;
  color: #666;
}

.btn-login {
  width: 400rpx;
  height: 88rpx;
  line-height: 88rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  font-size: 32rpx;
  font-weight: bold;
  border-radius: 44rpx;
  border: none;
}

.user-view {
  padding: 32rpx;
}

.user-header {
  display: flex;
  align-items: center;
  background: #fff;
  padding: 32rpx;
  border-radius: 24rpx;
  margin-bottom: 32rpx;
}

.avatar {
  width: 100rpx;
  height: 100rpx;
  border-radius: 50%;
  background: #e9ecef;
}

.user-info {
  margin-left: 24rpx;
}

.nickname {
  display: block;
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
}

.welcome {
  display: block;
  font-size: 26rpx;
  color: #666;
  margin-top: 8rpx;
}

.quick-actions {
  display: flex;
  gap: 24rpx;
  margin-bottom: 32rpx;
}

.action-card {
  flex: 1;
  background: #fff;
  padding: 32rpx 24rpx;
  border-radius: 24rpx;
  text-align: center;
}

.action-icon {
  display: block;
  font-size: 48rpx;
  margin-bottom: 16rpx;
}

.action-title {
  display: block;
  font-size: 30rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 8rpx;
}

.action-desc {
  display: block;
  font-size: 24rpx;
  color: #666;
}

.section {
  background: #fff;
  padding: 32rpx;
  border-radius: 24rpx;
}

.section-title {
  display: block;
  font-size: 32rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 24rpx;
}

.empty-state {
  text-align: center;
  padding: 40rpx 0;
  color: #999;
}

.btn-discover {
  margin-top: 24rpx;
  padding: 16rpx 48rpx;
  background: #007aff;
  color: #fff;
  border-radius: 32rpx;
  font-size: 28rpx;
  border: none;
}
</style>
