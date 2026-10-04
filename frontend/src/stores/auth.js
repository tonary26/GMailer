import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import api from '@/services/api'

const parseStoredUser = () => {
  try {
    return JSON.parse(localStorage.getItem('gmailer_user'))
  } catch {
    return null
  }
}

const getErrorMessage = (error) => {
  if (error.code === 'ECONNABORTED') return 'Сервер отвечает слишком долго. Попробуйте ещё раз.'
  if (!error.response) return 'Не удалось связаться с сервером. Проверьте, запущен ли API.'
  return error.response.data?.message || 'Что-то пошло не так. Попробуйте ещё раз.'
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('gmailer_token'))
  const user = ref(parseStoredUser())
  const loading = ref(false)
  const error = ref('')
  const isAuthenticated = computed(() => Boolean(token.value))

  const persistSession = (payload) => {
    const safeUser = payload.user
      ? {
          id: payload.user.id,
          email: payload.user.email,
          created_at: payload.user.created_at,
        }
      : null
    token.value = payload.token
    user.value = safeUser
    localStorage.setItem('gmailer_token', payload.token)
    localStorage.setItem('gmailer_user', JSON.stringify(safeUser))
  }

  const login = async (credentials) => {
    loading.value = true
    error.value = ''
    try {
      const { data } = await api.post('/auth/login', credentials)
      persistSession(data)
      return true
    } catch (requestError) {
      error.value = getErrorMessage(requestError)
      return false
    } finally {
      loading.value = false
    }
  }

  const register = async (credentials) => {
    loading.value = true
    error.value = ''
    try {
      await api.post('/auth/register', credentials)
      const { data } = await api.post('/auth/login', credentials)
      persistSession(data)
      return true
    } catch (requestError) {
      error.value = getErrorMessage(requestError)
      return false
    } finally {
      loading.value = false
    }
  }

  const logout = () => {
    token.value = null
    user.value = null
    error.value = ''
    localStorage.removeItem('gmailer_token')
    localStorage.removeItem('gmailer_user')
  }

  const clearError = () => {
    error.value = ''
  }
  return { token, user, loading, error, isAuthenticated, login, register, logout, clearError }
})
