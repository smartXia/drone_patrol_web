<template>
  <div id="app">
    <el-container class="app-container">
      <el-header class="app-header">
        <div class="header-content">
          <div class="header-left">
            <h1>无人机设备监控系统</h1>
            <el-menu
              :default-active="$route.path"
              mode="horizontal"
              background-color="transparent"
              text-color="rgba(255, 255, 255, 0.8)"
              active-text-color="#ffffff"
              @select="handleMenuSelect"
              class="nav-menu"
              style="display: flex; gap: 0;"
            >
              <el-menu-item index="/">
                <el-icon><Monitor /></el-icon>
                <span>设备监控</span>
              </el-menu-item>
              <el-menu-item index="/redis">
                <el-icon><DataBoard /></el-icon>
                <span>Redis 管理</span>
              </el-menu-item>
              <el-menu-item index="/error-codes">
                <el-icon><Document /></el-icon>
                <span>错误码查询</span>
              </el-menu-item>
              <el-menu-item index="/live-stream">
                <el-icon><VideoPlay /></el-icon>
                <span>直播流</span>
              </el-menu-item>
            </el-menu>
          </div>
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
import { useRouter } from 'vue-router'
import { useMqttStore } from '@/stores/mqtt'
import { Monitor, DataBoard, Document, VideoPlay } from '@element-plus/icons-vue'

const router = useRouter()
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

const handleMenuSelect = (index) => {
  console.log('菜单选择:', index)
  router.push(index)
}
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

.header-left {
  display: flex;
  align-items: center;
  gap: 32px;
}

.nav-menu {
  border-bottom: none;
}

.nav-menu .el-menu-item {
  border-bottom: 2px solid transparent;
  transition: all 0.3s;
  display: flex !important;
  align-items: center;
  gap: 8px;
}

.nav-menu .el-menu-item:hover,
.nav-menu .el-menu-item.is-active {
  background-color: rgba(255, 255, 255, 0.1);
  border-bottom-color: #ffffff;
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
