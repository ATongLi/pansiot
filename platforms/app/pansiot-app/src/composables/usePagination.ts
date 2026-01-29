/**
 * 分页相关的组合式函数
 */

import { ref, computed } from 'vue';
import { DEFAULT_PAGE_SIZE } from '@/utils/constants';

export interface PaginationParams {
  page: number;
  pageSize: number;
}

/**
 * 使用分页
 */
export function usePagination(initialPageSize = DEFAULT_PAGE_SIZE) {
  const page = ref(1);
  const pageSize = ref(initialPageSize);
  const total = ref(0);

  /**
   * 分页参数
   */
  const paginationParams = computed<PaginationParams>(() => ({
    page: page.value,
    pageSize: pageSize.value,
  }));

  /**
   * 总页数
   */
  const totalPages = computed(() => {
    return Math.ceil(total.value / pageSize.value);
  });

  /**
   * 是否是第一页
   */
  const isFirstPage = computed(() => {
    return page.value === 1;
  });

  /**
   * 是否是最后一页
   */
  const isLastPage = computed(() => {
    return page.value >= totalPages.value;
  });

  /**
   * 设置总数
   */
  const setTotal = (value: number) => {
    total.value = value;
  };

  /**
   * 跳转到指定页
   */
  const goToPage = (targetPage: number) => {
    if (targetPage < 1 || targetPage > totalPages.value) {
      return;
    }
    page.value = targetPage;
  };

  /**
   * 上一页
   */
  const prevPage = () => {
    if (!isFirstPage.value) {
      page.value--;
    }
  };

  /**
   * 下一页
   */
  const nextPage = () => {
    if (!isLastPage.value) {
      page.value++;
    }
  };

  /**
   * 重置分页
   */
  const reset = () => {
    page.value = 1;
    total.value = 0;
  };

  /**
   * 分页信息
   */
  const paginationInfo = computed(() => {
    const start = (page.value - 1) * pageSize.value + 1;
    const end = Math.min(page.value * pageSize.value, total.value);
    return {
      start,
      end,
      total: total.value,
    };
  });

  return {
    page,
    pageSize,
    total,
    totalPages,
    isFirstPage,
    isLastPage,
    paginationParams,
    paginationInfo,
    setTotal,
    goToPage,
    prevPage,
    nextPage,
    reset,
  };
}
