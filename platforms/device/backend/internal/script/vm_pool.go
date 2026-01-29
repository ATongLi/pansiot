package script

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"pansiot-device/internal/core"

	"github.com/dop251/goja"
)

// VMPool VM 池
type VMPool struct {
	pool         chan *PooledVM
	factory      func() *goja.Runtime
	maxSize      int
	maxIdle      time.Duration
	maxLifetime  time.Duration
	mu           sync.RWMutex
	stats        VMPoolStats
	sandbox      *Sandbox
	storage      core.Storage // 存储层，用于设置沙箱环境
}

// PooledVM 池化的 VM
type PooledVM struct {
	VM         *goja.Runtime
	CreatedAt  time.Time
	LastUsed   time.Time
	UsageCount int64
}

// VMPoolStats VM 池统计
type VMPoolStats struct {
	TotalCreated int64
	TotalReused  int64
	TotalExpired int64
	CurrentSize  int
	CurrentIdle  int
}

// NewVMPool 创建 VM 池
func NewVMPool(size int, maxIdle, maxLifetime time.Duration, sandbox *Sandbox) *VMPool {
	pool := &VMPool{
		pool:        make(chan *PooledVM, size),
		maxSize:     size,
		maxIdle:     maxIdle,
		maxLifetime: maxLifetime,
		sandbox:     sandbox,
	}

	// VM 创建工厂函数
	pool.factory = func() *goja.Runtime {
		vm := goja.New()
		atomic.AddInt64(&pool.stats.TotalCreated, 1)
		return vm
	}

	// 启动清理协程
	go pool.cleanupRoutine()

	return pool
}

// Get 获取 VM（从池中复用或新建）
func (p *VMPool) Get() *PooledVM {
	select {
	case vm := <-p.pool:
		// 从池中获取到 VM
		atomic.AddInt64(&vm.UsageCount, 1)
		atomic.AddInt64(&p.stats.TotalReused, 1)
		vm.LastUsed = time.Now()

		// 检查 VM 是否过期
		if time.Since(vm.CreatedAt) > p.maxLifetime {
			p.Destroy(vm)
			atomic.AddInt64(&p.stats.TotalExpired, 1)
			return p.create()
		}

		return vm
	default:
		// 池为空，创建新 VM
		return p.create()
	}
}

// create 创建新 VM
func (p *VMPool) create() *PooledVM {
	vm := &PooledVM{
		VM:        p.factory(),
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
	}

	// 如果有沙箱，设置沙箱环境
	if p.sandbox != nil && p.storage != nil {
		if err := p.sandbox.SetupVM(vm.VM, p.storage); err != nil {
			log.Printf("[VM池] 设置沙箱环境失败: %v", err)
		}
	}

	p.mu.Lock()
	p.stats.CurrentSize++
	p.mu.Unlock()

	return vm
}

// Put 归还 VM 到池中
func (p *VMPool) Put(vm *PooledVM) error {
	if vm == nil || vm.VM == nil {
		return nil
	}

	vm.LastUsed = time.Now()

	select {
	case p.pool <- vm:
		return nil
	default:
		// 池已满，销毁 VM
		p.Destroy(vm)
		return nil
	}
}

// Reset 重置 VM 到初始状态
func (p *VMPool) Reset(vm *PooledVM) error {
	if vm == nil || vm.VM == nil {
		return nil
	}

	// 清理全局变量
	// 注意：Goja 不支持完全重置，这里只是清理一些基本状态
	// 如果需要完全隔离，应该销毁旧 VM 并创建新的

	return nil
}

// Destroy 销毁 VM
func (p *VMPool) Destroy(vm *PooledVM) {
	if vm == nil || vm.VM == nil {
		return
	}

	// Goja 的 VM 会被垃圾回收器自动清理
	// 这里只是移除引用
	vm.VM = nil

	p.mu.Lock()
	p.stats.CurrentSize--
	p.mu.Unlock()
}

// cleanupRoutine 定期清理过期 VM
func (p *VMPool) cleanupRoutine() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		p.Cleanup()
	}
}

// Cleanup 清理过期 VM
func (p *VMPool) Cleanup() {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	cleaned := 0

	// 检查池中的所有 VM
	for i := 0; i < len(p.pool); i++ {
		select {
		case vm := <-p.pool:
			// 检查是否过期
			if now.Sub(vm.LastUsed) > p.maxIdle || now.Sub(vm.CreatedAt) > p.maxLifetime {
				p.Destroy(vm)
				cleaned++
			} else {
				// 未过期，放回池中
				p.pool <- vm
			}
		default:
			break
		}
	}

	if cleaned > 0 {
		log.Printf("[VM池] 清理了 %d 个过期 VM", cleaned)
	}
}

// GetStats 获取统计信息
func (p *VMPool) GetStats() VMPoolStats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	stats := p.stats
	stats.CurrentIdle = len(p.pool)

	return stats
}

// Close 关闭 VM 池
func (p *VMPool) Close() {
	// 清空池中的所有 VM
	close(p.pool)
	for vm := range p.pool {
		p.Destroy(vm)
	}
}

// SetStorage 设置存储层
func (p *VMPool) SetStorage(storage core.Storage) {
	p.storage = storage
}

// GetStorage 获取存储层
func (p *VMPool) GetStorage() core.Storage {
	return p.storage
}
