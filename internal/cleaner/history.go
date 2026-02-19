package cleaner

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"win-cleaner/internal/model"
	"win-cleaner/pkg/datadir"
)

// getHistoryPath 获取历史文件路径 (~/.wincleaner/clean_history.json)
func getHistoryPath() string {
	return datadir.FilePath("clean_history.json")
}

// loadHistory 加载历史记录
func loadHistory() (*model.CleanHistory, error) {
	path := getHistoryPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &model.CleanHistory{}, nil
		}
		return nil, err
	}

	var history model.CleanHistory
	if err := json.Unmarshal(data, &history); err != nil {
		return &model.CleanHistory{}, nil
	}
	return &history, nil
}

// saveHistory 保存历史记录
func saveHistory(history *model.CleanHistory) error {
	path := getHistoryPath()

	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// RecordClean 记录一次清理
func RecordClean(result model.CleanResult) error {
	history, err := loadHistory()
	if err != nil {
		return err
	}

	now := time.Now()
	record := model.CleanRecord{
		Date:         now.Format("2006-01-02"),
		Time:         now.Format("15:04:05"),
		FreedSize:    result.FreedSize,
		CleanedCount: result.CleanedCount,
	}

	history.Records = append(history.Records, record)
	return saveHistory(history)
}

// GetCleanHistoryStats 获取清理历史统计
func GetCleanHistoryStats() (*model.CleanHistoryStats, error) {
	history, err := loadHistory()
	if err != nil {
		return nil, err
	}

	stats := &model.CleanHistoryStats{
		Records: history.Records,
	}

	if len(history.Records) == 0 {
		stats.LastCleanTime = ""
		stats.LastCleanAgo = "从未清理"
		return stats, nil
	}

	// 累计
	for _, r := range history.Records {
		stats.TotalFreed += r.FreedSize
		stats.TotalCount += r.CleanedCount
	}

	// 上次清理
	last := history.Records[len(history.Records)-1]
	stats.LastCleanTime = last.Date + " " + last.Time
	lastTime, err := time.Parse("2006-01-02 15:04:05", stats.LastCleanTime)
	if err == nil {
		stats.LastCleanAgo = formatDuration(time.Since(lastTime))
	}

	// 按天汇总（近30天）
	dailyMap := make(map[string]*model.DailyStat)
	for _, r := range history.Records {
		if d, ok := dailyMap[r.Date]; ok {
			d.FreedSize += r.FreedSize
			d.Count += r.CleanedCount
		} else {
			dailyMap[r.Date] = &model.DailyStat{
				Date:      r.Date,
				FreedSize: r.FreedSize,
				Count:     r.CleanedCount,
			}
		}
	}

	// 取近30天
	cutoff := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	for _, d := range dailyMap {
		if d.Date >= cutoff {
			stats.DailyStats = append(stats.DailyStats, *d)
		}
	}
	sort.Slice(stats.DailyStats, func(i, j int) bool {
		return stats.DailyStats[i].Date < stats.DailyStats[j].Date
	})

	// 按月汇总
	monthlyMap := make(map[string]*model.MonthlyStat)
	for _, r := range history.Records {
		month := r.Date[:7] // YYYY-MM
		if m, ok := monthlyMap[month]; ok {
			m.FreedSize += r.FreedSize
			m.Count += r.CleanedCount
		} else {
			monthlyMap[month] = &model.MonthlyStat{
				Month:     month,
				FreedSize: r.FreedSize,
				Count:     r.CleanedCount,
			}
		}
	}
	for _, m := range monthlyMap {
		stats.MonthlyStats = append(stats.MonthlyStats, *m)
	}
	sort.Slice(stats.MonthlyStats, func(i, j int) bool {
		return stats.MonthlyStats[i].Month < stats.MonthlyStats[j].Month
	})

	return stats, nil
}

// formatDuration 格式化时间间隔为中文
func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%d天%d小时前", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%d小时%d分钟前", hours, minutes)
	}
	if minutes > 0 {
		return fmt.Sprintf("%d分钟前", minutes)
	}
	return "刚刚"
}
