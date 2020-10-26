import { Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { LoginService } from '../../services';
import { UserService } from '../../services/users';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnInit, OnDestroy {
  private subscriptions = new Subscription();
  public userEmail: string;
  public showLogOut: boolean = false;
  public isAdmin = false;
  constructor(
    private router: Router,
    private snackBar: MatSnackBar,
    private loginService: LoginService,
    private userService: UserService
  ) {}

  ngOnInit() {
    this.loginService.userLoggedIn.subscribe(data => {
      if(data) {
        this.showLogOut = true;
      }
    });

    this.userService.userDetails.subscribe(data => {
      if(data) {
        this.isAdmin = JSON.parse(sessionStorage.getItem('IsAdmin'));
        this.showLogOut = true;
      }
    })


  }


  logout() {
          sessionStorage.clear();
          this.showLogOut = false;
          this.isAdmin = false;
          this.router.navigate(['login']);
  }

  goToAdmin() {
    this.router.navigate(['admin']);
  }
  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }
}
