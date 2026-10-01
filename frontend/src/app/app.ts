import { Component, inject } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { NavigationEnd, Router, RouterOutlet } from '@angular/router';
import { filter, map } from 'rxjs';
import { ReservationBar } from './layout/reservation-bar/reservation-bar';
import { ReservationModal } from './layout/reservation-modal/reservation-modal';
import { SiteFooter } from './layout/site-footer/site-footer';
import { SiteNav } from './layout/site-nav/site-nav';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, SiteNav, SiteFooter, ReservationBar, ReservationModal],
  templateUrl: './app.html',
  styleUrl: './app.css',
})
export class App {
  private readonly router = inject(Router);

  readonly showBar = toSignal(
    this.router.events.pipe(
      filter((event): event is NavigationEnd => event instanceof NavigationEnd),
      map(() => !this.router.url.startsWith('/login')),
    ),
    { initialValue: true },
  );
}
