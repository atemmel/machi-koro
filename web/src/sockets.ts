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
} from "@/models";
import { useStore } from "@/store";

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
          store.activePlayer = msg.player;
          store.phase = "roll";
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

const joinMessage = (code: string, player: string): string => {
  const openMessage: ClientRequest = {
    code,
    operand: BlankOperand,
    player,
    requestOperation: JoinOperation,
  };
  return JSON.stringify(openMessage);
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
