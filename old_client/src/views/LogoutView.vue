<script setup lang="ts">
import { useAuthStore } from "@/stores/auth.store";
import { onMounted } from "vue";
import { useRouter } from "vue-router";
import BlankLayout from "@/layouts/BlankLayout.vue";

const authStore = useAuthStore();
const router = useRouter();

onMounted(() => {
	authStore.logout();
	router.push({ name: "signin" });
});
</script>

<template>
  <BlankLayout>
  </BlankLayout>
</template>

