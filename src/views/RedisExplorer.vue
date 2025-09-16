<template>
  <el-container style="height: 100%;">
    <el-aside width="320px" class="redis-aside">
      <connection-list
        :profiles="profiles"
        :currentProfileId="currentProfileId"
        :connecting="connecting"
        :connected="connected"
        @save="handleSaveProfile"
        @delete="handleDeleteProfile"
        @connect="handleConnect"
        @disconnect="handleDisconnect"
        @select="handleSelectProfile"
      />
    </el-aside>
    <el-container>
      <el-header class="redis-header">
        <key-search-bar
          :pattern="pattern"
          :typeFilter="typeFilter"
          :loading="scanning"
          @search="handleSearch"
          @refresh="handleRefresh"
          @reset="handleReset"
        />
        <div style="margin-left:auto; display:flex; align-items:center; gap:8px;">
          <el-button @click="gotoDashboard">返回仪表盘</el-button>
          <el-button type="primary" @click="gotoDashboardOpenMqtt">仪表盘打开MQTT</el-button>
        </div>
      </el-header>
      <el-main class="redis-main">
        <el-row :gutter="12" style="height: 100%;">
          <el-col :span="6" style="height: 100%;">
            <key-tree
              :items="keys"
              :hasMore="hasMore"
              :loading="scanning"
              @loadMore="handleLoadMore"
              @select="handleSelectKey"
            />
          </el-col>
          <el-col :span="18" style="height: 100%;">
            <component
              :is="currentDetailComponent"
              v-if="selectedKey"
              :key="selectedKey + ':' + selectedType"
              :keyName="selectedKey"
              :keyType="selectedType"
              @rename="handleRenameKey"
              @expire="handleExpireKey"
              @persist="handlePersistKey"
              @delete="handleDeleteKey"
              @refresh="handleRefreshKey"
            />
            <el-empty v-else description="请选择一个键以查看详情" />
          </el-col>
        </el-row>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRedisStore } from '@/stores/redis'
import { useRouter } from 'vue-router'
import ConnectionList from '@/components/redis/ConnectionList.vue'
import KeySearchBar from '@/components/redis/KeySearchBar.vue'
import KeyList from '@/components/redis/KeyList.vue'
import KeyTree from '@/components/redis/KeyTree.vue'
import KeyDetailString from '@/components/redis/KeyDetailString.vue'
import KeyDetailHash from '@/components/redis/KeyDetailHash.vue'
import KeyDetailList from '@/components/redis/KeyDetailList.vue'
import KeyDetailSet from '@/components/redis/KeyDetailSet.vue'
import KeyDetailZSet from '@/components/redis/KeyDetailZSet.vue'

const store = useRedisStore()
const router = useRouter()

const profiles = computed(() => store.profiles)
const currentProfileId = computed(() => store.currentProfileId)
const connecting = computed(() => store.connecting)
const connected = computed(() => store.connected)

const pattern = computed(() => store.pattern)
const typeFilter = computed(() => store.typeFilter)
const scanning = computed(() => store.scanning)
const keys = computed(() => store.keys)
const hasMore = computed(() => store.hasMore)

const selectedKey = computed(() => store.selectedKey)
const selectedType = computed(() => store.selectedType)

const currentDetailComponent = computed(() => {
  switch (selectedType.value) {
    case 'string':
      return KeyDetailString
    case 'hash':
      return KeyDetailHash
    case 'list':
      return KeyDetailList
    case 'set':
      return KeyDetailSet
    case 'zset':
      return KeyDetailZSet
    default:
      return null
  }
})

function handleSaveProfile(payload) { store.saveProfile(payload) }
function handleDeleteProfile(id) { store.deleteProfile(id) }
function handleConnect(id) { store.connect(id) }
function handleDisconnect() { store.disconnect() }
function handleSelectProfile(id) { store.selectProfile(id) }

function handleSearch(params) { store.search(params) }
function handleRefresh() { store.search({ pattern: pattern.value, type: typeFilter.value }) }
function handleReset() { store.resetSearch() }
function handleLoadMore() { store.loadMore() }

function handleSelectKey(item) { store.selectKey(item) }
function handleDeleteKey(key) { store.deleteKey(key) }
function handleRenameKey({ key, newKey }) { store.renameKey(key, newKey) }
function handleExpireKey({ key, ttl }) { store.expireKey(key, ttl) }
function handlePersistKey(key) { store.persistKey(key) }
function handleRefreshKey() { store.refreshSelected() }

function gotoDashboard() { router.push({ name: 'Dashboard' }) }
function gotoDashboardOpenMqtt() { router.push({ name: 'Dashboard', query: { openMqtt: 1 } }) }

onMounted(() => {
  store.init()
})
</script>

<style scoped>
.redis-aside { border-right: 1px solid var(--el-border-color); padding: 8px; background: var(--el-fill-color-blank); }
.redis-header { border-bottom: 1px solid var(--el-border-color); padding: 8px 12px; background: var(--el-fill-color-blank); }
.redis-main { padding: 8px; background: var(--el-fill-color-lighter); }
</style>


