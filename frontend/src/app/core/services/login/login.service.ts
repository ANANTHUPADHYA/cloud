import { environment } from 'src/environments/environment';
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { LoginResponse, UserParams, RegisterResponse } from './models';
import { Observable, of, Subject } from 'rxjs';
import { urlConstants } from '../../rest-api-configuration';

@Injectable({
  providedIn: 'root'
})
export class LoginService {
private baseUrl: string;
public userLoggedIn = new Subject();

  constructor(private http: HttpClient) {
    this.baseUrl = environment.url;
   }

   loginUser({email, password}): Observable<LoginResponse> {
    return this.http.put<LoginResponse>(`${this.baseUrl}${urlConstants.LOGIN}`, {emailAddress: email, password});
  }


  registerUser(userDetails: UserParams):  Observable<RegisterResponse> {
    return this.http.post<RegisterResponse>(`${this.baseUrl}${urlConstants.REGISTER}`, userDetails); 
  }


}
