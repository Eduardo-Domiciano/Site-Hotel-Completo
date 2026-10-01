import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ReservationService } from '../../core/services/reservation';

@Component({
  selector: 'app-reservation-bar',
  imports: [FormsModule],
  templateUrl: './reservation-bar.html',
  styleUrl: './reservation-bar.css',
})
export class ReservationBar {
  readonly reservation = inject(ReservationService);

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

  get adultsModel(): number {
    return this.reservation.adults();
  }
  set adultsModel(value: number) {
    this.reservation.adults.set(Number(value) || 1);
  }

  get childrenModel(): number {
    return this.reservation.children();
  }
  set childrenModel(value: number) {
    this.reservation.children.set(Number(value) || 0);
  }

  open(): void {
    this.reservation.openModal();
  }
}
