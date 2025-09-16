<template>
  <div class="bar">
    <el-input v-model="local.pattern" placeholder="模式，如 user:*" style="width: 320px;" />
    <el-select v-model="local.type" placeholder="类型" clearable style="width: 160px; margin-left: 8px;">
      <el-option label="全部" value="" />
      <el-option label="string" value="string" />
      <el-option label="hash" value="hash" />
      <el-option label="list" value="list" />
      <el-option label="set" value="set" />
      <el-option label="zset" value="zset" />
    </el-select>
    <el-button type="primary" :loading="loading" style="margin-left: 8px;" @click="onSearch">搜索</el-button>
    <el-button :loading="loading" @click="emit('refresh')">刷新</el-button>
    <el-button @click="emit('reset')">重置</el-button>
  </div>
  </template>

<script setup>
import { reactive, watchEffect } from 'vue'
const props = defineProps({ pattern: String, typeFilter: String, loading: Boolean })
const emit = defineEmits(['search','reset','refresh'])
const local = reactive({ pattern: props.pattern || '*', type: props.typeFilter || '' })
watchEffect(() => { local.pattern = props.pattern; local.type = props.typeFilter })
function onSearch() { emit('search', { pattern: local.pattern, type: local.type }) }
</script>

<style scoped>
.bar { display: flex; align-items: center; }
</style>


