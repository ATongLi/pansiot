// 角色和权限相关类型定义

export interface Role {
  id: number;
  tenant_id: number;
  name: string;
  code: string;
  description?: string;
  is_system_role: boolean;
  status: 'ACTIVE' | 'INACTIVE';
  created_at: string;
  updated_at: string;
}

export interface CreateRoleRequest {
  name: string;
  code: string;
  description?: string;
}

export interface UpdateRoleRequest {
  name?: string;
  description?: string;
  status?: 'ACTIVE' | 'INACTIVE';
}

export interface Permission {
  id: number;
  module_code: string;
  module_name: string;
  action_code: string;
  action_name: string;
  description?: string;
}

export interface RolePermission {
  role_id: number;
  permission_id: number;
}

export interface AssignPermissionsRequest {
  permission_ids: number[];
}

export interface FeatureModule {
  id: number;
  module_code: string;
  module_name: string;
  module_type: 'SYSTEM_DEFAULT' | 'OPTIONAL';
  description?: string;
  is_enabled: boolean;
}

export interface UserWithRoles {
  user_id: number;
  role_id: number;
  assigned_at: string;
  assigned_by: number;
}

export interface UserRole {
  id: number;
  user_id: number;
  role_id: number;
  role: Role;
  assigned_at: string;
  assigned_by: number;
}
