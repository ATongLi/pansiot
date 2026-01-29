package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"pansiot-cloud/internal/middleware"
	"pansiot-cloud/internal/models"
	"pansiot-cloud/pkg/logger"
	"pansiot-cloud/pkg/utils"
)

type AuthService struct {
	db  *gorm.DB
	rdb *redis.Client
}

// NewAuthService 创建认证服务
func NewAuthService(db *gorm.DB, rdb *redis.Client) *AuthService {
	return &AuthService{db: db, rdb: rdb}
}


// Register 用户注册
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest, ipAddress, userAgent string) (*RegisterResponse, error) {
	// 1. 检查用户名或邮箱是否已存在
	var count int64
	if err := s.db.WithContext(ctx).Model(&models.User{}).
		Where("username = ? OR email = ?", req.Username, req.Email).
		Count(&count).Error; err != nil {
		return nil, fmt.Errorf("检查用户失败: %w", err)
	}
	if count > 0 {
		return nil, ErrUserExists
	}

	// 2. 确定租户
	var tenantID int64
	var tenantType string

	if req.TenantSerial != "" {
		// 加入已有企业
		tenantSvc := NewTenantService(s.db)
		tenant, err := tenantSvc.GetTenantBySerialNumber(ctx, req.TenantSerial)
		if err != nil {
			if errors.Is(err, ErrTenantNotFound) {
				return nil, errors.New("企业序列码不存在")
			}
			return nil, fmt.Errorf("查询企业失败: %w", err)
		}
		tenantID = tenant.ID
		tenantType = tenant.TenantType
	} else if req.TenantName != "" {
		// 创建新企业
		tenantSvc := NewTenantService(s.db)
		createReq := &CreateSubTenantRequest{
			Name:          req.TenantName,
			Industry:      req.Industry,
			ContactPerson: req.Username,
			ContactPhone:  req.Phone,
			ContactEmail:  req.Email,
			Username:      req.Username,
			Email:         req.Email,
			Password:      req.Password,
		}

		// 这里需要父租户ID，对于平台级注册使用默认集成商ID
		// TODO: 从配置获取默认集成商ID或根据邀请码确定
		var parentTenantID int64 = 1 // 默认父租户ID

		result, err := tenantSvc.CreateSubTenant(ctx, parentTenantID, createReq, 0)
		if err != nil {
			return nil, fmt.Errorf("创建企业失败: %w", err)
		}

		tenantID = result.TenantID
		tenantType = "TERMINAL"

		// 生成JWT token
		accessToken, refreshToken, err := middleware.GenerateTokenPair(result.UserID, result.Username, result.TenantID, tenantType)
		if err != nil {
			return nil, fmt.Errorf("生成token失败: %w", err)
		}

		// 直接返回结果（用户和租户已在CreateSubTenant中创建）
		return &RegisterResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    2 * 60 * 60,
			UserID:       result.UserID,
			Username:     result.Username,
			Email:        req.Email,
			RealName:     req.Username,
			TenantID:     result.TenantID,
			TenantName:   req.TenantName,
			SerialNumber: result.SerialNumber,
			TenantType:   tenantType,
		}, nil
	} else {
		return nil, errors.New("必须提供企业序列码或企业名称")
	}

	// 3. 哈希密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 4. 创建用户
	user := &models.User{
		TenantID:     tenantID,
		Username:     req.Username,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
		RealName:     req.Username,
		Status:       "ACTIVE",
	}

	if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
		logger.Error(fmt.Sprintf("创建用户失败: error=%v", err))
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 5. 获取租户信息
	var tenant models.Tenant
	if err := s.db.WithContext(ctx).First(&tenant, tenantID).Error; err != nil {
		return nil, fmt.Errorf("查询租户失败: %w", err)
	}

	// 6. 生成JWT token
	accessToken, refreshToken, err := middleware.GenerateTokenPair(user.ID, user.Username, tenant.ID, tenant.TenantType)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %w", err)
	}

	// 7. 记录日志
	logger.Info(fmt.Sprintf("用户注册成功: username=%s, tenant_id=%d, user_id=%d", req.Username, tenantID, user.ID))

	return &RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    2 * 60 * 60,
		UserID:       user.ID,
		Username:     user.Username,
		Email:        user.Email,
		RealName:     user.RealName,
		TenantID:     tenant.ID,
		TenantName:   tenant.Name,
		SerialNumber: tenant.SerialNumber,
		TenantType:   tenant.TenantType,
	}, nil
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	RealName     string `json:"real_name"`
	TenantID     int64  `json:"tenant_id"`
	TenantName   string `json:"tenant_name"`
	SerialNumber string `json:"serial_number"`
	TenantType   string `json:"tenant_type"`
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *LoginRequest, ipAddress, userAgent string) (*LoginResponse, error) {
	// 1. 查找用户
	var user models.User
	err := s.db.WithContext(ctx).Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 2. 验证密码
	if ok, _ := utils.VerifyPassword(req.Password, user.PasswordHash); !ok {
		return nil, errors.New("密码错误")
	}

	// 3. 查询租户信息
	var tenant models.Tenant
	err = s.db.WithContext(ctx).First(&tenant, user.TenantID).Error
	if err != nil {
		return nil, fmt.Errorf("查询租户失败: %w", err)
	}

	// 4. 生成JWT token
	accessToken, refreshToken, err := middleware.GenerateTokenPair(user.ID, user.Username, tenant.ID, tenant.TenantType)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %w", err)
	}

	// 5. 记录日志
	logger.Info(fmt.Sprintf("用户 %s 登录成功", user.Username))

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    2 * 60 * 60,
		UserID:       user.ID,
		Username:     user.Username,
		Email:        user.Email,
		RealName:     user.RealName,
		TenantID:     tenant.ID,
		TenantName:   tenant.Name,
		SerialNumber: tenant.SerialNumber,
		TenantType:   tenant.TenantType,
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *AuthService) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	var user models.User
	err := s.db.WithContext(ctx).First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
