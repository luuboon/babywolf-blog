import { Injectable } from '@angular/core';
import { CanActivate, Router, UrlTree } from '@angular/router';
import { Observable, from } from 'rxjs';
import { map } from 'rxjs/operators';
import { AuthService } from '../services/auth.service';

@Injectable({
  providedIn: 'root'
})
export class AuthGuard implements CanActivate {
  constructor(private authService: AuthService, private router: Router) {}

  canActivate(): Observable<boolean | UrlTree> {
    return from(this.authService.checkIsAdmin()).pipe(
      map((isAdmin) => {
        if (isAdmin) {
          return true;
        }
        return this.router.parseUrl('/login');
      })
    );
  }
}
