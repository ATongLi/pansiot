package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// AuditLogService 审计日志服务
type AuditLogService struct {
	db  *gorm.DB
	rdb *redis.Client
}

// NewAuditLogService 创建审计日志服务
func NewAuditLogService(db *gorm.DB, rdb *redis.Client) *AuditLogService {
	return &AuditLogService{
		db:  db,
		rdb: rdb,
	}
}

// AuditLog 审计日志模型
type AuditLog struct {
	ID           uint            `gorm:"primaryKey" json:"id"`
	TenantID     uint            `gorm:"not null;index:idx_tenant_module" json:"tenant_id"`
	ModuleCode   string          `gorm:"size:50;index:idx_tenant_module" json:"module_code"`
	ActionType   string          `gorm:"size:20;index" json:"action_type"` // CREATE, UPDATE, DELETE, LOGIN, LOGOUT, etc.
	EntityType   string          `gorm:"size:50;index:idx_entity" json:"entity_type"`
	EntityID     *uint           `gorm:"index:idx_entity" json:"entity_id"`
	ActionDetail json.RawMessage `gorm:"type:json" json:"action_detail"`
	OperatorID   uint            `gorm:"not null;index:idx_operator" json:"operator_id"`
	OperatorName string          `gorm:"size:100" json:"operator_name"`
	IPAddress    string          `gorm:"size:45" json:"ip_address"`
	UserAgent    string          `gorm:"size:500" json:"user_agent"`
	Status       string          `gorm:"size:20;index" json:"status"` // SUCCESS, FAILED, PARTIAL
	ErrorMessage string          `gorm:"type:text" json:"error_message"`
	CreatedAt    time.Time       `gorm:"not null;index:idx_created_at" json:"created_at"`
}

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_logs"
}

// ActionDetail 操作详情
type ActionDetail struct {
	Before  interface{} `json:"before"`
	After   interface{} `json:"after"`
	Changes []Change    `json:"changes,omitempty"`
}

// Change 字段变更记录
type Change struct {
	Field string `json:"field"`
	Old   string `json:"old"`
	New   string `json:"new"`
}

// Record 记录审计日志
func (s *AuditLogService) Record(ctx context.Context, log *AuditLog) error {
	// 设置创建时间
	log.CreatedAt = time.Now()

	// 默认状态为成功
	if log.Status == "" {
		log.Status = "SUCCESS"
	}

	// 写入数据库
	if err := s.db.Create(log).Error; err != nil {
		return fmt.Errorf("记录审计日志失败: %w", err)
	}

	return nil
}

// RecordAsync 异步记录审计日志（使用Redis队列）
func (s *AuditLogService) RecordAsync(ctx context.Context, log *AuditLog) error {
	// 序列化日志
	data, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("序列化审计日志失败: %w", err)
	}

	// 推送到Redis队列
	queueKey := "audit_log_queue"
	if err := s.rdb.RPush(ctx, queueKey, data).Err(); err != nil {
		return fmt.Errorf("推送审计日志到队列失败: %w", err)
	}

	return nil
}

// ProcessQueue 处理审计日志队列（由后台任务调用）
func (s *AuditLogService) ProcessQueue(ctx context.Context, batchSize int) error {
	queueKey := "audit_log_queue"

	for i := 0; i < batchSize; i++ {
		// 从队列左侧弹出
		result, err := s.rdb.LPop(ctx, queueKey).Result()
		if err != nil {
			if err == redis.Nil {
				// 队列为空
				return nil
			}
			return fmt.Errorf("从队列弹出日志失败: %w", err)
		}

		// 反序列化
		var log AuditLog
		if err := json.Unmarshal([]byte(result), &log); err != nil {
			continue // 跳过无法解析的日志
		}

		// 写入数据库
		if err := s.db.Create(&log).Error; err != nil {
			// 写入失败，推送到失败队列
			failQueueKey := "audit_log_queue_failed"
			s.rdb.RPush(ctx, failQueueKey, result)
		}
	}

	return nil
}

// Query 查询审计日志
func (s *AuditLogService) Query(ctx context.Context, filter *AuditLogFilter) ([]AuditLog, int64, error) {
	var logs []AuditLog
	var total int64

	query := s.db.Model(&AuditLog{})

	// 租户过滤
	if filter.TenantID != nil {
		query = query.Where("tenant_id = ?", *filter.TenantID)
	}

	// 模块过滤
	if filter.ModuleCode != "" {
		query = query.Where("module_code = ?", filter.ModuleCode)
	}

	// 操作类型过滤
	if filter.ActionType != "" {
		query = query.Where("action_type = ?", filter.ActionType)
	}

	// 实体类型过滤
	if filter.EntityType != "" {
		query = query.Where("entity_type = ?", filter.EntityType)
	}

	// 实体ID过滤
	if filter.EntityID != nil {
		query = query.Where("entity_id = ?", *filter.EntityID)
	}

	// 操作人过滤
	if filter.OperatorID != nil {
		query = query.Where("operator_id = ?", *filter.OperatorID)
	}

	// 时间范围过滤
	if filter.StartTime != nil {
		query = query.Where("created_at >= ?", *filter.StartTime)
	}
	if filter.EndTime != nil {
		query = query.Where("created_at <= ?", *filter.EndTime)
	}

	// 状态过滤
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// 关键字搜索
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("entity_type LIKE ? OR operator_name LIKE ?", keyword, keyword)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询审计日志总数失败: %w", err)
	}

	// 分页查询
	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// 排序
	if filter.OrderBy != "" {
		query = query.Order(filter.OrderBy)
	} else {
		query = query.Order("created_at DESC")
	}

	// 执行查询
	if err := query.Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("查询审计日志失败: %w", err)
	}

	return logs, total, nil
}

// AuditLogFilter 审计日志查询过滤器
type AuditLogFilter struct {
	TenantID    *uint      `json:"tenant_id"`
	ModuleCode  string     `json:"module_code"`
	ActionType  string     `json:"action_type"`
	EntityType  string     `json:"entity_type"`
	EntityID    *uint      `json:"entity_id"`
	OperatorID  *uint      `json:"operator_id"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	Status      string     `json:"status"`
	Keyword     string     `json:"keyword"`
	Page        int        `json:"page"`
	PageSize    int        `json:"page_size"`
	OrderBy     string     `json:"order_by"`
}

// GetByID 根据ID获取审计日志
func (s *AuditLogService) GetByID(ctx context.Context, id uint) (*AuditLog, error) {
	var log AuditLog
	if err := s.db.Where("id = ?", id).First(&log).Error; err != nil {
		return nil, fmt.Errorf("查询审计日志失败: %w", err)
	}
	return &log, nil
}

// Export 导出审计日志
func (s *AuditLogService) Export(ctx context.Context, filter *AuditLogFilter) ([]AuditLog, error) {
	// 设置大页面大小用于导出
	filter.Page = 0
	filter.PageSize = 10000 // 最多导出10000条

	logs, _, err := s.Query(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 脱敏处理
	for i := range logs {
		logs[i] = s.sanitizeLog(logs[i])
	}

	return logs, nil
}

// sanitizeLog 脱敏处理
func (s *AuditLogService) sanitizeLog(log AuditLog) AuditLog {
	// 脱敏ActionDetail中的敏感信息
	if log.ActionDetail != nil {
		var detail ActionDetail
		if err := json.Unmarshal(log.ActionDetail, &detail); err == nil {
			detail = s.sanitizeActionDetail(detail)
			if data, err := json.Marshal(detail); err == nil {
				log.ActionDetail = data
			}
		}
	}

	return log
}

// sanitizeActionDetail 脱敏ActionDetail
func (s *AuditLogService) sanitizeActionDetail(detail ActionDetail) ActionDetail {
	// 对敏感字段进行脱敏
	sensitiveFields := []string{"password", "token", "secret", "key"}

	// 脱敏Before
	if detail.Before != nil {
		if beforeMap, ok := detail.Before.(map[string]interface{}); ok {
			for _, field := range sensitiveFields {
				if _, exists := beforeMap[field]; exists {
					beforeMap[field] = "***"
				}
			}
			detail.Before = beforeMap
		}
	}

	// 脱敏After
	if detail.After != nil {
		if afterMap, ok := detail.After.(map[string]interface{}); ok {
			for _, field := range sensitiveFields {
				if _, exists := afterMap[field]; exists {
					afterMap[field] = "***"
				}
			}
			detail.After = afterMap
		}
	}

	// 脱敏Changes
	for i := range detail.Changes {
		for _, field := range sensitiveFields {
			if detail.Changes[i].Field == field {
				detail.Changes[i].Old = "***"
				detail.Changes[i].New = "***"
			}
		}
	}

	return detail
}

// CleanupOldLogs 清理旧日志（定时任务调用）
func (s *AuditLogService) CleanupOldLogs(ctx context.Context, retentionMonths int) (int64, error) {
	cutoffDate := time.Now().AddDate(0, -retentionMonths, 0)

	result := s.db.Where("created_at < ?", cutoffDate).Delete(&AuditLog{})
	if result.Error != nil {
		return 0, fmt.Errorf("清理旧日志失败: %w", result.Error)
	}

	return result.RowsAffected, nil
}

// GetStatistics 获取审计日志统计
func (s *AuditLogService) GetStatistics(ctx context.Context, tenantID uint, days int) (map[string]int64, error) {
	since := time.Now().AddDate(0, 0, -days)

	var stats []struct {
		ActionType string
		Count      int64
	}

	err := s.db.Model(&AuditLog{}).
		Select("action_type, count(*) as count").
		Where("tenant_id = ? AND created_at >= ?", tenantID, since).
		Group("action_type").
		Scan(&stats).Error

	if err != nil {
		return nil, fmt.Errorf("查询审计日志统计失败: %w", err)
	}

	result := make(map[string]int64)
	for _, stat := range stats {
		result[stat.ActionType] = stat.Count
	}

	return result, nil
}
