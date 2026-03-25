import { Injectable, inject } from '@angular/core';
import { from, Observable, map } from 'rxjs';
import { SupabaseService } from './supabase.service';

export interface UploadResponse {
  message: string;
  url: string;
}

@Injectable({
  providedIn: 'root'
})
export class StorageService {
  private supabase = inject(SupabaseService);

  uploadFile(file: File, bucket: string = 'blog_assets'): Observable<UploadResponse> {
    const fileName = `${Date.now()}_${file.name.replace(/\s/g, '_')}`;
    const filePath = `${fileName}`;

    return from(this.supabase.client.storage.from(bucket).upload(filePath, file)).pipe(
      map(res => {
        if (res.error) throw res.error;
        
        const { data } = this.supabase.client.storage.from(bucket).getPublicUrl(filePath);
        return {
          message: 'Archivo subido con éxito',
          url: data.publicUrl
        };
      })
    );
  }
}
