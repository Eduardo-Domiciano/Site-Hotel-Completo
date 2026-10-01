import { Component, computed, inject, signal } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { hotel } from '../../core/data/hotel';
import { AuthService } from '../../core/services/auth';

@Component({
  selector: 'app-site-nav',
  imports: [RouterLink, RouterLinkActive],
  templateUrl: './site-nav.html',
  styleUrl: './site-nav.css',
})
export class SiteNav {
  readonly hotel = hotel;
  readonly menuOpen = signal(false);
  private readonly auth = inject(AuthService);
  readonly loggedIn = this.auth.isLoggedIn;

  readonly links = computed(() => {
    const base = [
      { path: '/', label: 'Home', exact: true },
      { path: '/sobre', label: 'Sobre', exact: false },
      { path: '/acomodacoes', label: 'Acomodações', exact: false },
      { path: '/trabalhe-conosco', label: 'Trabalhe Conosco', exact: false },
    ];
    if (this.loggedIn()) {
      return [...base, { path: '/perfil', label: 'Perfil', exact: false }];
    }
    return [...base, { path: '/login', label: 'Login', exact: false }];
  });

  toggleMenu(): void {
    this.menuOpen.update((open) => !open);
  }

  closeMenu(): void {
    this.menuOpen.set(false);
  }
}
