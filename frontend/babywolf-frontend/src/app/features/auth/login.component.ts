import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators, FormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { AuthService } from '../../core/services/auth.service';
import { SupabaseService } from '../../core/services/supabase.service';

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
    private router: Router,
    private supabase: SupabaseService
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
        
        // Check if MFA is required using Authenticator Assurance Level
        const { data: aalData, error: aalError } = 
          await this.supabase.client.auth.mfa.getAuthenticatorAssuranceLevel();
        
        if (aalError) {
          console.error('AAL check error:', aalError);
          // Proceed with normal login if AAL check fails
          this.router.navigate(['/']);
          return;
        }

        // If current level is aal1 but next level requires aal2, user needs MFA
        if (aalData.currentLevel === 'aal1' && aalData.nextLevel === 'aal2') {
          // User has verified TOTP factor – needs to enter code
          const { data: mfaData, error: mfaError } = await this.auth.getMfaFactors();
          if (mfaError) throw mfaError;

          const totpFactor = mfaData?.totp?.find((f: any) => f.status === 'verified');
          if (!totpFactor) {
            this.errorMsg = 'MFA requerido pero no hay métodos verificados.';
            this.loading = false;
            return;
          }

          this.mfaFactorId = totpFactor.id;
          this.requiresMfa = true;
        } else {
          // Login successful, no MFA needed or already aal2
          this.router.navigate(['/']);
        }
      } else {
        // Step 2: Verify MFA Code
        if (this.mfaCode.length < 6) {
          this.errorMsg = 'El código debe tener 6 dígitos.';
          this.loading = false;
          return;
        }

        const { error } = await this.auth.verifyMfaLogin(this.mfaFactorId, this.mfaCode);
        if (error) {
          this.errorMsg = 'Código incorrecto. Intenta nuevamente.';
          this.mfaCode = '';
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
