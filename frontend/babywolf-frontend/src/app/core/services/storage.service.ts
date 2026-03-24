import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';

export interface UploadResponse {
  message: string;
  url: string;
}

@Injectable({
  providedIn: 'root'
})
export class StorageService {
  private http = inject(HttpClient);
  private apiUrl = `${environment.apiUrl}/admin/upload`;

  uploadFile(file: File, folder: string = 'uploads'): Observable<UploadResponse> {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('folder', folder);

    return this.http.post<UploadResponse>(this.apiUrl, formData);
  }
}
