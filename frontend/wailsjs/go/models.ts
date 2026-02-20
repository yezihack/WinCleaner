export namespace model {
	
	export class MonthlyStat {
	    month: string;
	    freed_size: number;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new MonthlyStat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.month = source["month"];
	        this.freed_size = source["freed_size"];
	        this.count = source["count"];
	    }
	}
	export class DailyStat {
	    date: string;
	    freed_size: number;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new DailyStat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.freed_size = source["freed_size"];
	        this.count = source["count"];
	    }
	}
	export class CleanRecord {
	    date: string;
	    time: string;
	    freed_size: number;
	    cleaned_count: number;
	
	    static createFrom(source: any = {}) {
	        return new CleanRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.time = source["time"];
	        this.freed_size = source["freed_size"];
	        this.cleaned_count = source["cleaned_count"];
	    }
	}
	export class CleanHistoryStats {
	    records: CleanRecord[];
	    daily_stats: DailyStat[];
	    monthly_stats: MonthlyStat[];
	    last_clean_time: string;
	    last_clean_ago: string;
	    total_freed: number;
	    total_count: number;
	
	    static createFrom(source: any = {}) {
	        return new CleanHistoryStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.records = this.convertValues(source["records"], CleanRecord);
	        this.daily_stats = this.convertValues(source["daily_stats"], DailyStat);
	        this.monthly_stats = this.convertValues(source["monthly_stats"], MonthlyStat);
	        this.last_clean_time = source["last_clean_time"];
	        this.last_clean_ago = source["last_clean_ago"];
	        this.total_freed = source["total_freed"];
	        this.total_count = source["total_count"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class CleanResult {
	    freed_size: number;
	    cleaned_count: number;
	    failed_count: number;
	
	    static createFrom(source: any = {}) {
	        return new CleanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.freed_size = source["freed_size"];
	        this.cleaned_count = source["cleaned_count"];
	        this.failed_count = source["failed_count"];
	    }
	}
	
	export class DiskInfo {
	    device: string;
	    mountpoint: string;
	    fstype: string;
	    total: number;
	    used: number;
	    free: number;
	    used_percent: number;
	
	    static createFrom(source: any = {}) {
	        return new DiskInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.device = source["device"];
	        this.mountpoint = source["mountpoint"];
	        this.fstype = source["fstype"];
	        this.total = source["total"];
	        this.used = source["used"];
	        this.free = source["free"];
	        this.used_percent = source["used_percent"];
	    }
	}
	export class LargeFileInfo {
	    path: string;
	    size: number;
	    ext: string;
	
	    static createFrom(source: any = {}) {
	        return new LargeFileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.size = source["size"];
	        this.ext = source["ext"];
	    }
	}
	export class DiskScanResult {
	    files: LargeFileInfo[];
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new DiskScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.files = this.convertValues(source["files"], LargeFileInfo);
	        this.count = source["count"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GPUInfo {
	    name: string;
	    type: string;
	    type_label: string;
	    vram: number;
	    driver_ver: string;
	    resolution: string;
	
	    static createFrom(source: any = {}) {
	        return new GPUInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.type_label = source["type_label"];
	        this.vram = source["vram"];
	        this.driver_ver = source["driver_ver"];
	        this.resolution = source["resolution"];
	    }
	}
	export class GPUResult {
	    gpus: GPUInfo[];
	
	    static createFrom(source: any = {}) {
	        return new GPUResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.gpus = this.convertValues(source["gpus"], GPUInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class JunkItem {
	    path: string;
	    size: number;
	    category: string;
	
	    static createFrom(source: any = {}) {
	        return new JunkItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.size = source["size"];
	        this.category = source["category"];
	    }
	}
	
	export class MemOptDailyStat {
	    date: string;
	    freed_mb: number;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new MemOptDailyStat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.freed_mb = source["freed_mb"];
	        this.count = source["count"];
	    }
	}
	export class MemOptMonthlyStat {
	    month: string;
	    freed_mb: number;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new MemOptMonthlyStat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.month = source["month"];
	        this.freed_mb = source["freed_mb"];
	        this.count = source["count"];
	    }
	}
	export class MemOptRecord {
	    date: string;
	    time: string;
	    freed_mb: number;
	    before_percent: number;
	    after_percent: number;
	
	    static createFrom(source: any = {}) {
	        return new MemOptRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.time = source["time"];
	        this.freed_mb = source["freed_mb"];
	        this.before_percent = source["before_percent"];
	        this.after_percent = source["after_percent"];
	    }
	}
	export class MemOptStats {
	    recent_records: MemOptRecord[];
	    daily_stats: MemOptDailyStat[];
	    monthly_stats: MemOptMonthlyStat[];
	    last_opt_time: string;
	    last_opt_ago: string;
	    total_freed_mb: number;
	    total_count: number;
	
	    static createFrom(source: any = {}) {
	        return new MemOptStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.recent_records = this.convertValues(source["recent_records"], MemOptRecord);
	        this.daily_stats = this.convertValues(source["daily_stats"], MemOptDailyStat);
	        this.monthly_stats = this.convertValues(source["monthly_stats"], MemOptMonthlyStat);
	        this.last_opt_time = source["last_opt_time"];
	        this.last_opt_ago = source["last_opt_ago"];
	        this.total_freed_mb = source["total_freed_mb"];
	        this.total_count = source["total_count"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MemoryOptResult {
	    before_used: number;
	    after_used: number;
	    freed_mb: number;
	    before_percent: number;
	    after_percent: number;
	
	    static createFrom(source: any = {}) {
	        return new MemoryOptResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.before_used = source["before_used"];
	        this.after_used = source["after_used"];
	        this.freed_mb = source["freed_mb"];
	        this.before_percent = source["before_percent"];
	        this.after_percent = source["after_percent"];
	    }
	}
	
	export class NetDailyStat {
	    date: string;
	    sent: number;
	    recv: number;
	
	    static createFrom(source: any = {}) {
	        return new NetDailyStat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.sent = source["sent"];
	        this.recv = source["recv"];
	    }
	}
	export class NetMonthlyStat {
	    month: string;
	    sent: number;
	    recv: number;
	
	    static createFrom(source: any = {}) {
	        return new NetMonthlyStat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.month = source["month"];
	        this.sent = source["sent"];
	        this.recv = source["recv"];
	    }
	}
	export class NetTrafficInfo {
	    total_sent: number;
	    total_recv: number;
	    up_speed: number;
	    down_speed: number;
	
	    static createFrom(source: any = {}) {
	        return new NetTrafficInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_sent = source["total_sent"];
	        this.total_recv = source["total_recv"];
	        this.up_speed = source["up_speed"];
	        this.down_speed = source["down_speed"];
	    }
	}
	export class ProcessNetInfo {
	    name: string;
	    count: number;
	    sent: number;
	    recv: number;
	
	    static createFrom(source: any = {}) {
	        return new ProcessNetInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.count = source["count"];
	        this.sent = source["sent"];
	        this.recv = source["recv"];
	    }
	}
	export class NetTrafficResult {
	    overview: NetTrafficInfo;
	    processes: ProcessNetInfo[];
	
	    static createFrom(source: any = {}) {
	        return new NetTrafficResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.overview = this.convertValues(source["overview"], NetTrafficInfo);
	        this.processes = this.convertValues(source["processes"], ProcessNetInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class NetYearlyStat {
	    year: string;
	    sent: number;
	    recv: number;
	
	    static createFrom(source: any = {}) {
	        return new NetYearlyStat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.year = source["year"];
	        this.sent = source["sent"];
	        this.recv = source["recv"];
	    }
	}
	export class NetTrafficStats {
	    daily_stats: NetDailyStat[];
	    monthly_stats: NetMonthlyStat[];
	    yearly_stats: NetYearlyStat[];
	    total_sent: number;
	    total_recv: number;
	
	    static createFrom(source: any = {}) {
	        return new NetTrafficStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.daily_stats = this.convertValues(source["daily_stats"], NetDailyStat);
	        this.monthly_stats = this.convertValues(source["monthly_stats"], NetMonthlyStat);
	        this.yearly_stats = this.convertValues(source["yearly_stats"], NetYearlyStat);
	        this.total_sent = source["total_sent"];
	        this.total_recv = source["total_recv"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class ProcessInfo {
	    pid: number;
	    name: string;
	    cpu_percent: number;
	    mem_rss: number;
	    mem_percent: number;
	    status: string;
	    username: string;
	
	    static createFrom(source: any = {}) {
	        return new ProcessInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pid = source["pid"];
	        this.name = source["name"];
	        this.cpu_percent = source["cpu_percent"];
	        this.mem_rss = source["mem_rss"];
	        this.mem_percent = source["mem_percent"];
	        this.status = source["status"];
	        this.username = source["username"];
	    }
	}
	
	export class RealtimeStats {
	    cpu_percent: number;
	    mem_percent: number;
	    net_up_speed: number;
	    net_down_speed: number;
	
	    static createFrom(source: any = {}) {
	        return new RealtimeStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cpu_percent = source["cpu_percent"];
	        this.mem_percent = source["mem_percent"];
	        this.net_up_speed = source["net_up_speed"];
	        this.net_down_speed = source["net_down_speed"];
	    }
	}
	export class ScanResult {
	    category: string;
	    items: JunkItem[];
	    size: number;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new ScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.category = source["category"];
	        this.items = this.convertValues(source["items"], JunkItem);
	        this.size = source["size"];
	        this.count = source["count"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SystemInfo {
	    os: string;
	    hostname: string;
	    cpu_usage: number;
	    mem_total: number;
	    mem_used: number;
	    mem_percent: number;
	    disk_total: number;
	    disk_used: number;
	    disk_percent: number;
	    public_ip: string;
	    ip_location: string;
	    ip_operator: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.os = source["os"];
	        this.hostname = source["hostname"];
	        this.cpu_usage = source["cpu_usage"];
	        this.mem_total = source["mem_total"];
	        this.mem_used = source["mem_used"];
	        this.mem_percent = source["mem_percent"];
	        this.disk_total = source["disk_total"];
	        this.disk_used = source["disk_used"];
	        this.disk_percent = source["disk_percent"];
	        this.public_ip = source["public_ip"];
	        this.ip_location = source["ip_location"];
	        this.ip_operator = source["ip_operator"];
	    }
	}
	export class UpdateInfo {
	    current_version: string;
	    latest_version: string;
	    has_update: boolean;
	    release_url: string;
	    release_notes: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.current_version = source["current_version"];
	        this.latest_version = source["latest_version"];
	        this.has_update = source["has_update"];
	        this.release_url = source["release_url"];
	        this.release_notes = source["release_notes"];
	    }
	}

}

