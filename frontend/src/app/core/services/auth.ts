import { HttpClient } from '@angular/common/http';
import { Injectable, computed, inject, signal } from '@angular/core';
import { Router } from '@angular/router';
import { Observable, map, tap } from 'rxjs';
import { environment } from '../../../environments/environment';
import { ApiUser, AuthPayload, TokenPair } from '../models/api';

const ACCESS_KEY = 'mh_access';
const REFRESH_KEY = 'mh_refresh';
const USER_KEY = 'mh_user';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly http = inject(HttpClient);
  private readonly router = inject(Router);
  private readonly userSignal = signal<ApiUser | null>(this.readUser());

  readonly user = this.userSignal.asReadonly();
  readonly isLoggedIn = computed(() => !!this.userSignal());

  register(body: {
    name: string;
    email: string;
    password: string;
    phone?: string;
    birthDate?: string;
  }): Observable<AuthPayload> {
    return this.http
      .post<{ data: AuthPayload }>(`${environment.apiUrl}/auth/register`, body)
      .pipe(
        map((res) => res.data),
        tap((data) => this.persist(data)),
      );
  }

  login(email: string, password: string): Observable<AuthPayload> {
    return this.http
      .post<{ data: AuthPayload }>(`${environment.apiUrl}/auth/login`, { email, password })
      .pipe(
        map((res) => res.data),
        tap((data) => this.persist(data)),
      );
  }

  refresh(): Observable<AuthPayload> {
    const refreshToken = localStorage.getItem(REFRESH_KEY) ?? '';
    return this.http
      .post<{ data: AuthPayload }>(`${environment.apiUrl}/auth/refresh`, { refreshToken })
      .pipe(
        map((res) => res.data),
        tap((data) => this.persist(data)),
      );
  }

  me(): Observable<ApiUser> {
    return this.http.get<{ data: ApiUser }>(`${environment.apiUrl}/me`).pipe(
      map((res) => res.data),
      tap((user) => {
        localStorage.setItem(USER_KEY, JSON.stringify(user));
        this.userSignal.set(user);
      }),
    );
  }

  updateProfile(body: { name: string; phone?: string; birthDate?: string }): Observable<ApiUser> {
    return this.http.put<{ data: ApiUser }>(`${environment.apiUrl}/me`, body).pipe(
      map((res) => res.data),
      tap((user) => {
        localStorage.setItem(USER_KEY, JSON.stringify(user));
        this.userSignal.set(user);
      }),
    );
  }

  changePassword(currentPassword: string, newPassword: string): Observable<void> {
    return this.http
      .put<{ data: { ok: boolean } }>(`${environment.apiUrl}/me/password`, {
        currentPassword,
        newPassword,
      })
      .pipe(map(() => undefined));
  }

  logout(): void {
    const refreshToken = localStorage.getItem(REFRESH_KEY);
    if (refreshToken) {
      this.http.post(`${environment.apiUrl}/auth/logout`, { refreshToken }).subscribe({
        error: () => undefined,
      });
    }
    localStorage.removeItem(ACCESS_KEY);
    localStorage.removeItem(REFRESH_KEY);
    localStorage.removeItem(USER_KEY);
    this.userSignal.set(null);
    void this.router.navigateByUrl('/login');
  }

  accessToken(): string | null {
    return localStorage.getItem(ACCESS_KEY);
  }

  private persist(data: AuthPayload): void {
    localStorage.setItem(ACCESS_KEY, data.tokens.accessToken);
    localStorage.setItem(REFRESH_KEY, data.tokens.refreshToken);
    localStorage.setItem(USER_KEY, JSON.stringify(data.user));
    this.userSignal.set(data.user);
  }

  private readUser(): ApiUser | null {
    const raw = localStorage.getItem(USER_KEY);
    if (!raw) return null;
    try {
      return JSON.parse(raw) as ApiUser;
    } catch {
      return null;
    }
  }
}

export type { TokenPair };
