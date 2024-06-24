<script lang="ts">
import { defineComponent } from "vue";
import TwoColumnLayout from "@/layouts/TwoColumnLayout.vue";
import { useRoute } from "vue-router";
import { getUserById } from "@/user";
import { sendChat, retrieveChat } from "@/chat";
import Message from "@/components/Message.vue";
import type { IUser, IMessage } from "@/types";

export default defineComponent({
  name: "Home",
  components: {
    TwoColumnLayout,
    Message,
  },
  data() {
    return {
      loading: true,
      user: {} as IUser,
      targetUser: {} as IUser,
      messages: [] as IMessage[],
      userinput: "",
    };
  },
  async mounted() {
    const route = useRoute();
    const id = route.params.id as string;
    await this.fetchUser(id);
    await this.fetchMessages(id);
  },
  methods: {
    async fetchUser(id: string) {
      try {
        const response = await getUserById(id);
        this.targetUser = response;
      } catch (error) {
        console.error(error);
      } finally {
        this.loading = false;
      }
    },
    scrollToBot() {
      const scrollableElement = document.getElementById("chatbox");
      if (scrollableElement) {
        scrollableElement.scrollTop = scrollableElement.scrollHeight;
      }
    },
    async fetchMessages(targetUserId: string) {
      try {
        const response = await retrieveChat(targetUserId).data;
        this.messages = response;
        this.scrollToBot();
      } catch (error) {
        console.error("Failed to fetch messages:", error);
      }
    },
    async send() {
      if (this.userinput === "") return;
      try {
        const response = await sendChat(this.targetUser.Id);
        this.userinput = "";
        this.scrollToBot();
      } catch (error) {
        console.error("Failed to send message:", error);
      }
    },
  },
});
</script>

<template>
  <TwoColumnLayout>
    <div class="flex items-end flex-row gap-2">
      <h2 class="text-3xl font-bold mb-4">Chat</h2>
      <h5 class="text-xl font-bold text-gray-500 mb-4">{{ targetUser.username }}</h5>
    </div>
    <div class="w-full h-full">
      <div id="chatbox"
        class="overflow-auto no-scrollbar h-[calc(100svh-20em)] border border-gray-600 flex flex-col gap-1 p-6">
        <Message v-for="message in messages" :key="message.id" :message="message" :me="message.userId === user.id" />
      </div>
      <div class="relative">
        <textarea v-model="userinput" rows="4"
          class="text-sm border border-gray-600 border-1 mt-4 w-full bg-neutral-900 px-2 py-1 resize-none"
          type="text"></textarea>
        <svg xmlns="http://www.w3.org/2000/svg" @click.stop="send" width="24" height="24" viewBox="0 0 24 24"
          fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
          class="icon hover:cursor-pointer absolute bottom-4 right-3 icon-tabler icons-tabler-outline icon-tabler-send">
          <path stroke="none" d="M0 0h24v24H0z" fill="none" />
          <path d="M10 14l11 -11" />
          <path d="M21 3l-6.5 18a.55 .55 0 0 1 -1 0l-3.5 -7l-7 -3.5a.55 .55 0 0 1 0 -1l18 -6.5" />
        </svg>
      </div>
    </div>
    <template #aside>
    </template>
  </TwoColumnLayout>
</template>

<style scoped lang="css">
.no-scrollbar::-webkit-scrollbar {
  display: none;
}

.no-scrollbar {
  -ms-overflow-style: none;
  /* IE and Edge */
  scrollbar-width: none;
  /* Firefox */
}
</style>

