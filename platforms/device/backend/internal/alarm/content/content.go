package content

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"pansiot-device/internal/core"
)

// Renderer 内容渲染器接口
type Renderer interface {
	// Render 渲染内容
	Render(content interface{}, storage core.Storage) (string, error)

	// Validate 验证内容
	Validate(content interface{}) error
}

// TemplateRenderer 模板渲染器
// 支持变量占位符替换：{var:ID} -> 实际变量值
type TemplateRenderer struct {
	// 变量占位符的正则表达式
	varPattern *regexp.Regexp
}

// NewTemplateRenderer 创建模板渲染器
func NewTemplateRenderer() *TemplateRenderer {
	return &TemplateRenderer{
		varPattern: regexp.MustCompile(`\{var:(\d+)\}`),
	}
}

// Render 渲染内容
func (tr *TemplateRenderer) Render(content interface{}, storage core.Storage) (string, error) {
	switch c := content.(type) {
	case string:
		// 静态文本，直接返回
		return c, nil

	case *TemplateContent:
		return tr.renderTemplate(c, storage)

	case *TextLibraryContent:
		return tr.renderTextLibrary(c, storage)

	default:
		return "", fmt.Errorf("不支持的内容类型: %T", content)
	}
}

// renderTemplate 渲染模板内容
func (tr *TemplateRenderer) renderTemplate(tc *TemplateContent, storage core.Storage) (string, error) {
	result := tc.Template

	// 查找所有变量占位符
	matches := tr.varPattern.FindAllStringSubmatch(result, -1)

	// 替换每个占位符
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		// 提取变量ID
		varIDStr := match[1]
		varID, err := strconv.ParseUint(varIDStr, 10, 64)
		if err != nil {
			continue // 跳过无效的变量ID
		}

		// 读取变量值
		variable, err := storage.ReadVar(varID)
		if err != nil {
			// 变量读取失败，使用默认值
			result = strings.Replace(result, match[0], "[N/A]", 1)
			continue
		}

		// 格式化变量值
		valueStr := formatValue(variable.Value)
		result = strings.Replace(result, match[0], valueStr, 1)
	}

	return result, nil
}

// renderTextLibrary 渲染文本库内容
func (tr *TemplateRenderer) renderTextLibrary(tlc *TextLibraryContent, storage core.Storage) (string, error) {
	// 从文本库获取基础文本
	baseText := tlc.GetText()

	// 如果没有变量，直接返回
	if len(tlc.VariableIDs) == 0 {
		return baseText, nil
	}

	// 创建临时模板内容进行渲染
	tc := &TemplateContent{
		Template:    baseText,
		VariableIDs: tlc.VariableIDs,
	}

	return tr.renderTemplate(tc, storage)
}

// Validate 验证内容
func (tr *TemplateRenderer) Validate(content interface{}) error {
	switch c := content.(type) {
	case string:
		if c == "" {
			return fmt.Errorf("静态内容不能为空")
		}
		return nil

	case *TemplateContent:
		if c.Template == "" {
			return fmt.Errorf("模板内容不能为空")
		}
		return nil

	case *TextLibraryContent:
		if c.TextLibraryID == "" {
			return fmt.Errorf("文本库ID不能为空")
		}
		return nil

	default:
		return fmt.Errorf("不支持的内容类型: %T", content)
	}
}

// TemplateContent 模板内容
type TemplateContent struct {
	Template    string   // 模板文本（支持{var:ID}占位符）
	VariableIDs []uint64 // 涉及的变量ID列表
}

// TextLibraryContent 文本库内容
type TextLibraryContent struct {
	TextLibraryID string   // 文本库ID
	Language      string   // 语言（如"zh-CN", "en-US"）
	VariableIDs   []uint64 // 涉及的变量ID列表
}

// GetText 获取文本库的文本
// 实际实现需要从文本库存储中读取
func (tlc *TextLibraryContent) GetText() string {
	// TODO: 从文本库存储中读取
	return tlc.TextLibraryID
}

// TextLibrary 文本库
// 存储多语言的文本内容
type TextLibrary struct {
	Entries map[string]*TextLibraryEntry // key: textID
}

// TextLibraryEntry 文本库条目
type TextLibraryEntry struct {
	TextID   string                 // 文本ID
	Versions map[string]string      // key: language, value: text
	Metadata map[string]interface{} // 元数据
}

// NewTextLibrary 创建新的文本库
func NewTextLibrary() *TextLibrary {
	return &TextLibrary{
		Entries: make(map[string]*TextLibraryEntry),
	}
}

// AddEntry 添加文本条目
func (lib *TextLibrary) AddEntry(entry *TextLibraryEntry) {
	lib.Entries[entry.TextID] = entry
}

// GetText 获取指定语言和ID的文本
func (lib *TextLibrary) GetText(textID, language string) (string, error) {
	entry, ok := lib.Entries[textID]
	if !ok {
		return "", fmt.Errorf("文本ID不存在: %s", textID)
	}

	text, ok := entry.Versions[language]
	if !ok {
		// 如果指定语言不存在，尝试使用默认语言
		text, ok = entry.Versions["default"]
		if !ok {
			return "", fmt.Errorf("语言[%s]不存在且无默认文本", language)
		}
	}

	return text, nil
}

// SetText 设置文本
func (lib *TextLibrary) SetText(textID, language, text string) {
	entry, ok := lib.Entries[textID]
	if !ok {
		entry = &TextLibraryEntry{
			TextID:   textID,
			Versions: make(map[string]string),
			Metadata: make(map[string]interface{}),
		}
		lib.Entries[textID] = entry
	}

	entry.Versions[language] = text
}

// formatValue 格式化变量值
func formatValue(value interface{}) string {
	if value == nil {
		return "null"
	}

	switch v := value.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32:
		return fmt.Sprintf("%.2f", v)
	case float64:
		return fmt.Sprintf("%.2f", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ContentRenderer 内容渲染器
// 封装了模板渲染器和文本库
type ContentRenderer struct {
	templateRenderer *TemplateRenderer
	textLibrary      *TextLibrary
	defaultLanguage  string
}

// NewContentRenderer 创建内容渲染器
func NewContentRenderer(textLibrary *TextLibrary, defaultLanguage string) *ContentRenderer {
	return &ContentRenderer{
		templateRenderer: NewTemplateRenderer(),
		textLibrary:      textLibrary,
		defaultLanguage:  defaultLanguage,
	}
}

// RenderStatic 渲染静态内容
func (cr *ContentRenderer) RenderStatic(text string) string {
	return text
}

// RenderDynamic 渲染动态内容（模板）
func (cr *ContentRenderer) RenderDynamic(template string, storage core.Storage) (string, error) {
	tc := &TemplateContent{
		Template: template,
	}
	return cr.templateRenderer.renderTemplate(tc, storage)
}

// RenderFromTextLibrary 从文本库渲染内容
func (cr *ContentRenderer) RenderFromTextLibrary(textID, language string, storage core.Storage) (string, error) {
	if language == "" {
		language = cr.defaultLanguage
	}

	tlc := &TextLibraryContent{
		TextLibraryID: textID,
		Language:      language,
	}

	return cr.templateRenderer.renderTextLibrary(tlc, storage)
}

// GetTextLibrary 获取文本库
func (cr *ContentRenderer) GetTextLibrary() *TextLibrary {
	return cr.textLibrary
}

// SetDefaultLanguage 设置默认语言
func (cr *ContentRenderer) SetDefaultLanguage(language string) {
	cr.defaultLanguage = language
}
