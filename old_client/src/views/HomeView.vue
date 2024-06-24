<script lang="ts">
import { defineComponent } from "vue";
import ThreeColumnLayout from "@/layouts/ThreeColumnLayout.vue";
import Post from "@/components/Post.vue";
import { useAuthStore } from "@/stores/auth.store";
import { getPosts, postPost } from "@/post";
import { getUserLikes } from "@/user";
import type { IPost } from "@/types";

export default defineComponent({
  name: "Home",
  components: {
    ThreeColumnLayout,
    Post,
  },
  data() {
    return {
      loading: true,
      postcontent: "",
      shownewpost: false,
      posts: [] as IPost[],
      likedPosts: [] as string[],
    };
  },
  computed: {
    isAuthenticated() {
      const authStore = useAuthStore();
      return authStore.isAuthenticated;
    },
    UserId() {
      const authStore = useAuthStore();
      return authStore.UserId;
    },
  },
  mounted() {
    this.fetchPosts();
    this.fetchLikes();
  },
  methods: {
    async newPost() {
      try {
        const response = await postPost(this.postcontent);
        this.postcontent = ""
        this.shownewpost = false
        await this.fetchPosts()
      } catch (error) {
        console.error(error);
      }
    },
    async fetchPosts() {
      try {
        const response = await getPosts();
        this.posts = response;
      } catch (error) {
        console.error(error);
      } finally {
        this.loading = false;
      }
    },
    async fetchLikes() {
      try {
        const id = this.UserId;
        const response = await getUserLikes(id);
        this.likedPosts = response;
      } catch (error) {
        console.error(error);
      } finally {
        this.loading = false;
      }
    },
  },
});
</script>

<template>
  <ThreeColumnLayout>
    <h2 class="text-3xl font-bold mb-4">Home</h2>
    <div v-if="isAuthenticated">
      <div v-if="loading">

      </div>
      <div v-else>
        <div v-show="shownewpost">
          <div @click.prevent="shownewpost = false" class="fixed top-0 left-0 w-full h-full bg-black opacity-50"></div>
          <form @submit.prevent="newPost" class="fixed top-1/3 left-1/2 w-[36rem] z-20 h-96 translate-x-[-50%]">
            <div class="grid gap-y-4">
              <div>
                <textarea rows="4" type="text" v-model="postcontent" id="postcontent" name="postcontent"
                  class="py-3 px-4 resize-none block w-full border border-gray-200 rounded-lg text-sm focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:pointer-events-none dark:bg-neutral-900 dark:border-neutral-700 dark:text-neutral-400 dark:placeholder-neutral-500 dark:focus:ring-neutral-600"
                  required />
              </div>
            </div>
            <button type="submit"
              class=" w-full py-3 px-4 mt-3 inline-flex justify-center items-center gap-x-2 text-sm font-semibold rounded-lg border border-transparent bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-50 disabled:pointer-events-none">Post</button>
          </form>
        </div>
        <div>
          <div v-if="posts.length === 0" class="opacity-50 italic">no posts currently, try again later</div>
          <div class="flex flex-col gap-1">
            <Post @update-likes="fetchLikes" v-for="post in posts" :key="post.Id" :post="post"
              :liked="likedPosts.indexOf(post.Id) > -1" />
          </div>
        </div>
      </div>
    </div>
    <div v-else>

            <Post @update-likes="fetchLikes" v-for="post in posts" :key="post.Id" :post="post"
              :liked="likedPosts.indexOf(post.Id) > -1" />
    </div>

    <template #aside>
      <div @click.prevent="shownewpost = true" v-if="isAuthenticated"><svg xmlns="http://www.w3.org/2000/svg" width="24"
          height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
          stroke-linejoin="round" class="icon cursor-pointer icon-tabler icons-tabler-outline icon-tabler-pencil-plus">
          <path stroke="none" d="M0 0h24v24H0z" fill="none" />
          <path d="M4 20h4l10.5 -10.5a2.828 2.828 0 1 0 -4 -4l-10.5 10.5v4" />
          <path d="M13.5 6.5l4 4" />
          <path d="M16 19h6" />
          <path d="M19 16v6" />
        </svg></div>
      <div v-else></div>
    </template>
  </ThreeColumnLayout>
</template>
