package action

import (
	"context"
	"fmt"
	"log"

	"pansiot-device/internal/alarm/engine"
	"pansiot-device/internal/alarm/rule"
)

// PageJumperHandler 页面跳转处理器（日志模拟）
type PageJumperHandler struct {
	// 预留：未来可接入路由系统
}

// Handle 执行页面跳转
func (h *PageJumperHandler) Handle(ctx context.Context, action *rule.Action, alarm *engine.ActiveAlarm) error {
	// 解析参数
	page, _ := action.Params["page"].(string)
	params, _ := action.Params["params"].(map[string]interface{})

	if page == "" {
		return fmt.Errorf("缺少 page 参数")
	}

	log.Printf("[页面跳转] 目标=%s, 参数=%v, 报警=%s",
		page, params, alarm.Rule.Name)

	// TODO: 未来接入路由系统
	// - 集成前端路由跳转
	// - 支持传递参数（如设备ID、报警ID等）
	// - 支持新窗口打开
	//
	// 示例集成代码：
	// router.Navigate(page, params)

	return nil
}

// Validate 验证动作参数
func (h *PageJumperHandler) Validate(action *rule.Action) error {
	// 检查 page 参数
	page, ok := action.Params["page"]
	if !ok {
		return fmt.Errorf("缺少必需参数: page")
	}

	pagePath, ok := page.(string)
	if !ok {
		return fmt.Errorf("参数 page 类型错误，期望 string")
	}

	if pagePath == "" {
		return fmt.Errorf("参数 page 不能为空")
	}

	return nil
}
