<template>
  <view class="custom-nav-bar" :style="{ height: navBarHeight + 'px', paddingTop: statusBarHeight + 'px' }">
    <view class="nav-bar-content">
      <!-- 返回按钮 -->
      <view v-if="showBack" class="nav-back" @click="handleBack">
        <uni-icons type="back" size="20" :color="color"></uni-icons>
      </view>

      <!-- 标题 -->
      <view class="nav-title" :style="{ color: color }">
        {{ title }}
      </view>

      <!-- 右侧操作按钮 -->
      <view v-if="showRight" class="nav-right" @click="handleRightClick">
        <slot name="right">
          <uni-icons v-if="rightIcon" :type="rightIcon" size="20" :color="color"></uni-icons>
        </slot>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  title?: string;
  color?: string;
  backgroundColor?: string;
  showBack?: boolean;
  showRight?: boolean;
  rightIcon?: string;
}

const props = withDefaults(defineProps<Props>(), {
  title: '',
  color: '#000000',
  backgroundColor: '#ffffff',
  showBack: true,
  showRight: false,
  rightIcon: '',
});

const emit = defineEmits<{
  back: [];
  rightClick: [];
}>();

// 状态栏高度
const statusBarHeight = computed(() => {
  const systemInfo = uni.getSystemInfoSync();
  return systemInfo.statusBarHeight || 20;
});

// 导航栏高度 (44px)
const navBarHeight = computed(() => {
  return statusBarHeight.value + 44;
});

// 返回
const handleBack = () => {
  emit('back');
  uni.navigateBack();
};

// 右侧按钮点击
const handleRightClick = () => {
  emit('rightClick');
};
</script>

<style lang="scss" scoped>
.custom-nav-bar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 999;
  background-color: v-bind(backgroundColor);

  .nav-bar-content {
    position: relative;
    height: 44px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .nav-back {
    position: absolute;
    left: 0;
    bottom: 0;
    padding: 0 20rpx;
    height: 44px;
    display: flex;
    align-items: center;
    z-index: 1;
  }

  .nav-title {
    font-size: 32rpx;
    font-weight: 500;
    max-width: 400rpx;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .nav-right {
    position: absolute;
    right: 0;
    bottom: 0;
    padding: 0 20rpx;
    height: 44px;
    display: flex;
    align-items: center;
    z-index: 1;
  }
}
</style>
