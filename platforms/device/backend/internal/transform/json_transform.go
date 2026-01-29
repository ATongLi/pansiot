package transform

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"pansiot-device/internal/core"
)

// JSONTransform JSON数据转换配置
type JSONTransform struct {
	JSONPath string        // JSONPath表达式，如 "$.data.temperature"
	DataType core.DataType // 目标类型
}

// JSONTransformer JSON转换器
type JSONTransformer struct{}

// NewJSONTransformer 创建JSON转换器
func NewJSONTransformer() *JSONTransformer {
	return &JSONTransformer{}
}

// Transform 执行JSON数据转换
func (jt *JSONTransformer) Transform(data []byte, config JSONTransform) (interface{}, error) {
	// 1. 解析JSON
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v", err)
	}

	// 2. 使用JSONPath提取值
	value, err := jt.extractByPath(jsonData, config.JSONPath)
	if err != nil {
		return nil, fmt.Errorf("JSONPath提取失败: %v", err)
	}

	// 3. 类型转换
	result, err := jt.convertType(value, config.DataType)
	if err != nil {
		return nil, fmt.Errorf("类型转换失败: %v", err)
	}

	return result, nil
}

// extractByPath 根据JSONPath提取值
func (jt *JSONTransformer) extractByPath(data interface{}, jsonPath string) (interface{}, error) {
	// 简化的JSONPath实现
	if jsonPath == "$" {
		return data, nil
	}

	// 去掉前缀 "$."
	path := strings.TrimPrefix(jsonPath, "$.")

	// 按点分割路径
	parts := strings.Split(path, ".")

	current := data
	for _, part := range parts {
		if part == "" {
			continue
		}

		// 处理数组索引，如 sensors[0]
		if strings.Contains(part, "[") {
			var err error
			current, err = jt.extractFromArray(current, part)
			if err != nil {
				return nil, err
			}
			continue
		}

		// 处理对象字段
		if m, ok := current.(map[string]interface{}); ok {
			var exists bool
			current, exists = m[part]
			if !exists {
				return nil, fmt.Errorf("字段不存在: %s", part)
			}
		} else {
			return nil, fmt.Errorf("不是对象类型，无法提取字段: %s", part)
		}
	}

	return current, nil
}

// extractFromArray 从数组中提取元素
func (jt *JSONTransformer) extractFromArray(data interface{}, part string) (interface{}, error) {
	// 解析 sensors[0] 格式
	openBracket := strings.Index(part, "[")
	closeBracket := strings.Index(part, "]")

	if openBracket == -1 || closeBracket == -1 {
		return nil, fmt.Errorf("无效的数组索引格式: %s", part)
	}

	fieldName := part[:openBracket]
	indexStr := part[openBracket+1 : closeBracket]

	// 获取字段
	m, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("不是对象类型，无法访问字段: %s", fieldName)
	}

	arr, exists := m[fieldName]
	if !exists {
		return nil, fmt.Errorf("数组字段不存在: %s", fieldName)
	}

	// 获取数组元素
	slice, ok := arr.([]interface{})
	if !ok {
		return nil, fmt.Errorf("不是数组类型: %s", fieldName)
	}

	// 解析索引
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return nil, fmt.Errorf("无效的数组索引: %s", indexStr)
	}

	if index < 0 || index >= len(slice) {
		return nil, fmt.Errorf("数组索引越界: %d (长度: %d)", index, len(slice))
	}

	return slice[index], nil
}

// convertType 类型转换
func (jt *JSONTransformer) convertType(value interface{}, targetType core.DataType) (interface{}, error) {
	switch targetType {
	case core.DataTypeFloat64:
		return jt.toFloat64(value)
	case core.DataTypeFloat32:
		return jt.toFloat32(value)
	case core.DataTypeInt64:
		return jt.toInt64(value)
	case core.DataTypeInt32:
		return jt.toInt32(value)
	case core.DataTypeInt16:
		return jt.toInt16(value)
	case core.DataTypeInt8:
		return jt.toInt8(value)
	case core.DataTypeString:
		return jt.toString(value)
	case core.DataTypeBool:
		return jt.toBool(value)
	default:
		return nil, fmt.Errorf("不支持的目标类型: %v", targetType)
	}
}

func (jt *JSONTransformer) toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("无法将字符串转换为float64: %s", v)
		}
		return f, nil
	case bool:
		if v {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		return 0, fmt.Errorf("无法转换为float64: %T", value)
	}
}

func (jt *JSONTransformer) toFloat32(value interface{}) (float32, error) {
	f64, err := jt.toFloat64(value)
	if err != nil {
		return 0, err
	}
	return float32(f64), nil
}

func (jt *JSONTransformer) toInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case int32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("无法将字符串转换为int64: %s", v)
		}
		return i, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("无法转换为int64: %T", value)
	}
}

func (jt *JSONTransformer) toInt32(value interface{}) (int32, error) {
	i64, err := jt.toInt64(value)
	if err != nil {
		return 0, err
	}
	return int32(i64), nil
}

func (jt *JSONTransformer) toInt16(value interface{}) (int16, error) {
	i64, err := jt.toInt64(value)
	if err != nil {
		return 0, err
	}
	return int16(i64), nil
}

func (jt *JSONTransformer) toInt8(value interface{}) (int8, error) {
	i64, err := jt.toInt64(value)
	if err != nil {
		return 0, err
	}
	return int8(i64), nil
}

func (jt *JSONTransformer) toString(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case int, int64, int32, int16, int8:
		return fmt.Sprintf("%d", v), nil
	case float64, float32:
		return fmt.Sprintf("%f", v), nil
	case bool:
		return strconv.FormatBool(v), nil
	case nil:
		return "", nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func (jt *JSONTransformer) toBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return false, fmt.Errorf("无法将字符串转换为bool: %s", v)
		}
		return b, nil
	case int, int64, int32, int16, int8:
		return v != 0, nil
	case float64, float32:
		return v != 0.0, nil
	default:
		return false, fmt.Errorf("无法转换为bool: %T", value)
	}
}
