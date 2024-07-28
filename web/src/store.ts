import { defineStore } from "pinia";
import { ref } from "vue";
import type { Game } from "./models";

export const useStore = defineStore("store", () => {
  const game = ref<Game | undefined>();
  return {
    game,
  };
});
