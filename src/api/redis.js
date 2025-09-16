import axios from 'axios'

const client = axios.create({ baseURL: 'http://127.0.0.1:8001', timeout: 15000 })

let activeConn = { host: '127.0.0.1', port: 6379, db: 0, password: '' }

export function setConnection(conn) {
  activeConn = { host: conn.host, port: conn.port, db: conn.db || 0, password: conn.password || '' }
}

function withConn(body = {}) {
  return { ...activeConn, ...body }
}

export async function testConnect(conn) {
  const { data } = await client.post('/connect/test', conn)
  return data
}

export async function scan({ cursor = '0', match, count = 100, type }) {
  const { data } = await client.post('/scan', withConn({ cursor, match, count, type }))
  return data
}

export async function keyType(key) {
  const { data } = await client.post('/type', withConn({ key }))
  return data.type
}

export async function keyTTL(key) {
  const { data } = await client.post('/ttl', withConn({ key }))
  return data.ttl
}

export async function keyExpire(key, ttl) {
  const { data } = await client.post('/expire', withConn({ key, ttl }))
  return data.ok
}

export async function keyPersist(key) {
  const { data } = await client.post('/persist', withConn({ key }))
  return data.ok
}

export async function keyRename(key, newKey) {
  const { data } = await client.post('/rename', withConn({ key, newKey }))
  return data.ok
}

export async function keyDelete(key) {
  const { data } = await client.post('/del', withConn({ key }))
  return data.ok
}

// string
export async function strGet(key) {
  const { data } = await client.post('/get', withConn({ key }))
  return data.value
}
export async function strSet(key, value) {
  const { data } = await client.post('/set', withConn({ key, value }))
  return data.ok
}

// hash
export async function hGetAllPaged(key, offset = 0, limit = 100) {
  const { data } = await client.post('/hash/getall', withConn({ key, offset, limit }))
  return data
}
export async function hSet(key, field, value) {
  const { data } = await client.post('/hash/set', withConn({ key, field, value }))
  return data.ok
}
export async function hDel(key, field) {
  const { data } = await client.post('/hash/del', withConn({ key, field }))
  return data.ok
}

// list
export async function lRange(key, start = 0, stop = 100) {
  const { data } = await client.post('/list/range', withConn({ key, start, stop }))
  return data.items
}
export async function lPush(key, value) {
  const { data } = await client.post('/list/lpush', withConn({ key, value }))
  return data.ok
}
export async function rPush(key, value) {
  const { data } = await client.post('/list/rpush', withConn({ key, value }))
  return data.ok
}
export async function lPop(key) {
  const { data } = await client.post('/list/lpop', withConn({ key }))
  return data.item
}
export async function rPop(key) {
  const { data } = await client.post('/list/rpop', withConn({ key }))
  return data.item
}
export async function lSet(key, index, value) {
  const { data } = await client.post('/list/set', withConn({ key, index, value }))
  return data.ok
}
export async function lRem(key, value) {
  const { data } = await client.post('/list/lrem', withConn({ key, value }))
  return data.ok
}

// set
export async function sScan(key, cursor = '0', match, count = 100) {
  const { data } = await client.post('/set/scan', withConn({ key, cursor, match, count }))
  return data
}
export async function sAdd(key, member) {
  const { data } = await client.post('/set/sadd', withConn({ key, member }))
  return data.ok
}
export async function sRem(key, member) {
  const { data } = await client.post('/set/srem', withConn({ key, member }))
  return data.ok
}

// zset
export async function zRange(key, offset = 0, limit = 100) {
  const { data } = await client.post('/zset/range', withConn({ key, offset, limit }))
  return data.items
}
export async function zAdd(key, member, score) {
  const { data } = await client.post('/zset/zadd', withConn({ key, member, score }))
  return data.ok
}
export async function zRem(key, member) {
  const { data } = await client.post('/zset/zrem', withConn({ key, member }))
  return data.ok
}
export async function zIncrBy(key, member, increment) {
  const { data } = await client.post('/zset/zincrby', withConn({ key, member, increment }))
  return data.score
}


