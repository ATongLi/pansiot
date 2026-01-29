<template>
  <view class="workspace-page">
    <view class="header">
      <text class="title">工作台</text>
    </view>
    <view class="content">
      <view class="welcome-card card">
        <text class="welcome-text">欢迎使用 PansIot 移动端</text>
        <text class="sub-text">工业物联网平台</text>
      </view>

      <view class="quick-actions">
        <view class="action-item card">
          <text class="action-icon">📊</text>
          <text class="action-text">数据统计</text>
        </view>
        <view class="action-item card">
          <text class="action-icon">🔧</text>
          <text class="action-text">设备管理</text>
        </view>
        <view class="action-item card">
          <text class="action-icon">⚠️</text>
          <text class="action-text">告警中心</text>
        </view>
        <view class="action-item card">
          <text class="action-icon">📝</text>
          <text class="action-text">工程管理</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app';

onLoad(() => {
  console.log('Workspace Page Loaded');
});
</script>

<style lang="scss" scoped>
.workspace-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 100rpx 30rpx 50rpx;
}

.title {
  font-size: 48rpx;
  font-weight: bold;
  color: #ffffff;
}

.content {
  padding: 30rpx;
}

.welcome-card {
  padding: 60rpx 30rpx;
  text-align: center;
  margin-bottom: 30rpx;
}

.welcome-text {
  display: block;
  font-size: 40rpx;
  font-weight: bold;
  color: #333333;
  margin-bottom: 20rpx;
}

.sub-text {
  display: block;
  font-size: 28rpx;
  color: #999999;
}

.quick-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 20rpx;
}

.action-item {
  width: calc(50% - 10rpx);
  padding: 40rpx 20rpx;
  text-align: center;
}

.action-icon {
  display: block;
  font-size: 64rpx;
  margin-bottom: 20rpx;
}

.action-text {
  display: block;
  font-size: 28rpx;
  color: #333333;
}
</style>
