import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';
import { User } from '../models/user.model';

@Injectable({
  providedIn: 'root'
})
export class UserServiceService {

  private url = "users"
  constructor(private http: HttpClient) { }

  public GetOneUserByUsername(username: string): Observable<User>{
    let headers = new HttpHeaders({
      "Content-Type" : "application/json",
      "Authorization" : "Bearer " + localStorage.getItem("authToken"),
    });

    let options = {headers:headers};
    return this.http.get<User>(`${environment.baseApiUrl}/${this.url}/getOne/${username}`,options)
  }

}
