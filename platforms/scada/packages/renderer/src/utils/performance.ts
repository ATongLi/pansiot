/**
 * Performance Utilities
 *
 * 性能优化工具函数
 *
 * FE-006-29: 性能优化
 */

import { useCallback, useRef, useEffect, useState } from 'react';

/**
 * 防抖函数
 * @param func 要防抖的函数
 * @param wait 等待时间（毫秒）
 * @returns 防抖后的函数
 */
export function debounce<T extends (...args: any[]) => any>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let timeout: NodeJS.Timeout | null = null;

  return function executedFunction(...args: Parameters<T>) {
    const later = () => {
      timeout = null;
      func(...args);
    };

    if (timeout) {
      clearTimeout(timeout);
    }
    timeout = setTimeout(later, wait);
  };
}

/**
 * 节流函数
 * @param func 要节流的函数
 * @param limit 时间限制（毫秒）
 * @returns 节流后的函数
 */
export function throttle<T extends (...args: any[]) => any>(
  func: T,
  limit: number
): (...args: Parameters<T>) => void {
  let inThrottle: boolean = false;

  return function executedFunction(...args: Parameters<T>) {
    if (!inThrottle) {
      func(...args);
      inThrottle = true;
      setTimeout(() => (inThrottle = false), limit);
    }
  };
}

/**
 * 性能优化的 useCallback Hook
 * 使用防抖优化回调函数
 */
export function useDebouncedCallback<T extends (...args: any[]) => any>(
  callback: T,
  delay: number
): T {
  const debouncedRef = useRef<T>();

  useEffect(() => {
    debouncedRef.current = debounce(callback, delay) as T;
  }, [callback, delay]);

  return debouncedRef.current as T;
}

/**
 * 性能优化的 useCallback Hook
 * 使用节流优化回调函数
 */
export function useThrottledCallback<T extends (...args: any[]) => any>(
  callback: T,
  limit: number
): T {
  const throttledRef = useRef<T>();

  useEffect(() => {
    throttledRef.current = throttle(callback, limit) as T;
  }, [callback, limit]);

  return throttledRef.current as T;
}

/**
 * RAF (Request Animation Frame) 节流
 * 用于优化滚动、缩放等高频事件
 */
export function rafThrottle<T extends (...args: any[]) => any>(callback: T): T {
  let requestId: number | null = null;

  return function executedFunction(...args: Parameters<T>) {
    if (requestId !== null) {
      cancelAnimationFrame(requestId);
    }

    requestId = requestAnimationFrame(() => {
      callback(...args);
      requestId = null;
    });
  } as T;
}

/**
 * 批量更新 Hook
 * 将多个状态更新合并为单次渲染
 */
export function useBatchUpdate<T>(fn: () => T, deps: any[] = []): T {
  return useCallback(fn, deps);
}

/**
 * 性能监控工具
 */
export class PerformanceMonitor {
  private marks: Map<string, number> = new Map();
  private measures: Map<string, number[]> = new Map();

  /**
   * 开始标记
   */
  mark(name: string): void {
    this.marks.set(name, performance.now());
  }

  /**
   * 结束标记并测量时间
   */
  measure(name: string): number {
    const startTime = this.marks.get(name);
    if (!startTime) {
      console.warn(`PerformanceMonitor: mark "${name}" not found`);
      return 0;
    }

    const endTime = performance.now();
    const duration = endTime - startTime;

    const measures = this.measures.get(name) || [];
    measures.push(duration);
    this.measures.set(name, measures);

    console.log(`PerformanceMonitor: ${name} took ${duration.toFixed(2)}ms`);
    return duration;
  }

  /**
   * 获取平均测量时间
   */
  getAverage(name: string): number {
    const measures = this.measures.get(name);
    if (!measures || measures.length === 0) {
      return 0;
    }

    const sum = measures.reduce((a, b) => a + b, 0);
    return sum / measures.length;
  }

  /**
   * 获取测量统计
   */
  getStats(name: string): { min: number; max: number; avg: number; count: number } {
    const measures = this.measures.get(name);
    if (!measures || measures.length === 0) {
      return { min: 0, max: 0, avg: 0, count: 0 };
    }

    const min = Math.min(...measures);
    const max = Math.max(...measures);
    const avg = measures.reduce((a, b) => a + b, 0) / measures.length;
    const count = measures.length;

    return { min, max, avg, count };
  }

  /**
   * 清除标记和测量
   */
  clear(name?: string): void {
    if (name) {
      this.marks.delete(name);
      this.measures.delete(name);
    } else {
      this.marks.clear();
      this.measures.clear();
    }
  }
}

/**
 * 全局性能监控实例
 */
export const performanceMonitor = new PerformanceMonitor();

/**
 * 性能测量装饰器
 */
export function measurePerformance(name: string) {
  return function (
    target: any,
    propertyKey: string,
    descriptor: PropertyDescriptor
  ) {
    const originalMethod = descriptor.value;

    descriptor.value = function (...args: any[]) {
      performanceMonitor.mark(`${name}-start`);
      const result = originalMethod.apply(this, args);
      performanceMonitor.measure(`${name}-start`);
      return result;
    };

    return descriptor;
  };
}

/**
 * 虚拟滚动辅助函数
 * 计算可见区域的起始和结束索引
 */
export function getVisibleRange(
  scrollTop: number,
  containerHeight: number,
  itemHeight: number,
  totalItems: number
): { startIndex: number; endIndex: number } {
  const startIndex = Math.floor(scrollTop / itemHeight);
  const visibleItemCount = Math.ceil(containerHeight / itemHeight);
  const endIndex = Math.min(startIndex + visibleItemCount, totalItems - 1);

  return { startIndex: Math.max(0, startIndex), endIndex };
}

/**
 * 懒加载 Hook
 * 延迟加载组件或资源
 */
export function useLazyLoad<T>(
  loader: () => Promise<T>,
  deps: any[] = []
): { data: T | null; loading: boolean; error: Error | null; reload: () => void } {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await loader();
      setData(result);
    } catch (err) {
      setError(err as Error);
    } finally {
      setLoading(false);
    }
  }, deps);

  useEffect(() => {
    load();
  }, [load]);

  return { data, loading, error, reload: load };
}

/**
 * 内存使用监控
 */
export function getMemoryUsage(): {
  usedJSHeapSize: number;
  totalJSHeapSize: number;
  jsHeapSizeLimit: number;
} | null {
  if (typeof performance !== 'undefined' && (performance as any).memory) {
    const memory = (performance as any).memory;
    return {
      usedJSHeapSize: memory.usedJSHeapSize,
      totalJSHeapSize: memory.totalJSHeapSize,
      jsHeapSizeLimit: memory.jsHeapSizeLimit,
    };
  }
  return null;
}

/**
 * 批量更新 Hook (React 18+)
 */
export function useBatchedUpdates() {
  const [_, setState] = useState<object>();

  return useCallback((fn: () => void) => {
    setState((prev) => {
      fn();
      return prev;
    });
  }, []);
}
