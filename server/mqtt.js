// 本地 MQTT 配置管理服务（基础框架）
const path = require('path')
const fs = require('fs')
const express = require('express')
const cors = require('cors')

// lowdb v6 (Esm to CJS via dynamic import)
let Low
let JSONFile

async function getLowdb() {
  if (!Low || !JSONFile) {
    const lowdb = await import('lowdb')
    const lowdbNode = await import('lowdb/node')
    Low = lowdb.Low
    JSONFile = lowdbNode.JSONFile
  }
  return { Low, JSONFile }
}

async function createDb(filePath) {
  const { Low, JSONFile } = await getLowdb()
  const adapter = new JSONFile(filePath)
  const db = new Low(adapter, { profiles: [], currentProfileId: null })
  await db.read()
  // 初始化默认结构
  db.data = db.data || { profiles: [], currentProfileId: null }
  await db.write()
  return db
}

async function bootstrap() {
  const app = express()
  app.use(cors())
  app.use(express.json())

  const dataDir = path.join(__dirname, 'data')
  if (!fs.existsSync(dataDir)) fs.mkdirSync(dataDir, { recursive: true })
  const dbFile = path.join(dataDir, 'mqtt-profiles.json')
  const db = await createDb(dbFile)

  // 健康检查与版本信息
  app.get('/api/health', (req, res) => {
    res.json({ code: 0, message: 'ok', data: { service: 'mqtt-local', version: '1.0.0' } })
  })

  // ===== 配置 CRUD =====
  app.get('/api/mqtt/profiles', async (req, res) => {
    await db.read()
    res.json({ code: 0, message: 'ok', data: db.data.profiles })
  })

  app.get('/api/mqtt/profiles/:id', async (req, res) => {
    const { id } = req.params
    await db.read()
    const profile = db.data.profiles.find(p => p.id === id)
    if (!profile) return res.status(404).json({ code: 404, message: 'not found' })
    res.json({ code: 0, message: 'ok', data: profile })
  })

  app.post('/api/mqtt/profiles', async (req, res) => {
    const { name, config, isDefault } = req.body || {}
    if (!name || !config || !config.host || !config.port || !config.protocol) {
      return res.status(400).json({ code: 400, message: 'invalid payload' })
    }
    const id = (await import('nanoid')).nanoid()
    const now = Date.now()
    const profile = { id, name, config, isDefault: !!isDefault, updatedAt: now }
    await db.read()
    db.data.profiles.push(profile)
    if (profile.isDefault) {
      db.data.profiles.forEach(p => { if (p.id !== id) p.isDefault = false })
      db.data.currentProfileId = id
    }
    await db.write()
    res.json({ code: 0, message: 'ok', data: { id } })
  })

  app.put('/api/mqtt/profiles/:id', async (req, res) => {
    const { id } = req.params
    const { name, config, isDefault } = req.body || {}
    await db.read()
    const idx = db.data.profiles.findIndex(p => p.id === id)
    if (idx === -1) return res.status(404).json({ code: 404, message: 'not found' })
    const now = Date.now()
    const merged = {
      ...db.data.profiles[idx],
      ...(name !== undefined ? { name } : {}),
      ...(config !== undefined ? { config } : {}),
      ...(isDefault !== undefined ? { isDefault: !!isDefault } : {}),
      updatedAt: now
    }
    db.data.profiles[idx] = merged
    if (merged.isDefault) {
      db.data.profiles.forEach(p => { if (p.id !== id) p.isDefault = false })
      db.data.currentProfileId = id
    }
    await db.write()
    res.json({ code: 0, message: 'ok', data: true })
  })

  app.delete('/api/mqtt/profiles/:id', async (req, res) => {
    const { id } = req.params
    await db.read()
    const before = db.data.profiles.length
    db.data.profiles = db.data.profiles.filter(p => p.id !== id)
    if (db.data.currentProfileId === id) db.data.currentProfileId = null
    const removed = before !== db.data.profiles.length
    await db.write()
    res.json({ code: 0, message: 'ok', data: removed })
  })

  app.post('/api/mqtt/profiles/:id/default', async (req, res) => {
    const { id } = req.params
    await db.read()
    const exists = db.data.profiles.some(p => p.id === id)
    if (!exists) return res.status(404).json({ code: 404, message: 'not found' })
    db.data.profiles.forEach(p => { p.isDefault = p.id === id })
    db.data.currentProfileId = id
    await db.write()
    res.json({ code: 0, message: 'ok', data: true })
  })

  app.get('/api/mqtt/profiles-default', async (req, res) => {
    await db.read()
    const d = db.data.profiles.find(p => p.isDefault) || null
    res.json({ code: 0, message: 'ok', data: d })
  })

  // ===== 测试连接 =====
  app.post('/api/mqtt/test', async (req, res) => {
    const payload = req.body || {}
    const { profileId, config } = payload
    let cfg = config
    await db.read()
    if (!cfg && profileId) {
      const p = db.data.profiles.find(x => x.id === profileId)
      if (!p) return res.status(404).json({ code: 404, message: 'profile not found' })
      cfg = p.config
    }
    if (!cfg || !cfg.host || !cfg.port || !cfg.protocol) {
      return res.status(400).json({ code: 400, message: 'invalid config' })
    }

    try {
      const mqtt = require('mqtt')
      const protocol = cfg.protocol
      const pathSuffix = cfg.path || '/mqtt'
      const url = `${protocol}://${cfg.host}:${cfg.port}${pathSuffix}`
      const client = mqtt.connect(url, {
        clientId: cfg.clientId || `test_${Date.now()}`,
        username: cfg.username,
        password: cfg.password,
        clean: cfg.clean !== false,
        connectTimeout: Math.min(cfg.connectTimeout || 10000, 20000),
        keepalive: cfg.keepalive || 60,
        protocolVersion: 4,
        rejectUnauthorized: !!cfg.rejectUnauthorized,
        reconnectPeriod: 0
      })

      let settled = false
      const timer = setTimeout(() => {
        if (settled) return
        settled = true
        try { client.end(true) } catch (e) {}
        res.status(504).json({ code: 504, message: 'connect timeout' })
      }, Math.min(cfg.connectTimeout || 10000, 20000))

      client.on('connect', () => {
        if (settled) return
        settled = true
        clearTimeout(timer)
        try { client.end(true) } catch (e) {}
        res.json({ code: 0, message: 'ok', data: { url } })
      })
      client.on('error', (err) => {
        if (settled) return
        settled = true
        clearTimeout(timer)
        try { client.end(true) } catch (e) {}
        res.status(502).json({ code: 502, message: err.message || 'connect error' })
      })
      client.on('close', () => {
        // 如果先触发 close 且未 settle，当作失败
        if (settled) return
        settled = true
        clearTimeout(timer)
        res.status(502).json({ code: 502, message: 'connection closed' })
      })
    } catch (e) {
      res.status(500).json({ code: 500, message: e.message || 'internal error' })
    }
  })

  const PORT = process.env.MQTT_PROFILE_PORT || 7077
  app.listen(PORT, () => {
    console.log(`[mqtt-local] running at http://127.0.0.1:${PORT}`)
  })
}

bootstrap().catch((err) => {
  console.error('[mqtt-local] bootstrap failed:', err)
  process.exit(1)
})


