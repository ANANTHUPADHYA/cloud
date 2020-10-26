import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FilesComponent } from './files/files.component';
import { MaterialModule } from '../material/material.module';
import { AdminRoutes } from './admin.routes';
import { RouterModule } from '@angular/router';



@NgModule({
  declarations: [FilesComponent],
  imports: [
    CommonModule,
    MaterialModule,
    RouterModule.forChild(AdminRoutes)
  ]
})
export class AdminModule { }
