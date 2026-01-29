# ADR-008-05: 采用模块化架构组织代码

## Status
Accepted

## Date
2026-01-28

## Context
移动端应用包含多个业务模块(登录、设备、工作台、看板、消息、我的),需要选择合适的代码组织方式来保证代码的可维护性、可扩展性和团队协作效率。

## Decision
采用 **模块化架构** 按业务模块组织代码。

### Rationale
1. **职责清晰**: 每个模块独立负责自己的业务逻辑
2. **并行开发**: 不同模块可以并行开发,互不干扰
3. **易于维护**: 模块边界清晰,修改影响范围小
4. **易于扩展**: 新增功能只需新增模块
5. **代码复用**: 模块内部组件可复用
6. **团队协作**: 团队成员可以负责不同模块

### Alternatives Considered

#### 1. 按技术分层组织 (pages/components/stores)
- **优势**:
  - 结构简单直观
  - 适合小型项目
- **劣势**:
  - 业务逻辑分散
  - 难以维护
  - 模块间耦合高
- **结论**: 不选择,不适合中大型项目

#### 2. 单体应用 (所有代码放一起)
- **优势**:
  - 最简单
- **劣势**:
  - 难以维护
  - 难以并行开发
  - 代码冲突多
- **结论**: 不选择

#### 3. Micro-frontends (微前端)
- **优势**:
  - 模块完全独立
  - 技术栈灵活
- **劣势**:
  - 过度设计
  - 复杂度高
  - 性能开销
- **结论**: 不选择,移动端不适合

## Consequences

### Positive
1. **清晰的模块边界**: 每个模块独立,职责明确
2. **支持并行开发**: 团队成员可以开发不同模块
3. **易于维护**: 修改某个模块不影响其他模块
4. **易于扩展**: 新增功能只需新增模块
5. **代码复用**: 模块内部组件高度复用
6. **降低耦合**: 模块间通过 API 和 Store 通信

### Negative
1. **初期搭建成本**: 需要设计模块结构和通信机制
2. **模块间通信**: 需要统一的通信规范

### Mitigation Strategies
1. **初期成本**:
   - 提供模块模板
   - 提供脚手架工具
   - 建立开发规范

2. **模块间通信**:
   - 使用 Pinia Store 进行状态共享
   - 使用事件总线进行事件通信
   - 使用路由参数进行数据传递

## Implementation

### 模块结构设计

```
src/modules/
├── auth/                    # 认证模块
│   ├── pages/              # 模块页面
│   │   ├── login.vue
│   │   └── register.vue
│   ├── components/         # 模块组件
│   │   ├── LoginForm.vue
│   │   └── RegisterForm.vue
│   ├── services/           # 模块服务
│   │   └── auth.service.ts
│   ├── stores/             # 模块状态
│   │   └── auth.store.ts
│   ├── types/              # 模块类型
│   │   └── auth.types.ts
│   └── index.ts            # 模块导出
│
├── device/                  # 设备模块
│   ├── pages/
│   │   ├── list.vue
│   │   └── detail.vue
│   ├── components/
│   │   ├── DeviceCard.vue
│   │   ├── DeviceFilter.vue
│   │   └── DeviceStatus.vue
│   ├── services/
│   │   └── device.service.ts
│   ├── stores/
│   │   └── device.store.ts
│   ├── types/
│   │   └── device.types.ts
│   └── index.ts
│
├── workspace/               # 工作台模块
├── dashboard/               # 看板模块
├── message/                 # 消息模块
└── profile/                 # 我的模块
```

### 模块规范

#### 1. 模块导出 (modules/auth/index.ts)
```typescript
// 导出组件
export { default as LoginForm } from './components/LoginForm.vue';
export { default as RegisterForm } from './components/RegisterForm.vue';

// 导出 Service
export { authApi } from './services/auth.service';

// 导出 Store
export { useAuthStore } from './stores/auth.store';

// 导出类型
export type * from './types/auth.types';
```

#### 2. 模块 Service
```typescript
// modules/device/services/device.service.ts
import { request } from '@/api/client/request';
import type { Device, DeviceListParams, DeviceListResult } from '../types/device.types';

/**
 * 设备服务
 */
export const deviceService = {
  /**
   * 获取设备列表
   */
  async getList(params: DeviceListParams): Promise<DeviceListResult> {
    const res = await request.get<DeviceListResult>('/api/devices', params);
    return res.data;
  },

  /**
   * 获取设备详情
   */
  async getDetail(id: string): Promise<Device> {
    const res = await request.get<Device>(`/api/devices/${id}`);
    return res.data;
  },

  /**
   * 控制设备
   */
  async control(id: string, command: any): Promise<void> {
    await request.post(`/api/devices/${id}/control`, command);
  },
};
```

#### 3. 模块 Store
```typescript
// modules/device/stores/device.store.ts
import { defineStore } from 'pinia';
import { deviceService } from '../services/device.service';
import type { Device, DeviceListParams } from '../types/device.types';

interface DeviceState {
  list: Device[];
  currentDevice: Device | null;
  loading: boolean;
  total: number;
}

export const useDeviceStore = defineStore('device', {
  state: (): DeviceState => ({
    list: [],
    currentDevice: null,
    loading: false,
    total: 0,
  }),

  actions: {
    /**
     * 加载设备列表
     */
    async fetchList(params: DeviceListParams) {
      this.loading = true;
      try {
        const res = await deviceService.getList(params);
        this.list = res.list;
        this.total = res.total;
      } finally {
        this.loading = false;
      }
    },

    /**
     * 加载设备详情
     */
    async fetchDetail(id: string) {
      this.loading = true;
      try {
        const res = await deviceService.getDetail(id);
        this.currentDevice = res;
      } finally {
        this.loading = false;
      }
    },
  },
});
```

#### 4. 模块组件
```vue
<!-- modules/device/components/DeviceCard.vue -->
<script setup lang="ts">
import { computed } from 'vue';
import type { Device } from '../types/device.types';

interface Props {
  device: Device;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'click', device: Device): void;
  (e: 'toggle', device: Device): void;
}>();

const statusText = computed(() => {
  const statusMap = {
    normal: '正常',
    warning: '警告',
    error: '故障',
    offline: '离线',
  };
  return statusMap[props.device.status] || '未知';
});

const handleClick = () => {
  emit('click', props.device);
};

const handleToggle = () => {
  emit('toggle', props.device);
};
</script>

<template>
  <view class="device-card" @click="handleClick">
    <view class="device-name">{{ device.name }}</view>
    <view class="device-status" :class="`status-${device.status}`">
      {{ statusText }}
    </view>
  </view>
</template>

<style scoped lang="scss">
.device-card {
  padding: 20rpx;
  background: #fff;
  border-radius: 16rpx;
  margin-bottom: 20rpx;

  .device-name {
    font-size: 32rpx;
    font-weight: bold;
  }

  .device-status {
    margin-top: 10rpx;
    font-size: 24rpx;

    &.status-normal { color: #52c41a; }
    &.status-warning { color: #faad14; }
    &.status-error { color: #f5222d; }
    &.status-offline { color: #8c8c8c; }
  }
}
</style>
```

### 模块间通信

#### 1. 通过共享 Store
```typescript
// stores/tenant.store.ts - 全局租户状态
export const useTenantStore = defineStore('tenant', {
  state: () => ({
    currentTenant: null as Tenant | null,
  }),
  actions: {
    switchTenant(tenantId: string) {
      // 切换租户后,所有模块重新加载数据
      this.currentTenant = this.tenantList.find(t => t.id === tenantId);
    },
  },
});

// modules/device/stores/device.store.ts
export const useDeviceStore = defineStore('device', {
  actions: {
    async fetchList(params: DeviceListParams) {
      // 从租户 Store 获取当前租户
      const tenantStore = useTenantStore();
      const tenantId = tenantStore.currentTenant?.id;

      // 调用 API 时带上租户 ID
      const res = await deviceService.getList({
        ...params,
        tenantId,
      });

      this.list = res.list;
    },
  },
});
```

#### 2. 通过路由传参
```typescript
// 从设备列表跳转到设备详情
// pages/tabbar/device.vue
const goToDetail = (device: Device) => {
  uni.navigateTo({
    url: `/pages/device/detail?id=${device.id}`,
  });
};

// pages/device/detail.vue
const onLoad = (options: { id: string }) => {
  const deviceId = options.id;
  // 加载设备详情
  deviceStore.fetchDetail(deviceId);
};
```

#### 3. 通过事件总线 (可选)
```typescript
// utils/event-bus.ts
import { ref } from 'vue';

type EventCallback = (...args: any[]) => void;

class EventBus {
  private events: Record<string, EventCallback[]> = {};

  on(event: string, callback: EventCallback) {
    if (!this.events[event]) {
      this.events[event] = [];
    }
    this.events[event].push(callback);
  }

  emit(event: string, ...args: any[]) {
    const callbacks = this.events[event];
    if (callbacks) {
      callbacks.forEach(callback => callback(...args));
    }
  }

  off(event: string, callback?: EventCallback) {
    if (!callback) {
      delete this.events[event];
    } else {
      const callbacks = this.events[event];
      if (callbacks) {
        const index = callbacks.indexOf(callback);
        if (index > -1) {
          callbacks.splice(index, 1);
        }
      }
    }
  }
}

export const eventBus = new EventBus();

// 使用: 设备状态变化通知
eventBus.on('device:status-change', (device: Device) => {
  console.log('设备状态变化', device);
});

eventBus.emit('device:status-change', device);
```

### 模块复用策略

#### 1. 通用组件放在 components/common
```
src/components/common/
├── AppNavBar.vue          # 所有模块可用的导航栏
├── PageContainer.vue      # 页面容器
├── Loading.vue            # 加载指示器
└── EmptyState.vue         # 空状态
```

#### 2. 业务组件放在各自模块
```
src/modules/device/components/
├── DeviceCard.vue         # 仅设备模块使用
├── DeviceFilter.vue       # 仅设备模块使用
└── DeviceStatus.vue       # 仅设备模块使用
```

#### 3. 跨模块共享组件
如果多个模块需要某个组件,将其提升到 `components/common` 或 `components/business`

## 最佳实践

### 1. 模块独立性
- 每个模块应尽可能独立
- 避免模块间直接调用
- 通过 Store 和 API 进行通信

### 2. 模块导出
- 每个模块应有 `index.ts` 统一导出
- 明确模块的公共 API

### 3. 模块命名
- 模块名使用单数形式 (auth 而不是 auths)
- 模块内组件、Service、Store 使用模块名前缀
  - `auth.store.ts`
  - `auth.service.ts`
  - `auth.types.ts`

### 4. 依赖管理
- 模块间依赖尽量减少
- 优先依赖全局 Store 而不是其他模块 Store
- 使用依赖注入而不是直接引用

## Related Decisions
- ADR-008-01: 采用 UniApp 作为跨平台移动开发框架
- ADR-008-02: 使用 Pinia 作为状态管理方案
