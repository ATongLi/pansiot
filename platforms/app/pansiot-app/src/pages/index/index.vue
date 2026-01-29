<template>
  <view class="splash-page">
    <view class="logo-container">
      <image class="logo" src="/static/images/logo.png" mode="aspectFit"></image>
      <text class="app-name">PansIot</text>
      <text class="app-slogan">工业物联网移动平台</text>
    </view>
    <view class="version-info">
      <text class="version">版本 1.0.0</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app';

onLoad(() => {
  console.log('Splash Page Loaded');

  // 2秒后自动跳转到登录页
  setTimeout(() => {
    uni.redirectTo({
      url: '/pages/auth/login/index',
    });
  }, 2000);
});
</script>

<style lang="scss" scoped>
.splash-page {
  width: 100vw;
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.logo-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.logo {
  width: 200rpx;
  height: 200rpx;
  margin-bottom: 40rpx;
}

.app-name {
  font-size: 64rpx;
  font-weight: bold;
  color: #ffffff;
  margin-bottom: 20rpx;
}

.app-slogan {
  font-size: 32rpx;
  color: rgba(255, 255, 255, 0.8);
}

.version-info {
  position: absolute;
  bottom: 100rpx;
}

.version {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.6);
}
</style>
