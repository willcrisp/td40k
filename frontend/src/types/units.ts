export type FootprintShape = "round" | "oval" | "hull";
export type UnitStatus = "alive" | "in_reserves" | "dead";

export interface UnitStats {
  movement: number;
  toughness: number;
  save: number;
  invulnerable_save?: number;
  wounds: number;
  leadership: number;
  objective_control: number;
}

export interface Footprint {
  x: number;
  y: number;
  has_base: boolean;
}

export interface BoardPosition {
  x: number;
  y: number;
  facing: number;
}

export interface Unit {
  name: string;
  faction: string;
  keywords: string[];
  stats: UnitStats;
  footprint: Footprint;
  position: BoardPosition;
  status: UnitStatus;
  current_wounds: number;
}
