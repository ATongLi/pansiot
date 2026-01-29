# ADR-008-003: 采用模块化分层架构

## Status
Accepted

## Date
2026-01-28

## Context
移动端应用包含 6 个核心业务模块:
1. **Auth** - 登录认证
2. **Device** - 设备管理
3. **Workspace** - 工作台
4. **Dashboard** - 看板
5. **Message** - 消息
6. **Profile** - 个人中心

每个模块都有独立的页面、状态管理、API 接口和业务逻辑。

可选的架构方案:
- **模块化分层架构** (推荐)
- **单体架构** (所有代码混在一起)
- **微前端架构** (过度设计)

## Decision
采用 **模块化分层架构**。

### 架构设计

#### 三层架构

```
Presentation Layer (展示层)
    ↓
Business Layer (业务层)
    ↓
Data Layer (数据层)
```

#### 模块化组织

每个业务模块独立组织:

```
device-module/
├── pages/          # 展示层 - 页面
├── components/     # 展示层 - 组件
├── stores/         # 业务层 - 状态
├── composables/    # 业务层 - 逻辑
├── services/       # 业务层 - 服务
└── api/            # 数据层 - 接口
```

### 理由

1. **职责清晰**:
   - 展示层只负责 UI 渲染
   - 业务层负责逻辑和状态
   - 数据层负责接口和存储

2. **模块独立**:
   - 每个业务模块独立开发
   - 降低模块间耦合
   - 便于团队分工

3. **易于测试**:
   - 每层独立测试
   - Mock 依赖简单
   - 单元测试覆盖率高

4. **可维护性强**:
   - 代码组织清晰
   - 定位问题快速
   - 重构风险低

5. **可扩展性好**:
   - 新增模块简单
   - 不影响现有模块
   - 支持功能按需加载

## Consequences

### 正面影响

1. **代码组织清晰**:
```
src/
├── pages/device/          # 设备相关页面
│   ├── list/
│   └── detail/
├── stores/device.ts       # 设备状态
├── composables/useDevice.ts  # 设备逻辑
├── services/DeviceService.ts # 设备服务
└── api/device.ts          # 设备API
```

2. **团队协作友好**:
   - 开发 A 可以专注 Device 模块
   - 开发 B 可以专注 Message 模块
   - 减少代码冲突

3. **测试覆盖率高**:
   - 每层独立测试
   - Mock 简单
   - 测试运行快

4. **性能优化**:
   - 支持按需加载 (路由懒加载)
   - 支持分包加载 (UniApp 分包)
   - 减小主包体积

### 负面影响

1. **初期文件数量多**:
   - **缓解**: 使用脚手架自动生成模板
   - **缓解**: 文件组织清晰,长期维护更简单

2. **学习成本**:
   - **缓解**: 架构简单,易于理解
   - **缓解**: 提供开发文档和示例

## Implementation

### 目录结构

```
src/
├── api/                    # Data Layer - API 封装
│   ├── auth.ts
│   ├── device.ts
│   ├── workspace.ts
│   ├── dashboard.ts
│   ├── message.ts
│   └── profile.ts
│
├── components/             # Presentation Layer - 组件
│   ├── common/            # 通用组件
│   │   ├── CustomNavBar/
│   │   ├── PageContainer/
│   │   └── Loading/
│   └── business/          # 业务组件
│       ├── auth/
│       ├── device/
│       └── ...
│
├── pages/                  # Presentation Layer - 页面
│   ├── index/             # 启动页
│   ├── auth/              # 登录模块
│   ├── device/            # 设备模块
│   ├── workspace/         # 工作台模块
│   ├── dashboard/         # 看板模块
│   ├── message/           # 消息模块
│   └── profile/           # 个人中心模块
│
├── stores/                 # Business Layer - 状态
│   ├── auth.ts
│   ├── device.ts
│   ├── workspace.ts
│   ├── dashboard.ts
│   ├── message.ts
│   ├── profile.ts
│   └── app.ts
│
├── composables/           # Business Layer - 组合式函数
│   ├── useAuth.ts
│   ├── useDevice.ts
│   ├── useRequest.ts
│   └── usePagination.ts
│
├── services/              # Business Layer - 业务服务
│   ├── AuthService.ts
│   ├── DeviceService.ts
│   └── ...
│
├── utils/                 # Data Layer - 工具类
│   ├── request.ts
│   ├── storage.ts
│   ├── validator.ts
│   └── format.ts
│
├── types/                 # Data Layer - 类型定义
│   ├── auth.d.ts
│   ├── device.d.ts
│   └── common.d.ts
│
└── styles/                # 全局样式
    ├── variables.scss
    └── mixins.scss
```

### 模块示例

```typescript
// api/device.ts - 数据层
import { request } from '@/utils/request'
import type { Device, DeviceListParams } from '@/types/device'

export const deviceApi = {
  getList(params: DeviceListParams) {
    return request.get<Device[]>('/api/devices', { params })
  },

  getDetail(id: string) {
    return request.get<Device>(`/api/devices/${id}`)
  }
}
```

```typescript
// stores/device.ts - 业务层 (状态)
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { deviceApi } from '@/api/device'

export const useDeviceStore = defineStore('device', () => {
  const devices = ref<Device[]>([])

  const fetchDevices = async (params: DeviceListParams) => {
    const res = await deviceApi.getList(params)
    devices.value = res
  }

  return {
    devices,
    fetchDevices
  }
})
```

```typescript
// composables/useDevice.ts - 业务层 (逻辑)
import { useDeviceStore } from '@/stores/device'
import { useRequest } from '@/composables/useRequest'

export function useDevice() {
  const deviceStore = useDeviceStore()

  const { loading, error, execute } = useRequest(
    deviceStore.fetchDevices
  )

  return {
    devices: deviceStore.devices,
    loading,
    error,
    fetchDevices: execute
  }
}
```

```vue
<!-- pages/device/list/index.vue - 展示层 -->
<template>
  <view>
    <DeviceCard
      v-for="device in devices"
      :key="device.id"
      :device="device"
    />
    <Loading v-if="loading" />
  </view>
</template>

<script setup lang="ts">
import { useDevice } from '@/composables/useDevice'

const { devices, loading, fetchDevices } = useDevice()

onMounted(() => {
  fetchDevices({ page: 1, pageSize: 20 })
})
</script>
```

### 分层原则

| 层级 | 职责 | 可调用 | 不可调用 |
|------|------|--------|----------|
| **Presentation** | UI 渲染、用户交互 | Business | Data |
| **Business** | 业务逻辑、状态管理 | Data | Presentation |
| **Data** | API、存储、工具 | External | Business, Presentation |

## Alternatives Considered

### 单体架构
**优点**:
- 文件结构简单

**缺点**:
- 代码组织混乱
- 难以维护
- 团队协作困难

**拒绝理由**: 代码混乱,难以长期维护

### 微前端架构
**优点**:
- 模块完全独立

**缺点**:
- 过度设计
- 复杂度高
- 性能开销大

**拒绝理由**: 移动端应用不需要微前端,复杂度过高

## Related Decisions

- **ADR-008-001**: 选择 UniApp 作为跨平台框架
- **ADR-008-002**: 选择 Pinia 作为状态管理方案

## References

- Clean Architecture: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
- Vue 3 风格指南: https://cn.vuejs.org/style-guide/

---

**决策者**: Claude Code
**审核者**: 待审核
**生效日期**: 2026-01-28
