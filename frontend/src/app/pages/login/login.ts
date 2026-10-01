import { HttpErrorResponse } from '@angular/common/http';
import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { loginBackgrounds } from '../../core/data/hotel';
import { AuthService } from '../../core/services/auth';
import { PageHeader } from '../../layout/page-header/page-header';

@Component({
  selector: 'app-login',
  imports: [PageHeader, FormsModule],
  templateUrl: './login.html',
  styleUrl: './login.css',
})
export class Login {
  private readonly auth = inject(AuthService);
  private readonly router = inject(Router);

  readonly background =
    loginBackgrounds[Math.floor(Math.random() * loginBackgrounds.length)];
  mode: 'login' | 'signup' = 'login';
  message = '';
  busy = false;

  login = { email: '', password: '' };
  signup = { name: '', birth: '', email: '', password: '', confirm: '' };

  submitLogin(): void {
    if (!this.login.email || !this.login.password) {
      this.message = 'Preencha e-mail e senha.';
      return;
    }
    this.busy = true;
    this.auth.login(this.login.email, this.login.password).subscribe({
      next: () => {
        this.busy = false;
        void this.router.navigateByUrl('/perfil');
      },
      error: (err: HttpErrorResponse) => {
        this.busy = false;
        this.message = err.error?.error?.message ?? 'Falha no login.';
      },
    });
  }

  submitSignup(): void {
    if (this.signup.password !== this.signup.confirm) {
      this.message = 'As senhas não coincidem.';
      return;
    }
    if (!this.signup.name || !this.signup.email || !this.signup.password) {
      this.message = 'Preencha nome, e-mail e senha.';
      return;
    }
    if (this.signup.password.length < 8) {
      this.message = 'A senha deve ter pelo menos 8 caracteres.';
      return;
    }
    this.busy = true;
    this.auth
      .register({
        name: this.signup.name,
        email: this.signup.email,
        password: this.signup.password,
        birthDate: this.signup.birth || undefined,
      })
      .subscribe({
        next: () => {
          this.busy = false;
          void this.router.navigateByUrl('/perfil');
        },
        error: (err: HttpErrorResponse) => {
          this.busy = false;
          this.message = err.error?.error?.message ?? 'Falha no cadastro.';
        },
      });
  }
}
