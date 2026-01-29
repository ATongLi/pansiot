<template>
  <view class="page-container" :style="{ backgroundColor: bgColor }">
    <view class="container-content" :style="{ padding: padding }">
      <slot></slot>
    </view>
  </view>
</template>

<script setup lang="ts">
interface Props {
  bgColor?: string;
  padding?: string;
  minHeight?: string;
}

withDefaults(defineProps<Props>(), {
  bgColor: '#f5f5f5',
  padding: '30rpx',
  minHeight: '100vh',
});
</script>

<style lang="scss" scoped>
.page-container {
  width: 100%;
  min-height: v-bind(minHeight);
  box-sizing: border-box;

  .container-content {
    width: 100%;
    box-sizing: border-box;
  }
}
</style>
