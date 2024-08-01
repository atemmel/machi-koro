<script setup lang="ts">
import { useStore } from "@/store";
import { ref, computed, onMounted } from "vue";
import { join, start } from "@/sockets";
import { useRouter } from "vue-router";
import Board from "@/components/Board.vue";
import Die from "@/components/Die.vue";
import type { Card } from "@/models";

const router = useRouter();
const store = useStore();
const game = computed(() => store.game);
const joined = ref(false);
const name = ref("");

onMounted(() => {
  if (store.game == undefined) {
    router.replace("/");
  }
});

const onSubmit = () => {
  if (store.game) {
    store.player = name.value;
    join(store.game.code, name.value).then(() => {
      joined.value = true;
    });
  }
};

const onStart = () => {
  start();
};

const cards: Card[] = [
  {
    name: "Wheat Field",
    description: "Get 1 coin from the bank. (anyone's turn)",
    purchaseCost: 1,
    dieTriggers: [1],
    suit: "blue",
    archetype: "food",
    effect: "gain (any)",
    effectSource: "",
  },
  {
    name: "Ranch",
    description: "Get 1 coin from the bank. (anyone's turn)",
    purchaseCost: 1,
    dieTriggers: [1],
    suit: "blue",
    archetype: "food",
    effect: "gain (any)",
    effectSource: "",
  },
];
</script>
<template>
  <div v-if="store.inLobby">
    <form @submit.prevent="onSubmit" v-if="!joined">
      <input v-model="name" />
      <input type="submit" />
    </form>
    <button v-else-if="store.isOwner" @click="onStart">Start game</button>
    <ul>
      <li v-for="(p, idx) of game?.players" :key="idx">{{ p }}</li>
    </ul>
  </div>
  <div v-else="store.inProgress">
    <h1>game is in progress, current player {{ store.activePlayer }}</h1>
    <template v-if="store.isMyTurn">
      <die v-if="store.isRollPhase" :eyes="[]" />
      <board :cards="cards" v-if="store.isBuyPhase" />
    </template>
  </div>
</template>
