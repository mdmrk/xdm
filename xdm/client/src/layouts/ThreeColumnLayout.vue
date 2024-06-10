<script setup>
import AppNavigation from "@/components/AppNavigation.vue";
import AppFooter from "@/components/AppFooter.vue";
import AppLogo from "@/components/AppLogo.vue";
</script>

<template>
  <div class="three-column-layout">
    <header>
      <AppLogo />
      <AppNavigation />
    </header>

    <main>
      <slot />
    </main>

    <aside>
      <slot name="aside" />
      <AppFooter />
    </aside>
  </div>
</template>

<style scoped>
.three-column-layout {
  display: grid;
  grid-template-columns: 1fr 3fr 1fr;
  grid-template-areas: "header main aside";
}

.three-column-layout>header {
  grid-area: header;
  @apply flex flex-col ml-7 mt-7;
}

.three-column-layout>main {
  grid-area: main;
  margin-top: 10px;
  padding: 20px;
}

.three-column-layout>aside {
  grid-area: aside;
  margin-top: 10px;
  padding: 20px;
}
</style>
