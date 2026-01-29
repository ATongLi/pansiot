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

type RoleController struct {
	roleService *service.RoleService
}

// NewRoleController 创建角色控制器
func NewRoleController(roleService *service.RoleService) *RoleController {
	return &RoleController{
		roleService: roleService,
	}
}

// ListRoles 获取角色列表
// @Summary 获取角色列表
// @Description 分页查询角色列表，支持关键词搜索和状态筛选
// @Tags 角色管理
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Param is_system query bool false "是否系统角色"
// @Param status query string false "角色状态"
// @Success 200 {object} response.Response{data=service.ListRolesResponse}
// @Router /api/v1/roles [get]
func (ctrl *RoleController) ListRoles(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req service.ListRolesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误")
		return
	}

	result, err := ctrl.roleService.ListRoles(c.Request.Context(), tenantID, &req)
	if err != nil {
		logger.Error(fmt.Sprintf("获取角色列表失败: tenant_id=%d, error=%v", tenantID, err))
		response.Error(c, 500, "获取角色列表失败")
		return
	}

	response.SuccessWithPagination(c, result.Total, result.Page, result.PageSize, result.Roles)
}

// CreateRole 创建角色
// @Summary 创建角色
// @Description 创建新角色并分配权限
// @Tags 角色管理
// @Security Bearer
// @Param request body service.CreateRoleRequest true "角色信息"
// @Success 200 {object} response.Response{data=models.Role}
// @Router /api/v1/roles [post]
func (ctrl *RoleController) CreateRole(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	operatorID := middleware.GetUserID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req service.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误: "+err.Error())
		return
	}

	role, err := ctrl.roleService.CreateRole(c.Request.Context(), tenantID, &req, operatorID)
	if err != nil {
		logger.Error(fmt.Sprintf("创建角色失败: tenant_id=%d, error=%v", tenantID, err))
		if err == service.ErrRoleAlreadyExists {
			response.Error(c, response.UserExistsCode, "角色已存在")
		} else {
			response.Error(c, 500, err.Error())
		}
		return
	}

	response.Success(c, role)
}

// GetRole 获取角色详情
// @Summary 获取角色详情
// @Description 获取指定角色的详细信息，包括权限列表
// @Tags 角色管理
// @Security Bearer
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response{data=models.Role}
// @Router /api/v1/roles/:id [get]
func (ctrl *RoleController) GetRole(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.InvalidParamsCode, "角色ID格式错误")
		return
	}

	role, err := ctrl.roleService.GetRole(c.Request.Context(), tenantID, roleID)
	if err != nil {
		logger.Error(fmt.Sprintf("获取角色详情失败: role_id=%d, tenant_id=%d, error=%v", roleID, tenantID, err))
		if err == service.ErrRoleNotFound {
			response.Error(c, 404, "角色不存在")
		} else {
			response.Error(c, 500, "获取角色详情失败")
		}
		return
	}

	response.Success(c, role)
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Description 更新角色信息和权限
// @Tags 角色管理
// @Security Bearer
// @Param id path int true "角色ID"
// @Param request body service.UpdateRoleRequest true "角色信息"
// @Success 200 {object} response.Response{data=models.Role}
// @Router /api/v1/roles/:id [put]
func (ctrl *RoleController) UpdateRole(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	operatorID := middleware.GetUserID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.InvalidParamsCode, "角色ID格式错误")
		return
	}

	var req service.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误: "+err.Error())
		return
	}

	role, err := ctrl.roleService.UpdateRole(c.Request.Context(), tenantID, roleID, &req, operatorID)
	if err != nil {
		logger.Error(fmt.Sprintf("更新角色失败: role_id=%d, tenant_id=%d, error=%v", roleID, tenantID, err))
		if err == service.ErrRoleNotFound {
			response.Error(c, 404, "角色不存在")
		} else if err == service.ErrCannotModifySystemRole {
			response.Error(c, 400, "系统预设角色不能修改")
		} else if err == service.ErrRoleAlreadyExists {
			response.Error(c, response.UserExistsCode, "角色名称已存在")
		} else {
			response.Error(c, 500, err.Error())
		}
		return
	}

	response.Success(c, role)
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 删除指定角色
// @Tags 角色管理
// @Security Bearer
// @Param id path int true "角色ID"
// @Param request body service.DeleteRoleRequest true "删除选项"
// @Success 200 {object} response.Response
// @Router /api/v1/roles/:id [delete]
func (ctrl *RoleController) DeleteRole(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	operatorID := middleware.GetUserID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.InvalidParamsCode, "角色ID格式错误")
		return
	}

	var req service.DeleteRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误: "+err.Error())
		return
	}

	err = ctrl.roleService.DeleteRole(c.Request.Context(), tenantID, roleID, &req, operatorID)
	if err != nil {
		logger.Error(fmt.Sprintf("删除角色失败: role_id=%d, tenant_id=%d, error=%v", roleID, tenantID, err))
		if err == service.ErrRoleNotFound {
			response.Error(c, 404, "角色不存在")
		} else if err == service.ErrCannotModifySystemRole {
			response.Error(c, 400, "系统预设角色不能删除")
		} else if err == service.ErrRoleInUse {
			response.Error(c, 400, "角色正在使用中，请先移除用户或选择强制删除")
		} else {
			response.Error(c, 500, err.Error())
		}
		return
	}

	response.Success(c, gin.H{
		"message": "删除成功",
	})
}

// GetRolePermissions 获取角色的所有权限
// @Summary 获取角色权限
// @Description 获取指定角色的所有权限
// @Tags 角色管理
// @Security Bearer
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response{data=[]models.Permission}
// @Router /api/v1/roles/:id/permissions [get]
func (ctrl *RoleController) GetRolePermissions(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.InvalidParamsCode, "角色ID格式错误")
		return
	}

	permissions, err := ctrl.roleService.GetRolePermissions(c.Request.Context(), tenantID, roleID)
	if err != nil {
		logger.Error(fmt.Sprintf("获取角色权限失败: role_id=%d, tenant_id=%d, error=%v", roleID, tenantID, err))
		if err == service.ErrRoleNotFound {
			response.Error(c, 404, "角色不存在")
		} else {
			response.Error(c, 500, "获取角色权限失败")
		}
		return
	}

	response.Success(c, permissions)
}

// AssignPermissionsToRole 为角色分配权限
// @Summary 分配角色权限
// @Description 为指定角色分配一个或多个权限
// @Tags 角色管理
// @Security Bearer
// @Param id path int true "角色ID"
// @Param request body object{permission_ids=[]int64} true "权限ID列表"
// @Success 200 {object} response.Response
// @Router /api/v1/roles/:id/permissions [put]
func (ctrl *RoleController) AssignPermissionsToRole(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	operatorID := middleware.GetUserID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.InvalidParamsCode, "角色ID格式错误")
		return
	}

	var req struct {
		PermissionIDs []int64 `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.InvalidParamsCode, "参数错误")
		return
	}

	if len(req.PermissionIDs) == 0 {
		response.Error(c, response.InvalidParamsCode, "权限ID列表不能为空")
		return
	}

	err = ctrl.roleService.AssignPermissionsToRole(c.Request.Context(), tenantID, roleID, req.PermissionIDs, operatorID)
	if err != nil {
		logger.Error(fmt.Sprintf("分配角色权限失败: role_id=%d, tenant_id=%d, error=%v", roleID, tenantID, err))
		if err == service.ErrRoleNotFound {
			response.Error(c, 404, "角色不存在")
		} else if err == service.ErrCannotModifySystemRole {
			response.Error(c, 400, "系统预设角色不能修改")
		} else {
			response.Error(c, 500, err.Error())
		}
		return
	}

	response.Success(c, gin.H{
		"message": "权限分配成功",
	})
}

// GetAllPermissions 获取所有可用权限
// @Summary 获取所有权限
// @Description 获取当前租户的所有可用权限列表
// @Tags 角色管理
// @Security Bearer
// @Success 200 {object} response.Response{data=[]models.Permission}
// @Router /api/v1/permissions [get]
func (ctrl *RoleController) GetAllPermissions(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	if tenantID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	permissions, err := ctrl.roleService.GetAllPermissions(c.Request.Context(), tenantID)
	if err != nil {
		logger.Error(fmt.Sprintf("获取权限列表失败: tenant_id=%d, error=%v", tenantID, err))
		response.Error(c, 500, "获取权限列表失败")
		return
	}

	response.Success(c, permissions)
}
