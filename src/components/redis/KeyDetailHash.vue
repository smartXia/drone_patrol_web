<template>
  <el-card shadow="never" class="detail">
    <template #header>
      <div class="head">
        <div class="title">{{ keyName }} <el-tag size="small">hash</el-tag></div>
        <div class="ops">
          <el-button size="small" @click="toggleEdit" :type="editing ? 'warning' : 'default'">{{ editing ? '完成编辑' : '编辑' }}</el-button>
          <el-button size="small" @click="onRename">重命名</el-button>
          <el-button size="small" @click="onExpire">设置TTL</el-button>
          <el-button size="small" @click="emit('persist', keyName)">移除TTL</el-button>
          <el-button size="small" type="primary" @click="addField" :disabled="!editing">新增字段</el-button>
        </div>
      </div>
    </template>
    <el-table :data="items" height="calc(100% - 10px)">
      <el-table-column prop="field" label="Field" min-width="200">
        <template #default="{ row }">
          <el-input v-model="row.field" placeholder="字段名" :disabled="!editing" />
        </template>
      </el-table-column>
      <el-table-column prop="value" label="Value" min-width="300">
        <template #default="{ row }">
          <el-input v-model="row.value" :disabled="!editing" />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button size="small" type="primary" @click="saveField(row)" :disabled="!editing">保存</el-button>
          <el-popconfirm title="删除该字段？" @confirm="removeField(row)">
            <template #reference>
              <el-button size="small" type="danger" :disabled="!editing">删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessageBox } from 'element-plus'
import { hGetAllPaged, hSet, hDel } from '@/api/redis'
const props = defineProps({ keyName: String, keyType: String })
const emit = defineEmits(['rename','expire','persist','refresh'])
const items = ref([])
const editing = ref(false)

async function load() {
  const res = await hGetAllPaged(props.keyName, 0, 200)
  items.value = Object.entries(res).map(([field, value]) => ({ field, value }))
}
function addField() {
  items.value.unshift({ field: '', value: '' })
}
async function saveField(row) {
  await hSet(props.keyName, row.field, row.value)
  emit('refresh')
}
async function removeField(row) {
  await hDel(props.keyName, row.field)
  items.value = items.value.filter(i => !(i.field === row.field))
  emit('refresh')
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


