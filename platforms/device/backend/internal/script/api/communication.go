package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/dop251/goja"
)

// CommunicationAPI 通讯API (HTTP、MQTT、Modbus)
// Phase 4: 通讯功能API
type CommunicationAPI struct {
	httpClient       *http.Client
	mqttClient       interface{}              // 等待集成
	modbusClient     interface{}              // 等待集成
	subscriptions    map[string]*MQTTSubscription
	subscriptionsMu  sync.RWMutex
	maxResponseSize  int64                    // 最大响应体大小(字节)
}

// HTTPResponse HTTP响应对象
type HTTPResponse struct {
	Status     int                    // HTTP状态码
	StatusText string                 // 状态文本
	Headers    map[string]string      // 响应头
	Body       string                 // 响应体
	Duration   int64                  // 请求耗时(毫秒)
	Error      string                 // 错误信息(如果有)
}

// MQTTSubscription MQTT订阅
type MQTTSubscription struct {
	ID       string      // 订阅ID
	Topic    string      // 订阅主题
	QoS      byte        // 服务质量
	Callback goja.Value  // 回调函数
}

// NewCommunicationAPI 创建通讯API实例
func NewCommunicationAPI() *CommunicationAPI {
	return &CommunicationAPI{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		subscriptions:   make(map[string]*MQTTSubscription),
		maxResponseSize: 10 * 1024 * 1024, // 10MB
	}
}

// InjectToVM 注入API到VM
func (api *CommunicationAPI) InjectToVM(vm *goja.Runtime) error {
	// 注入 HTTP API
	httpAPI := api.createHTTPAPI(vm)
	vm.Set("HTTP", httpAPI)

	// 注入 MQTT API (预留)
	mqttAPI := api.createMQTTAPI(vm)
	vm.Set("MQTT", mqttAPI)

	// 注入 Modbus API (预留)
	modbusAPI := api.createModbusAPI(vm)
	vm.Set("Modbus", modbusAPI)

	return nil
}

// ==================== HTTP API ====================

// createHTTPAPI 创建HTTP API
func (api *CommunicationAPI) createHTTPAPI(vm *goja.Runtime) map[string]interface{} {
	return map[string]interface{}{
		"get": func(call goja.FunctionCall) goja.Value {
			url := call.Argument(0).String()
			response := api.httpGet(url)
			return vm.ToValue(response)
		},
		"post": func(call goja.FunctionCall) goja.Value {
			url := call.Argument(0).String()
			options := call.Argument(1).Export()
			response := api.httpPost(url, options)
			return vm.ToValue(response)
		},
		"put": func(call goja.FunctionCall) goja.Value {
			url := call.Argument(0).String()
			options := call.Argument(1).Export()
			response := api.httpPut(url, options)
			return vm.ToValue(response)
		},
		"delete": func(call goja.FunctionCall) goja.Value {
			url := call.Argument(0).String()
			response := api.httpDelete(url)
			return vm.ToValue(response)
		},
		"request": func(call goja.FunctionCall) goja.Value {
			options := call.Argument(0).Export()
			response := api.httpRequest(options)
			return vm.ToValue(response)
		},
	}
}

// httpGet 执行 GET 请求
func (api *CommunicationAPI) httpGet(url string) map[string]interface{} {
	return api.doRequest("GET", url, nil, nil)
}

// httpPost 执行 POST 请求
func (api *CommunicationAPI) httpPost(url string, options interface{}) map[string]interface{} {
	opts := api.parseRequestOptions(options)
	return api.doRequest("POST", url, opts.Data, opts)
}

// httpPut 执行 PUT 请求
func (api *CommunicationAPI) httpPut(url string, options interface{}) map[string]interface{} {
	opts := api.parseRequestOptions(options)
	return api.doRequest("PUT", url, opts.Data, opts)
}

// httpDelete 执行 DELETE 请求
func (api *CommunicationAPI) httpDelete(url string) map[string]interface{} {
	return api.doRequest("DELETE", url, nil, nil)
}

// httpRequest 执行通用请求
func (api *CommunicationAPI) httpRequest(options interface{}) map[string]interface{} {
	opts := api.parseRequestOptions(options)
	return api.doRequest(opts.Method, opts.URL, opts.Data, opts)
}

// RequestOptions 请求选项
type RequestOptions struct {
	Method  string                 // HTTP方法
	URL     string                 // 请求URL
	Data    interface{}            // 请求数据
	Headers map[string]string      // 请求头
	Timeout int                    // 超时(毫秒)
}

// parseRequestOptions 解析请求选项
func (api *CommunicationAPI) parseRequestOptions(options interface{}) *RequestOptions {
	opts := &RequestOptions{
		Method:  "GET",
		Timeout: 30000, // 默认30秒
	}

	if options == nil {
		return opts
	}

	optMap, ok := options.(map[string]interface{})
	if !ok {
		return opts
	}

	if v, ok := optMap["url"].(string); ok {
		opts.URL = v
	}
	if v, ok := optMap["method"].(string); ok {
		opts.Method = v
	}
	if v, ok := optMap["data"]; ok {
		opts.Data = v
	}
	if v, ok := optMap["headers"].(map[string]interface{}); ok {
		headers := make(map[string]string)
		for k, val := range v {
			if strVal, ok := val.(string); ok {
				headers[k] = strVal
			}
		}
		opts.Headers = headers
	}
	if v, ok := optMap["timeout"].(int64); ok {
		opts.Timeout = int(v)
	}

	return opts
}

// doRequest 执行HTTP请求
func (api *CommunicationAPI) doRequest(method, url string, data interface{}, opts *RequestOptions) map[string]interface{} {
	startTime := time.Now()

	// 构建请求
	var bodyReader io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return api.errorResponse(err, startTime)
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return api.errorResponse(err, startTime)
	}

	// 设置默认头
	req.Header.Set("Accept", "application/json")
	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// 设置自定义头
	if opts != nil && opts.Headers != nil {
		for k, v := range opts.Headers {
			req.Header.Set(k, v)
		}
	}

	// 设置超时
	timeout := 30 * time.Second
	if opts != nil && opts.Timeout > 0 {
		timeout = time.Duration(opts.Timeout) * time.Millisecond
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 执行请求
	resp, err := api.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return api.errorResponse(err, startTime)
	}
	defer resp.Body.Close()

	// 限制响应体大小
	limitedReader := io.LimitReader(resp.Body, api.maxResponseSize+1)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return api.errorResponse(err, startTime)
	}

	// 检查是否超过大小限制
	if int64(len(body)) > api.maxResponseSize {
		return map[string]interface{}{
			"status":     resp.StatusCode,
			"statusText": resp.Status,
			"headers":    make(map[string]string),
			"body":       "",
			"duration":   time.Since(startTime).Milliseconds(),
			"error":      fmt.Sprintf("响应体过大: %d 字节 (最大 %d 字节)", len(body), api.maxResponseSize),
		}
	}

	// 构建响应头
	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	return map[string]interface{}{
		"status":     resp.StatusCode,
		"statusText": resp.Status,
		"headers":    headers,
		"body":       string(body),
		"duration":   time.Since(startTime).Milliseconds(),
		"error":      nil,
	}
}

// errorResponse 错误响应
func (api *CommunicationAPI) errorResponse(err error, startTime time.Time) map[string]interface{} {
	return map[string]interface{}{
		"status":     0,
		"statusText": "Error",
		"headers":    make(map[string]string),
		"body":       "",
		"duration":   time.Since(startTime).Milliseconds(),
		"error":      err.Error(),
	}
}

// ==================== MQTT API (预留) ====================

// createMQTTAPI 创建MQTT API
func (api *CommunicationAPI) createMQTTAPI(vm *goja.Runtime) map[string]interface{} {
	return map[string]interface{}{
		"publish": func(call goja.FunctionCall) goja.Value {
			topic := call.Argument(0).String()
			message := call.Argument(1).String()
			err := api.mqttPublish(topic, message)
			if err != nil {
				panic(vm.NewGoError(err))
			}
			return goja.Undefined()
		},
		"subscribe": func(call goja.FunctionCall) goja.Value {
			topic := call.Argument(0).String()
			callback := call.Argument(1)
			subId, err := api.mqttSubscribe(topic, callback)
			if err != nil {
				panic(vm.NewGoError(err))
			}
			return vm.ToValue(subId)
		},
		"unsubscribe": func(call goja.FunctionCall) goja.Value {
			subId := call.Argument(0).String()
			err := api.mqttUnsubscribe(subId)
			if err != nil {
				panic(vm.NewGoError(err))
			}
			return goja.Undefined()
		},
		"isConnected": func(call goja.FunctionCall) goja.Value {
			connected := api.mqttIsConnected()
			return vm.ToValue(connected)
		},
	}
}

// mqttPublish 发布MQTT消息 (预留实现)
func (api *CommunicationAPI) mqttPublish(topic, message string) error {
	// TODO: 集成 internal/adapter/mqtt.go
	// 需要等待 adapter 模块完成 MQTT 适配器
	return fmt.Errorf("MQTT功能暂未实现,等待adapter模块集成")
}

// mqttSubscribe 订阅MQTT主题 (预留实现)
func (api *CommunicationAPI) mqttSubscribe(topic string, callback goja.Value) (string, error) {
	// TODO: 集成 internal/adapter/mqtt.go
	// 需要等待 adapter 模块完成 MQTT 适配器
	return "", fmt.Errorf("MQTT功能暂未实现,等待adapter模块集成")
}

// mqttUnsubscribe 取消订阅 (预留实现)
func (api *CommunicationAPI) mqttUnsubscribe(subId string) error {
	// TODO: 集成 internal/adapter/mqtt.go
	// 需要等待 adapter 模块完成 MQTT 适配器
	return fmt.Errorf("MQTT功能暂未实现,等待adapter模块集成")
}

// mqttIsConnected 检查MQTT连接状态 (预留实现)
func (api *CommunicationAPI) mqttIsConnected() bool {
	// TODO: 集成 internal/adapter/mqtt.go
	// 需要等待 adapter 模块完成 MQTT 适配器
	return false
}

// ==================== Modbus API (预留) ====================

// createModbusAPI 创建Modbus API (预留)
func (api *CommunicationAPI) createModbusAPI(vm *goja.Runtime) map[string]interface{} {
	return map[string]interface{}{
		"readHoldingRegisters": func(call goja.FunctionCall) goja.Value {
			device := call.Argument(0).String()
			address := int(call.Argument(1).ToInteger())
			count := int(call.Argument(2).ToInteger())
			values, err := api.modbusReadHoldingRegisters(device, address, count)
			if err != nil {
				panic(vm.NewGoError(err))
			}
			return vm.ToValue(values)
		},
		"readInputRegisters": func(call goja.FunctionCall) goja.Value {
			device := call.Argument(0).String()
			address := int(call.Argument(1).ToInteger())
			count := int(call.Argument(2).ToInteger())
			values, err := api.modbusReadInputRegisters(device, address, count)
			if err != nil {
				panic(vm.NewGoError(err))
			}
			return vm.ToValue(values)
		},
		"writeSingleRegister": func(call goja.FunctionCall) goja.Value {
			device := call.Argument(0).String()
			address := int(call.Argument(1).ToInteger())
			value := int(call.Argument(2).ToInteger())
			err := api.modbusWriteSingleRegister(device, address, value)
			if err != nil {
				panic(vm.NewGoError(err))
			}
			return goja.Undefined()
		},
		"writeMultipleRegisters": func(call goja.FunctionCall) goja.Value {
			device := call.Argument(0).String()
			address := int(call.Argument(1).ToInteger())
			values := call.Argument(2).Export()
			err := api.modbusWriteMultipleRegisters(device, address, values)
			if err != nil {
				panic(vm.NewGoError(err))
			}
			return goja.Undefined()
		},
		"readCoils": func(call goja.FunctionCall) goja.Value {
			device := call.Argument(0).String()
			address := int(call.Argument(1).ToInteger())
			count := int(call.Argument(2).ToInteger())
			coils, err := api.modbusReadCoils(device, address, count)
			if err != nil {
				panic(vm.NewGoError(err))
			}
			return vm.ToValue(coils)
		},
		"writeCoil": func(call goja.FunctionCall) goja.Value {
			device := call.Argument(0).String()
			address := int(call.Argument(1).ToInteger())
			value := call.Argument(2).ToBoolean()
			err := api.modbusWriteCoil(device, address, value)
			if err != nil {
				panic(vm.NewGoError(err))
			}
			return goja.Undefined()
		},
	}
}

// Modbus 方法 (预留实现)
func (api *CommunicationAPI) modbusReadHoldingRegisters(device string, address, count int) ([]int, error) {
	return nil, fmt.Errorf("Modbus功能暂未实现,等待adapter模块实现")
}

func (api *CommunicationAPI) modbusReadInputRegisters(device string, address, count int) ([]int, error) {
	return nil, fmt.Errorf("Modbus功能暂未实现,等待adapter模块实现")
}

func (api *CommunicationAPI) modbusWriteSingleRegister(device string, address, value int) error {
	return fmt.Errorf("Modbus功能暂未实现,等待adapter模块实现")
}

func (api *CommunicationAPI) modbusWriteMultipleRegisters(device string, address int, values interface{}) error {
	return fmt.Errorf("Modbus功能暂未实现,等待adapter模块实现")
}

func (api *CommunicationAPI) modbusReadCoils(device string, address, count int) ([]bool, error) {
	return nil, fmt.Errorf("Modbus功能暂未实现,等待adapter模块实现")
}

func (api *CommunicationAPI) modbusWriteCoil(device string, address int, value bool) error {
	return fmt.Errorf("Modbus功能暂未实现,等待adapter模块实现")
}
