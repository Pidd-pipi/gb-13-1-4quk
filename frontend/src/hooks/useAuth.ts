import { computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'

// useAuth：跨页面复用的登录态 hook（登录/登出/权限判断）
export function useAuth() {
  const auth = useAuthStore()
  const isLoggedIn = computed(() => auth.isLoggedIn)
  const user = computed(() => auth.user)
  const isAdmin = computed(() => auth.user?.role === 'admin')

  function login(account: string, password: string) {
    return auth.login(account, password)
  }

  function logout() {
    auth.logout()
  }

  return { isLoggedIn, user, isAdmin, login, logout }
}
