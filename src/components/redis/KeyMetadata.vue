<template>
  <el-card shadow="never" class="metadata-card">
    <template #header>
      <div class="header">
        <span>键信息</span>
        <el-button size="small" @click="refresh" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </template>
    
    <div class="metadata-content" v-loading="loading">
      <div class="metadata-grid">
        <div class="metadata-item">
          <div class="label">键名</div>
          <div class="value key-name">{{ keyName }}</div>
        </div>
        
        <div class="metadata-item">
          <div class="label">类型</div>
          <div class="value">
            <el-tag :type="getTypeColor(metadata.type)" size="small">
              {{ metadata.type || 'unknown' }}
            </el-tag>
          </div>
        </div>
        
        <div class="metadata-item">
          <div class="label">TTL</div>
          <div class="value">
            <span v-if="metadata.ttl === -1" class="no-expiry">永不过期</span>
            <span v-else-if="metadata.ttl === -2" class="expired">已过期</span>
            <span v-else class="ttl-value">
              {{ formatTTL(metadata.ttl) }}
              <el-tag v-if="metadata.ttl > 0" size="small" type="warning">
                {{ metadata.ttl }}秒
              </el-tag>
            </span>
          </div>
        </div>
        
        <div class="metadata-item" v-if="metadata.length !== undefined">
          <div class="label">长度</div>
          <div class="value">{{ formatNumber(metadata.length) }}</div>
        </div>
        
        <div class="metadata-item" v-if="metadata.memoryUsage > 0">
          <div class="label">内存使用</div>
          <div class="value">{{ formatBytes(metadata.memoryUsage) }}</div>
        </div>
        
        <div class="metadata-item" v-if="metadata.encoding">
          <div class="label">编码</div>
          <div class="value">
            <el-tag size="small" type="info">{{ metadata.encoding }}</el-tag>
          </div>
        </div>
        
        <div class="metadata-item" v-if="metadata.refCount !== undefined">
          <div class="label">引用计数</div>
          <div class="value">{{ metadata.refCount }}</div>
        </div>
        
        <div class="metadata-item" v-if="metadata.idleTime !== undefined">
          <div class="label">空闲时间</div>
          <div class="value">{{ formatDuration(metadata.idleTime) }}</div>
        </div>
        
        <div class="metadata-item">
          <div class="label">存在状态</div>
          <div class="value">
            <el-tag :type="metadata.exists ? 'success' : 'danger'" size="small">
              {{ metadata.exists ? '存在' : '不存在' }}
            </el-tag>
          </div>
        </div>
      </div>
    </div>
  </el-card>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { keyMetadata } from '@/api/redis'

const props = defineProps({
  keyName: {
    type: String,
    required: true
  }
})

const loading = ref(false)
const metadata = ref({})

const getTypeColor = (type) => {
  const colors = {
    string: 'primary',
    list: 'success',
    set: 'warning',
    zset: 'danger',
    hash: 'info',
    none: 'info'
  }
  return colors[type] || 'info'
}

const formatTTL = (ttl) => {
  if (ttl === -1) return '永不过期'
  if (ttl === -2) return '已过期'
  if (ttl <= 0) return '已过期'
  
  const days = Math.floor(ttl / 86400)
  const hours = Math.floor((ttl % 86400) / 3600)
  const minutes = Math.floor((ttl % 3600) / 60)
  const seconds = ttl % 60
  
  const parts = []
  if (days > 0) parts.push(`${days}天`)
  if (hours > 0) parts.push(`${hours}小时`)
  if (minutes > 0) parts.push(`${minutes}分钟`)
  if (seconds > 0) parts.push(`${seconds}秒`)
  
  return parts.join(' ') || '即将过期'
}

const formatBytes = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatNumber = (num) => {
  return num.toLocaleString()
}

const formatDuration = (seconds) => {
  if (seconds < 60) return `${seconds}秒`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}分钟`
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}小时`
  return `${Math.floor(seconds / 86400)}天`
}

const loadMetadata = async () => {
  if (!props.keyName) return
  
  loading.value = true
  try {
    metadata.value = await keyMetadata(props.keyName)
  } catch (error) {
    console.error('获取键元数据失败:', error)
    metadata.value = {}
  } finally {
    loading.value = false
  }
}

const refresh = () => {
  loadMetadata()
}

// 监听keyName变化
watch(() => props.keyName, () => {
  loadMetadata()
}, { immediate: true })

onMounted(() => {
  loadMetadata()
})
</script>

<style scoped>
.metadata-card {
  margin-bottom: 16px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.metadata-content {
  min-height: 200px;
}

.metadata-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 16px;
}

.metadata-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 12px;
  background: var(--el-fill-color-lighter);
  border-radius: 6px;
  border: 1px solid var(--el-border-color-light);
}

.label {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  font-weight: 500;
}

.value {
  font-size: 14px;
  color: var(--el-text-color-primary);
  word-break: break-all;
}

.key-name {
  font-family: monospace;
  background: var(--el-fill-color-blank);
  padding: 4px 8px;
  border-radius: 4px;
  border: 1px solid var(--el-border-color);
}

.ttl-value {
  display: flex;
  align-items: center;
  gap: 8px;
}

.no-expiry {
  color: var(--el-color-success);
  font-weight: 500;
}

.expired {
  color: var(--el-color-danger);
  font-weight: 500;
}

@media (max-width: 768px) {
  .metadata-grid {
    grid-template-columns: 1fr;
  }
}
</style>
