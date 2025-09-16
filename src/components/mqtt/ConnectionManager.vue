<template>
  <el-drawer
    v-model="visible"
    title="MQTT 连接管理"
    size="50%"
  >
    <div class="conn-manager">
      <div class="left">
        <el-form :model="form" label-width="120px">
          <el-form-item label="协议">
            <el-select v-model="form.protocol" placeholder="选择协议" style="width: 100%">
              <el-option label="ws" value="ws" />
              <el-option label="wss" value="wss" />
            </el-select>
          </el-form-item>
          <el-form-item label="主机">
            <el-input v-model="form.host" placeholder="例如 127.0.0.1" />
          </el-form-item>
          <el-form-item label="端口">
            <el-input-number v-model="form.port" :min="1" :max="65535" />
          </el-form-item>
          <el-form-item label="路径">
            <el-input v-model="form.path" placeholder="/mqtt" />
          </el-form-item>
          <el-form-item label="用户名">
            <el-input v-model="form.username" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="form.password" type="password" show-password />
          </el-form-item>
          <el-form-item label="Client ID">
            <el-input v-model="form.clientId" placeholder="留空自动生成" />
          </el-form-item>
          <el-form-item label="clean">
            <el-switch v-model="form.clean" />
          </el-form-item>
          <el-form-item label="keepalive">
            <el-input-number v-model="form.keepalive" :min="0" :max="600" />
          </el-form-item>
          <el-form-item label="重连周期(ms)">
            <el-input-number v-model="form.reconnectPeriod" :min="0" :step="1000" />
          </el-form-item>
          <el-form-item label="超时(ms)">
            <el-input-number v-model="form.connectTimeout" :min="1000" :step="1000" />
          </el-form-item>
          <el-form-item label="证书校验">
            <el-switch v-model="form.rejectUnauthorized" />
          </el-form-item>

          <el-space>
            <el-button type="primary" @click="onQuickConnect" :loading="mqtt.isConnecting">快速连接</el-button>
            <el-button @click="onTest">测试连接</el-button>
            <el-button type="warning" @click="onDisconnect" :disabled="!mqtt.isConnected">断开</el-button>
          </el-space>
        </el-form>
      </div>

      <div class="right">
        <div class="toolbar">
          <el-input v-model="profileName" placeholder="配置名称" style="width: 220px; margin-right: 8px;" />
          <el-checkbox v-model="isDefault">设为默认</el-checkbox>
          <el-button type="success" @click="onSave">保存配置</el-button>
        </div>

        <el-table :data="mqtt.profiles" height="520px" style="width: 100%; margin-top: 10px;">
          <el-table-column prop="name" label="名称" width="200" />
          <el-table-column label="目标">
            <template #default="{ row }">
              {{ row.config.protocol }}://{{ row.config.host }}:{{ row.config.port }}{{ row.config.path || '/mqtt' }}
            </template>
          </el-table-column>
          <el-table-column label="默认" width="80">
            <template #default="{ row }">
              <el-tag v-if="row.isDefault" type="success">默认</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="260">
            <template #default="{ row }">
              <el-button size="small" @click="useProfile(row.id)">使用</el-button>
              <el-button size="small" @click="setDefault(row.id)">设为默认</el-button>
              <el-button size="small" type="danger" @click="remove(row.id)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </el-drawer>
</template>

<script setup>
import { ref, reactive, onMounted, defineExpose } from 'vue'
import { ElMessage } from 'element-plus'
import { useMqttStore } from '../../stores/mqtt'
import { generateClientId } from '../../utils/mqtt'

const mqtt = useMqttStore()
const visible = ref(false)
const profileName = ref('默认配置')
const isDefault = ref(false)

const form = reactive({
  protocol: 'ws',
  host: '',
  port: undefined,
  path: '/mqtt',
  username: '',
  password: '',
  clientId: '',
  clean: true,
  keepalive: 60,
  reconnectPeriod: 0,
  connectTimeout: 30000,
  rejectUnauthorized: false
})

function open() {
  visible.value = true
}

defineExpose({ open })

onMounted(async () => {
  try {
    await mqtt.loadProfiles()
    // 用默认配置填充表单
    if (mqtt.currentConfig) Object.assign(form, mqtt.currentConfig)
  } catch (e) {}
})

async function onQuickConnect() {
  try {
    const cfg = { ...form, clientId: form.clientId || generateClientId('web_client') }
    await mqtt.connectWithConfig(cfg)
    ElMessage.success('连接成功')
  } catch (e) {
    ElMessage.error(`连接失败: ${e.message || e}`)
  }
}

async function onTest() {
  try {
    await mqtt.testConnection({ config: { ...form, clientId: form.clientId || generateClientId('web_client') } })
    ElMessage.success('测试通过')
  } catch (e) {
    ElMessage.error(`测试失败: ${e.message || e}`)
  }
}

function onDisconnect() {
  mqtt.disconnect()
  ElMessage.info('已断开')
}

async function onSave() {
  if (!profileName.value) {
    ElMessage.warning('请输入配置名称')
    return
  }
  try {
    const payload = {
      name: profileName.value,
      isDefault: isDefault.value,
      config: { ...form, clientId: form.clientId || generateClientId('web_client') }
    }
    await mqtt.saveProfile(payload)
    ElMessage.success('保存成功')
  } catch (e) {
    ElMessage.error(`保存失败: ${e.message || e}`)
  }
}

async function useProfile(id) {
  try {
    await mqtt.selectProfile(id)
    Object.assign(form, mqtt.currentConfig)
    ElMessage.success('已选择配置')
  } catch (e) {
    ElMessage.error(e.message || '选择失败')
  }
}

async function setDefault(id) {
  try {
    await mqtt.setDefaultProfile(id)
    ElMessage.success('已设为默认')
  } catch (e) {
    ElMessage.error(e.message || '设置失败')
  }
}

async function remove(id) {
  try {
    await mqtt.deleteProfile(id)
    ElMessage.success('已删除')
  } catch (e) {
    ElMessage.error(e.message || '删除失败')
  }
}
</script>

<style scoped>
.conn-manager {
  display: flex;
  gap: 16px;
}
.left {
  flex: 1;
  min-width: 360px;
}
.right {
  width: 48%;
}
.toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>


