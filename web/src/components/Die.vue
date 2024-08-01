<script setup lang="ts">
import Dice from "@/components/Dice.vue";
import type { Eyes } from "@/models";
defineProps<{
  eyes: Eyes[];
}>();
</script>
<template>
  <dice :id="'dice-' + idx" :eyes="e" v-for="(e, idx) in eyes" />
</template>
