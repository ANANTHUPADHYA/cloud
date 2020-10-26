import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { UserRoutes } from './user.routes';
import { UploadComponent } from './upload/upload.component';
import { MaterialModule } from '../material/material.module';
import { MatTableModule } from '@angular/material/table';
import { UserComponent } from './user/user.component';



@NgModule({
  declarations: [UploadComponent, UserComponent],
  imports: [
    CommonModule,
    RouterModule.forChild(UserRoutes),
    MaterialModule,
    MatTableModule
  ]
})
export class UserModule { }
