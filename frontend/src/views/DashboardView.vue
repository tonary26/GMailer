<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import AppLogo from '@/components/AppLogo.vue'
import AppIcon from '@/components/AppIcon.vue'
import api from '@/services/api'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const menuOpen = ref(false)
const composerOpen = ref(false)
const activeSection = ref('overview')
const search = ref('')
const contacts = ref([])
const contactsLoading = ref(true)
const contactsError = ref('')
const menuButton = ref(null)
const sidebar = ref(null)
const composer = ref(null)
const composerFirstField = ref(null)
const composerTrigger = ref(null)

const navigation = computed(() => [
  { id: 'overview', label: 'Обзор', icon: 'grid' },
  { id: 'campaigns', label: 'Рассылки', icon: 'send' },
  {
    id: 'contacts',
    label: 'Контакты',
    icon: 'users',
    count: contactsLoading.value || contactsError.value ? null : contacts.value.length,
  },
  { id: 'analytics', label: 'Аналитика', icon: 'chart' },
])

const filteredContacts = computed(() => {
  const query = search.value.trim().toLocaleLowerCase('ru-RU')
  if (!query) return contacts.value

  return contacts.value.filter((contact) =>
    `${contact.name ?? ''} ${contact.email ?? ''}`.toLocaleLowerCase('ru-RU').includes(query),
  )
})

const initials = computed(() => (auth.user?.email?.[0] || 'G').toUpperCase())
const sectionTitle = computed(
  () => navigation.value.find((item) => item.id === activeSection.value)?.label || 'Обзор',
)

const pluralize = (value, forms) => {
  const modulo100 = Math.abs(value) % 100
  const modulo10 = modulo100 % 10
  if (modulo100 > 10 && modulo100 < 20) return forms[2]
  if (modulo10 > 1 && modulo10 < 5) return forms[1]
  if (modulo10 === 1) return forms[0]
  return forms[2]
}

const contactsSummary = computed(() => {
  if (contactsLoading.value) return 'Загружаем список…'
  if (contactsError.value) return 'Данные недоступны'

  const count = contacts.value.length
  return `${count.toLocaleString('ru-RU')} ${pluralize(count, ['контакт', 'контакта', 'контактов'])}`
})

const formatContactDate = (value) => {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return new Intl.DateTimeFormat('ru-RU', { dateStyle: 'medium' }).format(date)
}

const contactInitial = (contact) => (contact.name?.[0] || contact.email?.[0] || '?').toUpperCase()

const loadContacts = async () => {
  contactsLoading.value = true
  contactsError.value = ''
  try {
    const { data } = await api.get('/contact/list')
    contacts.value = Array.isArray(data.contacts) ? data.contacts : []
  } catch {
    contacts.value = []
    contactsError.value =
      'Не удалось загрузить контакты. Проверьте соединение и попробуйте ещё раз.'
  } finally {
    contactsLoading.value = false
  }
}

const selectSection = (id) => {
  activeSection.value = id
  menuOpen.value = false
  search.value = ''
}

const openComposer = (event) => {
  composerTrigger.value = event?.currentTarget || document.activeElement
  composerOpen.value = true
  nextTick(() => composerFirstField.value?.focus())
}

const closeComposer = () => {
  composerOpen.value = false
  nextTick(() => composerTrigger.value?.focus())
}

const toggleMenu = () => {
  menuOpen.value = !menuOpen.value
  if (menuOpen.value) nextTick(() => sidebar.value?.focus())
}

const closeMenu = (restoreFocus = false) => {
  menuOpen.value = false
  if (restoreFocus) nextTick(() => menuButton.value?.focus())
}

const handleKeydown = (event) => {
  if (event.key === 'Escape') {
    if (composerOpen.value) closeComposer()
    else if (menuOpen.value) closeMenu(true)
    return
  }
  if (event.key !== 'Tab' || !composerOpen.value || !composer.value) return
  const focusable = [
    ...composer.value.querySelectorAll('button:not(:disabled), input:not(:disabled)'),
  ]
  if (!focusable.length) return
  const first = focusable[0]
  const last = focusable[focusable.length - 1]
  if (event.shiftKey && document.activeElement === first) {
    event.preventDefault()
    last.focus()
  } else if (!event.shiftKey && document.activeElement === last) {
    event.preventDefault()
    first.focus()
  }
}

document.addEventListener('keydown', handleKeydown)
onBeforeUnmount(() => document.removeEventListener('keydown', handleKeydown))
onMounted(loadContacts)

const logout = () => {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <div class="dashboard-shell">
    <header class="topbar">
      <div class="topbar__brand">
        <button
          ref="menuButton"
          class="icon-button mobile-menu"
          type="button"
          :aria-label="menuOpen ? 'Закрыть меню' : 'Открыть меню'"
          aria-controls="dashboard-navigation"
          :aria-expanded="menuOpen"
          @click="toggleMenu"
        >
          <AppIcon name="menu" :size="23" />
        </button>
        <AppLogo />
      </div>
      <label v-if="activeSection === 'contacts'" class="search-box">
        <AppIcon name="search" :size="19" />
        <span class="sr-only">Поиск по контактам</span>
        <input v-model="search" type="search" placeholder="Поиск по контактам" />
        <kbd>/</kbd>
      </label>
      <div v-else aria-hidden="true"></div>
      <div class="account">
        <div class="account__copy">
          <strong>{{ auth.user?.email || 'Аккаунт' }}</strong
          ><span>Рабочее пространство</span>
        </div>
        <span class="avatar" :title="auth.user?.email">{{ initials }}</span>
      </div>
    </header>

    <div v-if="menuOpen" class="nav-scrim" @click="closeMenu(true)"></div>
    <aside
      id="dashboard-navigation"
      ref="sidebar"
      class="sidebar"
      :class="{ 'sidebar--open': menuOpen }"
      tabindex="-1"
    >
      <button class="compose-button" type="button" @click="openComposer">
        <AppIcon name="plus" :size="23" />
        Новая рассылка
      </button>
      <nav aria-label="Основная навигация">
        <button
          v-for="item in navigation"
          :key="item.id"
          class="nav-item"
          :class="{ 'nav-item--active': activeSection === item.id }"
          type="button"
          @click="selectSection(item.id)"
        >
          <AppIcon :name="item.icon" />
          <span>{{ item.label }}</span>
          <small v-if="item.count">{{ item.count.toLocaleString('ru-RU') }}</small>
        </button>
      </nav>
      <div class="sidebar__bottom">
        <button class="nav-item" type="button" @click="selectSection('settings')">
          <AppIcon name="settings" /><span>Настройки</span>
        </button>
        <button class="nav-item" type="button" @click="logout">
          <AppIcon name="logout" /><span>Выйти</span>
        </button>
      </div>
    </aside>

    <main class="workspace">
      <div class="workspace__heading">
        <div>
          <h1>{{ sectionTitle }}</h1>
          <p v-if="activeSection === 'overview'">Здесь появятся данные о ваших рассылках.</p>
        </div>
      </div>

      <section
        v-if="activeSection === 'overview'"
        class="content-panel full-panel empty-state empty-state--large"
      >
        <span class="empty-state__icon"><AppIcon name="send" :size="30" /></span>
        <h2>Данных о рассылках пока нет</h2>
        <p>После подключения рассылок здесь появятся их актуальные статусы и результаты.</p>
      </section>

      <section v-else-if="activeSection === 'campaigns'" class="content-panel full-panel">
        <div class="panel-heading">
          <div>
            <h2>Все рассылки</h2>
            <p>Данные появятся после подключения API рассылок</p>
          </div>
          <button class="small-primary" type="button" @click="openComposer">
            <AppIcon name="plus" :size="17" /> Создать
          </button>
        </div>
        <div class="empty-state">
          <AppIcon name="send" :size="26" />
          <h3>Рассылок пока нет</h3>
          <p>Здесь будут отображаться реальные рассылки, когда для них появится API.</p>
        </div>
      </section>

      <section v-else-if="activeSection === 'contacts'" class="content-panel full-panel">
        <div class="panel-heading">
          <div>
            <h2>Контакты</h2>
            <p>{{ contactsSummary }}</p>
          </div>
        </div>
        <div v-if="contactsLoading" class="empty-state" aria-live="polite">
          <span class="loading-indicator" aria-hidden="true"></span>
          <h3>Загружаем контакты</h3>
        </div>
        <div v-else-if="contactsError" class="empty-state" role="alert">
          <AppIcon name="clock" :size="26" />
          <h3>Контакты недоступны</h3>
          <p>{{ contactsError }}</p>
          <button class="retry-button" type="button" @click="loadContacts">Повторить</button>
        </div>
        <div v-else-if="filteredContacts.length" class="contacts-list">
          <article v-for="contact in filteredContacts" :key="contact.id" class="contact-row">
            <span class="contact-avatar">{{ contactInitial(contact) }}</span>
            <div>
              <strong>{{ contact.name }}</strong
              ><span>{{ contact.email }}</span>
            </div>
            <time v-if="formatContactDate(contact.created_at)" :datetime="contact.created_at">{{
              formatContactDate(contact.created_at)
            }}</time>
          </article>
        </div>
        <div v-else class="empty-state">
          <AppIcon :name="search ? 'search' : 'users'" :size="26" />
          <h3>{{ search ? 'Ничего не найдено' : 'Контактов пока нет' }}</h3>
          <p>
            {{
              search
                ? 'Попробуйте изменить поисковый запрос.'
                : 'Добавленные через API контакты появятся здесь.'
            }}
          </p>
        </div>
      </section>

      <section v-else class="content-panel full-panel empty-state empty-state--large">
        <span class="empty-state__icon"
          ><AppIcon :name="activeSection === 'analytics' ? 'chart' : 'settings'" :size="30"
        /></span>
        <h2>
          {{
            activeSection === 'analytics'
              ? 'Аналитика появится здесь'
              : 'Настройки рабочего пространства'
          }}
        </h2>
        <p>
          {{
            activeSection === 'analytics'
              ? 'Подключите эндпоинты статистики, чтобы видеть реальные открытия и переходы.'
              : 'Раздел готов к подключению параметров отправителя и SMTP.'
          }}
        </p>
      </section>
    </main>

    <div v-if="composerOpen" class="composer-scrim" @click.self="closeComposer">
      <section
        ref="composer"
        class="composer"
        role="dialog"
        aria-modal="true"
        aria-labelledby="composer-title"
      >
        <div class="composer__header">
          <div>
            <h2 id="composer-title">Новая рассылка</h2>
            <p>Черновик сохраняется только в интерфейсе</p>
          </div>
          <button
            class="icon-button"
            type="button"
            aria-label="Закрыть окно"
            @click="closeComposer"
          >
            <AppIcon name="close" />
          </button>
        </div>
        <div class="field">
          <label for="campaign-title">Название</label>
          <div class="field__control">
            <input
              id="campaign-title"
              ref="composerFirstField"
              placeholder="Например, Новости октября"
            />
          </div>
        </div>
        <div class="field">
          <label for="campaign-subject">Тема письма</label>
          <div class="field__control">
            <input id="campaign-subject" placeholder="Что увидит получатель" />
          </div>
        </div>
        <div class="composer__notice">
          Для сохранения рассылки нужен API-эндпоинт. Сейчас это демонстрация будущего сценария.
        </div>
        <div class="composer__actions">
          <button type="button" @click="closeComposer">Отмена</button
          ><button type="button" disabled>Сохранить черновик</button>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}
.dashboard-shell {
  min-height: 100vh;
  background: var(--canvas);
}
.topbar {
  position: fixed;
  z-index: 30;
  inset: 0 0 auto 0;
  height: 72px;
  display: grid;
  grid-template-columns: 238px minmax(260px, 680px) 1fr;
  align-items: center;
  gap: 22px;
  padding: 0 24px;
  background: var(--canvas);
}
.topbar__brand {
  display: flex;
  align-items: center;
  gap: 8px;
}
.mobile-menu {
  display: none !important;
}
.icon-button {
  display: grid;
  width: 40px;
  height: 40px;
  place-items: center;
  border: 0;
  border-radius: 50%;
  color: #4f545a;
  background: transparent;
  cursor: pointer;
}
.icon-button:hover {
  background: #e8ecf2;
}
.search-box {
  display: flex;
  height: 48px;
  align-items: center;
  gap: 13px;
  padding: 0 16px;
  border-radius: 24px;
  color: #5f6368;
  background: #eaf1fb;
  transition:
    background 0.16s ease,
    box-shadow 0.16s ease;
}
.search-box:focus-within {
  background: white;
  box-shadow: var(--shadow);
}
.search-box input {
  min-width: 0;
  flex: 1;
  border: 0;
  outline: 0;
  color: var(--ink);
  background: transparent;
}
.search-box input::placeholder {
  color: #5f6368;
}
.search-box kbd {
  padding: 2px 7px;
  border: 1px solid #c4cad3;
  border-radius: 5px;
  color: #5f6368;
  background: rgba(255, 255, 255, 0.55);
  font:
    11px/1.5 'Manrope',
    sans-serif;
}
.account {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
}
.account__copy {
  display: grid;
  justify-items: end;
  font-size: 12px;
}
.account__copy strong {
  max-width: 210px;
  overflow: hidden;
  text-overflow: ellipsis;
}
.account__copy span {
  margin-top: 2px;
  color: var(--muted);
  font-size: 11px;
}
.avatar {
  display: grid;
  flex: 0 0 38px;
  height: 38px;
  place-items: center;
  border-radius: 50%;
  color: white;
  background: #7c4dff;
  font-size: 14px;
  font-weight: 700;
}
.sidebar {
  position: fixed;
  z-index: 25;
  inset: 72px auto 0 0;
  width: 256px;
  display: flex;
  flex-direction: column;
  padding: 14px 12px 20px;
  background: var(--canvas);
}
.compose-button {
  display: flex;
  width: 190px;
  min-height: 56px;
  align-items: center;
  gap: 13px;
  margin: 0 0 20px 5px;
  padding: 0 20px;
  border: 0;
  border-radius: 17px;
  color: #001d35;
  background: var(--surface-blue);
  box-shadow: 0 1px 2px rgba(60, 64, 67, 0.08);
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  transition:
    box-shadow 0.16s ease,
    transform 0.16s ease;
}
.compose-button:hover {
  box-shadow: 0 5px 12px rgba(60, 64, 67, 0.17);
  transform: translateY(-1px);
}
.sidebar nav {
  display: grid;
  gap: 2px;
}
.nav-item {
  display: grid;
  width: 100%;
  min-height: 43px;
  grid-template-columns: 22px 1fr auto;
  align-items: center;
  gap: 13px;
  padding: 0 16px;
  border: 0;
  border-radius: 0 22px 22px 0;
  color: #3c4043;
  background: transparent;
  font-size: 13px;
  text-align: left;
  cursor: pointer;
}
.nav-item:hover {
  background: #e9edf3;
}
.nav-item--active {
  color: #001d35;
  background: #d3e3fd;
  font-weight: 700;
}
.nav-item--active:hover {
  background: #d3e3fd;
}
.nav-item small {
  font-size: 11px;
  font-weight: 650;
}
.sidebar__bottom {
  display: grid;
  gap: 2px;
  margin-top: auto;
  padding-top: 14px;
  border-top: 1px solid var(--line);
}
.workspace {
  min-height: 100vh;
  margin-left: 256px;
  padding: 94px 18px 28px 0;
}
.workspace::before {
  content: '';
  position: fixed;
  z-index: 0;
  inset: 72px 18px 18px 256px;
  border-radius: 20px;
  background: var(--surface);
}
.workspace > * {
  position: relative;
  z-index: 1;
}
.workspace__heading {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 20px;
  padding: 20px 28px 22px;
}
.workspace__heading h1 {
  margin: 0 0 6px;
  font-size: 26px;
  font-weight: 650;
  letter-spacing: -0.03em;
}
.workspace__heading p {
  margin: 0;
  color: var(--muted);
  font-size: 13px;
}
.content-panel {
  border: 1px solid var(--line);
  border-radius: 14px;
  background: white;
}
.panel-heading {
  min-height: 82px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  padding: 18px 20px;
  border-bottom: 1px solid var(--line);
}
.panel-heading h2 {
  margin: 0 0 5px;
  font-size: 15px;
  font-weight: 700;
}
.panel-heading p {
  margin: 0;
  color: var(--muted);
  font-size: 11px;
}
.full-panel {
  margin: 0 28px 28px;
  overflow: hidden;
}
.small-primary {
  display: flex;
  min-height: 36px;
  align-items: center;
  gap: 7px;
  padding: 0 14px;
  border: 0;
  border-radius: 18px;
  color: white;
  background: var(--blue);
  font-size: 11px;
  font-weight: 700;
  cursor: pointer;
}
.retry-button {
  min-height: 36px;
  margin-top: 18px;
  padding: 0 16px;
  border: 1px solid var(--line-strong);
  border-radius: 18px;
  color: var(--blue);
  background: white;
  font-size: 11px;
  font-weight: 700;
  cursor: pointer;
}
.contacts-list {
  display: grid;
}
.contact-row {
  display: grid;
  min-height: 68px;
  grid-template-columns: 40px 1fr auto;
  align-items: center;
  gap: 13px;
  padding: 9px 20px;
  border-bottom: 1px solid #edf0f2;
}
.contact-row:last-child {
  border-bottom: 0;
}
.contact-row:hover {
  background: #f8fafd;
}
.contact-avatar {
  display: grid;
  width: 36px;
  height: 36px;
  place-items: center;
  border-radius: 50%;
  color: #174ea6;
  background: #d2e3fc;
  font-size: 12px;
  font-weight: 700;
}
.contact-row div {
  display: grid;
  gap: 4px;
}
.contact-row strong {
  font-size: 12px;
}
.contact-row span,
.contact-row time {
  color: var(--muted);
  font-size: 11px;
}
.empty-state {
  display: grid;
  justify-items: center;
  padding: 50px 20px;
  color: var(--muted);
  text-align: center;
}
.empty-state h3,
.empty-state h2 {
  margin: 14px 0 6px;
  color: var(--ink);
  font-size: 15px;
}
.empty-state p {
  max-width: 480px;
  margin: 0;
  font-size: 12px;
  line-height: 1.6;
}
.empty-state--large {
  min-height: 400px;
  align-content: center;
}
.empty-state__icon {
  display: grid;
  width: 58px;
  height: 58px;
  place-items: center;
  border-radius: 18px;
  color: var(--blue);
  background: #eaf2ff;
}
.loading-indicator {
  width: 24px;
  height: 24px;
  border: 2px solid #d2e3fc;
  border-top-color: var(--blue);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
.composer-scrim {
  position: fixed;
  z-index: 50;
  inset: 0;
  display: grid;
  place-items: center;
  padding: 20px;
  background: rgba(32, 33, 36, 0.35);
}
.composer {
  width: min(520px, 100%);
  display: grid;
  gap: 18px;
  padding: 24px;
  border-radius: 18px;
  background: white;
  box-shadow: 0 20px 60px rgba(32, 33, 36, 0.25);
  animation: composerIn 0.24s cubic-bezier(0.2, 0.8, 0.2, 1);
}
@keyframes composerIn {
  from {
    opacity: 0.5;
    transform: translateY(12px) scale(0.98);
    filter: blur(2px);
  }
}
.composer__header {
  display: flex;
  justify-content: space-between;
  gap: 20px;
}
.composer__header h2 {
  margin: 0 0 5px;
  font-size: 19px;
}
.composer__header p {
  margin: 0;
  color: var(--muted);
  font-size: 11px;
}
.composer__notice {
  padding: 12px 14px;
  border-radius: 9px;
  color: #6b4e10;
  background: var(--amber-bg);
  font-size: 11px;
  line-height: 1.5;
}
.composer__actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
.composer__actions button {
  min-height: 38px;
  padding: 0 15px;
  border: 1px solid var(--line-strong);
  border-radius: 19px;
  background: white;
  font-size: 11px;
  font-weight: 700;
  cursor: pointer;
}
.composer__actions button:disabled {
  border-color: transparent;
  color: #8a8f96;
  background: #e7e9ec;
  cursor: not-allowed;
}
.nav-scrim {
  display: none;
}

@media (max-width: 1100px) {
  .topbar {
    grid-template-columns: 220px minmax(240px, 1fr) auto;
  }
  .account__copy {
    display: none;
  }
}

@media (max-width: 760px) {
  .topbar {
    height: 66px;
    grid-template-columns: auto 1fr auto;
    gap: 10px;
    padding: 0 14px;
  }
  .topbar__brand :deep(.brand__name) {
    display: none;
  }
  .mobile-menu {
    display: grid !important;
  }
  .search-box {
    height: 42px;
    padding: 0 13px;
  }
  .search-box kbd {
    display: none;
  }
  .account .avatar {
    width: 34px;
    height: 34px;
    flex-basis: 34px;
  }
  .sidebar {
    inset: 66px auto 0 0;
    width: min(282px, 86vw);
    padding-top: 18px;
    background: white;
    box-shadow: 10px 0 30px rgba(32, 33, 36, 0.14);
    transform: translateX(-105%);
    transition: transform 0.22s cubic-bezier(0.2, 0.8, 0.2, 1);
  }
  .sidebar--open {
    transform: translateX(0);
  }
  .nav-scrim {
    position: fixed;
    z-index: 20;
    inset: 66px 0 0 0;
    display: block;
    background: rgba(32, 33, 36, 0.3);
  }
  .workspace {
    margin-left: 0;
    padding: 79px 0 12px;
  }
  .workspace::before {
    inset: 66px 0 0;
    border-radius: 18px 18px 0 0;
  }
  .workspace__heading {
    padding: 20px 18px;
  }
  .workspace__heading h1 {
    font-size: 23px;
  }
  .workspace__heading p {
    max-width: 250px;
    line-height: 1.45;
  }
  .full-panel {
    margin: 0 14px 20px;
  }
  .panel-heading {
    padding: 16px;
  }
  .composer {
    padding: 20px;
  }
}

@media (max-width: 430px) {
  .topbar__brand :deep(.brand) {
    display: none;
  }
  .workspace__heading {
    align-items: center;
  }
  .workspace__heading p {
    display: none;
  }
}
</style>
