package action

import (
	"context"
	"log"

	"pansiot-device/internal/alarm/engine"
	"pansiot-device/internal/alarm/rule"
)

// PopupNotifierHandler 弹窗通知处理器（日志模拟，使用静态文本）
type PopupNotifierHandler struct {
	// 预留：未来可接入 WebSocket 前端通讯模块
}

// Handle 执行弹窗通知
func (h *PopupNotifierHandler) Handle(ctx context.Context, action *rule.Action, alarm *engine.ActiveAlarm) error {
	// 解析参数（使用静态文本）
	title, _ := action.Params["title"].(string)
	message, _ := action.Params["message"].(string)

	// 如果没有提供消息，使用报警规则的 TriggerMsg
	if message == "" && alarm.Rule.TriggerMsg.Content != "" {
		message = alarm.Rule.TriggerMsg.Content
	}

	// 默认标题
	if title == "" {
		title = "报警通知"
	}

	log.Printf("[弹窗通知] 标题=%s, 消息=%s, 报警=%s",
		title, message, alarm.Rule.Name)

	// TODO: 未来接入 WebSocket
	// - 通过 FrontendCommunicator.Broadcast 发送弹窗消息
	// - 支持前端展示报警弹窗
	// - 支持确认按钮交互
	// - 集成文本库模块 (internal/nls) 实现多语言支持
	//
	// 未来集成文本库后的示例：
	// textlib.GetTextOrDefault(textID, language, defaultMessage, args)

	return nil
}

// Validate 验证动作参数
func (h *PopupNotifierHandler) Validate(action *rule.Action) error {
	// 所有参数都是可选的
	// title: 弹窗标题（可选，默认"报警通知"）
	// message: 弹窗消息（可选，默认使用报警规则的 TriggerMsg）

	return nil
}
