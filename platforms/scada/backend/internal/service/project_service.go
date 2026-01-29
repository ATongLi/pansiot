/**
 * Scada 工程管理服务
 * 提供工程业务逻辑处理
 */

package service

import (
	"encoding/json"
	"os"
	"pansiot-scada/internal/model"
	"pansiot-scada/internal/security"
	"time"

	"gorm.io/gorm"
)

/**
 * ProjectService 工程服务
 */
type ProjectService struct {
	db *gorm.DB
}

/**
 * 创建新的 ProjectService 实例
 */
func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{
		db: db,
	}
}

/**
 * CreateProject 创建新工程
 */
func (s *ProjectService) CreateProject(req model.CreateProjectRequest) (*model.Project, error) {
	// 1. 生成UUID
	projectID := security.GenerateUUID()

	// 2. 创建时间戳
	now := security.GetCurrentTimestamp()

	// 3. 构建工程对象
	project := &model.Project{
		Version:   "1.0.0",
		ProjectID: projectID,
		Metadata: model.ProjectMetadata{
			Name:        req.Metadata.Name,
			Author:      req.Metadata.Author,
			Description: req.Metadata.Description,
			Category:    req.Metadata.Category,
			Platform:    req.Metadata.Platform,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		Security: model.ProjectSecurity{
			Encrypted:     req.Security.Encrypted,
			FileSignature: "", // 稍后生成
		},
		Canvas: model.CanvasConfig{
			Width:  1920,
			Height: 1080,
		},
		Components: []model.Component{},
	}

	// 4. 处理加密
	if req.Security.Encrypted {
		// 哈希密码
		hash, err := security.HashPassword(req.Security.Password)
		if err != nil {
			return nil, err
		}
		project.Security.PasswordHash = hash

		// 生成签名密钥
		signatureKey := req.Security.Password + projectID

		// 序列化工程内容
		content, err := json.Marshal(project)
		if err != nil {
			return nil, err
		}

		// 加密内容
		encryptedContent, err := security.EncryptData(content, req.Security.Password)
		if err != nil {
			return nil, err
		}

		project.EncryptedContent = encryptedContent
		project.Security.FileSignature = security.GenerateSignature(content, signatureKey)
	} else {
		// 未加密：直接生成签名
		signatureKey := "pantool-" + projectID
		content, _ := json.Marshal(project)
		project.Security.FileSignature = security.GenerateSignature(content, signatureKey)
	}

	// 5. 序列化为JSON
	jsonData, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return nil, err
	}

	// 6. 写入文件
	if err := os.WriteFile(req.SavePath, jsonData, 0644); err != nil {
		return nil, err
	}

	// 7. 添加到最近工程列表
	recentProject := model.RecentProject{
		ProjectID:   projectID,
		Name:        req.Metadata.Name,
		Category:    req.Metadata.Category,
		FilePath:    req.SavePath,
		LastOpened:  time.Now(),
		IsEncrypted: req.Security.Encrypted,
	}

	if err := s.AddOrUpdateRecentProject(recentProject); err != nil {
		// 记录错误但不影响工程创建
		println("Warning: failed to add to recent projects:", err.Error())
	}

	return project, nil
}

/**
 * OpenProject 打开工程
 */
func (s *ProjectService) OpenProject(filePath string, password string) (*model.Project, error) {
	// 1. 读取文件
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 2. 解析JSON
	var project model.Project
	if err := json.Unmarshal(jsonData, &project); err != nil {
		return nil, err
	}

	// 3. 如果加密，解密内容
	if project.Security.Encrypted {
		if password == "" {
			return nil, model.ErrInvalidPassword
		}

		// 验证密码
		if !security.VerifyPassword(password, project.Security.PasswordHash) {
			return nil, model.ErrInvalidPassword
		}

		// 解密内容
		decryptedData, err := security.DecryptData(project.EncryptedContent, password)
		if err != nil {
			return nil, model.ErrInvalidPassword
		}

		// 反序列化解密后的内容
		if err := json.Unmarshal(decryptedData, &project); err != nil {
			return nil, err
		}

		// 验证签名
		signatureKey := password + project.ProjectID
		if !security.VerifySignature(decryptedData, project.Security.FileSignature, signatureKey) {
			return nil, model.ErrInvalidSignature
		}
	} else {
		// 未加密：验证签名
		signatureKey := "pantool-" + project.ProjectID
		if !security.VerifySignature(jsonData, project.Security.FileSignature, signatureKey) {
			return nil, model.ErrInvalidSignature
		}
	}

	// 4. 更新最近工程列表
	recentProject := model.RecentProject{
		ProjectID:   project.ProjectID,
		Name:        project.Metadata.Name,
		Category:    project.Metadata.Category,
		FilePath:    filePath,
		LastOpened:  time.Now(),
		IsEncrypted: project.Security.Encrypted,
	}

	if err := s.AddOrUpdateRecentProject(recentProject); err != nil {
		println("Warning: failed to update recent projects:", err.Error())
	}

	return &project, nil
}

/**
 * SaveProject 保存工程
 */
func (s *ProjectService) SaveProject(project model.Project) (string, error) {
	// 1. 更新时间戳
	project.Metadata.UpdatedAt = security.GetCurrentTimestamp()

	// 2. 查询文件路径（从最近工程列表）
	var recentProject model.RecentProject
	if err := s.db.Where("project_id = ?", project.ProjectID).First(&recentProject).Error; err == nil {
		filePath := recentProject.FilePath

		// 3. 处理加密
		if project.Security.Encrypted {
			// 这里假设密码在创建时已经设置
			// 实际应该从session或其他地方获取
			// 简化实现：重新加密

			signatureKey := "pantool-" + project.ProjectID // 实际应该使用真实密码
			content, _ := json.Marshal(project)

			// 注意：这里简化了，实际需要保存密码
			encryptedContent, err := security.EncryptData(content, "saved-password")
			if err != nil {
				return "", err
			}

			project.EncryptedContent = encryptedContent
			project.Security.FileSignature = security.GenerateSignature(content, signatureKey)
		}

		// 4. 序列化并保存
		jsonData, err := json.MarshalIndent(project, "", "  ")
		if err != nil {
			return "", err
		}

		if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
			return "", err
		}

		return filePath, nil
	}

	return "", gorm.ErrRecordNotFound
}

/**
 * ValidatePassword 验证密码
 */
func (s *ProjectService) ValidatePassword(filePath string, password string) (bool, error) {
	// 读取文件
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	// 解析
	var project model.Project
	if err := json.Unmarshal(jsonData, &project); err != nil {
		return false, err
	}

	// 检查是否加密
	if !project.Security.Encrypted {
		return true, nil // 未加密工程
	}

	// 验证密码
	return security.VerifyPassword(password, project.Security.PasswordHash), nil
}

/**
 * GetRecentProjects 获取最近工程列表
 */
func (s *ProjectService) GetRecentProjects() ([]model.RecentProject, error) {
	var projects []model.RecentProject

	// 按最后打开时间倒序，限制50个
	if err := s.db.Order("last_opened DESC").Limit(50).Find(&projects).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

/**
 * AddOrUpdateRecentProject 添加或更新最近工程
 */
func (s *ProjectService) AddOrUpdateRecentProject(project model.RecentProject) error {
	// 使用UPSERT（ON CONFLICT）
	return s.db.Save(&project).Error
}

/**
 * RemoveRecentProject 删除最近工程
 */
func (s *ProjectService) RemoveRecentProject(projectID string) error {
	return s.db.Where("project_id = ?", projectID).Delete(&model.RecentProject{}).Error
}
