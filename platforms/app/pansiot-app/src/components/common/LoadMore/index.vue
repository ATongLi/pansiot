<template>
  <view class="load-more">
    <!-- 加载中 -->
    <view v-if="status === 'loading'" class="load-more-loading">
      <view class="load-more-spinner"></view>
      <text class="load-more-text">{{ loadingText }}</text>
    </view>

    <!-- 加载完成 -->
    <view v-else-if="status === 'success'" class="load-more-success">
      <text class="load-more-text">{{ successText }}</text>
    </view>

    <!-- 加载失败 -->
    <view v-else-if="status === 'error'" class="load-more-error" @click="handleRetry">
      <text class="load-more-text">{{ errorText }}</text>
    </view>

    <!-- 没有更多 -->
    <view v-else-if="status === 'nomore'" class="load-more-nomore">
      <text class="load-more-text">{{ nomoreText }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
interface Props {
  status?: 'loading' | 'success' | 'error' | 'nomore';
  loadingText?: string;
  successText?: string;
  errorText?: string;
  nomoreText?: string;
}

withDefaults(defineProps<Props>(), {
  status: 'loading',
  loadingText: '加载中...',
  successText: '加载成功',
  errorText: '加载失败,点击重试',
  nomoreText: '没有更多了',
});

const emit = defineEmits<{
  retry: [];
}>();

const handleRetry = () => {
  emit('retry');
};
</script>

<style lang="scss" scoped>
.load-more {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40rpx 0;
  min-height: 80rpx;
}

.load-more-loading,
.load-more-success,
.load-more-error,
.load-more-nomore {
  display: flex;
  align-items: center;
  justify-content: center;
}

.load-more-spinner {
  width: 40rpx;
  height: 40rpx;
  border: 3rpx solid rgba(0, 0, 0, 0.1);
  border-top-color: #007aff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-right: 16rpx;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.load-more-text {
  font-size: 24rpx;
  color: #999999;
}

.load-more-error {
  cursor: pointer;

  .load-more-text {
    color: #dd524d;
  }
}
</style>
