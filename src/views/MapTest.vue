<template>
  <div class="map-test">
    <h2>高德地图API测试页面</h2>
    
    <div class="test-section">
      <h3>API状态检查</h3>
      <div class="status-info">
        <p><strong>AMap对象:</strong> {{ amapStatus.hasAMap ? '✅ 已加载' : '❌ 未加载' }}</p>
        <p><strong>窗口标志:</strong> {{ amapStatus.hasWindowFlag ? '✅ 已设置' : '❌ 未设置' }}</p>
        <p><strong>API密钥:</strong> {{ amapStatus.key }}</p>
        <p><strong>API版本:</strong> {{ amapStatus.version }}</p>
      </div>
      
      <el-button @click="checkStatus" type="primary">刷新状态</el-button>
      <el-button @click="testMap" type="success" :disabled="!amapStatus.loaded">测试地图</el-button>
    </div>
    
    <div class="test-section">
      <h3>网络请求测试</h3>
      <el-button @click="testNetwork" type="info">测试网络连接</el-button>
      <div v-if="networkResult">
        <p><strong>网络状态:</strong> {{ networkResult.status }}</p>
        <p><strong>响应时间:</strong> {{ networkResult.responseTime }}ms</p>
      </div>
    </div>
    
    <div class="test-section">
      <h3>地图容器</h3>
      <div id="test-map" class="test-map"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getAMapStatus, waitForAMap, AMAP_CONFIG } from '@/config/amap'

const amapStatus = ref({})
const networkResult = ref(null)

// 检查API状态
const checkStatus = () => {
  amapStatus.value = getAMapStatus()
  console.log('📊 API状态:', amapStatus.value)
}

// 测试网络连接
const testNetwork = async () => {
  try {
    const startTime = Date.now()
    const response = await fetch(`https://webapi.amap.com/maps?v=2.0&key=${AMAP_CONFIG.key}`)
    const endTime = Date.now()
    
    networkResult.value = {
      status: response.ok ? '成功' : '失败',
      responseTime: endTime - startTime
    }
    
    ElMessage.success('网络连接正常')
  } catch (error) {
    networkResult.value = {
      status: '失败',
      responseTime: 0
    }
    ElMessage.error('网络连接失败: ' + error.message)
  }
}

// 测试地图初始化
const testMap = async () => {
  try {
    const AMap = await waitForAMap()
    console.log('✅ 开始测试地图初始化...')
    
    const map = new AMap.Map('test-map', {
      zoom: 15,
      center: [118.654765, 31.971145]
    })
    
    ElMessage.success('地图初始化成功')
  } catch (error) {
    console.error('❌ 地图测试失败:', error)
    ElMessage.error('地图测试失败: ' + error.message)
  }
}

onMounted(() => {
  checkStatus()
  
  // 定期检查状态
  setInterval(() => {
    checkStatus()
  }, 2000)
})
</script>

<style scoped>
.map-test {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.test-section {
  margin: 20px 0;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
}

.status-info {
  margin: 10px 0;
}

.status-info p {
  margin: 5px 0;
}

.test-map {
  width: 100%;
  height: 400px;
  border: 1px solid #ccc;
  border-radius: 4px;
}
</style>
