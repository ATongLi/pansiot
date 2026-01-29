import React from 'react';
import { Result, Button } from 'antd';
import { useNavigate } from 'react-router-dom';

const EditorApp: React.FC = () => {
  const navigate = useNavigate();

  return (
    <div style={{ padding: '50px', textAlign: 'center' }}>
      <Result
        status="info"
        title="Web工程编辑器"
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

export default EditorApp;
