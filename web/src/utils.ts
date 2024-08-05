import type { Eyes } from "./models";

export const randomDiceRoll = (): Eyes => {
  return (Math.floor(Math.random() * 6) + 1) as Eyes;
};

export const someOtherDice = (eyes: Eyes): Eyes => {
  if (eyes == 1) {
    return 6;
  }
  return eyes - 1 as Eyes;
};
