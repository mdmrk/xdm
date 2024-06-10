<script setup lang="ts">
import { ref } from "vue";
import { useAuthStore } from "@/stores/auth.store";
import BlankLayout from "@/layouts/BlankLayout.vue";

const credentials = ref({
	alias: "",
	username: "",
	password: "",
	confirmPassword: "",
});

const error = ref<string | null>(null);
const authStore = useAuthStore();

const signup = async () => {
	if (credentials.value.password !== credentials.value.confirmPassword) {
		error.value = "Passwords do not match";
		return;
	}

	error.value = null;

	try {
		await authStore.signup({
			alias: credentials.value.alias,
			username: credentials.value.username,
			password: credentials.value.password,
		});
	} catch (err) {
		error.value = authStore.error;
	}
};
</script>

<template>
  <BlankLayout class="max-w-xl">
    <div class="bg-white border border-gray-200 rounded-xl shadow-sm dark:bg-neutral-900 dark:border-neutral-700">
      <div class="p-4 sm:p-7">
        <div class="text-center">
          <h1 class="block text-2xl font-bold text-gray-800 dark:text-white">Sign up</h1>
          <p class="mt-2 text-sm text-gray-600 dark:text-neutral-400">
            Already have an account?
            <RouterLink class="text-blue-600 decoration-2 hover:underline font-medium dark:text-blue-500" to="/signin">
              Sign in here
            </RouterLink>
          </p>
        </div>
        <div class="mt-5">
          <form @submit.prevent="signup">
            <div class="grid gap-y-4">
              <div>
                <label for="alias" class="block text-sm mb-2 dark:text-white">Alias</label>
                <input v-model="credentials.alias" id="alias" name="alias" type="text" class="py-3 px-4 block w-full border border-gray-200 rounded-lg text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-neutral-900 dark:border-neutral-700 dark:text-neutral-400 dark:placeholder-neutral-500 dark:focus:ring-neutral-600" required>
              </div>

              <div>
                <label for="username" class="block text-sm mb-2 dark:text-white">Username</label>
                <input v-model="credentials.username" id="username" name="username" type="text" class="py-3 px-4 block w-full border border-gray-200 rounded-lg text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-neutral-900 dark:border-neutral-700 dark:text-neutral-400 dark:placeholder-neutral-500 dark:focus:ring-neutral-600" required>
              </div>

              <div>
                <label for="password" class="block text-sm mb-2 dark:text-white">Password</label>
                <input v-model="credentials.password" type="password" id="password" name="password" class="py-3 px-4 block w-full border border-gray-200 rounded-lg text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-neutral-900 dark:border-neutral-700 dark:text-neutral-400 dark:placeholder-neutral-500 dark:focus:ring-neutral-600" required>
              </div>

              <div>
                <label for="confirm-password" class="block text-sm mb-2 dark:text-white">Confirm Password</label>
                <input v-model="credentials.confirmPassword" type="password" id="confirm-password" name="confirm-password" class="py-3 px-4 block w-full border border-gray-200 rounded-lg text-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-neutral-900 dark:border-neutral-700 dark:text-neutral-400 dark:placeholder-neutral-500 dark:focus:ring-neutral-600" required>
              </div>

              <div v-if="error" class="text-red-500 text-sm mt-2">{{ error }}</div>

              <button type="submit" class="w-full py-3 px-4 inline-flex justify-center items-center gap-x-2 text-sm font-semibold rounded-lg border border-transparent bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-50 disabled:pointer-events-none" :disabled="authStore.loading">
                Sign up
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </BlankLayout>
</template>

