import { Component, OnInit, ViewChild } from '@angular/core';
import { MatTableDataSource } from '@angular/material/table';
import { Subscription } from 'rxjs';
import { MatSort } from '@angular/material/sort';
import { UserService } from 'src/app/core/services/users/user.service';
import { Router } from '@angular/router';
import { MatSnackBar } from '@angular/material/snack-bar';

@Component({
  selector: 'app-files',
  templateUrl: './files.component.html',
  styleUrls: ['./files.component.scss']
})
export class FilesComponent implements OnInit {

  public displayedColumns: string[] = ['fileName', 'description','firstname', 'lastname','modifiedDate', 'uploadedDate', 'delete' ];
  public dataSource = new MatTableDataSource([]);
  private subscriptions = new Subscription();
  public email: string;
  public userDetails;
  public showSpinner = false;
   @ViewChild(MatSort, {static: false}) sort: MatSort;
  constructor(private userService: UserService, private router: Router, private snackBar: MatSnackBar) { }

  getUserDetails() {
    this.showSpinner = true;
    this.userService.getAllUsers().subscribe(data => {
      if(data.success) {

        this.userDetails =  JSON.parse(sessionStorage.getItem('userDetails'));

       
        this.userService.userDetails.next(true);
        data.users.forEach(element => {
          
          element.files.forEach(file => {
            const obj = {
              firstname: element.firstname,
              lastname: element.lastname,
            }
            obj['fileName'] = file.fileName;
            obj['modifiedDate']= file.modifiedDate;
            obj['uploadedDate'] = file.uploadedDate;
            obj['downloadUrl'] = file.downloadUrl;
            obj['description'] = file.description;
            obj['email'] = element.username;
            this.dataSource.data.push(obj);
          });

        });
        
        this.dataSource.sort = this.sort;
        this.showSpinner = false;
      }
    });
  }

  deleteFile(email, filename) {
    this.showSpinner = true;
    this.userService.deleteFile(email, filename).subscribe(response => {
      this.openSnackBar(response.Message, 'mat-warn');
    this.dataSource.data = [];
      this.getUserDetails();
    })
  }

  openSnackBar(message: string, className: string) {
    this.snackBar.open(message, '', {
      duration: 5000,
      panelClass: ['mat-toolbar', className]
    });
  }

  ngOnInit(): void {
    this.getUserDetails();
  }
}
