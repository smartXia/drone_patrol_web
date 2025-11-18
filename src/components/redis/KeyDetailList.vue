<template>
  <el-card shadow="never" class="detail">
    <template #header>
      <div class="head">
        <div class="title">{{ keyName }} <el-tag size="small">list</el-tag></div>
        <div class="ops">
          <el-button size="small" @click="toggleEdit" :type="editing ? 'warning' : 'default'">{{ editing ? '完成编辑' : '编辑' }}</el-button>
          <el-button size="small" @click="onRename">重命名</el-button>
          <el-button size="small" @click="onExpire">设置TTL</el-button>
          <el-button size="small" @click="emit('persist', keyName)">移除TTL</el-button>
        </div>
      </div>
    </template>
    <div class="bar">
      <el-input v-model="pushValue" style="width: 60%;" placeholder="值" :disabled="!editing" />
      <el-button size="small" type="primary" @click="lpush" :disabled="!editing">LPUSH</el-button>
      <el-button size="small" type="primary" @click="rpush" :disabled="!editing">RPUSH</el-button>
      <el-button size="small" @click="lpop" :disabled="!editing">LPOP</el-button>
      <el-button size="small" @click="rpop" :disabled="!editing">RPOP</el-button>
    </div>
    <el-table :data="items" height="calc(100% - 80px)">
      <el-table-column type="index" width="70" />
      <el-table-column label="Value">
        <template #default="{ row, $index }">
          <div style="display:flex; gap:8px; align-items:center;">
            <el-input v-model="items[$index]" :disabled="!editing" />
            <el-button size="small" type="primary" @click="lset($index, items[$index])" :disabled="!editing">保存</el-button>
            <el-popconfirm title="删除该元素？" @confirm="lrem(items[$index])">
              <template #reference>
                <el-button size="small" type="danger" :disabled="!editing">删除</el-button>
              </template>
            </el-popconfirm>
          </div>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessageBox } from 'element-plus'
import { lRange, lPush, rPush, lPop, rPop, lSet, lRem } from '@/api/redis'
const props = defineProps({ keyName: String, keyType: String })
const emit = defineEmits(['rename','expire','persist','refresh'])
const items = ref([])
const pushValue = ref('')
const editing = ref(false)

async function load() { items.value = await lRange(props.keyName, 0, 100) }
async function lpush() { await lPush(props.keyName, pushValue.value); pushValue.value = ''; await load(); emit('refresh') }
async function rpush() { await rPush(props.keyName, pushValue.value); pushValue.value = ''; await load(); emit('refresh') }
async function lpop() { await lPop(props.keyName); await load(); emit('refresh') }
async function rpop() { await rPop(props.keyName); await load(); emit('refresh') }
async function lset(index, value) { await lSet(props.keyName, index, value); await load(); emit('refresh') }
async function lrem(value) { await lRem(props.keyName, value); await load(); emit('refresh') }

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
onMounted(load)
function toggleEdit() { editing.value = !editing.value }
</script>

<style scoped>
.detail { height: 100%; }
.head { display: flex; justify-content: space-between; align-items: center; }
.bar { margin-bottom: 8px; display: flex; align-items: center; gap: 8px; }
</style>


