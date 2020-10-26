import { Injectable } from '@angular/core';
import { FileParams, FileResponse } from './models';
import { Observable } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { urlConstants } from '../../rest-api-configuration';

@Injectable({
  providedIn: 'root'
})
export class FilesService {
  private baseUrl: string;
    constructor(private http: HttpClient) {
      this.baseUrl = environment.url;
     }

  uploadFile(fileparams: FileParams): Observable<FileResponse> {
    const formData: FormData = new FormData();
    formData.append('file', fileparams.file);
    //formData.append('description', fileparams.description);
      return this.http.put<FileResponse>(`${this.baseUrl}${urlConstants.UPLOAD}/${fileparams.UserID}/upload`, formData);
    }
  
  
  editFile(fileparams: FileParams): Observable<FileResponse> {
    const formData: FormData = new FormData();
    //formData.append('file', fileparams.file, fileparams.filename);
    formData.append('description', fileparams.description);
    return this.http.patch<FileResponse>(`${this.baseUrl}${urlConstants.EDIT}/${fileparams.UserID}/file-update?file=${fileparams.filename}`, formData);
  }

  downloadFile(userId, filename): Observable<any> {    
      return this.http.get<any>(`${this.baseUrl}${urlConstants.DOWNLOAD}/${userId}/download?file=${filename}`);
  }
}
