package core

import (
	"time"
)

// DataType 定义变量数据类型
type DataType int

const (
	DataTypeBool    DataType = iota // 布尔型
	DataTypeInt8                    // 8位整数
	DataTypeInt16                   // 16位整数
	DataTypeInt32                   // 32位整数
	DataTypeInt64                   // 64位整数
	DataTypeUint8                   // 无符号8位整数
	DataTypeUint16                  // 无符号16位整数
	DataTypeUint32                  // 无符号32位整数
	DataTypeUint64                  // 无符号64位整数
	DataTypeFloat32                 // 32位浮点数
	DataTypeFloat64                 // 64位浮点数
	DataTypeString                  // 字符串
	DataTypeBytes                   // 字节数组
)

// String 返回数据类型的字符串表示
func (dt DataType) String() string {
	switch dt {
	case DataTypeBool:
		return "bool"
	case DataTypeInt8:
		return "int8"
	case DataTypeInt16:
		return "int16"
	case DataTypeInt32:
		return "int32"
	case DataTypeInt64:
		return "int64"
	case DataTypeUint8:
		return "uint8"
	case DataTypeUint16:
		return "uint16"
	case DataTypeUint32:
		return "uint32"
	case DataTypeUint64:
		return "uint64"
	case DataTypeFloat32:
		return "float32"
	case DataTypeFloat64:
		return "float64"
	case DataTypeString:
		return "string"
	case DataTypeBytes:
		return "bytes"
	default:
		return "unknown"
	}
}

// QualityCode 定义数据质量码
type QualityCode uint8

const (
	QualityGood         QualityCode = 0 // 数据良好
	QualityBad          QualityCode = 1 // 数据无效
	QualityUncertain    QualityCode = 2 // 数据不确定
	QualityDisconnected QualityCode = 3 // 设备断开
	QualityTimeout      QualityCode = 4 // 采集超时
	QualityOverflow     QualityCode = 5 // 数据溢出
)

// Variable 表示一个实时变量
type Variable struct {
	ID          uint64      // 数字ID，全局唯一
	StringID    string      // 字符串ID，语义化标识
	Name        string      // 变量名称
	Description string      // 变量描述
	DataType    DataType    // 数据类型
	Value       interface{} // 当前值
	Quality     QualityCode // 质量码
	Timestamp   time.Time   // 最后更新时间
	DeviceID    string      // 所属设备ID
	Unit        string      // 单位
	MinValue    interface{} // 最小值（用于校验）
	MaxValue    interface{} // 最大值（用于校验）
}

// SetValue 设置变量值，并进行类型检查
func (v *Variable) SetValue(value interface{}) error {
	// TODO: 实现类型检查和转换
	v.Value = value
	v.Timestamp = time.Now()
	v.Quality = QualityGood
	return nil
}

// Device 表示一个设备
type Device struct {
	ID             string            // 设备ID
	Name           string            // 设备名称
	Description    string            // 设备描述
	Protocol       string            // 协议类型 (modbus, opcua, mqtt, profinet)
	IPAddress      string            // IP地址
	Port           int               // 端口号
	SlaveID        uint8             // 从站ID (Modbus使用)
	ConnectTimeout time.Duration     // 连接超时时间
	ReadTimeout    time.Duration     // 读取超时时间
	Enable         bool              // 是否启用
	Tags           map[string]string // 扩展标签
}

// CollectionTask 表示一个采集任务
type CollectionTask struct {
	ID           string   // 任务ID，格式: TASK_<频率>_<设备ID>_<协议>
	Frequency    int      // 采集频率，单位：毫秒
	DeviceID     string   // 设备ID
	ProtocolType string   // 协议类型
	VariableIDs  []uint64 // 包含的变量ID列表
	Priority     int      // 优先级 (1-10，数字越大优先级越高)
	Timeout      int      // 超时时间，单位：毫秒
	Enable       bool     // 是否启用
}

// AlarmLevel 定义报警级别
type AlarmLevel int

const (
	AlarmLevelLow      AlarmLevel = 1 // 低级报警
	AlarmLevelMedium   AlarmLevel = 2 // 中级报警
	AlarmLevelHigh     AlarmLevel = 3 // 高级报警
	AlarmLevelCritical AlarmLevel = 4 // 严重报警
)

// AlarmState 定义报警状态
type AlarmState int

const (
	AlarmStateInactive     AlarmState = 0 // 未激活
	AlarmStateActive       AlarmState = 1 // 激活
	AlarmStateAcknowledged AlarmState = 2 // 已确认
	AlarmStateCleared      AlarmState = 3 // 已清除
)

// AlarmEvent 表示一个报警事件
type AlarmEvent struct {
	ID          string     // 报警ID
	VariableID  uint64     // 触发报警的变量ID
	Level       AlarmLevel // 报警级别
	State       AlarmState // 报警状态
	Condition   string     // 报警条件表达式 (如 "temp > 30")
	Message     string     // 报警消息
	TriggerTime time.Time  // 触发时间
	AckTime     time.Time  // 确认时间
	ClearTime   time.Time  // 清除时间
	AckUser     string     // 确认用户
}

// AlarmRule 定义报警规则
type AlarmRule struct {
	ID         string     // 规则ID
	Name       string     // 规则名称
	VariableID uint64     // 监控的变量ID
	Condition  string     // 报警条件表达式
	Level      AlarmLevel // 报警级别
	Message    string     // 报警消息模板
	Enable     bool       // 是否启用
	DelayTime  int        // 延时时间，单位：毫秒 (避免瞬时抖动)
}

// SubscriptionMode 订阅模式
type SubscriptionMode int

const (
	SubscriptionModeExact    SubscriptionMode = iota // 精确匹配变量ID
	SubscriptionModePrefix                           // 前缀匹配（如 "DV-PLC001-"）
	SubscriptionModeWildcard                         // 通配符匹配（如 "DV-PLC*-TEMP*"）
	SubscriptionModeDevice                           // 按设备订阅
)

// Subscription 表示一个订阅
type Subscription struct {
	SubscriberID string               // 订阅者ID
	VariableIDs  []uint64             // 订阅的变量ID列表
	Callback     func(VariableUpdate) // 回调函数
}

// VariableUpdate 表示变量更新事件
type VariableUpdate struct {
	VariableID uint64
	Value      interface{}
	Quality    QualityCode
	Timestamp  time.Time
}

// CollectionStats 表示采集统计信息
type CollectionStats struct {
	TotalCount      int64         // 总采集次数
	SuccessCount    int64         // 成功次数
	FailureCount    int64         // 失败次数
	LastCollectTime time.Time     // 最后采集时间
	AvgDuration     time.Duration // 平均耗时
}

// DeviceStats 表示设备统计信息
type DeviceStats struct {
	Online             bool      // 在线状态
	LastConnectTime    time.Time // 最后连接时间
	LastDisconnectTime time.Time // 最后断开时间
	ReconnectCount     int       // 重连次数
	TotalBytes         uint64    // 总字节数
}
