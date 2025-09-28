# API参考文档

## 概述

本文档描述了无人机监控系统Go后端的所有API端点。

## 基础信息

- **基础URL**: `http://localhost:18080`
- **内容类型**: `application/json`
- **响应格式**: 统一JSON格式

### 统一响应格式

```json
{
  "code": 0,           // 状态码: 0=成功, 1=失败
  "message": "ok",     // 消息描述
  "data": {}           // 响应数据(可选)
}
```

## 健康检查和工具API

### 健康检查
```http
GET /api/health
```

**响应示例:**
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "service": "unified-go",
    "version": "1.0.0"
  }
}
```

### 获取错误码列表
```http
GET /api/error-codes
```

**响应示例:**
```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "code": "1001",
      "description": "系统初始化失败",
      "category": "系统错误",
      "severity": "严重"
    }
  ]
}
```

### 网络Ping测试
```http
GET /api/network/ping?host=example.com&port=80
```

**参数:**
- `host` (string, required): 目标主机
- `port` (int, required): 目标端口

**响应示例:**
```json
{
  "code": 0,
  "message": "ping成功",
  "data": {
    "host": "example.com",
    "port": 80,
    "ping": true
  }
}
```

## 设备管理API

### 获取设备列表
```http
GET /api/devices
```

**响应示例:**
```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": "device-uuid",
      "name": "无人机001",
      "sn": "8UUXN2B00A00ST",
      "type": "drone",
      "status": "online",
      "lastSeen": "2024-01-01T12:00:00Z",
      "createdAt": 1704067200000,
      "updatedAt": 1704067200000,
      "isCurrent": true,
      "isGateway": false
    }
  ]
}
```

### 获取当前设备信息
```http
GET /api/devices/current
```

**响应示例:**
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "device": {
      "id": "device-uuid",
      "name": "无人机001",
      "sn": "8UUXN2B00A00ST",
      "type": "drone",
      "status": "online",
      "lastSeen": "2024-01-01T12:00:00Z"
    },
    "gateway": {
      "id": "gateway-uuid",
      "name": "网关001",
      "sn": "8UUXN2B00A00GW",
      "type": "gateway",
      "status": "online",
      "lastSeen": "2024-01-01T12:00:00Z"
    }
  }
}
```

### 创建设备
```http
POST /api/devices
Content-Type: application/json

{
  "name": "无人机001",
  "sn": "8UUXN2B00A00ST",
  "type": "drone",
  "status": "offline"
}
```

**响应示例:**
```json
{
  "code": 0,
  "message": "设备创建成功",
  "data": {
    "id": "device-uuid"
  }
}
```

### 更新设备
```http
PUT /api/devices/{device_id}
Content-Type: application/json

{
  "name": "无人机001-更新",
  "status": "online"
}
```

### 删除设备
```http
DELETE /api/devices/{device_id}
```

### 设置当前设备
```http
POST /api/devices/{device_id}/set-current
```

### 设置网关设备
```http
POST /api/devices/{device_id}/set-gateway
```

### 清空所有设备
```http
DELETE /api/devices/clear
```

### 删除默认设备
```http
DELETE /api/devices/remove-defaults
```

## MQTT配置管理API

### 获取MQTT配置列表
```http
GET /api/mqtt/profiles
```

**响应示例:**
```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": "profile-uuid",
      "name": "默认配置",
      "config": {
        "broker": "localhost",
        "port": 1883,
        "username": "user",
        "password": "pass",
        "clientId": "client001"
      },
      "isDefault": true,
      "updatedAt": 1704067200000
    }
  ]
}
```

### 创建MQTT配置
```http
POST /api/mqtt/profiles
Content-Type: application/json

{
  "name": "新配置",
  "config": {
    "broker": "localhost",
    "port": 1883,
    "username": "user",
    "password": "pass",
    "clientId": "client001"
  },
  "isDefault": false
}
```

### 获取单个MQTT配置
```http
GET /api/mqtt/profiles/{pid}
```

### 更新MQTT配置
```http
PUT /api/mqtt/profiles/{pid}
Content-Type: application/json

{
  "name": "更新配置",
  "config": {
    "broker": "newhost",
    "port": 1883
  },
  "isDefault": true
}
```

### 删除MQTT配置
```http
DELETE /api/mqtt/profiles/{pid}
```

### 设置默认MQTT配置
```http
POST /api/mqtt/profiles/{pid}/default
```

### 测试MQTT连接
```http
POST /api/mqtt/test
Content-Type: application/json

{
  "broker": "localhost",
  "port": 1883,
  "username": "user",
  "password": "pass",
  "clientId": "test-client"
}
```

**响应示例:**
```json
{
  "code": 0,
  "message": "连接成功",
  "data": {
    "connected": true
  }
}
```

## Redis代理API

### 测试Redis连接
```http
POST /api/redis/connect/test
Content-Type: application/json

{
  "host": "localhost",
  "port": 6379,
  "password": "",
  "db": 0
}
```

### 扫描Redis键
```http
POST /scan
Content-Type: application/json

{
  "key": "prefix*"
}
```

**响应示例:**
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "keys": ["key1", "key2", "key3"],
    "cursor": 0
  }
}
```

### 获取键类型
```http
POST /type
Content-Type: application/json

{
  "key": "mykey"
}
```

### 获取键TTL
```http
POST /ttl
Content-Type: application/json

{
  "key": "mykey"
}
```

### 设置键过期
```http
POST /expire?seconds=3600
Content-Type: application/json

{
  "key": "mykey"
}
```

### 持久化键
```http
POST /persist
Content-Type: application/json

{
  "key": "mykey"
}
```

### 重命名键
```http
POST /rename?newKey=newkey
Content-Type: application/json

{
  "key": "mykey"
}
```

### 删除键
```http
POST /del
Content-Type: application/json

{
  "key": "mykey"
}
```

### 字符串操作

#### 获取字符串值
```http
POST /get
Content-Type: application/json

{
  "key": "mykey"
}
```

#### 设置字符串值
```http
POST /set
Content-Type: application/json

{
  "key": "mykey",
  "value": "myvalue"
}
```

### 哈希操作

#### 获取哈希所有字段
```http
POST /hash/getall
Content-Type: application/json

{
  "key": "myhash"
}
```

#### 设置哈希字段
```http
POST /hash/set
Content-Type: application/json

{
  "key": "myhash",
  "value": {
    "field1": "value1",
    "field2": "value2"
  }
}
```

#### 删除哈希字段
```http
POST /hash/del
Content-Type: application/json

{
  "key": "myhash",
  "field": "field1"
}
```

### 列表操作

#### 获取列表范围
```http
POST /list/range
Content-Type: application/json

{
  "key": "mylist",
  "start": 0,
  "stop": -1
}
```

#### 左推入列表
```http
POST /list/lpush
Content-Type: application/json

{
  "key": "mylist",
  "value": "item1"
}
```

#### 右推入列表
```http
POST /list/rpush
Content-Type: application/json

{
  "key": "mylist",
  "value": "item2"
}
```

#### 左弹出列表
```http
POST /list/lpop
Content-Type: application/json

{
  "key": "mylist"
}
```

#### 右弹出列表
```http
POST /list/rpop
Content-Type: application/json

{
  "key": "mylist"
}
```

#### 设置列表元素
```http
POST /list/set
Content-Type: application/json

{
  "key": "mylist",
  "index": 0,
  "value": "newitem"
}
```

#### 删除列表元素
```http
POST /list/lrem
Content-Type: application/json

{
  "key": "mylist",
  "value": "item1",
  "count": 1
}
```

### 集合操作

#### 扫描集合
```http
POST /set/scan
Content-Type: application/json

{
  "key": "myset"
}
```

#### 添加集合成员
```http
POST /set/sadd
Content-Type: application/json

{
  "key": "myset",
  "member": "member1"
}
```

#### 删除集合成员
```http
POST /set/srem
Content-Type: application/json

{
  "key": "myset",
  "member": "member1"
}
```

### 有序集合操作

#### 获取有序集合范围
```http
POST /zset/range
Content-Type: application/json

{
  "key": "myzset"
}
```

#### 添加有序集合成员
```http
POST /zset/zadd
Content-Type: application/json

{
  "key": "myzset",
  "member": "member1",
  "score": 1.0
}
```

#### 删除有序集合成员
```http
POST /zset/zrem
Content-Type: application/json

{
  "key": "myzset",
  "member": "member1"
}
```

#### 增加有序集合分数
```http
POST /zset/zincrby
Content-Type: application/json

{
  "key": "myzset",
  "member": "member1",
  "score": 1.5
}
```

## 错误码

| 错误码 | 描述 |
|--------|------|
| 0 | 成功 |
| 1 | 失败 |

## 状态码

| HTTP状态码 | 描述 |
|------------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 500 | 服务器内部错误 |

## 注意事项

1. 所有POST请求的Content-Type必须为`application/json`
2. 路径参数使用`{param}`格式
3. 查询参数使用`?key=value`格式
4. 响应时间通常在100ms以内
5. 支持CORS跨域请求
