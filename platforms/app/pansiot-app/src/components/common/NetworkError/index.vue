<template>
  <view class="network-error">
    <uni-icons type="info" size="100" color="#dd524d"></uni-icons>
    <text class="error-text">{{ text }}</text>
    <slot name="action">
      <button class="error-action" type="primary" size="mini" @click="handleRetry">
        {{ retryText }}
      </button>
    </slot>
  </view>
</template>

<script setup lang="ts">
interface Props {
  text?: string;
  retryText?: string;
}

withDefaults(defineProps<Props>(), {
  text: '网络请求失败',
  retryText: '点击重试',
});

const emit = defineEmits<{
  retry: [];
}>();

const handleRetry = () => {
  emit('retry');
};
</script>

<style lang="scss" scoped>
.network-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100rpx 60rpx;
}

.error-text {
  font-size: 28rpx;
  color: #dd524d;
  margin: 40rpx 0;
}

.error-action {
  margin-top: 20rpx;
}
</style>
