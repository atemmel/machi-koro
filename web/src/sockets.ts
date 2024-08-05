import type { Eyes } from "@/models";
import type { ClientRequest, ServerResponse } from "@/models";
import {
  JoinOperation,
  BlankOperand,
  PongResponse,
  PingRequest,
  JoinAnnouncement,
  StartRequest,
  InProgress,
  LeaveAnnouncement,
  StartAnnouncement,
  OwnerAssign,
  TurnChangeAnnouncement,
  RollRequest,
  RollAnnouncement,
  BuyAnnouncement,
  BuyRequest,
} from "@/models";
import { useStore } from "@/store";
import { toRaw } from "vue";
import { someOtherDice } from "./utils";

let socket: WebSocket;

export const join = (code: string, player: string): Promise<void> => {
  return new Promise((resolve, reject) => {
    const store = useStore();
    // Create WebSocket connection.
    socket = new WebSocket("/api/ws");

    // Connection opened
    socket.addEventListener("open", () => {
      socket.send(joinMessage(code, player));
      resolve();
    });

    // Listen for messages
    socket.addEventListener("message", (event) => {
      const msg = JSON.parse(event.data) as ServerResponse;
      console.log("Message from server ", msg);
      if (store.game == undefined) {
        return;
      }
      switch (msg.responseCode) {
        case PingRequest: {
          const code = store.game?.code || "";
          const player = store.player;
          socket.send(pongMessage(code, player));
          break;
        }
        case JoinAnnouncement: {
          store.game.players.push(msg.player);
          break;
        }
        case LeaveAnnouncement: {
          store.game.players = store.game.players.filter((p) => {
            return p != msg.player;
          });
          break;
        }
        case StartAnnouncement: {
          store.game.state = InProgress;
          break;
        }
        case OwnerAssign: {
          store.owner = true;
          break;
        }
        case TurnChangeAnnouncement: {
          console.log("Active player is now", msg.player);
          store.activePlayer = msg.player;
		  setTimeout(() => {
          store.phase = "roll";
		  });
          break;
        }
        case RollAnnouncement: {
          if (store.eyes[0] == msg.operands[0]) {
            store.eyes[0] = someOtherDice(store.eyes[0]);
          }
          setTimeout(() => {
            store.eyes = msg.operands as Eyes[];
            setTimeout(() => {
              store.phase = "buy";
            }, 2800);
          }, 200);
          break;
        }
        case BuyAnnouncement: {
          const cardIdx = msg.operands[0];
          store.availableCards[cardIdx].count -= 1;
          console.log(msg.player, store.player);
          if (msg.player == store.player) {
            if (store.boughtCards[cardIdx] == undefined) {
              store.boughtCards[cardIdx] = JSON.parse(
                JSON.stringify(toRaw(store.availableCards[cardIdx])),
              );
              store.boughtCards[cardIdx].count = 0;
            }
            store.boughtCards[cardIdx].count += 1;
          }
          break;
        }
      }
    });

    socket.addEventListener("close", () => {
      console.log("The connection has been closed successfully.");
    });

    socket.addEventListener("error", (event) => {
      console.log("WebSocket error: ", event);
      reject();
    });
  });
};

export const start = () => {
  const store = useStore();
  const code = store.game?.code || "";
  const player = store.player;

  const startMessage: ClientRequest = {
    player,
    code,
    requestOperation: StartRequest,
    operand: BlankOperand,
  };
  socket.send(JSON.stringify(startMessage));
};

export const askForRoll = (n: number) => {
  const store = useStore();
  const code = store.game?.code || "";
  const player = store.player;

  const rollMessage: ClientRequest = {
    player,
    code,
    requestOperation: RollRequest,
    operand: n,
  };
  socket.send(JSON.stringify(rollMessage));
};

const joinMessage = (code: string, player: string): string => {
  const openMessage: ClientRequest = {
    code,
    operand: BlankOperand,
    player,
    requestOperation: JoinOperation,
  };
  return JSON.stringify(openMessage);
};

export const buyMessage = (idx: number) => {
  const store = useStore();
  const code = store.game?.code || "";
  const player = store.player;

  const rollMessage: ClientRequest = {
    player,
    code,
    requestOperation: BuyRequest,
    operand: idx,
  };
  socket.send(JSON.stringify(rollMessage));
};

const pongMessage = (code: string, player: string): string => {
  const resp: ClientRequest = {
    code,
    player,
    operand: BlankOperand,
    requestOperation: PongResponse,
  };
  return JSON.stringify(resp);
};
