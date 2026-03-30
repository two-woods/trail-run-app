// API 请求封装
const BASE_URL = 'http://localhost:8080/api/v1'

const request = (options) => {
  return new Promise((resolve, reject) => {
    uni.request({
      url: BASE_URL + options.url,
      method: options.method || 'GET',
      data: options.data || {},
      header: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + uni.getStorageSync('token'),
        ...options.header
      },
      timeout: 10000,
      success: (res) => {
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve(res.data)
        } else if (res.statusCode === 401) {
          uni.removeStorageSync('token')
          uni.reLaunch({ url: '/pages/login/login' })
          reject(new Error('Unauthorized'))
        } else {
          reject(res.data)
        }
      },
      fail: (err) => {
        uni.showToast({ title: '网络错误', icon: 'none' })
        reject(err)
      }
    })
  })
}

export const api = {
  // Auth
  register: (data) => request({ url: '/auth/register', method: 'POST', data }),
  login: (data) => request({ url: '/auth/login', method: 'POST', data }),

  // User
  getCurrentUser: () => request({ url: '/users/me' }),

  // Races
  getRaces: (params) => request({ url: '/races', method: 'GET', data: params }),
  getRace: (id) => request({ url: `/races/${id}` }),
  getRacePlan: (id) => request({ url: `/races/${id}/plan` }),
  getEquipment: (raceId) => request({ url: `/races/${raceId}/equipment` }),
  checkEquipment: (raceId, data) => request({ url: `/races/${raceId}/equipment/check`, method: 'POST', data }),

  // Results
  getResults: () => request({ url: '/results' }),
  createResult: (data) => request({ url: '/results', method: 'POST', data }),
  getResult: (id) => request({ url: `/results/${id}` }),
  generateImage: (id, data) => request({ url: `/results/${id}/image`, method: 'POST', data }),
}
