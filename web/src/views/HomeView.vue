<script setup lang="ts">
import { RouterLink, useRouter } from "vue-router";
import { ref } from "vue";
import { get, post } from "@/api";
import { useStore } from "@/store";
import type { Game } from "@/models";

const router = useRouter();
const store = useStore();

const code = ref("");
const loading = ref(false);

const enter = (g: Game) => {
  store.game = g;
  router.push("/game");
};

const onJoin = () => {
  get<Game>(`/api/games/${code.value}`)
    .then((g) => {
      enter(g);
    })
    .catch((s) => {
      console.log("no game", s);
      if (s == 404) {
        //TODO: handle...
      }
    });
};

const onHost = () => {
  post<Game>(`/api/games`)
    .then((g) => {
      enter(g);
    })
    .catch((s) => {
      console.log("no game", s);
    });
};
</script>

<template>
  <main>
    <h1>homeview</h1>
    <div v-if="!loading">
      <!-- TODO ska vara form -->
      <input v-model="code" />
      <a href="javascript:;" @click="onJoin"><div>Join</div></a>
      <a href="javascript:;" @click="onHost"><div>Host (as Player)</div></a>
      <router-link to=""><div>Host (as Spectator)</div></router-link>
    </div>
  </main>
</template>
