/**
 * Scada 硬件平台数据模型
 * 定义硬件平台配置的数据结构
 */

package model

/**
 * 平台类型枚举
 */
type PlatformType string

const (
	PlatformTypeBOX     PlatformType = "box"
	PlatformTypeHMI     PlatformType = "hmi"
	PlatformTypeGATEWAY PlatformType = "gateway"
)

/**
 * 硬件平台配置
 */
type HardwarePlatformConfig struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Type        PlatformType `json:"type"`
	Resolution  string       `json:"resolution,omitempty"`
	Enabled     bool         `json:"enabled"`
	Description string       `json:"description,omitempty"`
}

/**
 * 获取所有硬件平台配置（Mock数据）
 */
func GetAllHardwarePlatforms() []HardwarePlatformConfig {
	return []HardwarePlatformConfig{
		{
			ID:          "box1",
			Name:        "BOX1",
			Type:        PlatformTypeBOX,
			Resolution:  "1920x1080",
			Enabled:     true,
			Description: "标准BOX型号",
		},
		{
			ID:          "hmi01",
			Name:        "HMI01",
			Type:        PlatformTypeHMI,
			Resolution:  "1280x800",
			Enabled:     true,
			Description: "HMI触摸屏",
		},
		{
			ID:          "tbox1",
			Name:        "TBOX1",
			Type:        PlatformTypeGATEWAY,
			Resolution:  "1024x600",
			Enabled:     true,
			Description: "网关盒子",
		},
	}
}
