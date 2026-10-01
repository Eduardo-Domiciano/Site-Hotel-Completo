import { Occupancy, Room } from '../models/room';

export const rooms: Room[] = [
  {
    id: 'express-duplo',
    name: 'Express Duplo',
    description: 'Modelo com 02 camas de solteiro, frigobar, TV e ar-condicionado.',
    image: '/assets/img/twin-bed.jpg',
    maxAdults: 2,
    maxChildren: 2,
    prices: {
      1: { from: 200, to: 180 },
      2: { from: 269, to: 249 },
    },
  },
  {
    id: 'express-triplo',
    name: 'Express Triplo',
    description: 'Modelo com 03 camas de solteiro, frigobar, TV e ar-condicionado.',
    image: '/assets/img/triplo.jpg',
    maxAdults: 3,
    maxChildren: 2,
    prices: {
      1: { from: 250, to: 210 },
      2: { from: 269, to: 249 },
      3: { from: 349, to: 329 },
    },
  },
  {
    id: 'luxo-casal',
    name: 'Luxo Casal',
    description: 'Modelo com 01 cama de casal, frigobar, TV e ar-condicionado.',
    image: '/assets/img/double-bed.jpg',
    maxAdults: 2,
    maxChildren: 2,
    prices: {
      1: { from: 249, to: 229 },
      2: { from: 299, to: 279 },
    },
  },
  {
    id: 'luxo-familia',
    name: 'Luxo Família',
    description: 'Modelo com 01 cama de casal, 01 cama de solteiro, frigobar, TV e ar-condicionado.',
    image: '/assets/img/double_bed_and_single_bed.jpg',
    maxAdults: 3,
    maxChildren: 2,
    prices: {
      1: { from: 249, to: 229 },
      2: { from: 299, to: 279 },
      3: { from: 369, to: 349 },
    },
  },
  {
    id: 'flat',
    name: 'Flat',
    description:
      'Modelo com 01 cama de casal, 02 camas de solteiro, mini-cozinha, sofá, TV e ar-condicionado.',
    image: '/assets/img/flat.jpg',
    maxAdults: 4,
    maxChildren: 4,
    prices: {
      1: { from: 330, to: 310 },
      2: { from: 410, to: 390 },
      3: { from: 530, to: 510 },
      4: { from: 600, to: 580 },
    },
  },
  {
    id: 'master',
    name: 'Master',
    description: 'Modelo com 01 cama de casal, sofá, TV, ar-condicionado e hidromassagem.',
    image: '/assets/img/master.jpg',
    maxAdults: 2,
    maxChildren: 2,
    prices: {
      1: { from: 400, to: 380 },
      2: { from: 520, to: 500 },
    },
  },
];

export function occupancyOf(adults: number): Occupancy {
  if (adults >= 4) return 4;
  if (adults === 3) return 3;
  if (adults === 2) return 2;
  return 1;
}

export function nightlyPrice(room: Room, adults: number): number {
  const occupancy = occupancyOf(Math.min(adults, room.maxAdults));
  const tier = room.prices[occupancy] ?? room.prices[1];
  return tier.to;
}

export function nightlyFrom(room: Room, adults: number): number {
  const occupancy = occupancyOf(Math.min(adults, room.maxAdults));
  const tier = room.prices[occupancy] ?? room.prices[1];
  return tier.from;
}
