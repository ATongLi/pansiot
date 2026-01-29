package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"pansiot-device/internal/core"
	"pansiot-device/internal/transform"
)

// MQTTTopicConfig MQTT主题配置
type MQTTTopicConfig struct {
	Topic      string
	VariableID uint64
	JSONPath   string
	DataType   core.DataType
}

// MQTTDeviceConfig MQTT设备配置
type MQTTDeviceConfig struct {
	BrokerURL    string
	ClientID     string
	Username     string
	Password     string
	QoS          byte
	CleanSession bool
	Topics       []MQTTTopicConfig
	Timeout      time.Duration
}

// MQTTAdapter MQTT协议适配器
type MQTTAdapter struct {
	*BaseAdapter
	client      mqtt.Client
	config      MQTTDeviceConfig
	transformer *transform.JSONTransformer
	storage     core.Storage // 可选：直接写入存储层
	callback    UpdateCallback
	mu          sync.RWMutex
	lastValues  map[uint64]*core.Variable // 缓存最新值
	stopChan    chan struct{}
}

// UpdateCallback 变量更新回调函数
type UpdateCallback func(variableID uint64, value interface{}, quality core.QualityCode, timestamp time.Time)

// NewMQTTAdapter 创建MQTT适配器
func NewMQTTAdapter(device *core.Device, config MQTTDeviceConfig) *MQTTAdapter {
	return &MQTTAdapter{
		BaseAdapter: NewBaseAdapter(device),
		config:      config,
		transformer: transform.NewJSONTransformer(),
		lastValues:  make(map[uint64]*core.Variable),
		stopChan:    make(chan struct{}),
	}
}

// SetStorage 设置存储层（可选）
func (ma *MQTTAdapter) SetStorage(storage core.Storage) {
	ma.storage = storage
}

// SetUpdateCallback 设置更新回调（可选）
func (ma *MQTTAdapter) SetUpdateCallback(callback UpdateCallback) {
	ma.callback = callback
}

// Connect 连接到MQTT Broker并订阅主题
func (ma *MQTTAdapter) Connect(ctx context.Context, device *core.Device) error {
	// 配置MQTT客户端选项
	opts := mqtt.NewClientOptions()
	opts.AddBroker(ma.config.BrokerURL)
	opts.SetClientID(ma.config.ClientID)
	opts.SetCleanSession(ma.config.CleanSession)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(30 * time.Second)
	opts.SetConnectTimeout(ma.config.Timeout)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(10 * time.Second)

	// 设置用户名和密码
	if ma.config.Username != "" {
		opts.SetUsername(ma.config.Username)
	}
	if ma.config.Password != "" {
		opts.SetPassword(ma.config.Password)
	}

	// 设置连接成功回调
	opts.OnConnect = func(client mqtt.Client) {
		log.Printf("[MQTTAdapter] 已连接到Broker: %s", ma.config.BrokerURL)

		// 订阅所有配置的主题
		for _, topicConfig := range ma.config.Topics {
			token := client.Subscribe(topicConfig.Topic, ma.config.QoS, ma.handleMessage)
			if token.Wait() && token.Error() != nil {
				log.Printf("[MQTTAdapter] 订阅主题失败 %s: %v", topicConfig.Topic, token.Error())
			} else {
				log.Printf("[MQTTAdapter] 已订阅主题: %s", topicConfig.Topic)
			}
		}
	}

	// 设置连接丢失回调
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Printf("[MQTTAdapter] 连接丢失: %v", err)
		ma.mu.Lock()
		ma.connected = false
		ma.mu.Unlock()
	}

	// 创建客户端
	ma.client = mqtt.NewClient(opts)

	// 连接到Broker
	token := ma.client.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("连接MQTT Broker失败: %v", token.Error())
	}

	// 等待连接完成
	for !ma.client.IsConnected() {
		select {
		case <-time.After(100 * time.Millisecond):
			continue
		case <-time.After(5 * time.Second):
			return fmt.Errorf("连接超时")
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	ma.mu.Lock()
	ma.connected = true
	ma.mu.Unlock()

	log.Printf("[MQTTAdapter] 已连接到设备: %s", device.ID)
	return nil
}

// handleMessage 处理接收到的MQTT消息
func (ma *MQTTAdapter) handleMessage(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	payload := msg.Payload()

	log.Printf("[MQTTAdapter] 收到消息 - 主题: %s, 载荷大小: %d bytes", topic, len(payload))

	// 查找对应的主题配置
	var topicConfig *MQTTTopicConfig
	for i := range ma.config.Topics {
		if ma.config.Topics[i].Topic == topic {
			topicConfig = &ma.config.Topics[i]
			break
		}
	}

	if topicConfig == nil {
		log.Printf("[MQTTAdapter] 未找到主题配置: %s", topic)
		return
	}

	// 使用JSON转换器解析数据
	config := transform.JSONTransform{
		JSONPath: topicConfig.JSONPath,
		DataType: topicConfig.DataType,
	}

	value, err := ma.transformer.Transform(payload, config)
	if err != nil {
		log.Printf("[MQTTAdapter] 数据转换失败: %v", err)
		return
	}

	// 创建变量
	variable := &core.Variable{
		ID:        topicConfig.VariableID,
		Value:     value,
		Quality:   core.QualityGood,
		Timestamp: time.Now(),
		DeviceID:  ma.device.ID,
	}

	// 更新缓存
	ma.mu.Lock()
	ma.lastValues[topicConfig.VariableID] = variable
	ma.IncrementReadCount(true)
	ma.mu.Unlock()

	// 直接写入存储层
	if ma.storage != nil {
		if err := ma.storage.WriteVar(variable); err != nil {
			log.Printf("[MQTTAdapter] 写入存储失败: %v", err)
		}
	}

	// 调用回调函数
	if ma.callback != nil {
		ma.callback(variable.ID, variable.Value, variable.Quality, variable.Timestamp)
	}

	log.Printf("[MQTTAdapter] 变量更新成功 - ID: %d, Value: %v", variable.ID, variable.Value)
}

// Disconnect 断开MQTT连接
func (ma *MQTTAdapter) Disconnect(ctx context.Context) error {
	ma.mu.Lock()
	defer ma.mu.Unlock()

	if !ma.connected {
		return nil
	}

	// 取消订阅所有主题
	for _, topicConfig := range ma.config.Topics {
		token := ma.client.Unsubscribe(topicConfig.Topic)
		token.Wait()
	}

	// 断开连接
	ma.client.Disconnect(250)
	ma.connected = false
	close(ma.stopChan)

	log.Printf("[MQTTAdapter] 已断开连接: %s", ma.device.ID)
	return nil
}

// ReadVariable 读取单个变量（返回缓存值）
func (ma *MQTTAdapter) ReadVariable(ctx context.Context, variableID uint64) (*core.Variable, error) {
	ma.mu.RLock()
	defer ma.mu.RUnlock()

	if !ma.connected {
		return nil, fmt.Errorf("MQTT未连接")
	}

	variable, exists := ma.lastValues[variableID]
	if !exists {
		return nil, fmt.Errorf("变量尚未接收数据: %d", variableID)
	}

	return variable, nil
}

// ReadVariables 批量读取变量（返回缓存值）
func (ma *MQTTAdapter) ReadVariables(ctx context.Context, variableIDs []uint64) ([]*core.Variable, error) {
	ma.mu.RLock()
	defer ma.mu.RUnlock()

	if !ma.connected {
		return nil, fmt.Errorf("MQTT未连接")
	}

	variables := make([]*core.Variable, 0, len(variableIDs))
	for _, variableID := range variableIDs {
		variable, exists := ma.lastValues[variableID]
		if !exists {
			return nil, fmt.Errorf("变量尚未接收数据: %d", variableID)
		}
		variables = append(variables, variable)
	}

	return variables, nil
}

// WriteVariable 写入单个变量（发布到MQTT主题）
func (ma *MQTTAdapter) WriteVariable(ctx context.Context, variableID uint64, value interface{}) error {
	ma.mu.Lock()
	defer ma.mu.Unlock()

	if !ma.connected {
		return fmt.Errorf("MQTT未连接")
	}

	// 查找变量对应的主题
	var topicConfig *MQTTTopicConfig
	for i := range ma.config.Topics {
		if ma.config.Topics[i].VariableID == variableID {
			topicConfig = &ma.config.Topics[i]
			break
		}
	}

	if topicConfig == nil {
		return fmt.Errorf("未找到变量配置: %d", variableID)
	}

	// 构建JSON消息
	payload := map[string]interface{}{
		"variableID": variableID,
		"value":      value,
		"timestamp":  time.Now().Unix(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("JSON编码失败: %v", err)
	}

	// 发布消息
	token := ma.client.Publish(topicConfig.Topic, ma.config.QoS, false, data)
	if token.Wait() && token.Error() != nil {
		ma.IncrementWriteCount(false)
		return fmt.Errorf("发布消息失败: %v", token.Error())
	}

	ma.IncrementWriteCount(true)
	log.Printf("[MQTTAdapter] 已发布消息 - 主题: %s, 值: %v", topicConfig.Topic, value)
	return nil
}

// GetProtocol 获取协议类型
func (ma *MQTTAdapter) GetProtocol() string {
	return "mqtt"
}
