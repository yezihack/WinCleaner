<template>
  <div class="process-page">
    <div class="page-header">
      <h2>进程管理</h2>
      <p class="page-sub">管理系统运行中的进程，查看资源占用</p>
    </div>

    <div class="toolbar">
      <button class="refresh-btn" :class="{ loading: loading }" :disabled="loading" @click="loadProcesses">
        <span class="btn-icon">🔄</span>
        <span>{{ loading ? '刷新中...' : '刷新' }}</span>
      </button>
      <div class="search-wrap">
        <span class="search-icon">🔍</span>
        <input v-model="keyword" class="search-input" placeholder="搜索进程名..." />
      </div>
      <span class="proc-count">{{ filteredList.length }} 个进程</span>
    </div>

    <div class="table-wrap">
      <table class="proc-table">
        <thead>
          <tr>
            <th>PID</th>
            <th>进程名</th>
            <th class="col-cpu">CPU</th>
            <th class="col-mem">内存</th>
            <th>用户</th>
            <th class="col-op">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in filteredList" :key="row.pid" class="proc-row">
            <td class="td-pid">{{ row.pid }}</td>
            <td class="td-name">{{ row.name }}</td>
            <td class="td-cpu">
              <div class="bar-cell">
                <div class="bar-fill" :style="{ width: Math.min(row.cpu_percent, 100) + '%', background: cpuColor(row.cpu_percent) }"></div>
                <span class="bar-text">{{ row.cpu_percent.toFixed(1) }}%</span>
              </div>
            </td>
            <td class="td-mem">
              <div class="bar-cell">
                <div class="bar-fill" :style="{ width: Math.min(row.mem_percent, 100) + '%', background: memColor(row.mem_percent) }"></div>
                <span class="bar-text">{{ formatBytes(row.mem_rss) }}</span>
              </div>
            </td>
            <td class="td-user">{{ row.username }}</td>
            <td class="td-op">
              <button class="kill-btn" @click="handleKill(row)">结束</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { api, type ProcessInfo } from '@/api/backend'

const loading = ref(false)
const keyword = ref('')
const processList = ref<ProcessInfo[]>([])

const filteredList = computed(() => {
  const kw = keyword.value.toLowerCase()
  if (!kw) return processList.value
  return processList.value.filter(p => p.name.toLowerCase().includes(kw))
})

const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const cpuColor = (pct: number) => {
  if (pct < 10) return '#22c55e'
  if (pct < 50) return '#f59e0b'
  return '#ef4444'
}

const memColor = (pct: number) => {
  if (pct < 5) return '#22c55e'
  if (pct < 20) return '#f59e0b'
  return '#ef4444'
}

const loadProcesses = async () => {
  loading.value = true
  try {
    processList.value = await api.getProcessList()
  } catch {
    ElMessage.error('获取进程列表失败')
  } finally {
    loading.value = false
  }
}

const handleKill = async (row: ProcessInfo) => {
  try {
    await ElMessageBox.confirm(`确定结束进程 "${row.name}" (PID: ${row.pid})？`, '结束进程', { type: 'warning' })
  } catch { return }
  try {
    await api.killProcess(row.pid)
    ElMessage.success(`已结束进程 ${row.name}`)
    await loadProcesses()
  } catch {
    ElMessage.error('结束进程失败，可能权限不足')
  }
}

onMounted(() => { loadProcesses() })
</script>

<style scoped>
.process-page { padding: 0 4px; }
.page-header { margin-bottom: 20px; }
.page-header h2 { font-size: 22px; font-weight: 600; color: #1a1a2e; margin: 0; }
.page-sub { font-size: 13px; color: #64748b; margin: 4px 0 0; }

.toolbar { display: flex; align-items: center; gap: 12px; margin-bottom: 16px; }

.refresh-btn {
  display: flex; align-items: center; gap: 6px;
  padding: 8px 16px; border: 1px solid #e2e8f0; background: #fff;
  border-radius: 8px; font-size: 13px; color: #1e293b; cursor: pointer; transition: all 0.2s ease;
}
.refresh-btn:hover:not(:disabled) { border-color: #3b82f6; color: #3b82f6; }
.refresh-btn.loading { opacity: 0.6; cursor: not-allowed; }
.btn-icon { font-size: 14px; }

.search-wrap { display: flex; align-items: center; gap: 8px; background: #fff; border: 1px solid #e2e8f0; border-radius: 8px; padding: 0 12px; }
.search-icon { font-size: 14px; color: #94a3b8; }
.search-input { border: none; outline: none; background: transparent; font-size: 13px; color: #1e293b; width: 180px; padding: 8px 0; }

.proc-count { margin-left: auto; font-size: 13px; color: #94a3b8; }

.table-wrap { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; overflow: hidden; }

.proc-table { width: 100%; border-collapse: collapse; font-size: 13px; }
.proc-table th {
  background: #f8fafc; padding: 10px 12px; text-align: left;
  font-weight: 600; color: #64748b; font-size: 12px;
  border-bottom: 1px solid #e2e8f0;
}
.proc-table td { padding: 10px 12px; border-bottom: 1px solid #f8fafc; color: #1e293b; }
.proc-row:last-child td { border-bottom: none; }
.proc-row:hover td { background: #f8fafc; }

.col-cpu, .col-mem { width: 180px; }
.col-op { width: 70px; text-align: center; }

.td-pid { color: #64748b; font-family: monospace; }
.td-name { font-weight: 500; }

.bar-cell { position: relative; height: 20px; background: #f1f5f9; border-radius: 4px; overflow: hidden; }
.bar-fill { height: 100%; border-radius: 4px; transition: width 0.3s ease; }
.bar-text { position: absolute; right: 6px; top: 50%; transform: translateY(-50%); font-size: 11px; color: #475569; }

.td-user { color: #64748b; font-size: 12px; }

.kill-btn {
  padding: 4px 12px; border: 1px solid #fee2e2; background: #fff;
  border-radius: 6px; font-size: 12px; color: #ef4444; cursor: pointer;
  transition: all 0.15s ease;
}
.kill-btn:hover { background: #ef4444; color: #fff; border-color: #ef4444; }
</style>
