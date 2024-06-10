<script setup lang="ts">
import BlankLayout from "@/layouts/BlankLayout.vue";
import { ref } from "vue";
import { useAuthStore } from "@/stores/auth.store";
import { useRouter } from "vue-router";

const router = useRouter();
const credentials = ref({
	username: "",
	password: "",
});

const error = ref<string | null>(null);
const authStore = useAuthStore();

const signin = async () => {
	error.value = null;

	try {
		await authStore.signin({
			username: credentials.value.username,
			password: credentials.value.password,
		});
		router.push({ path: "/" });
	} catch (err) {
		error.value = authStore.error || "An error occurred during sign in.";
	}
};
</script>

<template>
  <BlankLayout class="max-w-xl">
    <div class="bg-white border border-gray-200 rounded-xl shadow-sm dark:bg-neutral-900 dark:border-neutral-700">
      <div class="p-4 sm:p-7">
        <div class="text-center">
          <h1 class="block text-2xl font-bold text-gray-800 dark:text-white">Sign in</h1>
          <p class="mt-2 text-sm text-gray-600 dark:text-neutral-400">
            Don't have an account yet?
            <RouterLink class="text-blue-600 decoration-2 hover:underline font-medium dark:text-blue-500" to="/signup">
              Sign up here
            </RouterLink>
          </p>
        </div>

        <div class="mt-5">
          <form @submit.prevent="signin">
            <div class="grid gap-y-4">
              <div>
                <label for="username" class="block text-sm mb-2 dark:text-white">Username</label>
                <div class="relative">
                  <input type="text" v-model="credentials.username" id="username" name="username"
                    class="py-3 px-4 block w-full border border-gray-200 rounded-lg text-sm focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:pointer-events-none dark:bg-neutral-900 dark:border-neutral-700 dark:text-neutral-400 dark:placeholder-neutral-500 dark:focus:ring-neutral-600"
                    required>
                </div>
              </div>
              <div>
                <div class="flex justify-between items-center">
                  <label for="password" class="block text-sm mb-2 dark:text-white">Password</label>
                  <a class="text-sm text-blue-600 decoration-2 hover:underline font-medium" href="/recover-account">
                    Forgot password?
                  </a>
                </div>
                <div class="relative">
                  <input type="password" v-model="credentials.password" id="password" name="password"
                    class="py-3 px-4 block w-full border border-gray-200 rounded-lg text-sm focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:pointer-events-none dark:bg-neutral-900 dark:border-neutral-700 dark:text-neutral-400 dark:placeholder-neutral-500 dark:focus:ring-neutral-600"
                    required>
                </div>
              </div>

              <div class="flex items-center">
                <div class="flex">
                  <input id="remember-me" name="remember-me" type="checkbox"
                    class="shrink-0 mt-0.5 border-gray-200 rounded text-blue-600 focus:ring-blue-500 dark:bg-neutral-800 dark:border-neutral-700 dark:checked:bg-blue-500 dark:checked:border-blue-500 dark:focus:ring-offset-gray-800">
                </div>
                <div class="ms-3">
                  <label for="remember-me" class="text-sm dark:text-white">Remember me</label>
                </div>
              </div>

              <button type="submit"
                class="w-full py-3 px-4 inline-flex justify-center items-center gap-x-2 text-sm font-semibold rounded-lg border border-transparent bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-50 disabled:pointer-events-none"
                :disabled="authStore.loading">Sign in</button>

              <div v-if="error" class="text-red-600 mt-4">
                {{ error }}
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </BlankLayout>
</template>
