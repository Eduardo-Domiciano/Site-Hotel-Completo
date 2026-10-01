import { Routes } from '@angular/router';
import { authGuard } from './core/guards/auth.guard';
import { About } from './pages/about/about';
import { Career } from './pages/career/career';
import { Home } from './pages/home/home';
import { Login } from './pages/login/login';
import { Profile } from './pages/profile/profile';
import { Rooms } from './pages/rooms/rooms';

export const routes: Routes = [
  { path: '', component: Home, title: 'Mikhalateia Hotel' },
  { path: 'sobre', component: About, title: 'Sobre | Mikhalateia Hotel' },
  { path: 'acomodacoes', component: Rooms, title: 'Acomodações | Mikhalateia Hotel' },
  { path: 'login', component: Login, title: 'Login | Mikhalateia Hotel' },
  {
    path: 'perfil',
    component: Profile,
    canActivate: [authGuard],
    title: 'Perfil | Mikhalateia Hotel',
  },
  { path: 'trabalhe-conosco', component: Career, title: 'Trabalhe Conosco | Mikhalateia Hotel' },
  { path: '**', redirectTo: '' },
];
