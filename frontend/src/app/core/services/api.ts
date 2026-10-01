import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable, map } from 'rxjs';
import { environment } from '../../../environments/environment';
import { ApiRoom } from '../models/api';

export interface ApiReservation {
  id: string;
  guestName: string;
  guestPhone: string;
  guestEmail: string;
  checkIn: string;
  checkOut: string;
  nights: number;
  total: string | number;
  status: string;
  rooms: Array<{
    roomTypeId: string;
    adults: number;
    children: number;
    nightlyPrice: string | number;
    lineTotal: string | number;
  }>;
  createdAt: string;
}

@Injectable({ providedIn: 'root' })
export class ApiService {
  private readonly http = inject(HttpClient);
  private readonly base = environment.apiUrl;

  listRooms(): Observable<ApiRoom[]> {
    return this.http.get<{ data: ApiRoom[] }>(`${this.base}/rooms`).pipe(map((r) => r.data));
  }

  createReservation(
    body: {
      guestName: string;
      guestPhone: string;
      guestEmail: string;
      checkIn: string;
      checkOut: string;
      rooms: Array<{ roomTypeId: string; adults: number; children: number }>;
    },
    asGuest = false,
  ): Observable<ApiReservation> {
    const url = asGuest ? `${this.base}/reservations/guest` : `${this.base}/reservations`;
    return this.http.post<{ data: ApiReservation }>(url, body).pipe(map((r) => r.data));
  }

  listReservations(): Observable<ApiReservation[]> {
    return this.http
      .get<{ data: ApiReservation[] }>(`${this.base}/reservations`)
      .pipe(map((r) => r.data));
  }

  createEventQuote(body: Record<string, unknown>): Observable<unknown> {
    return this.http.post(`${this.base}/event-quotes`, body).pipe(map((r: any) => r.data));
  }

  createJobApplication(form: FormData): Observable<unknown> {
    return this.http.post(`${this.base}/job-applications`, form).pipe(map((r: any) => r.data));
  }
}
