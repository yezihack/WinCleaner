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
import { ElMessage } from 'element-plus'
import { api, type RealtimeStats } from '@/api/backend'

const route = useRoute()
const optimizing = ref(false)

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

onMounted(() => {
  fetchStats()
  timer = setInterval(fetchStats, 2000)
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
}

.app-aside .el-menu {
  flex: 1;
  border-right: none;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 20px 16px;
  color: #409eff;
  font-size: 18px;
  font-weight: bold;
}

.app-main {
  background-color: #f0f2f5;
  padding: 20px;
  overflow-y: auto;
}

/* 实时状态栏 */
.realtime-bar {
  padding: 12px 16px;
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
  padding: 12px 0;
}

.opt-btn {
  width: 72px;
  height: 72px;
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
