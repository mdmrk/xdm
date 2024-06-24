<script setup>
import ThreeColumnLayout from "@/layouts/ThreeColumnLayout.vue";
import { useAuthStore } from "@/stores/auth.store";

const authStore = useAuthStore();
const isAuthenticated = authStore.isAuthenticated;
</script>
<template>
  <div class="text-sm">
    <div v-if="isAuthenticated" class="flex flex-col gap-1">
      <RouterLink :to="{ name: 'profile', params: { id: authStore.UserId } }">Profile</RouterLink>
      <RouterLink to="/logout">Logout</RouterLink>
    </div>
    <div v-else>
      <RouterLink to="/signin">Login</RouterLink>
    </div>
  </div>
</template>
