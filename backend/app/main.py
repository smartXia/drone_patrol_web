# -*- coding: utf-8 -*-
import sys
import io

# 设置标准输出编码为UTF-8
sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8')

from fastapi import FastAPI, HTTPException, Body, WebSocket, WebSocketDisconnect
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
import asyncio
import threading
from typing import Any, Dict, List, Optional, Tuple

from fastapi.middleware.cors import CORSMiddleware

# 数据库路径
DB_PATH = os.path.join(os.path.dirname(__file__), '..', 'data', 'backend.db')

# MQTT 代理管理器
class MQTTProxy:
    def __init__(self):
        self.mqtt_client = None
        self.websocket_clients = set()
        self.is_connected = False
        self._loop = None
        
    def set_event_loop(self, loop):
        """设置事件循环引用"""
        self._loop = loop
        
    def _schedule_broadcast(self, message):
        """线程安全地调度广播消息"""
        if self._loop and not self._loop.is_closed():
            asyncio.run_coroutine_threadsafe(
                self.broadcast_to_websockets(message), 
                self._loop
            )
        
    def connect_mqtt(self, host: str, port: int, username: str = None, password: str = None, client_id: str = None):
        """连接到 MQTT 服务器"""
        try:
            # 如果已经连接，先断开
            if self.mqtt_client and self.is_connected:
                print("MQTT 已连接，先断开旧连接")
                self.mqtt_client.disconnect()
                self.mqtt_client.loop_stop()
                self.mqtt_client = None
                self.is_connected = False
            
            # 创建新的 MQTT 客户端
            self.mqtt_client = mqtt.Client(client_id or f"proxy_{uuid.uuid4().hex[:8]}")
            
            if username:
                self.mqtt_client.username_pw_set(username, password)
            
            # 设置回调函数
            self.mqtt_client.on_connect = self.on_mqtt_connect
            self.mqtt_client.on_disconnect = self.on_mqtt_disconnect
            self.mqtt_client.on_message = self.on_mqtt_message
            
            # 连接到 MQTT 服务器
            self.mqtt_client.connect(host, port, 60)
            self.mqtt_client.loop_start()
            
            return True
        except Exception as e:
            print(f"MQTT 连接失败: {e}")
            # 发送连接失败消息
            self._schedule_broadcast({
                "type": "connect_result",
                "success": False,
                "message": f"MQTT 连接失败: {e}"
            })
            return False
    
    def on_mqtt_connect(self, client, userdata, flags, rc):
        """MQTT 连接成功回调"""
        if rc == 0:
            self.is_connected = True
            print("MQTT 连接成功")
            # 通知所有 WebSocket 客户端连接成功
            self._schedule_broadcast({
                "type": "connect_result",
                "success": True,
                "message": "MQTT 连接成功"
            })
        else:
            print(f"MQTT 连接失败，错误码: {rc}")
            # 通知所有 WebSocket 客户端连接失败
            self._schedule_broadcast({
                "type": "connect_result",
                "success": False,
                "message": f"MQTT 连接失败，错误码: {rc}"
            })
    
    def on_mqtt_disconnect(self, client, userdata, rc):
        """MQTT 断开连接回调"""
        self.is_connected = False
        print("MQTT 连接断开")
        self._schedule_broadcast({
            "type": "mqtt_disconnected",
            "message": "MQTT 连接断开"
        })
    
    def on_mqtt_message(self, client, userdata, msg):
        """接收到 MQTT 消息回调"""
        print(f"收到MQTT消息: {msg.topic}")
        message = {
            "type": "mqtt_message",
            "topic": msg.topic,
            "payload": msg.payload.decode('utf-8'),
            "qos": msg.qos,
            "retain": msg.retain
        }
        print(f"广播消息到 {len(self.websocket_clients)} 个WebSocket客户端")
        self._schedule_broadcast(message)
    
    def subscribe(self, topic: str, qos: int = 0):
        """订阅 MQTT 主题"""
        print(f"尝试订阅主题: {topic}, QoS: {qos}")
        print(f"MQTT客户端状态: {self.mqtt_client is not None}, 连接状态: {self.is_connected}")
        
        if self.mqtt_client and self.is_connected:
            result = self.mqtt_client.subscribe(topic, qos)
            print(f"订阅结果: {result}")
            if result[0] == mqtt.MQTT_ERR_SUCCESS:
                print(f"订阅成功: {topic}")
                success_message = {
                    "type": "subscription_success",
                    "topic": topic,
                    "qos": qos
                }
                print(f"发送订阅成功消息: {success_message}")
                self._schedule_broadcast(success_message)
                return True
            else:
                print(f"订阅失败: {topic}, 错误码: {result[0]}")
        else:
            print(f"无法订阅: MQTT客户端={self.mqtt_client is not None}, 连接状态={self.is_connected}")
        return False
    
    def publish(self, topic: str, payload: str, qos: int = 0, retain: bool = False):
        """发布 MQTT 消息"""
        if self.mqtt_client and self.is_connected:
            result = self.mqtt_client.publish(topic, payload, qos, retain)
            if result.rc == mqtt.MQTT_ERR_SUCCESS:
                return True
        return False
    
    def add_websocket_client(self, websocket: WebSocket):
        """添加 WebSocket 客户端"""
        self.websocket_clients.add(websocket)
    
    def remove_websocket_client(self, websocket: WebSocket):
        """移除 WebSocket 客户端"""
        self.websocket_clients.discard(websocket)
    
    async def broadcast_to_websockets(self, message: dict):
        """向所有 WebSocket 客户端广播消息"""
        if self.websocket_clients:
            disconnected = set()
            for websocket in self.websocket_clients:
                try:
                    await websocket.send_text(json.dumps(message))
                except:
                    disconnected.add(websocket)
            
            # 移除断开的连接
            for websocket in disconnected:
                self.websocket_clients.discard(websocket)
    
    def disconnect(self):
        """断开 MQTT 连接"""
        if self.mqtt_client:
            self.mqtt_client.loop_stop()
            self.mqtt_client.disconnect()
            self.is_connected = False

# 全局 MQTT 代理实例
mqtt_proxy = MQTTProxy()

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

@app.on_event("startup")
async def startup_event():
    """应用启动时设置事件循环引用"""
    mqtt_proxy.set_event_loop(asyncio.get_running_loop())

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

@app.get("/api/error-codes")
async def get_error_codes():
    """获取错误码列表"""
    try:
        # 读取错误码文档
        error_codes_file = os.path.join(os.path.dirname(__file__), '..', '..', 'public', 'docs', 'dji-error-codes.md')
        with open(error_codes_file, 'r', encoding='utf-8') as f:
            content = f.read()
        
        lines = content.split('\n')
        error_codes = []
        
        for i, line in enumerate(lines[1:], 1):  # 跳过标题行
            if line.strip():
                parts = line.split('\t')
                if len(parts) >= 2:
                    code = parts[0].strip()
                    description = parts[1].strip()
                    
                    # 根据错误码范围确定分类
                    code_num = int(code) if code.isdigit() else 0
                    category = get_error_category(code_num)
                    severity = get_error_severity(description)
                    
                    error_codes.append({
                        "code": code,
                        "description": description,
                        "category": category,
                        "severity": severity
                    })
        
        return {
            "code": 0,
            "message": "ok",
            "data": error_codes
        }
    except Exception as e:
        return {
            "code": 1,
            "message": f"获取错误码失败: {str(e)}",
            "data": []
        }

def get_error_category(code_num):
    """根据错误码确定分类"""
    if 312000 <= code_num < 315000:
        return "upgrade"
    elif 314000 <= code_num < 317000:
        return "flight"
    elif 315000 <= code_num < 316000:
        return "communication"
    elif 316000 <= code_num < 317000:
        return "battery"
    elif 317000 <= code_num < 319000:
        return "media"
    elif 319000 <= code_num < 322000:
        return "system"
    elif 321000 <= code_num < 325000:
        return "flight"
    elif 325000 <= code_num < 327000:
        return "network"
    elif 327000 <= code_num < 329000:
        return "camera"
    elif 336000 <= code_num < 339000:
        return "flight"
    elif 338000 <= code_num < 339000:
        return "flight"
    elif 514000 <= code_num < 515000:
        return "weather"
    else:
        return "other"

def get_error_severity(description):
    """根据描述确定严重程度"""
    if any(keyword in description for keyword in ["失败", "异常", "错误", "无法"]):
        return "error"
    elif any(keyword in description for keyword in ["超时", "请重试", "稍后"]):
        return "warning"
    elif any(keyword in description for keyword in ["成功", "完成"]):
        return "success"
    else:
        return "info"

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

# WebSocket 端点
@app.websocket("/ws/mqtt")
async def websocket_mqtt_proxy(websocket: WebSocket):
    """MQTT WebSocket 代理端点"""
    await websocket.accept()
    mqtt_proxy.add_websocket_client(websocket)
    
    try:
        while True:
            # 接收来自前端的消息
            data = await websocket.receive_text()
            message = json.loads(data)
            
            if message["type"] == "connect":
                # 连接 MQTT 服务器
                config = message["config"]
                
                # 检查是否已经连接到相同的服务器
                if (mqtt_proxy.is_connected and 
                    hasattr(mqtt_proxy, '_last_config') and
                    mqtt_proxy._last_config == config):
                    print("MQTT 已连接到相同配置，跳过重复连接")
                    await websocket.send_text(json.dumps({
                        "type": "connect_result",
                        "success": True
                    }))
                else:
                    # 先发送连接请求，不等待结果
                    mqtt_proxy.connect_mqtt(
                        host=config["host"],
                        port=config["port"],
                        username=config.get("username"),
                        password=config.get("password"),
                        client_id=config.get("clientId")
                    )
                    
                    # 保存配置
                    mqtt_proxy._last_config = config
                    
                    # 不立即发送结果，等待 on_mqtt_connect 回调
                    print(f"MQTT 连接请求已发送到 {config['host']}:{config['port']}")
                
            elif message["type"] == "subscribe":
                # 订阅主题
                topic = message["topic"]
                qos = message.get("qos", 0)
                success = mqtt_proxy.subscribe(topic, qos)
                
                await websocket.send_text(json.dumps({
                    "type": "subscribe_result",
                    "topic": topic,
                    "success": success
                }))
                
                if not success:
                    print(f"订阅失败: {topic}")
                
            elif message["type"] == "publish":
                # 发布消息
                topic = message["topic"]
                payload = message["payload"]
                qos = message.get("qos", 0)
                retain = message.get("retain", False)
                success = mqtt_proxy.publish(topic, payload, qos, retain)
                
                await websocket.send_text(json.dumps({
                    "type": "publish_result",
                    "topic": topic,
                    "success": success
                }))
                
            elif message["type"] == "disconnect":
                # 断开 MQTT 连接
                mqtt_proxy.disconnect()
                await websocket.send_text(json.dumps({
                    "type": "disconnect_result",
                    "success": True
                }))
                
    except WebSocketDisconnect:
        mqtt_proxy.remove_websocket_client(websocket)
    except Exception as e:
        print(f"WebSocket 错误: {e}")
        mqtt_proxy.remove_websocket_client(websocket)

if __name__ == '__main__':
    import uvicorn
    uvicorn.run(app, host='127.0.0.1', port=8080)
