import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as api from '@/api/redis'

export const useRedisStore = defineStore('redis', () => {
  const profiles = ref([])
  const currentProfileId = ref(null)
  const connecting = ref(false)
  const connected = ref(false)

  const pattern = ref('*')
  const typeFilter = ref('')
  const scanning = ref(false)
  const cursor = ref('0')
  const keys = ref([])
  const hasMore = ref(false)

  const selectedKey = ref('')
  const selectedType = ref('')

  function init() {
    // 从 localStorage 读取 profiles
    const raw = localStorage.getItem('redis_profiles')
    profiles.value = raw ? JSON.parse(raw) : []
    const last = localStorage.getItem('redis_current_profile_id')
    currentProfileId.value = last || (profiles.value[0]?.id || null)
  }

  function persistProfiles() {
    localStorage.setItem('redis_profiles', JSON.stringify(profiles.value))
  }

  function selectProfile(id) {
    currentProfileId.value = id
    localStorage.setItem('redis_current_profile_id', id || '')
  }

  async function saveProfile(profile) {
    const exists = profile.id
    if (!exists) profile.id = `p_${Date.now()}`
    const idx = profiles.value.findIndex(p => p.id === profile.id)
    if (idx >= 0) profiles.value[idx] = { ...profiles.value[idx], ...profile }
    else profiles.value.push({ ...profile })
    persistProfiles()
  }

  async function deleteProfile(id) {
    profiles.value = profiles.value.filter(p => p.id !== id)
    if (currentProfileId.value === id) currentProfileId.value = null
    persistProfiles()
  }

  async function connect(id) {
    if (!id) id = currentProfileId.value
    if (!id) return
    connecting.value = true
    try {
      const prof = profiles.value.find(p => p.id === id)
      if (!prof) return
      const { host, port, db, password } = prof
      const res = await api.testConnect({ host, port, db, password })
      if (res?.ok) {
        api.setConnection({ host, port, db, password })
        connected.value = true
        currentProfileId.value = id
        // 自动加载：连接后立即扫描所有 key
        pattern.value = pattern.value || '*'
        typeFilter.value = ''
        cursor.value = '0'
        keys.value = []
        await loadMore()
      }
    } finally {
      connecting.value = false
    }
  }

  async function disconnect() {
    connected.value = false
  }

  async function search(params) {
    if (params?.pattern !== undefined) pattern.value = params.pattern
    if (params?.type !== undefined) typeFilter.value = params.type
    cursor.value = '0'
    keys.value = []
    await loadMore()
  }

  async function resetSearch() {
    pattern.value = '*'
    typeFilter.value = ''
    await search({})
  }

  async function loadMore() {
    if (!connected.value) return
    scanning.value = true
    try {
      const res = await api.scan({ cursor: cursor.value, match: pattern.value || undefined, count: 100, type: typeFilter.value || undefined })
      cursor.value = res.cursor
      const mapped = res.keys
      keys.value = keys.value.concat(mapped)
      hasMore.value = cursor.value !== '0'
    } finally {
      scanning.value = false
    }
  }

  function selectKey(item) {
    selectedKey.value = item?.key || ''
    selectedType.value = item?.type || ''
  }

  async function deleteKey(key) {
    await api.keyDelete(key)
    keys.value = keys.value.filter(i => i.key !== key)
  }

  async function renameKey(key, newKey) { await api.keyRename(key, newKey); await refreshSelected() }
  async function expireKey(key, ttl) { await api.keyExpire(key, ttl); await refreshSelected() }
  async function persistKey(key) { await api.keyPersist(key); await refreshSelected() }
  async function refreshSelected() {
    if (!selectedKey.value) return
    // 可以在此刷新 TTL/Type，如需要
  }

  return {
    profiles, currentProfileId, connecting, connected,
    pattern, typeFilter, scanning, cursor, keys, hasMore,
    selectedKey, selectedType,
    init, saveProfile, deleteProfile, selectProfile,
    connect, disconnect,
    search, resetSearch, loadMore,
    selectKey, deleteKey, renameKey, expireKey, persistKey, refreshSelected,
  }
})


