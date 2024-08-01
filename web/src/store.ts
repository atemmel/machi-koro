import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { InProgress, Lobby } from "@/models";
import type { Game } from "@/models";

export const useStore = defineStore("store", () => {
  const game = ref<Game | undefined>();
  const player = ref("");
  const owner = ref(false);
  const activePlayer = ref("");
  const phase = ref<"roll" | "buy">("roll");

  const inProgress = computed((): boolean => {
    return game.value?.state == InProgress;
  });

  const inLobby = computed((): boolean => {
    return game.value?.state == Lobby;
  });

  const isOwner = computed((): boolean => {
    return owner.value;
  });

  const isMyTurn = computed((): boolean => {
    return activePlayer.value == player.value;
  });

  const isRollPhase = computed((): boolean => {
    return phase.value == "roll";
  });

  const isBuyPhase = computed((): boolean => {
    return phase.value == "buy";
  });

  return {
    activePlayer,
    game,
    inLobby,
    inProgress,
    isBuyPhase,
    isMyTurn,
    isOwner,
    isRollPhase,
    owner,
    phase,
    player,
  };
});
