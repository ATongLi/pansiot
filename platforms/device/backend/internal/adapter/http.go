package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"pansiot-device/internal/core"
	"pansiot-device/internal/transform"
)

// HTTPEndpointConfig HTTP端点配置
type HTTPEndpointConfig struct {
	Path        string
	Method      string
	VariableID  uint64
	JSONPath    string
	DataType    core.DataType
	QueryParams map[string]string
}

// HTTPDeviceConfig HTTP设备配置
type HTTPDeviceConfig struct {
	BaseURL      string
	PollInterval time.Duration // 轮询间隔
	Timeout      time.Duration // 请求超时
	Headers      map[string]string
	Auth         AuthConfig
	Endpoints    []HTTPEndpointConfig
}

// AuthConfig 认证配置
type AuthConfig struct {
	Type     string // "basic", "bearer", "none"
	Username string
	Password string
	Token    string
}

// HTTPAdapter HTTP协议适配器（轮询模式）
type HTTPAdapter struct {
	*BaseAdapter
	client      *http.Client
	config      HTTPDeviceConfig
	transformer *transform.JSONTransformer
	storage     core.Storage
	callback    UpdateCallback
	mu          sync.RWMutex
	ticker      *time.Ticker
	stopChan    chan struct{}
	lastValues  map[uint64]*core.Variable
	running     atomic.Bool
}

// NewHTTPAdapter 创建HTTP适配器
func NewHTTPAdapter(device *core.Device, config HTTPDeviceConfig) *HTTPAdapter {
	return &HTTPAdapter{
		BaseAdapter: NewBaseAdapter(device),
		config:      config,
		transformer: transform.NewJSONTransformer(),
		stopChan:    make(chan struct{}),
		lastValues:  make(map[uint64]*core.Variable),
		running:     atomic.Bool{},
	}
}

// SetStorage 设置存储层（可选）
func (ha *HTTPAdapter) SetStorage(storage core.Storage) {
	ha.storage = storage
}

// SetUpdateCallback 设置更新回调（可选）
func (ha *HTTPAdapter) SetUpdateCallback(callback UpdateCallback) {
	ha.callback = callback
}

// Connect 连接到HTTP设备（启动轮询）
func (ha *HTTPAdapter) Connect(ctx context.Context, device *core.Device) error {
	// 配置HTTP客户端
	ha.client = &http.Client{
		Timeout: ha.config.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			DisableCompression:  false,
		},
	}

	// 启动轮询协程
	go ha.pollingLoop()

	ha.mu.Lock()
	ha.connected = true
	ha.mu.Unlock()

	log.Printf("[HTTPAdapter] 已启动轮询: %s", device.ID)
	return nil
}

// pollingLoop 轮询循环
func (ha *HTTPAdapter) pollingLoop() {
	ha.ticker = time.NewTicker(ha.config.PollInterval)
	defer ha.ticker.Stop()
	ha.running.Store(true)

	// 立即执行一次轮询
	ha.pollAllEndpoints()

	for {
		select {
		case <-ha.ticker.C:
			if ha.running.Load() {
				ha.pollAllEndpoints()
			}
		case <-ha.stopChan:
			log.Printf("[HTTPAdapter] 轮询已停止: %s", ha.device.ID)
			return
		}
	}
}

// pollAllEndpoints 轮询所有端点
func (ha *HTTPAdapter) pollAllEndpoints() {
	for _, endpoint := range ha.config.Endpoints {
		if err := ha.pollEndpoint(endpoint); err != nil {
			log.Printf("[HTTPAdapter] 轮询端点失败 %s: %v", endpoint.Path, err)
			ha.IncrementReadCount(false)
		}
	}
}

// pollEndpoint 轮询单个端点
func (ha *HTTPAdapter) pollEndpoint(endpoint HTTPEndpointConfig) error {
	// 构建请求URL
	url := ha.config.BaseURL + endpoint.Path

	// 创建请求
	req, err := http.NewRequest(endpoint.Method, url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 添加查询参数
	if len(endpoint.QueryParams) > 0 {
		q := req.URL.Query()
		for key, value := range endpoint.QueryParams {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	// 添加请求头
	for key, value := range ha.config.Headers {
		req.Header.Set(key, value)
	}

	// 添加认证
	if ha.config.Auth.Type == "basic" {
		req.SetBasicAuth(ha.config.Auth.Username, ha.config.Auth.Password)
	} else if ha.config.Auth.Type == "bearer" {
		req.Header.Set("Authorization", "Bearer "+ha.config.Auth.Token)
	}

	// 设置超时
	ctx, cancel := context.WithTimeout(context.Background(), ha.config.Timeout)
	defer cancel()
	req = req.WithContext(ctx)

	// 发送请求
	startTime := time.Now()
	resp, err := ha.client.Do(req)
	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("请求超时")
		}
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP状态码: %d", resp.StatusCode)
	}

	// 读取响应体
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	elapsed := time.Since(startTime)
	log.Printf("[HTTPAdapter] 收到响应 - 路径: %s, 大小: %d bytes, 耗时: %v",
		endpoint.Path, len(data), elapsed)

	// 使用JSON转换器解析数据
	config := transform.JSONTransform{
		JSONPath: endpoint.JSONPath,
		DataType: endpoint.DataType,
	}

	value, err := ha.transformer.Transform(data, config)
	if err != nil {
		return fmt.Errorf("数据转换失败: %v", err)
	}

	// 创建变量
	variable := &core.Variable{
		ID:        endpoint.VariableID,
		Value:     value,
		Quality:   core.QualityGood,
		Timestamp: time.Now(),
		DeviceID:  ha.device.ID,
	}

	// 更新缓存
	ha.mu.Lock()
	ha.lastValues[endpoint.VariableID] = variable
	ha.IncrementReadCount(true)
	ha.mu.Unlock()

	// 直接写入存储层
	if ha.storage != nil {
		if err := ha.storage.WriteVar(variable); err != nil {
			log.Printf("[HTTPAdapter] 写入存储失败: %v", err)
		}
	}

	// 调用回调函数
	if ha.callback != nil {
		ha.callback(variable.ID, variable.Value, variable.Quality, variable.Timestamp)
	}

	log.Printf("[HTTPAdapter] 变量更新成功 - ID: %d, Value: %v", variable.ID, variable.Value)
	return nil
}

// Disconnect 断开HTTP连接（停止轮询）
func (ha *HTTPAdapter) Disconnect(ctx context.Context) error {
	ha.running.Store(false)

	ha.mu.Lock()
	defer ha.mu.Unlock()

	if !ha.connected {
		return nil
	}

	close(ha.stopChan)

	if ha.ticker != nil {
		ha.ticker.Stop()
	}

	if ha.client != nil {
		ha.client.CloseIdleConnections()
	}

	ha.connected = false
	log.Printf("[HTTPAdapter] 已停止轮询: %s", ha.device.ID)
	return nil
}

// ReadVariable 读取单个变量（返回缓存值）
func (ha *HTTPAdapter) ReadVariable(ctx context.Context, variableID uint64) (*core.Variable, error) {
	ha.mu.RLock()
	defer ha.mu.RUnlock()

	if !ha.connected {
		return nil, fmt.Errorf("HTTP未连接")
	}

	variable, exists := ha.lastValues[variableID]
	if !exists {
		return nil, fmt.Errorf("变量尚未采集数据: %d", variableID)
	}

	return variable, nil
}

// ReadVariables 批量读取变量（返回缓存值）
func (ha *HTTPAdapter) ReadVariables(ctx context.Context, variableIDs []uint64) ([]*core.Variable, error) {
	ha.mu.RLock()
	defer ha.mu.RUnlock()

	if !ha.connected {
		return nil, fmt.Errorf("HTTP未连接")
	}

	variables := make([]*core.Variable, 0, len(variableIDs))
	for _, variableID := range variableIDs {
		variable, exists := ha.lastValues[variableID]
		if !exists {
			return nil, fmt.Errorf("变量尚未采集数据: %d", variableID)
		}
		variables = append(variables, variable)
	}

	return variables, nil
}

// WriteVariable 写入单个变量（发送HTTP请求）
func (ha *HTTPAdapter) WriteVariable(ctx context.Context, variableID uint64, value interface{}) error {
	ha.mu.RLock()
	defer ha.mu.RUnlock()

	if !ha.connected {
		return fmt.Errorf("HTTP未连接")
	}

	// 查找变量对应的端点
	var endpoint *HTTPEndpointConfig
	for i := range ha.config.Endpoints {
		if ha.config.Endpoints[i].VariableID == variableID {
			endpoint = &ha.config.Endpoints[i]
			break
		}
	}

	if endpoint == nil {
		return fmt.Errorf("未找到变量配置: %d", variableID)
	}

	// 构建请求体
	payload := map[string]interface{}{
		"variableID": variableID,
		"value":      value,
		"timestamp":  time.Now().Unix(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("JSON编码失败: %v", err)
	}

	// 创建POST请求
	url := ha.config.BaseURL + endpoint.Path
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// 设置超时
	ctx, cancel := context.WithTimeout(context.Background(), ha.config.Timeout)
	defer cancel()
	req = req.WithContext(ctx)

	// 发送请求
	resp, err := ha.client.Do(req)
	if err != nil {
		ha.IncrementWriteCount(false)
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ha.IncrementWriteCount(false)
		return fmt.Errorf("HTTP状态码: %d", resp.StatusCode)
	}

	ha.IncrementWriteCount(true)
	log.Printf("[HTTPAdapter] 已发送数据 - 端点: %s, 值: %v", endpoint.Path, value)
	return nil
}

// GetProtocol 获取协议类型
func (ha *HTTPAdapter) GetProtocol() string {
	return "http"
}
