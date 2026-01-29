/**
 * 应用主入口
 */
import { createSSRApp } from 'vue';
import { createPinia } from 'pinia';
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate';
import App from './App.vue';

export function createApp() {
  const app = createSSRApp(App);

  // 创建 Pinia
  const pinia = createPinia();
  pinia.use(piniaPluginPersistedstate);
  app.use(pinia);

  // 恢复用户状态
  // import { useUserStore } from './stores/user.store';
  // const userStore = useUserStore();
  // userStore.restoreUser();

  return {
    app,
    pinia,
  };
}
