/**
 * 租户状态管理
 */
import { defineStore } from 'pinia';
import type { Tenant } from '@/types/global';

interface TenantState {
  currentTenant: Tenant | null;
  tenantList: Tenant[];
}

export const useTenantStore = defineStore('tenant', {
  state: (): TenantState => ({
    currentTenant: null,
    tenantList: [],
  }),

  actions: {
    /**
     * 切换租户
     */
    switchTenant(tenantId: string) {
      const tenant = this.tenantList.find((t) => t.id === tenantId);
      if (tenant) {
        this.currentTenant = tenant;
        uni.setStorageSync('currentTenant', tenant);
      }
    },

    /**
     * 设置租户列表
     */
    setTenantList(list: Tenant[]) {
      this.tenantList = list;
    },

    /**
     * 从本地存储恢复租户信息
     */
    restoreTenant() {
      const tenant = uni.getStorageSync('currentTenant');
      if (tenant) {
        this.currentTenant = tenant;
      }
    },
  },

  // 持久化配置
  persist: {
    key: 'tenant-store',
    storage: {
      getItem: (key) => uni.getStorageSync(key),
      setItem: (key, value) => uni.setStorageSync(key, value),
    },
  },
});
