<template>
  <el-card shadow="never" class="detail">
    <template #header>
      <div class="head">
        <div class="title">{{ keyName }} <el-tag size="small">string</el-tag></div>
        <div class="ops">
          <el-button size="small" @click="toggleEdit" :type="editing ? 'warning' : 'default'">{{ editing ? '完成编辑' : '编辑' }}</el-button>
          <el-button size="small" @click="onRename">重命名</el-button>
          <el-button size="small" @click="onExpire">设置TTL</el-button>
          <el-button size="small" @click="emit('persist', keyName)">移除TTL</el-button>
          <el-button size="small" type="primary" @click="save" :disabled="!editing" :loading="saving">保存</el-button>
        </div>
      </div>
    </template>
    <el-input v-model="value" type="textarea" :rows="20" placeholder="值..." :disabled="!editing" />
    <div style="margin-top:8px; text-align:right;">
      <el-popconfirm title="删除该键？" @confirm="emit('delete', keyName)">
        <template #reference>
          <el-button size="small" type="danger" :disabled="!editing">删除键</el-button>
        </template>
      </el-popconfirm>
    </div>
  </el-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessageBox } from 'element-plus'
import { strGet, strSet } from '@/api/redis'
const props = defineProps({ keyName: String, keyType: String })
const emit = defineEmits(['rename','expire','persist','refresh','delete'])
const value = ref('')
const editing = ref(false)
const saving = ref(false)

async function load() { value.value = await strGet(props.keyName) || '' }
async function save() {
  saving.value = true
  try {
    await strSet(props.keyName, value.value)
  } finally { saving.value = false; emit('refresh') }
}
function onRename() {
  ElMessageBox.prompt('输入新的 Key 名称', '重命名', { inputValue: props.keyName })
    .then(({ value: nv }) => emit('rename', { key: props.keyName, newKey: nv }))
    .catch(() => {})
}
function onExpire() {
  ElMessageBox.prompt('设置 TTL（秒），-1 取消', 'TTL', { inputValue: '3600' })
    .then(({ value: ttl }) => emit('expire', { key: props.keyName, ttl: Number(ttl) }))
    .catch(() => {})
}
function toggleEdit() { editing.value = !editing.value }
onMounted(load)
</script>

<style scoped>
.detail { height: 100%; }
.head { display: flex; justify-content: space-between; align-items: center; }
.ops > * { margin-left: 6px; }
</style>


