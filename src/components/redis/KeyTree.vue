<template>
  <div class="key-tree">
    <el-scrollbar height="100%" ref="scrollbarRef" @scroll="onScroll">
      <el-tree
      :data="treeData"
      :props="defaultProps"
      node-key="id"
      highlight-current
      :default-expanded-keys="[]"
      :expand-on-click-node="true"
      @node-click="onNodeClick"
    >
      <template #default="{ data }">
        <span class="node-label">
          <el-icon v-if="data.isFolder" style="margin-right:6px;"><Folder /></el-icon>
          <el-icon v-else style="margin-right:6px;"><Document /></el-icon>
          <span class="text" :title="data.fullKey || data.label">{{ data.label }}</span>
          <el-tag v-if="!data.isFolder && data.type" size="small" style="margin-left:8px;">{{ data.type }}</el-tag>
          <span v-if="!data.isFolder && (data.ttl !== undefined)" class="ttl">TTL: {{ data.ttl }}</span>
        </span>
      </template>
      </el-tree>
    </el-scrollbar>
    <div class="more">
      <el-button v-if="hasMore" :loading="loading" @click="$emit('loadMore')">加载更多</el-button>
      <el-text v-else type="info">没有更多了</el-text>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { Folder, Document } from '@element-plus/icons-vue'

const props = defineProps({
  items: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
  hasMore: { type: Boolean, default: false },
  delimiter: { type: String, default: ':' },
})
const emit = defineEmits(['select','loadMore'])

const defaultProps = { label: 'label', children: 'children' }

function buildTree(keys, delimiter) {
  const root = {}
  for (const item of keys) {
    const parts = String(item.key).split(delimiter)
    let node = root
    for (let i = 0; i < parts.length; i++) {
      const part = parts[i]
      node.children = node.children || {}
      if (!node.children[part]) node.children[part] = { __meta: null, children: {} }
      node = node.children[part]
      if (i === parts.length - 1) {
        node.__meta = item
      }
    }
  }
  function toArray(map, path = []) {
    const arr = []
    for (const name of Object.keys(map)) {
      const entry = map[name]
      const isFolder = Object.keys(entry.children || {}).length > 0 && !entry.__meta
      const id = [...path, name].join(':')
      const node = {
        id,
        label: name,
        isFolder,
        children: [],
      }
      if (entry.__meta) {
        node.isFolder = false
        node.type = entry.__meta.type
        node.ttl = entry.__meta.ttl
        node.fullKey = entry.__meta.key
      }
      const childrenArr = toArray(entry.children || {}, [...path, name])
      if (childrenArr.length) node.children = childrenArr
      arr.push(node)
    }
    // 排序：文件夹在前，名称升序
    arr.sort((a, b) => (a.isFolder === b.isFolder ? a.label.localeCompare(b.label) : a.isFolder ? -1 : 1))
    return arr
  }
  return toArray(root.children || {})
}

const treeData = computed(() => buildTree(props.items || [], props.delimiter))
const scrollbarRef = ref()

function onNodeClick(node) {
  if (!node.isFolder && node.fullKey) {
    emit('select', { key: node.fullKey, type: node.type, ttl: node.ttl })
  }
}

function onScroll() {
  const wrap = scrollbarRef.value?.wrapRef
  if (!wrap) return
  const nearBottom = wrap.scrollTop + wrap.clientHeight >= wrap.scrollHeight - 20
  if (nearBottom && !props.loading && props.hasMore) emit('loadMore')
}
</script>

<style scoped>
.key-tree { height: 100%; display: flex; flex-direction: column; background: var(--el-fill-color-blank); border: 1px solid var(--el-border-color-lighter); border-radius: 8px; overflow: hidden; }
.node-label { display: inline-flex; align-items: center; }
.node-label .text { font-size: var(--font-size-small); max-width: 360px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.node-label .ttl { margin-left: 8px; color: var(--el-text-color-secondary); font-size: 12px; }
.more { padding: 8px; text-align: center; border-top: 1px solid var(--el-border-color); }
</style>


