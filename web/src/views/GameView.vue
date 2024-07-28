<script setup lang="ts">
import { useStore } from "@/store";
import { ref } from "vue";
import { computed } from "vue";
import { onMounted } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
const store = useStore();

const loading = ref(false);
const game = computed(() => store.game);

onMounted(() => {
  if (store.game == undefined) {
    router.replace("/");
  } else {
  }
});
</script>
<template>
  <ul v-if="!loading">
    <li v-for="(p, idx) of game?.players" :key="idx">
      {{ p.name }}
    </li>
  </ul>
</template>
