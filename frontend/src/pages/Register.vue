<template>
  <div class="register-page">
    <van-nav-bar title="注册学生账号" left-arrow fixed placeholder @click-left="$router.back()" />
    <van-form @submit="onRegister">
      <van-cell-group inset>
        <van-field v-model="form.student_no" label="学号" placeholder="学号" :rules="[{ required: true, message: '请输入学号' }]" />
        <van-field v-model="form.email" label="学校邮箱" placeholder="学校邮箱" :rules="[{ required: true, message: '请输入学校邮箱' }]" />
        <van-field v-model="form.password" type="password" label="密码" placeholder="至少 6 位" :rules="[{ required: true, message: '请输入密码' }]" />
        <van-field v-model="form.code" label="验证码" placeholder="6 位验证码" :rules="[{ required: true, message: '请输入验证码' }]">
          <template #button>
            <van-button size="small" type="primary" plain :disabled="countdown > 0" @click="onSendCode">
              {{ countdown > 0 ? `${countdown}s` : '获取验证码' }}
            </van-button>
          </template>
        </van-field>
      </van-cell-group>
      <div class="actions">
        <van-button round block type="primary" native-type="submit" :loading="loading">注册</van-button>
      </div>
    </van-form>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { showSuccessToast, showToast } from 'vant'
import { sendCode, register } from '@/api/auth'
import { setToken } from '@/utils/request'
import { useAuthStore } from '@/stores/authStore'

const router = useRouter()
const auth = useAuthStore()
const loading = ref(false)
const countdown = ref(0)
let timer: number | undefined

const form = reactive({ student_no: '', email: '', password: '', code: '' })

async function onSendCode() {
  if (!form.email) {
    showToast('请先输入学校邮箱')
    return
  }
  try {
    const res = await sendCode(form.email)
    if (res.dev_code) {
      form.code = res.dev_code
      showToast(`验证码已发送：${res.dev_code}（开发环境）`)
    } else {
      showToast('验证码已发送')
    }
    countdown.value = 60
    timer = window.setInterval(() => {
      countdown.value -= 1
      if (countdown.value <= 0) {
        window.clearInterval(timer)
      }
    }, 1000)
  } catch {
    /* toast */
  }
}

async function onRegister() {
  loading.value = true
  try {
    const res = await register(form)
    setToken(res.token)
    auth.fetchMe()
    showSuccessToast('注册成功')
    router.push('/profile')
  } catch {
    /* toast */
  } finally {
    loading.value = false
  }
}

onUnmounted(() => {
  if (timer) window.clearInterval(timer)
})
</script>

<style scoped>
.actions {
  padding: 24px 16px;
}
</style>
