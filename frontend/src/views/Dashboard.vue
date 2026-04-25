<template>
  <div class="dashboard-page">
    <div class="page-header">
      <h2>系统概览</h2>
      <p class="page-sub">实时监控您的系统资源状态</p>
    </div>

    <div class="gauge-row">
      <div class="gauge-card">
        <div class="gauge-wrap" :style="{ '--pct': info.cpu_usage + '%', '--color': cpuColor }">
          <svg class="gauge-svg" viewBox="0 0 200 200">
            <circle class="gauge-track" cx="100" cy="100" r="88" />
            <circle class="gauge-fill" cx="100" cy="100" r="88" />
          </svg>
          <div class="gauge-center">
            <div class="gauge-val">{{ Math.round(info.cpu_usage) }}<span>%</span></div>
            <div class="gauge-lbl">CPU</div>
          </div>
        </div>
      </div>

      <div class="gauge-card">
        <div class="gauge-wrap" :style="{ '--pct': info.mem_percent + '%', '--color': memColor }">
          <svg class="gauge-svg" viewBox="0 0 200 200">
            <circle class="gauge-track" cx="100" cy="100" r="88" />
            <circle class="gauge-fill" cx="100" cy="100" r="88" />
          </svg>
          <div class="gauge-center">
            <div class="gauge-val">{{ Math.round(info.mem_percent) }}<span>%</span></div>
            <div class="gauge-lbl">内存</div>
          </div>
        </div>
        <div class="gauge-sub">{{ formatBytes(info.mem_used) }} / {{ formatBytes(info.mem_total) }}</div>
      </div>

      <div class="gauge-card">
        <div class="gauge-wrap" :style="{ '--pct': info.disk_percent + '%', '--color': diskColor }">
          <svg class="gauge-svg" viewBox="0 0 200 200">
            <circle class="gauge-track" cx="100" cy="100" r="88" />
            <circle class="gauge-fill" cx="100" cy="100" r="88" />
          </svg>
          <div class="gauge-center">
            <div class="gauge-val">{{ Math.round(info.disk_percent) }}<span>%</span></div>
            <div class="gauge-lbl">磁盘</div>
          </div>
        </div>
        <div class="gauge-sub">{{ formatBytes(info.disk_used) }} / {{ formatBytes(info.disk_total) }}</div>
      </div>
    </div>

    <div class="info-grid">
      <div class="info-card gpu-card">
        <div class="info-header">
          <h3>显卡信息</h3>
          <el-tag v-if="gpuLoaded && gpus.length === 0" type="info" size="small">未检测到</el-tag>
        </div>
        <div v-if="gpuLoading" class="skeleton-wrap" />
        <div v-else-if="gpus.length > 0" class="gpu-list">
          <div v-for="(gpu, idx) in gpus" :key="idx" class="gpu-item">
            <div class="gpu-top">
              <el-tag :type="gpu.type === 'discrete' ? 'danger' : gpu.type === 'integrated' ? 'warning' : 'info'" size="small" effect="dark">
                {{ gpu.type_label }}
              </el-tag>
              <span class="gpu-name">{{ gpu.name }}</span>
            </div>
            <div class="gpu-meta">
              <div class="meta-item">
                <span class="meta-key">显存</span>
                <span class="meta-val">{{ gpu.vram > 0 ? formatBytes(gpu.vram) : '共享内存' }}</span>
              </div>
              <div class="meta-item">
                <span class="meta-key">驱动</span>
                <span class="meta-val">{{ gpu.driver_ver || '-' }}</span>
              </div>
              <div class="meta-item">
                <span class="meta-key">分辨率</span>
                <span class="meta-val">{{ gpu.resolution || '-' }}</span>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="empty-tip">未检测到独立显卡或核显</div>
      </div>

      <div class="info-card sys-card">
        <div class="info-header">
          <h3>系统信息</h3>
        </div>
        <div class="sys-grid">
          <div class="sys-item">
            <span class="sys-icon">🖥️</span>
            <div class="sys-info">
              <div class="sys-val">{{ info.hostname || '-' }}</div>
              <div class="sys-key">主机名</div>
            </div>
          </div>
          <div class="sys-item">
            <span class="sys-icon">💿</span>
            <div class="sys-info">
              <div class="sys-val">{{ info.os || '-' }}</div>
              <div class="sys-key">操作系统</div>
            </div>
          </div>
          <div class="sys-item">
            <span class="sys-icon">🌐</span>
            <div class="sys-info">
              <div class="sys-val">{{ info.public_ip || '获取中...' }}</div>
              <div class="sys-key">公网 IP</div>
            </div>
          </div>
          <div class="sys-item">
            <span class="sys-icon">📍</span>
            <div class="sys-info">
              <div class="sys-val">{{ info.ip_location || '-' }}</div>
              <div class="sys-key">归属地</div>
            </div>
          </div>
          <div class="sys-item">
            <span class="sys-icon">📡</span>
            <div class="sys-info">
              <div class="sys-val">{{ info.ip_operator || '-' }}</div>
              <div class="sys-key">运营商</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="refresh-row">
      <button class="refresh-btn" :class="{ loading }" :disabled="loading" @click="refresh">
        <span class="refresh-icon">🔄</span>
        <span>{{ loading ? '刷新中...' : '刷新数据' }}</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { api, type SystemInfo, type GPUInfo } from '@/api/backend'

const loading = ref(false)
const gpuLoading = ref(false)
const gpuLoaded = ref(false)
const gpus = ref<GPUInfo[]>([])
const info = reactive<SystemInfo>({
  os: '', hostname: '', cpu_usage: 0,
  mem_total: 0, mem_used: 0, mem_percent: 0,
  disk_total: 0, disk_used: 0, disk_percent: 0,
  public_ip: '', ip_location: '', ip_operator: '',
})

const cpuColor = computed(() => pctColor(info.cpu_usage))
const memColor = computed(() => pctColor(info.mem_percent))
const diskColor = computed(() => pctColor(info.disk_percent))

const pctColor = (pct: number) => {
  if (pct < 50) return '#22c55e'
  if (pct < 80) return '#f59e0b'
  return '#ef4444'
}

const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const refresh = async () => {
  loading.value = true
  try {
    const data = await api.getSystemInfo()
    Object.assign(info, data)
  } catch {
    ElMessage.error('获取系统信息失败')
  } finally {
    loading.value = false
  }
}

const loadGPU = async () => {
  gpuLoading.value = true
  try {
    const result = await api.getGPUInfo()
    gpus.value = result.gpus || []
    gpuLoaded.value = true
  } catch {
    gpus.value = []
    gpuLoaded.value = true
  } finally {
    gpuLoading.value = false
  }
}

onMounted(() => {
  refresh()
  loadGPU()
})
</script>

<style scoped>
.dashboard-page { padding: 0 4px; }

.page-header { margin-bottom: 24px; }
.page-header h2 { font-size: 22px; font-weight: 600; color: #1a1a2e; margin: 0; }
.page-sub { font-size: 13px; color: #64748b; margin: 4px 0 0; }

.gauge-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-bottom: 20px;
}

.gauge-card {
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  border-radius: 16px;
  padding: 28px 20px 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.gauge-wrap {
  position: relative;
  width: 160px;
  height: 160px;
}

.gauge-center {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.gauge-val {
  font-size: 36px;
  font-weight: 700;
  color: #f8fafc;
  line-height: 1;
}

.gauge-val span { font-size: 18px; font-weight: 400; }

.gauge-lbl {
  font-size: 12px;
  color: #94a3b8;
  margin-top: 4px;
}

.gauge-svg {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  transform: rotate(-90deg);
}

.gauge-track {
  fill: none;
  stroke: #334155;
  stroke-width: 10;
}

.gauge-fill {
  fill: none;
  stroke: var(--color, #22c55e);
  stroke-width: 10;
  stroke-linecap: round;
  stroke-dasharray: 553;
  stroke-dashoffset: calc(553 - (553 * var(--pct, 0%) / 100));
  transition: stroke-dashoffset 1s ease, stroke 0.5s ease;
}

.gauge-sub {
  margin-top: 12px;
  font-size: 13px;
  color: #94a3b8;
}

.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
}

.info-card {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 20px;
}

.info-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.info-header h3 {
  font-size: 15px;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.gpu-list { display: flex; flex-direction: column; gap: 12px; }

.gpu-item {
  background: #f8fafc;
  border-radius: 8px;
  padding: 12px 14px;
}

.gpu-top {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.gpu-name {
  font-size: 14px;
  font-weight: 500;
  color: #1e293b;
}

.gpu-meta {
  display: flex;
  gap: 16px;
}

.meta-item { display: flex; flex-direction: column; gap: 2px; }
.meta-key { font-size: 11px; color: #94a3b8; }
.meta-val { font-size: 13px; color: #1e293b; font-weight: 500; }

.empty-tip { font-size: 13px; color: #94a3b8; text-align: center; padding: 20px 0; }

.sys-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }

.sys-item {
  display: flex;
  align-items: center;
  gap: 10px;
  background: #f8fafc;
  border-radius: 8px;
  padding: 10px 12px;
}

.sys-icon { font-size: 20px; }

.sys-val { font-size: 13px; font-weight: 600; color: #1e293b; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.sys-key { font-size: 11px; color: #94a3b8; margin-top: 2px; }

.refresh-row { display: flex; justify-content: flex-end; }

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 20px;
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  color: #1e293b;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.refresh-btn:hover:not(:disabled) {
  border-color: #3b82f6;
  color: #3b82f6;
}

.refresh-btn.loading { opacity: 0.6; cursor: not-allowed; }

.refresh-icon { font-size: 14px; }

.skeleton-wrap { height: 80px; background: linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite; border-radius: 8px; }

@keyframes shimmer { 0% { background-position: 200% 0; } 100% { background-position: -200% 0; } }
</style>
