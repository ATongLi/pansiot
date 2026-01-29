package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"pansiot-cloud/internal/models"
	"pansiot-cloud/pkg/logger"
	"pansiot-cloud/pkg/utils"
)

var (
ErrUserExists       = errors.New("用户已存在")
	ErrUserNotFound        = errors.New("用户不存在")
	ErrUserCannotBeDeleted = errors.New("用户不能被删除（系统管理员）")
	ErrLastAdmin         = errors.New("最后一个系统管理员不能删除或禁用")
)

type UserService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// ListUsersRequest 用户列表请求
type ListUsersRequest struct {
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=20"`
	Keyword   string `form:"keyword"`    // 搜索关键词（用户名、邮箱、手机号、真实姓名）
	Status    string `form:"status"`     // 用户状态筛选
	RoleID    int64  `form:"role_id"`    // 按角色筛选
	SortBy    string `form:"sort_by,default=created_at"` // 排序字段
	SortOrder string `form:"sort_order,default=desc"`     // 排序方向
}

// ListUsersResponse 用户列表响应
type ListUsersResponse struct {
	Total int64          `json:"total"`
	Users []models.User  `json:"users"`
	Page  int            `json:"page"`
	PageSize int         `json:"page_size"`
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers(ctx context.Context, tenantID int64, req *ListUsersRequest) (*ListUsersResponse, error) {
	// 构建查询
	query := s.db.WithContext(ctx).Model(&models.User{}).Where("tenant_id = ? AND deleted_at IS NULL", tenantID)

	// 关键词搜索
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR phone LIKE ? OR real_name LIKE ?",
			keyword, keyword, keyword, keyword)
	}

	// 状态筛选
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 角色筛选
	if req.RoleID > 0 {
		// 子查询：查找拥有该角色的用户
		subQuery := s.db.Table("user_roles").
			Select("user_id").
			Where("role_id = ?", req.RoleID)
		query = query.Where("id IN (?)", subQuery)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 排序
	sortBy := req.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := req.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	// 分页
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100 // 最大每页100条
	}

	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	// 查询数据
	var users []models.User
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	// 清除敏感信息
	for i := range users {
		users[i].PasswordHash = ""
	}

	return &ListUsersResponse{
		Total:    total,
		Users:    users,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username         string `json:"username" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Phone            string `json:"phone"`
	PhoneCountryCode string `json:"phone_country_code"`
	Password         string `json:"password" binding:"required,min=8"`
	RealName         string `json:"real_name"`
	Avatar           string `json:"avatar"`
	RoleIDs          []int64 `json:"role_ids"` // 分配的角色列表
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, tenantID int64, req *CreateUserRequest, operatorID int64) (*models.User, error) {
	// 1. 检查用户名是否已存在
	var count int64
	if err := s.db.WithContext(ctx).Model(&models.User{}).
		Where("tenant_id = ? AND username = ?", tenantID, req.Username).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrUserExists
	}

	// 2. 检查邮箱是否已存在
	if err := s.db.WithContext(ctx).Model(&models.User{}).
		Where("tenant_id = ? AND email = ?", tenantID, req.Email).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrUserExists
	}

	// 3. 检查配额（如果有）
	// TODO: 调用配额检查
	// hasEnough, err := CheckQuota(ctx, tenantID, QuotaUsers, 1)
	// if !hasEnough {
	// 	return nil, errors.New("用户配额已用完")
	// }

	// 4. 密码加密
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 5. 创建用户
	user := &models.User{
		TenantID:         tenantID,
		Username:         req.Username,
		Email:            req.Email,
		Phone:            req.Phone,
		PhoneCountryCode: req.PhoneCountryCode,
		PasswordHash:     passwordHash,
		RealName:         req.RealName,
		Avatar:           req.Avatar,
		Status:           "ACTIVE",
	}

	// 6. 使用事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 7. 分配角色
	if len(req.RoleIDs) > 0 {
		for _, roleID := range req.RoleIDs {
			// 验证角色属于该租户
			var role models.Role
			if err := tx.Where("id = ? AND tenant_id = ?", roleID, tenantID).First(&role).Error; err != nil {
				tx.Rollback()
				return nil, errors.New("角色不存在")
			}

			userRole := &models.UserRole{
				UserID: user.ID,
				RoleID: roleID,
			}
			if err := tx.Create(userRole).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	} else {
		// 如果没有指定角色，分配默认角色
		if err := s.assignDefaultRoleInTx(tx, tenantID, user.ID); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 8. 更新配额使用量
	// TODO: UpdateQuotaUsage(ctx, tenantID, QuotaUsers, 1)

	// 9. 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 10. 记录审计日志
	// middleware.AuditCustomOperation(ctx, tenantID, operatorID, "", "USER_MANAGEMENT", "CREATE", "USER", &user.ID, nil, user, nil)

	logger.Info(fmt.Sprintf("创建用户成功: user_id=%d, username=%s, tenant_id=%d", user.ID, user.Username, tenantID))

	// 清除密码
	user.PasswordHash = ""
	return user, nil
}

// assignDefaultRoleInTx 在事务中分配默认角色
func (s *UserService) assignDefaultRoleInTx(tx *gorm.DB, tenantID, userID int64) error {
	var role models.Role
	err := tx.Where("tenant_id = ? AND role_code = ?", tenantID, "NORMAL_USER").First(&role).Error
	if err != nil {
		// 如果没有默认角色，创建一个
		role = models.Role{
			TenantID:    tenantID,
			RoleCode:    "NORMAL_USER",
			RoleName:    "普通用户",
			IsSystem:    false,
			IsDeletable: true,
		}
		if err := tx.Create(&role).Error; err != nil {
			return err
		}
	}

	userRole := &models.UserRole{
		UserID: userID,
		RoleID: role.ID,
	}
	return tx.Create(userRole).Error
}

// GetUser 获取用户详情
func (s *UserService) GetUser(ctx context.Context, tenantID, userID int64) (*models.User, error) {
	var user models.User
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ? AND deleted_at IS NULL", tenantID, userID).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 清除密码
	user.PasswordHash = ""

	return &user, nil
}

// GetUserWithRoles 获取用户及其角色
func (s *UserService) GetUserWithRoles(ctx context.Context, tenantID, userID int64) (*models.User, []int64, error) {
	// 获取用户
	user, err := s.GetUser(ctx, tenantID, userID)
	if err != nil {
		return nil, nil, err
	}

	// 获取用户角色
	var userRoles []models.UserRole
	if err := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&userRoles).Error; err != nil {
		return nil, nil, err
	}

	roleIDs := make([]int64, len(userRoles))
	for i, ur := range userRoles {
		roleIDs[i] = ur.RoleID
	}

	return user, roleIDs, nil
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email            string   `json:"email" binding:"omitempty,email"`
	Phone            string   `json:"phone"`
	PhoneCountryCode string   `json:"phone_country_code"`
	RealName         string   `json:"real_name"`
	Avatar           string   `json:"avatar"`
	Status           string   `json:"status" binding:"omitempty,oneof=ACTIVE SUSPENDED"`
	RoleIDs          []int64  `json:"role_ids"` // 更新角色列表
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(ctx context.Context, tenantID, userID int64, req *UpdateUserRequest, operatorID int64) (*models.User, error) {
	// 1. 检查用户是否存在
	user, err := s.GetUser(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}

	// 2. 如果更新邮箱或手机号，检查是否已被其他用户使用
	if req.Email != "" && req.Email != user.Email {
		var count int64
		if err := s.db.WithContext(ctx).Model(&models.User{}).
			Where("tenant_id = ? AND email = ? AND id != ?", tenantID, req.Email, userID).Count(&count).Error; err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, errors.New("邮箱已被其他用户使用")
		}
	}

	// 3. 使用事务更新
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 4. 更新用户基本信息
	updates := make(map[string]interface{})
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.PhoneCountryCode != "" {
		updates["phone_country_code"] = req.PhoneCountryCode
	}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if len(updates) > 0 {
		if err := tx.Model(user).Updates(updates).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 6. 更新角色（如果提供）
	if req.RoleIDs != nil {
		// 删除旧的角色关联
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// 添加新的角色关联
		for _, roleID := range req.RoleIDs {
			// 验证角色属于该租户
			var role models.Role
			if err := tx.Where("id = ? AND tenant_id = ?", roleID, tenantID).First(&role).Error; err != nil {
				tx.Rollback()
				return nil, errors.New("角色不存在")
			}

			userRole := &models.UserRole{
				UserID: userID,
				RoleID: roleID,
			}
			if err := tx.Create(userRole).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// 7. 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 8. 重新查询用户信息
	user, err = s.GetUser(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}

	// 9. 记录审计日志
	// middleware.AuditCustomOperation(ctx, tenantID, operatorID, "", "USER_MANAGEMENT", "UPDATE", "USER", &userID, oldData, user, nil)

	logger.Info(fmt.Sprintf("更新用户成功: user_id=%d, tenant_id=%d", userID, tenantID))

	return user, nil
}

// DeleteUser 删除用户（软删除）
func (s *UserService) DeleteUser(ctx context.Context, tenantID, userID int64, operatorID int64) error {
	// 1. 检查用户是否存在
	user, err := s.GetUser(ctx, tenantID, userID)
	if err != nil {
		return err
	}

	// 2. 检查是否是系统管理员角色
	// 获取用户的角色
	var userRoles []models.UserRole
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		return err
	}

	for _, ur := range userRoles {
		var role models.Role
		if err := s.db.WithContext(ctx).First(&role, ur.RoleID).Error; err != nil {
			continue
		}
		if role.IsSystem && role.RoleCode == "SYSTEM_ADMIN" {
			// 检查是否是最后一个系统管理员
			var adminCount int64
			s.db.WithContext(ctx).Table("user_roles").
				Joins("JOIN roles ON roles.id = user_roles.role_id").
				Where("roles.role_code = ? AND roles.tenant_id = ?", "SYSTEM_ADMIN", tenantID).
				Count(&adminCount)

			if adminCount <= 1 {
				return ErrLastAdmin
			}
		}
	}

	// 3. 软删除用户
	now := time.Now()
	if err := s.db.WithContext(ctx).
		Model(user).
		Update("deleted_at", now).Error; err != nil {
		return err
	}

	// 4. 记录审计日志
	// middleware.AuditCustomOperation(ctx, tenantID, operatorID, "", "USER_MANAGEMENT", "DELETE", "USER", &userID, user, nil, nil)

	logger.Info(fmt.Sprintf("删除用户成功: user_id=%d, tenant_id=%d", userID, tenantID))

	return nil
}

// BatchDeleteUsersRequest 批量删除用户请求
type BatchDeleteUsersRequest struct {
	UserIDs []int64 `json:"user_ids" binding:"required"`
}

// BatchDeleteUsers 批量删除用户
func (s *UserService) BatchDeleteUsers(ctx context.Context, tenantID int64, req *BatchDeleteUsersRequest, operatorID int64) (int, error) {
	successCount := 0
	var lastErr error

	for _, userID := range req.UserIDs {
		if err := s.DeleteUser(ctx, tenantID, userID, operatorID); err != nil {
			logger.Error(fmt.Sprintf("删除用户失败: user_id=%d, error=%v", userID, err))
			lastErr = err
			continue
		}
		successCount++
	}

	if successCount == 0 {
		return 0, lastErr
	}

	return successCount, nil
}

// UpdateUserStatusRequest 更新用户状态请求
type UpdateUserStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=ACTIVE SUSPENDED"`
}

// UpdateUserStatus 批量更新用户状态
func (s *UserService) UpdateUserStatus(ctx context.Context, tenantID int64, userIDs []int64, status string, operatorID int64) (int, error) {
	// 检查是否有系统管理员
	if status == "SUSPENDED" || status == "DELETED" {
		for _, userID := range userIDs {
			var userRoles []models.UserRole
			if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
				continue
			}

			for _, ur := range userRoles {
				var role models.Role
				if err := s.db.WithContext(ctx).First(&role, ur.RoleID).Error; err != nil {
					continue
				}
				if role.IsSystem && role.RoleCode == "SYSTEM_ADMIN" {
					var adminCount int64
					s.db.WithContext(ctx).Table("user_roles").
						Joins("JOIN roles ON roles.id = user_roles.role_id").
						Where("roles.role_code = ? AND roles.tenant_id = ?", "SYSTEM_ADMIN", tenantID).
						Count(&adminCount)

					if adminCount <= 1 {
						return 0, ErrLastAdmin
					}
				}
			}
		}
	}

	// 批量更新状态
	result := s.db.WithContext(ctx).
		Model(&models.User{}).
		Where("tenant_id = ? AND id IN ?", tenantID, userIDs).
		Update("status", status)

	if result.Error != nil {
		return 0, result.Error
	}

	logger.Info(fmt.Sprintf("批量更新用户状态: tenant_id=%d, user_ids=%v, status=%s, affected=%d",
		tenantID, userIDs, status, result.RowsAffected))

	return int(result.RowsAffected), nil
}

// ResetUserPasswordRequest 重置用户密码请求
type ResetUserPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ResetUserPassword 重置用户密码（管理员功能）
func (s *UserService) ResetUserPassword(ctx context.Context, tenantID, userID int64, req *ResetUserPasswordRequest, operatorID int64) error {
	// 1. 检查用户是否存在
	user, err := s.GetUser(ctx, tenantID, userID)
	if err != nil {
		return err
	}

	// 2. 加密新密码
	passwordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// 3. 更新密码
	if err := s.db.WithContext(ctx).Model(user).Update("password_hash", passwordHash).Error; err != nil {
		return err
	}

	// 4. 记录审计日志
	// middleware.AuditCustomOperation(ctx, tenantID, operatorID, "", "USER_MANAGEMENT", "UPDATE", "USER_PASSWORD", &userID, nil, map[string]interface{}{"user_id": userID}, nil)

	logger.Info(fmt.Sprintf("重置用户密码成功: user_id=%d, tenant_id=%d, operator_id=%d", userID, tenantID, operatorID))

	return nil
}

// GetUserRoles 获取用户的所有角色
func (s *UserService) GetUserRoles(ctx context.Context, tenantID, userID int64) ([]models.Role, error) {
	// 验证用户属于该租户
	var user models.User
	if err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, userID).
		First(&user).Error; err != nil {
		return nil, ErrUserNotFound
	}

	// 查询用户的角色
	var roles []models.Role
	err := s.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.tenant_id = ? AND roles.deleted_at IS NULL", userID, tenantID).
		Find(&roles).Error

	if err != nil {
		return nil, err
	}

	return roles, nil
}

// AssignRolesRequest 分配角色请求
type AssignRolesRequest struct {
	RoleIDs []int64 `json:"role_ids" binding:"required"`
}

// AssignRoles 为用户分配角色
func (s *UserService) AssignRoles(ctx context.Context, tenantID, userID int64, req *AssignRolesRequest, operatorID int64) error {
	// 1. 验证用户存在
	_, err := s.GetUser(ctx, tenantID, userID)
	if err != nil {
		return err
	}

	// 2. 使用事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 3. 删除旧的角色关联
	if err := tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 4. 添加新的角色关联
	for _, roleID := range req.RoleIDs {
		// 验证角色属于该租户
		var role models.Role
		if err := tx.Where("id = ? AND tenant_id = ?", roleID, tenantID).First(&role).Error; err != nil {
			tx.Rollback()
			return errors.New("角色不存在")
		}

		userRole := &models.UserRole{
			UserID: userID,
			RoleID: roleID,
		}
		if err := tx.Create(userRole).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 5. 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// 6. 清除权限缓存
	// middleware.InvalidateUserPermissionsCache(ctx, userID, tenantID)

	logger.Info(fmt.Sprintf("分配用户角色成功: user_id=%d, role_ids=%v", userID, req.RoleIDs))

	return nil
}
