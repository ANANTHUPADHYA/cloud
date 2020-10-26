import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { AuthGuardGuard } from './core/auth-guard/auth-guard.guard';

const routes: Routes = [
  {
    path: 'login',
    loadChildren: () => import('./login/login.module').then(m => m.LoginModule)
  },
  {
    path: 'admin',
    loadChildren:() => import('./admin/admin.module').then(m => m.AdminModule),
    canActivate: [AuthGuardGuard]
  },
  {
    path: 'home',
    loadChildren:() => import('./user/user.module').then(m => m.UserModule),
    canActivate: [AuthGuardGuard]
  },
  { path: '', redirectTo: 'login', pathMatch: 'full' }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
