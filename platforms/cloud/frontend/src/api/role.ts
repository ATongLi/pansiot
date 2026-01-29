import axios from 'axios';

export const roleApi = {
  listRoles: (params: any) => axios.get('/api/v1/roles', { params }),
  getRole: (id: number) => axios.get(`/api/v1/roles/${id}`),
  createRole: (data: any) => axios.post('/api/v1/roles', data),
  updateRole: (id: number, data: any) => axios.put(`/api/v1/roles/${id}`, data),
  deleteRole: (id: number) => axios.delete(`/api/v1/roles/${id}`),
  getRolePermissions: (id: number) => axios.get(`/api/v1/roles/${id}/permissions`),
  assignPermissionsToRole: (id: number, data: any) => axios.put(`/api/v1/roles/${id}/permissions`, data),
  getAllPermissions: () => axios.get('/api/v1/roles/permissions/all'),
};
