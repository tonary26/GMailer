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
const contactModalOpen = ref(false)
const activeSection = ref('overview')
const search = ref('')
const contacts = ref([])
const contactsLoading = ref(true)
const contactsError = ref('')
const mailings = ref([])
const mailingsLoading = ref(true)
const mailingsError = ref('')
const mailingProgress = ref({})
const progressError = ref('')
const selectedMailing = ref(null)
const mailingDetailsOpen = ref(false)
const mailingDetailsLoading = ref(false)
const mailingDetailsError = ref('')
const startingMailingId = ref(null)
const menuButton = ref(null)
const sidebar = ref(null)
const composer = ref(null)
const composerFirstField = ref(null)
const composerTrigger = ref(null)
const mailingDetailsDialog = ref(null)
const mailingDetailsTrigger = ref(null)
const campaignTitle = ref('')
const campaignSubject = ref('')
const campaignBody = ref('')
const campaignTitleError = ref('')
const campaignSubjectError = ref('')
const campaignBodyError = ref('')
const campaignSubmitError = ref('')
const campaignSubmitting = ref(false)
const contactDialog = ref(null)
const contactNameField = ref(null)
const contactModalTrigger = ref(null)
const contactName = ref('')
const contactEmail = ref('')
const contactNameError = ref('')
const contactEmailError = ref('')
const contactSubmitError = ref('')
const contactSubmitting = ref(false)
let progressTimer = null

const navigation = computed(() => [
  { id: 'overview', label: 'Обзор', icon: 'grid' },
  {
    id: 'campaigns',
    label: 'Рассылки',
    icon: 'send',
    count: mailingsLoading.value || mailingsError.value ? null : mailings.value.length,
  },
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

const mailingsSummary = computed(() => {
  if (mailingsLoading.value) return 'Загружаем список…'
  if (mailingsError.value) return 'Данные недоступны'

  const count = mailings.value.length
  return `${count.toLocaleString('ru-RU')} ${pluralize(count, ['рассылка', 'рассылки', 'рассылок'])}`
})

const runningMailingsCount = computed(
  () => mailings.value.filter((mailing) => mailing.status === 'running').length,
)

const deliveryTotals = computed(() =>
  Object.values(mailingProgress.value).reduce(
    (totals, item) => ({
      total: totals.total + item.total,
      pending: totals.pending + item.pending + item.sending,
      sent: totals.sent + item.sent,
      failed: totals.failed + item.failed,
    }),
    { total: 0, pending: 0, sent: 0, failed: 0 },
  ),
)

const statusMeta = {
  draft: { label: 'Черновик', className: 'status--draft' },
  running: { label: 'Запущена', className: 'status--running' },
  paused: { label: 'Приостановлена', className: 'status--paused' },
  done: { label: 'Завершена', className: 'status--done' },
}

const mailingStatus = (status) =>
  statusMeta[status] || { label: status, className: 'status--draft' }

const progressFor = (mailingId) => mailingProgress.value[mailingId] || null

const progressPercent = (mailingId) => {
  const progress = progressFor(mailingId)
  if (!progress?.total) return 0
  return Math.min(100, Math.round(((progress.sent + progress.failed) / progress.total) * 100))
}

const deliverySummary = (mailingId) => {
  const progress = progressFor(mailingId)
  if (!progress?.total) return 'Готовим очередь отправки…'
  const delivered = `${progress.sent.toLocaleString('ru-RU')} из ${progress.total.toLocaleString('ru-RU')} отправлено`
  return progress.failed
    ? `${delivered}, ошибок: ${progress.failed.toLocaleString('ru-RU')}`
    : delivered
}

const formatContactDate = (value) => {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return new Intl.DateTimeFormat('ru-RU', { dateStyle: 'medium' }).format(date)
}

const formatMailingDate = (value, includeTime = false) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return new Intl.DateTimeFormat('ru-RU', {
    dateStyle: 'medium',
    ...(includeTime ? { timeStyle: 'short' } : {}),
  }).format(date)
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

const loadMailings = async () => {
  mailingsLoading.value = true
  mailingsError.value = ''
  try {
    const { data } = await api.get('/mailings/list')
    mailings.value = Array.isArray(data.mailings) ? data.mailings : []
  } catch {
    mailings.value = []
    mailingsError.value =
      'Не удалось загрузить рассылки. Проверьте соединение и попробуйте ещё раз.'
  } finally {
    mailingsLoading.value = false
  }
}

const loadMailingProgress = async (silent = false) => {
  if (!silent) progressError.value = ''
  try {
    const { data } = await api.get('/mailings/progress')
    const items = Array.isArray(data.progress) ? data.progress : []
    mailingProgress.value = Object.fromEntries(items.map((item) => [item.mailing_id, item]))

    const statuses = new Map(items.map((item) => [item.mailing_id, item.status]))
    mailings.value = mailings.value.map((mailing) =>
      statuses.has(mailing.id) ? { ...mailing, status: statuses.get(mailing.id) } : mailing,
    )
    if (selectedMailing.value && statuses.has(selectedMailing.value.id)) {
      selectedMailing.value = {
        ...selectedMailing.value,
        status: statuses.get(selectedMailing.value.id),
      }
    }
  } catch {
    if (!silent) progressError.value = 'Не удалось получить актуальный прогресс отправки.'
  }
}

const selectSection = (id) => {
  activeSection.value = id
  menuOpen.value = false
  search.value = ''
}

const openComposer = (event) => {
  composerTrigger.value = event?.currentTarget || document.activeElement
  campaignTitle.value = ''
  campaignSubject.value = ''
  campaignBody.value = ''
  campaignTitleError.value = ''
  campaignSubjectError.value = ''
  campaignBodyError.value = ''
  campaignSubmitError.value = ''
  composerOpen.value = true
  nextTick(() => composerFirstField.value?.focus())
}

const closeComposer = () => {
  if (campaignSubmitting.value) return
  composerOpen.value = false
  nextTick(() => composerTrigger.value?.focus())
}

const validateCampaign = () => {
  campaignTitleError.value = campaignTitle.value.trim() ? '' : 'Укажите название рассылки.'
  campaignSubjectError.value = campaignSubject.value.trim() ? '' : 'Укажите тему письма.'
  campaignBodyError.value = campaignBody.value.trim() ? '' : 'Добавьте текст письма.'
  return !campaignTitleError.value && !campaignSubjectError.value && !campaignBodyError.value
}

const submitCampaign = async () => {
  campaignSubmitError.value = ''
  if (!validateCampaign()) return

  campaignSubmitting.value = true
  try {
    const { data } = await api.post('/mailings/create', {
      title: campaignTitle.value.trim(),
      subject: campaignSubject.value.trim(),
      body_template: campaignBody.value.trim(),
    })
    if (data.mailing) mailings.value.unshift(data.mailing)
    composerOpen.value = false
    activeSection.value = 'campaigns'
    nextTick(() => composerTrigger.value?.focus())
  } catch (error) {
    campaignSubmitError.value =
      error.response?.data?.message ||
      (error.response
        ? 'Не удалось сохранить рассылку. Проверьте данные и попробуйте ещё раз.'
        : 'Нет соединения с сервером. Проверьте, запущен ли API.')
  } finally {
    campaignSubmitting.value = false
  }
}

const openMailingDetails = async (mailing, event) => {
  mailingDetailsTrigger.value = event?.currentTarget || document.activeElement
  selectedMailing.value = mailing
  mailingDetailsError.value = ''
  mailingDetailsOpen.value = true
  mailingDetailsLoading.value = true
  try {
    const { data } = await api.get(`/mailings/${mailing.id}/get`)
    if (data.mailing) selectedMailing.value = data.mailing
  } catch (error) {
    mailingDetailsError.value =
      error.response?.data?.message || 'Не удалось загрузить актуальные данные рассылки.'
  } finally {
    mailingDetailsLoading.value = false
    nextTick(() => mailingDetailsDialog.value?.querySelector('button')?.focus())
  }
}

const closeMailingDetails = () => {
  if (startingMailingId.value) return
  mailingDetailsOpen.value = false
  nextTick(() => mailingDetailsTrigger.value?.focus())
}

const startMailing = async (mailing) => {
  if (!mailing || mailing.status !== 'draft' || startingMailingId.value) return
  startingMailingId.value = mailing.id
  mailingDetailsError.value = ''
  try {
    await api.post(`/mailings/${mailing.id}/start`)
    const startedAt = new Date().toISOString()
    const update = { status: 'running', started_at: startedAt }
    mailings.value = mailings.value.map((item) =>
      item.id === mailing.id ? { ...item, ...update } : item,
    )
    if (selectedMailing.value?.id === mailing.id) {
      selectedMailing.value = { ...selectedMailing.value, ...update }
    }
    await loadMailingProgress()
  } catch (error) {
    mailingDetailsError.value =
      error.response?.data?.message || 'Не удалось запустить рассылку. Попробуйте ещё раз.'
  } finally {
    startingMailingId.value = null
  }
}

const resetContactForm = () => {
  contactName.value = ''
  contactEmail.value = ''
  contactNameError.value = ''
  contactEmailError.value = ''
  contactSubmitError.value = ''
}

const openContactModal = (event) => {
  contactModalTrigger.value = event?.currentTarget || document.activeElement
  resetContactForm()
  contactModalOpen.value = true
  nextTick(() => contactNameField.value?.focus())
}

const closeContactModal = () => {
  if (contactSubmitting.value) return
  contactModalOpen.value = false
  nextTick(() => contactModalTrigger.value?.focus())
}

const validateContact = () => {
  const name = contactName.value.trim()
  const email = contactEmail.value.trim()
  contactNameError.value = name ? '' : 'Укажите имя контакта.'
  contactEmailError.value = !email
    ? 'Укажите email.'
    : /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
      ? ''
      : 'Проверьте формат email.'
  return !contactNameError.value && !contactEmailError.value
}

const submitContact = async () => {
  contactSubmitError.value = ''
  if (!validateContact()) return

  contactSubmitting.value = true
  try {
    const { data } = await api.post('/contact/create', {
      name: contactName.value.trim(),
      email: contactEmail.value.trim().toLocaleLowerCase('en-US'),
    })
    if (data.contact) contacts.value.unshift(data.contact)
    contactModalOpen.value = false
    nextTick(() => contactModalTrigger.value?.focus())
  } catch (error) {
    contactSubmitError.value =
      error.response?.data?.message ||
      (error.response
        ? 'Не удалось добавить контакт. Проверьте данные и попробуйте ещё раз.'
        : 'Нет соединения с сервером. Проверьте, запущен ли API.')
  } finally {
    contactSubmitting.value = false
  }
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
    if (mailingDetailsOpen.value) closeMailingDetails()
    else if (contactModalOpen.value) closeContactModal()
    else if (composerOpen.value) closeComposer()
    else if (menuOpen.value) closeMenu(true)
    return
  }
  const activeDialog = mailingDetailsOpen.value
    ? mailingDetailsDialog.value
    : contactModalOpen.value
      ? contactDialog.value
      : composer.value
  if (event.key !== 'Tab' || !activeDialog) return
  const focusable = [
    ...activeDialog.querySelectorAll(
      'button:not(:disabled), input:not(:disabled), textarea:not(:disabled)',
    ),
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
onBeforeUnmount(() => {
  document.removeEventListener('keydown', handleKeydown)
  if (progressTimer) window.clearInterval(progressTimer)
})
onMounted(async () => {
  await Promise.all([loadContacts(), loadMailings()])
  await loadMailingProgress()
  progressTimer = window.setInterval(() => {
    if (mailings.value.some((mailing) => mailing.status === 'running')) {
      loadMailingProgress(true)
    }
  }, 4000)
})

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
        <button class="nav-item" type="button" @click="logout">
          <AppIcon name="logout" /><span>Выйти</span>
        </button>
      </div>
    </aside>

    <main class="workspace">
      <div class="workspace__heading">
        <div>
          <h1>{{ sectionTitle }}</h1>
          <p v-if="activeSection === 'overview'">Актуальное состояние ваших рассылок.</p>
        </div>
      </div>

      <template v-if="activeSection === 'overview'">
        <section
          v-if="mailingsLoading"
          class="content-panel full-panel empty-state empty-state--large"
        >
          <span class="loading-indicator" aria-hidden="true"></span>
          <h2>Загружаем данные</h2>
        </section>
        <section
          v-else-if="mailingsError"
          class="content-panel full-panel empty-state empty-state--large"
          role="alert"
        >
          <AppIcon name="clock" :size="26" />
          <h2>Данные недоступны</h2>
          <p>{{ mailingsError }}</p>
          <button class="retry-button" type="button" @click="loadMailings">Повторить</button>
        </section>
        <section v-else-if="mailings.length" class="overview-grid">
          <article class="overview-metric">
            <span>Всего рассылок</span>
            <strong>{{ mailings.length.toLocaleString('ru-RU') }}</strong>
          </article>
          <article class="overview-metric">
            <span>Сейчас запущено</span>
            <strong>{{ runningMailingsCount.toLocaleString('ru-RU') }}</strong>
          </article>
          <div class="content-panel overview-latest">
            <div class="panel-heading">
              <div>
                <h2>Последние рассылки</h2>
                <p>{{ mailingsSummary }}</p>
              </div>
              <button class="text-button" type="button" @click="selectSection('campaigns')">
                Показать все
              </button>
            </div>
            <div class="mailings-list">
              <button
                v-for="mailing in mailings.slice(0, 5)"
                :key="mailing.id"
                class="mailing-row"
                type="button"
                @click="openMailingDetails(mailing, $event)"
              >
                <span class="mailing-row__main"
                  ><strong>{{ mailing.title }}</strong
                  ><small>{{ mailing.subject }}</small
                  ><small
                    v-if="['running', 'done'].includes(mailing.status) && progressFor(mailing.id)"
                    class="mailing-row__progress"
                    >{{ deliverySummary(mailing.id) }}</small
                  ></span
                >
                <span class="status-badge" :class="mailingStatus(mailing.status).className">{{
                  mailingStatus(mailing.status).label
                }}</span>
                <time :datetime="mailing.created_at">{{
                  formatMailingDate(mailing.created_at)
                }}</time>
              </button>
            </div>
          </div>
        </section>
        <section v-else class="content-panel full-panel empty-state empty-state--large">
          <span class="empty-state__icon"><AppIcon name="send" :size="30" /></span>
          <h2>Рассылок пока нет</h2>
          <p>Создайте первую рассылку, добавьте письмо и запустите отправку по списку контактов.</p>
          <button class="empty-state__action" type="button" @click="openComposer">
            Создать рассылку
          </button>
        </section>
      </template>

      <section v-else-if="activeSection === 'campaigns'" class="content-panel full-panel">
        <div class="panel-heading">
          <div>
            <h2>Все рассылки</h2>
            <p>{{ mailingsSummary }}</p>
          </div>
          <button class="small-primary" type="button" @click="openComposer">
            <AppIcon name="plus" :size="17" /> Создать
          </button>
        </div>
        <div v-if="mailingsLoading" class="empty-state" aria-live="polite">
          <span class="loading-indicator" aria-hidden="true"></span>
          <h3>Загружаем рассылки</h3>
        </div>
        <div v-else-if="mailingsError" class="empty-state" role="alert">
          <AppIcon name="clock" :size="26" />
          <h3>Рассылки недоступны</h3>
          <p>{{ mailingsError }}</p>
          <button class="retry-button" type="button" @click="loadMailings">Повторить</button>
        </div>
        <div v-else-if="mailings.length" class="mailings-list">
          <button
            v-for="mailing in mailings"
            :key="mailing.id"
            class="mailing-row"
            type="button"
            @click="openMailingDetails(mailing, $event)"
          >
            <span class="mailing-row__main"
              ><strong>{{ mailing.title }}</strong
              ><small>{{ mailing.subject }}</small
              ><small
                v-if="['running', 'done'].includes(mailing.status) && progressFor(mailing.id)"
                class="mailing-row__progress"
                >{{ deliverySummary(mailing.id) }}</small
              ></span
            >
            <span class="status-badge" :class="mailingStatus(mailing.status).className">{{
              mailingStatus(mailing.status).label
            }}</span>
            <time :datetime="mailing.created_at">{{ formatMailingDate(mailing.created_at) }}</time>
            <AppIcon name="arrow" :size="17" />
          </button>
        </div>
        <div v-else class="empty-state">
          <AppIcon name="send" :size="26" />
          <h3>Рассылок пока нет</h3>
          <p>Создайте черновик, чтобы подготовить первую отправку.</p>
          <button class="empty-state__action" type="button" @click="openComposer">
            Создать рассылку
          </button>
        </div>
      </section>

      <section v-else-if="activeSection === 'contacts'" class="content-panel full-panel">
        <div class="panel-heading">
          <div>
            <h2>Контакты</h2>
            <p>{{ contactsSummary }}</p>
          </div>
          <button class="small-primary" type="button" @click="openContactModal">
            <AppIcon name="plus" :size="17" /> Добавить контакт
          </button>
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
                : 'Добавьте первый контакт, чтобы подготовить список получателей.'
            }}
          </p>
          <button
            v-if="!search"
            class="empty-state__action"
            type="button"
            @click="openContactModal"
          >
            Добавить контакт
          </button>
        </div>
      </section>

      <template v-else-if="activeSection === 'analytics'">
        <section v-if="progressError" class="content-panel full-panel empty-state" role="alert">
          <AppIcon name="clock" :size="26" />
          <h2>Статистика недоступна</h2>
          <p>{{ progressError }}</p>
          <button class="retry-button" type="button" @click="loadMailingProgress()">
            Повторить
          </button>
        </section>
        <section v-else-if="deliveryTotals.total" class="overview-grid delivery-overview">
          <article class="overview-metric">
            <span>Отправлено писем</span>
            <strong>{{ deliveryTotals.sent.toLocaleString('ru-RU') }}</strong>
          </article>
          <article class="overview-metric">
            <span>Ожидают отправки</span>
            <strong>{{ deliveryTotals.pending.toLocaleString('ru-RU') }}</strong>
          </article>
          <article class="overview-metric">
            <span>Ошибки доставки</span>
            <strong>{{ deliveryTotals.failed.toLocaleString('ru-RU') }}</strong>
          </article>
        </section>
        <section v-else class="content-panel full-panel empty-state empty-state--large">
          <span class="empty-state__icon"><AppIcon name="chart" :size="30" /></span>
          <h2>Статистика появится после запуска</h2>
          <p>Создайте рассылку и запустите её, чтобы увидеть результаты доставки.</p>
        </section>
      </template>
    </main>

    <div v-if="composerOpen" class="composer-scrim" @click.self="closeComposer">
      <form
        ref="composer"
        class="composer"
        role="dialog"
        aria-modal="true"
        aria-labelledby="composer-title"
        novalidate
        @submit.prevent="submitCampaign"
      >
        <div class="composer__header">
          <div>
            <h2 id="composer-title">Новая рассылка</h2>
            <p>После сохранения рассылка появится в списке</p>
          </div>
          <button
            class="icon-button"
            type="button"
            aria-label="Закрыть окно"
            :disabled="campaignSubmitting"
            @click="closeComposer"
          >
            <AppIcon name="close" />
          </button>
        </div>
        <div class="field" :class="{ 'field--error': campaignTitleError }">
          <label for="campaign-title">Название</label>
          <div class="field__control">
            <input
              id="campaign-title"
              ref="composerFirstField"
              v-model="campaignTitle"
              maxlength="200"
              placeholder="Например, Новости октября"
              :aria-invalid="Boolean(campaignTitleError)"
              aria-describedby="campaign-title-error"
              @input="campaignTitleError = ''"
            />
          </div>
          <p v-if="campaignTitleError" id="campaign-title-error" class="field__error">
            {{ campaignTitleError }}
          </p>
        </div>
        <div class="field" :class="{ 'field--error': campaignSubjectError }">
          <label for="campaign-subject">Тема письма</label>
          <div class="field__control">
            <input
              id="campaign-subject"
              v-model="campaignSubject"
              maxlength="300"
              placeholder="Что увидит получатель"
              :aria-invalid="Boolean(campaignSubjectError)"
              aria-describedby="campaign-subject-error"
              @input="campaignSubjectError = ''"
            />
          </div>
          <p v-if="campaignSubjectError" id="campaign-subject-error" class="field__error">
            {{ campaignSubjectError }}
          </p>
        </div>
        <div class="field" :class="{ 'field--error': campaignBodyError }">
          <label for="campaign-body">Текст письма</label>
          <div class="field__control">
            <textarea
              id="campaign-body"
              v-model="campaignBody"
              rows="7"
              maxlength="20000"
              placeholder="Напишите сообщение для получателей"
              :aria-invalid="Boolean(campaignBodyError)"
              aria-describedby="campaign-body-error"
              @input="campaignBodyError = ''"
            ></textarea>
          </div>
          <p v-if="campaignBodyError" id="campaign-body-error" class="field__error">
            {{ campaignBodyError }}
          </p>
        </div>
        <p v-if="campaignSubmitError" class="form-error" role="alert">{{ campaignSubmitError }}</p>
        <div class="composer__actions">
          <button type="button" :disabled="campaignSubmitting" @click="closeComposer">
            Отмена
          </button>
          <button class="action-primary" type="submit" :disabled="campaignSubmitting">
            <span v-if="campaignSubmitting" class="button-spinner" aria-hidden="true"></span>
            {{ campaignSubmitting ? 'Сохраняем…' : 'Сохранить черновик' }}
          </button>
        </div>
      </form>
    </div>

    <div v-if="mailingDetailsOpen" class="composer-scrim" @click.self="closeMailingDetails">
      <section
        ref="mailingDetailsDialog"
        class="composer mailing-details"
        role="dialog"
        aria-modal="true"
        aria-labelledby="mailing-details-title"
      >
        <div class="composer__header">
          <div>
            <span
              v-if="selectedMailing"
              class="status-badge"
              :class="mailingStatus(selectedMailing.status).className"
              >{{ mailingStatus(selectedMailing.status).label }}</span
            >
            <h2 id="mailing-details-title">{{ selectedMailing?.title || 'Рассылка' }}</h2>
            <p v-if="selectedMailing">
              Создана {{ formatMailingDate(selectedMailing.created_at, true) }}
            </p>
          </div>
          <button
            class="icon-button"
            type="button"
            aria-label="Закрыть окно"
            :disabled="Boolean(startingMailingId)"
            @click="closeMailingDetails"
          >
            <AppIcon name="close" />
          </button>
        </div>
        <div v-if="mailingDetailsLoading" class="details-loading" aria-live="polite">
          <span class="loading-indicator" aria-hidden="true"></span> Обновляем данные…
        </div>
        <template v-if="selectedMailing">
          <div
            v-if="['running', 'done'].includes(selectedMailing.status)"
            class="delivery-progress"
            aria-live="polite"
          >
            <div class="delivery-progress__copy">
              <strong>Доставка писем</strong>
              <span>{{ deliverySummary(selectedMailing.id) }}</span>
            </div>
            <div
              class="delivery-progress__track"
              role="progressbar"
              aria-label="Прогресс доставки"
              :aria-valuenow="progressPercent(selectedMailing.id)"
              aria-valuemin="0"
              aria-valuemax="100"
            >
              <span :style="{ width: `${progressPercent(selectedMailing.id)}%` }"></span>
            </div>
          </div>
          <dl class="details-list">
            <div>
              <dt>Тема письма</dt>
              <dd>{{ selectedMailing.subject }}</dd>
            </div>
            <div>
              <dt>Запущена</dt>
              <dd>{{ formatMailingDate(selectedMailing.started_at, true) }}</dd>
            </div>
          </dl>
          <div class="message-preview">
            <span>Текст письма</span>
            <p>{{ selectedMailing.body_template }}</p>
          </div>
        </template>
        <p v-if="mailingDetailsError" class="form-error" role="alert">{{ mailingDetailsError }}</p>
        <p
          v-if="selectedMailing?.status === 'draft' && !contactsLoading && !contacts.length"
          id="no-contacts-hint"
          class="form-hint"
        >
          Перед запуском добавьте хотя бы один контакт.
        </p>
        <div class="composer__actions">
          <button type="button" :disabled="Boolean(startingMailingId)" @click="closeMailingDetails">
            Закрыть
          </button>
          <button
            v-if="selectedMailing?.status === 'draft'"
            class="action-primary"
            type="button"
            :disabled="
              Boolean(startingMailingId) ||
              mailingDetailsLoading ||
              (!contactsLoading && !contacts.length)
            "
            :aria-describedby="
              !contactsLoading && !contacts.length ? 'no-contacts-hint' : undefined
            "
            @click="startMailing(selectedMailing)"
          >
            <span v-if="startingMailingId" class="button-spinner" aria-hidden="true"></span>
            {{ startingMailingId ? 'Запускаем…' : 'Запустить рассылку' }}
          </button>
        </div>
      </section>
    </div>

    <div v-if="contactModalOpen" class="composer-scrim" @click.self="closeContactModal">
      <form
        ref="contactDialog"
        class="composer contact-modal"
        role="dialog"
        aria-modal="true"
        aria-labelledby="contact-modal-title"
        novalidate
        @submit.prevent="submitContact"
      >
        <div class="composer__header">
          <div>
            <h2 id="contact-modal-title">Новый контакт</h2>
            <p>Он сразу появится в списке получателей</p>
          </div>
          <button
            class="icon-button"
            type="button"
            aria-label="Закрыть окно"
            :disabled="contactSubmitting"
            @click="closeContactModal"
          >
            <AppIcon name="close" />
          </button>
        </div>

        <div class="field" :class="{ 'field--error': contactNameError }">
          <label for="contact-name">Имя</label>
          <div class="field__control">
            <input
              id="contact-name"
              ref="contactNameField"
              v-model="contactName"
              type="text"
              autocomplete="name"
              maxlength="200"
              placeholder="Например, Анна Смирнова"
              :aria-invalid="Boolean(contactNameError)"
              aria-describedby="contact-name-error"
              @input="contactNameError = ''"
            />
          </div>
          <p v-if="contactNameError" id="contact-name-error" class="field__error">
            {{ contactNameError }}
          </p>
        </div>

        <div class="field" :class="{ 'field--error': contactEmailError }">
          <label for="contact-email">Email</label>
          <div class="field__control">
            <input
              id="contact-email"
              v-model="contactEmail"
              type="email"
              inputmode="email"
              autocomplete="email"
              maxlength="320"
              placeholder="anna@example.com"
              :aria-invalid="Boolean(contactEmailError)"
              aria-describedby="contact-email-error"
              @input="contactEmailError = ''"
            />
          </div>
          <p v-if="contactEmailError" id="contact-email-error" class="field__error">
            {{ contactEmailError }}
          </p>
        </div>

        <p v-if="contactSubmitError" class="form-error" role="alert">
          {{ contactSubmitError }}
        </p>

        <div class="composer__actions">
          <button type="button" :disabled="contactSubmitting" @click="closeContactModal">
            Отмена
          </button>
          <button class="action-primary" type="submit" :disabled="contactSubmitting">
            <span v-if="contactSubmitting" class="button-spinner" aria-hidden="true"></span>
            {{ contactSubmitting ? 'Добавляем…' : 'Добавить контакт' }}
          </button>
        </div>
      </form>
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
.overview-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
  margin: 0 28px 28px;
}
.overview-metric {
  min-width: 0;
  padding: 22px;
  border: 1px solid var(--line);
  border-radius: 14px;
  background: white;
}
.overview-metric span {
  display: block;
  margin-bottom: 8px;
  color: var(--muted);
  font-size: 11px;
  font-weight: 600;
}
.overview-metric strong {
  font-size: 28px;
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.03em;
}
.overview-latest {
  grid-column: 1 / -1;
  overflow: hidden;
}
.delivery-overview {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}
.text-button {
  border: 0;
  color: var(--blue);
  background: transparent;
  font-size: 11px;
  font-weight: 700;
  cursor: pointer;
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
.small-primary:hover {
  background: var(--blue-hover);
}
.empty-state__action {
  min-height: 38px;
  margin-top: 20px;
  padding: 0 17px;
  border: 0;
  border-radius: 19px;
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
.mailings-list {
  display: grid;
}
.mailing-row {
  display: grid;
  min-width: 0;
  min-height: 74px;
  grid-template-columns: minmax(160px, 1fr) auto 110px auto;
  align-items: center;
  gap: 16px;
  padding: 11px 20px;
  border: 0;
  border-bottom: 1px solid #edf0f2;
  color: var(--ink);
  background: white;
  text-align: left;
  cursor: pointer;
}
.mailing-row:last-child {
  border-bottom: 0;
}
.mailing-row:hover {
  background: #f8fafd;
}
.mailing-row__main {
  min-width: 0;
  display: grid;
  gap: 5px;
}
.mailing-row__main strong,
.mailing-row__main small {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.mailing-row__main strong {
  font-size: 12px;
}
.mailing-row__main small,
.mailing-row time {
  color: var(--muted);
  font-size: 11px;
}
.mailing-row__main .mailing-row__progress {
  color: var(--blue);
  font-weight: 600;
}
.mailing-row time {
  font-variant-numeric: tabular-nums;
  text-align: right;
}
.status-badge {
  width: fit-content;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 10px;
  font-weight: 700;
  white-space: nowrap;
}
.status--draft {
  color: #4f545a;
  background: #eceff3;
}
.status--running,
.status--done {
  color: var(--green);
  background: var(--green-bg);
}
.status--paused {
  color: var(--amber);
  background: var(--amber-bg);
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
.field {
  display: grid;
  gap: 7px;
}
.field label {
  color: #3c4043;
  font-size: 11px;
  font-weight: 700;
}
.field__control {
  border: 1px solid var(--line-strong);
  border-radius: 10px;
  background: white;
  transition:
    border-color 0.16s ease,
    box-shadow 0.16s ease;
}
.field__control:focus-within {
  border-color: var(--blue);
  box-shadow: 0 0 0 3px rgba(11, 87, 208, 0.12);
}
.field input,
.field textarea {
  width: 100%;
  border: 0;
  border-radius: inherit;
  outline: 0;
  color: var(--ink);
  background: transparent;
  font-size: 12px;
}
.field input {
  min-height: 44px;
  padding: 0 13px;
}
.field textarea {
  min-height: 126px;
  padding: 12px 13px;
  line-height: 1.55;
  resize: vertical;
}
.field input::placeholder,
.field textarea::placeholder {
  color: #777c82;
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
.composer__actions .action-primary {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  border-color: var(--blue);
  color: white;
  background: var(--blue);
}
.composer__actions .action-primary:hover:not(:disabled) {
  border-color: var(--blue-hover);
  background: var(--blue-hover);
}
.contact-modal {
  gap: 16px;
}
.mailing-details {
  max-height: min(720px, calc(100vh - 40px));
  overflow-y: auto;
}
.mailing-details .composer__header > div {
  min-width: 0;
}
.mailing-details .composer__header h2 {
  overflow-wrap: anywhere;
  margin-top: 10px;
}
.details-loading {
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--muted);
  font-size: 11px;
}
.details-loading .loading-indicator {
  width: 18px;
  height: 18px;
}
.delivery-progress {
  display: grid;
  gap: 10px;
  padding: 14px;
  border-radius: 10px;
  background: var(--surface-soft);
}
.delivery-progress__copy {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 14px;
  font-size: 11px;
}
.delivery-progress__copy span {
  color: var(--muted);
  text-align: right;
}
.delivery-progress__track {
  height: 6px;
  overflow: hidden;
  border-radius: 999px;
  background: #dfe5ec;
}
.delivery-progress__track span {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: var(--blue);
}
.details-list {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(130px, auto);
  gap: 12px;
  margin: 0;
}
.details-list div {
  min-width: 0;
  padding: 12px 14px;
  border-radius: 10px;
  background: var(--surface-soft);
}
.details-list dt,
.message-preview > span {
  margin-bottom: 5px;
  color: var(--muted);
  font-size: 10px;
  font-weight: 700;
}
.details-list dd {
  margin: 0;
  overflow-wrap: anywhere;
  font-size: 12px;
}
.message-preview {
  min-width: 0;
  padding: 14px;
  border: 1px solid var(--line);
  border-radius: 10px;
}
.message-preview p {
  max-height: 240px;
  margin: 8px 0 0;
  overflow-y: auto;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
  font-size: 12px;
  line-height: 1.6;
}
.form-hint {
  margin: 0;
  color: var(--amber);
  font-size: 11px;
  line-height: 1.5;
}
.field--error .field__control {
  border-color: var(--red);
}
.field__error {
  margin: 6px 0 0;
  color: var(--red);
  font-size: 11px;
}
.form-error {
  margin: 0;
  padding: 11px 13px;
  border-radius: 9px;
  color: #8c1d18;
  background: var(--red-bg);
  font-size: 11px;
  line-height: 1.5;
}
.button-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.45);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
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
  .overview-grid {
    margin: 0 14px 20px;
  }
  .delivery-overview {
    grid-template-columns: 1fr;
  }
  .panel-heading {
    padding: 16px;
  }
  .composer {
    padding: 20px;
  }
  .mailing-row {
    grid-template-columns: minmax(0, 1fr) auto auto;
    gap: 10px;
    padding: 12px 16px;
  }
  .mailing-row time {
    display: none;
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
  .overview-grid {
    grid-template-columns: 1fr;
  }
  .overview-latest {
    grid-column: auto;
  }
  .details-list {
    grid-template-columns: 1fr;
  }
  .delivery-progress__copy {
    align-items: flex-start;
    flex-direction: column;
    gap: 4px;
  }
  .delivery-progress__copy span {
    text-align: left;
  }
}
</style>
