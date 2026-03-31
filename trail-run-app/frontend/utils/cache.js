/**
 * 离线缓存策略实现
 * 使用 uni-storage 进行数据持久化缓存
 */

const CACHE_PREFIX = 'trail_run_'
const DEFAULT_EXPIRE = 30 * 60 * 1000 // 30分钟

/**
 * 设置缓存
 * @param {string} key - 缓存键
 * @param {any} value - 缓存值
 * @param {number} expire - 过期时间(ms)，默认30分钟
 */
export function setCache(key, value, expire = DEFAULT_EXPIRE) {
  const data = {
    value,
    expire: Date.now() + expire,
    version: 1
  }
  uni.setStorageSync(CACHE_PREFIX + key, JSON.stringify(data))
}

/**
 * 获取缓存
 * @param {string} key - 缓存键
 * @param {any} defaultValue - 默认值
 * @returns {any}
 */
export function getCache(key, defaultValue = null) {
  const str = uni.getStorageSync(CACHE_PREFIX + key)
  if (!str) return defaultValue

  try {
    const data = JSON.parse(str)
    // 检查过期
    if (data.expire && Date.now() > data.expire) {
      removeCache(key)
      return defaultValue
    }
    return data.value
  } catch (e) {
    return defaultValue
  }
}

/**
 * 移除缓存
 * @param {string} key - 缓存键
 */
export function removeCache(key) {
  uni.removeStorageSync(CACHE_PREFIX + key)
}

/**
 * 清空所有应用缓存
 */
export function clearCache() {
  const keys = ['races', 'user', 'equipment', 'results']
  keys.forEach(key => removeCache(key))
}

/**
 * 缓存用户信息
 */
export function cacheUser(user) {
  setCache('user', user, 60 * 60 * 1000) // 1小时
}

/**
 * 获取缓存的用户信息
 */
export function getCachedUser() {
  return getCache('user')
}

/**
 * 缓存比赛列表
 */
export function cacheRaces(races) {
  setCache('races', races, 10 * 60 * 1000) // 10分钟
}

/**
 * 获取缓存的比赛列表
 */
export function getCachedRaces() {
  return getCache('races')
}

/**
 * 缓存比赛详情
 */
export function cacheRaceDetail(raceId, detail) {
  setCache(`race_${raceId}`, detail, 15 * 60 * 1000) // 15分钟
}

/**
 * 获取缓存的比赛详情
 */
export function getCachedRaceDetail(raceId) {
  return getCache(`race_${raceId}`)
}

/**
 * 缓存装备清单
 */
export function cacheEquipment(raceId, equipment) {
  setCache(`equipment_${raceId}`, equipment, 30 * 60 * 1000)
}

/**
 * 获取缓存的装备清单
 */
export function getCachedEquipment(raceId) {
  return getCache(`equipment_${raceId}`)
}

/**
 * 缓存我的比赛列表
 */
export function cacheResults(results) {
  setCache('results', results, 5 * 60 * 1000) // 5分钟
}

/**
 * 获取缓存的比赛结果
 */
export function getCachedResults() {
  return getCache('results')
}

export default {
  setCache,
  getCache,
  removeCache,
  clearCache,
  cacheUser,
  getCachedUser,
  cacheRaces,
  getCachedRaces,
  cacheRaceDetail,
  getCachedRaceDetail,
  cacheEquipment,
  getCachedEquipment,
  cacheResults,
  getCachedResults
}
