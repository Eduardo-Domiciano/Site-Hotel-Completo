import { computed, Injectable, inject, signal } from '@angular/core';
import { nightlyPrice } from '../data/rooms';
import { Room, SelectedRoom } from '../models/room';
import { AuthService } from './auth';
import { ApiService } from './api';

@Injectable({ providedIn: 'root' })
export class ReservationService {
  private readonly api = inject(ApiService);
  private readonly auth = inject(AuthService);

  readonly open = signal(false);
  readonly checkIn = signal('');
  readonly checkOut = signal('');
  readonly adults = signal(1);
  readonly children = signal(0);
  readonly guestName = signal('');
  readonly guestPhone = signal('');
  readonly guestEmail = signal('');
  readonly selectedRooms = signal<SelectedRoom[]>([]);
  readonly confirmation = signal<string | null>(null);
  readonly error = signal<string | null>(null);
  readonly submitting = signal(false);

  readonly nights = computed(() => {
    const start = this.checkIn();
    const end = this.checkOut();
    if (!start || !end) return 1;
    const ms = new Date(end).getTime() - new Date(start).getTime();
    const days = Math.round(ms / 86_400_000);
    return Math.max(1, days);
  });

  readonly total = computed(() => {
    const nights = this.nights();
    return this.selectedRooms().reduce(
      (sum, item) => sum + nightlyPrice(item.room, item.adults) * nights,
      0,
    );
  });

  openModal(room?: Room): void {
    this.confirmation.set(null);
    this.error.set(null);
    if (room) {
      this.addRoom(room);
    }
    this.open.set(true);
  }

  closeModal(): void {
    this.open.set(false);
  }

  addRoom(room: Room): void {
    if (this.selectedRooms().length >= 3) return;
    const already = this.selectedRooms().some((item) => item.room.id === room.id);
    if (already) return;
    this.selectedRooms.update((list) => [
      ...list,
      {
        uid: `${room.id}-${Date.now()}`,
        room,
        adults: Math.min(this.adults(), room.maxAdults),
        children: Math.min(this.children(), room.maxChildren),
      },
    ]);
  }

  removeRoom(uid: string): void {
    this.selectedRooms.update((list) => list.filter((item) => item.uid !== uid));
  }

  updateOccupancy(uid: string, field: 'adults' | 'children', value: number): void {
    this.selectedRooms.update((list) =>
      list.map((item) => {
        if (item.uid !== uid) return item;
        const max = field === 'adults' ? item.room.maxAdults : item.room.maxChildren;
        return { ...item, [field]: Math.min(max, Math.max(field === 'adults' ? 1 : 0, value)) };
      }),
    );
  }

  confirm(): void {
    this.error.set(null);
    if (!this.guestName() || !this.guestPhone() || !this.guestEmail()) {
      this.error.set('Preencha nome, telefone e e-mail.');
      return;
    }
    if (!this.checkIn() || !this.checkOut()) {
      this.error.set('Informe check-in e check-out.');
      return;
    }
    if (this.selectedRooms().length === 0) {
      this.error.set('Escolha pelo menos um quarto.');
      return;
    }

    const body = {
      guestName: this.guestName(),
      guestPhone: this.guestPhone(),
      guestEmail: this.guestEmail(),
      checkIn: this.checkIn(),
      checkOut: this.checkOut(),
      rooms: this.selectedRooms().map((item) => ({
        roomTypeId: item.room.id,
        adults: item.adults,
        children: item.children,
      })),
    };

    this.submitting.set(true);
    const asGuest = !this.auth.isLoggedIn();
    this.api.createReservation(body, asGuest).subscribe({
      next: (res) => {
        this.submitting.set(false);
        this.confirmation.set(
          `Reserva ${res.id} confirmada para ${res.guestName}. Total: R$ ${res.total}.`,
        );
      },
      error: (err) => {
        this.submitting.set(false);
        this.error.set(err?.error?.error?.message ?? 'Não foi possível criar a reserva.');
      },
    });
  }

  resetBooking(): void {
    this.selectedRooms.set([]);
    this.guestName.set('');
    this.guestPhone.set('');
    this.guestEmail.set('');
    this.confirmation.set(null);
  }
}
