import axios from 'axios';

export const tenantApi = {
  getCurrentTenant: () => axios.get('/api/v1/tenants/me'),
  updateCurrentTenant: (data: any) => axios.put('/api/v1/tenants/me', data),
  getTenantStats: () => axios.get('/api/v1/tenants/stats'),
  listSubTenants: (params: any) => axios.get('/api/v1/tenants/subs', { params }),
  getSubTenant: (id: number) => axios.get(`/api/v1/tenants/subs/${id}`),
  createSubTenant: (data: any) => axios.post('/api/v1/tenants/subs', data),
};
