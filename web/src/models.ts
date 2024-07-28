export type Player = {
	name: string;
};

export type Game = {
  code: string;
  state: number;
  players: Player[];
};
