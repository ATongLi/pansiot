<template>
  <scroll-view
    class="pull-refresh"
    scroll-y
    refresher-enabled
    :refresher-triggered="refreshing"
    :refresher-threshold="100"
    @refresherrefresh="handleRefresh"
    @refresherrestore="handleRestore"
  >
    <slot></slot>
  </scroll-view>
</template>

<script setup lang="ts">
import { ref } from 'vue';

interface Props {
  disabled?: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  refresh: [];
}>();

const refreshing = ref(false);

/**
 * 触发刷新
 */
const handleRefresh = () => {
  if (props.disabled) return;

  refreshing.value = true;
  emit('refresh');
};

/**
 * 刷新完成
 */
const handleRestore = () => {
  refreshing.value = false;
};

/**
 * 停止刷新
 */
const stopRefresh = () => {
  refreshing.value = false;
};

defineExpose({
  stopRefresh,
});
</script>

<style lang="scss" scoped>
.pull-refresh {
  width: 100%;
  height: 100%;
}
</style>
