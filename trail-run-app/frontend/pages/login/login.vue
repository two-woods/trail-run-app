<template>
  <view class="container">
    <view class="header">
      <text class="title">越野跑</text>
      <text class="subtitle">记录每一次突破</text>
    </view>

    <view class="form">
      <view class="input-group">
        <input
          type="text"
          v-model="email"
          placeholder="邮箱"
          class="input"
        />
      </view>

      <view class="input-group">
        <input
          :type="showPassword ? 'text' : 'password'"
          v-model="password"
          placeholder="密码"
          class="input"
        />
        <text class="toggle-password" @click="showPassword = !showPassword">
          {{ showPassword ? '隐藏' : '显示' }}
        </text>
      </view>

      <view class="input-group" v-if="isRegister">
        <input
          type="text"
          v-model="nickname"
          placeholder="昵称"
          class="input"
        />
      </view>

      <button class="btn-primary" @click="handleSubmit" :loading="loading">
        {{ isRegister ? '注册' : '登录' }}
      </button>

      <text class="switch-mode" @click="isRegister = !isRegister">
        {{ isRegister ? '已有账号？立即登录' : '没有账号？立即注册' }}
      </text>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { api } from '@/services/api'

const email = ref('')
const password = ref('')
const nickname = ref('')
const showPassword = ref(false)
const loading = ref(false)
const isRegister = ref(false)

async function handleSubmit() {
  if (!email.value || !password.value) {
    uni.showToast({ title: '请填写完整信息', icon: 'none' })
    return
  }

  if (isRegister.value && !nickname.value) {
    uni.showToast({ title: '请输入昵称', icon: 'none' })
    return
  }

  loading.value = true

  try {
    let res
    if (isRegister.value) {
      res = await api.register({
        email: email.value,
        password: password.value,
        nickname: nickname.value
      })
    } else {
      res = await api.login({
        email: email.value,
        password: password.value
      })
    }

    // 保存token
    uni.setStorageSync('token', res.token)
    uni.setStorageSync('user', {
      id: res.user_id,
      nickname: res.nickname,
      avatar_url: res.avatar_url
    })

    uni.showToast({
      title: isRegister.value ? '注册成功' : '登录成功',
      icon: 'success'
    })

    // 跳转到首页
    setTimeout(() => {
      uni.reLaunch({ url: '/pages/index/index' })
    }, 1500)
  } catch (e) {
    uni.showToast({
      title: e.error?.message || (isRegister.value ? '注册失败' : '登录失败'),
      icon: 'none'
    })
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.container {
  padding: 60rpx 40rpx;
  min-height: 100vh;
  background: linear-gradient(180deg, #f8f9fa 0%, #e9ecef 100%);
}

.header {
  text-align: center;
  margin-bottom: 80rpx;
  padding-top: 100rpx;
}

.title {
  display: block;
  font-size: 56rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 16rpx;
}

.subtitle {
  display: block;
  font-size: 28rpx;
  color: #666;
}

.form {
  background: #fff;
  border-radius: 24rpx;
  padding: 48rpx 32rpx;
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.08);
}

.input-group {
  margin-bottom: 32rpx;
  position: relative;
}

.input {
  width: 100%;
  height: 88rpx;
  padding: 0 32rpx;
  border: 2rpx solid #e9ecef;
  border-radius: 16rpx;
  font-size: 28rpx;
  box-sizing: border-box;
}

.input:focus {
  border-color: #007aff;
}

.toggle-password {
  position: absolute;
  right: 32rpx;
  top: 50%;
  transform: translateY(-50%);
  font-size: 26rpx;
  color: #666;
  z-index: 10;
}

.btn-primary {
  width: 100%;
  height: 88rpx;
  line-height: 88rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  font-size: 32rpx;
  font-weight: bold;
  border-radius: 16rpx;
  border: none;
  margin-top: 16rpx;
}

.switch-mode {
  display: block;
  text-align: center;
  margin-top: 32rpx;
  font-size: 26rpx;
  color: #007aff;
}
</style>
