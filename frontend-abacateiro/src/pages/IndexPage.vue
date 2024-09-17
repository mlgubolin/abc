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
import { ref } from "vue";
import { useUserStore } from "src/stores/user-store";
import { useRouter } from "vue-router";
import { useQuasar } from "quasar";

const userStore = useUserStore();
const router = useRouter();
const $q = useQuasar();

const token = ref(userStore.token);

const username = ref("");
const password = ref("");

const handleLogin = async () => {
  try {
    // Here you would typically call an API to validate credentials
    // For this example, we'll just simulate a successful login
    await userStore.login(token.value);
    $q.notify({
      color: "positive",
      message: "Login successful",
      icon: "check",
    });
    router.push("/users");
  } catch (error) {
    console.error("Login failed:", error);
    $q.notify({
      color: "negative",
      message: "Login failed. Please try again.",
      icon: "warning",
    });
  }
};

defineOptions({
  name: "IndexPage",
});
</script>

<style scoped>
.login-card {
  width: 100%;
  max-width: 400px;
}
</style>
