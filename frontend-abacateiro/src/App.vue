<template>
  <router-view />
</template>

<script setup>
import { createPinia } from "pinia";
import { createApp } from "vue";
import App from "./App.vue";

const pinia = createPinia();
const app = createApp(App);
app.use(pinia);
// pinia.use(piniaPluginPersistedstate);

// ✅ works because the pinia instance is now active
// const userStore = useUserStore();

defineOptions({
  name: "App",
});
</script>
