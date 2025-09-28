<template>
  <div class="conn">
    <div class="conn-head">
      <span>连接管理</span>
      <el-button type="primary" size="small" @click="openEdit()">新建</el-button>
    </div>
    <el-scrollbar height="calc(100% - 40px)">
      <div class="conn-list">
        <div
          v-for="p in profiles"
          :key="p.id"
          class="conn-item"
          :class="{ active: p.id === currentProfileId }"
          @click="emit('select', p.id)"
        >
          <div class="conn-item-main">
            <div class="title">{{ p.name }} <small>{{ p.host }}:{{ p.port }}</small></div>
            <div class="ops">
              <el-button v-if="connected && p.id === currentProfileId" size="small" @click.stop="emit('disconnect')" :loading="connecting">断开</el-button>
              <el-button v-else size="small" type="success" @click.stop="emit('connect', p.id)" :loading="connecting">连接</el-button>
              <el-button size="small" @click.stop="openEdit(p)">编辑</el-button>
              <el-popconfirm title="删除该连接？" @confirm="emit('delete', p.id)">
                <template #reference>
                  <el-button size="small" type="danger">删除</el-button>
                </template>
              </el-popconfirm>
            </div>
          </div>
        </div>
      </div>
    </el-scrollbar>

    <el-dialog v-model="editVisible" title="连接配置" width="520">
      <el-form :model="form" label-width="90px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="Host"><el-input v-model="form.host" placeholder="127.0.0.1" /></el-form-item>
        <el-form-item label="Port"><el-input v-model.number="form.port" placeholder="6379" /></el-form-item>
        <el-form-item label="DB"><el-input v-model.number="form.db" placeholder="0" /></el-form-item>
        <el-form-item label="密码"><el-input v-model="form.password" type="password" show-password /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible=false">取消</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  profiles: { type: Array, default: () => [] },
  currentProfileId: { type: String, default: '' },
  connecting: { type: Boolean, default: false },
  connected: { type: Boolean, default: false },
})
const emit = defineEmits(['save','delete','connect','disconnect','select'])

const editVisible = ref(false)
const form = ref({ id: '', name: '', host: '127.0.0.1', port: 6379, db: 0, password: '' })

function openEdit(p) {
  if (p) form.value = { ...p }
  else form.value = { id: '', name: '', host: '127.0.0.1', port: 6379, db: 0, password: '' }
  editVisible.value = true
}

function save() {
  if (!form.value.name || !form.value.host || !form.value.port) return
  emit('save', { ...form.value })
  editVisible.value = false
}
</script>

<style scoped>
.conn { height: 100%; display: flex; flex-direction: column; }
.conn-head { display: flex; justify-content: space-between; align-items: center; padding: 8px 4px; }
.conn-list { padding: 4px; }
.conn-item { border: 1px solid var(--el-border-color); border-radius: 6px; padding: 8px; margin-bottom: 8px; cursor: pointer; }
.conn-item.active { border-color: var(--el-color-primary); background: var(--el-color-primary-light-9); }
.conn-item-main { display: flex; justify-content: space-between; align-items: center; }
.title small { color: var(--el-text-color-secondary); margin-left: 8px; }
.ops > * { margin-left: 6px; }
</style>


