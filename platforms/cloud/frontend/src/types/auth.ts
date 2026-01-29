// 认证相关类型定义

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  tenantName: string;
  industry: string;
  username: string;
  email: string;
  phone?: string;
  phoneCountryCode?: string;
  password: string;
  verificationCode: string;
}

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  user: User;
  tenant: Tenant;
}

export interface User {
  id: number;
  tenant_id: number;
  username: string;
  email: string;
  phone?: string;
  real_name?: string;
  avatar?: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface Tenant {
  id: number;
  serial_number: string;
  name: string;
  tenant_type: 'INTEGRATOR' | 'TERMINAL';
  industry: string;
  parent_tenant_id?: number;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

export interface ResetPasswordRequest {
  email: string;
  verification_code: string;
  new_password: string;
}

export interface SendVerificationCodeRequest {
  email?: string;
  phone?: string;
  phone_country_code?: string;
}
