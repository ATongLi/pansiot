/**
 * 应用状态管理
 */
import { defineStore } from 'pinia';

export const useAppStore = defineStore('app', {
  state: () => ({
    theme: 'light' as 'light' | 'dark',
    language: 'zh-CN',
    networkStatus: true as boolean,
  }),

  actions: {
    /**
     * 设置主题
     */
    setTheme(theme: 'light' | 'dark') {
      this.theme = theme;
    },

    /**
     * 设置语言
     */
    setLanguage(language: string) {
      this.language = language;
    },

    /**
     * 设置网络状态
     */
    setNetworkStatus(status: boolean) {
      this.networkStatus = status;
    },
  },
});
