/**
 * Scada 工程数据模型
 * 定义工程管理功能的所有数据结构
 */

package model

import (
	"errors"
	"time"
)

// 错误定义
var (
	ErrInvalidPassword  = errors.New("invalid password")
	ErrInvalidSignature = errors.New("invalid signature")
)

/**
 * 工程分类
 */
type ProjectCategory string

const (
	Category1 ProjectCategory = "分类1"
	Category2 ProjectCategory = "分类2"
	Custom    ProjectCategory = "自定义分类"
)

/**
 * 硬件平台
 */
type HardwarePlatform string

const (
	HMIModel1    HardwarePlatform = "HMI型号1"
	HMIModel2    HardwarePlatform = "HMI型号2"
	GatewayModel1 HardwarePlatform = "网关型号1"
)

/**
 * 工程元数据
 */
type ProjectMetadata struct {
	Name        string          `json:"name"`
	Author      string          `json:"author,omitempty"`
	Description string          `json:"description,omitempty"`
	Category    string          `json:"category"`
	Platform    HardwarePlatform `json:"platform"`
	CreatedAt   string          `json:"createdAt"`
	UpdatedAt   string          `json:"updatedAt"`
}

/**
 * 工程安全配置
 */
type ProjectSecurity struct {
	Encrypted     bool   `json:"encrypted"`
	Password      string `json:"password,omitempty"`      // 临时字段，用于接收用户密码（创建/打开时），不序列化到文件
	PasswordHash  string `json:"passwordHash,omitempty"`
	DeviceBinding string `json:"deviceBinding,omitempty"`
	FileSignature string `json:"fileSignature"`
	KEKVersion    string `json:"kekVersion,omitempty"`
	// 双重加密机制 - DEK
	UserEncryptedDEK       string `json:"userEncrypted,omitempty"`
	OfficialEncryptedDEK   string `json:"officialEncrypted,omitempty"`
}

/**
 * 画布配置
 */
type CanvasConfig struct {
	Width           int     `json:"width"`
	Height          int     `json:"height"`
	BackgroundColor string  `json:"backgroundColor,omitempty"`
}

/**
 * 数据绑定
 */
type DataBinding struct {
	ComponentID string `json:"componentId"`
	Property    string `json:"property"`
	Source      string `json:"source"`
}

/**
 * 组件定义
 */
type Component struct {
	ID            string         `json:"id"`
	Type          string         `json:"type"`
	X             int            `json:"x"`
	Y             int            `json:"y"`
	Width         int            `json:"width"`
	Height        int            `json:"height"`
	Properties    map[string]any `json:"properties"`
	DataBindings  []DataBinding  `json:"dataBindings,omitempty"`
}

/**
 * 工程主数据结构
 */
type Project struct {
	Version          string             `json:"version"`
	ProjectID        string             `json:"projectId"`
	Metadata         ProjectMetadata    `json:"metadata"`
	Security         ProjectSecurity    `json:"security"`
	Canvas           CanvasConfig       `json:"canvas"`
	Components       []Component        `json:"components"`
	EncryptedContent string             `json:"encryptedContent,omitempty"`
}

/**
 * 最近工程（数据库模型）
 */
type RecentProject struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID   string    `gorm:"uniqueIndex;not null" json:"projectId"`
	Name        string    `gorm:"not null" json:"name"`
	Category    string    `json:"category,omitempty"`
	FilePath    string    `gorm:"not null" json:"filePath"`
	LastOpened  time.Time `gorm:"not null;index" json:"lastOpened"`
	IsEncrypted bool      `gorm:"default:false" json:"isEncrypted"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

/**
 * 自定义分类（数据库模型）
 */
type CustomCategory struct {
	CategoryID string    `gorm:"primaryKey" json:"categoryId"`
	Name       string    `gorm:"not null" json:"name"`
	Color      string    `gorm:"default:'#2196F3'" json:"color"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

/**
 * 应用配置（数据库模型）
 */
type AppConfig struct {
	Key       string    `gorm:"primaryKey" json:"key"`
	Value     string    `gorm:"not null" json:"value"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

/**
 * 用户偏好（数据库模型）
 */
type UserPreference struct {
	Key       string    `gorm:"primaryKey" json:"key"`
	Value     string    `gorm:"not null" json:"value"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

/**
 * 审计日志（数据库模型）
 */
type AuditLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Timestamp time.Time `gorm:"not null;index" json:"timestamp"`
	Operator  string    `gorm:"not null;index" json:"operator"`
	Operation string    `gorm:"not null" json:"operation"`
	TargetID  string    `json:"targetId,omitempty"`
	Details   string    `json:"details,omitempty"`
	Result    string    `json:"result,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

/**
 * KEK版本管理（数据库模型）
 */
type KEKVersion struct {
	VersionID   string    `gorm:"primaryKey" json:"versionId"`
	Version     string    `gorm:"not null" json:"version"`
	CreatedAt   time.Time `gorm:"not null" json:"createdAt"`
	RotatedAt   *time.Time `json:"rotatedAt,omitempty"`
	Status      string    `gorm:"default:'active'" json:"status"` // active, rotated, expired
}

/**
 * 创建工程请求
 */
type CreateProjectRequest struct {
	Metadata  ProjectMetadata `json:"metadata"`
	Security  ProjectSecurity `json:"security"`
	SavePath  string          `json:"savePath"`
}

/**
 * 打开工程请求
 */
type OpenProjectRequest struct {
	FilePath string `json:"filePath"`
	Password string `json:"password,omitempty"`
}

/**
 * 验证密码请求
 */
type ValidatePasswordRequest struct {
	FilePath string `json:"filePath"`
	Password string `json:"password"`
}

/**
 * 保存工程请求
 */
type SaveProjectRequest struct {
	Project Project `json:"project"`
}
