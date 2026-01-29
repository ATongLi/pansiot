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
	ErrTenantNotFound       = errors.New("企业不存在")
	ErrTenantAlreadyExists   = errors.New("租户已存在")
	ErrInvalidSerialNumber   = errors.New("企业序列码无效")
	ErrCannotModifyTenant    = errors.New("不能修改租户类型和上级租户")
	ErrNotIntegrator         = errors.New("只有集成商可以创建子租户")
	ErrQuotaExceeded         = errors.New("配额已用完")
)

type TenantService struct {
	db *gorm.DB
}

func NewTenantService(db *gorm.DB) *TenantService {
	return &TenantService{
		db: db,
	}
}

// GetTenantByID 根据ID获取租户信息
func (s *TenantService) GetTenantByID(ctx context.Context, tenantID int64) (*models.Tenant, error) {
	var tenant models.Tenant
	err := s.db.WithContext(ctx).Where("id = ?", tenantID).First(&tenant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTenantNotFound
		}
		logger.Error(fmt.Sprintf("获取租户信息失败: tenant_id=%d, error=%v", tenantID, err))
		return nil, err
	}
	return &tenant, nil
}

// GetTenantBySerialNumber 根据序列号获取租户信息
func (s *TenantService) GetTenantBySerialNumber(ctx context.Context, serialNumber string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := s.db.WithContext(ctx).Where("serial_number = ?", serialNumber).First(&tenant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTenantNotFound
		}
		logger.Error(fmt.Sprintf("获取租户信息失败: serial_number=%s, error=%v", serialNumber, err))
		return nil, err
	}
	return &tenant, nil
}

// UpdateTenantRequest 更新租户信息请求
type UpdateTenantRequest struct {
	Name          string `json:"name" binding:"required,max=200"`
	Industry      string `json:"industry" binding:"max=100"`
	ContactPerson string `json:"contact_person" binding:"max=100"`
	ContactPhone  string `json:"contact_phone" binding:"max=20"`
	ContactEmail  string `json:"contact_email" binding:"omitempty,email"`
}

// UpdateTenant 更新租户信息
func (s *TenantService) UpdateTenant(ctx context.Context, tenantID int64, req *UpdateTenantRequest, operatorID int64) error {
	// 1. 获取租户
	_, err := s.GetTenantByID(ctx, tenantID)
	if err != nil {
		return err
	}

	// 2. 更新基本信息
	updates := map[string]interface{}{
		"name":           req.Name,
		"industry":       req.Industry,
		"contact_person": req.ContactPerson,
		"contact_phone":  req.ContactPhone,
		"contact_email":  req.ContactEmail,
		"updated_at":     time.Now(),
	}

	if err := s.db.WithContext(ctx).Model(&models.Tenant{}).
		Where("id = ?", tenantID).
		Updates(updates).Error; err != nil {
		logger.Error(fmt.Sprintf("更新租户信息失败: tenant_id=%d, error=%v", tenantID, err))
		return err
	}

	logger.Info(fmt.Sprintf("更新租户信息成功: tenant_id=%d, operator_id=%d", tenantID, operatorID))
	return nil
}

// CreateSubTenantRequest 创建子租户请求
type CreateSubTenantRequest struct {
	Name          string `json:"name" binding:"required,max=200"`
	Industry      string `json:"industry" binding:"max=100"`
	ContactPerson string `json:"contact_person" binding:"required,max=100"`
	ContactPhone  string `json:"contact_phone" binding:"required,max=20"`
	ContactEmail  string `json:"contact_email" binding:"omitempty,email"`
	Username      string `json:"username" binding:"required,min=3,max=50"`
	Email         string `json:"email" binding:"required,email"`
	Password      string `json:"password" binding:"required,min=8"`
}

// CreateSubTenantResponse 创建子租户响应
type CreateSubTenantResponse struct {
	TenantID     int64  `json:"tenant_id"`
	SerialNumber string `json:"serial_number"`
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
}

// CreateSubTenant 创建子租户（仅集成商）
func (s *TenantService) CreateSubTenant(ctx context.Context, parentTenantID int64, req *CreateSubTenantRequest, operatorID int64) (*CreateSubTenantResponse, error) {
	// 1. 获取父租户信息，验证是否是集成商
	parentTenant, err := s.GetTenantByID(ctx, parentTenantID)
	if err != nil {
		return nil, err
	}
	if parentTenant.TenantType != "INTEGRATOR" {
		return nil, ErrNotIntegrator
	}

	// 2. 检查配额（可选）
	// TODO: 检查子租户数量配额

	// 3. 使用事务创建租户和管理员用户
	var result CreateSubTenantResponse
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建子租户
		subTenant := &models.Tenant{
			SerialNumber:   "", // 将由数据库自动生成
			TenantType:     "TERMINAL",
			ParentTenantID: &parentTenantID,
			Name:           req.Name,
			Industry:       req.Industry,
			ContactPerson:  req.ContactPerson,
			ContactPhone:   req.ContactPhone,
			ContactEmail:   req.ContactEmail,
			Status:         "ACTIVE",
		}

		if err := tx.Create(subTenant).Error; err != nil {
			return err
		}

		// 重新获取以获取自动生成的serial_number
		if err := tx.First(subTenant, subTenant.ID).Error; err != nil {
			return err
		}

		// 创建管理员用户
		passwordHash, err := utils.HashPassword(req.Password)
		if err != nil {
			return err
		}

		adminUser := &models.User{
			TenantID:         subTenant.ID,
			Username:         req.Username,
			Email:            req.Email,
			PasswordHash:     passwordHash,
			RealName:         req.ContactPerson,
			Status:           "ACTIVE",
		}

		if err := tx.Create(adminUser).Error; err != nil {
			return fmt.Errorf("创建管理员用户失败: %w", err)
		}

		// 初始化子租户的角色和权限
		if err := s.initTenantRolesAndPermissions(tx, subTenant.ID, adminUser.ID); err != nil {
			return fmt.Errorf("初始化角色权限失败: %w", err)
		}

		result.TenantID = subTenant.ID
		result.SerialNumber = subTenant.SerialNumber
		result.UserID = adminUser.ID
		result.Username = adminUser.Username

		return nil
	})

	if err != nil {
		logger.Error(fmt.Sprintf("创建子租户失败: parent_tenant_id=%d, error=%v", parentTenantID, err))
		return nil, err
	}

	logger.Info(fmt.Sprintf("创建子租户成功: parent_tenant_id=%d, tenant_id=%d, serial_number=%s, operator_id=%d",
		parentTenantID, result.TenantID, result.SerialNumber, operatorID))

	return &result, nil
}

// initTenantRolesAndPermissions 初始化租户的角色和权限
func (s *TenantService) initTenantRolesAndPermissions(tx *gorm.DB, tenantID, userID int64) error {
	// 1. 创建系统管理员角色
	systemAdminRole := &models.Role{
		TenantID:    tenantID,
		RoleName:    "系统管理员",
		RoleCode:    "SYSTEM_ADMIN",
		IsSystem:    true,
		Description: "系统预设角色，拥有所有权限",
	}
	if err := tx.Create(systemAdminRole).Error; err != nil {
		return err
	}

	// 2. 创建所有权限
	permissions := s.createAllPermissions(tenantID)
	if err := tx.Create(&permissions).Error; err != nil {
		return err
	}

	// 3. 为系统管理员角色分配所有权限
	rolePermissions := make([]models.RolePermission, len(permissions))
	for i, perm := range permissions {
		rolePermissions[i] = models.RolePermission{
			RoleID:       systemAdminRole.ID,
			PermissionID: perm.ID,
		}
	}
	if err := tx.Create(&rolePermissions).Error; err != nil {
		return err
	}

	// 4. 将用户分配到系统管理员角色
	userRole := &models.UserRole{
		UserID: userID,
		RoleID: systemAdminRole.ID,
	}
	if err := tx.Create(userRole).Error; err != nil {
		return err
	}

	return nil
}

// createAllPermissions 创建所有权限
func (s *TenantService) createAllPermissions(tenantID int64) []models.Permission {
	features := []string{
		models.FeatureSystemConfig,
		models.FeatureOrganizationMgmt,
		models.FeatureUserMgmt,
		models.FeatureRoleMgmt,
		models.FeatureDeviceMgmt,
		models.FeatureDataView,
		models.FeatureAlertMgmt,
		models.FeatureQuotaMgmt,
		models.FeatureAuditLogView,
	}

	actions := []string{
		models.ActionView,
		models.ActionCreate,
		models.ActionEdit,
		models.ActionDelete,
		models.ActionExport,
		models.ActionImport,
	}

	permissions := make([]models.Permission, 0, len(features)*len(actions))
	for _, feature := range features {
		for _, action := range actions {
			permissions = append(permissions, models.Permission{
				TenantID:    tenantID,
				FeatureCode: feature,
				ActionCode:  action,
			})
		}
	}

	return permissions
}

// ListSubTenantsRequest 获取子租户列表请求
type ListSubTenantsRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
	Keyword  string `form:"keyword"`    // 搜索企业名称
	Status   string `form:"status"`     // 租户状态
}

// ListSubTenantsResponse 获取子租户列表响应
type ListSubTenantsResponse struct {
	Tenants    []models.Tenant `json:"tenants"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
}

// ListSubTenants 获取子租户列表（仅集成商）
func (s *TenantService) ListSubTenants(ctx context.Context, integratorID int64, req *ListSubTenantsRequest) (*ListSubTenantsResponse, error) {
	// 构建查询：查询所有parent_tenant_id = integratorID的租户
	query := s.db.WithContext(ctx).Model(&models.Tenant{}).Where("parent_tenant_id = ?", integratorID)

	// 关键词搜索
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("company_name LIKE ?", keyword)
	}

	// 状态筛选
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		logger.Error(fmt.Sprintf("统计子租户数量失败: integrator_id=%d, error=%v", integratorID, err))
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
	var tenants []models.Tenant
	if err := query.Find(&tenants).Error; err != nil {
		logger.Error(fmt.Sprintf("获取子租户列表失败: integrator_id=%d, error=%v", integratorID, err))
		return nil, err
	}

	return &ListSubTenantsResponse{
		Tenants:  tenants,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetTenantStats 获取租户统计信息
func (s *TenantService) GetTenantStats(ctx context.Context, tenantID int64) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 用户数量
	var userCount int64
	if err := s.db.WithContext(ctx).Model(&models.User{}).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Count(&userCount).Error; err != nil {
		return nil, err
	}
	stats["user_count"] = userCount

	// 角色数量
	var roleCount int64
	if err := s.db.WithContext(ctx).Model(&models.Role{}).
		Where("tenant_id = ?", tenantID).
		Count(&roleCount).Error; err != nil {
		return nil, err
	}
	stats["role_count"] = roleCount

	// 子租户数量（仅集成商）
	var subTenantCount int64
	if err := s.db.WithContext(ctx).Model(&models.Tenant{}).
		Where("parent_tenant_id = ?", tenantID).
		Count(&subTenantCount).Error; err != nil {
		return nil, err
	}
	stats["sub_tenant_count"] = subTenantCount

	// 设备数量（待实现）
	stats["device_count"] = 0

	return stats, nil
}
