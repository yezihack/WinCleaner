<template>
  <div class="app-shell">
    <aside class="sidebar">
      <div class="sidebar-top">
        <div class="brand">
          <div class="brand-icon">⚡</div>
          <div class="brand-text">
            <span class="brand-name">WinCleaner</span>
            <span class="brand-sub">系统清理工具</span>
          </div>
        </div>

        <nav class="nav-list">
          <router-link v-for="item in navItems" :key="item.path" :to="item.path" class="nav-item" :class="{ active: route.path === item.path }">
            <span class="nav-icon">{{ item.icon }}</span>
            <span class="nav-label">{{ item.label }}</span>
          </router-link>
        </nav>
      </div>

      <div class="sidebar-bottom">
        <div class="quick-action">
          <button class="optimize-btn" :class="{ loading: optimizing }" :disabled="optimizing" @click="handleQuickOptimize">
            <span class="optimize-icon">🧠</span>
            <span class="optimize-text">{{ optimizing ? '优化中...' : '一键优化' }}</span>
          </button>
        </div>

        <div class="realtime-stats">
          <div class="stat-row">
            <span class="stat-label">CPU</span>
            <div class="stat-bar">
              <div class="stat-fill" :style="{ width: stats.cpu_percent + '%', background: pctColor(stats.cpu_percent) }"></div>
            </div>
            <span class="stat-val">{{ stats.cpu_percent.toFixed(0) }}%</span>
          </div>
          <div class="stat-row">
            <span class="stat-label">内存</span>
            <div class="stat-bar">
              <div class="stat-fill" :style="{ width: stats.mem_percent + '%', background: pctColor(stats.mem_percent) }"></div>
            </div>
            <span class="stat-val">{{ stats.mem_percent.toFixed(0) }}%</span>
          </div>
          <div class="net-row">
            <span class="net-up">↑ {{ formatSpeed(stats.net_up_speed) }}</span>
            <span class="net-down">↓ {{ formatSpeed(stats.net_down_speed) }}</span>
          </div>
        </div>

        <div class="footer">
          <span class="version" :class="{ update: hasUpdate }" :title="updateTip" @click="handleVersionClick">
            v{{ appVersion }}{{ hasUpdate ? ' ⬆' : '' }}
          </span>
          <a class="github" href="https://github.com/yezihack/WinCleaner" target="_blank">
            <svg viewBox="0 0 16 16" width="14" height="14" fill="currentColor">
              <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
            </svg>
          </a>
        </div>
      </div>
    </aside>

    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElNotification } from 'element-plus'
import { api, type RealtimeStats } from '@/api/backend'

const route = useRoute()
const optimizing = ref(false)
const appVersion = ref('0.2.0')
const hasUpdate = ref(false)
const releaseURL = ref('')
const updateTip = ref('')

const stats = reactive<RealtimeStats>({
  cpu_percent: 0,
  mem_percent: 0,
  net_up_speed: 0,
  net_down_speed: 0,
})

const navItems = [
  { path: '/', label: '系统概览', icon: '📊' },
  { path: '/cleaner', label: '垃圾清理', icon: '🗑️' },
  { path: '/memory', label: '内存优化', icon: '🧠' },
  { path: '/process', label: '进程管理', icon: '📈' },
  { path: '/network', label: '流量监控', icon: '🌐' },
  { path: '/disk', label: '磁盘管理', icon: '💾' },
  { path: '/port', label: '端口管理', icon: '🔌' },
]

let timer: ReturnType<typeof setInterval> | null = null

const fetchStats = async () => {
  try {
    const data = await api.getRealtimeStats()
    Object.assign(stats, data)
  } catch { /* silent */ }
}

const pctColor = (pct: number) => {
  if (pct < 50) return '#22c55e'
  if (pct < 80) return '#f59e0b'
  return '#ef4444'
}

const formatSpeed = (bps: number): string => {
  if (bps < 1024) return bps + ' B/s'
  if (bps < 1024 * 1024) return (bps / 1024).toFixed(1) + ' KB/s'
  return (bps / 1024 / 1024).toFixed(1) + ' MB/s'
}

const handleQuickOptimize = async () => {
  if (optimizing.value) return
  optimizing.value = true
  try {
    const result = await api.optimizeMemory()
    ElMessage.success(`释放了 ${result.freed_mb.toFixed(1)} MB 内存`)
  } catch {
    ElMessage.error('内存优化失败')
  } finally {
    optimizing.value = false
  }
}

const handleVersionClick = () => {
  if (hasUpdate.value && releaseURL.value) {
    window.open(releaseURL.value, '_blank')
  }
}

const checkAppVersion = async () => {
  try {
    appVersion.value = await api.getAppVersion()
  } catch { /* silent */ }
}

const checkUpdate = async () => {
  try {
    const info = await api.checkUpdate()
    appVersion.value = info.current_version
    if (info.has_update) {
      hasUpdate.value = true
      releaseURL.value = info.release_url
      updateTip.value = `发现新版本 v${info.latest_version}，点击前往下载`
      ElNotification({
        title: '发现新版本',
        message: `v${info.current_version} → v${info.latest_version}，点击查看更新`,
        type: 'warning',
        duration: 0,
        onClick: () => { window.open(info.release_url, '_blank') },
      })
    } else {
      updateTip.value = '已是最新版本'
    }
  } catch { /* silent */ }
}

onMounted(() => {
  fetchStats()
  timer = setInterval(fetchStats, 2000)
  checkAppVersion()
  checkUpdate()
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
html, body, #app { height: 100%; font-family: 'Microsoft YaHei', sans-serif; }
</style>

<style scoped>
.app-shell {
  display: flex;
  height: 100%;
  background: #f0f2f5;
}

.sidebar {
  width: 200px;
  background: linear-gradient(180deg, #0f172a 0%, #1e293b 100%);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.sidebar-top {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 18px 16px 14px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
}

.brand-icon {
  font-size: 24px;
  filter: drop-shadow(0 0 8px rgba(34,197,94,0.5));
}

.brand-text {
  display: flex;
  flex-direction: column;
}

.brand-name {
  font-size: 16px;
  font-weight: 700;
  color: #f8fafc;
  letter-spacing: 0.5px;
}

.brand-sub {
  font-size: 10px;
  color: #64748b;
  margin-top: 1px;
}

.nav-list {
  display: flex;
  flex-direction: column;
  padding: 12px 8px;
  gap: 2px;
  overflow-y: auto;
  flex: 1;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  color: #94a3b8;
  text-decoration: none;
  font-size: 13px;
  transition: all 0.15s ease;
}

.nav-item:hover {
  background: rgba(255,255,255,0.06);
  color: #e2e8f0;
}

.nav-item.active {
  background: rgba(34,197,94,0.15);
  color: #22c55e;
}

.nav-icon { font-size: 16px; flex-shrink: 0; }
.nav-label { font-weight: 500; }

.sidebar-bottom {
  border-top: 1px solid rgba(255,255,255,0.06);
  padding: 12px;
}

.quick-action {
  display: flex;
  justify-content: center;
  margin-bottom: 12px;
}

.optimize-btn {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px;
  background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%);
  border: none;
  border-radius: 10px;
  color: #fff;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 12px rgba(34,197,94,0.3);
}

.optimize-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 16px rgba(34,197,94,0.4);
}

.optimize-btn.loading {
  opacity: 0.8;
  cursor: not-allowed;
}

.optimize-icon { font-size: 16px; }
.optimize-text { font-size: 13px; }

.realtime-stats {
  background: rgba(0,0,0,0.2);
  border-radius: 8px;
  padding: 10px;
  margin-bottom: 10px;
}

.stat-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.stat-row:last-of-type { margin-bottom: 4px; }

.stat-label {
  font-size: 10px;
  color: #64748b;
  width: 26px;
  flex-shrink: 0;
}

.stat-bar {
  flex: 1;
  height: 4px;
  background: rgba(255,255,255,0.1);
  border-radius: 2px;
  overflow: hidden;
}

.stat-fill {
  height: 100%;
  border-radius: 2px;
  transition: width 0.5s ease, background 0.3s ease;
}

.stat-val {
  font-size: 10px;
  color: #94a3b8;
  width: 28px;
  text-align: right;
  flex-shrink: 0;
}

.net-row {
  display: flex;
  justify-content: space-between;
  padding-top: 2px;
}

.net-up, .net-down {
  font-size: 10px;
  font-weight: 500;
}

.net-up { color: #22c55e; }
.net-down { color: #3b82f6; }

.footer {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding-top: 8px;
  border-top: 1px solid rgba(255,255,255,0.04);
}

.version {
  font-size: 11px;
  color: #475569;
  cursor: default;
}

.version.update {
  color: #f59e0b;
  cursor: pointer;
  animation: blink 2s ease-in-out infinite;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.github {
  color: #475569;
  display: flex;
  align-items: center;
  transition: color 0.2s ease;
}

.github:hover { color: #fff; }

.main-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}
</style>
