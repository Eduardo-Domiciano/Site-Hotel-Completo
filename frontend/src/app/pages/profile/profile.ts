import { CurrencyPipe } from '@angular/common';
import { Component, HostListener, OnInit, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { rooms } from '../../core/data/rooms';
import { PageHeader } from '../../layout/page-header/page-header';
import { AuthService } from '../../core/services/auth';
import { ApiReservation, ApiService } from '../../core/services/api';

interface ReservationRoomDetail {
  roomTypeId: string;
  name: string;
  description: string;
  image: string;
  adults: number;
  children: number;
  nightlyPrice: string | number;
  lineTotal: string | number;
}

interface ReservationDetail {
  id: string;
  stay: string;
  nightsLabel: string;
  roomsLabel: string;
  total: string | number;
  statusLabel: string;
  checkInLabel: string;
  checkOutLabel: string;
  nights: number;
  rooms: ReservationRoomDetail[];
}

@Component({
  selector: 'app-profile',
  imports: [PageHeader, FormsModule, CurrencyPipe],
  templateUrl: './profile.html',
  styleUrl: './profile.css',
})
export class Profile implements OnInit {
  private readonly auth = inject(AuthService);
  private readonly api = inject(ApiService);

  readonly loading = signal(true);
  readonly selected = signal<ReservationDetail | null>(null);
  message = '';
  error = '';

  form = { name: '', email: '', phone: '', birthDate: '' };
  password = { current: '', next: '', confirm: '' };
  readonly reservations = signal<ReservationDetail[]>([]);

  ngOnInit(): void {
    this.auth.me().subscribe({
      next: (user) => {
        this.form = {
          name: user.name,
          email: user.email,
          phone: user.phone ?? '',
          birthDate: user.birthDate ?? '',
        };
        this.loading.set(false);
      },
      error: () => {
        this.error = 'Não foi possível carregar o perfil.';
        this.loading.set(false);
      },
    });
    this.api.listReservations().subscribe({
      next: (list) => {
        this.reservations.set(list.map((r) => this.toDetail(r)));
      },
      error: () => undefined,
    });
  }

  saveProfile(): void {
    this.message = '';
    this.error = '';
    this.auth
      .updateProfile({
        name: this.form.name,
        phone: this.form.phone || undefined,
        birthDate: this.form.birthDate || undefined,
      })
      .subscribe({
        next: () => {
          this.message = 'Perfil atualizado.';
        },
        error: (err) => {
          this.error = err?.error?.error?.message ?? 'Falha ao atualizar perfil.';
        },
      });
  }

  savePassword(): void {
    this.message = '';
    this.error = '';
    if (this.password.next !== this.password.confirm) {
      this.error = 'As novas senhas não coincidem.';
      return;
    }
    this.auth.changePassword(this.password.current, this.password.next).subscribe({
      next: () => {
        this.message = 'Senha alterada. Faça login novamente se necessário.';
        this.password = { current: '', next: '', confirm: '' };
      },
      error: (err) => {
        this.error = err?.error?.error?.message ?? 'Falha ao alterar senha.';
      },
    });
  }

  logout(): void {
    this.auth.logout();
  }

  openReservation(item: ReservationDetail): void {
    this.selected.set(item);
  }

  closeReservation(): void {
    this.selected.set(null);
  }

  onBackdrop(event: MouseEvent): void {
    if (event.target === event.currentTarget) {
      this.closeReservation();
    }
  }

  @HostListener('document:keydown.escape')
  onEscape(): void {
    this.closeReservation();
  }

  private toDetail(reservation: ApiReservation): ReservationDetail {
    const checkInLabel = this.formatDayAndDate(reservation.checkIn);
    const checkOutLabel = this.formatDayAndDate(reservation.checkOut);
    const roomDetails = reservation.rooms.map((line) => {
      const catalog = rooms.find((room) => room.id === line.roomTypeId);
      return {
        roomTypeId: line.roomTypeId,
        name: catalog?.name ?? line.roomTypeId,
        description: catalog?.description ?? '',
        image: catalog?.image ?? '',
        adults: line.adults,
        children: line.children,
        nightlyPrice: line.nightlyPrice,
        lineTotal: line.lineTotal,
      };
    });
    return {
      id: reservation.id,
      stay: `${this.formatDate(reservation.checkIn)} → ${this.formatDate(reservation.checkOut)}`,
      nightsLabel: reservation.nights === 1 ? '1 noite' : `${reservation.nights} noites`,
      roomsLabel: roomDetails.map((room) => room.name).join(', ') || 'Quarto',
      total: reservation.total,
      statusLabel: this.statusLabel(reservation.status),
      checkInLabel,
      checkOutLabel,
      nights: reservation.nights,
      rooms: roomDetails,
    };
  }

  private formatDayAndDate(iso: string): string {
    const [year, month, day] = iso.slice(0, 10).split('-').map(Number);
    if (!year || !month || !day) return iso;
    const date = new Date(year, month - 1, day);
    const weekday = new Intl.DateTimeFormat('pt-BR', { weekday: 'long' }).format(date);
    const label = weekday.charAt(0).toUpperCase() + weekday.slice(1);
    return `${label}, ${this.formatDate(iso)}`;
  }

  private formatDate(iso: string): string {
    const [year, month, day] = iso.slice(0, 10).split('-');
    if (!year || !month || !day) return iso;
    return `${day}/${month}/${year}`;
  }

  private statusLabel(status: string): string {
    const labels: Record<string, string> = {
      confirmed: 'Confirmada',
      cancelled: 'Cancelada',
      canceled: 'Cancelada',
      pending: 'Pendente',
    };
    return labels[status] ?? status;
  }
}
