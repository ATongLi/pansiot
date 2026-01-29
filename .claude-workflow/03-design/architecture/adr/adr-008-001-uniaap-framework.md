# ADR-008-001: 选择 UniApp 作为跨平台移动端框架

## Status
Accepted

## Date
2026-01-28

## Context
我们需要为工业物联网平台开发移动端应用,需要支持以下平台:
- iOS (iPhone/iPad)
- Android (智能手机/平板)
- 微信小程序
- H5 (移动浏览器)

核心需求包括:
1. **多端兼容**: 一套代码运行在多个平台
2. **开发效率**: 快速迭代,降低开发和维护成本
3. **性能要求**: 页面响应流畅,用户体验良好
4. **生态完善**: 组件库丰富,社区活跃
5. **团队技能**: 团队具备 Vue.js 经验

可选方案包括:
- **UniApp**: DCloud 推出的跨平台框架
- **React Native / Flutter**: 原生性能更好但学习成本高
- **Taro**: 京东出品,主要面向小程序
- **Cordova/PhoneGap**: 老牌混合应用框架,性能较差

## Decision
选择 **UniApp** 作为移动端跨平台框架。

### 理由

1. **真正的跨平台**:
   - 一套代码可编译到 iOS、Android、Web、以及各种小程序(微信/支付宝/百度/头条等)
   - 使用条件编译处理平台差异
   - 统一的 API 封装,屏蔽底层差异

2. **Vue 生态**:
   - 基于 Vue 3,团队熟悉度高
   - 支持 Composition API
   - TypeScript 支持良好
   - 可复用 Vue 生态的组件和工具

3. **性能优化**:
   - 原生渲染,性能接近原生应用
   - 优化的页面加载和渲染机制
   - 支持原生模块扩展

4. **开发体验**:
   - HBuilderX 可视化开发,也可使用 CLI
   - 热更新快速调试
   - 丰富的插件市场
   - 完善的文档和社区支持

5. **成本优势**:
   - 一套代码多端运行,大幅降低开发和维护成本
   - 开发效率高,快速迭代
   - 降低人力成本(不需要分别招聘 iOS/Android 开发)

## Consequences

### 正面影响

1. **开发效率提升**:
   - 一套代码多端发布
   - 减少 70% 以上的重复代码
   - 快速响应用户需求

2. **维护成本降低**:
   - 统一的代码库
   - 统一的版本管理
   - Bug 修复一次,多端生效

3. **团队协作友好**:
   - Vue 技术栈,学习成本低
   - 组件化开发,便于分工
   - 代码复用率高

4. **生态完善**:
   - UniUI 官方组件库
   - 插件市场资源丰富
   - 社区活跃,问题解决快

### 负面影响

1. **性能略逊于原生**:
   - **缓解**: UniApp 性能已足够好,对于工业物联网应用完全满足
   - **优化**: 使用虚拟列表、懒加载、分包加载等优化手段

2. **平台特性支持有限**:
   - **缓解**: UniApp 支持原生插件扩展
   - **缓解**: 使用条件编译处理平台差异
   - **缓解**: 大部分常用功能已封装

3. **依赖 UniApp 版本更新**:
   - **缓解**: DCloud 团队活跃,版本更新及时
   - **缓解**: 社区成熟,问题解决快

4. **调试复杂度**:
   - **缓解**: 多端测试需要真机或模拟器
   - **缓解**: HBuilderX 提供良好的调试工具

### 技术选型配套

基于 UniApp,我们选择以下配套技术:

```
框架: UniApp (Vue 3 + TypeScript)
状态管理: Pinia (Vue 3 官方推荐)
UI 组件: UniUI (官方组件库)
构建工具: Vite (快速热更新)
代码规范: ESLint + Prettier
```

## Alternatives Considered

### React Native
**优点**:
- 原生性能更好
- Facebook 支持,生态成熟
- 可使用 React 生态

**缺点**:
- 需要学习 React/React Native
- 不支持小程序(需要额外框架)
- 开发效率不如 UniApp
- iOS 和 Android 部分功能需要分别实现

**拒绝理由**: 团队不具备 React 经验,开发成本高,不支持小程序

### Flutter
**优点**:
- 性能优秀,接近原生
- UI 渲染能力强
- Google 支持

**缺点**:
- 需要学习 Dart 语言
- 不支持小程序
- 生态不如 UniApp 丰富
- 团队学习成本高

**拒绝理由**: 学习成本高,不支持小程序,团队无 Dart 经验

### Taro
**优点**:
- 京东出品,质量可靠
- 支持 React/Vue
- 小程序性能好

**缺点**:
- 主要面向小程序,对 App 支持不如 UniApp
- 多端支持不如 UniApp 完善
- 社区资源不如 UniApp

**拒绝理由**: App 支持不如 UniApp,不是我们的主要需求

## Implementation

### 项目初始化

```bash
# 使用 Vue 3 + TypeScript 模板创建项目
npx degit dcloudio/uni-preset-vue#vite-ts pansiot-app

# 或使用 HBuilderX 可视化创建
```

### 目录结构

```
src/
├── api/              # API 封装
├── components/       # 组件
├── pages/           # 页面
├── stores/          # Pinia Store
├── composables/     # Composables
├── utils/           # 工具函数
├── types/           # TypeScript 类型
└── styles/          # 全局样式
```

### 条件编译示例

```vue
<template>
  <view>
    <!-- #ifdef H5 -->
    <web-view src="https://example.com"></web-view>
    <!-- #endif -->

    <!-- #ifdef MP-WEIXIN -->
    <button open-type="getUserInfo">获取用户信息</button>
    <!-- #endif -->

    <!-- #ifdef APP-PLUS -->
    <button @click="scan">扫码</button>
    <!-- #endif -->
  </view>
</template>

<script setup lang="ts">
// #ifdef APP-PLUS
const scan = () => {
  plus.barcode.scan()
}
// #endif
</script>
```

## Related Decisions

- **ADR-008-002**: 选择 Pinia 作为状态管理方案
- **ADR-008-003**: 采用模块化分层架构

## References

- UniApp 官方文档: https://uniapp.dcloud.net.cn/
- Vue 3 官方文档: https://cn.vuejs.org/
- UniApp 插件市场: https://ext.dcloud.net.cn/

---

**决策者**: Claude Code
**审核者**: 待审核
**生效日期**: 2026-01-28
