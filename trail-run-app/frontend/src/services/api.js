// API 请求封装
const BASE_URL = '/api/v1'

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

// 带缓存的请求
const cachedRequest = async (options, cacheKey, cacheTime = 10 * 60 * 1000) => {
  // GET请求优先使用缓存
  if (options.method === 'GET' || !options.method) {
    const cached = uni.getStorageSync('trail_run_' + cacheKey)
    if (cached) {
      try {
        const data = JSON.parse(cached)
        if (data.expire > Date.now()) {
          return data.value
        }
      } catch (e) {}
    }
  }

  const result = await request(options)

  // 缓存结果
  if (options.method === 'GET' || !options.method) {
    const cacheData = {
      value: result,
      expire: Date.now() + cacheTime
    }
    uni.setStorageSync('trail_run_' + cacheKey, JSON.stringify(cacheData))
  }

  return result
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

  // Hotels
  searchHotels: (params) => request({ url: '/hotels/search', method: 'GET', data: params }),
  getHotels: (raceId) => request({ url: `/races/${raceId}/hotels` }),
  getFavoriteHotels: (raceId) => request({ url: '/hotels/favorites', method: 'GET', data: { race_id: raceId } }),
  createFavoriteHotel: (data) => request({ url: '/hotels/favorites', method: 'POST', data }),
  deleteFavoriteHotel: (id) => request({ url: `/hotels/favorites/${id}`, method: 'DELETE' }),
}
