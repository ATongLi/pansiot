package action

import (
	"context"
	"fmt"
	"log"

	"pansiot-device/internal/alarm/engine"
	"pansiot-device/internal/alarm/rule"
)

// SoundPlayerHandler 声音播放处理器（日志模拟）
type SoundPlayerHandler struct {
	// 预留：未来可接入真实的声音播放库
	// 例如: github.com/faiface/beep 或 github.com/hajimehoshi/oto
}

// Handle 执行声音播放
func (h *SoundPlayerHandler) Handle(ctx context.Context, action *rule.Action, alarm *engine.ActiveAlarm) error {
	// 参数解析（简化版）
	soundFile, _ := action.Params["file"].(string)
	continuous, _ := action.Params["continuous"].(bool)
	volume, _ := action.Params["volume"].(float64)

	if volume == 0 {
		volume = 1.0 // 默认音量
	}

	log.Printf("[声音播放] 文件=%s, 持续=%v, 音量=%.2f, 报警=%s",
		soundFile, continuous, volume, alarm.Rule.Name)

	// TODO: 未来接入真实声音播放库
	// - 使用 beep、oto 或其他音频库
	// - 支持循环播放和停止控制
	// - 支持音量调节
	// - 支持多通道混音

	return nil
}

// Validate 验证动作参数
func (h *SoundPlayerHandler) Validate(action *rule.Action) error {
	// 检查 file 参数
	file, ok := action.Params["file"]
	if !ok {
		return fmt.Errorf("缺少必需参数: file")
	}

	filePath, ok := file.(string)
	if !ok {
		return fmt.Errorf("参数 file 类型错误，期望 string")
	}

	if filePath == "" {
		return fmt.Errorf("参数 file 不能为空")
	}

	// 可选参数：volume
	if volume, ok := action.Params["volume"]; ok {
		v, ok := volume.(float64)
		if !ok {
			return fmt.Errorf("参数 volume 类型错误，期望 float64")
		}
		if v < 0 || v > 1 {
			return fmt.Errorf("参数 volume 超出范围 [0, 1]: %.2f", v)
		}
	}

	return nil
}
