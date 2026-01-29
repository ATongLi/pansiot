// 租户相关类型定义

export interface Tenant {
  id: number;
  serial_number: string;
  name: string;
  tenant_type: 'INTEGRATOR' | 'TERMINAL';
  industry: string;
  contact_person?: string;
  contact_phone?: string;
  contact_email?: string;
  parent_tenant_id?: number;
  status: 'ACTIVE' | 'SUSPENDED' | 'DELETED';
  expire_date?: string;
  max_sub_tenants: number;
  max_users: number;
  max_devices: number;
  max_storage_gb: number;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

export interface CreateTenantRequest {
  name: string;
  tenant_type: 'INTEGRATOR' | 'TERMINAL';
  industry: string;
  contact_person?: string;
  contact_phone?: string;
  contact_email?: string;
  parent_tenant_id?: number;
  max_sub_tenants?: number;
  max_users?: number;
  max_devices?: number;
  max_storage_gb?: number;
}

export interface UpdateTenantRequest {
  name?: string;
  industry?: string;
  contact_person?: string;
  contact_phone?: string;
  contact_email?: string;
}

export interface UpgradeToIntegratorRequest {
  max_sub_tenants: number;
  max_users: number;
  max_devices: number;
  max_storage_gb: number;
}

export interface TenantStats {
  total_users: number;
  active_users: number;
  total_sub_tenants: number;
  total_devices: number;
  used_storage_gb: number;
}
