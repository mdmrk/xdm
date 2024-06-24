<template>
  <p class="text-sm bg-neutral-900 px-2 py-1 w-auto max-w-56 break-words"
    :class="{ 'self-end': me, 'self-start': !me, 'text-right': me, 'text-left': !me }">
    {{ message.Body }}
  </p>
</template>

<script lang="ts">
import { RouterLink } from "vue-router";
import type { IMessage } from "@/types";
import { getUserById } from "@/user";
import { likePost, unlikePost } from "@/post";
import type { PropType } from "vue";

export default {
	name: "Message",
	props: {
		message: {
			type: Object as PropType<IMessage>,
			required: true,
		},
		me: {
			type: Boolean,
			required: false,
			default: true,
		},
	},
};
</script>
