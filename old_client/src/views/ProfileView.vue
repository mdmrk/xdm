<script lang="ts">
import { defineComponent } from "vue";
import ThreeColumnLayout from "@/layouts/ThreeColumnLayout.vue";
import { useRoute } from "vue-router";
import { getUserById, followUser, unfollowUser } from "@/user";
import type { IUser } from "@/types";
import { useAuthStore } from "@/stores/auth.store";

export default defineComponent({
  name: "Home",
  components: {
    ThreeColumnLayout,
  },
  data() {
    return {
      loading: true,
      user: {},
      following: false,
    };
  },
  mounted() {
    const route = useRoute();
    const id = route.params.id;
    this.fetchUser(id);
  },
  computed: {
    followIcon() {
      if (this.user.Id === "") {
        return 2;
      }
      return this.following ? 1 : 0;
    },
    UserId() {
      const authStore = useAuthStore();
      return authStore.UserId;
    },
  },
  methods: {
    async fetchUser(id: string) {
      try {
        const response = await getUserById(id);
        this.user = response;
      } catch (error) {
        console.error("Error fetching posts:", error);
      } finally {
        this.loading = false;
      }
    },
    async followUser() {
      try {
        const response = await followUser(this.user.Id);
      } catch (error) {
        console.error("Error fetching posts:", error);
      }
    },
    async unfollowUser() {
      try {
        const response = await unfollowUser(this.user.Id);
      } catch (error) {
        console.error("Error fetching posts:", error);
      }
    },
    async likePost() {
      try {
        const response = await likePost(this.post.Id);
        this.$emit("updateLikes");
      } catch (error) {
        console.error("Error fetching posts:", error);
      }
    },
    async unlikePost() {
      try {
        const response = await unlikePost(this.post.Id);
        this.$emit("updateLikes");
      } catch (error) {
        console.error("Error fetching posts:", error);
      }
    },
  },
});
</script>

<template>
  <ThreeColumnLayout>
    <h2 class="text-3xl font-bold mb-4">Profile</h2>
    <div class="flex flex-row justify-between">
      <div class="font-bold">{{ this.user.Alias }}</div>
      <div class="flex flex-row gap-2">
        <svg @click="this.$router.push('/chat/' + this.$route.params.id)" v-if="this.UserId !== this.$route.params.id"
          xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
          stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
          class="icon hover:cursor-pointer icon-tabler icons-tabler-outline icon-tabler-message">
          <path stroke="none" d="M0 0h24v24H0z" fill="none" />
          <path d="M8 9h8" />
          <path d="M8 13h6" />
          <path d="M18 4a3 3 0 0 1 3 3v8a3 3 0 0 1 -3 3h-5l-5 3v-3h-2a3 3 0 0 1 -3 -3v-8a3 3 0 0 1 3 -3h12z" />
        </svg>
        <div>
          <svg v-if="followIcon === 1" @click.prevent="unfollowUser" xmlns="http://www.w3.org/2000/svg" width="24"
            height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
            stroke-linejoin="round" class="icon icon-tabler icons-tabler-outline icon-tabler-user-check">
            <path stroke="none" d="M0 0h24v24H0z" fill="none" />
            <path d="M8 7a4 4 0 1 0 8 0a4 4 0 0 0 -8 0" />
            <path d="M6 21v-2a4 4 0 0 1 4 -4h4" />
            <path d="M15 19l2 2l4 -4" />
          </svg>
          <svg v-else-if="followIcon === 0" @click.prevent="followUser" xmlns="http://www.w3.org/2000/svg" width="24"
            height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
            stroke-linejoin="round" class="icon icon-tabler icons-tabler-outline icon-tabler-user-plus">
            <path stroke="none" d="M0 0h24v24H0z" fill="none" />
            <path d="M8 7a4 4 0 1 0 8 0a4 4 0 0 0 -8 0" />
            <path d="M16 19h6" />
            <path d="M19 16v6" />
            <path d="M6 21v-2a4 4 0 0 1 4 -4h4" />
          </svg>
        </div>
      </div>
    </div>
    <div class="text-gray-500">@{{ this.user.Username }}</div>
    <div class="mt-2 italic text-gray-500 text-xs">last seen - {{ this.user.Seen }}</div>

    <template #aside>
    </template>
  </ThreeColumnLayout>
</template>
