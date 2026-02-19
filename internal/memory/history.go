package memory

import (
	"encoding/json"
	"os"
	"sort"
	"strconv"
	"time"

	"win-cleaner/internal/model"
	"win-cleaner/pkg/datadir"
)

func memHistoryPath() string {
	return datadir.FilePath("mem_opt_history.json")
}

// RecordOptimize 记录一次优化
func RecordOptimize(result *model.MemoryOptResult) error {
	filePath := memHistoryPath()

	var history model.MemOptHistory
	data, err := os.ReadFile(filePath)
	if err == nil {
		_ = json.Unmarshal(data, &history)
	}

	now := time.Now()
	record := model.MemOptRecord{
		Date:          now.Format("2006-01-02"),
		Time:          now.Format("15:04:05"),
		FreedMB:       result.FreedMB,
		BeforePercent: result.BeforePercent,
		AfterPercent:  result.AfterPercent,
	}
	history.Records = append(history.Records, record)

	// 只保留最近 90 天
	cutoff := now.AddDate(0, 0, -90).Format("2006-01-02")
	var filtered []model.MemOptRecord
	for _, r := range history.Records {
		if r.Date >= cutoff {
			filtered = append(filtered, r)
		}
	}
	history.Records = filtered

	out, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(memHistoryPath(), out, 0644)
}

// GetMemOptStats 获取优化历史统计
func GetMemOptStats() (*model.MemOptStats, error) {
	var history model.MemOptHistory
	data, err := os.ReadFile(memHistoryPath())
	if err != nil {
		// 没有历史记录，返回空
		return &model.MemOptStats{}, nil
	}
	if err := json.Unmarshal(data, &history); err != nil {
		return &model.MemOptStats{}, nil
	}

	stats := &model.MemOptStats{}
	if len(history.Records) == 0 {
		return stats, nil
	}

	// 总计
	for _, r := range history.Records {
		stats.TotalFreedMB += r.FreedMB
		stats.TotalCount++
	}

	// 上次优化
	last := history.Records[len(history.Records)-1]
	stats.LastOptTime = last.Date + " " + last.Time
	lastTime, err := time.Parse("2006-01-02 15:04:05", stats.LastOptTime)
	if err == nil {
		dur := time.Since(lastTime)
		if dur.Hours() < 1 {
			stats.LastOptAgo = "刚刚"
		} else if dur.Hours() < 24 {
			hours := int(dur.Hours())
			stats.LastOptAgo = strconv.Itoa(hours) + " 小时前"
		} else {
			days := int(dur.Hours() / 24)
			stats.LastOptAgo = strconv.Itoa(days) + " 天前"
		}
	}

	// 按天汇总（近 30 天）
	dailyMap := make(map[string]*model.MemOptDailyStat)
	cutoff30 := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	for _, r := range history.Records {
		if r.Date < cutoff30 {
			continue
		}
		if d, ok := dailyMap[r.Date]; ok {
			d.FreedMB += r.FreedMB
			d.Count++
		} else {
			dailyMap[r.Date] = &model.MemOptDailyStat{
				Date:    r.Date,
				FreedMB: r.FreedMB,
				Count:   1,
			}
		}
	}
	for _, d := range dailyMap {
		stats.DailyStats = append(stats.DailyStats, *d)
	}
	sort.Slice(stats.DailyStats, func(i, j int) bool {
		return stats.DailyStats[i].Date < stats.DailyStats[j].Date
	})

	// 按月汇总
	monthlyMap := make(map[string]*model.MemOptMonthlyStat)
	for _, r := range history.Records {
		month := r.Date[:7]
		if m, ok := monthlyMap[month]; ok {
			m.FreedMB += r.FreedMB
			m.Count++
		} else {
			monthlyMap[month] = &model.MemOptMonthlyStat{
				Month:   month,
				FreedMB: r.FreedMB,
				Count:   1,
			}
		}
	}
	for _, m := range monthlyMap {
		stats.MonthlyStats = append(stats.MonthlyStats, *m)
	}
	sort.Slice(stats.MonthlyStats, func(i, j int) bool {
		return stats.MonthlyStats[i].Month < stats.MonthlyStats[j].Month
	})

	// 最近 10 次记录（用于折线图）
	start := 0
	if len(history.Records) > 10 {
		start = len(history.Records) - 10
	}
	stats.RecentRecords = history.Records[start:]

	return stats, nil
}
