<template>
  <RouterLink class="brand" to="/" aria-label="GMailer - на главную">
    <span class="brand__mark" aria-hidden="true">
      <svg viewBox="0 0 32 24" role="img">
        <path
          d="M3 5.5 16 15 29 5.5"
          fill="none"
          stroke="#0b57d0"
          stroke-width="3.2"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
        <path
          d="M4.5 4h23A2.5 2.5 0 0 1 30 6.5v14A2.5 2.5 0 0 1 27.5 23h-23A2.5 2.5 0 0 1 2 20.5v-14A2.5 2.5 0 0 1 4.5 4Z"
          fill="none"
          stroke="#5f6368"
          stroke-width="2"
        />
        <path
          d="m3 21 9-9M29 21l-9-9"
          fill="none"
          stroke="#ea4335"
          stroke-width="2.2"
          stroke-linecap="round"
        />
      </svg>
    </span>
    <span class="brand__name">GMailer</span>
  </RouterLink>
</template>

<style scoped>
.brand {
  display: inline-flex;
  align-items: center;
  gap: 11px;
  color: var(--ink);
  text-decoration: none;
}
.brand__mark {
  display: grid;
  width: 34px;
  height: 34px;
  place-items: center;
}
.brand__mark svg {
  width: 30px;
}
.brand__name {
  color: #3c4043;
  font-size: 20px;
  font-weight: 600;
  letter-spacing: -0.02em;
}
</style>
