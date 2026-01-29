package service

import (
	"context"
	"fmt"
	mrand "math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"pansiot-cloud/internal/models"
	"pansiot-cloud/pkg/logger"
	"pansiot-cloud/pkg/utils"
)

// SendVerificationCodeRequest 发送验证码请求
type SendVerificationCodeRequest struct {
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	PhoneCountryCode string `json:"phone_country_code"`
}

// SendVerificationCode 发送验证码
func (s *AuthService) SendVerificationCode(ctx context.Context, req *SendVerificationCodeRequest) error {
	if req.Email == "" && req.Phone == "" {
		return fmt.Errorf("邮箱和手机号至少提供一个")
	}

	// 生成6位验证码
	code := fmt.Sprintf("%06d", mrand.Intn(1000000))

	// 存储到Redis
	cacheKey := ""
	if req.Email != "" {
		cacheKey = fmt.Sprintf("verification_code:email:%s", req.Email)
	} else {
		cacheKey = fmt.Sprintf("verification_code:phone:%s", req.Phone)
	}

	if err := s.rdb.Set(ctx, cacheKey, code, 5*time.Minute).Err(); err != nil {
		return fmt.Errorf("存储验证码失败: %w", err)
	}

	logger.Info(fmt.Sprintf("发送验证码: %s = %s", cacheKey, code))
	return nil
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	VerificationCode string `json:"verification_code"`
	NewPassword      string `json:"new_password"`
}

// ResetPassword 重置密码
func (s *AuthService) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	// 验证验证码
	cacheKey := ""
	if req.Email != "" {
		cacheKey = fmt.Sprintf("verification_code:email:%s", req.Email)
	} else {
		cacheKey = fmt.Sprintf("verification_code:phone:%s", req.Phone)
	}

	storedCode, err := s.rdb.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("验证码不存在或已过期")
		}
		return fmt.Errorf("获取验证码失败: %w", err)
	}

	if storedCode != req.VerificationCode {
		return fmt.Errorf("验证码错误")
	}

	// 查找用户
	var user models.User
	query := s.db.WithContext(ctx)
	if req.Email != "" {
		query = query.Where("email = ?", req.Email)
	} else {
		query = query.Where("phone = ?", req.Phone)
	}

	if err := query.First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("用户不存在")
		}
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 哈希新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("哈希密码失败: %w", err)
	}

	// 更新密码
	user.PasswordHash = hashedPassword
	if err := s.db.WithContext(ctx).Save(&user).Error; err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	// 删除验证码
	s.rdb.Del(ctx, cacheKey)

	logger.Info(fmt.Sprintf("用户 %s 重置密码成功", user.Username))
	return nil
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	UserID      int64  `json:"user_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(ctx context.Context, req *ChangePasswordRequest) error {
	// 查找用户
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, req.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("用户不存在")
		}
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 验证旧密码
	ok, _ := utils.VerifyPassword(req.OldPassword, user.PasswordHash)
	if !ok {
		return fmt.Errorf("旧密码错误")
	}

	// 哈希新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("哈希密码失败: %w", err)
	}

	// 更新密码
	user.PasswordHash = hashedPassword
	if err := s.db.WithContext(ctx).Save(&user).Error; err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	logger.Info(fmt.Sprintf("用户 %d 修改密码成功", req.UserID))
	return nil
}
