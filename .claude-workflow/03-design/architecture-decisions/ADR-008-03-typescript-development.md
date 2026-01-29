# ADR-008-03: 使用 TypeScript 进行全栈类型安全开发

## Status
Accepted

## Date
2026-01-28

## Context
移动端应用需要保证代码质量、减少运行时错误、提高可维护性。需要决定是否使用 TypeScript 以及如何使用。

## Decision
**全面使用 TypeScript** 进行开发,开启 strict 模式。

### Rationale
1. **类型安全**: 编译时类型检查,减少运行时错误
2. **IDE 支持**: 完善的代码提示、自动补全、重构功能
3. **可维护性**: 类型即文档,代码可读性高
4. **重构安全**: 类型检查保证重构的安全性
5. **团队协作**: 类型约束降低沟通成本
6. **现代前端**: TypeScript 已成为前端主流,生态成熟
7. **Vue 3 支持**: Vue 3 对 TypeScript 支持完美

### Alternatives Considered

#### 1. JavaScript
- **优势**:
  - 无需编译,直接运行
  - 学习成本低
- **劣势**:
  - 无类型检查,易出错
  - IDE 支持弱
  - 大型项目维护困难
- **结论**: 不选择

#### 2. JSDoc (Type Annotation via Comments)
- **优势**:
  - 保持 JavaScript 语法
  - 提供类型提示
- **劣势**:
  - 类型检查弱
  - 代码冗余
  - 维护成本高
- **结论**: 不选择

## Consequences

### Positive
1. **减少错误**: 编译时发现大部分类型错误,减少运行时错误
2. **提高效率**: IDE 智能提示,开发效率提升
3. **代码质量**: 类型约束强制代码规范化
4. **重构安全**: 修改代码时类型检查保证安全性
5. **文档化**: 类型定义即文档,无需额外注释
6. **团队协作**: 接口定义明确,降低协作成本

### Negative
1. **学习成本**: 团队需要学习 TypeScript(但学习曲线低)
2. **初期繁琐**: 需要编写类型定义,初期开发略慢
3. **编译时间**: 需要编译,但 Vite 速度快,影响小

### Mitigation Strategies
1. **学习成本**:
   - 提供文档和教程
   - 组织培训
   - 从简单项目开始实践

2. **初期繁琐**:
   - 使用类型推断减少显式类型
   - 使用泛型提高复用性
   - 建立类型定义模板库

3. **编译时间**:
   - 使用 Vite,启动和热更新极快
   - 增量编译
   - 合理配置 tsconfig.json

## Implementation

### TypeScript 配置

#### tsconfig.json
```json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "ESNext",
    "strict": true,                    // 开启严格模式
    "esModuleInterop": true,
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "preserve",
    "strictNullChecks": true,          // 严格空值检查
    "strictFunctionTypes": true,       // 严格函数类型
    "strictPropertyInitialization": true,
    "noImplicitAny": true,             // 禁止隐式 any
    "noImplicitThis": true,            // 禁止隐式 this
    "alwaysStrict": true,
    "noUnusedLocals": true,            // 检查未使用的变量
    "noUnusedParameters": true,        // 检查未使用的参数
    "noImplicitReturns": true,         // 检查函数返回值
    "noFallthroughCasesInSwitch": true
  },
  "include": [
    "src/**/*.ts",
    "src/**/*.d.ts",
    "src/**/*.tsx",
    "src/**/*.vue"
  ],
  "exclude": [
    "node_modules",
    "dist",
    "build"
  ]
}
```

### 类型定义规范

#### 1. 接口优先
```typescript
// ✅ 推荐: 使用 interface 定义对象类型
interface UserInfo {
  id: string;
  username: string;
  email?: string;
}

// ❌ 避免: 除非定义联合类型
type UserInfo = {
  id: string;
  username: string;
};
```

#### 2. 类型文件组织
```
src/
├── types/                  # 全局类型
│   ├── global.d.ts        # 全局类型声明
│   ├── vue-shim.d.ts      # Vue 类型声明
│   └── uniapp.d.ts        # UniApp 类型声明
├── api/
│   └── types/             # API 类型
│       └── api.types.ts
└── modules/
    └── {module}/
        └── types/         # 模块类型
            └── {module}.types.ts
```

#### 3. 类型定义示例
```typescript
// types/api.types.ts
export interface ApiResponse<T = any> {
  code: number;
  data: T;
  message: string;
}

export interface PageParams {
  page: number;
  pageSize: number;
}

export interface PageResult<T> {
  list: T[];
  total: number;
}

// api/modules/auth.api.ts
export interface LoginParams {
  username: string;
  password: string;
  tenantId?: string;
}

export interface LoginResult {
  token: string;
  refreshToken: string;
  expiresIn: number;
  user: UserInfo;
}
```

### 组件类型定义

#### Vue 3 Composition API + TypeScript
```vue
<script setup lang="ts">
// 定义 Props
interface Props {
  title: string;
  count?: number;
}
const props = withDefaults(defineProps<Props>(), {
  count: 0,
});

// 定义 Emits
interface Emits {
  (e: 'update', value: number): void;
  (e: 'delete', id: string): void;
}
const emit = defineEmits<Emits>();

// Ref 类型推断
const count = ref<number>(0);
const userInfo = ref<UserInfo | null>(null);

// Computed
const doubleCount = computed<number>(() => count.value * 2);

// 方法
const handleClick = (): void => {
  emit('update', count.value + 1);
};
</script>
```

### Store 类型定义

```typescript
// stores/user.store.ts
import { defineStore } from 'pinia';
import type { UserInfo, LoginParams } from '@/modules/auth/types';

interface UserState {
  token: string;
  userInfo: UserInfo | null;
  isLoggedIn: boolean;
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    token: '',
    userInfo: null,
    isLoggedIn: false,
  }),

  getters: {
    hasPermission: (state) => {
      return (permission: string): boolean => {
        return state.userInfo?.roles.some(role =>
          role.permissions.includes(permission)
        ) ?? false;
      };
    },
  },

  actions: {
    async login(params: LoginParams): Promise<void> {
      const res = await authApi.login(params);
      this.token = res.token;
      this.userInfo = res.user;
      this.isLoggedIn = true;
    },
  },
});
```

### 最佳实践

#### 1. 避免使用 any
```typescript
// ❌ 避免
function process(data: any) {
  return data.value;
}

// ✅ 推荐: 使用 unknown 或泛型
function process<T>(data: T): T {
  return data;
}

// ✅ 推荐: 使用具体类型
function process(data: UserInfo): string {
  return data.username;
}
```

#### 2. 使用类型推断
```typescript
// ❌ 冗余
const count: number = 0;

// ✅ 推荐: 利用类型推断
const count = 0;
```

#### 3. 使用泛型提高复用性
```typescript
// ✅ 推荐: 使用泛型
interface ApiResponse<T = any> {
  code: number;
  data: T;
  message: string;
}

const userResponse = await request.get<ApiResponse<UserInfo>>('/api/user');
const deviceResponse = await request.get<ApiResponse<Device[]>>('/api/devices');
```

## Related Decisions
- ADR-008-01: 采用 UniApp 作为跨平台移动开发框架
- ADR-008-02: 使用 Pinia 作为状态管理方案
