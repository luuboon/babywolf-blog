import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterModule],
  templateUrl: './register.component.html',
  styleUrls: ['./login.component.scss'] // Reusing login styles for simplicity
})
export class RegisterComponent {
  registerForm: FormGroup;
  errorMsg = '';
  successMsg = '';
  loading = false;

  constructor(
    private fb: FormBuilder,
    private auth: AuthService,
    private router: Router
  ) {
    this.registerForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]]
    });
  }

  async onSubmit() {
    if (this.registerForm.invalid) return;

    this.loading = true;
    this.errorMsg = '';
    this.successMsg = '';
    const { email, password } = this.registerForm.value;

    try {
      const { error } = await this.auth.signUp(email, password);
      if (error) {
        this.errorMsg = error.message;
      } else {
        this.successMsg = '¡Registro exitoso! Por favor, verifica tu correo.';
        // Optional: redirect to login after a few seconds
        setTimeout(() => this.router.navigate(['/login']), 3000);
      }
    } catch (err: any) {
      this.errorMsg = err.message || 'Registration failed';
    } finally {
      this.loading = false;
    }
  }
}
