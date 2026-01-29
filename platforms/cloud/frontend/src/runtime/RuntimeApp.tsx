import React from 'react';
import { Result, Button } from 'antd';
import { useNavigate } from 'react-router-dom';

const RuntimeApp: React.FC = () => {
  const navigate = useNavigate();

  return (
    <div style={{ padding: '50px', textAlign: 'center' }}>
      <Result
        status="info"
        title="Runtime运行环境"
        subTitle="此功能正在开发中，敬请期待..."
        extra={
          <Button type="primary" onClick={() => navigate('/dashboard')}>
            返回仪表板
          </Button>
        }
      />
    </div>
  );
};

export default RuntimeApp;
