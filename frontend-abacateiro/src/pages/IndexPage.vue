<template>
  <q-page class="flex flex-center">
    <q-card class="login-card" v-if="token.value != ''">
      <q-card-section>
        <div class="text-h6">Login</div>
      </q-card-section>

      <q-card-section>
        <q-form @submit="handleLogin" class="q-gutter-md">
          <q-input
            v-model="username"
            label="Username"
            :rules="[(val) => !!val || 'Username is required']"
          />

          <q-input
            v-model="password"
            label="Password"
            type="password"
            :rules="[(val) => !!val || 'Password is required']"
          />

          <div>
            <q-btn label="Login" type="submit" color="primary" />
          </div>
        </q-form>
      </q-card-section>
    </q-card>
    <img
      v-else
      alt="Quasar logo"
      src="~assets/quasar-logo-vertical.svg"
      style="width: 200px; height: 200px"
    />
  </q-page>
</template>

<script setup>
import { useUserStore } from "src/stores/userStore";
import { ref } from "vue";

defineOptions({
  name: "IndexPage",
});

const token = ref(useUserStore().token);
</script>

<style scoped>
.login-card {
  width: 100%;
  max-width: 400px;
}
</style>
