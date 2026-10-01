import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { catchError, switchMap, throwError } from 'rxjs';
import { AuthService } from '../services/auth';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const auth = inject(AuthService);
  const token = auth.accessToken();
  const withAuth =
    token && !req.url.includes('/auth/login') && !req.url.includes('/auth/register')
      ? req.clone({ setHeaders: { Authorization: `Bearer ${token}` } })
      : req;

  return next(withAuth).pipe(
    catchError((err: HttpErrorResponse) => {
      if (err.status !== 401 || req.url.includes('/auth/refresh') || req.url.includes('/auth/login')) {
        return throwError(() => err);
      }
      return auth.refresh().pipe(
        switchMap(() => {
          const retryToken = auth.accessToken();
          const retry = retryToken
            ? req.clone({ setHeaders: { Authorization: `Bearer ${retryToken}` } })
            : req;
          return next(retry);
        }),
        catchError((refreshErr) => {
          auth.logout();
          return throwError(() => refreshErr);
        }),
      );
    }),
  );
};
