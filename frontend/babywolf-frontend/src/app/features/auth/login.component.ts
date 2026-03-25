import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators, FormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, FormsModule, RouterModule],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent {
  loginForm: FormGroup;
  errorMsg = '';
  loading = false;
  
  // MFA
  requiresMfa = false;
  mfaCode = '';
  mfaFactorId = '';

  constructor(
    private fb: FormBuilder,
    private auth: AuthService,
    private router: Router
  ) {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', Validators.required]
    });
  }

  async onSubmit() {
    if (this.loginForm.invalid && !this.requiresMfa) return;

    this.loading = true;
    this.errorMsg = '';

    try {
      if (!this.requiresMfa) {
        // Step 1: Login with Email & Password
        const { email, password } = this.loginForm.value;
        const { data, error } = await this.auth.signIn(email, password);
        
        if (error) {
          this.errorMsg = error.message;
          this.loading = false;
          return;
        } 
        
        // Verifica si la sesión está pendiente por 2FA
        if (data?.session === null && data?.user) {
          // Requires MFA
          const { data: mfaData, error: mfaError } = await this.auth.getMfaFactors();
          if (mfaError) throw mfaError;

          const totpFactor = mfaData?.all.find((f: any) => f.status === 'verified' && f.factor_type === 'totp');
          if (!totpFactor) {
            this.errorMsg = 'MFA requerido pero no hay métodos verificados.';
            this.loading = false;
            return;
          }

          this.mfaFactorId = totpFactor.id;
          this.requiresMfa = true;
        } else {
          // Login normal
          this.router.navigate(['/']);
        }
      } else {
        // Step 2: Login with MFA Code
        if (this.mfaCode.length < 6) {
          this.errorMsg = 'El código debe tener 6 dígitos.';
          this.loading = false;
          return;
        }

        const { error } = await this.auth.verifyMfaLogin(this.mfaFactorId, this.mfaCode);
        if (error) {
          this.errorMsg = 'Código incorrecto. Intenta nuevamente.';
        } else {
          this.router.navigate(['/']);
        }
      }
    } catch (err: any) {
      this.errorMsg = err.message || 'Login failed';
    } finally {
      this.loading = false;
    }
  }
}
