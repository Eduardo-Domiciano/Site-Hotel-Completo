export interface Room {
  id: string;
  name: string;
  description: string;
  image: string;
  maxAdults: number;
  maxChildren: number;
  prices: {
    1: { from: number; to: number };
    2: { from: number; to: number };
    3?: { from: number; to: number };
    4?: { from: number; to: number };
  };
}

export interface SelectedRoom {
  uid: string;
  room: Room;
  adults: number;
  children: number;
}

export type Occupancy = 1 | 2 | 3 | 4;
