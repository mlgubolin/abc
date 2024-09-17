import { defineStore } from "pinia";

export const useUserStore = defineStore("user", {
  state: () => ({
    user: "",
    token: "",
    isLoggedIn: false,
  }),
  actions: {
    login(user) {
      this.user = user;
      this.isLoggedIn = true;
    },
    logout() {
      this.user = "";
      this.isLoggedIn = false;
    },
  },
  persist: true,
});
