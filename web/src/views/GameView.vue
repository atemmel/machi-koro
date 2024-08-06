<script setup lang="ts">
import { useStore } from "@/store";
import { ref, computed, onMounted } from "vue";
import { askForRoll, buyMessage, join, start } from "@/sockets";
import { useRouter } from "vue-router";
import Board from "@/components/Board.vue";
import Die from "@/components/Die.vue";
import { watch } from "vue";

const router = useRouter();
const store = useStore();
const game = computed(() => store.game);
const joined = ref(false);
const name = ref("");
const enableInput = ref(true);

onMounted(() => {
  if (store.game == undefined) {
    router.replace("/");
  } else {
    store.getAvailableCards();
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

const roll = (thisMany: number) => {
  enableInput.value = false;
  askForRoll(thisMany);
};

const cardClick = (idx: number) => {
  enableInput.value = false;
  buyMessage(idx);
};

watch(
  () => store.phase,
  () => {
    if (store.activePlayer != store.player) {
      return;
    }

    enableInput.value = true;
  },
);
</script>
<template>
  <div v-if="store.inLobby">
    <div>Game code {{ store.game?.code }}</div>
    <form @submit.prevent="onSubmit" v-if="!joined">
      <input v-model="name" />
      <input type="submit" />
    </form>
    <button v-else-if="store.isOwner" @click="onStart">Start game</button>
    <ul>
      <li v-for="(p, idx) of game?.players" :key="idx">{{ p }}</li>
    </ul>
  </div>
  <div v-else-if="store.inProgress">
    <h1>game is in progress, current player {{ store.activePlayer }}</h1>
    <div v-if="store.isRollPhase">
      <div>
        <die :eyes="store.eyes" />
      </div>
      <div v-if="store.isMyTurn">
        <button @click="roll(1)" v-if="enableInput">Roll 1 dice</button>
      </div>
    </div>
    <div v-if="store.isBuyPhase">
      <div>Available cards</div>
      <board :cards="store.availableCards" @cardClick="cardClick" />
    </div>
    <div>
      <div>Your cards</div>
      <board :cards="store.boughtCards" />
    </div>
  </div>
</template>
