package monitor

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"win-cleaner/internal/model"
	"win-cleaner/pkg/datadir"

	"github.com/shirou/gopsutil/v3/net"
)

var (
	nhMu           sync.Mutex
	prevSampleSent uint64
	prevSampleRecv uint64
	prevSampleTime time.Time
	sampleInited   bool
)

func getNetHistoryPath() string {
	return datadir.FilePath("net_history.json")
}

func loadNetHistory() (*model.NetTrafficHistory, error) {
	path := getNetHistoryPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &model.NetTrafficHistory{}, nil
		}
		return nil, err
	}
	var h model.NetTrafficHistory
	if err := json.Unmarshal(data, &h); err != nil {
		return &model.NetTrafficHistory{}, nil
	}
	return &h, nil
}

func saveNetHistory(h *model.NetTrafficHistory) error {
	path := getNetHistoryPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.Marshal(h)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// RecordNetTrafficSample 采样并记录当前流量增量
func RecordNetTrafficSample() {
	nhMu.Lock()
	defer nhMu.Unlock()

	counters, err := net.IOCounters(false)
	if err != nil || len(counters) == 0 {
		return
	}

	currentSent := counters[0].BytesSent
	currentRecv := counters[0].BytesRecv
	now := time.Now()

	if !sampleInited {
		prevSampleSent = currentSent
		prevSampleRecv = currentRecv
		prevSampleTime = now
		sampleInited = true
		return
	}

	deltaSent := currentSent - prevSampleSent
	deltaRecv := currentRecv - prevSampleRecv

	prevSampleSent = currentSent
	prevSampleRecv = currentRecv
	prevSampleTime = now

	if deltaSent == 0 && deltaRecv == 0 {
		return
	}

	h, err := loadNetHistory()
	if err != nil {
		return
	}

	record := model.NetTrafficRecord{
		Timestamp: now.Format("2006-01-02 15:04"),
		Date:      now.Format("2006-01-02"),
		Sent:      deltaSent,
		Recv:      deltaRecv,
	}

	// 合并同一分钟的记录
	if len(h.Records) > 0 {
		last := &h.Records[len(h.Records)-1]
		if last.Timestamp == record.Timestamp {
			last.Sent += record.Sent
			last.Recv += record.Recv
			_ = saveNetHistory(h)
			return
		}
	}

	h.Records = append(h.Records, record)

	// 只保留最近 90 天的记录
	cutoff := now.AddDate(0, 0, -90).Format("2006-01-02")
	filtered := h.Records[:0]
	for _, r := range h.Records {
		if r.Date >= cutoff {
			filtered = append(filtered, r)
		}
	}
	h.Records = filtered

	_ = saveNetHistory(h)
}

// GetNetTrafficStats 获取流量历史统计
func GetNetTrafficStats() (*model.NetTrafficStats, error) {
	nhMu.Lock()
	defer nhMu.Unlock()

	h, err := loadNetHistory()
	if err != nil {
		return nil, err
	}

	stats := &model.NetTrafficStats{}

	// 按天
	dailyMap := make(map[string]*model.NetDailyStat)
	for _, r := range h.Records {
		if d, ok := dailyMap[r.Date]; ok {
			d.Sent += r.Sent
			d.Recv += r.Recv
		} else {
			dailyMap[r.Date] = &model.NetDailyStat{
				Date: r.Date, Sent: r.Sent, Recv: r.Recv,
			}
		}
		stats.TotalSent += r.Sent
		stats.TotalRecv += r.Recv
	}

	// 近 30 天
	cutoff30 := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	for _, d := range dailyMap {
		if d.Date >= cutoff30 {
			stats.DailyStats = append(stats.DailyStats, *d)
		}
	}
	sort.Slice(stats.DailyStats, func(i, j int) bool {
		return stats.DailyStats[i].Date < stats.DailyStats[j].Date
	})

	// 按月
	monthlyMap := make(map[string]*model.NetMonthlyStat)
	for _, r := range h.Records {
		month := r.Date[:7]
		if m, ok := monthlyMap[month]; ok {
			m.Sent += r.Sent
			m.Recv += r.Recv
		} else {
			monthlyMap[month] = &model.NetMonthlyStat{
				Month: month, Sent: r.Sent, Recv: r.Recv,
			}
		}
	}
	for _, m := range monthlyMap {
		stats.MonthlyStats = append(stats.MonthlyStats, *m)
	}
	sort.Slice(stats.MonthlyStats, func(i, j int) bool {
		return stats.MonthlyStats[i].Month < stats.MonthlyStats[j].Month
	})

	// 按年
	yearlyMap := make(map[string]*model.NetYearlyStat)
	for _, r := range h.Records {
		year := r.Date[:4]
		if y, ok := yearlyMap[year]; ok {
			y.Sent += r.Sent
			y.Recv += r.Recv
		} else {
			yearlyMap[year] = &model.NetYearlyStat{
				Year: year, Sent: r.Sent, Recv: r.Recv,
			}
		}
	}
	for _, y := range yearlyMap {
		stats.YearlyStats = append(stats.YearlyStats, *y)
	}
	sort.Slice(stats.YearlyStats, func(i, j int) bool {
		return stats.YearlyStats[i].Year < stats.YearlyStats[j].Year
	})

	return stats, nil
}
