import "./assets/main.css";

import { createApp } from "vue";
import { createPinia } from "pinia";
import piniaPluginPersistedstate from "pinia-plugin-persistedstate";
import { useAuthStore } from "./stores/auth.store";
import instance from "@/axios";

import App from "./App.vue";
import router from "./router";

const app = createApp(App);
const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);
app.use(pinia);
app.use(router);

const authStore = useAuthStore();
if (authStore.token) {
	instance.defaults.headers.common["Authorization"] =
		`Bearer ${authStore.token}`;
}

app.mount("#app");
