import { Component, OnInit, ViewChild } from '@angular/core';
import { MatTableDataSource } from '@angular/material/table';
import { Subscription } from 'rxjs';
import { MatSort } from '@angular/material/sort';
import { UserService } from 'src/app/core/services/users';
import { Router } from '@angular/router';
import { MatSnackBar } from '@angular/material/snack-bar';
import { FilesService } from 'src/app/core/services/files';

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.scss']
})
export class UserComponent implements OnInit {
  public displayedColumns: string[] = ['file_name', 'description','updated_at','created_at','download', 'edit', 'delete' ];
  public dataSource = new MatTableDataSource();
  private subscriptions = new Subscription();
  public email: string;
  public userDetails;
  public showSpinner = false;
   @ViewChild(MatSort, {static: false}) sort: MatSort;
  constructor(private userService: UserService, private fileService: FilesService, private router: Router, private snackBar: MatSnackBar) { }

  getUserDetails() {
    this.showSpinner = true;
    this.userService.getUserDetails(this.email).subscribe(data => {
      if(data.UserID) {
        this.userDetails = data;
        const files = [];
        for (var key in data.files) {
          if (data.files.hasOwnProperty(key)) {
              files.push(data.files[key])
          }
        }      
        this.userService.userDetails.next(true);
        this.dataSource.data = files;
        this.dataSource.sort = this.sort;
        this.showSpinner = false;
      }
    });
  }

  deleteFile(filename) {
    this.showSpinner = true;
    this.userService.deleteFile(this.email, filename).subscribe(response => {
      this.openSnackBar(response.Message, 'mat-warn');
      this.getUserDetails();
    })
  }

  downloadFile(file){ 
    this.showSpinner = true;
    this.fileService.downloadFile(sessionStorage.getItem('UserID'), file.file_name).subscribe(response => {
      window.open(response.presignedURL);
      this.showSpinner = false;
    })
  }

  openSnackBar(message: string, className: string) {
    this.snackBar.open(message, '', {
      duration: 5000,
      panelClass: ['mat-toolbar', className]
    });
  }

  goToEdit(file) {
    sessionStorage.setItem('file', JSON.stringify(file));
    this.router.navigate(['home/upload',{ name: file.file_name}])
  }
  ngOnInit(): void {
    this.email = sessionStorage.getItem('UserID');
    this.getUserDetails();
  }

}
