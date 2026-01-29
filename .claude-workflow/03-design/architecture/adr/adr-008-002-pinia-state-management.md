# ADR-008-002: 选择 Pinia 作为状态管理方案

## Status
Accepted

## Date
2026-01-28

## Context
移动端应用需要管理全局状态,包括:
- 用户认证状态 (Token, UserInfo)
- 设备列表和状态
- 消息列表和未读数
- 工作台数据
- 看板图表数据
- 个人中心设置

可选方案:
- **Pinia**: Vue 3 官方推荐的状态管理库
- **Vuex**: Vue 2 时代的状态管理方案
- **组合式函数 (Composables)**: 使用 Vue 3 Composition API 自定义 hooks

## Decision
选择 **Pinia** 作为状态管理方案。

### 理由

1. **Vue 3 官方推荐**:
   - Pinia 是 Vue 3 官方推荐的状态管理库
   - Vuex 5 将基于 Pinia 重构
   - 与 Vue 3 Composition API 深度集成

2. **TypeScript 友好**:
   - 完整的类型推断
   - 无需手动定义类型
   - API 设计优秀,类型安全

3. **API 简洁**:
   - 无需 mutations,直接修改状态
   - 去除复杂的命名空间
   - 支持组合式函数风格

4. **开发体验好**:
   - DevTools 支持
   - 模块化设计,每个 Store 独立
   - 支持 HMR (热模块替换)
   - 代码更简洁

5. **轻量级**:
   - 体积小 (约 1KB)
   - 性能优秀
   - 无额外依赖

### 与 Vuex 对比

| 特性 | Pinia | Vuex |
|------|-------|------|
| TypeScript 支持 | ✅ 完整类型推断 | ⚠️ 需要手动定义类型 |
| mutations | ❌ 不需要 | ✅ 必须 |
| actions | ✅ 异步操作 | ✅ 异步操作 |
| getters | ✅ 计算属性 | ✅ 计算属性 |
| 命名空间 | ✅ 默认支持 | ⚠️ 需要手动配置 |
| 模块化 | ✅ 天然模块化 | ⚠️ 需要手动配置 |
| Vue 3 支持 | ✅ 原生支持 | ⚠️ Vuex 4 支持,但不完善 |
| 代码量 | 更少 | 更多 |

## Consequences

### 正面影响

1. **代码更简洁**:
```typescript
// Pinia Store
export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: '',
    userInfo: null
  }),
  actions: {
    login(username, password) {
      // 直接修改 state,无需 mutation
      this.token = 'new_token'
    }
  }
})
```

```typescript
// Vuex (需要 mutation)
const store = createStore({
  state: {
    token: '',
    userInfo: null
  },
  mutations: {
    SET_TOKEN(state, token) {
      state.token = token
    }
  },
  actions: {
    login({ commit }, token) {
      commit('SET_TOKEN', token)
    }
  }
})
```

2. **TypeScript 类型安全**:
```typescript
// Pinia - 完整类型推断
const authStore = useAuthStore()
authStore.token // string 类型
authStore.login('admin', 'password') // 自动类型检查

// Vuex - 需要手动定义类型
const store = useStore()
store.state.auth.token // any 类型,需要断言
```

3. **模块化设计**:
- 每个 Store 文件独立
- 自动按文件名命名
- 无需手动注册模块

4. **开发体验提升**:
- DevTools 完美支持
- 代码提示更智能
- 重构更安全

### 负面影响

1. **生态相对较新**:
   - **缓解**: Pinia 已成熟,社区广泛采用
   - **缓解**: 官方推荐,文档完善

2. **团队需要学习**:
   - **缓解**: API 简单,学习成本低
   - **缓解**: 与 Vue 3 Composition API 一致

## Implementation

### Store 模块设计

```
src/stores/
├── index.ts              # 导出所有 stores
├── auth.ts               # 认证状态
├── device.ts             # 设备状态
├── workspace.ts          # 工作台状态
├── dashboard.ts          # 看板状态
├── message.ts            # 消息状态
├── profile.ts            # 个人中心状态
└── app.ts                # 应用全局状态
```

### Store 示例

```typescript
// src/stores/auth.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserInfo, LoginParams } from '@/types/auth'

export const useAuthStore = defineStore('auth', () => {
  // State
  const token = ref<string>('')
  const userInfo = ref<UserInfo | null>(null)

  // Getters
  const isLoggedIn = computed(() => !!token.value)

  // Actions
  const login = async (params: LoginParams) => {
    const res = await authApi.login(params)
    token.value = res.token
    userInfo.value = res.userInfo
  }

  const logout = () => {
    token.value = ''
    userInfo.value = null
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    login,
    logout
  }
})
```

### 在组件中使用

```vue
<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// 直接使用
console.log(authStore.token)
console.log(authStore.isLoggedIn)

// 调用 action
authStore.login('admin', 'password')
</script>
```

### 持久化状态

```typescript
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

// Store 配置持久化
export const useAuthStore = defineStore('auth', () => {
  // ...
}, {
  persist: {
    key: 'auth',
    storage: {
      getItem: uni.getStorageSync,
      setItem: uni.setStorageSync
    }
  }
})
```

## Alternatives Considered

### Vuex 4
**优点**:
- 成熟稳定,生态完善
- 团队可能已熟悉

**缺点**:
- 需要手动定义类型
- 代码更冗长 (需要 mutations)
- TypeScript 支持不如 Pinia

**拒绝理由**: Pinia 是 Vue 3 官方推荐,代码更简洁,类型安全更好

### 仅使用 Composables
**优点**:
- 无需额外依赖
- 更灵活

**缺点**:
- 缺少状态追踪
- 缺少 DevTools 支持
- 缺少持久化等高级功能

**拒绝理由**: Pinia 提供完整的 DevTools 支持和状态管理功能

## Related Decisions

- **ADR-008-001**: 选择 UniApp 作为跨平台框架
- **ADR-008-003**: 采用模块化分层架构

## References

- Pinia 官方文档: https://pinia.vuejs.org/zh/
- Vue 3 状态管理: https://cn.vuejs.org/guide/scaling-up/state-management.html

---

**决策者**: Claude Code
**审核者**: 待审核
**生效日期**: 2026-01-28
