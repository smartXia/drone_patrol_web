from fastapi import FastAPI, HTTPException, Body
from fastapi.responses import JSONResponse
from pydantic import BaseModel
import socket
import time
import json
import os
import sqlite3
import uuid
import redis
import paho.mqtt.client as mqtt
from typing import Any, Dict, List, Optional, Tuple

from fastapi.middleware.cors import CORSMiddleware

# 数据库路径
DB_PATH = os.path.join(os.path.dirname(__file__), '..', 'data', 'backend.db')

def get_db():
    """获取SQLite数据库连接"""
    conn = sqlite3.connect(DB_PATH)
    conn.row_factory = sqlite3.Row
    # 创建表
    conn.execute("""
        CREATE TABLE IF NOT EXISTS mqtt_profiles (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            config TEXT NOT NULL,
            is_default INTEGER DEFAULT 0,
            updated_at INTEGER
        )
    """)
    conn.execute("CREATE INDEX IF NOT EXISTS idx_mqtt_profiles_default ON mqtt_profiles(is_default)")
    conn.commit()
    return conn

class MqttProfilePayload(BaseModel):
    name: str
    config: dict
    isDefault: bool = False

class RedisConnectPayload(BaseModel):
    host: str
    port: int = 6379
    password: str = None
    db: int = 0

app = FastAPI(title="Unified Backend (Python)")

# CORS（允许前端直接请求本服务）
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


def _decode(value: Any) -> Any:
    if isinstance(value, bytes):
        try:
            return value.decode("utf-8", errors="ignore")
        except Exception:
            return value.hex()
    if isinstance(value, list):
        return [_decode(v) for v in value]
    if isinstance(value, tuple):
        return tuple(_decode(v) for v in value)
    if isinstance(value, dict):
        return { _decode(k): _decode(v) for k, v in value.items() }
    return value


def _get_redis(payload: Dict[str, Any]) -> redis.Redis:
    host = payload.get("host", "127.0.0.1")
    port = int(payload.get("port", 6379))
    db = int(payload.get("db", 0))
    password = payload.get("password") or None
    return redis.Redis(host=host, port=port, db=db, password=password, socket_timeout=5)

@app.get('/api/health')
def health():
    """健康检查"""
    return {"code": 0, "message": "ok", "data": {"service": "unified-python", "version": "1.0.0"}}

@app.get('/api/network/ping')
def ping(host: str, port: int):
    """网络连通性检查"""
    start = time.time()
    s = socket.socket()
    s.settimeout(5)
    try:
        s.connect((host, port))
        s.close()
        return {
            "code": 0, 
            "message": "ok", 
            "data": {
                "host": host, 
                "port": str(port), 
                "success": True, 
                "duration": int((time.time() - start) * 1000)
            }
        }
    except Exception as e:
        return JSONResponse(
            status_code=502, 
            content={
                "code": 502, 
                "message": str(e), 
                "data": {
                    "host": host, 
                    "port": str(port), 
                    "success": False, 
                    "duration": int((time.time() - start) * 1000)
                }
            }
        )

# MQTT 配置管理
@app.get('/api/mqtt/profiles')
def list_profiles():
    """获取MQTT配置列表"""
    conn = get_db()
    try:
        rows = conn.execute('SELECT * FROM mqtt_profiles ORDER BY is_default DESC, updated_at DESC').fetchall()
        data = []
        for r in rows:
            data.append({
                "id": r['id'],
                "name": r['name'],
                "config": json.loads(r['config']),
                "isDefault": bool(r['is_default']),
                "updatedAt": r['updated_at']
            })
        return {"code": 0, "message": "ok", "data": data}
    finally:
        conn.close()

@app.post('/api/mqtt/profiles')
def create_profile(payload: MqttProfilePayload):
    """创建MQTT配置"""
    pid = uuid.uuid4().hex
    now = int(time.time() * 1000)
    conn = get_db()
    try:
        if payload.isDefault:
            conn.execute('UPDATE mqtt_profiles SET is_default = 0')
        conn.execute(
            'INSERT INTO mqtt_profiles (id, name, config, is_default, updated_at) VALUES (?, ?, ?, ?, ?)',
            (pid, payload.name, json.dumps(payload.config, ensure_ascii=False), 1 if payload.isDefault else 0, now)
        )
        conn.commit()
        return {"code": 0, "message": "ok", "data": {"id": pid}}
    finally:
        conn.close()

@app.get('/api/mqtt/profiles/{pid}')
def get_profile(pid: str):
    """获取单个MQTT配置"""
    conn = get_db()
    try:
        r = conn.execute('SELECT * FROM mqtt_profiles WHERE id = ?', (pid,)).fetchone()
        if not r:
            raise HTTPException(status_code=404, detail="Profile not found")
        return {
            "code": 0, 
            "message": "ok", 
            "data": {
                "id": r['id'],
                "name": r['name'],
                "config": json.loads(r['config']),
                "isDefault": bool(r['is_default']),
                "updatedAt": r['updated_at']
            }
        }
    finally:
        conn.close()

@app.put('/api/mqtt/profiles/{pid}')
def update_profile(pid: str, payload: dict):
    """更新MQTT配置"""
    conn = get_db()
    try:
        r = conn.execute('SELECT * FROM mqtt_profiles WHERE id = ?', (pid,)).fetchone()
        if not r:
            raise HTTPException(status_code=404, detail="Profile not found")
        
        name = payload.get('name', r['name'])
        config = json.dumps(payload.get('config', json.loads(r['config'])), ensure_ascii=False)
        is_default = 1 if payload.get('isDefault', bool(r['is_default'])) else 0
        now = int(time.time() * 1000)
        
        if is_default:
            conn.execute('UPDATE mqtt_profiles SET is_default = 0')
        
        conn.execute(
            'UPDATE mqtt_profiles SET name = ?, config = ?, is_default = ?, updated_at = ? WHERE id = ?',
            (name, config, is_default, now, pid)
        )
        conn.commit()
        return {"code": 0, "message": "ok", "data": True}
    finally:
        conn.close()

@app.delete('/api/mqtt/profiles/{pid}')
def delete_profile(pid: str):
    """删除MQTT配置"""
    conn = get_db()
    try:
        conn.execute('DELETE FROM mqtt_profiles WHERE id = ?', (pid,))
        conn.commit()
        return {"code": 0, "message": "ok", "data": True}
    finally:
        conn.close()

@app.post('/api/mqtt/profiles/{pid}/default')
def set_default(pid: str):
    """设置默认MQTT配置"""
    conn = get_db()
    try:
        exists = conn.execute('SELECT 1 FROM mqtt_profiles WHERE id = ?', (pid,)).fetchone()
        if not exists:
            raise HTTPException(status_code=404, detail="Profile not found")
        
        conn.execute('UPDATE mqtt_profiles SET is_default = (id = ?)', (pid,))
        conn.commit()
        return {"code": 0, "message": "ok", "data": True}
    finally:
        conn.close()

@app.post('/api/mqtt/test')
def test_mqtt_connection(payload: dict):
    """测试MQTT连接"""
    config = payload.get('config', {})
    host = config.get('host', 'localhost')
    port = config.get('port', 1883)
    protocol = config.get('protocol', 'tcp')
    username = config.get('username')
    password = config.get('password')
    
    try:
        # 创建MQTT客户端
        client = mqtt.Client()
        if username:
            client.username_pw_set(username, password)
        
        # 设置连接超时
        client.connect_async(host, port, 60)
        client.loop_start()
        
        # 等待连接结果
        time.sleep(2)
        
        if client.is_connected():
            client.disconnect()
            client.loop_stop()
            return {"code": 0, "message": "连接成功", "data": {"connected": True}}
        else:
            client.loop_stop()
            return {"code": 1, "message": "连接失败", "data": {"connected": False}}
            
    except Exception as e:
        return {"code": 1, "message": f"连接错误: {str(e)}", "data": {"connected": False}}

# Redis 代理接口
@app.post('/api/redis/connect/test')
@app.post('/connect/test')
def test_redis_connection(payload: RedisConnectPayload):
    """测试Redis连接（兼容前端：返回 { ok: boolean, message? }）"""
    try:
        r = redis.Redis(host=payload.host, port=payload.port, password=payload.password, db=payload.db, socket_timeout=5)
        r.ping()
        return {"ok": True}
    except Exception as e:
        return {"ok": False, "message": str(e)}

@app.post('/scan')
def redis_scan(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    cursor = payload.get("cursor", "0")
    match = payload.get("match")
    count = int(payload.get("count", 100))
    type_filter = payload.get("type")
    try:
        cursor_int = int(cursor)
    except Exception:
        cursor_int = 0
    try:
        if match:
            cursor_int, keys = r.scan(cursor=cursor_int, match=match, count=count)
        else:
            cursor_int, keys = r.scan(cursor=cursor_int, count=count)
        keys = [_decode(k) for k in keys]
        items = []
        for k in keys:
            try:
                kt = _decode(r.type(k))
                if isinstance(kt, bytes):
                    kt = kt.decode()
                if type_filter and kt != type_filter:
                    continue
                ttl = r.ttl(k)
            except Exception:
                kt = "string"
                ttl = -1
            items.append({"key": k, "type": kt, "ttl": ttl})
        next_cursor = str(cursor_int)
        return {"keys": items, "cursor": next_cursor, "hasMore": next_cursor != "0"}
    except Exception as e:
        return JSONResponse(status_code=500, content={"message": str(e)})


@app.post('/type')
def redis_type(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    t = _decode(r.type(key))
    if isinstance(t, bytes):
        t = t.decode()
    return {"type": t}


@app.post('/ttl')
def redis_ttl(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    return {"ttl": r.ttl(key)}


@app.post('/expire')
def redis_expire(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    ttl = int(payload.get("ttl", -1))
    return {"ok": bool(r.expire(key, ttl))}


@app.post('/persist')
def redis_persist(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    return {"ok": bool(r.persist(key))}


@app.post('/rename')
def redis_rename(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    new_key = payload.get("newKey")
    return {"ok": bool(r.rename(key, new_key))}


@app.post('/del')
def redis_del(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    return {"ok": bool(r.delete(key))}


# string
@app.post('/get')
def str_get(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    v = r.get(key)
    return {"value": _decode(v) if v is not None else None}


@app.post('/set')
def str_set(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    value = payload.get("value")
    return {"ok": bool(r.set(key, value))}


# hash
@app.post('/hash/getall')
def h_getall(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    offset = int(payload.get("offset", 0))
    limit = int(payload.get("limit", 100))
    data = r.hgetall(key)
    # data is dict field->value
    items_all = [{"field": _decode(f), "value": _decode(v)} for f, v in data.items()]
    total = len(items_all)
    items = items_all[offset: offset + limit]
    return {"items": items, "total": total, "offset": offset, "limit": limit}


@app.post('/hash/set')
def h_set(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    field = payload.get("field")
    value = payload.get("value")
    return {"ok": bool(r.hset(key, field, value))}


@app.post('/hash/del')
def h_del(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    field = payload.get("field")
    return {"ok": bool(r.hdel(key, field))}


# list
@app.post('/list/range')
def l_range(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    start = int(payload.get("start", 0))
    stop = int(payload.get("stop", 100))
    vals = r.lrange(key, start, stop)
    return {"items": [_decode(v) for v in vals]}


@app.post('/list/lpush')
def l_lpush(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    value = payload.get("value")
    return {"ok": bool(r.lpush(key, value))}


@app.post('/list/rpush')
def l_rpush(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    value = payload.get("value")
    return {"ok": bool(r.rpush(key, value))}


@app.post('/list/lpop')
def l_lpop(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    v = r.lpop(key)
    return {"item": _decode(v) if v is not None else None}


@app.post('/list/rpop')
def l_rpop(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    v = r.rpop(key)
    return {"item": _decode(v) if v is not None else None}


@app.post('/list/set')
def l_set(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    index = int(payload.get("index", 0))
    value = payload.get("value")
    return {"ok": r.lset(key, index, value) is True or True}


@app.post('/list/lrem')
def l_lrem(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    value = payload.get("value")
    # count=0 -> remove all occurrences
    removed = r.lrem(key, 0, value)
    return {"ok": removed >= 0}


# set
@app.post('/set/scan')
def s_scan(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    cursor = payload.get("cursor", "0")
    match = payload.get("match")
    count = int(payload.get("count", 100))
    try:
        cursor_int = int(cursor)
    except Exception:
        cursor_int = 0
    if match:
        cursor_int, members = r.sscan(name=key, cursor=cursor_int, match=match, count=count)
    else:
        cursor_int, members = r.sscan(name=key, cursor=cursor_int, count=count)
    items = [_decode(m) for m in members]
    next_cursor = str(cursor_int)
    return {"items": items, "cursor": next_cursor, "hasMore": next_cursor != "0"}


@app.post('/set/sadd')
def s_sadd(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    member = payload.get("member")
    return {"ok": bool(r.sadd(key, member))}


@app.post('/set/srem')
def s_srem(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    member = payload.get("member")
    return {"ok": bool(r.srem(key, member))}


# zset
@app.post('/zset/range')
def z_range(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    offset = int(payload.get("offset", 0))
    limit = int(payload.get("limit", 100))
    stop = offset + limit - 1
    pairs = r.zrange(key, offset, stop, withscores=True)
    items = [{"member": _decode(m), "score": float(s)} for m, s in pairs]
    return {"items": items, "offset": offset, "limit": limit}


@app.post('/zset/zadd')
def z_add(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    member = payload.get("member")
    score = float(payload.get("score", 0))
    return {"ok": bool(r.zadd(key, {member: score}))}


@app.post('/zset/zrem')
def z_rem(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    member = payload.get("member")
    return {"ok": bool(r.zrem(key, member))}


@app.post('/zset/zincrby')
def z_incrby(payload: Dict[str, Any] = Body(...)):
    r = _get_redis(payload)
    key = payload.get("key")
    member = payload.get("member")
    inc = float(payload.get("increment", 1))
    new_score = r.zincrby(key, inc, member)
    return {"score": float(new_score)}

if __name__ == '__main__':
    import uvicorn
    uvicorn.run(app, host='127.0.0.1', port=8080)
