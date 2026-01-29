import React, { useState, useEffect } from 'react';
import { Table, Button, Space, Modal, Form, Input, Tag, message, Popconfirm, Card, Transfer } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { observer } from 'mobx-react-lite';
import { roleApi } from '../../api/role';
import type { Role, Permission } from '../../types/role';

const RoleManagement: React.FC = observer(() => {
  const [roles, setRoles] = useState<Role[]>([]);
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingRole, setEditingRole] = useState<Role | null>(null);
  const [selectedPermissions, setSelectedPermissions] = useState<string[]>([]);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchRoles();
    fetchPermissions();
  }, []);

  const fetchRoles = async () => {
    setLoading(true);
    try {
      const response = await roleApi.listRoles({ page: 1, page_size: 100 });
      setRoles(response.data.roles || []);
    } catch (error) {
      message.error('获取角色列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchPermissions = async () => {
    try {
      const response = await roleApi.getAllPermissions();
      setPermissions(response.data || []);
    } catch (error) {
      console.error('获取权限列表失败', error);
    }
  };

  const handleCreate = () => {
    setEditingRole(null);
    setSelectedPermissions([]);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = async (role: Role) => {
    try {
      const response = await roleApi.getRolePermissions(role.id);
      const permIds = (response.data || []).map((p: Permission) => p.id.toString());
      setSelectedPermissions(permIds);
      setEditingRole(role);
      form.setFieldsValue(role);
      setModalVisible(true);
    } catch (error) {
      message.error('获取角色权限失败');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await roleApi.deleteRole(id);
      message.success('删除成功');
      fetchRoles();
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      const permissionIds = selectedPermissions.map(id => parseInt(id));
      
      if (editingRole) {
        await roleApi.updateRole(editingRole.id, {
          ...values,
          permission_ids: permissionIds,
        });
        message.success('更新成功');
      } else {
        await roleApi.createRole({
          ...values,
          permission_ids: permissionIds,
        });
        message.success('创建成功');
      }
      setModalVisible(false);
      fetchRoles();
    } catch (error) {
      message.error(editingRole ? '更新失败' : '创建失败');
    }
  };

  const handlePermissionChange = (targetKeys: React.Key[]) => {
    setSelectedPermissions(targetKeys as string[]);
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '角色名称',
      dataIndex: 'role_name',
      key: 'role_name',
    },
    {
      title: '角色代码',
      dataIndex: 'role_code',
      key: 'role_code',
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '系统角色',
      dataIndex: 'is_system_role',
      key: 'is_system_role',
      render: (isSystem: boolean) => (
        <Tag color={isSystem ? 'blue' : 'default'}>
          {isSystem ? '是' : '否'}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleString('zh-CN'),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Role) => (
        <Space>
          <Button 
            type="link" 
            icon={<EditOutlined />} 
            onClick={() => handleEdit(record)}
            disabled={record.is_system_role}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个角色吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
            disabled={record.is_system_role}
          >
            <Button 
              type="link" 
              danger 
              icon={<DeleteOutlined />}
              disabled={record.is_system_role}
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const permissionDataSource = permissions.map(perm => ({
    key: perm.id.toString(),
    title: `${perm.module_code} - ${perm.action_code}`,
    description: perm.id.toString(),
  }));

  return (
    <Card title="角色管理">
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新建角色
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={roles}
        loading={loading}
        rowKey="id"
        pagination={false}
      />
      <Modal
        title={editingRole ? '编辑角色' : '新建角色'}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={() => setModalVisible(false)}
        width={800}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            label="角色名称"
            name="role_name"
            rules={[{ required: true, message: '请输入角色名称' }]}
          >
            <Input placeholder="请输入角色名称" />
          </Form.Item>
          <Form.Item
            label="角色代码"
            name="role_code"
            rules={[{ required: true, message: '请输入角色代码' }]}
          >
            <Input placeholder="请输入角色代码（英文）" />
          </Form.Item>
          <Form.Item
            label="描述"
            name="description"
          >
            <Input.TextArea placeholder="请输入角色描述" rows={3} />
          </Form.Item>
          <Form.Item label="权限分配">
            <Transfer
              dataSource={permissionDataSource}
              titles={['可用权限', '已选权限']}
              targetKeys={selectedPermissions}
              onChange={handlePermissionChange}
              render={item => item.title}
              listStyle={{
                width: 300,
                height: 400,
              }}
              showSearch
              filterOption={(inputValue, item) =>
                item.title?.indexOf(inputValue) !== -1
              }
            />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
});

export default RoleManagement;
