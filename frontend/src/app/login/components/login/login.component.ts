import { MatSnackBar } from '@angular/material/snack-bar';
import { LoginService } from '../../../core/services/login';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { Subscription } from 'rxjs';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit, OnDestroy {
  public loginForm: FormGroup;
  private subscriptions = new Subscription();
  constructor(private formBuilder: FormBuilder, private loginService: LoginService, private router: Router, private snackBar: MatSnackBar) {
   this.loginForm = this.formBuilder.group({
      email: ['', Validators.email],
      password: ['', Validators.required]
    });
   }
  ngOnInit() {
    const token = sessionStorage.getItem('token');
    if(token) {
      this.router.navigate(['home']);
      this.loginService.userLoggedIn.next(true);
    }
  }

  openSnackBar(message: string, className: string) {
    this.snackBar.open(message, '', {
      duration: 5000,
      panelClass: ['mat-toolbar', className]
    });
  }


  onSubmit() {
    console.log('value', this.loginForm.value);
    this.subscriptions.add(this.loginService.loginUser(this.loginForm.value).subscribe(response => {
      if (response.EmailAddress) {
       sessionStorage.setItem('EmailAddress', response.EmailAddress);
       sessionStorage.setItem('IsAdmin', JSON.stringify(response.IsAdmin));
       sessionStorage.setItem('UserID', response.UserID);
       this.loginService.userLoggedIn.next(true);
       this.router.navigate(['/home']);
      }
    }, error => {
      this.openSnackBar(error.error.error, 'mat-warn');
    }));
  }

  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  goToRegister() {
    this.router.navigate(['login/register']);
  }
}
