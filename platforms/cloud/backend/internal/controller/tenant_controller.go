package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"pansiot-cloud/internal/middleware"
	"pansiot-cloud/internal/service"
	"pansiot-cloud/pkg/logger"
	"pansiot-cloud/pkg/response"
)

type TenantController struct {
	tenantService *service.TenantService
}

// NewTenantController 创建租户控制器
func NewTenantController(tenantService *service.TenantService) *TenantController {
	return &TenantController{
		tenantService: tenantService,
	}
}

// GetCurrentTenant 获取当前租户信息
// @Summary 获取当前租户信息
// @Description 获取登录用户所属租户的详细信息
// @Tags 租户管理
// @Security Bearer
// @Success 200 {object} response.Response{data=models.Tenant}
// @Router /api/v1/tenants/me [get]
func (ctrl *TenantController) GetCurrentTenant(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	tenant, err := ctrl.tenantService.GetTenantByID(c.Request.Context(), tenantID)
	if err != nil {
		logger.Error(fmt.Sprintf("获取租户信息失败: tenant_id=%d, error=%v", tenantID, err))
		response.Error(c, 500, "获取租户信息失败")
		return
	}

	response.Success(c, tenant)
}

// UpdateCurrentTenant 更新当前租户信息
// @Summary 更新当前租户信息
// @Description 更新当前租户的基本信息
// @Tags 租户管理
// @Security Bearer
// @Param request body service.UpdateTenantRequest true "租户信息"
// @Success 200 {object} response.Response
// @Router /api/v1/tenants/me [put]
func (ctrl *TenantController) UpdateCurrentTenant(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	operatorID := middleware.GetUserID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req service.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误: "+err.Error())
		return
	}

	err := ctrl.tenantService.UpdateTenant(c.Request.Context(), tenantID, &req, operatorID)
	if err != nil {
		logger.Error(fmt.Sprintf("更新租户信息失败: tenant_id=%d, error=%v", tenantID, err))
		response.Error(c, 500, "更新租户信息失败")
		return
	}

	response.Success(c, gin.H{
		"message": "更新成功",
	})
}

// GetTenantStats 获取租户统计信息
// @Summary 获取租户统计信息
// @Description 获取当前租户的统计信息
// @Tags 租户管理
// @Security Bearer
// @Success 200 {object} response.Response{data=object}
// @Router /api/v1/tenants/stats [get]
func (ctrl *TenantController) GetTenantStats(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	stats, err := ctrl.tenantService.GetTenantStats(c.Request.Context(), tenantID)
	if err != nil {
		logger.Error(fmt.Sprintf("获取租户统计信息失败: tenant_id=%d, error=%v", tenantID, err))
		response.Error(c, 500, "获取统计信息失败")
		return
	}

	response.Success(c, stats)
}

// ListSubTenants 获取子租户列表（仅集成商）
// @Summary 获取子租户列表
// @Description 集成商获取其所有子租户列表
// @Tags 租户管理
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Param status query string false "租户状态"
// @Success 200 {object} response.Response{data=service.ListSubTenantsResponse}
// @Router /api/v1/tenants/subs [get]
func (ctrl *TenantController) ListSubTenants(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	tenantType := middleware.GetTenantType(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	// 验证是否是集成商
	if tenantType != "INTEGRATOR" {
		response.Error(c, response.PermissionDenyCode, "只有集成商可以查看子租户列表")
		return
	}

	var req service.ListSubTenantsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误")
		return
	}

	result, err := ctrl.tenantService.ListSubTenants(c.Request.Context(), tenantID, &req)
	if err != nil {
		logger.Error(fmt.Sprintf("获取子租户列表失败: tenant_id=%d, error=%v", tenantID, err))
		response.Error(c, 500, "获取子租户列表失败")
		return
	}

	response.SuccessWithPagination(c, result.Total, result.Page, result.PageSize, result.Tenants)
}

// CreateSubTenant 创建子租户（仅集成商）
// @Summary 创建子租户
// @Description 集成商创建新的子租户
// @Tags 租户管理
// @Security Bearer
// @Param request body service.CreateSubTenantRequest true "子租户信息"
// @Success 200 {object} response.Response{data=service.CreateSubTenantResponse}
// @Router /api/v1/tenants/subs [post]
func (ctrl *TenantController) CreateSubTenant(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	tenantType := middleware.GetTenantType(c)
	operatorID := middleware.GetUserID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	// 验证是否是集成商
	if tenantType != "INTEGRATOR" {
		response.Error(c, response.PermissionDenyCode, "只有集成商可以创建子租户")
		return
	}

	var req service.CreateSubTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误: "+err.Error())
		return
	}

	result, err := ctrl.tenantService.CreateSubTenant(c.Request.Context(), tenantID, &req, operatorID)
	if err != nil {
		logger.Error(fmt.Sprintf("创建子租户失败: tenant_id=%d, error=%v", tenantID, err))
		if err == service.ErrNotIntegrator {
			response.Error(c, response.PermissionDenyCode, "只有集成商可以创建子租户")
		} else {
			response.Error(c, 500, err.Error())
		}
		return
	}

	response.Success(c, result)
}

// GetSubTenant 获取子租户详情（仅集成商）
// @Summary 获取子租户详情
// @Description 集成商获取指定子租户的详细信息
// @Tags 租户管理
// @Security Bearer
// @Param id path int true "子租户ID"
// @Success 200 {object} response.Response{data=models.Tenant}
// @Router /api/v1/tenants/subs/:id [get]
func (ctrl *TenantController) GetSubTenant(c *gin.Context) {
	integratorID := middleware.GetTenantID(c)
	tenantType := middleware.GetTenantType(c)
	if integratorID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	// 验证是否是集成商
	if tenantType != "INTEGRATOR" {
		response.Error(c, response.PermissionDenyCode, "只有集成商可以查看子租户详情")
		return
	}

	subTenantID := c.Param("id")
	if subTenantID == "" {
		response.Error(c, response.InvalidParamsCode, "子租户ID不能为空")
		return
	}

	// 验证子租户是否属于当前集成商
	subTenant, err := ctrl.tenantService.GetTenantByID(c.Request.Context(), integratorID)
	if err != nil {
		response.Error(c, 404, "子租户不存在")
		return
	}

	if subTenant.ParentTenantID == nil || *subTenant.ParentTenantID != integratorID {
		response.Error(c, response.PermissionDenyCode, "无权访问此子租户")
		return
	}

	response.Success(c, subTenant)
}
