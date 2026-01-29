<template>
  <view class="tabbar-page">
    <view class="header">
      <text class="title">消息</text>
    </view>
    <view class="content">
      <view class="empty-state">
        <text class="empty-icon">💬</text>
        <text class="empty-text">消息中心功能</text>
        <text class="empty-hint">敬请期待...</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app';

onLoad(() => {
  console.log('Message Page Loaded');
});
</script>

<style lang="scss" scoped>
.tabbar-page {
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

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 200rpx 0;
}

.empty-icon {
  font-size: 120rpx;
  margin-bottom: 30rpx;
}

.empty-text {
  font-size: 32rpx;
  color: #333333;
  margin-bottom: 20rpx;
}

.empty-hint {
  font-size: 28rpx;
  color: #999999;
}
</style>
