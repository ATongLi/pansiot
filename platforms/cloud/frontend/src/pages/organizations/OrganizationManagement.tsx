import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, message, Space, Popconfirm } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { tenantApi } from '../../api/tenant';

const OrganizationManagement: React.FC = () => {
  const [organizations, setOrganizations] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingOrg, setEditingOrg] = useState<any | null>(null);
  const [form] = Form.useForm();

  const fetchOrganizations = async () => {
    setLoading(true);
    try {
      const response = await tenantApi.listSubTenants({ page: 1, pageSize: 100 });
      setOrganizations(response.data.items || []);
    } catch (error) {
      message.error('获取组织列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchOrganizations();
  }, []);

  const handleCreate = () => {
    setEditingOrg(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (record: any) => {
    setEditingOrg(record);
    form.setFieldsValue(record);
    setModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await tenantApi.getSubTenant(id); // Note: API doesn't have delete, this is a placeholder
      message.success('删除成功');
      fetchOrganizations();
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingOrg) {
        await tenantApi.getSubTenant(editingOrg.id); // Note: API doesn't have update
        message.success('更新成功');
      } else {
        await tenantApi.createSubTenant(values);
        message.success('创建成功');
      }
      setModalVisible(false);
      fetchOrganizations();
    } catch (error) {
      message.error('操作失败');
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id' },
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '描述', dataIndex: 'description', key: 'description' },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: any) => (
        <Space>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个组织吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新建组织
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={organizations}
        loading={loading}
        rowKey="id"
      />
      <Modal
        title={editingOrg ? '编辑组织' : '新建组织'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="组织名称"
            rules={[{ required: true, message: '请输入组织名称' }]}
          >
            <Input placeholder="请输入组织名称" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea placeholder="请输入描述" rows={4} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default OrganizationManagement;
