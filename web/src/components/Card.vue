<script setup lang="ts">
import type { Card } from "@/models";
import { computed } from "vue";
const props = defineProps<{
  card: Card;
}>();

const classes = computed(() => {
  return {
    purple: props.card.suit == "purple",
    red: props.card.suit == "red",
    green: props.card.suit == "green",
    blue: props.card.suit == "blue",
  };
});
</script>
<template>
  <button class="card" :class="classes">
    {{ card.name }}
  </button>
</template>
<style scoped>
.card {
  border: none;
  padding: 1rem;
  font-size: 2rem;
  text-align: center;
  color: var(--color-text);
}

.card:hover {
  cursor: pointer;
  outline: 1px solid white;
}

.purple {
  background-color: purple;
}

.red {
  background-color: red;
}

.green {
  background-color: green;
}

.blue {
  background-color: blue;
}
</style>
