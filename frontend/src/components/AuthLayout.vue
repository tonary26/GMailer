<script setup>
import AppLogo from './AppLogo.vue'

defineProps({ title: String })
</script>

<template>
  <main class="auth-page">
    <section class="auth-story" aria-labelledby="auth-story-title">
      <AppLogo />
      <div class="auth-story__content">
        <h1 id="auth-story-title">Письма доходят.<br /><em>Вы видите главное.</em></h1>
        <p>
          Соберите аудиторию, запустите рассылку и следите за доставкой в одном спокойном рабочем
          пространстве.
        </p>
      </div>
      <p class="auth-story__footer">GMailer · панель управления рассылками</p>
    </section>

    <section class="auth-panel">
      <div class="auth-panel__mobile-logo"><AppLogo /></div>
      <div class="auth-card">
        <h2>{{ title }}</h2>
        <slot />
      </div>
    </section>
  </main>
</template>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: grid;
  grid-template-columns: minmax(420px, 1.08fr) minmax(420px, 0.92fr);
  background: var(--surface);
}
.auth-story {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  padding: 38px clamp(38px, 6vw, 88px);
  background: #eaf2ff;
}
.auth-story__content {
  max-width: 610px;
  margin: auto 0;
}
.auth-story h1 {
  max-width: 650px;
  margin: 0;
  color: #202124;
  font-size: clamp(42px, 5.4vw, 76px);
  font-weight: 600;
  letter-spacing: -0.045em;
  line-height: 1.03;
  text-wrap: balance;
}
.auth-story h1 em {
  color: var(--blue);
  font-style: normal;
}
.auth-story__content > p:last-child {
  max-width: 560px;
  margin: 28px 0 0;
  color: #4f5967;
  font-size: 17px;
  line-height: 1.7;
}
.auth-story__footer {
  margin: 24px 0 0;
  color: #515963;
  font-size: 12px;
}
.auth-panel {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 42px clamp(32px, 6vw, 92px);
  background: var(--surface);
}
.auth-panel__mobile-logo {
  display: none;
}
.auth-card {
  width: min(430px, 100%);
  margin: auto;
}
.auth-card h2 {
  margin: 0 0 34px;
  font-size: clamp(30px, 3vw, 40px);
  font-weight: 600;
  letter-spacing: -0.035em;
}
@media (max-width: 900px) {
  .auth-page {
    grid-template-columns: 1fr;
  }
  .auth-story {
    display: none;
  }
  .auth-panel {
    min-height: 100vh;
    padding: 28px 24px;
  }
  .auth-panel__mobile-logo {
    display: block;
    margin-bottom: 50px;
  }
  .auth-card {
    margin: auto;
  }
}
</style>
