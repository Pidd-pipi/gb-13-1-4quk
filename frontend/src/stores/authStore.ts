import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import * as authApi from '@/api/auth'
import * as userApi from '@/api/user'
import { clearToken, getToken, setToken } from '@/utils/request'
import type { User } from '@/types'

const USER_KEY = 'campusbooks_user'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(getToken())
  const user = ref<User | null>(readStoredUser())

  const isLoggedIn = computed(() => !!token.value)

  function persistUser(u: User | null) {
    user.value = u
    if (u) {
      localStorage.setItem(USER_KEY, JSON.stringify(u))
    } else {
      localStorage.removeItem(USER_KEY)
    }
  }

  function readStoredUser(): User | null {
    try {
      const raw = localStorage.getItem(USER_KEY)
      return raw ? (JSON.parse(raw) as User) : null
    } catch {
      return null
    }
  }

  async function login(account: string, password: string) {
    const res = await authApi.login({ account, password })
    token.value = res.token
    setToken(res.token)
    persistUser(res.user)
    return res.user
  }

  async function fetchMe() {
    const me = await userApi.getMe()
    persistUser(me)
    return me
  }

  function logout() {
    token.value = ''
    clearToken()
    persistUser(null)
  }

  return { token, user, isLoggedIn, login, fetchMe, logout }
})
