<template>
  <view class="empty-state">
    <image v-if="image" class="empty-image" :src="image" mode="aspectFit"></image>
    <uni-icons v-else type="inbox" size="100" color="#c0c0c0"></uni-icons>
    <text class="empty-text">{{ text }}</text>
    <slot name="action">
      <button v-if="showAction" class="empty-action" type="primary" size="mini" @click="handleAction">
        {{ actionText }}
      </button>
    </slot>
  </view>
</template>

<script setup lang="ts">
interface Props {
  image?: string;
  text?: string;
  showAction?: boolean;
  actionText?: string;
}

const props = withDefaults(defineProps<Props>(), {
  image: '',
  text: '暂无数据',
  showAction: false,
  actionText: '重新加载',
});

const emit = defineEmits<{
  action: [];
}>();

const handleAction = () => {
  emit('action');
};
</script>

<style lang="scss" scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100rpx 60rpx;
}

.empty-image {
  width: 200rpx;
  height: 200rpx;
  margin-bottom: 40rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #999999;
  margin-bottom: 40rpx;
}

.empty-action {
  margin-top: 20rpx;
}
</style>
