import { Routes } from '@angular/router';
import { UserComponent } from './user/user.component';
import { UploadComponent } from './upload/upload.component';

export const UserRoutes: Routes = [
  {path: 'upload', component: UploadComponent},
  {path: '', component: UserComponent},
];

