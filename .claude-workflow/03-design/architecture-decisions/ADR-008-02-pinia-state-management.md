# ADR-008-02: 使用 Pinia 作为状态管理方案

## Status
Accepted

## Date
2026-01-28

## Context
移动端应用需要管理全局状态(用户信息、租户信息、应用配置等)和模块状态(设备列表、消息列表等)。需要选择一个状态管理方案来统一管理应用状态。

## Decision
使用 **Pinia** 作为状态管理方案。

### Rationale
1. **Vue 3 官方推荐**: Pinia 是 Vue 官方推荐的状态管理库,是 Vuex 的继任者
2. **TypeScript 友好**: 完美的 TypeScript 支持,类型推断准确
3. **API 简洁**: 相比 Vuex,API 更简洁直观,学习成本低
4. **Composition API**: 原生支持 Composition API,与 Vue 3 风格一致
5. **轻量级**: 体积小(~1KB),性能优秀
6. **模块化**: 天然支持模块化,Store 之间可以相互调用
7. **DevTools**: 完美的 Vue DevTools 集成
8. **持久化**: 支持插件,易于实现状态持久化

### Alternatives Considered

#### 1. Vuex
- **优势**:
  - Vue 官方状态管理(旧版)
  - 生态成熟,文档完善
- **劣势**:
  - TypeScript 支持不友好
  - API 冗余,学习曲线陡峭
  - 对 Vue 3 Composition API 支持不好
- **结论**: 不选择,已被 Pinia 取代

#### 2. Vuex 5 (Vuex Next)
- **优势**:
  - Vuex 的下一代版本
- **劣势**:
  - 仍处于开发阶段,不稳定
  - 不如直接使用 Pinia
- **结论**: 不选择

#### 3. 无状态管理 (使用 props/emit)
- **优势**:
  - 简单直接,无需额外库
- **劣势**:
  - 跨组件通信困难
  - 状态分散,难以管理
  - 不适合复杂应用
- **结论**: 不选择

## Consequences

### Positive
1. **类型安全**: 完整的 TypeScript 支持,编译时类型检查
2. **开发体验**: API 简洁直观,代码可读性高
3. **模块化**: 天然模块化,Store 可以按业务模块划分
4. **可维护性**: 代码结构清晰,易于维护和扩展
5. **性能优秀**: 轻量级,性能优于 Vuex
6. **DevTools**: 完美的开发工具集成,调试方便

### Negative
1. **学习成本**: 团队需要学习新 API(但学习曲线低)
2. **生态系统**: 相比 Vuex,插件和社区资源稍少(但快速增长)

### Mitigation Strategies
1. **学习成本**:
   - 提供详细的 Pinia 使用文档
   - 提供代码示例和最佳实践
   - 组织团队培训和分享

2. **生态问题**:
   - Pinia 生态已足够成熟
   - 常用功能都有对应插件
   - 不满足的功能可以自己开发插件

## Implementation

### Store 结构
```typescript
// stores/app.store.ts - 应用状态
export const useAppStore = defineStore('app', {
  state: () => ({
    theme: 'light',
    language: 'zh-CN',
    networkStatus: true,
  }),
  actions: {
    setTheme(theme: string) { ... },
    setLanguage(lang: string) { ... },
    setNetworkStatus(status: boolean) { ... },
  },
});

// stores/user.store.ts - 用户状态
export const useUserStore = defineStore('user', {
  state: () => ({
    token: '',
    userInfo: null as UserInfo | null,
    isLoggedIn: false,
  }),
  getters: {
    hasPermission: (state) => (permission: string) => { ... },
  },
  actions: {
    async login(params: LoginParams) { ... },
    async logout() { ... },
  },
  persist: true, // 持久化
});

// stores/tenant.store.ts - 租户状态
export const useTenantStore = defineStore('tenant', {
  state: () => ({
    currentTenant: null as Tenant | null,
    tenantList: [] as Tenant[],
  }),
  actions: {
    switchTenant(tenantId: string) { ... },
  },
});
```

### 模块化 Store
每个业务模块有自己的 Store:
- `modules/auth/stores/auth.store.ts`
- `modules/device/stores/device.store.ts`
- `modules/workspace/stores/workspace.store.ts`
- `modules/dashboard/stores/dashboard.store.ts`
- `modules/message/stores/message.store.ts`
- `modules/profile/stores/profile.store.ts`

### 持久化策略
使用 `pinia-plugin-persistedstate` 插件实现状态持久化:
- **用户状态**: 持久化到本地存储
- **租户状态**: 持久化到本地存储
- **应用状态**: 持久化部分配置(主题、语言)

## Related Decisions
- ADR-008-01: 采用 UniApp 作为跨平台移动开发框架
- ADR-008-03: 使用 TypeScript 开发
