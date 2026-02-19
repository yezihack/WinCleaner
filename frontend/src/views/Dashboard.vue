<template>
  <div class="dashboard">
    <h2 style="margin-bottom: 20px; color: #303133;">系统概览</h2>

    <el-row :gutter="20">
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header><span>CPU 使用率</span></template>
          <el-progress
            type="dashboard"
            :percentage="Math.round(info.cpu_usage)"
            :color="progressColor"
          />
          <p class="stat-text">&nbsp;</p>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card shadow="hover">
          <template #header><span>内存使用</span></template>
          <el-progress
            type="dashboard"
            :percentage="Math.round(info.mem_percent)"
            :color="progressColor"
          />
          <p class="stat-text">
            {{ formatBytes(info.mem_used) }} / {{ formatBytes(info.mem_total) }}
          </p>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card shadow="hover">
          <template #header><span>磁盘使用 (C:)</span></template>
          <el-progress
            type="dashboard"
            :percentage="Math.round(info.disk_percent)"
            :color="progressColor"
          />
          <p class="stat-text">
            {{ formatBytes(info.disk_used) }} / {{ formatBytes(info.disk_total) }}
          </p>
        </el-card>
      </el-col>
    </el-row>

    <el-row style="margin-top: 20px;">
      <el-col :span="24">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header-row">
              <span>显卡信息</span>
              <el-tag v-if="gpuLoaded && gpus.length === 0" type="info" size="small">未检测到显卡</el-tag>
            </div>
          </template>

          <div v-if="gpuLoading" v-loading="true" style="height: 80px;" />

          <div v-else-if="gpus.length > 0" class="gpu-list">
            <div v-for="(gpu, idx) in gpus" :key="idx" class="gpu-item">
              <div class="gpu-header">
                <el-tag
                  :type="gpu.type === 'discrete' ? 'danger' : gpu.type === 'integrated' ? 'warning' : 'info'"
                  size="small"
                  effect="dark"
                >
                  {{ gpu.type_label }}
                </el-tag>
                <span class="gpu-name">{{ gpu.name }}</span>
              </div>
              <el-descriptions :column="3" size="small" border style="margin-top: 8px;">
                <el-descriptions-item label="显存">
                  {{ gpu.vram > 0 ? formatBytes(gpu.vram) : '共享内存' }}
                </el-descriptions-item>
                <el-descriptions-item label="驱动版本">
                  {{ gpu.driver_ver || '-' }}
                </el-descriptions-item>
                <el-descriptions-item label="分辨率">
                  {{ gpu.resolution || '-' }}
                </el-descriptions-item>
              </el-descriptions>
              <div v-if="gpu.vram > 0" style="margin-top: 8px;">
                <span class="stat-text">显存容量</span>
                <el-progress
                  :percentage="100"
                  :stroke-width="16"
                  :format="() => formatBytes(gpu.vram)"
                  :color="gpu.type === 'discrete' ? '#f56c6c' : '#e6a23c'"
                />
              </div>
            </div>
          </div>

          <el-empty v-else description="未检测到显卡设备" :image-size="60" />
        </el-card>
      </el-col>
    </el-row>

    <el-row style="margin-top: 20px;">
      <el-col :span="24">
        <el-card shadow="hover">
          <template #header><span>系统信息</span></template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="主机名">{{ info.hostname }}</el-descriptions-item>
            <el-descriptions-item label="操作系统">{{ info.os }}</el-descriptions-item>
            <el-descriptions-item label="公网 IP">{{ info.public_ip || '获取中...' }}</el-descriptions-item>
            <el-descriptions-item label="归属地">{{ info.ip_location || '-' }}</el-descriptions-item>
            <el-descriptions-item label="运营商">{{ info.ip_operator || '-' }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>

    <div style="margin-top: 16px; text-align: right;">
      <el-button type="primary" :loading="loading" @click="refresh">
        <el-icon><Refresh /></el-icon> 刷新
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { api, type SystemInfo, type GPUInfo } from '@/api/backend'

const loading = ref(false)
const gpuLoading = ref(false)
const gpuLoaded = ref(false)
const gpus = ref<GPUInfo[]>([])
const info = reactive<SystemInfo>({
  os: '',
  hostname: '',
  cpu_usage: 0,
  mem_total: 0,
  mem_used: 0,
  mem_percent: 0,
  disk_total: 0,
  disk_used: 0,
  disk_percent: 0,
  public_ip: '',
  ip_location: '',
  ip_operator: '',
})

const progressColor = (percentage: number) => {
  if (percentage < 50) return '#67c23a'
  if (percentage < 80) return '#e6a23c'
  return '#f56c6c'
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
  } catch (e) {
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
  } catch (e) {
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
.dashboard .el-card {
  text-align: center;
}
.stat-text {
  margin-top: 8px;
  color: #909399;
  font-size: 13px;
}
.card-header-row {
  display: flex;
  align-items: center;
  gap: 8px;
}
.gpu-list {
  text-align: left;
}
.gpu-item {
  padding: 12px 0;
}
.gpu-item + .gpu-item {
  border-top: 1px solid #ebeef5;
}
.gpu-header {
  display: flex;
  align-items: center;
  gap: 8px;
}
.gpu-name {
  font-size: 15px;
  font-weight: 500;
  color: #303133;
}
</style>
