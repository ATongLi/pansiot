/**
 * 请求相关的组合式函数
 */

import { ref, computed } from 'vue';

export interface RequestResult<T> {
  data: ref<T | null>;
  loading: ref<boolean>;
  error: ref<any>;
  execute: () => Promise<void>;
  refresh: () => Promise<void>;
}

/**
 * 使用请求
 */
export function useRequest<T>(apiFunc: () => Promise<T>): RequestResult<T> {
  const data = ref<T | null>(null);
  const loading = ref(false);
  const error = ref<any>(null);

  /**
   * 执行请求
   */
  const execute = async () => {
    loading.value = true;
    error.value = null;

    try {
      const result = await apiFunc();
      data.value = result as T;
    } catch (err) {
      error.value = err;
      console.error('[useRequest Error]', err);
    } finally {
      loading.value = false;
    }
  };

  /**
   * 刷新请求
   */
  const refresh = async () => {
    await execute();
  };

  return {
    data,
    loading,
    error,
    execute,
    refresh,
  };
}

/**
 * 使用加载状态
 */
export function useLoading(initialValue = false) {
  const loading = ref(initialValue);

  const setLoading = (value: boolean) => {
    loading.value = value;
  };

  const startLoading = () => {
    loading.value = true;
  };

  const stopLoading = () => {
    loading.value = false;
  };

  return {
    loading,
    setLoading,
    startLoading,
    stopLoading,
  };
}
