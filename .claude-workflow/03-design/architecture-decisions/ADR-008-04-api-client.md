# ADR-008-04: 封装 uni.request 作为 HTTP 客户端

## Status
Accepted

## Date
2026-01-28

## Context
移动端应用需要与云平台 API 进行通信,需要选择合适的 HTTP 客户端。同时需要处理请求拦截、响应拦截、错误处理、Token 管理等功能。

## Decision
**封装 uni.request** 作为 HTTP 客户端,而不是使用 Axios。

### Rationale
1. **UniApp 内置**: uni.request 是 UniApp 内置 API,无需额外依赖
2. **多端兼容**: 官方保证多端兼容性,适配各平台差异
3. **轻量级**: 无需引入第三方库,减少包体积
4. **符合最佳实践**: UniApp 官方推荐使用 uni.request
5. **可封装**: 可以封装成类似 Axios 的 API,使用体验一致

### Alternatives Considered

#### 1. Axios
- **优势**:
  - 功能丰富,API 友好
  - 拦截器、请求取消等功能完善
  - 浏览器端生态成熟
- **劣势**:
  - 在小程序端有兼容性问题
  - 需要适配器才能工作
  - 增加包体积(~40KB)
  - 不是 UniApp 官方推荐
- **结论**: 不选择

#### 2. Fly.js
- **优势**:
  - 支持多端
  - API 类似 Axios
- **劣势**:
  - 社区不如 uni.request 成熟
  - 维护不活跃
  - 增加依赖
- **结论**: 不选择

#### 3. 原生 uni.request
- **优势**:
  - 官方内置,无需依赖
  - 多端兼容性有保证
- **劣势**:
  - API 不够友好
  - 需要手动处理拦截器
- **结论**: 选择,但需要封装

## Consequences

### Positive
1. **多端兼容**: 官方保证多端兼容性
2. **包体积小**: 无额外依赖,减少包体积
3. **性能好**: 轻量级,性能优秀
4. **可维护**: 封装层统一管理,易于维护

### Negative
1. **需要封装**: 需要自行封装拦截器等功能
2. **功能限制**: 某些高级功能需要自己实现

### Mitigation Strategies
1. **封装拦截器**:
   - 实现请求拦截器(添加 Token、设置 Header)
   - 实现响应拦截器(统一错误处理、Token 刷新)

2. **实现高级功能**:
   - 请求取消
   - 请求重试
   - 请求缓存
   - 并发控制

## Implementation

### HTTP 客户端封装

#### 1. 基础封装 (api/client/request.ts)
```typescript
import { interceptor } from './interceptor';
import type { RequestConfig, ApiResponse } from '../types/api.types';

class Request {
  private config: RequestConfig;

  constructor(config: RequestConfig) {
    this.config = config;
  }

  /**
   * 通用请求方法
   */
  request<T = any>(config: RequestConfig): Promise<ApiResponse<T>> {
    // 请求前拦截
    const requestConfig = interceptor.request(config);

    return new Promise((resolve, reject) => {
      uni.request({
        url: this.getFullUrl(requestConfig.url),
        method: requestConfig.method || 'GET',
        data: requestConfig.data,
        header: requestConfig.header,
        timeout: requestConfig.timeout || this.config.timeout,
        success: (res) => {
          // 响应拦截
          const response = interceptor.response(res, requestConfig);
          resolve(response);
        },
        fail: (err) => {
          // 错误拦截
          const error = interceptor.error(err, requestConfig);
          reject(error);
        },
      });
    });
  }

  /**
   * GET 请求
   */
  get<T = any>(url: string, data?: any, config?: RequestConfig): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: 'GET', data, ...config });
  }

  /**
   * POST 请求
   */
  post<T = any>(url: string, data?: any, config?: RequestConfig): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: 'POST', data, ...config });
  }

  /**
   * PUT 请求
   */
  put<T = any>(url: string, data?: any, config?: RequestConfig): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: 'PUT', data, ...config });
  }

  /**
   * DELETE 请求
   */
  delete<T = any>(url: string, data?: any, config?: RequestConfig): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: 'DELETE', data, ...config });
  }

  /**
   * 获取完整 URL
   */
  private getFullUrl(url: string): string {
    if (url.startsWith('http')) {
      return url;
    }
    return this.config.baseURL + url;
  }
}

export const request = new Request({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 30000,
});
```

#### 2. 拦截器 (api/client/interceptors.ts)
```typescript
import type { RequestConfig, ApiResponse } from '../types/api.types';
import { useUserStore } from '@/stores/user.store';
import { useTenantStore } from '@/stores/tenant.store';

class Interceptor {
  /**
   * 请求拦截器
   */
  request(config: RequestConfig): RequestConfig {
    // 添加 Token
    const userStore = useUserStore();
    if (userStore.token) {
      config.header = {
        ...config.header,
        'Authorization': `Bearer ${userStore.token}`,
      };
    }

    // 添加租户 ID
    const tenantStore = useTenantStore();
    if (tenantStore.currentTenant?.id) {
      config.header = {
        ...config.header,
        'X-Tenant-ID': tenantStore.currentTenant.id,
      };
    }

    // 添加通用 Header
    config.header = {
      'Content-Type': 'application/json',
      ...config.header,
    };

    // 请求时间戳(防止缓存)
    if (config.method === 'GET') {
      config.data = {
        ...config.data,
        _t: Date.now(),
      };
    }

    console.log(`[Request] ${config.method} ${config.url}`, config.data);
    return config;
  }

  /**
   * 响应拦截器
   */
  response<T>(res: any, config: RequestConfig): ApiResponse<T> {
    console.log(`[Response] ${config.url}`, res);

    const { statusCode, data } = res;

    // HTTP 状态码检查
    if (statusCode !== 200) {
      return {
        code: statusCode,
        data: null as any,
        message: '请求失败',
      };
    }

    // 业务状态码检查
    if (data.code !== 200) {
      // Token 过期,自动刷新
      if (data.code === 401) {
        this.handleTokenExpired();
      }

      return {
        code: data.code,
        data: null as any,
        message: data.message || '请求失败',
      };
    }

    return {
      code: data.code,
      data: data.data,
      message: data.message,
    };
  }

  /**
   * 错误拦截器
   */
  error(err: any, config: RequestConfig): never {
    console.error(`[Error] ${config.url}`, err);

    let message = '网络错误';

    // 网络超时
    if (err.errMsg?.includes('timeout')) {
      message = '请求超时';
    }
    // 网络异常
    else if (err.errMsg?.includes('fail')) {
      message = '网络连接失败';
    }

    uni.showToast({
      title: message,
      icon: 'none',
    });

    throw new Error(message);
  }

  /**
   * 处理 Token 过期
   */
  private async handleTokenExpired() {
    const userStore = useUserStore();

    try {
      // 刷新 Token
      const res = await uni.request({
        url: `${import.meta.env.VITE_API_BASE_URL}/api/auth/refresh`,
        method: 'POST',
        data: {
          refreshToken: userStore.refreshToken,
        },
      });

      if (res.data.code === 200) {
        userStore.token = res.data.data.token;
        // 重试原请求
      } else {
        // 刷新失败,跳转登录
        userStore.logout();
        uni.navigateTo({ url: '/pages/login' });
      }
    } catch (error) {
      userStore.logout();
      uni.navigateTo({ url: '/pages/login' });
    }
  }
}

export const interceptor = new Interceptor();
```

#### 3. 配置 (api/client/config.ts)
```typescript
import type { RequestConfig } from '../types/api.types';

export const requestConfig: RequestConfig = {
  baseURL: import.meta.env.VITE_API_BASE_URL || 'https://api.example.com',
  timeout: 30000,
  header: {
    'Content-Type': 'application/json',
  },
};
```

### API 模块封装

#### 认证 API (api/modules/auth.api.ts)
```typescript
import { request } from '../client/request';
import type { LoginParams, LoginResult, RegisterParams } from '../types/api.types';

export const authApi = {
  /**
   * 用户登录
   */
  login: (params: LoginParams) =>
    request.post<LoginResult>('/api/auth/login', params),

  /**
   * 用户注册
   */
  register: (params: RegisterParams) =>
    request.post<LoginResult>('/api/auth/register', params),

  /**
   * 刷新 Token
   */
  refreshToken: (refreshToken: string) =>
    request.post<{ token: string }>('/api/auth/refresh', { refreshToken }),

  /**
   * 登出
   */
  logout: () =>
    request.post('/api/auth/logout'),

  /**
   * 发送验证码
   */
  sendCode: (phone: string) =>
    request.post('/api/auth/send-code', { phone }),
};
```

### 类型定义 (api/types/api.types.ts)
```typescript
/**
 * 请求配置
 */
export interface RequestConfig {
  url: string;
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE';
  data?: any;
  header?: Record<string, string>;
  timeout?: number;
}

/**
 * API 响应
 */
export interface ApiResponse<T = any> {
  code: number;
  data: T;
  message: string;
}
```

### 组合式函数 (hooks/useRequest.ts)
```typescript
import { ref } from 'vue';
import type { ApiResponse } from '@/api/types/api.types';

export function useRequest<T>() {
  const loading = ref(false);
  const error = ref<Error | null>(null);
  const data = ref<T | null>(null);

  /**
   * 发起请求
   */
  const execute = async (requestFn: () => Promise<ApiResponse<T>>) => {
    loading.value = true;
    error.value = null;

    try {
      const res = await requestFn();
      data.value = res.data;
      return res;
    } catch (err) {
      error.value = err as Error;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  return {
    loading,
    error,
    data,
    execute,
  };
}
```

### 使用示例

```vue
<script setup lang="ts">
import { authApi } from '@/api/modules/auth.api';
import { useUserStore } from '@/stores/user.store';
import { useRequest } from '@/hooks/useRequest';

const userStore = useUserStore();
const { loading, execute } = useRequest<LoginResult>();

const handleLogin = async () => {
  try {
    const res = await execute(() =>
      authApi.login({
        username: '13800138000',
        password: '123456',
      })
    );

    if (res.data) {
      userStore.token = res.data.token;
      userStore.userInfo = res.data.user;

      uni.navigateTo({ url: '/pages/tabbar/workspace' });
    }
  } catch (error) {
    console.error('登录失败', error);
  }
};
</script>
```

## Related Decisions
- ADR-008-01: 采用 UniApp 作为跨平台移动开发框架
- ADR-008-03: 使用 TypeScript 开发
