<script setup>
import { reactive, ref, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import AuthLayout from '@/components/AuthLayout.vue'
import AppIcon from '@/components/AppIcon.vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const showPassword = ref(false)
const form = reactive({ email: '', password: '', confirmPassword: '' })
const errors = reactive({ email: '', password: '', confirmPassword: '' })

onBeforeUnmount(auth.clearError)

const submit = async () => {
  errors.email = ''
  errors.password = ''
  errors.confirmPassword = ''
  if (!form.email) errors.email = 'Укажите электронную почту.'
  else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email))
    errors.email = 'Введите адрес в формате name@example.com.'
  if (!form.password) errors.password = 'Придумайте пароль.'
  else if (form.password.length < 8) errors.password = 'Пароль должен содержать минимум 8 символов.'
  if (!form.confirmPassword) errors.confirmPassword = 'Повторите пароль.'
  else if (form.password !== form.confirmPassword) errors.confirmPassword = 'Пароли не совпадают.'
  if (Object.values(errors).some(Boolean)) return
  const success = await auth.register({ email: form.email, password: form.password })
  if (success) router.push('/dashboard')
}
</script>

<template>
  <AuthLayout title="Создайте аккаунт">
    <form class="auth-form" novalidate @submit.prevent="submit">
      <div class="field">
        <label for="register-email">Электронная почта</label>
        <div class="field__control">
          <input
            id="register-email"
            v-model.trim="form.email"
            type="email"
            autocomplete="email"
            placeholder="name@example.com"
            :disabled="auth.loading"
            :aria-invalid="Boolean(errors.email)"
            :aria-describedby="errors.email ? 'register-email-error' : undefined"
          />
        </div>
        <p v-if="errors.email" id="register-email-error" class="field__error">{{ errors.email }}</p>
      </div>
      <div class="field">
        <label for="register-password">Пароль</label>
        <div class="field__control">
          <input
            id="register-password"
            v-model="form.password"
            :type="showPassword ? 'text' : 'password'"
            autocomplete="new-password"
            placeholder="Минимум 8 символов"
            :disabled="auth.loading"
            :aria-invalid="Boolean(errors.password)"
            :aria-describedby="
              errors.password
                ? 'register-password-hint register-password-error'
                : 'register-password-hint'
            "
          />
          <button
            class="field__action"
            type="button"
            :aria-label="showPassword ? 'Скрыть пароль' : 'Показать пароль'"
            @click="showPassword = !showPassword"
          >
            <AppIcon :name="showPassword ? 'eye-off' : 'eye'" />
          </button>
        </div>
        <p v-if="errors.password" id="register-password-error" class="field__error">
          {{ errors.password }}
        </p>
      </div>
      <div class="field">
        <label for="register-confirm">Повторите пароль</label>
        <div class="field__control">
          <input
            id="register-confirm"
            v-model="form.confirmPassword"
            :type="showPassword ? 'text' : 'password'"
            autocomplete="new-password"
            placeholder="Введите пароль ещё раз"
            :disabled="auth.loading"
            :aria-invalid="Boolean(errors.confirmPassword)"
            :aria-describedby="errors.confirmPassword ? 'register-confirm-error' : undefined"
          />
        </div>
        <p v-if="errors.confirmPassword" id="register-confirm-error" class="field__error">
          {{ errors.confirmPassword }}
        </p>
      </div>
      <p id="register-password-hint" class="field__hint">Используйте не менее 8 символов.</p>
      <p v-if="auth.error" class="form-error" role="alert">{{ auth.error }}</p>
      <button class="primary-button" type="submit" :disabled="auth.loading">
        <span v-if="auth.loading" class="button-spinner"></span>
        {{ auth.loading ? 'Создаём аккаунт…' : 'Создать аккаунт' }}
        <AppIcon v-if="!auth.loading" name="arrow" />
      </button>
      <p class="auth-switch">Уже есть аккаунт? <RouterLink to="/login">Войти</RouterLink></p>
    </form>
  </AuthLayout>
</template>

<style scoped>
.auth-form {
  display: grid;
  gap: 19px;
}
.field {
  display: grid;
  gap: 8px;
}
.field label {
  color: #303134;
  font-size: 13px;
  font-weight: 650;
}
.field__control {
  position: relative;
}
.field input {
  width: 100%;
  height: 52px;
  padding: 0 46px 0 15px;
  border: 1px solid var(--line-strong);
  border-radius: 9px;
  color: var(--ink);
  background: var(--surface);
  caret-color: var(--blue);
  transition:
    border-color 0.16s ease,
    box-shadow 0.16s ease;
}
.field input::placeholder {
  color: #5f6368;
}
.field input:hover {
  border-color: #8d949d;
}
.field input:focus {
  border-color: var(--blue);
  outline: none;
  box-shadow: 0 0 0 1px var(--blue);
}
.field__action {
  position: absolute;
  top: 50%;
  right: 8px;
  display: grid;
  width: 38px;
  height: 38px;
  place-items: center;
  border: 0;
  border-radius: 50%;
  color: #5f6368;
  background: transparent;
  cursor: pointer;
  transform: translateY(-50%);
}
.field__action:hover {
  background: var(--surface-soft);
}
.field__hint {
  margin: -10px 0 0;
  color: #5f6368;
  font-size: 12px;
  line-height: 1.5;
}
.field__error {
  margin: 0;
  color: var(--red);
  font-size: 12px;
  line-height: 1.45;
}
.field input[aria-invalid='true'] {
  border-color: var(--red);
}
.form-error {
  display: flex;
  align-items: flex-start;
  gap: 9px;
  margin: 0;
  padding: 12px 14px;
  border-radius: 9px;
  color: #8c1d18;
  background: var(--red-bg);
  font-size: 13px;
  line-height: 1.45;
}
.form-error::before {
  content: '!';
  display: grid;
  flex: 0 0 19px;
  height: 19px;
  place-items: center;
  border: 1px solid currentColor;
  border-radius: 50%;
  font-size: 11px;
  font-weight: 700;
}
.primary-button {
  display: inline-flex;
  min-height: 50px;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border: 0;
  border-radius: 25px;
  color: white;
  background: var(--blue);
  box-shadow: 0 4px 11px rgba(11, 87, 208, 0.2);
  font-weight: 700;
  cursor: pointer;
  transition:
    background 0.16s ease,
    box-shadow 0.16s ease,
    transform 0.16s ease;
}
.primary-button:hover {
  background: var(--blue-hover);
  box-shadow: 0 7px 16px rgba(11, 87, 208, 0.24);
  transform: translateY(-1px);
}
.primary-button:disabled {
  color: #7a7f86;
  background: #e4e7eb;
  box-shadow: none;
  cursor: not-allowed;
  transform: none;
}
.auth-switch {
  margin: 7px 0 0;
  color: var(--muted);
  font-size: 14px;
  text-align: center;
}
.auth-switch a {
  font-weight: 700;
  text-decoration: none;
}
.auth-switch a:hover {
  text-decoration: underline;
}
.button-spinner {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, 0.45);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.65s linear infinite;
}
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
