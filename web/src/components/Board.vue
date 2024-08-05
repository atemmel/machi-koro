<script setup lang="ts">
import type { Card } from "@/models";
import CardVue from "@/components/Card.vue";
defineProps<{
  cards: Record<number, Card>;
}>();

const emit = defineEmits<{
  (e: "cardClick", id: number): void;
}>();

const click = (idx: number) => {
  emit("cardClick", idx);
};
</script>
<template>
  <div class="grid-container">
    <card-vue
      :card="card"
      class="card"
      v-for="(card, idx) in cards"
      :key="idx"
      @click="click(idx)"
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
