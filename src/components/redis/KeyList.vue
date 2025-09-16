<template>
  <div class="keys">
    <el-table :data="items" height="100%" size="small" border stripe @row-click="row => emit('select', row)">
      <el-table-column prop="key" label="Key" min-width="320" show-overflow-tooltip />
      <el-table-column prop="type" label="类型" width="100" />
      <el-table-column prop="ttl" label="TTL" width="100" />
      <el-table-column label="操作" width="180">
        <template #default="{ row }">
          <el-button size="small" @click.stop="emit('select', row)">打开</el-button>
          <el-popconfirm title="删除该键？" @confirm="emit('delete', row.key)">
            <template #reference>
              <el-button size="small" type="danger">删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
    <div class="more">
      <el-button v-if="hasMore" :loading="loading" size="small" @click="emit('loadMore')">加载更多</el-button>
      <el-text v-else type="info">没有更多了</el-text>
    </div>
  </div>
  </template>

<script setup>
const props = defineProps({ items: Array, hasMore: Boolean, loading: Boolean })
const emit = defineEmits(['loadMore','select','delete'])
</script>

<style scoped>
.keys { height: 100%; display: flex; flex-direction: column; background: var(--el-fill-color-blank); border: 1px solid var(--el-border-color-lighter); border-radius: 8px; overflow: hidden; }
.more { padding: 8px; text-align: center; }
</style>


