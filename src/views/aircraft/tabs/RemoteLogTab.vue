<template>
  <div style="padding: 8px 0; color: var(--el-text-color-regular); font-size: 13px; display:flex; flex-direction: column; gap:10px;">
    <!-- 文件上传管理 -->
    <div style="display: flex; flex-direction: column; gap: 10px;">
      
      <!-- 1. 获取设备可上传的文件列表 -->
      <el-card shadow="never">
      <template #header>
        <div style="display: flex; align-items: center; justify-content: space-between; cursor: pointer;" @click="toggleFileListCollapse">
          <div style="display: flex; align-items: center; gap: 8px;">
            <el-icon :style="{ transform: fileListCollapsed ? 'rotate(0deg)' : 'rotate(90deg)', transition: 'transform 0.3s' }">
              <ArrowRight />
            </el-icon>
            <span style="font-weight: 600; font-size: 14px;">1. 获取设备可上传的文件列表</span>
            <el-tooltip 
              v-if="!fileListCollapsed" 
              content="点击展开配置区域，编辑请求参数"
              placement="top"
            >
              <el-tag size="small" type="info" style="cursor: help;">配置参数</el-tag>
            </el-tooltip>
          </div>
        </div>
      </template>
      
      <el-collapse-transition>
        <div v-show="!fileListCollapsed">
          <div style="margin-bottom: 15px;">
            <div style="margin-bottom: 8px; color: var(--el-text-color-regular); font-size: 13px;">
              请求参数 (JSON格式) - 模块ID: 0=飞机类, 1=负载类, 2=遥控器类, 3=机场类
            </div>
            <el-input
              v-model="requestJson"
              type="textarea"
              :rows="8"
              placeholder='请输入JSON格式的请求参数，例如：
{
  "method": "fileupload_list",
  "module_list": ["0", "3"],
  "extra_params": {}
}'
              style="font-family: monospace; font-size: 12px;"
            />
          </div>
          <div style="display: flex; gap: 8px; justify-content: flex-end;">
            <el-button type="primary" :loading="sending" @click="handleRequest">
              发送请求
            </el-button>
            <el-button @click="resetConfig">重置</el-button>
            <el-button @click="loadPreset('fileupload_list')">文件列表</el-button>
          </div>
    </div>
      </el-collapse-transition>
    </el-card>

    <!-- 发起日志文件上传 -->
    <el-card shadow="never" style="margin-top: 10px;">
      <template #header>
        <div style="display: flex; align-items: center; justify-content: space-between; cursor: pointer;" @click="toggleFileUploadCollapse">
          <div style="display: flex; align-items: center; gap: 8px;">
            <el-icon :style="{ transform: fileUploadCollapsed ? 'rotate(0deg)' : 'rotate(90deg)', transition: 'transform 0.3s' }">
              <ArrowRight />
            </el-icon>
            <span style="font-weight: 600; font-size: 14px;">2. 发起日志文件上传</span>
            <el-tooltip 
              v-if="!fileUploadCollapsed" 
              content="点击展开配置区域，编辑上传参数"
              placement="top"
            >
              <el-tag size="small" type="info" style="cursor: help;">配置参数</el-tag>
            </el-tooltip>
          </div>
        </div>
      </template>
      
      <el-collapse-transition>
        <div v-show="!fileUploadCollapsed">
          <div style="margin-bottom: 15px;">
            <div style="margin-bottom: 8px; color: var(--el-text-color-regular); font-size: 13px; display: flex; align-items: center; gap: 8px;">
              <span>上传参数 (JSON格式)</span>
              <el-button 
                size="small" 
                type="primary" 
                text 
                @click="toggleUploadHelp"
                style="padding: 0; font-size: 12px; color: var(--el-color-primary);"
              >
                {{ showUploadHelp ? '隐藏' : '显示' }}字段说明
              </el-button>
            </div>
            
            <!-- 字段说明提示 -->
            <el-collapse-transition>
              <div v-show="showUploadHelp" style="margin-bottom: 10px; padding: 10px; background: var(--el-bg-color-page); border-radius: 4px; border: 1px solid var(--el-border-color-light);">
                <div style="color: var(--el-text-color-regular); font-size: 12px; line-height: 1.6;">
                  <div><strong>字段说明：</strong></div>
                  <div><strong>bucket:</strong> 对象存储桶名称 (text)</div>
                  <div><strong>region:</strong> 数据中心所在的地域 (text)</div>
                  <div><strong>credentials:</strong> 凭证信息 (struct)</div>
                  <div style="margin-left: 20px;">• access_key_id: 访问密钥ID (text)</div>
                  <div style="margin-left: 20px;">• access_key_secret: 秘密访问密钥 (text)</div>
                  <div style="margin-left: 20px;">• expire: 访问密钥过期时间 (int, 秒)</div>
                  <div style="margin-left: 20px;">• security_token: 会话凭证 (text)</div>
                  <div><strong>endpoint:</strong> 对外服务的访问域名 (text)</div>
                  <div><strong>provider:</strong> 云厂商枚举值 (ali=阿里云, aws=亚马逊云, minio=minio)</div>
                  <div><strong>params:</strong> 参数结构 (struct)</div>
                  <div style="margin-left: 20px;"><strong>files:</strong> 文件列表 (array)</div>
                  <div style="margin-left: 40px;">• object_key: 文件在对象存储桶的Key (text)</div>
                  <div style="margin-left: 40px;">• module: 日志所属模块 (text)</div>
                  <div style="margin-left: 40px;">• list: 日志列表 (array)</div>
                  <div style="margin-left: 60px;">• boot_index: 日志索引 (int)</div>
                </div>
              </div>
            </el-collapse-transition>
            <el-input
              v-model="uploadJson"
              type="textarea"
              :rows="12"
              placeholder='请输入JSON格式的上传参数，点击上方"显示字段说明"查看详细参数说明

示例：
{
  "bucket": "your-bucket-name",
  "region": "cn-hangzhou", 
  "credentials": {
    "access_key_id": "your-access-key",
    "access_key_secret": "your-secret-key",
    "expire": 3600,
    "security_token": "your-token"
  },
  "endpoint": "https://oss-cn-hangzhou.aliyuncs.com",
  "provider": "ali",
  "params": {
    "files": [
      {
        "object_key": "logs/device_sn_001/",
        "module": "0",
        "list": [
          {
            "boot_index": 12345
          }
        ]
      }
    ]
  }
}'
              style="font-family: monospace; font-size: 12px;"
            />
          </div>
          <div style="display: flex; gap: 8px; justify-content: flex-end;">
            <el-button type="primary" :loading="uploadSending" @click="handleUpload">
              发起上传
            </el-button>
            <el-button @click="resetUploadConfig">重置</el-button>
            <el-button @click="loadUploadPreset('basic')">基础配置</el-button>
            <el-button @click="loadUploadPreset('minio')">MinIO配置</el-button>
          </div>
    </div>
      </el-collapse-transition>
    </el-card>
    
    <el-collapse-transition>
      <el-row v-show="!fileListCollapsed" :gutter="10">
        <!-- 文件列表请求信息（左侧） -->
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <div style="display:flex; align-items:center; justify-content: space-between;">
                <span>文件列表请求（down）</span>
              <el-tag size="small" type="info">{{ requestTopic || 'topic 未生成' }}</el-tag>
            </div>
          </template>
          <pre v-if="requestPayload" class="json-block">{{ format(requestPayload) }}</pre>
          <div v-else style="color:#909399;">尚未发送请求</div>
        </el-card>
      </el-col>
        
        <!-- 文件列表响应数据（右侧） -->
        <el-col :span="12">
          <el-card shadow="never" v-if="replyPayload">
            <template #header>
              <div style="display:flex; align-items:center; justify-content: space-between;">
                <span>文件列表响应（up）</span>
                <el-tag size="small" type="success">{{ replyTopic || '响应主题' }}</el-tag>
              </div>
            </template>
            
            <!-- 如果有文件列表数据，显示表格 -->
            <div v-if="fileListData.length > 0">
              <div style="margin-bottom: 10px; color: var(--el-text-color-regular); font-size: 13px;">
                文件列表 ({{ fileListData.length }} 个文件)
              </div>
              <el-table 
                :data="fileListData" 
                size="small" 
                height="300" 
                border 
                stripe
                style="margin-bottom: 15px;"
              >
                <el-table-column prop="index" label="序号" width="60" />
                <el-table-column prop="name" label="文件名" min-width="200" />
                <el-table-column prop="size" label="文件大小" width="120">
                  <template #default="scope">
                    {{ formatFileSize(scope.row.size) }}
                  </template>
                </el-table-column>
                <el-table-column prop="type" label="文件类型" width="100">
                  <template #default="scope">
                    <el-tag :type="getFileTypeColor(scope.row.type || '未知')" size="small">
                      {{ scope.row.type || '未知' }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="module" label="模块" width="100">
                  <template #default="scope">
                    <el-tag size="small" type="info" :title="`模块ID: ${scope.row.module || '--'}`">
                      {{ getModuleDescription(scope.row.module) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="timestamp" label="时间戳" width="120">
                  <template #default="scope">
                    {{ formatTimestamp(scope.row.timestamp) }}
                  </template>
                </el-table-column>
              </el-table>
            </div>
            
            <!-- 原始响应数据（保持原格式） -->
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
              <span style="color: var(--el-text-color-regular); font-size: 13px;">原始响应数据</span>
              <el-button size="small" type="primary" @click="showDetailModal">
                查看详细表格
              </el-button>
            </div>
            <pre class="json-block">{{ format(replyPayload) }}</pre>
          </el-card>
          
          <!-- 等待响应状态 -->
          <el-card shadow="never" v-else>
            <template #header>
              <div style="display:flex; align-items:center; justify-content: space-between;">
                <span>文件列表响应（up）</span>
                <el-tag size="small" :type="sending ? 'warning' : 'info'">
                  {{ sending ? '等待响应…' : '未收到响应' }}
                </el-tag>
              </div>
            </template>
            <div style="color:#909399; text-align: center; padding: 20px;">
              {{ sending ? '正在等待服务器响应…' : '请先发送请求获取文件列表' }}
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-collapse-transition>

    <!-- 上传请求和响应 -->
    <el-collapse-transition>
      <el-row v-show="!fileUploadCollapsed" :gutter="10" style="margin-top: 10px;">
        <!-- 上传请求信息（左侧） -->
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <div style="display:flex; align-items:center; justify-content: space-between;">
                <span>上传请求（down）</span>
                <el-tag size="small" type="info">{{ uploadRequestTopic || 'topic 未生成' }}</el-tag>
              </div>
            </template>
            <pre v-if="uploadRequestPayload" class="json-block">{{ format(uploadRequestPayload) }}</pre>
            <div v-else style="color:#909399;">尚未发送上传请求</div>
          </el-card>
        </el-col>
        
        <!-- 上传响应数据（右侧） -->
        <el-col :span="12">
          <el-card shadow="never" v-if="uploadReplyPayload">
            <template #header>
              <div style="display:flex; align-items:center; justify-content: space-between;">
                <span>上传响应（up）</span>
                <el-tag size="small" type="success">{{ uploadReplyTopic || '响应主题' }}</el-tag>
              </div>
            </template>
            
            <!-- 上传结果信息 -->
            <div v-if="uploadResultData.length > 0">
              <div style="margin-bottom: 10px; color: var(--el-text-color-regular); font-size: 13px;">
                上传执行结果
              </div>
              <el-table 
                :data="uploadResultData" 
                size="small" 
                height="200" 
                border 
                stripe
                style="margin-bottom: 15px;"
              >
                <el-table-column prop="result" label="返回码" width="80">
                  <template #default="scope">
                    <el-tag :type="scope.row.result === 0 ? 'success' : 'danger'" size="small">
                      {{ scope.row.result }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="result_text" label="执行状态" width="100">
                  <template #default="scope">
                    <el-tag :type="scope.row.result === 0 ? 'success' : 'danger'" size="small">
                      {{ scope.row.result_text }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="bid" label="业务ID" width="200">
                  <template #default="scope">
                    <span style="font-family: monospace; font-size: 12px;">{{ scope.row.bid }}</span>
                  </template>
                </el-table-column>
                <el-table-column prop="tid" label="事务ID" width="200">
                  <template #default="scope">
                    <span style="font-family: monospace; font-size: 12px;">{{ scope.row.tid }}</span>
                  </template>
                </el-table-column>
                <el-table-column prop="timestamp" label="响应时间" width="180">
                  <template #default="scope">
                    {{ formatTimestamp(scope.row.timestamp) }}
                  </template>
                </el-table-column>
                <el-table-column prop="method" label="方法" width="120">
                  <template #default="scope">
                    <el-tag size="small" type="info">{{ scope.row.method }}</el-tag>
                  </template>
                </el-table-column>
              </el-table>
            </div>
            
            <!-- 原始响应数据 -->
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
              <span style="color: var(--el-text-color-regular); font-size: 13px;">原始响应数据</span>
            </div>
            <pre class="json-block">{{ format(uploadReplyPayload) }}</pre>
          </el-card>
          
          <!-- 等待上传响应状态 -->
          <el-card shadow="never" v-else>
            <template #header>
              <div style="display:flex; align-items:center; justify-content: space-between;">
                <span>上传响应（up）</span>
                <el-tag size="small" :type="uploadSending ? 'warning' : 'info'">
                  {{ uploadSending ? '等待响应…' : '未收到响应' }}
                </el-tag>
            </div>
          </template>
            <div style="color:#909399; text-align: center; padding: 20px;">
              {{ uploadSending ? '正在等待设备响应…' : '请先发送上传请求' }}
            </div>
        </el-card>
      </el-col>
    </el-row>
    </el-collapse-transition>

    <!-- 详细数据模态框 -->
    <el-dialog
      v-model="detailModalVisible"
      title="响应数据详细表格"
      width="80%"
      :before-close="handleCloseModal"
    >
      <div v-if="detailTableData.length > 0">
        <div style="margin-bottom: 15px; color: var(--el-text-color-regular);">
          共找到 {{ detailTableData.length }} 个设备的数据
        </div>
        <el-table 
          :data="detailTableData" 
          size="small" 
          height="500" 
          border 
          stripe
        >
          <el-table-column prop="device_sn" label="设备SN" width="180" />
          <el-table-column prop="module" label="模块" width="100">
            <template #default="scope">
              <el-tag size="small" type="info" :title="`模块ID: ${scope.row.module}`">
                {{ getModuleDescription(scope.row.module) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="result" label="结果" width="80">
            <template #default="scope">
              <el-tag :type="scope.row.result === 0 ? 'success' : 'danger'" size="small">
                {{ scope.row.result === 0 ? '成功' : '失败' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="file_count" label="文件数量" width="100" />
          <el-table-column prop="total_size" label="总大小" width="120">
            <template #default="scope">
              {{ formatFileSize(scope.row.total_size) }}
            </template>
          </el-table-column>
          <el-table-column prop="time_range" label="时间范围" min-width="200">
            <template #default="scope">
              {{ scope.row.time_range }}
            </template>
          </el-table-column>
        </el-table>

        <!-- 文件列表详情 -->
        <el-collapse style="margin-top: 20px;">
          <el-collapse-item 
            v-for="(device, index) in detailTableData" 
            :key="index"
            :title="`设备 ${device.device_sn} (${getModuleDescription(device.module)}) - ${device.file_count} 个文件`"
            :name="index"
          >
            <el-table 
              :data="device.file_list" 
              size="small" 
              border 
              stripe
            >
              <el-table-column prop="boot_index" label="启动索引" width="120" />
              <el-table-column prop="start_time" label="开始时间" width="180">
                <template #default="scope">
                  {{ formatTimestamp(scope.row.start_time) }}
                </template>
              </el-table-column>
              <el-table-column prop="end_time" label="结束时间" width="180">
                <template #default="scope">
                  {{ formatTimestamp(scope.row.end_time) }}
                </template>
              </el-table-column>
              <el-table-column prop="size" label="文件大小" width="120">
                <template #default="scope">
                  {{ formatFileSize(scope.row.size) }}
                </template>
              </el-table-column>
              <el-table-column prop="duration" label="持续时间" width="120">
                <template #default="scope">
                  {{ formatDuration(scope.row.start_time, scope.row.end_time) }}
                </template>
              </el-table-column>
            </el-table>
          </el-collapse-item>
        </el-collapse>
      </div>
      
      <div v-else style="text-align: center; padding: 40px; color: var(--el-text-color-regular);">
        暂无数据或数据格式不正确
      </div>

      <template #footer>
        <el-button @click="detailModalVisible = false">关闭</el-button>
        <el-button type="primary" @click="exportToCSV">导出CSV</el-button>
      </template>
    </el-dialog>

    <!-- 上传状态更新 -->
    <el-card style="margin-top: 15px;">
      <template #header>
        <div style="display: flex; align-items: center; justify-content: space-between; cursor: pointer;" @click="toggleStatusUpdateCollapse">
          <div style="display: flex; align-items: center; gap: 8px;">
            <el-icon :style="{ transform: statusUpdateCollapsed ? 'rotate(0deg)' : 'rotate(90deg)', transition: 'transform 0.3s' }">
              <ArrowRight />
            </el-icon>
            <span style="font-weight: 600; font-size: 14px;">3. 上传状态更新</span>
            <el-tooltip 
              v-if="!statusUpdateCollapsed" 
              content="点击展开配置区域，编辑状态更新参数"
              placement="top"
            >
              <el-tag size="small" type="info" style="cursor: help;">配置参数</el-tag>
            </el-tooltip>
          </div>
        </div>
      </template>
      
      <el-collapse-transition>
        <div v-show="!statusUpdateCollapsed">
          <div style="margin-bottom: 15px;">
            <div style="margin-bottom: 8px; color: var(--el-text-color-regular); font-size: 13px; display: flex; align-items: center; gap: 8px;">
              <span>状态更新参数 (JSON格式)</span>
              <el-button 
                size="small" 
                type="primary" 
                text 
                @click="toggleStatusUpdateHelp"
                style="padding: 0; font-size: 12px; color: var(--el-color-primary);"
              >
                {{ showStatusUpdateHelp ? '隐藏' : '显示' }}字段说明
              </el-button>
            </div>
            
            <!-- 字段说明提示 -->
            <el-collapse-transition>
              <div v-show="showStatusUpdateHelp" style="margin-bottom: 10px; padding: 10px; background: var(--el-bg-color-page); border-radius: 4px; border: 1px solid var(--el-border-color-light);">
                <div style="color: var(--el-text-color-regular); font-size: 12px; line-height: 1.6;">
                  <div><strong>字段说明：</strong></div>
                  <div><strong>status:</strong> 上传状态 (enum_string) - "cancel": 取消</div>
                  <div><strong>module_list:</strong> 日志所属模块列表 (array) - 模块ID数组</div>
                </div>
              </div>
            </el-collapse-transition>
            
            <el-input
              v-model="statusUpdateJson"
              type="textarea"
              :rows="8"
              placeholder='请输入JSON格式的状态更新参数

示例：
{
  "status": "cancel",
  "module_list": ["0", "3"]
}'
              style="font-family: monospace; font-size: 12px;"
            />
          </div>
          <div style="display: flex; gap: 8px; justify-content: flex-end;">
            <el-button type="primary" :loading="statusUpdateSending" @click="handleStatusUpdate">
              发送状态更新
            </el-button>
            <el-button @click="resetStatusUpdateConfig">
              重置
            </el-button>
            <el-button @click="loadStatusUpdatePreset">
              加载预设
            </el-button>
          </div>
        </div>
      </el-collapse-transition>
    </el-card>

    <!-- 文件上传进度通知 -->
    <el-card style="margin-top: 15px;">
      <template #header>
        <div style="display: flex; align-items: center; justify-content: space-between; cursor: pointer;" @click="toggleProgressCollapse">
          <div style="display: flex; align-items: center; gap: 8px;">
            <el-icon :style="{ transform: progressCollapsed ? 'rotate(0deg)' : 'rotate(90deg)', transition: 'transform 0.3s' }">
              <ArrowRight />
            </el-icon>
            <span style="font-weight: 600; font-size: 14px;">4. 文件上传进度通知</span>
            <div v-if="!progressCollapsed" style="display: flex; align-items: center; gap: 8px;">
              <el-tag size="small" type="success">实时监控</el-tag>
              <el-tag size="small" type="info">Topic: thing/product/{gateway_sn}/events</el-tag>
              <el-tag size="small" type="warning">Direction: up</el-tag>
              <el-tag size="small" type="primary">Method: fileupload_progress</el-tag>
            </div>
          </div>
        </div>
      </template>
      
      <el-collapse-transition>
        <div v-show="!progressCollapsed">
          <!-- 方法说明 -->
          <div style="margin-bottom: 15px; padding: 12px; background: var(--el-bg-color-page); border-radius: 4px; border: 1px solid var(--el-border-color-light);">
            <div style="color: var(--el-text-color-regular); font-size: 13px; line-height: 1.6;">
              <div style="margin-bottom: 8px;"><strong>方法说明：</strong></div>
              <div><strong>Topic:</strong> thing/product/{gateway_sn}/events</div>
              <div><strong>Direction:</strong> up (设备端向服务端发送)</div>
              <div><strong>Method:</strong> fileupload_progress</div>
          <div style="margin-top: 8px; color: var(--el-text-color-placeholder); font-size: 12px;">
            请手动订阅此Topic来接收实时的文件上传进度通知
          </div>
          
          <!-- 手动订阅按钮 -->
          <div style="margin-top: 10px;">
            <el-button 
              type="success" 
              @click="subscribeProgressEvents"
              :loading="progressSubscribing"
              size="small"
            >
              <el-icon><Connection /></el-icon>
              手动订阅进度通知
            </el-button>
          </div>
            </div>
          </div>
          
          <!-- 进度数据表格 -->
          <div v-if="progressData.length > 0" style="margin-bottom: 15px;">
            <el-table :data="progressData" style="width: 100%" size="small">
              <el-table-column prop="device_sn" label="设备SN" width="120" />
              <el-table-column prop="module" label="模块" width="80">
                <template #default="scope">
                  {{ getModuleDescription(scope.row.module) }}
                </template>
              </el-table-column>
              <el-table-column prop="key" label="存储Key" show-overflow-tooltip />
              <el-table-column prop="size" label="文件大小" width="100">
                <template #default="scope">
                  {{ formatFileSize(scope.row.size) }}
                </template>
              </el-table-column>
              <el-table-column prop="progress.progress" label="进度" width="120">
                <template #default="scope">
                  <div style="display: flex; align-items: center; gap: 8px;">
                    <el-progress 
                      :percentage="scope.row.progress?.progress || 0" 
                      :status="getProgressStatus(scope.row.progress)"
                      :stroke-width="6"
                      style="flex: 1;"
                    />
                    <span style="font-size: 12px; color: var(--el-text-color-regular);">
                      {{ scope.row.progress?.progress || 0 }}%
                    </span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="progress.status" label="状态" width="80">
                <template #default="scope">
                  <el-tag 
                    :type="getStatusType(scope.row.progress?.status)" 
                    size="small"
                  >
                    {{ scope.row.progress?.status || '未知' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="progress.upload_rate" label="上传速率" width="100">
                <template #default="scope">
                  {{ formatUploadRate(scope.row.progress?.upload_rate) }}
                </template>
              </el-table-column>
              <el-table-column prop="progress.finish_time" label="完成时间" width="120">
                <template #default="scope">
                  {{ formatTime(scope.row.progress?.finish_time) }}
                </template>
              </el-table-column>
              <el-table-column prop="progress.current_step" label="步骤" width="100">
                <template #default="scope">
                  <div style="font-size: 12px;">
                    {{ scope.row.progress?.current_step || 0 }}/{{ scope.row.progress?.total_step || 0 }}
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="result" label="返回码" width="120">
                <template #default="scope">
                  <div v-if="scope.row.result !== null && scope.row.result !== undefined">
                    <el-tag 
                      :type="scope.row.result === 0 ? 'success' : 'danger'" 
                      size="small"
                    >
                      {{ scope.row.result }}
                    </el-tag>
                    <div v-if="scope.row.result !== 0" style="font-size: 11px; color: var(--el-color-danger); margin-top: 2px;">
                      上传失败
                    </div>
                  </div>
                  <div v-else style="color: var(--el-text-color-placeholder); font-size: 12px;">
                    --
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </div>
          
          <!-- 无数据提示 -->
          <div v-else style="text-align: center; padding: 40px; color: var(--el-text-color-regular);">
            <div style="margin-bottom: 10px;">
              <el-icon size="48" color="var(--el-color-info)">
                <Document />
              </el-icon>
            </div>
            <div style="font-size: 14px; margin-bottom: 8px;">暂无上传进度数据</div>
            <div style="font-size: 12px; color: var(--el-text-color-placeholder);">
              请先发起文件上传，系统将自动接收进度通知
            </div>
          </div>
          
          
          <!-- 最新fileupload_progress消息展示 -->
          <div v-if="allEventsMessages.length > 0" style="margin-top: 15px;">
            <div style="margin-bottom: 8px; color: var(--el-text-color-regular); font-size: 13px;">
              最新文件上传进度消息：
            </div>
            <div v-for="(msg, index) in allEventsMessages" :key="index" style="margin-bottom: 10px; padding: 8px; background: var(--el-bg-color-page); border-radius: 4px; border: 1px solid var(--el-border-color-light);">
              <div style="font-size: 12px; color: var(--el-text-color-regular); margin-bottom: 4px;">
                <strong>方法:</strong> {{ msg.method || '未知' }} | 
                <strong>时间:</strong> {{ formatTime(msg.timestamp) }} | 
                <strong>网关:</strong> {{ msg.gateway || '未知' }}
              </div>
              <pre style="font-size: 11px; margin: 0; max-height: 200px; overflow-y: auto;">{{ format(msg) }}</pre>
            </div>
          </div>
          
          <!-- 调试信息 -->
          <div style="margin-top: 15px; padding: 10px; background: var(--el-bg-color-page); border-radius: 4px; border: 1px solid var(--el-border-color-light);">
            <div style="font-size: 12px; color: var(--el-text-color-regular); margin-bottom: 8px;">
              <strong>调试信息：</strong>
            </div>
            <div style="font-size: 11px; color: var(--el-text-color-placeholder); line-height: 1.4;">
              <div>progressPayload: {{ progressPayload ? '有数据' : '无数据' }}</div>
              <div>最新进度消息: {{ allEventsMessages.length }}条</div>
              <div>progressData: {{ progressData.length }}条</div>
              <div v-if="progressPayload">progressPayload类型: {{ typeof progressPayload }}</div>
            </div>
  </div>
  
          <!-- 无原始数据提示 -->
          <div v-if="!progressPayload && allEventsMessages.length === 0" style="margin-top: 15px; text-align: center; padding: 20px; color: var(--el-text-color-placeholder); font-size: 12px;">
            等待接收events主题消息...
          </div>
        </div>
      </el-collapse-transition>
    </el-card>
    
    </div> <!-- 文件上传管理区域结束 -->
  </div>
</template>

<script setup>
import { computed, ref, reactive } from 'vue'
import { ElButton, ElTag, ElCard, ElTable, ElTableColumn, ElForm, ElFormItem, ElInput, ElRow, ElCol, ElDialog, ElCollapse, ElCollapseItem, ElIcon, ElCollapseTransition, ElTooltip, ElProgress } from 'element-plus'
import { ElMessage } from 'element-plus'
import { ArrowRight, Document, Connection } from '@element-plus/icons-vue'
import { useDeviceStore } from '@/stores/device'
import { useMqttProxyStore } from '@/stores/mqtt-proxy'

const props = defineProps({
  sending: Boolean,
  airportSn: String,
  requestTopic: String,
  replyTopic: String,
  requestPayload: [Object, String],
  replyPayload: [Object, String],
  uploadSending: Boolean,
  uploadRequestTopic: String,
  uploadReplyTopic: String,
  uploadRequestPayload: [Object, String],
  uploadReplyPayload: [Object, String],
  progressPayload: [Object, String],
  allEventsMessages: Array
})

const emit = defineEmits(['request', 'upload', 'statusUpdate', 'subscribeProgress'])

// 模态框状态
const detailModalVisible = ref(false)

// 折叠状态
const fileListCollapsed = ref(true)
const fileUploadCollapsed = ref(true)
const showUploadHelp = ref(false)
const statusUpdateCollapsed = ref(true)
const showStatusUpdateHelp = ref(false)
const progressCollapsed = ref(true)

// 请求JSON
const requestJson = ref(`{
  "method": "fileupload_list",
  "module_list": ["0", "3"],
  "extra_params": {}
}`)

// 上传JSON
const uploadJson = ref(`{
  "bucket": "your-bucket-name",
  "region": "cn-hangzhou", 
  "credentials": {
    "access_key_id": "your-access-key",
    "access_key_secret": "your-secret-key",
    "expire": 3600,
    "security_token": "your-token"
  },
  "endpoint": "https://oss-cn-hangzhou.aliyuncs.com",
  "provider": "ali",
  "params": {
    "files": [
      {
        "object_key": "logs/device_sn_001/",
        "module": "0",
        "list": [
          {
            "boot_index": 12345
          }
        ]
      }
    ]
  }
}`)

// 状态更新JSON
const statusUpdateJson = ref(`{
  "status": "cancel",
  "module_list": ["0", "3"]
}`)

// 上传状态
const uploadSending = ref(false)
const statusUpdateSending = ref(false)

// 进度订阅状态
const progressSubscribing = ref(false)

// 上传结果数据
const uploadResultData = computed(() => {
  if (!props.uploadReplyPayload) return []
  
  try {
    const data = typeof props.uploadReplyPayload === 'string' 
      ? JSON.parse(props.uploadReplyPayload) 
      : props.uploadReplyPayload
    
    if (data && data.data) {
      return [{
        result: data.data.result,
        result_text: data.data.result === 0 ? '成功' : '失败',
        bid: data.bid || '--',
        tid: data.tid || '--',
        timestamp: data.timestamp || Date.now(),
        method: data.method || '--'
      }]
    }
    
    return []
  } catch (error) {
    console.error('解析上传响应数据失败:', error)
    return []
  }
})

// 预设配置
const presets = {
  fileupload_list: `{
  "method": "fileupload_list",
  "module_list": ["0", "3"],
  "extra_params": {}
}` // 模块ID: 0=飞机类, 3=机场类
}

// 上传预设配置
const uploadPresets = {
  basic: `{
  "bucket": "your-bucket-name",
  "region": "cn-hangzhou", 
  "credentials": {
    "access_key_id": "your-access-key",
    "access_key_secret": "your-secret-key",
    "expire": 3600,
    "security_token": "your-token"
  },
  "endpoint": "https://oss-cn-hangzhou.aliyuncs.com",
  "provider": "ali",
  "params": {
    "files": [
      {
        "object_key": "logs/device_sn_001/",
        "module": "0",
        "list": [
          {
            "boot_index": 12345
          }
        ]
      }
    ]
  }
}`,
  minio: `{
  "bucket": "cloud-bucket",
  "region": "default", 
  "credentials": {
    "access_key_id": "fCgOzxx2SzRIGnXaOCp6",
    "access_key_secret": "RvUTCzIBFwj1bSFJIxnNudK59wTdJl68Je9rZ04x",
    "expire": 3600,
    "security_token": ""
  },
  "endpoint": "http://121.229.184.48:1300",
  "provider": "minio",
  "params": {
    "files": [
      {
        "object_key": "logs/device_sn_001/",
        "module": "0",
        "list": [
          {
            "boot_index": 12345
          }
        ]
      }
    ]
  }
}`
}

// 处理请求
const handleRequest = () => {
  try {
    // 验证JSON格式
    let payload
    try {
      payload = JSON.parse(requestJson.value)
    } catch (error) {
      ElMessage.error('JSON格式错误，请检查语法：' + error.message)
      return
    }

    // 验证必要字段
    if (!payload.method) {
      ElMessage.error('缺少必要字段：method')
      return
    }

    // 发送请求事件，传递解析好的参数
    emit('request', payload)
  } catch (error) {
    ElMessage.error('请求参数处理失败：' + error.message)
  }
}

// 重置配置
const resetConfig = () => {
  requestJson.value = `{
  "method": "fileupload_list",
  "module_list": ["0", "3"],
  "extra_params": {}
}`
}

// 加载预设配置
const loadPreset = (presetName) => {
  if (presets[presetName]) {
    requestJson.value = presets[presetName]
    ElMessage.success(`已加载预设配置：${presetName}`)
  }
}

// 切换折叠状态
const toggleFileListCollapse = () => {
  fileListCollapsed.value = !fileListCollapsed.value
}

// 切换上传折叠状态
const toggleFileUploadCollapse = () => {
  fileUploadCollapsed.value = !fileUploadCollapsed.value
}

// 切换上传帮助显示
const toggleUploadHelp = () => {
  showUploadHelp.value = !showUploadHelp.value
}

// 切换状态更新折叠状态
const toggleStatusUpdateCollapse = () => {
  statusUpdateCollapsed.value = !statusUpdateCollapsed.value
}

// 切换状态更新帮助显示
const toggleStatusUpdateHelp = () => {
  showStatusUpdateHelp.value = !showStatusUpdateHelp.value
}

// 切换进度折叠状态
const toggleProgressCollapse = () => {
  progressCollapsed.value = !progressCollapsed.value
}

// 手动订阅进度通知
const subscribeProgressEvents = async () => {
  try {
    progressSubscribing.value = true
    
    // 使用传入的机场序列号
    const gatewaySn = props.airportSn
    
    if (!gatewaySn || gatewaySn === '--') {
      ElMessage.error('缺少网关SN（机场SN）')
      return
    }
    
    // 获取MQTT代理store
    const mqttProxyStore = useMqttProxyStore()
    
    if (!mqttProxyStore.isConnected) {
      ElMessage.error('MQTT未连接，请先连接MQTT服务')
      return
    }
    
    // 订阅进度通知Topic
    const progressTopic = `thing/product/${gatewaySn}/events`
    console.log('=== 手动订阅进度通知Topic ===')
    console.log('Topic:', progressTopic)
    console.log('Gateway SN:', gatewaySn)
    
    // 这里需要调用父组件的方法来订阅
    // 由于RemoteLogTab是子组件，我们需要通过emit通知父组件
    emit('subscribeProgress', { topic: progressTopic, gatewaySn })
    
    ElMessage.success('已发送订阅请求，请查看控制台确认订阅结果')
    console.log('📡 订阅请求已发送，等待父组件处理...')
    
  } catch (error) {
    console.error('订阅进度通知失败:', error)
    ElMessage.error('订阅失败: ' + (error?.message || '未知错误'))
  } finally {
    progressSubscribing.value = false
  }
}

// 处理上传请求
const handleUpload = () => {
  try {
    // 验证JSON格式
    let payload
    try {
      payload = JSON.parse(uploadJson.value)
    } catch (error) {
      ElMessage.error('JSON格式错误，请检查语法：' + error.message)
      return
    }

    // 验证必要字段
    if (!payload.bucket || !payload.region || !payload.credentials) {
      ElMessage.error('缺少必要字段：bucket, region, credentials')
      return
    }

    // 发送上传请求事件
    emit('upload', payload)
    // 注意：uploadSending状态由父组件AircraftDetail.vue管理
  } catch (error) {
    ElMessage.error('上传参数处理失败：' + error.message)
  }
}

// 重置上传配置
const resetUploadConfig = () => {
  uploadJson.value = `{
  "bucket": "your-bucket-name",
  "region": "cn-hangzhou",
  "credentials": {
    "access_key_id": "your-access-key",
    "access_key_secret": "your-secret-key",
    "expire": 3600,
    "security_token": "your-token"
  },
  "endpoint": "https://oss-cn-hangzhou.aliyuncs.com",
  "provider": "ali",
  "params": {
    "files": [
      {
        "object_key": "logs/device_sn_001/",
        "module": "0",
        "list": [
          {
            "boot_index": 12345
          }
        ]
      }
    ]
  }
}`
}

// 加载上传预设配置
const loadUploadPreset = (presetName) => {
  if (uploadPresets[presetName]) {
    uploadJson.value = uploadPresets[presetName]
    ElMessage.success(`已加载上传预设配置：${presetName}`)
  }
}

// 处理状态更新请求
const handleStatusUpdate = () => {
  try {
    // 验证JSON格式
    const parsed = JSON.parse(statusUpdateJson.value)
    
    // 验证必要字段
    if (!parsed.status) {
      ElMessage.error('缺少必要字段: status')
      return
    }
    if (!parsed.module_list || !Array.isArray(parsed.module_list)) {
      ElMessage.error('缺少必要字段: module_list (数组)')
      return
    }
    
    // 发送状态更新请求
    emit('statusUpdate', parsed)
    ElMessage.success('状态更新请求已发送')
  } catch (e) {
    ElMessage.error('JSON格式错误: ' + e.message)
  }
}

// 重置状态更新配置
const resetStatusUpdateConfig = () => {
  statusUpdateJson.value = `{
  "status": "cancel",
  "module_list": ["0", "3"]
}`
}

// 状态更新预设
const statusUpdatePresets = {
  cancel: {
    status: "cancel",
    module_list: ["0", "3"]
  }
}

// 进度数据
const progressData = computed(() => {
  console.log('=== progressData computed 被调用 ===')
  console.log('props.progressPayload:', props.progressPayload)
  console.log('progressPayload类型:', typeof props.progressPayload)
  
  if (!props.progressPayload) {
    console.log('❌ progressPayload为空')
    return []
  }
  
  try {
    const data = typeof props.progressPayload === 'string' 
      ? JSON.parse(props.progressPayload) 
      : props.progressPayload
    
    console.log('解析后的数据:', data)
    console.log('数据结构检查:')
    console.log('- data:', !!data)
    console.log('- data.data:', !!data?.data)
    console.log('- data.data.result:', data?.data?.result)
    console.log('- 第一个文件的progress.result:', data?.data?.output?.ext?.files?.[0]?.progress?.result)
    console.log('- data.data.output:', !!data?.data?.output)
    console.log('- data.data.output.ext:', !!data?.data?.output?.ext)
    console.log('- data.data.output.ext.files:', !!data?.data?.output?.ext?.files)
    
    if (data && data.data && data.data.output && data.data.output.ext && data.data.output.ext.files) {
      console.log('✅ 数据结构正确，开始处理文件列表')
      console.log('文件列表:', data.data.output.ext.files)
      
      const result = data.data.output.ext.files.map(file => ({
        device_sn: file.device_sn || '--',
        module: file.module || '--',
        key: file.key || '--',
        size: file.size || 0,
        fingerprint: file.fingerprint || '--',
        result: file.progress?.result || data.data.result || null, // 优先使用file.progress.result
        progress: {
          current_step: file.progress?.current_step || 0,
          progress: file.progress?.progress || 0,
          status: file.progress?.status || 'unknown',
          total_step: file.progress?.total_step || 0,
          upload_rate: file.progress?.upload_rate || 0,
          finish_time: file.progress?.finish_time || null,
          result: file.progress?.result || null
        }
      }))
      
      console.log('处理后的结果:', result)
      console.log('每个文件的result字段:')
      result.forEach((file, index) => {
        console.log(`文件${index}: result = ${file.result}`)
      })
      return result
    } else {
      console.log('❌ 数据结构不正确')
      return []
    }
  } catch (e) {
    console.error('❌ 解析进度数据失败:', e)
    return []
  }
})

// 加载状态更新预设
const loadStatusUpdatePreset = () => {
  statusUpdateJson.value = JSON.stringify(statusUpdatePresets.cancel, null, 2)
}

// 详细表格数据
const detailTableData = computed(() => {
  if (!props.replyPayload) return []
  
  try {
    const data = typeof props.replyPayload === 'string' 
      ? JSON.parse(props.replyPayload) 
      : props.replyPayload
    
    if (data && data.data && data.data.files && Array.isArray(data.data.files)) {
      return data.data.files.map(device => {
        const fileList = device.list || []
        const totalSize = fileList.reduce((sum, file) => sum + (file.size || 0), 0)
        const startTimes = fileList.map(f => f.start_time).filter(t => t)
        const endTimes = fileList.map(f => f.end_time).filter(t => t)
        const minTime = startTimes.length ? Math.min(...startTimes) : 0
        const maxTime = endTimes.length ? Math.max(...endTimes) : 0
        
        return {
          device_sn: device.device_sn || '未知设备',
          module: device.module || '--',
          result: device.result || -1,
          file_count: fileList.length,
          total_size: totalSize,
          time_range: minTime && maxTime ? 
            `${formatTimestamp(minTime)} ~ ${formatTimestamp(maxTime)}` : '--',
          file_list: fileList.map(file => ({
            ...file,
            duration: file.start_time && file.end_time ? 
              formatDuration(file.start_time, file.end_time) : '--'
          }))
        }
      })
    }
    
    return []
  } catch (error) {
    console.error('解析详细数据失败:', error)
    return []
  }
})

// 显示详细模态框
const showDetailModal = () => {
  if (!props.replyPayload) {
    ElMessage.warning('暂无响应数据')
    return
  }
  detailModalVisible.value = true
}

// 关闭模态框
const handleCloseModal = () => {
  detailModalVisible.value = false
}

// 格式化持续时间
const formatDuration = (startTime, endTime) => {
  if (!startTime || !endTime) return '--'
  const duration = endTime - startTime
  const seconds = Math.floor(duration / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)
  
  if (days > 0) return `${days}天${hours % 24}小时`
  if (hours > 0) return `${hours}小时${minutes % 60}分钟`
  if (minutes > 0) return `${minutes}分钟${seconds % 60}秒`
  return `${seconds}秒`
}

// 导出CSV
const exportToCSV = () => {
  if (detailTableData.value.length === 0) {
    ElMessage.warning('暂无数据可导出')
    return
  }
  
  try {
    // 构建CSV数据
    let csvContent = '设备SN,模块,结果,文件数量,总大小,时间范围\n'
    
    detailTableData.value.forEach(device => {
      csvContent += `"${device.device_sn}","${device.module}","${device.result === 0 ? '成功' : '失败'}","${device.file_count}","${formatFileSize(device.total_size)}","${device.time_range}"\n`
      
      // 添加文件详情
      device.file_list.forEach(file => {
        csvContent += `  ,"文件","","","","启动索引:${file.boot_index},开始:${formatTimestamp(file.start_time)},结束:${formatTimestamp(file.end_time)},大小:${formatFileSize(file.size)}"\n`
      })
    })
    
    // 下载文件
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')
    const url = URL.createObjectURL(blob)
    link.setAttribute('href', url)
    link.setAttribute('download', `文件列表_${new Date().toISOString().slice(0, 19).replace(/:/g, '-')}.csv`)
    link.style.visibility = 'hidden'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    
    ElMessage.success('CSV文件已导出')
  } catch (error) {
    ElMessage.error('导出失败：' + error.message)
  }
}

// 计算文件列表数据
const fileListData = computed(() => {
  if (!props.replyPayload) return []
  
  try {
    const data = typeof props.replyPayload === 'string' 
      ? JSON.parse(props.replyPayload) 
      : props.replyPayload
    
    // 假设响应数据结构包含文件列表
    if (data && data.result === 0 && data.data && Array.isArray(data.data)) {
      return data.data.map((file, index) => {
        const fileName = file?.name || file?.file_name || '未知文件'
        return {
          index: index + 1,
          name: fileName,
          size: file?.size || file?.file_size || 0,
          type: getFileType(fileName),
          module: file?.module || file?.module_id || '--',
          timestamp: file?.timestamp || file?.create_time || Date.now()
        }
      })
    }
    
    // 如果没有标准结构，尝试其他可能的字段
    if (data && data.files && Array.isArray(data.files)) {
      return data.files.map((file, index) => {
        const fileName = file?.name || file?.file_name || '未知文件'
        return {
          index: index + 1,
          name: fileName,
          size: file?.size || file?.file_size || 0,
          type: getFileType(fileName),
          module: file?.module || file?.module_id || '--',
          timestamp: file?.timestamp || file?.create_time || Date.now()
        }
      })
    }
    
    return []
  } catch (error) {
    console.error('解析响应数据失败:', error)
    return []
  }
})

// 获取文件类型
const getFileType = (filename) => {
  if (!filename || typeof filename !== 'string') return '未知'
  const ext = filename.split('.').pop()?.toLowerCase()
  if (!ext) return '文件'
  const typeMap = {
    'log': '日志',
    'txt': '文本',
    'json': 'JSON',
    'xml': 'XML',
    'csv': 'CSV',
    'zip': '压缩包',
    'tar': '压缩包',
    'gz': '压缩包',
    'jpg': '图片',
    'jpeg': '图片',
    'png': '图片',
    'gif': '图片',
    'mp4': '视频',
    'avi': '视频',
    'mov': '视频'
  }
  return typeMap[ext] || '文件'
}

// 模块ID说明
const MODULE_DESCRIPTIONS = {
  "0": "飞机类",
  "1": "负载类", 
  "2": "遥控器类",
  "3": "机场类"
}

// 获取模块描述
const getModuleDescription = (moduleId) => {
  return MODULE_DESCRIPTIONS[moduleId] || `模块${moduleId}`
}

// 获取文件类型颜色
const getFileTypeColor = (type) => {
  const colorMap = {
    '日志': 'warning',
    '文本': 'info',
    'JSON': 'success',
    'XML': 'primary',
    'CSV': 'success',
    '压缩包': 'danger',
    '图片': 'success',
    '视频': 'primary',
    '文件': 'info'
  }
  return colorMap[type] || 'info'
}

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let size = bytes
  let unitIndex = 0
  
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex++
  }
  
  return `${size.toFixed(2)} ${units[unitIndex]}`
}

// 格式化上传速率
const formatUploadRate = (rate) => {
  if (!rate || rate === 0) return '0 B/s'
  return formatFileSize(rate) + '/s'
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return '--'
  return new Date(timestamp).toLocaleString()
}

// 获取进度状态
const getProgressStatus = (progress) => {
  if (!progress) return 'exception'
  if (progress.status === 'ok' || progress.status === 'completed') return 'success'
  if (progress.status === 'error' || progress.status === 'failed') return 'exception'
  if (progress.status === 'file_zip' || progress.status === 'uploading') return ''
  return ''
}

// 获取状态类型
const getStatusType = (status) => {
  if (status === 'ok') return 'success'
  if (status === 'error') return 'danger'
  if (status === 'file_zip') return 'warning'
  if (status === 'uploading') return 'primary'
  if (status === 'completed') return 'success'
  if (status === 'failed') return 'danger'
  return 'info'
}

// 格式化时间戳
const formatTimestamp = (timestamp) => {
  if (!timestamp) return '--'
  const date = new Date(timestamp * 1000) // 假设是秒级时间戳
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

const format = (obj) => {
  try {
    if (typeof obj === 'string') return JSON.stringify(JSON.parse(obj), null, 2)
    return JSON.stringify(obj, null, 2)
  } catch {
    return String(obj)
  }
}
</script>

<style scoped>
.json-block {
  margin: 0;
  font-family: 'Fira Code', 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-word;
}
</style>


