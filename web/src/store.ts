import { defineStore } from "pinia";
import { ref, computed } from "vue";
import type { Eyes } from "@/models";
import { InProgress, Lobby } from "@/models";
import type { Card, Game } from "@/models";
import { get } from "./api";

export const useStore = defineStore("store", () => {
  const game = ref<Game | undefined>();
  const player = ref("");
  const owner = ref(false);
  const activePlayer = ref("");
  const phase = ref<"roll" | "buy">("roll");
  const eyes = ref<Eyes[]>([1]);
  const availableCards = ref<Card[]>([]);
  const boughtCards = ref<Record<number, Card>>({});

  const sumEyes = computed(() => {
    let sum = 0;
    for (const x of eyes.value) {
      sum += x;
    }
    return sum;
  });
  const nDie = ref(1);

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

  const getAvailableCards = () => {
    get<Card[]>("/api/cards").then((c) => {
      availableCards.value = c;
    });
  };

  return {
    activePlayer,
    availableCards,
    boughtCards,
    eyes,
    game,
    getAvailableCards,
    inLobby,
    inProgress,
    isBuyPhase,
    isMyTurn,
    isOwner,
    isRollPhase,
    nDie,
    owner,
    phase,
    player,
    sumEyes,
  };
});
