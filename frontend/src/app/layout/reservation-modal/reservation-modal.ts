import { CurrencyPipe } from '@angular/common';
import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { reservationPolicies } from '../../core/data/hotel';
import { nightlyFrom, nightlyPrice, rooms } from '../../core/data/rooms';
import { Room } from '../../core/models/room';
import { ReservationService } from '../../core/services/reservation';

@Component({
  selector: 'app-reservation-modal',
  imports: [FormsModule, CurrencyPipe],
  templateUrl: './reservation-modal.html',
  styleUrl: './reservation-modal.css',
})
export class ReservationModal {
  readonly reservation = inject(ReservationService);
  readonly rooms = rooms;
  readonly policies = reservationPolicies;
  readonly priceOf = nightlyPrice;
  readonly fromOf = nightlyFrom;

  get nameModel(): string {
    return this.reservation.guestName();
  }
  set nameModel(value: string) {
    this.reservation.guestName.set(value);
  }

  get phoneModel(): string {
    return this.reservation.guestPhone();
  }
  set phoneModel(value: string) {
    this.reservation.guestPhone.set(value);
  }

  get emailModel(): string {
    return this.reservation.guestEmail();
  }
  set emailModel(value: string) {
    this.reservation.guestEmail.set(value);
  }

  get checkInModel(): string {
    return this.reservation.checkIn();
  }
  set checkInModel(value: string) {
    this.reservation.checkIn.set(value);
  }

  get checkOutModel(): string {
    return this.reservation.checkOut();
  }
  set checkOutModel(value: string) {
    this.reservation.checkOut.set(value);
  }

  add(room: Room): void {
    this.reservation.addRoom(room);
  }

  close(): void {
    this.reservation.closeModal();
  }

  confirm(): void {
    this.reservation.confirm();
  }

  occupancy(uid: string, field: 'adults' | 'children', value: string | number): void {
    this.reservation.updateOccupancy(uid, field, Number(value));
  }

  onBackdrop(event: MouseEvent): void {
    if (event.target === event.currentTarget) {
      this.close();
    }
  }
}
