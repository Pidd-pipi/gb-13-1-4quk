<template>
  <div class="login-page">
    <div class="login-header">
      <div class="logo"><van-icon name="book-o" size="42" color="#1989fa" /></div>
      <div class="slogan">校园二手书交易平台</div>
      <div class="sub">教材循环 · 省钱环保</div>
    </div>
    <van-form @submit="onLogin">
      <van-cell-group inset>
        <van-field v-model="account" label="学号/邮箱" placeholder="学号或学校邮箱" :rules="[{ required: true, message: '请输入学号或邮箱' }]" />
        <van-field v-model="password" type="password" label="密码" placeholder="密码" :rules="[{ required: true, message: '请输入密码' }]" />
      </van-cell-group>
      <div class="actions">
        <van-button round block type="primary" native-type="submit" :loading="loading">登录</van-button>
        <div class="links">
          <span @click="$router.push('/register')">注册新账号</span>
          <span class="demo">演示账号：zhang@campusbooks.local / student123</span>
        </div>
      </div>
    </van-form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showSuccessToast } from 'vant'
import { useAuth } from '@/hooks/useAuth'

const route = useRoute()
const router = useRouter()
const { login } = useAuth()
const account = ref('')
const password = ref('')
const loading = ref(false)

async function onLogin() {
  loading.value = true
  try {
    await login(account.value, password.value)
    showSuccessToast('登录成功')
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch {
    /* toast handled */
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: #fff;
  padding-top: 80px;
}
.login-header {
  text-align: center;
  margin-bottom: 32px;
}
.logo {
  margin-bottom: 12px;
}
.slogan {
  font-size: 22px;
  font-weight: 700;
}
.sub {
  color: #969799;
  font-size: 13px;
  margin-top: 6px;
}
.actions {
  padding: 24px 16px;
}
.links {
  margin-top: 14px;
  text-align: center;
  color: #1989fa;
  font-size: 14px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.demo {
  color: #969799;
  font-size: 12px;
}
</style>
