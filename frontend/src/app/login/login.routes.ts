import { LoginComponent } from './components/login/login.component'
import { Routes } from '@angular/router';
import { SignUpComponent } from './components/sign-up/sign-up.component';

export const LoginRoutes: Routes = [
  {path: '', component: LoginComponent},
  {path: 'register', component: SignUpComponent}
];

