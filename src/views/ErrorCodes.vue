<template>
  <div class="error-codes-container">
    <div class="header">
      <h1>DJI 错误码查询</h1>
      <div class="search-section">
        <el-input
          v-model="searchQuery"
          placeholder="搜索错误码或描述..."
          clearable
          @input="handleSearch"
          class="search-input"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-select
          v-model="selectedCategory"
          placeholder="选择分类"
          clearable
          @change="handleCategoryChange"
          class="category-select"
        >
          <el-option
            v-for="category in categories"
            :key="category.value"
            :label="category.label"
            :value="category.value"
          />
        </el-select>
      </div>
    </div>

    <div class="stats">
      <el-card class="stat-card">
        <div class="stat-item">
          <span class="stat-label">总错误码数</span>
          <span class="stat-value">{{ totalCount }}</span>
        </div>
      </el-card>
      <el-card class="stat-card">
        <div class="stat-item">
          <span class="stat-label">搜索结果</span>
          <span class="stat-value">{{ filteredCodes.length }}</span>
        </div>
      </el-card>
    </div>

    <div class="content">
      <el-table
        :data="paginatedCodes"
        stripe
        border
        height="600"
        :row-class-name="getRowClassName"
        @row-click="handleRowClick"
        class="error-table"
      >
        <el-table-column prop="code" label="错误码" width="120" fixed="left">
          <template #default="{ row }">
            <el-tag :type="getCodeType(row.code)" class="code-tag">
              {{ row.code }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="错误描述" min-width="400">
          <template #default="{ row }">
            <div class="description-cell">
              <span v-html="highlightText(row.description)"></span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="category" label="分类" width="120">
          <template #default="{ row }">
            <el-tag :type="getCategoryType(row.category)" size="small">
              {{ getCategoryLabel(row.category) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="severity" label="严重程度" width="100">
          <template #default="{ row }">
            <el-tag :type="getSeverityType(row.severity)" size="small">
              {{ getSeverityLabel(row.severity) }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[20, 50, 100, 200]"
          :total="filteredCodes.length"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="`错误码 ${selectedCode?.code} 详情`"
      width="60%"
      class="detail-dialog"
    >
      <div v-if="selectedCode" class="code-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="错误码">
            <el-tag :type="getCodeType(selectedCode.code)" size="large">
              {{ selectedCode.code }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="分类">
            <el-tag :type="getCategoryType(selectedCode.category)">
              {{ getCategoryLabel(selectedCode.category) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="严重程度">
            <el-tag :type="getSeverityType(selectedCode.severity)">
              {{ getSeverityLabel(selectedCode.severity) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="建议操作">
            <span>{{ getSuggestedAction(selectedCode.code) }}</span>
          </el-descriptions-item>
        </el-descriptions>
        
        <div class="description-section">
          <h4>详细描述</h4>
          <p>{{ selectedCode.description }}</p>
        </div>

        <div class="suggestions-section">
          <h4>解决建议</h4>
          <ul>
            <li v-for="suggestion in getSuggestions(selectedCode.code)" :key="suggestion">
              {{ suggestion }}
            </li>
          </ul>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import axios from 'axios'

// 响应式数据
const searchQuery = ref('')
const selectedCategory = ref('')
const currentPage = ref(1)
const pageSize = ref(50)
const detailDialogVisible = ref(false)
const selectedCode = ref(null)
const errorCodes = ref([])

// 分类定义
const categories = ref([
  { value: 'upgrade', label: '设备升级' },
  { value: 'flight', label: '飞行任务' },
  { value: 'communication', label: '通信异常' },
  { value: 'battery', label: '电池相关' },
  { value: 'rtk', label: 'RTK定位' },
  { value: 'camera', label: '相机相关' },
  { value: 'media', label: '媒体文件' },
  { value: 'system', label: '系统异常' },
  { value: 'network', label: '网络相关' },
  { value: 'weather', label: '天气环境' },
  { value: 'other', label: '其他' }
])

// 计算属性
const totalCount = computed(() => errorCodes.value.length)

const filteredCodes = computed(() => {
  let filtered = errorCodes.value

  // 按分类过滤
  if (selectedCategory.value) {
    filtered = filtered.filter(code => code.category === selectedCategory.value)
  }

  // 按搜索关键词过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(code => 
      code.code.includes(query) || 
      code.description.toLowerCase().includes(query)
    )
  }

  return filtered
})

const paginatedCodes = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredCodes.value.slice(start, end)
})

// 方法
const loadErrorCodes = async () => {
  try {
    const response = await axios.get('http://127.0.0.1:8080/api/error-codes')
    if (response.data.code === 0) {
      errorCodes.value = response.data.data
    } else {
      throw new Error(response.data.message)
    }
  } catch (error) {
    console.error('加载错误码失败:', error)
    // 使用示例数据
    errorCodes.value = [
      {
        code: '312014',
        description: '设备升级中，请勿重复操作',
        category: 'upgrade',
        severity: 'warning'
      },
      {
        code: '314001',
        description: '飞行任务下发失败，请稍后重试',
        category: 'flight',
        severity: 'error'
      }
    ]
  }
}

const getCodeCategory = (code) => {
  const codeNum = parseInt(code)
  if (codeNum >= 312000 && codeNum < 315000) return 'upgrade'
  if (codeNum >= 314000 && codeNum < 317000) return 'flight'
  if (codeNum >= 315000 && codeNum < 316000) return 'communication'
  if (codeNum >= 316000 && codeNum < 317000) return 'battery'
  if (codeNum >= 317000 && codeNum < 319000) return 'media'
  if (codeNum >= 319000 && codeNum < 322000) return 'system'
  if (codeNum >= 321000 && codeNum < 325000) return 'flight'
  if (codeNum >= 325000 && codeNum < 327000) return 'network'
  if (codeNum >= 327000 && codeNum < 329000) return 'camera'
  if (codeNum >= 336000 && codeNum < 339000) return 'flight'
  if (codeNum >= 338000 && codeNum < 339000) return 'flight'
  if (codeNum >= 514000 && codeNum < 515000) return 'weather'
  return 'other'
}

const getCodeSeverity = (code, description) => {
  if (description.includes('失败') || description.includes('异常') || description.includes('错误')) {
    return 'error'
  }
  if (description.includes('超时') || description.includes('无法') || description.includes('请重试')) {
    return 'warning'
  }
  if (description.includes('成功') || description.includes('完成')) {
    return 'success'
  }
  return 'info'
}

const getCategoryLabel = (category) => {
  const cat = categories.value.find(c => c.value === category)
  return cat ? cat.label : '其他'
}

const getSeverityLabel = (severity) => {
  const labels = {
    error: '严重',
    warning: '警告',
    success: '正常',
    info: '信息'
  }
  return labels[severity] || '未知'
}

const getCodeType = (code) => {
  const codeNum = parseInt(code)
  if (codeNum >= 312000 && codeNum < 315000) return 'primary'
  if (codeNum >= 314000 && codeNum < 317000) return 'danger'
  if (codeNum >= 315000 && codeNum < 316000) return 'warning'
  if (codeNum >= 316000 && codeNum < 317000) return 'success'
  if (codeNum >= 317000 && codeNum < 319000) return 'info'
  return 'default'
}

const getCategoryType = (category) => {
  const types = {
    upgrade: 'primary',
    flight: 'danger',
    communication: 'warning',
    battery: 'success',
    rtk: 'info',
    camera: 'primary',
    media: 'info',
    system: 'danger',
    network: 'warning',
    weather: 'success',
    other: 'default'
  }
  return types[category] || 'default'
}

const getSeverityType = (severity) => {
  const types = {
    error: 'danger',
    warning: 'warning',
    success: 'success',
    info: 'info'
  }
  return types[severity] || 'default'
}

const getRowClassName = ({ row }) => {
  return `severity-${row.severity}`
}

const highlightText = (text) => {
  if (!searchQuery.value) return text
  const regex = new RegExp(`(${searchQuery.value})`, 'gi')
  return text.replace(regex, '<mark>$1</mark>')
}

const handleSearch = () => {
  currentPage.value = 1
}

const handleCategoryChange = () => {
  currentPage.value = 1
}

const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1
}

const handleCurrentChange = (val) => {
  currentPage.value = val
}

const handleRowClick = (row) => {
  selectedCode.value = row
  detailDialogVisible.value = true
}

const getSuggestedAction = (code) => {
  const actions = {
    '312014': '等待升级完成',
    '314001': '检查网络连接后重试',
    '315000': '重启机场设备',
    '316009': '给飞行器充电',
    '319002': '重启机场系统'
  }
  return actions[code] || '请根据错误描述采取相应措施'
}

const getSuggestions = (code) => {
  const suggestions = {
    '312014': [
      '等待当前升级任务完成',
      '检查设备状态指示灯',
      '确保网络连接稳定'
    ],
    '314001': [
      '检查网络连接状态',
      '确认机场设备在线',
      '稍后重新下发任务'
    ],
    '315000': [
      '重启机场设备',
      '检查电源连接',
      '联系技术支持'
    ],
    '316009': [
      '给飞行器充电至50%以上',
      '检查电池状态',
      '等待充电完成'
    ]
  }
  return suggestions[code] || [
    '根据错误描述检查相关设备',
    '重启相关设备后重试',
    '如问题持续请联系技术支持'
  ]
}

// 生命周期
onMounted(() => {
  loadErrorCodes()
})
</script>

<style scoped>
.error-codes-container {
  padding: 20px;
  background-color: #f5f5f5;
  min-height: 100vh;
}

.header {
  background: white;
  padding: 24px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.header h1 {
  margin: 0 0 20px 0;
  color: #303133;
  font-size: 28px;
  font-weight: 600;
}

.search-section {
  display: flex;
  gap: 16px;
  align-items: center;
}

.search-input {
  flex: 1;
  max-width: 400px;
}

.category-select {
  width: 200px;
}

.stats {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.stat-card {
  flex: 1;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
}

.stat-label {
  color: #606266;
  font-size: 14px;
}

.stat-value {
  color: #409eff;
  font-size: 24px;
  font-weight: 600;
}

.content {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.error-table {
  width: 100%;
}

.code-tag {
  font-weight: 600;
  font-family: 'Courier New', monospace;
}

.description-cell {
  line-height: 1.6;
}

.description-cell mark {
  background-color: #fff3cd;
  padding: 2px 4px;
  border-radius: 3px;
}

.severity-error {
  background-color: #fef0f0;
}

.severity-warning {
  background-color: #fdf6ec;
}

.severity-success {
  background-color: #f0f9ff;
}

.severity-info {
  background-color: #f4f4f5;
}

.pagination {
  padding: 16px;
  text-align: center;
  background: #fafafa;
}

.detail-dialog .code-detail {
  padding: 16px 0;
}

.description-section,
.suggestions-section {
  margin-top: 24px;
}

.description-section h4,
.suggestions-section h4 {
  margin: 0 0 12px 0;
  color: #303133;
  font-size: 16px;
  font-weight: 600;
}

.description-section p {
  margin: 0;
  line-height: 1.6;
  color: #606266;
}

.suggestions-section ul {
  margin: 0;
  padding-left: 20px;
}

.suggestions-section li {
  margin-bottom: 8px;
  line-height: 1.6;
  color: #606266;
}

:deep(.el-table__row) {
  cursor: pointer;
}

:deep(.el-table__row:hover) {
  background-color: #f5f7fa;
}

:deep(.el-table__row.severity-error:hover) {
  background-color: #fef0f0;
}

:deep(.el-table__row.severity-warning:hover) {
  background-color: #fdf6ec;
}

:deep(.el-table__row.severity-success:hover) {
  background-color: #f0f9ff;
}

:deep(.el-table__row.severity-info:hover) {
  background-color: #f4f4f5;
}
</style>
