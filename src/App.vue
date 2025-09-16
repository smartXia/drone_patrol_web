<template>
  <div id="app">
    <el-container class="app-container">
      <el-header class="app-header">
        <div class="header-content">
          <h1>无人机设备监控系统</h1>
          <div class="connection-status">
            <el-tag :type="connectionStatus.type" size="large">
              {{ connectionStatus.text }}
            </el-tag>
          </div>
        </div>
      </el-header>
      
      <el-main class="app-main">
        <router-view />
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useMqttStore } from '@/stores/mqtt'

const mqttStore = useMqttStore()

const connectionStatus = computed(() => {
  if (mqttStore.isConnected) {
    return { type: 'success', text: '已连接' }
  } else if (mqttStore.isConnecting) {
    return { type: 'warning', text: '连接中...' }
  } else {
    return { type: 'danger', text: '未连接' }
  }
})
</script>

<style scoped>
.app-container {
  height: 100vh;
}

.app-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 0 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}

.header-content h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
}

.connection-status {
  display: flex;
  align-items: center;
}

.app-main {
  background-color: #f5f7fa;
  padding: 20px;
}
</style>
