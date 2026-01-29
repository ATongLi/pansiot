package adapter

import (
	"fmt"
	"time"

	"pansiot-device/internal/core"
)

// AdapterCreator 适配器创建函数类型
type AdapterCreator func(*core.Device) (core.ProtocolAdapter, error)

// AdapterFactory 适配器工厂
type AdapterFactory struct {
	creators map[string]AdapterCreator
}

// NewAdapterFactory 创建适配器工厂
func NewAdapterFactory() *AdapterFactory {
	factory := &AdapterFactory{
		creators: make(map[string]AdapterCreator),
	}

	// 注册内置适配器
	factory.Register("mock", NewMockAdapterCreator())
	factory.Register("mqtt", NewMQTTAdapterCreator())
	factory.Register("http", NewHTTPAdapterCreator())

	// TODO: 注册其他适配器
	// factory.Register("modbus", NewModbusAdapterCreator())
	// factory.Register("opcua", NewOPCUAAdapterCreator())

	return factory
}

// Register 注册适配器创建器
func (af *AdapterFactory) Register(protocol string, creator AdapterCreator) {
	af.creators[protocol] = creator
}

// Create 根据设备配置创建适配器
func (af *AdapterFactory) Create(device *core.Device) (core.ProtocolAdapter, error) {
	creator, exists := af.creators[device.Protocol]
	if !exists {
		return nil, fmt.Errorf("unsupported protocol: %s", device.Protocol)
	}

	return creator(device)
}

// CreateByConfig 根据协议类型和配置创建适配器
func (af *AdapterFactory) CreateByConfig(protocol string, config interface{}) (core.ProtocolAdapter, error) {
	creator, exists := af.creators[protocol]
	if !exists {
		return nil, fmt.Errorf("unsupported protocol: %s", protocol)
	}

	// 这里需要根据config创建device
	// 简化实现：创建一个基本的device
	device := &core.Device{
		Protocol: protocol,
		ID:       fmt.Sprintf("device-%s", protocol),
	}

	return creator(device)
}

// ListSupportedProtocols 列出支持的协议列表
func (af *AdapterFactory) ListSupportedProtocols() []string {
	protocols := make([]string, 0, len(af.creators))
	for protocol := range af.creators {
		protocols = append(protocols, protocol)
	}
	return protocols
}

// NewMockAdapterCreator 创建Mock适配器的工厂函数
func NewMockAdapterCreator() AdapterCreator {
	return func(device *core.Device) (core.ProtocolAdapter, error) {
		config := MockAdapterConfig{
			DeviceID:     device.ID,
			Protocol:     device.Protocol,
			DataDelay:    0,
			ErrorRate:    0.0,
			ValueRange:   [2]float64{0, 100},
			AutoIncrement: false,
		}

		return NewMockAdapter(config), nil
	}
}

// TODO: 添加其他适配器的工厂函数
// func NewModbusAdapterCreator() AdapterCreator { ... }
// func NewOPCUAAdapterCreator() AdapterCreator { ... }

// NewMQTTAdapterCreator 创建MQTT适配器的工厂函数
func NewMQTTAdapterCreator() AdapterCreator {
	return func(device *core.Device) (core.ProtocolAdapter, error) {
		// 从Device.Tags解析MQTT配置
		config, err := parseMQTTConfig(device.Tags)
		if err != nil {
			return nil, fmt.Errorf("解析MQTT配置失败: %v", err)
		}

		return NewMQTTAdapter(device, config), nil
	}
}

// NewHTTPAdapterCreator 创建HTTP适配器的工厂函数
func NewHTTPAdapterCreator() AdapterCreator {
	return func(device *core.Device) (core.ProtocolAdapter, error) {
		// 从Device.Tags解析HTTP配置
		config, err := parseHTTPConfig(device.Tags)
		if err != nil {
			return nil, fmt.Errorf("解析HTTP配置失败: %v", err)
		}

		return NewHTTPAdapter(device, config), nil
	}
}

// parseMQTTConfig 从Tags解析MQTT配置（简化实现）
func parseMQTTConfig(tags map[string]string) (MQTTDeviceConfig, error) {
	config := MQTTDeviceConfig{
		BrokerURL:    tags["broker_url"],
		ClientID:     tags["client_id"],
		Username:     tags["username"],
		Password:     tags["password"],
		QoS:          0,
		CleanSession: true,
		Timeout:      5 * time.Second,
	}

	// 解析QoS
	if qosStr, exists := tags["qos"]; exists {
		fmt.Sscanf(qosStr, "%d", &config.QoS)
	}

	// TODO: 解析topics配置（实际应用中应该从配置文件或数据库加载）
	// 这里提供一个默认的示例配置
	config.Topics = []MQTTTopicConfig{
		{
			Topic:      "sensor/data",
			VariableID: 100001,
			JSONPath:   "$.temperature",
			DataType:   core.DataTypeFloat64,
		},
	}

	return config, nil
}

// parseHTTPConfig 从Tags解析HTTP配置（简化实现）
func parseHTTPConfig(tags map[string]string) (HTTPDeviceConfig, error) {
	config := HTTPDeviceConfig{
		BaseURL: tags["base_url"],
		Headers: make(map[string]string),
	}

	// 解析轮询间隔
	if intervalStr, exists := tags["poll_interval"]; exists {
		interval, err := time.ParseDuration(intervalStr)
		if err == nil {
			config.PollInterval = interval
		}
	}
	if config.PollInterval == 0 {
		config.PollInterval = 10 * time.Second
	}

	// 解析超时
	if timeoutStr, exists := tags["timeout"]; exists {
		timeout, err := time.ParseDuration(timeoutStr)
		if err == nil {
			config.Timeout = timeout
		}
	}
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
	}

	// TODO: 解析endpoints配置（实际应用中应该从配置文件或数据库加载）
	// 这里提供一个默认的示例配置
	config.Endpoints = []HTTPEndpointConfig{
		{
			Path:      "/api/telemetry",
			Method:    "GET",
			VariableID: 200001,
			JSONPath:  "$.temperature",
			DataType:  core.DataTypeFloat64,
		},
	}

	return config, nil
}
