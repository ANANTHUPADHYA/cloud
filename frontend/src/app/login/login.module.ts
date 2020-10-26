import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LoginComponent } from './components/login/login.component';
import { SignUpComponent } from './components/sign-up/sign-up.component';
import { MaterialModule } from '../material/material.module';
import { RouterModule } from '@angular/router';
import { LoginRoutes } from './login.routes';



@NgModule({
  declarations: [LoginComponent, SignUpComponent],
  imports: [
    CommonModule,
    MaterialModule,
    RouterModule.forChild(LoginRoutes)
  ]
})
export class LoginModule { }
