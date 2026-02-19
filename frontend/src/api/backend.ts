// Wails 绑定调用封装
// 实际运行时由 Wails 自动生成绑定，这里提供类型定义和开发时的 mock

export interface SystemInfo {
  os: string
  hostname: string
  cpu_usage: number
  mem_total: number
  mem_used: number
  mem_percent: number
  disk_total: number
  disk_used: number
  disk_percent: number
  public_ip: string
  ip_location: string
  ip_operator: string
}

export interface ScanResult {
  category: string
  items: JunkItem[]
  size: number
  count: number
}

export interface JunkItem {
  path: string
  size: number
  category: string
}

export interface CleanResult {
  freed_size: number
  cleaned_count: number
  failed_count: number
}

export interface MemoryOptResult {
  before_used: number
  after_used: number
  freed_mb: number
  before_percent: number
  after_percent: number
}

export interface ProcessInfo {
  pid: number
  name: string
  cpu_percent: number
  mem_rss: number
  mem_percent: number
  status: string
  username: string
}

export interface GPUInfo {
  name: string
  type: 'discrete' | 'integrated' | 'none'
  type_label: string
  vram: number
  driver_ver: string
  resolution: string
}

export interface GPUResult {
  gpus: GPUInfo[]
}

export interface DailyStat {
  date: string
  freed_size: number
  count: number
}

export interface MonthlyStat {
  month: string
  freed_size: number
  count: number
}

export interface CleanHistoryStats {
  records: { date: string; time: string; freed_size: number; cleaned_count: number }[]
  daily_stats: DailyStat[]
  monthly_stats: MonthlyStat[]
  last_clean_time: string
  last_clean_ago: string
  total_freed: number
  total_count: number
}

export interface RealtimeStats {
  cpu_percent: number
  mem_percent: number
  net_up_speed: number
  net_down_speed: number
}

export interface NetTrafficInfo {
  total_sent: number
  total_recv: number
  up_speed: number
  down_speed: number
}

export interface ProcessNetInfo {
  name: string
  count: number
  sent: number
  recv: number
}

export interface NetTrafficResult {
  overview: NetTrafficInfo
  processes: ProcessNetInfo[]
}

export interface NetDailyStat {
  date: string
  sent: number
  recv: number
}

export interface NetMonthlyStat {
  month: string
  sent: number
  recv: number
}

export interface NetYearlyStat {
  year: string
  sent: number
  recv: number
}

export interface NetTrafficStats {
  daily_stats: NetDailyStat[]
  monthly_stats: NetMonthlyStat[]
  yearly_stats: NetYearlyStat[]
  total_sent: number
  total_recv: number
}

export interface DiskInfo {
  device: string
  mountpoint: string
  fstype: string
  total: number
  used: number
  free: number
  used_percent: number
}

export interface LargeFileInfo {
  path: string
  size: number
  ext: string
}

export interface DiskScanResult {
  files: LargeFileInfo[]
  count: number
}

export interface MemOptRecord {
  date: string
  time: string
  freed_mb: number
  before_percent: number
  after_percent: number
}

export interface MemOptDailyStat {
  date: string
  freed_mb: number
  count: number
}

export interface MemOptMonthlyStat {
  month: string
  freed_mb: number
  count: number
}

export interface MemOptStats {
  recent_records: MemOptRecord[] | null
  daily_stats: MemOptDailyStat[] | null
  monthly_stats: MemOptMonthlyStat[] | null
  last_opt_time: string
  last_opt_ago: string
  total_freed_mb: number
  total_count: number
}

declare global {
  interface Window {
    go: {
      app: {
        App: {
          GetSystemInfo(): Promise<SystemInfo>
          ScanJunk(): Promise<ScanResult[]>
          CleanJunk(categories: string[]): Promise<CleanResult>
          OptimizeMemory(): Promise<MemoryOptResult>
          GetProcessList(): Promise<ProcessInfo[]>
          KillProcess(pid: number): Promise<void>
          GetGPUInfo(): Promise<GPUResult>
          GetCleanHistory(): Promise<CleanHistoryStats>
          GetRealtimeStats(): Promise<RealtimeStats>
          GetNetTraffic(): Promise<NetTrafficResult>
          GetNetTrafficStats(): Promise<NetTrafficStats>
          GetDiskList(): Promise<DiskInfo[]>
          ScanLargeFiles(root: string, minSizeMB: number): Promise<DiskScanResult>
          GetMemOptStats(): Promise<MemOptStats>
        }
      }
    }
  }
}

export const api = {
  getSystemInfo: (): Promise<SystemInfo> =>
    window.go.app.App.GetSystemInfo(),

  scanJunk: (): Promise<ScanResult[]> =>
    window.go.app.App.ScanJunk(),

  cleanJunk: (categories: string[]): Promise<CleanResult> =>
    window.go.app.App.CleanJunk(categories),

  optimizeMemory: (): Promise<MemoryOptResult> =>
    window.go.app.App.OptimizeMemory(),

  getProcessList: (): Promise<ProcessInfo[]> =>
    window.go.app.App.GetProcessList(),

  killProcess: (pid: number): Promise<void> =>
    window.go.app.App.KillProcess(pid),

  getGPUInfo: (): Promise<GPUResult> =>
    window.go.app.App.GetGPUInfo(),

  getCleanHistory: (): Promise<CleanHistoryStats> =>
    window.go.app.App.GetCleanHistory(),

  getRealtimeStats: (): Promise<RealtimeStats> =>
    window.go.app.App.GetRealtimeStats(),

  getNetTraffic: (): Promise<NetTrafficResult> =>
    window.go.app.App.GetNetTraffic(),

  getNetTrafficStats: (): Promise<NetTrafficStats> =>
    window.go.app.App.GetNetTrafficStats(),

  getDiskList: (): Promise<DiskInfo[]> =>
    window.go.app.App.GetDiskList(),

  scanLargeFiles: (root: string, minSizeMB: number): Promise<DiskScanResult> =>
    window.go.app.App.ScanLargeFiles(root, minSizeMB),

  getMemOptStats: (): Promise<MemOptStats> =>
    window.go.app.App.GetMemOptStats(),
}
