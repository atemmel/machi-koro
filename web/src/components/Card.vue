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
    <div class="count">
      {{ card.count }}
    </div>
    <div class="center">
      {{ card.name }}
    </div>
  </button>
</template>
<style scoped>
.card {
  border: none;
  padding: 1rem;
  font-size: 2rem;
  text-align: left;
  color: var(--color-text);
  position: relative;
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

.center {
  text-align: center;
}

.count {
  position: absolute;
  background: gray;
  border-radius: 50%;
  display: inline-block;
  width: 2rem;
  line-height: 2rem;
  text-align: center;
  padding: 0.25rem;

  right: 0;
  top: 0;
  margin-right: -1rem;
  margin-top: -1rem;
  z-index: 10;
}
</style>
