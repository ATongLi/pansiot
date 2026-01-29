package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"pansiot-cloud/internal/models"
	"pansiot-cloud/pkg/logger"
)

var (
	ErrRoleNotFound           = errors.New("角色不存在")
	ErrRoleAlreadyExists      = errors.New("角色已存在")
	ErrRoleInUse              = errors.New("角色正在使用中，无法删除")
	ErrLastSystemAdmin        = errors.New("最后一个系统管理员角色不能删除或禁用")
	ErrCannotModifySystemRole = errors.New("系统预设角色不能修改或删除")
)

// RoleService 角色服务
type RoleService struct {
	db *gorm.DB
}

// NewRoleService 创建角色服务
func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{
		db: db,
	}
}

// ListRolesRequest 获取角色列表请求
type ListRolesRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
	Keyword  string `form:"keyword"`   // 搜索关键词
	IsSystem *bool  `form:"is_system"` // 是否系统角色
	Status   string `form:"status"`    // 角色状态
}

// ListRolesResponse 获取角色列表响应
type ListRolesResponse struct {
	Roles    []models.Role `json:"roles"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// ListRoles 获取角色列表
func (s *RoleService) ListRoles(ctx context.Context, tenantID int64, req *ListRolesRequest) (*ListRolesResponse, error) {
	// 构建查询
	query := s.db.WithContext(ctx).Model(&models.Role{}).Where("tenant_id = ?", tenantID)

	// 关键词搜索
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("role_name LIKE ? OR role_code LIKE ?", keyword, keyword)
	}

	// 是否系统角色筛选
	if req.IsSystem != nil {
		query = query.Where("is_system = ?", *req.IsSystem)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		logger.Error(fmt.Sprintf("统计角色数量失败: tenant_id=%d, error=%v", tenantID, err))
		return nil, err
	}

	// 分页
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	// 排序和查询
	query = query.Order("created_at DESC")
	var roles []models.Role
	if err := query.Find(&roles).Error; err != nil {
		logger.Error(fmt.Sprintf("获取角色列表失败: tenant_id=%d, error=%v", tenantID, err))
		return nil, err
	}

	return &ListRolesResponse{
		Roles:    roles,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	RoleName      string  `json:"role_name" binding:"required,max=100"`
	RoleCode      string  `json:"role_code" binding:"required,max=50"`
	Description   string  `json:"description" binding:"max=500"`
	PermissionIDs []int64 `json:"permission_ids"`
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(ctx context.Context, tenantID int64, req *CreateRoleRequest, operatorID int64) (*models.Role, error) {
	// 1. 检查角色代码是否已存在
	var count int64
	if err := s.db.WithContext(ctx).Model(&models.Role{}).
		Where("tenant_id = ? AND role_code = ?", tenantID, req.RoleCode).
		Count(&count).Error; err != nil {
		logger.Error(fmt.Sprintf("检查角色代码失败: tenant_id=%d, role_code=%s, error=%v", tenantID, req.RoleCode, err))
		return nil, err
	}
	if count > 0 {
		return nil, ErrRoleAlreadyExists
	}

	// 2. 使用事务创建角色和权限关联
	var role *models.Role
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建角色
		role = &models.Role{
			TenantID:    tenantID,
			RoleName:    req.RoleName,
			RoleCode:    req.RoleCode,
			Description: req.Description,
			IsSystem:    false,
			IsDeletable: true,
		}

		if err := tx.Create(role).Error; err != nil {
			return fmt.Errorf("创建角色失败: %w", err)
		}

		// 分配权限
		if len(req.PermissionIDs) > 0 {
			// 验证权限是否属于该租户
			var permissionCount int64
			if err := tx.Model(&models.Permission{}).
				Where("tenant_id = ? AND id IN ?", tenantID, req.PermissionIDs).
				Count(&permissionCount).Error; err != nil {
				return err
			}
			if int(permissionCount) != len(req.PermissionIDs) {
				return errors.New("部分权限不存在或不属于该租户")
			}

			// 创建角色权限关联
			rolePermissions := make([]models.RolePermission, len(req.PermissionIDs))
			for i, permID := range req.PermissionIDs {
				rolePermissions[i] = models.RolePermission{
					RoleID:       role.ID,
					PermissionID: permID,
				}
			}
			if err := tx.Create(&rolePermissions).Error; err != nil {
				return fmt.Errorf("分配权限失败: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		logger.Error(fmt.Sprintf("创建角色失败: tenant_id=%d, error=%v", tenantID, err))
		return nil, err
	}

	logger.Info(fmt.Sprintf("创建角色成功: tenant_id=%d, role_id=%d, role_code=%s, operator_id=%d",
		tenantID, role.ID, role.RoleCode, operatorID))

	return role, nil
}

// GetRole 获取角色详情
func (s *RoleService) GetRole(ctx context.Context, tenantID, roleID int64) (*models.Role, error) {
	var role models.Role
	err := s.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", roleID, tenantID).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		logger.Error(fmt.Sprintf("获取角色详情失败: role_id=%d, tenant_id=%d, error=%v", roleID, tenantID, err))
		return nil, err
	}
	return &role, nil
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	RoleName      string  `json:"role_name" binding:"required,max=100"`
	Description   string  `json:"description" binding:"max=500"`
	PermissionIDs []int64 `json:"permission_ids"`
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(ctx context.Context, tenantID, roleID int64, req *UpdateRoleRequest, operatorID int64) (*models.Role, error) {
	// 1. 获取角色
	role, err := s.GetRole(ctx, tenantID, roleID)
	if err != nil {
		return nil, err
	}

	// 2. 检查是否是系统角色
	if role.IsSystem {
		return nil, ErrCannotModifySystemRole
	}

	// 3. 检查角色名称是否与其他角色冲突
	var count int64
	if err := s.db.WithContext(ctx).Model(&models.Role{}).
		Where("tenant_id = ? AND role_name = ? AND id != ?", tenantID, req.RoleName, roleID).
		Count(&count).Error; err != nil {
		logger.Error(fmt.Sprintf("检查角色名称失败: tenant_id=%d, role_id=%d, error=%v", tenantID, roleID, err))
		return nil, err
	}
	if count > 0 {
		return nil, ErrRoleAlreadyExists
	}

	// 4. 使用事务更新角色和权限
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新基本信息
		updates := map[string]interface{}{
			"role_name":   req.RoleName,
			"description": req.Description,
			"updated_at":  time.Now(),
		}

		if err := tx.Model(&models.Role{}).
			Where("id = ? AND tenant_id = ?", roleID, tenantID).
			Updates(updates).Error; err != nil {
			return fmt.Errorf("更新角色失败: %w", err)
		}

		// 更新权限关联
		if req.PermissionIDs != nil {
			// 删除旧的权限关联
			if err := tx.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
				return fmt.Errorf("删除旧权限关联失败: %w", err)
			}

			// 添加新的权限关联
			if len(req.PermissionIDs) > 0 {
				// 验证权限是否属于该租户
				var permissionCount int64
				if err := tx.Model(&models.Permission{}).
					Where("tenant_id = ? AND id IN ?", tenantID, req.PermissionIDs).
					Count(&permissionCount).Error; err != nil {
					return err
				}
				if int(permissionCount) != len(req.PermissionIDs) {
					return errors.New("部分权限不存在或不属于该租户")
				}

				rolePermissions := make([]models.RolePermission, len(req.PermissionIDs))
				for i, permID := range req.PermissionIDs {
					rolePermissions[i] = models.RolePermission{
						RoleID:       roleID,
						PermissionID: permID,
					}
				}
				if err := tx.Create(&rolePermissions).Error; err != nil {
					return fmt.Errorf("分配权限失败: %w", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		logger.Error(fmt.Sprintf("更新角色失败: tenant_id=%d, role_id=%d, error=%v", tenantID, roleID, err))
		return nil, err
	}

	// 重新获取更新后的角色
	role, err = s.GetRole(ctx, tenantID, roleID)
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("更新角色成功: tenant_id=%d, role_id=%d, operator_id=%d", tenantID, roleID, operatorID))
	return role, nil
}

// DeleteRoleRequest 删除角色请求
type DeleteRoleRequest struct {
	Force bool `json:"force"` // 是否强制删除
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(ctx context.Context, tenantID, roleID int64, req *DeleteRoleRequest, operatorID int64) error {
	// 1. 获取角色
	role, err := s.GetRole(ctx, tenantID, roleID)
	if err != nil {
		return err
	}

	// 2. 检查是否是系统角色
	if role.IsSystem {
		return ErrCannotModifySystemRole
	}

	// 3. 检查是否有用户使用该角色
	var userRoleCount int64
	if err := s.db.WithContext(ctx).Model(&models.UserRole{}).
		Where("role_id = ?", roleID).
		Count(&userRoleCount).Error; err != nil {
		logger.Error(fmt.Sprintf("检查角色使用情况失败: role_id=%d, error=%v", roleID, err))
		return err
	}

	if userRoleCount > 0 && !req.Force {
		return ErrRoleInUse
	}

	// 4. 使用事务删除角色和相关数据
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 如果是强制删除，先删除用户角色关联
		if req.Force && userRoleCount > 0 {
			if err := tx.Where("role_id = ?", roleID).Delete(&models.UserRole{}).Error; err != nil {
				return fmt.Errorf("删除用户角色关联失败: %w", err)
			}
		}

		// 删除角色权限关联
		if err := tx.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
			return fmt.Errorf("删除角色权限关联失败: %w", err)
		}

		// 软删除角色
		if err := tx.Delete(&models.Role{}, roleID).Error; err != nil {
			return fmt.Errorf("删除角色失败: %w", err)
		}

		return nil
	})

	if err != nil {
		logger.Error(fmt.Sprintf("删除角色失败: tenant_id=%d, role_id=%d, error=%v", tenantID, roleID, err))
		return err
	}

	logger.Info(fmt.Sprintf("删除角色成功: tenant_id=%d, role_id=%d, operator_id=%d", tenantID, roleID, operatorID))
	return nil
}

// GetRolePermissions 获取角色的所有权限
func (s *RoleService) GetRolePermissions(ctx context.Context, tenantID, roleID int64) ([]models.Permission, error) {
	// 1. 验证角色存在
	_, err := s.GetRole(ctx, tenantID, roleID)
	if err != nil {
		return nil, err
	}

	// 2. 查询角色的权限
	var permissions []models.Permission
	err = s.db.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ? AND permissions.tenant_id = ?", roleID, tenantID).
		Find(&permissions).Error

	if err != nil {
		logger.Error(fmt.Sprintf("获取角色权限失败: role_id=%d, tenant_id=%d, error=%v", roleID, tenantID, err))
		return nil, err
	}

	return permissions, nil
}

// AssignPermissionsToRole 为角色分配权限
func (s *RoleService) AssignPermissionsToRole(ctx context.Context, tenantID, roleID int64, permissionIDs []int64, operatorID int64) error {
	// 1. 获取角色
	role, err := s.GetRole(ctx, tenantID, roleID)
	if err != nil {
		return err
	}

	// 2. 检查是否是系统角色
	if role.IsSystem {
		return ErrCannotModifySystemRole
	}

	// 3. 验证权限是否属于该租户
	var permissionCount int64
	if err := s.db.WithContext(ctx).Model(&models.Permission{}).
		Where("tenant_id = ? AND id IN ?", tenantID, permissionIDs).
		Count(&permissionCount).Error; err != nil {
		logger.Error(fmt.Sprintf("验证权限失败: role_id=%d, tenant_id=%d, error=%v", roleID, tenantID, err))
		return err
	}
	if int(permissionCount) != len(permissionIDs) {
		return errors.New("部分权限不存在或不属于该租户")
	}

	// 4. 使用事务更新权限
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除旧的权限关联
		if err := tx.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
			return fmt.Errorf("删除旧权限关联失败: %w", err)
		}

		// 添加新的权限关联
		rolePermissions := make([]models.RolePermission, len(permissionIDs))
		for i, permID := range permissionIDs {
			rolePermissions[i] = models.RolePermission{
				RoleID:       roleID,
				PermissionID: permID,
			}
		}
		if err := tx.Create(&rolePermissions).Error; err != nil {
			return fmt.Errorf("分配权限失败: %w", err)
		}

		return nil
	})

	if err != nil {
		logger.Error(fmt.Sprintf("分配角色权限失败: tenant_id=%d, role_id=%d, error=%v", tenantID, roleID, err))
		return err
	}

	logger.Info(fmt.Sprintf("分配角色权限成功: tenant_id=%d, role_id=%d, permission_count=%d, operator_id=%d",
		tenantID, roleID, len(permissionIDs), operatorID))
	return nil
}

// GetAllPermissions 获取所有可用权限
func (s *RoleService) GetAllPermissions(ctx context.Context, tenantID int64) ([]models.Permission, error) {
	var permissions []models.Permission
	err := s.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("feature_code, action_code").
		Find(&permissions).Error

	if err != nil {
		logger.Error(fmt.Sprintf("获取权限列表失败: tenant_id=%d, error=%v", tenantID, err))
		return nil, err
	}

	return permissions, nil
}
