import { Injectable } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { urlConstants } from '../../rest-api-configuration';
import { UserResponse, DeleteResponse, AllUsers } from './models';

@Injectable({
  providedIn: 'root'
})
export class UserService {
private baseUrl: string;
public userDetails = new Subject();
  constructor(private http: HttpClient) {
    this.baseUrl = environment.url;
   }
  getUserDetails(userID): Observable<UserResponse> {
    return this.http.get<UserResponse>(`${this.baseUrl}${urlConstants.GET_USER}/${userID}`);
  }

  deleteFile(userID, filename):  Observable<DeleteResponse> {
    return this.http.delete<DeleteResponse>(`${this.baseUrl}${urlConstants.DELETE}/${userID}/file?file=${filename}`);
  }

  getAllUsers():Observable<AllUsers>{
    return this.http.get<AllUsers>(`${this.baseUrl}${urlConstants.GET_USER}`);
  }

}
