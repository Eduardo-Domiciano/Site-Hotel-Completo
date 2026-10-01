import { Component, inject } from '@angular/core';
import { hotel, presentationItems } from '../../core/data/hotel';
import { ReservationService } from '../../core/services/reservation';

@Component({
  selector: 'app-home',
  templateUrl: './home.html',
  styleUrl: './home.css',
})
export class Home {
  readonly hotel = hotel;
  readonly items = presentationItems;
  private readonly reservation = inject(ReservationService);

  openReservation(): void {
    this.reservation.openModal();
  }
}
