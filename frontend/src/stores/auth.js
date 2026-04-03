import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const isAuthenticated = ref(!!token.value)

  async function login(password) {
    const res = await api.login(password)
    token.value = res.data.token
    localStorage.setItem('token', token.value)
    isAuthenticated.value = true
    router.push('/dashboard')
  }

  function logout() {
    token.value = ''
    localStorage.removeItem('token')
    isAuthenticated.value = false
    router.push('/login')
  }

  return { token, isAuthenticated, login, logout }
})
