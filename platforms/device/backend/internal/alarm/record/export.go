package record

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"pansiot-device/internal/core"
)

// exportToJSON 导出为JSON格式
func (rm *RecordManager) exportToJSON(records []*AlarmRecord, outputPath string) error {
	if outputPath == "" {
		// 默认输出路径
		outputPath = filepath.Join(rm.config.ExportDir, fmt.Sprintf("alarm_records_%s.json", time.Now().Format("20060102_150405")))
	}

	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 序列化数据
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化数据失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// exportToCSV 导出为CSV格式
func (rm *RecordManager) exportToCSV(records []*AlarmRecord, options *RecordExportOptions) error {
	if options == nil {
		options = &RecordExportOptions{}
	}

	// 确定输出路径
	outputPath := options.OutputPath
	if outputPath == "" {
		outputPath = filepath.Join(rm.config.ExportDir, fmt.Sprintf("alarm_records_%s.csv", time.Now().Format("20060102_150405")))
	}

	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 创建文件
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	if options.IncludeHeader {
		header := []string{
			"记录ID",
			"规则ID",
			"规则名称",
			"事件类型",
			"报警级别",
			"报警状态",
			"报警类别",
			"触发时间",
			"确认时间",
			"恢复时间",
			"持续时间",
			"报警消息",
			"阈值",
			"触发值",
			"责任人",
			"确认用户",
			"云上报",
			"存储位置",
			"创建时间",
		}
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("写入表头失败: %w", err)
		}
	}

	// 写入数据行
	for _, record := range records {
		row, err := rm.recordToCSVRow(record)
		if err != nil {
			return fmt.Errorf("转换记录失败: %w", err)
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("写入数据行失败: %w", err)
		}
	}

	return nil
}

// recordToCSVRow 将报警记录转换为CSV行
func (rm *RecordManager) recordToCSVRow(record *AlarmRecord) ([]string, error) {
	row := []string{
		record.RecordID,
		record.RuleID,
		record.RuleName,
		record.EventType.String(),
		rm.levelToString(record.Level),
		rm.stateToString(record.State),
		record.Category,
		record.TriggerTime.Format("2006-01-02 15:04:05"),
	}

	// 确认时间
	if record.AckTime != nil {
		row = append(row, record.AckTime.Format("2006-01-02 15:04:05"))
	} else {
		row = append(row, "")
	}

	// 恢复时间
	if record.RecoverTime != nil {
		row = append(row, record.RecoverTime.Format("2006-01-02 15:04:05"))
	} else {
		row = append(row, "")
	}

	// 持续时间
	if record.Duration != nil {
		row = append(row, record.Duration.String())
	} else {
		row = append(row, "")
	}

	// 报警消息
	row = append(row, record.AlarmMessage)

	// 阈值
	if record.Threshold != nil {
		row = append(row, fmt.Sprintf("%v", record.Threshold))
	} else {
		row = append(row, "")
	}

	// 触发值
	if record.TriggerValue != nil {
		row = append(row, fmt.Sprintf("%v", record.TriggerValue))
	} else {
		row = append(row, "")
	}

	// 责任人
	if record.ResponsibleUsers != nil && len(record.ResponsibleUsers) > 0 {
		row = append(row, fmt.Sprintf("%v", record.ResponsibleUsers))
	} else {
		row = append(row, "")
	}

	// 确认用户
	row = append(row, record.AckUser)

	// 云上报
	row = append(row, strconv.FormatBool(record.CloudReported))

	// 存储位置
	row = append(row, record.StorageLocation.String())

	// 创建时间
	row = append(row, record.CreatedAt.Format("2006-01-02 15:04:05"))

	return row, nil
}

// levelToString 级别转字符串
func (rm *RecordManager) levelToString(level core.AlarmLevel) string {
	switch level {
	case 1:
		return "提示"
	case 2:
		return "一般"
	case 3:
		return "重要"
	case 4:
		return "严重"
	case 5:
		return "紧急"
	default:
		return fmt.Sprintf("级别%d", level)
	}
}

// stateToString 状态转字符串
func (rm *RecordManager) stateToString(state core.AlarmState) string {
	switch state {
	case core.AlarmStateInactive:
		return "未激活"
	case core.AlarmStateActive:
		return "激活"
	case core.AlarmStateAcknowledged:
		return "已确认"
	case core.AlarmStateCleared:
		return "已清除"
	default:
		return fmt.Sprintf("状态%d", state)
	}
}

// ExportToExcel 导出为Excel格式（预留接口，暂未实现）
func (rm *RecordManager) ExportToExcel(records []*AlarmRecord, outputPath string) error {
	// TODO: 实现 Excel 导出功能
	// 可使用 github.com/xuri/excelize/v2 或类似库
	return fmt.Errorf("Excel导出功能暂未实现")
}

// ExportByDateRange 按日期范围导出
func (rm *RecordManager) ExportByDateRange(startDate, endDate time.Time, format ExportFormat, outputPath string) error {
	query := NewRecordQuery().
		WithStartTime(startDate).
		WithEndTime(endDate)

	_, err := rm.Query(query)
	if err != nil {
		return fmt.Errorf("查询记录失败: %w", err)
	}

	options := &RecordExportOptions{
		Format:        format,
		StartTime:     &startDate,
		EndTime:       &endDate,
		IncludeHeader: true,
		OutputPath:    outputPath,
	}

	return rm.Export(options)
}

// ExportByRule 按规则ID导出
func (rm *RecordManager) ExportByRule(ruleID string, format ExportFormat, outputPath string) error {
	query := NewRecordQuery().
		WithRuleIDs(ruleID)

	_, err := rm.Query(query)
	if err != nil {
		return fmt.Errorf("查询记录失败: %w", err)
	}

	options := &RecordExportOptions{
		Format:        format,
		IncludeHeader: true,
		OutputPath:    outputPath,
	}

	return rm.Export(options)
}

// ExportByLevel 按级别导出
func (rm *RecordManager) ExportByLevel(level core.AlarmLevel, format ExportFormat, outputPath string) error {
	query := NewRecordQuery().
		WithLevels(level)

	_, err := rm.Query(query)
	if err != nil {
		return fmt.Errorf("查询记录失败: %w", err)
	}

	options := &RecordExportOptions{
		Format:        format,
		IncludeHeader: true,
		OutputPath:    outputPath,
	}

	return rm.Export(options)
}

// ExportDailyReport 导出日报
func (rm *RecordManager) ExportDailyReport(date time.Time) error {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	outputPath := filepath.Join(rm.config.ExportDir, fmt.Sprintf("daily_report_%s.json", date.Format("20060102")))

	return rm.ExportByDateRange(startOfDay, endOfDay, ExportFormatJSON, outputPath)
}

// ExportWeeklyReport 导出周报
func (rm *RecordManager) ExportWeeklyReport(date time.Time) error {
	// 计算周一
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	startOfWeek := date.AddDate(0, 0, -weekday+1)
	startOfDay := time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
	endOfWeek := startOfDay.Add(7 * 24 * time.Hour)

	outputPath := filepath.Join(rm.config.ExportDir, fmt.Sprintf("weekly_report_%s.json", date.Format("20060102")))

	return rm.ExportByDateRange(startOfDay, endOfWeek, ExportFormatJSON, outputPath)
}

// ExportMonthlyReport 导出月报
func (rm *RecordManager) ExportMonthlyReport(date time.Time) error {
	startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	outputPath := filepath.Join(rm.config.ExportDir, fmt.Sprintf("monthly_report_%s.json", date.Format("200601")))

	return rm.ExportByDateRange(startOfMonth, endOfMonth, ExportFormatJSON, outputPath)
}

// ExportStatistics 导出统计信息
func (rm *RecordManager) ExportStatistics(outputPath string) error {
	stats, err := rm.GetStats()
	if err != nil {
		return fmt.Errorf("获取统计信息失败: %w", err)
	}

	if outputPath == "" {
		outputPath = filepath.Join(rm.config.ExportDir, fmt.Sprintf("statistics_%s.json", time.Now().Format("20060102_150405")))
	}

	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化统计信息失败: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// GenerateReport 生成综合报告（包含记录和统计）
func (rm *RecordManager) GenerateReport(startDate, endDate time.Time, outputPath string) error {
	if outputPath == "" {
		outputPath = filepath.Join(rm.config.ExportDir, fmt.Sprintf("report_%s_to_%s.json",
			startDate.Format("20060102"), endDate.Format("20060102")))
	}

	// 查询记录
	query := NewRecordQuery().
		WithStartTime(startDate).
		WithEndTime(endDate)

	result, err := rm.Query(query)
	if err != nil {
		return fmt.Errorf("查询记录失败: %w", err)
	}

	// 获取统计信息
	stats, err := rm.GetStats()
	if err != nil {
		return fmt.Errorf("获取统计信息失败: %w", err)
	}

	// 构建报告
	report := map[string]interface{}{
		"start_time": startDate.Format(time.RFC3339),
		"end_time":   endDate.Format(time.RFC3339),
		"total":      result.Total,
		"statistics": stats,
		"records":    result.Records,
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化报告失败: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("写入报告失败: %w", err)
	}

	return nil
}
