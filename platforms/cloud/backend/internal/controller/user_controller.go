package controller

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"pansiot-cloud/internal/middleware"
	"pansiot-cloud/internal/service"
	"pansiot-cloud/pkg/logger"
	"pansiot-cloud/pkg/response"
)

type UserController struct {
	userService *service.UserService
}

// NewUserController 创建用户控制器
func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Description 分页查询用户列表，支持关键词搜索和状态筛选
// @Tags 用户管理
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Param status query string false "用户状态"
// @Param role_id query int false "角色ID"
// @Success 200 {object} response.Response{data=service.ListUsersResponse}
// @Router /api/v1/users [get]
func (ctrl *UserController) ListUsers(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req service.ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误")
		return
	}

	result, err := ctrl.userService.ListUsers(c.Request.Context(), tenantID, &req)
	if err != nil {
		logger.Error(fmt.Sprintf("获取用户列表失败: tenant_id=%d, error=%v", tenantID, err))
		response.Error(c, 500, "获取用户列表失败")
		return
	}

	response.SuccessWithPagination(c, result.Total, result.Page, result.PageSize, result.Users)
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建新用户并分配角色
// @Tags 用户管理
// @Security Bearer
// @Param request body service.CreateUserRequest true "用户信息"
// @Success 200 {object} response.Response{data=models.User}
// @Router /api/v1/users [post]
func (ctrl *UserController) CreateUser(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	operatorID := middleware.GetUserID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误: "+err.Error())
		return
	}

	user, err := ctrl.userService.CreateUser(c.Request.Context(), tenantID, &req, operatorID)
	if err != nil {
		logger.Error(fmt.Sprintf("创建用户失败: tenant_id=%d, error=%v", tenantID, err))
		if err == service.ErrUserExists {
			response.Error(c, response.UserExistsCode, "用户已存在")
		} else {
			response.Error(c, 500, err.Error())
		}
		return
	}

	response.Success(c, user)
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Description 获取指定用户的详细信息
// @Tags 用户管理
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=models.User}
// @Router /api/v1/users/:id [get]
func (ctrl *UserController) GetUser(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.InvalidParamsCode, "用户ID格式错误")
		return
	}

	user, err := ctrl.userService.GetUser(c.Request.Context(), tenantID, userID)
	if err != nil {
		logger.Error(fmt.Sprintf("获取用户详情失败: user_id=%d, tenant_id=%d, error=%v", userID, tenantID, err))
		if err == service.ErrUserNotFound {
			response.Error(c, 404, "用户不存在")
		} else {
			response.Error(c, 500, "获取用户详情失败")
		}
		return
	}

	response.Success(c, user)
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Description 更新用户信息
// @Tags 用户管理
// @Security Bearer
// @Param id path int true "用户ID"
// @Param request body service.UpdateUserRequest true "用户信息"
// @Success 200 {object} response.Response{data=models.User}
// @Router /api/v1/users/:id [put]
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	operatorID := middleware.GetUserID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.InvalidParamsCode, "用户ID格式错误")
		return
	}

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误: "+err.Error())
		return
	}

	user, err := ctrl.userService.UpdateUser(c.Request.Context(), tenantID, userID, &req, operatorID)
	if err != nil {
		logger.Error(fmt.Sprintf("更新用户失败: user_id=%d, tenant_id=%d, error=%v", userID, tenantID, err))
		if err == service.ErrUserNotFound {
			response.Error(c, 404, "用户不存在")
		} else {
			response.Error(c, 500, err.Error())
		}
		return
	}

	response.Success(c, user)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 软删除指定用户
// @Tags 用户管理
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Router /api/v1/users/:id [delete]
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	operatorID := middleware.GetUserID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.InvalidParamsCode, "用户ID格式错误")
		return
	}

	err = ctrl.userService.DeleteUser(c.Request.Context(), tenantID, userID, operatorID)
	if err != nil {
		logger.Error(fmt.Sprintf("删除用户失败: user_id=%d, tenant_id=%d, error=%v", userID, tenantID, err))
		if err == service.ErrUserNotFound {
			response.Error(c, 404, "用户不存在")
		} else if err == service.ErrLastAdmin {
			response.Error(c, 400, "最后一个系统管理员不能删除")
		} else {
			response.Error(c, 500, err.Error())
		}
		return
	}

	response.Success(c, gin.H{
		"message": "删除成功",
	})
}

// BatchDeleteUsers 批量删除用户
// @Summary 批量删除用户
// @Description 批量软删除多个用户
// @Tags 用户管理
// @Security Bearer
// @Param request body service.BatchDeleteUsersRequest true "用户ID列表"
// @Success 200 {object} response.Response{data=object{success_count=int}}
// @Router /api/v1/users/batch-delete [post]
func (ctrl *UserController) BatchDeleteUsers(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	operatorID := middleware.GetUserID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req service.BatchDeleteUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误")
		return
	}

	if len(req.UserIDs) == 0 {
		response.Error(c, response.InvalidParamsCode, "用户ID列表不能为空")
		return
	}

	if len(req.UserIDs) > 100 {
		response.Error(c, response.InvalidParamsCode, "批量删除最多支持100个用户")
		return
	}

	successCount, err := ctrl.userService.BatchDeleteUsers(c.Request.Context(), tenantID, &req, operatorID)
	if err != nil {
		logger.Error(fmt.Sprintf("批量删除用户失败: tenant_id=%d, error=%v", tenantID, err))
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message":       "批量删除完成",
		"success_count": successCount,
	})
}

// UpdateUserStatus 批量更新用户状态
// @Summary 批量更新用户状态
// @Description 批量启用或禁用用户
// @Tags 用户管理
// @Security Bearer
// @Param request body object{user_ids=[]int64,status=string} true "用户ID列表和状态"
// @Success 200 {object} response.Response{data=object{affected_count=int}}
// @Router /api/v1/users/batch-update-status [post]
func (ctrl *UserController) UpdateUserStatus(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	operatorID := middleware.GetUserID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req struct {
		UserIDs []int64 `json:"user_ids" binding:"required"`
		Status  string  `json:"status" binding:"required,oneof=ACTIVE SUSPENDED"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误")
		return
	}

	if len(req.UserIDs) == 0 {
		response.Error(c, response.InvalidParamsCode, "用户ID列表不能为空")
		return
	}

	if len(req.UserIDs) > 100 {
		response.Error(c, response.InvalidParamsCode, "批量更新最多支持100个用户")
		return
	}

	affectedCount, err := ctrl.userService.UpdateUserStatus(c.Request.Context(), tenantID, req.UserIDs, req.Status, operatorID)
	if err != nil {
		logger.Error(fmt.Sprintf("批量更新用户状态失败: tenant_id=%d, error=%v", tenantID, err))
		if err == service.ErrLastAdmin {
			response.Error(c, 400, "最后一个系统管理员不能禁用")
		} else {
			response.Error(c, 500, err.Error())
		}
		return
	}

	response.Success(c, gin.H{
		"message":        "批量更新状态成功",
		"affected_count": affectedCount,
	})
}

// ResetUserPassword 重置用户密码
// @Summary 重置用户密码
// @Description 管理员重置指定用户的密码
// @Tags 用户管理
// @Security Bearer
// @Param id path int true "用户ID"
// @Param request body service.ResetUserPasswordRequest true "新密码"
// @Success 200 {object} response.Response
// @Router /api/v1/users/:id/reset-password [post]
func (ctrl *UserController) ResetUserPassword(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	operatorID := middleware.GetUserID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.InvalidParamsCode, "用户ID格式错误")
		return
	}

	var req service.ResetUserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误")
		return
	}

	err = ctrl.userService.ResetUserPassword(c.Request.Context(), tenantID, userID, &req, operatorID)
	if err != nil {
		logger.Error(fmt.Sprintf("重置用户密码失败: user_id=%d, tenant_id=%d, error=%v", userID, tenantID, err))
		if err == service.ErrUserNotFound {
			response.Error(c, 404, "用户不存在")
		} else {
			response.Error(c, 500, err.Error())
		}
		return
	}

	response.Success(c, gin.H{
		"message": "密码重置成功",
	})
}

// GetUserRoles 获取用户的所有角色
// @Summary 获取用户角色
// @Description 获取指定用户的所有角色
// @Tags 用户管理
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=[]models.Role}
// @Router /api/v1/users/:id/roles [get]
func (ctrl *UserController) GetUserRoles(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.InvalidParamsCode, "用户ID格式错误")
		return
	}

	roles, err := ctrl.userService.GetUserRoles(c.Request.Context(), tenantID, userID)
	if err != nil {
		logger.Error(fmt.Sprintf("获取用户角色失败: user_id=%d, tenant_id=%d, error=%v", userID, tenantID, err))
		if err == service.ErrUserNotFound {
			response.Error(c, 404, "用户不存在")
		} else {
			response.Error(c, 500, "获取用户角色失败")
		}
		return
	}

	response.Success(c, roles)
}

// AssignRoles 为用户分配角色
// @Summary 分配用户角色
// @Description 为指定用户分配一个或多个角色
// @Tags 用户管理
// @Security Bearer
// @Param id path int true "用户ID"
// @Param request body service.AssignRolesRequest true "角色ID列表"
// @Success 200 {object} response.Response
// @Router /api/v1/users/:id/roles [put]
func (ctrl *UserController) AssignRoles(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	operatorID := middleware.GetUserID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.InvalidParamsCode, "用户ID格式错误")
		return
	}

	var req service.AssignRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误")
		return
	}

	if len(req.RoleIDs) == 0 {
		response.Error(c, response.InvalidParamsCode, "角色ID列表不能为空")
		return
	}

	err = ctrl.userService.AssignRoles(c.Request.Context(), tenantID, userID, &req, operatorID)
	if err != nil {
		logger.Error(fmt.Sprintf("分配用户角色失败: user_id=%d, tenant_id=%d, error=%v", userID, tenantID, err))
		if err == service.ErrUserNotFound {
			response.Error(c, 404, "用户不存在")
		} else {
			response.Error(c, 500, err.Error())
		}
		return
	}

	response.Success(c, gin.H{
		"message": "角色分配成功",
	})
}

// GetCurrentUserRoles 获取当前用户的角色
// @Summary 获取当前用户角色
// @Description 获取登录用户的所有角色
// @Tags 用户管理
// @Security Bearer
// @Success 200 {object} response.Response{data=[]models.Role}
// @Router /api/v1/users/me/roles [get]
func (ctrl *UserController) GetCurrentUserRoles(c *gin.Context) {
	userID := middleware.GetUserID(c)
	tenantID := middleware.GetTenantID(c)
	if userID == 0 || tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	roles, err := ctrl.userService.GetUserRoles(c.Request.Context(), tenantID, userID)
	if err != nil {
		logger.Error(fmt.Sprintf("获取当前用户角色失败: user_id=%d, tenant_id=%d, error=%v", userID, tenantID, err))
		response.Error(c, 500, "获取用户角色失败")
		return
	}

	response.Success(c, roles)
}
