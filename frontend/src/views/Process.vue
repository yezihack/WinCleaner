<template>
  <div class="process">
    <h2 style="margin-bottom: 20px; color: #303133;">进程管理</h2>

    <el-card shadow="hover">
      <div class="action-bar">
        <el-button type="primary" :loading="loading" @click="loadProcesses">
          <el-icon><Refresh /></el-icon> 刷新
        </el-button>
        <el-input
          v-model="keyword"
          placeholder="搜索进程名..."
          clearable
          style="width: 240px;"
          :prefix-icon="Search"
        />
        <span class="process-count">共 {{ filteredList.length }} 个进程</span>
      </div>

      <el-table
        :data="filteredList"
        style="margin-top: 16px;"
        max-height="480"
        :default-sort="{ prop: 'mem_rss', order: 'descending' }"
        stripe
      >
        <el-table-column prop="pid" label="PID" width="80" sortable />
        <el-table-column prop="name" label="进程名" min-width="180" sortable show-overflow-tooltip />
        <el-table-column prop="cpu_percent" label="CPU" width="130" sortable>
          <template #default="{ row }">
            <div class="usage-cell">
              <el-progress
                :percentage="Math.min(row.cpu_percent, 100)"
                :stroke-width="14"
                :show-text="false"
                :color="cpuColor(row.cpu_percent)"
                style="flex: 1;"
              />
              <span class="usage-text">{{ row.cpu_percent.toFixed(1) }}%</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="mem_rss" label="内存" width="180" sortable>
          <template #default="{ row }">
            <div class="usage-cell">
              <el-progress
                :percentage="Math.min(row.mem_percent, 100)"
                :stroke-width="14"
                :show-text="false"
                :color="memColor(row.mem_percent)"
                style="flex: 1;"
              />
              <span class="usage-text">{{ formatBytes(row.mem_rss) }} ({{ row.mem_percent.toFixed(1) }}%)</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="username" label="用户" width="100" show-overflow-tooltip />
        <el-table-column label="操作" width="80" align="center">
          <template #default="{ row }">
            <el-button link type="danger" size="small" @click="handleKill(row)">
              结束
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
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
  if (pct < 10) return '#67c23a'
  if (pct < 50) return '#e6a23c'
  return '#f56c6c'
}

const memColor = (pct: number) => {
  if (pct < 5) return '#67c23a'
  if (pct < 20) return '#e6a23c'
  return '#f56c6c'
}

const loadProcesses = async () => {
  loading.value = true
  try {
    processList.value = await api.getProcessList()
  } catch (e) {
    ElMessage.error('获取进程列表失败')
  } finally {
    loading.value = false
  }
}

const handleKill = async (row: ProcessInfo) => {
  try {
    await ElMessageBox.confirm(
      `确定结束进程 "${row.name}" (PID: ${row.pid})？`,
      '结束进程',
      { type: 'warning' }
    )
  } catch {
    return
  }

  try {
    await api.killProcess(row.pid)
    ElMessage.success(`已结束进程 ${row.name}`)
    await loadProcesses()
  } catch (e) {
    ElMessage.error('结束进程失败，可能权限不足')
  }
}

onMounted(() => {
  loadProcesses()
})
</script>

<style scoped>
.action-bar {
  display: flex;
  align-items: center;
  gap: 12px;
}
.process-count {
  color: #909399;
  font-size: 13px;
  margin-left: auto;
}
.usage-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}
.usage-text {
  font-size: 12px;
  color: #606266;
  white-space: nowrap;
  min-width: 50px;
}
</style>
