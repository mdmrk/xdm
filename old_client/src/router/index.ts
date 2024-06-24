import { createRouter, createWebHistory } from "vue-router";
import HomeView from "@/views/HomeView.vue";
import SigninView from "@/views/SigninView.vue";
import SignupView from "@/views/SignupView.vue";
import LogoutView from "@/views/LogoutView.vue";
import ProfileView from "@/views/ProfileView.vue";
import ChatView from "@/views/ChatView.vue";
import { useAuthStore } from "@/stores/auth.store";

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: "/",
			name: "home",
			component: HomeView,
		},
		{
			path: "/signin",
			name: "signin",
			component: SigninView,
		},
		{
			path: "/signup",
			name: "signup",
			component: SignupView,
		},
		{
			path: "/logout",
			name: "logout",
			component: LogoutView,
			meta: { requiresAuth: true },
		},
		{
			path: "/profile/:id",
			name: "profile",
			component: ProfileView,
		},
		{
			path: "/chat/:id",
			name: "chat",
			component: ChatView,
			meta: { requiresAuth: true },
		},
	],
});

router.beforeEach((to, _, next) => {
	const authStore = useAuthStore();

	if (to.meta.requiresAuth && !authStore.isAuthenticated) {
		next({ name: "signin" });
	} else {
		next();
	}
});

export default router;
