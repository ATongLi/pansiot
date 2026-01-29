import axios from 'axios';

export const userApi = {
  listUsers: (params: any) => axios.get('/api/v1/users', { params }),
  getUser: (id: number) => axios.get(`/api/v1/users/${id}`),
  createUser: (data: any) => axios.post('/api/v1/users', data),
  updateUser: (id: number, data: any) => axios.put(`/api/v1/users/${id}`, data),
  deleteUser: (id: number) => axios.delete(`/api/v1/users/${id}`),
  resetUserPassword: (id: number, data: any) => axios.post(`/api/v1/users/${id}/reset-password`, data),
  getCurrentUserRoles: () => axios.get('/api/v1/users/me/roles'),
  getUserRoles: (id: number) => axios.get(`/api/v1/users/${id}/roles`),
  assignRoles: (id: number, data: any) => axios.put(`/api/v1/users/${id}/roles`, data),
};
