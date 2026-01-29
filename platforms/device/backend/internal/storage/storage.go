package storage

import (
	"pansiot-device/internal/core"
)

// NewStorage 创建新的存储实例
func NewStorage() core.Storage {
	return NewMemoryStorage()
}

// 确保 MemoryStorage 实现了 Storage 接口
var _ core.Storage = (*MemoryStorage)(nil)
