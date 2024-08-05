export type Game = {
  code: string;
  state: number;
  players: string[];
};

export type ClientRequest = {
  code: string;
  operand: number;
  requestOperation: number;
  player: string;
};

export type ServerResponse = {
  code: string;
  operands: number[];
  responseCode: number;
  player: string;
};

// Game states
export const Lobby = 0;
export const InProgress = 1;

// Client codes
export const JoinOperation = 0;
export const PongResponse = 1;
export const LeaveRequest = 2;
export const StartRequest = 3;
export const RollRequest = 4;
export const BuyRequest = 5;

// Server codes
export const JoinAnnouncement = 0;
export const PingRequest = 1;
export const LeaveAnnouncement = 2;
export const StartAnnouncement = 3;
export const OwnerAssign = 4;
export const TurnChangeAnnouncement = 5;
export const RollAnnouncement = 6;
export const BuyAnnouncement = 7;

// Operands
export const BlankOperand = 0;

// Cards
export type Suit = "red" | "green" | "blue" | "purple";

export type Archetype =
  | "food" // wheat field, apple orchard
  | "resource" // forest, mine
  | "animal" // ranch
  | "leisure" // cafe, family restaurant
  | "governance" // stadium, tv station, business centre
  | "market" // bakery, convenience store, farmers market
  | "crafting"; // cheese factory, furniture factory

export type Effect =
  | "gain"
  | "gain (any)"
  | "gain (multiplied)"
  | "take"
  | "all"
  | "choose"
  | "exchange";

export type Card = {
  name: string;
  description: string;

  purchaseCost: number;
  dieTriggers: number[];
  suit: Suit;
  archetype: Archetype;
  effect: Effect;
  effectSource: "" | Archetype;
  count: number;
};

export type Eyes = 1 | 2 | 3 | 4 | 5 | 6;
