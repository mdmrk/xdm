<template>
  <div class="text-sm bg-neutral-900 px-2 py-1">
    <div class="flex flex-col">
      <RouterLink :to="{ name: 'profile', params: { id: post.UserId } }" class="font-bold no-underline">{{ user.Alias }}
      </RouterLink>
      <div class="pb-4 pt-1 ml-4">{{ post.Body }}</div>
      <div class="flex flex-row gap-1 items-center">
        <div>
          <svg v-if="!liked" @click="likePost" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
            viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
            stroke-linejoin="round"
            class="icon hover:fill-white hover:cursor-pointer icon-tabler icons-tabler-outline icon-tabler-heart">
            <path stroke="none" d="M0 0h24v24H0z" fill="none" />
            <path d="M19.5 12.572l-7.5 7.428l-7.5 -7.428a5 5 0 1 1 7.5 -6.566a5 5 0 1 1 7.5 6.572" />
          </svg>
          <svg v-else @click="unlikePost" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
            fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
            class="icon fill-white hover:cursor-pointer icon-tabler icons-tabler-outline icon-tabler-heart">
            <path stroke="none" d="M0 0h24v24H0z" fill="none" />
            <path d="M19.5 12.572l-7.5 7.428l-7.5 -7.428a5 5 0 1 1 7.5 -6.566a5 5 0 1 1 7.5 6.572" />
          </svg>
        </div>
        <div class="select-none">{{ post.Likes }}</div>
        <div class="text-gray-500 ml-auto text-xs">{{ post.CreatedAt }}</div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { RouterLink } from "vue-router";
import type { Post } from "@/types";
import { getUserById } from "@/user";
import { likePost, unlikePost } from "@/post";
import type { PropType } from "vue";

export default {
  name: "Post",
  data() {
    return {
      user: {},
    };
  },
  props: {
    post: {
      type: Object as PropType<Post>,
      required: true,
    },
    liked: {
      type: Boolean,
      required: false,
      default: false,
    },
  },
  emits: {
    updateLikes() {
      return true;
    },
  },
  mounted() {
    this.fetchUser(this.post.UserId);
  },
  methods: {
    async fetchUser(id: string) {
      try {
        const response = await getUserById(id);
        this.$emit("updateLikes");
        this.user = response;
      } catch (error) {
        console.error(error);
      }
    },
    async likePost() {
      try {
        const response = await likePost(this.post.Id);
        this.$emit("updateLikes");
      } catch (error) {
        console.error(error);
      }
    },
    async unlikePost() {
      try {
        const response = await unlikePost(this.post.Id);
        this.$emit("updateLikes");
      } catch (error) {
        console.error(error);
      }
    },
  },
  components: {
    RouterLink,
  },
};
</script>
