import { CurrencyPipe } from '@angular/common';
import { Component, OnInit, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { DestroyRef } from '@angular/core';
import { meetingImages } from '../../core/data/hotel';
import { nightlyFrom, nightlyPrice, rooms as fallbackRooms } from '../../core/data/rooms';
import { Room } from '../../core/models/room';
import { ApiService } from '../../core/services/api';
import { PageHeader } from '../../layout/page-header/page-header';
import { ReservationService } from '../../core/services/reservation';

@Component({
  selector: 'app-rooms',
  imports: [PageHeader, CurrencyPipe, FormsModule],
  templateUrl: './rooms.html',
  styleUrl: './rooms.css',
})
export class Rooms implements OnInit {
  readonly rooms = signal<Room[]>(fallbackRooms);
  readonly adults = signal<Record<string, number>>(
    Object.fromEntries(fallbackRooms.map((room) => [room.id, 1])),
  );
  readonly slide = signal(0);
  readonly slides = meetingImages;
  readonly event = {
    type: '',
    name: '',
    email: '',
    phone: '',
    rooms: 0,
    people: 0,
    date: '',
    coffee: false,
    lunch: false,
    dinner: false,
    sent: false,
    error: '',
  };

  private readonly reservation = inject(ReservationService);
  private readonly api = inject(ApiService);

  constructor() {
    const id = window.setInterval(() => {
      this.slide.update((current) => (current + 1) % this.slides.length);
    }, 2500);
    inject(DestroyRef).onDestroy(() => window.clearInterval(id));
  }

  ngOnInit(): void {
    this.api.listRooms().subscribe({
      next: (list) => {
        const mapped: Room[] = list.map((r) => ({
          id: r.id,
          name: r.name,
          description: r.description,
          image: r.image,
          maxAdults: r.maxAdults,
          maxChildren: r.maxChildren,
          prices: {
            1: {
              from: Number(r.prices['1']?.from ?? 0),
              to: Number(r.prices['1']?.to ?? 0),
            },
            2: {
              from: Number(r.prices['2']?.from ?? 0),
              to: Number(r.prices['2']?.to ?? 0),
            },
            ...(r.prices['3']
              ? {
                  3: {
                    from: Number(r.prices['3'].from),
                    to: Number(r.prices['3'].to),
                  },
                }
              : {}),
            ...(r.prices['4']
              ? {
                  4: {
                    from: Number(r.prices['4'].from),
                    to: Number(r.prices['4'].to),
                  },
                }
              : {}),
          },
        }));
        this.rooms.set(mapped);
        this.adults.set(Object.fromEntries(mapped.map((room) => [room.id, 1])));
      },
      error: () => undefined,
    });
  }

  adultsOf(room: Room): number {
    return this.adults()[room.id] ?? 1;
  }

  setAdults(room: Room, value: string | number): void {
    const next = Math.min(room.maxAdults, Math.max(1, Number(value) || 1));
    this.adults.update((map) => ({ ...map, [room.id]: next }));
  }

  from(room: Room): number {
    return nightlyFrom(room, this.adultsOf(room));
  }

  price(room: Room): number {
    return nightlyPrice(room, this.adultsOf(room));
  }

  select(room: Room): void {
    this.reservation.adults.set(this.adultsOf(room));
    this.reservation.openModal(room);
  }

  requestQuote(): void {
    this.event.error = '';
    this.api
      .createEventQuote({
        eventType: this.event.type,
        name: this.event.name,
        email: this.event.email,
        phone: this.event.phone,
        rooms: this.event.rooms,
        people: this.event.people,
        eventAt: this.event.date ? new Date(this.event.date).toISOString() : '',
        coffee: this.event.coffee,
        lunch: this.event.lunch,
        dinner: this.event.dinner,
      })
      .subscribe({
        next: () => {
          this.event.sent = true;
        },
        error: (err) => {
          this.event.error = err?.error?.error?.message ?? 'Falha ao enviar orçamento.';
        },
      });
  }
}
