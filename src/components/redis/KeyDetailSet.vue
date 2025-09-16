<template>
  <el-card shadow="never" class="detail">
    <template #header>
      <div class="head">
        <div class="title">{{ keyName }} <el-tag size="small">set</el-tag></div>
        <div class="ops">
          <el-button size="small" @click="toggleEdit" :type="editing ? 'warning' : 'default'">{{ editing ? '完成编辑' : '编辑' }}</el-button>
          <el-button size="small" @click="onRename">重命名</el-button>
          <el-button size="small" @click="onExpire">设置TTL</el-button>
          <el-button size="small" @click="emit('persist', keyName)">移除TTL</el-button>
        </div>
      </div>
    </template>
    <div class="bar">
      <el-input v-model="member" style="width: 60%;" placeholder="成员" :disabled="!editing" />
      <el-button size="small" type="primary" @click="sadd" :disabled="!editing">SADD</el-button>
    </div>
    <el-table :data="items" height="calc(100% - 80px)">
      <el-table-column prop="member" label="Member" />
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-popconfirm title="移除该成员？" @confirm="srem(row)">
            <template #reference>
              <el-button size="small" type="danger" :disabled="!editing">移除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { sScan, sAdd, sRem } from '@/api/redis'
const props = defineProps({ keyName: String, keyType: String })
const emit = defineEmits(['rename','expire','persist','refresh'])
const items = ref([])
const member = ref('')
const editing = ref(false)

const sCursor = ref('0')
async function load() {
  const res = await sScan(props.keyName, sCursor.value, undefined, 200)
  sCursor.value = res.cursor
  items.value = res.items || []
}
async function sadd() { await sAdd(props.keyName, member.value); member.value=''; await load(); emit('refresh') }
async function srem(row) { await sRem(props.keyName, row.member); await load(); emit('refresh') }

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


