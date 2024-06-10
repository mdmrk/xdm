import { defineStore } from "pinia";
import axios from "axios";
import type { User, AuthResponse } from "@/types";
import instance from "@/axios";

interface AuthState {
	userId: string | null;
	token: string | null;
	loading: boolean;
	error: string | null;
}

export const useAuthStore = defineStore("auth", {
	persist: true,
	state: (): AuthState => ({
		userId: null,
		token: null,
		loading: false,
		error: null,
	}),

	getters: {
		isAuthenticated: (state: AuthState): boolean => !!state.token,
		UserId: (state: AuthState): string | null => state.userId,
	},

	actions: {
		async signup(credentials: {
			alias: string;
			username: string;
			password: string;
		}) {
			this.loading = true;
			this.error = null;
			try {
				const response = await instance.post<AuthResponse>(
					"/signup",
					credentials,
				);

				this.token = response.data.token;
				instance.defaults.headers.common["Authorization"] =
					`Bearer ${this.token}`;
			} catch (error) {
				this.error =
					axios.isAxiosError(error) && error.response
						? error.response.data
						: "Network error";
			} finally {
				this.loading = false;
			}
		},

		async signin(credentials: { username: string; password: string }) {
			this.loading = true;
			this.error = null;

			try {
				const response = await instance.post<AuthResponse>(
					"/signin",
					credentials,
				);
				this.token = response.data.token;
				this.userId = response.data.userId;
				instance.defaults.headers.common["Authorization"] =
					`Bearer ${this.token}`;
			} catch (error) {
				this.error =
					axios.isAxiosError(error) && error.response
						? error.response.data.message
						: "Network error";
			} finally {
				this.loading = false;
			}
		},

		logout() {
			this.userId = null;
			this.token = null;
			delete instance.defaults.headers.common["Authorization"];
		},
	},
});
