<template>
  <div class="command-executor">
    <div class="command-header">
      <h3>Redis 命令执行器</h3>
      <p class="description">输入Redis命令并执行，支持所有Redis命令</p>
    </div>
    
    <div class="command-input-section">
      <el-input
        v-model="command"
        placeholder="输入Redis命令，例如: LPUSH cut_job 123"
        @keyup.enter="executeCommand"
        :disabled="!connected || executing"
        class="command-input"
        size="large"
      >
        <template #prepend>
          <el-icon><Monitor /></el-icon>
        </template>
        <template #append>
          <el-button 
            type="primary" 
            @click="executeCommand"
            :loading="executing"
            :disabled="!connected || !command.trim()"
          >
            执行
          </el-button>
        </template>
      </el-input>
    </div>

    <!-- 常用命令快捷按钮 -->
    <div class="quick-commands">
      <h4>常用命令</h4>
      <div class="command-buttons">
        <el-button 
          v-for="cmd in quickCommands" 
          :key="cmd.name"
          @click="setCommand(cmd.command)"
          :disabled="!connected || executing"
          size="small"
        >
          {{ cmd.name }}
        </el-button>
      </div>
    </div>

    <!-- 执行历史 -->
    <div class="command-history" v-if="history.length > 0">
      <h4>执行历史</h4>
      <div class="history-list">
        <div 
          v-for="(item, index) in history" 
          :key="index"
          class="history-item"
          :class="{ 'error': item.error }"
        >
          <div class="history-command">
            <el-icon><Monitor /></el-icon>
            <code>{{ item.command }}</code>
            <span class="history-time">{{ formatTime(item.timestamp) }}</span>
          </div>
          <div class="history-result" v-if="item.result !== undefined">
            <strong>结果:</strong>
            <pre>{{ formatResult(item.result) }}</pre>
          </div>
          <div class="history-error" v-if="item.error">
            <strong>错误:</strong>
            <span class="error-text">{{ item.error }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 连接状态提示 -->
    <el-alert
      v-if="!connected"
      title="未连接到Redis"
      type="warning"
      description="请先连接到Redis服务器才能执行命令"
      show-icon
      :closable="false"
    />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Monitor } from '@element-plus/icons-vue'
import { useRedisStore } from '@/stores/redis'
import axios from 'axios'
import { API_BASE_URL } from '@/config.js'

const redisStore = useRedisStore()

// 响应式数据
const command = ref('')
const executing = ref(false)
const history = ref([])

// 计算属性
const connected = computed(() => redisStore.connected)

// 常用命令
const quickCommands = ref([
  { name: 'LPUSH cut_job 123', command: 'LPUSH cut_job 123' },
  { name: 'LRANGE cut_job 0 -1', command: 'LRANGE cut_job 0 -1' },
  { name: 'LLEN cut_job', command: 'LLEN cut_job' },
  { name: 'KEYS *', command: 'KEYS *' },
  { name: 'INFO', command: 'INFO' },
  { name: 'PING', command: 'PING' },
  { name: 'DBSIZE', command: 'DBSIZE' },
  { name: 'FLUSHDB', command: 'FLUSHDB' }
])

// 执行命令
const executeCommand = async () => {
  if (!command.value.trim() || !connected.value || executing.value) {
    return
  }

  executing.value = true
  
  try {
    const response = await axios.post(`${API_BASE_URL}/api/redis/command`, {
      command: command.value.trim()
    })

    const result = {
      command: command.value.trim(),
      result: response.data.data?.result,
      error: response.data.code !== 0 ? response.data.message : null,
      timestamp: new Date()
    }

    // 添加到历史记录
    history.value.unshift(result)
    
    // 限制历史记录数量
    if (history.value.length > 50) {
      history.value = history.value.slice(0, 50)
    }

    if (response.data.code === 0) {
      ElMessage.success('命令执行成功')
    } else {
      ElMessage.error(`命令执行失败: ${response.data.message}`)
    }

    // 清空输入框
    command.value = ''
    
  } catch (error) {
    const result = {
      command: command.value.trim(),
      result: null,
      error: error.response?.data?.message || error.message,
      timestamp: new Date()
    }
    
    history.value.unshift(result)
    ElMessage.error(`命令执行失败: ${error.message}`)
  } finally {
    executing.value = false
  }
}

// 设置命令
const setCommand = (cmd) => {
  command.value = cmd
}

// 格式化结果
const formatResult = (result) => {
  if (result === null || result === undefined) {
    return '(nil)'
  }
  
  if (typeof result === 'object') {
    return JSON.stringify(result, null, 2)
  }
  
  return String(result)
}

// 格式化时间
const formatTime = (timestamp) => {
  return timestamp.toLocaleTimeString()
}
</script>

<style scoped>
.command-executor {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 20px;
}

.command-header {
  margin-bottom: 20px;
}

.command-header h3 {
  margin: 0 0 8px 0;
  color: var(--el-text-color-primary);
}

.description {
  margin: 0;
  color: var(--el-text-color-regular);
  font-size: 14px;
}

.command-input-section {
  margin-bottom: 20px;
}

.command-input {
  width: 100%;
}

.quick-commands {
  margin-bottom: 20px;
}

.quick-commands h4 {
  margin: 0 0 12px 0;
  color: var(--el-text-color-primary);
  font-size: 14px;
}

.command-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.command-history {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.command-history h4 {
  margin: 0 0 12px 0;
  color: var(--el-text-color-primary);
  font-size: 14px;
}

.history-list {
  flex: 1;
  overflow-y: auto;
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  padding: 12px;
  background: var(--el-bg-color-page);
}

.history-item {
  margin-bottom: 16px;
  padding: 12px;
  border-radius: 4px;
  background: var(--el-bg-color);
  border-left: 3px solid var(--el-color-success);
}

.history-item.error {
  border-left-color: var(--el-color-danger);
}

.history-command {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  font-family: 'Courier New', monospace;
}

.history-command code {
  flex: 1;
  background: var(--el-fill-color-light);
  padding: 4px 8px;
  border-radius: 3px;
  font-size: 13px;
}

.history-time {
  color: var(--el-text-color-placeholder);
  font-size: 12px;
}

.history-result {
  margin-top: 8px;
}

.history-result strong {
  color: var(--el-text-color-primary);
  font-size: 13px;
}

.history-result pre {
  margin: 4px 0 0 0;
  padding: 8px;
  background: var(--el-fill-color-light);
  border-radius: 3px;
  font-size: 12px;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 200px;
  overflow-y: auto;
}

.history-error {
  margin-top: 8px;
}

.history-error strong {
  color: var(--el-color-danger);
  font-size: 13px;
}

.error-text {
  color: var(--el-color-danger);
  font-size: 13px;
}
</style>
