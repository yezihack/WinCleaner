<template>
  <div class="port-page">
    <div class="page-header">
      <h2>端口管理</h2>
      <p class="page-sub">查看监听端口及关联进程</p>
    </div>

    <div class="toolbar">
      <button class="refresh-btn" :class="{ loading: loading }" :disabled="loading" @click="fetchAllPorts">
        <span class="btn-icon">🔄</span>
        <span>{{ loading ? '加载中...' : '刷新' }}</span>
      </button>
      <div class="search-wrap">
        <span class="search-icon">🔍</span>
        <input
          v-model="selectedPortInput"
          class="search-input"
          placeholder="输入端口号..."
          type="text"
          inputmode="numeric"
          @keyup.enter.prevent="handleQuery"
        />
      </div>
      <button class="query-btn" @click="handleQuery">查询</button>
      <button class="kill-all-btn" :disabled="portList.length === 0 || killing" @click="handleKillAll">一键结束</button>
      <span class="port-count">{{ portList.length }} 个端口</span>
    </div>

    <div v-if="loadingText" class="loading-bar">
      <span>{{ loadingText }}</span>
    </div>

    <div v-if="queryPort > 0 && !loading" class="query-info">
      端口 {{ queryPort }} 共找到 {{ portList.length }} 个进程
    </div>

    <div class="table-wrap">
      <table class="port-table">
        <thead>
          <tr>
            <th>监听地址</th>
            <th>端口</th>
            <th>协议</th>
            <th>PID</th>
            <th>进程名</th>
            <th class="col-op">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in portList" :key="row.pid + row.listen_addr" class="port-row">
            <td class="td-addr">{{ row.listen_addr }}</td>
            <td class="td-port">{{ row.port }}</td>
            <td><span class="proto-tag">{{ row.proto.toUpperCase() }}</span></td>
            <td class="td-pid">{{ row.pid }}</td>
            <td class="td-name">{{ row.process_name }}</td>
            <td class="td-op">
              <button class="kill-btn" @click="handleKill(row)">结束</button>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="portList.length === 0 && !loading" class="empty-tip">暂无监听中的端口</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { api, type PortInfo } from '@/api/backend'

const loading = ref(false)
const killing = ref(false)
const selectedPortInput = ref('')
const queryPort = ref(0)
const loadingText = ref('')
const portList = ref<PortInfo[]>([])

const fetchAllPorts = async () => {
  loading.value = true
  loadingText.value = '正在加载端口...'
  queryPort.value = 0
  try {
    portList.value = await api.getListeningPorts(1)
  } catch {
    ElMessage.error('加载端口列表失败')
  } finally {
    loading.value = false
    loadingText.value = ''
  }
}

const handleQuery = async () => {
  const val = selectedPortInput.value.trim()
  if (!val) {
    fetchAllPorts()
    return
  }
  const port = parseInt(val, 10)
  if (isNaN(port) || port < 1 || port > 65535) {
    ElMessage.warning('请输入 1-65535 之间的端口号')
    return
  }
  queryPort.value = port
  loading.value = true
  loadingText.value = `正在查询端口 ${port}...`
  portList.value = []
  try {
    portList.value = await api.getPortList(port)
    if (portList.value.length === 0) ElMessage.info(`端口 ${port} 没有监听中的进程`)
  } catch {
    ElMessage.error('查询端口失败')
  } finally {
    loading.value = false
    loadingText.value = ''
  }
}

const handleKill = async (row: PortInfo) => {
  try {
    await ElMessageBox.confirm(`确定结束进程 "${row.process_name}" (PID: ${row.pid})？`, '结束进程', { type: 'warning' })
  } catch { return }
  killing.value = true
  try {
    await api.killProcessesByPort(row.port)
    ElMessage.success({ message: `进程 ${row.process_name}(${row.pid}) 已结束`, duration: 2000 })
    if (queryPort.value > 0) {
      portList.value = portList.value.filter(p => p.pid !== row.pid)
    } else {
      await fetchAllPorts()
    }
  } catch {
    ElMessage.error('结束进程失败，可能权限不足')
  } finally {
    killing.value = false
  }
}

const handleKillAll = async () => {
  if (portList.value.length === 0) return
  const pids = portList.value.map(p => `${p.process_name}(${p.pid})`).join(', ')
  try {
    await ElMessageBox.confirm(`确定结束端口 ${queryPort.value} 下的所有进程？\n${pids}`, '批量结束进程', { type: 'warning' })
  } catch { return }
  killing.value = true
  try {
    await api.killProcessesByPort(queryPort.value)
    ElMessage.success({ message: `端口 ${queryPort.value} 的进程已全部结束`, duration: 2000 })
    portList.value = []
  } catch {
    ElMessage.error('结束进程失败，可能权限不足')
  } finally {
    killing.value = false
  }
}

onMounted(() => {
  fetchAllPorts()
})
</script>

<style scoped>
.port-page { padding: 0 4px; }
.page-header { margin-bottom: 20px; }
.page-header h2 { font-size: 22px; font-weight: 600; color: #1a1a2e; margin: 0; }
.page-sub { font-size: 13px; color: #64748b; margin: 4px 0 0; }

.toolbar { display: flex; align-items: center; gap: 10px; margin-bottom: 14px; }

.refresh-btn {
  display: flex; align-items: center; gap: 6px;
  padding: 8px 14px; border: 1px solid #e2e8f0; background: #fff;
  border-radius: 8px; font-size: 13px; color: #1e293b; cursor: pointer; transition: all 0.2s ease;
}
.refresh-btn:hover:not(:disabled) { border-color: #3b82f6; color: #3b82f6; }
.refresh-btn.loading { opacity: 0.6; cursor: not-allowed; }
.btn-icon { font-size: 14px; }

.search-wrap { display: flex; align-items: center; gap: 6px; background: #fff; border: 1px solid #e2e8f0; border-radius: 8px; padding: 0 10px; }
.search-icon { font-size: 13px; color: #94a3b8; }
.search-input { border: none; outline: none; background: transparent; font-size: 13px; color: #1e293b; width: 120px; padding: 7px 0; font-family: monospace; }

.query-btn {
  padding: 8px 16px; border: none; border-radius: 8px;
  background: #0f172a; color: #fff; font-size: 13px; cursor: pointer; transition: all 0.2s ease;
}
.query-btn:hover { background: #1e293b; }

.kill-all-btn {
  padding: 8px 16px; border: 1px solid #fee2e2; background: #fff;
  border-radius: 8px; font-size: 13px; color: #ef4444; cursor: pointer; transition: all 0.2s ease;
}
.kill-all-btn:hover:not(:disabled) { background: #ef4444; color: #fff; border-color: #ef4444; }
.kill-all-btn:disabled { opacity: 0.5; cursor: not-allowed; }

.port-count { margin-left: auto; font-size: 13px; color: #94a3b8; }

.loading-bar { font-size: 13px; color: #64748b; margin-bottom: 10px; }
.query-info { font-size: 13px; color: #64748b; margin-bottom: 10px; }

.table-wrap { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; overflow: hidden; }

.port-table { width: 100%; border-collapse: collapse; font-size: 13px; }
.port-table th {
  background: #f8fafc; padding: 10px 12px; text-align: left;
  font-weight: 600; color: #64748b; font-size: 12px; border-bottom: 1px solid #e2e8f0;
}
.port-table td { padding: 10px 12px; border-bottom: 1px solid #f8fafc; color: #1e293b; }
.port-row:last-child td { border-bottom: none; }
.port-row:hover td { background: #f8fafc; }

.col-op { width: 70px; text-align: center; }
.td-addr { font-family: monospace; color: #64748b; }
.td-port { font-weight: 600; color: #3b82f6; }
.td-pid { font-family: monospace; color: #64748b; }
.td-name { font-weight: 500; }

.kill-btn {
  padding: 4px 12px; border: 1px solid #fee2e2; background: #fff;
  border-radius: 6px; font-size: 12px; color: #ef4444; cursor: pointer;
  transition: all 0.15s ease;
}
.kill-btn:hover { background: #ef4444; color: #fff; border-color: #ef4444; }

.empty-tip { text-align: center; padding: 40px 0; font-size: 13px; color: #94a3b8; }

.proto-tag {
  display: inline-block;
  padding: 2px 8px;
  background: #f1f5f9;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
  letter-spacing: 0.5px;
}
</style>
