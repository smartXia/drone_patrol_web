const express = require('express')
const cors = require('cors')
const bodyParser = require('body-parser')
const { createClient } = require('redis')

const app = express()
app.use(cors())
app.use(bodyParser.json({ limit: '2mb' }))

function getConn(reqBody) {
  const host = reqBody.host || '127.0.0.1'
  const port = reqBody.port || 6379
  const db = reqBody.db || 0
  const password = reqBody.password || undefined
  const url = password ? `redis://:${encodeURIComponent(password)}@${host}:${port}/${db}` : `redis://${host}:${port}/${db}`
  return { url }
}

async function withClient(reqBody, handler) {
  const { url } = getConn(reqBody)
  const client = createClient({ url })
  client.on('error', () => {})
  await client.connect()
  try {
    return await handler(client)
  } finally {
    try { await client.quit() } catch { /* noop */ }
  }
}

app.post('/connect/test', async (req, res) => {
  try {
    const out = await withClient(req.body, async (c) => {
      const pong = await c.ping()
      const infoStr = await c.info()
      const versionMatch = /redis_version:(.*)/.exec(infoStr)
      const roleMatch = /role:(.*)/.exec(infoStr)
      return { ok: true, pong, info: { redis_version: versionMatch?.[1] || '', role: roleMatch?.[1] || '' } }
    })
    res.json(out)
  } catch (e) {
    res.status(400).json({ ok: false, error: String(e.message || e) })
  }
})

app.post('/scan', async (req, res) => {
  try {
    const { cursor = '0', match, count = 100, type } = req.body
    const out = await withClient(req.body, async (c) => {
      const scanOpts = { MATCH: match, COUNT: count }
      let cur = cursor
      let keys = []
      const reply = await c.scan(cur, scanOpts)
      cur = reply.cursor
      keys = reply.keys
      const result = []
      for (const k of keys) {
        if (type) {
          const t = await c.type(k)
          if (t !== type) continue
        }
        const t = await c.type(k)
        const ttl = await c.ttl(k)
        result.push({ key: k, type: t, ttl })
      }
      return { cursor: String(cur), keys: result }
    })
    res.json(out)
  } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})

app.post('/type', async (req, res) => {
  try { const data = await withClient(req.body, (c) => c.type(req.body.key)); res.json({ type: data }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})
app.post('/ttl', async (req, res) => {
  try { const data = await withClient(req.body, (c) => c.ttl(req.body.key)); res.json({ ttl: data }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})
app.post('/expire', async (req, res) => {
  try { const data = await withClient(req.body, (c) => c.expire(req.body.key, req.body.ttl)); res.json({ ok: !!data }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})
app.post('/persist', async (req, res) => {
  try { const data = await withClient(req.body, (c) => c.persist(req.body.key)); res.json({ ok: !!data }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})
app.post('/rename', async (req, res) => {
  try { await withClient(req.body, (c) => c.rename(req.body.key, req.body.newKey)); res.json({ ok: true }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})
app.post('/del', async (req, res) => {
  try { const data = await withClient(req.body, (c) => c.del(req.body.key)); res.json({ ok: data > 0 }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})

// string
app.post('/get', async (req, res) => {
  try { const v = await withClient(req.body, (c) => c.get(req.body.key)); res.json({ value: v }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})
app.post('/set', async (req, res) => {
  try { const ok = await withClient(req.body, (c) => c.set(req.body.key, req.body.value)); res.json({ ok: !!ok }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})

// hash（分页采用 HSCAN）
app.post('/hash/getall', async (req, res) => {
  try {
    const { key, limit = 100 } = req.body
    const out = await withClient(req.body, async (c) => {
      let cursor = 0
      const items = []
      do {
        const reply = await c.hScan(key, cursor, { COUNT: limit })
        cursor = reply.cursor
        const arr = reply.tuples
        for (const [f, v] of arr) {
          items.push({ field: f, value: v })
          if (items.length >= limit) return { total: null, items }
        }
      } while (cursor !== 0)
      return { total: items.length, items }
    })
    res.json(out)
  } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})
app.post('/hash/set', async (req, res) => {
  try { const n = await withClient(req.body, (c) => c.hSet(req.body.key, req.body.field, req.body.value)); res.json({ ok: n >= 0 }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})
app.post('/hash/del', async (req, res) => {
  try { const n = await withClient(req.body, (c) => c.hDel(req.body.key, req.body.field)); res.json({ ok: n > 0 }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})

// list
app.post('/list/range', async (req, res) => {
  try { const items = await withClient(req.body, (c) => c.lRange(req.body.key, req.body.start ?? 0, req.body.stop ?? 100)); res.json({ items: items.map(v => ({ value: v })) }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})
app.post('/list/lpush', async (req, res) => {
  try { await withClient(req.body, (c) => c.lPush(req.body.key, req.body.value)); res.json({ ok: true }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})
app.post('/list/rpush', async (req, res) => {
  try { await withClient(req.body, (c) => c.rPush(req.body.key, req.body.value)); res.json({ ok: true }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})
app.post('/list/lpop', async (req, res) => {
  try { const v = await withClient(req.body, (c) => c.lPop(req.body.key)); res.json({ item: v }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})
app.post('/list/rpop', async (req, res) => {
  try { const v = await withClient(req.body, (c) => c.rPop(req.body.key)); res.json({ item: v }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})
app.post('/list/set', async (req, res) => {
  try { await withClient(req.body, (c) => c.lSet(req.body.key, req.body.index, req.body.value)); res.json({ ok: true }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})
app.post('/list/lrem', async (req, res) => {
  try { const n = await withClient(req.body, (c) => c.lRem(req.body.key, 1, req.body.value)); res.json({ ok: n > 0 }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})

// set
app.post('/set/scan', async (req, res) => {
  try {
    const { key, cursor = '0', match, count = 100 } = req.body
    const out = await withClient(req.body, async (c) => {
      const reply = await c.sScan(key, cursor, { MATCH: match, COUNT: count })
      return { cursor: String(reply.cursor), items: reply.members.map(m => ({ member: m })) }
    })
    res.json(out)
  } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})
app.post('/set/sadd', async (req, res) => {
  try { const n = await withClient(req.body, (c) => c.sAdd(req.body.key, req.body.member)); res.json({ ok: n > 0 }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})
app.post('/set/srem', async (req, res) => {
  try { const n = await withClient(req.body, (c) => c.sRem(req.body.key, req.body.member)); res.json({ ok: n > 0 }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})

// zset
app.post('/zset/range', async (req, res) => {
  try {
    const { key, offset = 0, limit = 100 } = req.body
    const items = await withClient(req.body, (c) => c.zRangeWithScores(key, offset, offset + limit - 1))
    res.json({ items: items.map(({ value, score }) => ({ member: value, score })) })
  } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})
app.post('/zset/zadd', async (req, res) => {
  try { const n = await withClient(req.body, (c) => c.zAdd(req.body.key, [{ score: req.body.score, value: req.body.member }])); res.json({ ok: n >= 0 }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})
app.post('/zset/zrem', async (req, res) => {
  try { const n = await withClient(req.body, (c) => c.zRem(req.body.key, req.body.member)); res.json({ ok: n > 0 }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})
app.post('/zset/zincrby', async (req, res) => {
  try { const s = await withClient(req.body, (c) => c.zIncrBy(req.body.key, req.body.increment, req.body.member)); res.json({ score: s }) } catch (e) { res.status(400).json({ error: String(e.message || e) }) }
})

const PORT = process.env.PORT || 8001
app.listen(PORT, () => console.log(`Redis API (Node) listening on http://127.0.0.1:${PORT}`))


