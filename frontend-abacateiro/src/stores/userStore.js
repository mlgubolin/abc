import { defineStore } from "pinia";
import axios from "axios";

export const useUserStore = defineStore("user", {
  state: () => ({
    token: null,
    userId: null,
  }),
  actions: {
    async login(username, password) {
      try {
        const response = await axios.post("http://localhost:8080/login", {
          username,
          password,
        });

        console.log("Resposta do servidor:", response.data);

        const { token, auth_id } = response.data;

        // Verifique se o token e auth_id estão presentes
        if (token && auth_id) {
          // Armazena o token e userId no estado
          this.token = token;
          this.userId = auth_id;

          // Armazena no local storage
          localStorage.setItem("token", token);
          localStorage.setItem("userId", auth_id);

          return true;
        } else {
          console.error("Token ou auth_id ausente na resposta do servidor");
          return false;
        }
      } catch (error) {
        console.error("Erro durante o login:", error);
        return false;
      }
    },
    logout() {
      this.token = null;
      this.userId = null;

      // Remove os dados do local storage
      localStorage.removeItem("token");
      localStorage.removeItem("userId");
    },
    loadFromLocalStorage() {
      const token = localStorage.getItem("token");
      const userId = localStorage.getItem("userId");

      if (token && userId) {
        this.token = token;
        this.userId = userId;
      }
    },
  },
});

// Certifique-se de chamar loadFromLocalStorage ao inicializar a aplicação
// import { useAuthStore } from 'src/stores/authStore';
// const authStore = useAuthStore();
// authStore.loadFromLocalStorage();
