<script setup lang="ts">
import type { Card } from "@/models";
import CardVue from "@/components/Card.vue";
defineProps<{
  cards: Card[];
}>();
</script>
<template>
  <div class="grid-container">
    <card-vue
      :card="card"
      class="card"
      v-for="(card, idx) in cards"
      :key="idx"
    />
  </div>
</template>
<style scoped>
.grid-container {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr 1fr 1fr;
  gap: 0.5rem;
}
</style>
