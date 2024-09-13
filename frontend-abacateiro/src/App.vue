<template>
  <router-view />
</template>

<script setup>

import { createPinia } from 'pinia'
import { createApp } from 'vue'
import App from './App.vue'

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)
app.mount('#app')

defineOptions({
  name: 'App'
});
</script>
