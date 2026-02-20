<template>
  <el-container class="app-container">
    <el-aside width="200px" class="app-aside">
      <div class="logo">
        <el-icon :size="28"><Monitor /></el-icon>
        <span>Win Cleaner</span>
      </div>
      <el-menu
        :default-active="route.path"
        router
        background-color="#1d2b36"
        text-color="#a0aec0"
        active-text-color="#409eff"
      >
        <el-menu-item index="/">
          <el-icon><Odometer /></el-icon>
          <span>系统概览</span>
        </el-menu-item>
        <el-menu-item index="/cleaner">
          <el-icon><Delete /></el-icon>
          <span>垃圾清理</span>
        </el-menu-item>
        <el-menu-item index="/memory">
          <el-icon><Cpu /></el-icon>
          <span>内存优化</span>
        </el-menu-item>
        <el-menu-item index="/process">
          <el-icon><DataLine /></el-icon>
          <span>进程管理</span>
        </el-menu-item>
        <el-menu-item index="/network">
          <el-icon><Upload /></el-icon>
          <span>流量监控</span>
        </el-menu-item>
        <el-menu-item index="/disk">
          <el-icon><FolderOpened /></el-icon>
          <span>磁盘管理</span>
        </el-menu-item>
      </el-menu>

      <!-- 快捷优化按钮 -->
      <div class="quick-opt">
        <div
          class="opt-btn"
          :class="{ 'opt-loading': optimizing }"
          @click="handleQuickOptimize"
        >
          <span class="opt-icon">🧠</span>
          <span class="opt-text">{{ optimizing ? '优化中' : '一键优化' }}</span>
        </div>
      </div>

      <!-- 版本 & GitHub -->
      <div class="app-footer">
        <span class="app-version" :title="updateTip" :class="{ 'has-update': hasUpdate }" @click="handleVersionClick">
          v{{ appVersion }}{{ hasUpdate ? ' ⬆' : '' }}
        </span>
        <a class="github-link" href="https://github.com/yezihack/WinCleaner" target="_blank" title="GitHub">
          <svg viewBox="0 0 16 16" width="16" height="16" fill="currentColor">
            <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
          </svg>
        </a>
      </div>

      <!-- 左下角实时状态 -->
      <div class="realtime-bar">
        <div class="rt-item">
          <span class="rt-label">CPU</span>
          <el-progress
            :percentage="Math.round(stats.cpu_percent)"
            :stroke-width="6"
            :show-text="false"
            :color="barColor(stats.cpu_percent)"
          />
          <span class="rt-value">{{ stats.cpu_percent.toFixed(0) }}%</span>
        </div>
        <div class="rt-item">
          <span class="rt-label">内存</span>
          <el-progress
            :percentage="Math.round(stats.mem_percent)"
            :stroke-width="6"
            :show-text="false"
            :color="barColor(stats.mem_percent)"
          />
          <span class="rt-value">{{ stats.mem_percent.toFixed(0) }}%</span>
        </div>
        <div class="rt-net">
          <span class="rt-net-item">↑ {{ formatSpeed(stats.net_up_speed) }}</span>
          <span class="rt-net-item">↓ {{ formatSpeed(stats.net_down_speed) }}</span>
        </div>
      </div>
    </el-aside>

    <el-main class="app-main">
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElNotification } from 'element-plus'
import { api, type RealtimeStats } from '@/api/backend'

const route = useRoute()
const optimizing = ref(false)
const appVersion = ref('1.0.0')
const hasUpdate = ref(false)
const latestVersion = ref('')
const releaseURL = ref('')
const updateTip = ref('')

const stats = reactive<RealtimeStats>({
  cpu_percent: 0,
  mem_percent: 0,
  net_up_speed: 0,
  net_down_speed: 0,
})

let timer: ReturnType<typeof setInterval> | null = null

const fetchStats = async () => {
  try {
    const data = await api.getRealtimeStats()
    Object.assign(stats, data)
  } catch {
    // 静默失败
  }
}

const barColor = (pct: number) => {
  if (pct < 50) return '#67c23a'
  if (pct < 80) return '#e6a23c'
  return '#f56c6c'
}

const formatSpeed = (bytesPerSec: number): string => {
  if (bytesPerSec < 1024) return bytesPerSec + ' B/s'
  if (bytesPerSec < 1024 * 1024) return (bytesPerSec / 1024).toFixed(1) + ' KB/s'
  return (bytesPerSec / 1024 / 1024).toFixed(1) + ' MB/s'
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
      latestVersion.value = info.latest_version
      releaseURL.value = info.release_url
      updateTip.value = `发现新版本 v${info.latest_version}，点击前往下载`
      ElNotification({
        title: '发现新版本',
        message: `v${info.current_version} → v${info.latest_version}，点击查看更新`,
        type: 'warning',
        duration: 0,
        onClick: () => {
          window.open(info.release_url, '_blank')
        },
      })
    } else {
      updateTip.value = '已是最新版本'
    }
  } catch {
    // 网络不通时静默
  }
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
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  height: 100%;
  font-family: 'Microsoft YaHei', sans-serif;
}

.app-container {
  height: 100%;
}

.app-aside {
  background-color: #1d2b36;
  border-right: 1px solid #2d3d4a;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.app-aside .el-menu {
  flex: 1;
  border-right: none;
  overflow: hidden;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 16px;
  color: #409eff;
  font-size: 18px;
  font-weight: bold;
}

.app-main {
  background-color: #f0f2f5;
  padding: 20px;
  overflow-y: auto;
}

/* 版本 & GitHub */
.app-footer {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 6px 16px;
  border-top: 1px solid #2d3d4a;
}

.app-version {
  font-size: 11px;
  color: #6b7d8e;
}

.app-version.has-update {
  color: #e6a23c;
  cursor: pointer;
  animation: blink 2s ease-in-out infinite;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.github-link {
  color: #8899a6;
  transition: color 0.2s;
  display: flex;
  align-items: center;
}

.github-link:hover {
  color: #fff;
}

/* 实时状态栏 */
.realtime-bar {
  padding: 8px 16px;
  border-top: 1px solid #2d3d4a;
  background-color: #162029;
}

.rt-item {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 6px;
}

.rt-label {
  font-size: 11px;
  color: #8899a6;
  width: 28px;
  flex-shrink: 0;
}

.rt-item .el-progress {
  flex: 1;
}

.rt-value {
  font-size: 11px;
  color: #a0aec0;
  width: 32px;
  text-align: right;
  flex-shrink: 0;
}

.rt-net {
  display: flex;
  justify-content: space-between;
  margin-top: 4px;
}

.rt-net-item {
  font-size: 11px;
  color: #67c23a;
}

.rt-net-item:last-child {
  color: #409eff;
}

/* 快捷优化按钮 */
.quick-opt {
  display: flex;
  justify-content: center;
  padding: 8px 0;
}

.opt-btn {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: linear-gradient(135deg, #409eff 0%, #67c23a 100%);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 12px rgba(64, 158, 255, 0.4);
}

.opt-btn:hover {
  transform: scale(1.08);
  box-shadow: 0 4px 16px rgba(64, 158, 255, 0.6);
}

.opt-btn:active {
  transform: scale(0.95);
}

.opt-btn.opt-loading {
  animation: pulse 1.2s ease-in-out infinite;
  pointer-events: none;
  opacity: 0.8;
}

.opt-icon {
  font-size: 22px;
  line-height: 1;
}

.opt-text {
  font-size: 10px;
  color: #fff;
  margin-top: 2px;
  font-weight: 500;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); box-shadow: 0 2px 12px rgba(64, 158, 255, 0.4); }
  50% { transform: scale(1.06); box-shadow: 0 4px 20px rgba(64, 158, 255, 0.7); }
}
</style>
