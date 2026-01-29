<template>
  <PageContainer bg-color="#ffffff">
    <view class="login-page">
      <!-- Logo 区域 -->
      <view class="login-logo">
        <image class="logo-image" src="/static/images/logo.png" mode="aspectFit"></image>
        <text class="logo-title">PansIot</text>
        <text class="logo-subtitle">工业物联网移动平台</text>
      </view>

      <!-- 登录表单 -->
      <view class="login-form">
        <!-- 用户名/手机号 -->
        <view class="form-item">
          <view class="form-item-icon">
            <uni-icons type="person" size="20" color="#999999"></uni-icons>
          </view>
          <input
            class="form-item-input"
            type="text"
            placeholder="请输入用户名/手机号"
            v-model="formData.username"
            placeholder-style="color: #cccccc"
          />
        </view>

        <!-- 密码 -->
        <view class="form-item">
          <view class="form-item-icon">
            <uni-icons type="locked" size="20" color="#999999"></uni-icons>
          </view>
          <input
            class="form-item-input"
            type="password"
            placeholder="请输入密码"
            v-model="formData.password"
            placeholder-style="color: #cccccc"
          />
        </view>

        <!-- 记住密码 & 忘记密码 -->
        <view class="form-options">
          <view class="form-remember" @click="handleRememberToggle">
            <uni-icons
              :type="rememberPassword ? 'checkbox' : 'circle'"
              :color="rememberPassword ? '#007aff' : '#cccccc'"
              size="18"
            ></uni-icons>
            <text class="form-remember-text">记住密码</text>
          </view>
          <text class="form-forgot" @click="handleForgotPassword">忘记密码?</text>
        </view>

        <!-- 登录按钮 -->
        <button class="login-button" type="primary" :loading="loading" @click="handleLogin">
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </view>

      <!-- 其他登录方式 -->
      <view class="login-other">
        <text class="login-other-text">其他登录方式</text>
        <view class="login-other-icons">
          <view class="login-other-icon" @click="handleWechatLogin">
            <uni-icons type="weixin" size="24" color="#09bb07"></uni-icons>
          </view>
        </view>
      </view>
    </view>
  </PageContainer>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { validateUsername, validatePassword, validatePhone } from '@/utils/validator';
import { useAuth } from '@/composables/useAuth';
import PageContainer from '@/components/common/PageContainer/index.vue';

interface FormData {
  username: string;
  password: string;
}

const formData = ref<FormData>({
  username: '',
  password: '',
});

const loading = ref(false);
const rememberPassword = ref(false);

const { login } = useAuth();

/**
 * 登录
 */
const handleLogin = async () => {
  const { username, password } = formData.value;

  // 表单验证
  if (!username) {
    uni.showToast({
      title: '请输入用户名/手机号',
      icon: 'none',
    });
    return;
  }

  if (!validateUsername(username) && !validatePhone(username)) {
    uni.showToast({
      title: '用户名/手机号格式不正确',
      icon: 'none',
    });
    return;
  }

  if (!password) {
    uni.showToast({
      title: '请输入密码',
      icon: 'none',
    });
    return;
  }

  if (!validatePassword(password)) {
    uni.showToast({
      title: '密码格式不正确',
      icon: 'none',
      duration: 2000,
    });
    return;
  }

  // 执行登录
  loading.value = true;
  try {
    const success = await login(username, password);

    if (success) {
      // 记住密码
      if (rememberPassword.value) {
        uni.setStorageSync('rememberedUsername', username);
      }

      uni.showToast({
        title: '登录成功',
        icon: 'success',
      });

      // 跳转到工作台
      setTimeout(() => {
        uni.switchTab({
          url: '/pages/tabbar/workspace',
        });
      }, 1000);
    } else {
      uni.showToast({
        title: '登录失败,请检查用户名和密码',
        icon: 'none',
      });
    }
  } finally {
    loading.value = false;
  }
};

/**
 * 切换记住密码
 */
const handleRememberToggle = () => {
  rememberPassword.value = !rememberPassword.value;
};

/**
 * 忘记密码
 */
const handleForgotPassword = () => {
  uni.showToast({
    title: '请联系管理员重置密码',
    icon: 'none',
  });
};

/**
 * 微信登录
 */
const handleWechatLogin = () => {
  uni.showToast({
    title: '微信登录功能开发中',
    icon: 'none',
  });
};

// 自动填充记住的用户名
const rememberedUsername = uni.getStorageSync('rememberedUsername');
if (rememberedUsername) {
  formData.value.username = rememberedUsername;
  rememberPassword.value = true;
}
</script>

<style lang="scss" scoped>
.login-page {
  min-height: 100vh;
  padding: 0 60rpx;
  display: flex;
  flex-direction: column;
}

.login-logo {
  padding-top: 120rpx;
  padding-bottom: 100rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.logo-image {
  width: 160rpx;
  height: 160rpx;
  margin-bottom: 40rpx;
}

.logo-title {
  font-size: 56rpx;
  font-weight: bold;
  color: #333333;
  margin-bottom: 16rpx;
}

.logo-subtitle {
  font-size: 28rpx;
  color: #999999;
}

.login-form {
  flex: 1;
}

.form-item {
  display: flex;
  align-items: center;
  height: 100rpx;
  background-color: #f5f5f5;
  border-radius: 50rpx;
  padding: 0 30rpx;
  margin-bottom: 30rpx;
}

.form-item-icon {
  margin-right: 20rpx;
}

.form-item-input {
  flex: 1;
  height: 100%;
  font-size: 28rpx;
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 60rpx;
}

.form-remember {
  display: flex;
  align-items: center;
}

.form-remember-text {
  margin-left: 10rpx;
  font-size: 24rpx;
  color: #666666;
}

.form-forgot {
  font-size: 24rpx;
  color: #007aff;
}

.login-button {
  height: 90rpx;
  border-radius: 45rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #ffffff;
  font-size: 32rpx;
  border: none;
}

.login-other {
  padding-bottom: 80rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.login-other-text {
  font-size: 24rpx;
  color: #999999;
  margin-bottom: 30rpx;
}

.login-other-icons {
  display: flex;
  gap: 60rpx;
}

.login-other-icon {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  background-color: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
