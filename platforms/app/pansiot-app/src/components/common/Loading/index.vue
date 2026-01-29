<template>
  <view v-if="show" class="loading-container">
    <view class="loading-mask" v-if="mask"></view>
    <view class="loading-content">
      <view class="loading-spinner"></view>
      <text v-if="text" class="loading-text">{{ text }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
interface Props {
  show?: boolean;
  text?: string;
  mask?: boolean;
}

withDefaults(defineProps<Props>(), {
  show: false,
  text: '加载中...',
  mask: true,
});
</script>

<style lang="scss" scoped>
.loading-container {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-mask {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.3);
}

.loading-content {
  position: relative;
  background-color: rgba(0, 0, 0, 0.7);
  border-radius: 16rpx;
  padding: 40rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.loading-spinner {
  width: 60rpx;
  height: 60rpx;
  border: 4rpx solid rgba(255, 255, 255, 0.3);
  border-top-color: #ffffff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.loading-text {
  margin-top: 20rpx;
  font-size: 28rpx;
  color: #ffffff;
}
</style>
