import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { environment } from "src/environments/environment";
import { User } from "../models/user.model";

@Injectable({
providedIn: 'root'
})
export class AuthService {
    constructor(private http: HttpClient) { }

    public registerRegular(user: User): Observable<User> {
        return this.http.post<User>(`${environment.baseApiUrl}/registerRegular`, user);
    }

    public registerBusiness(user: User): Observable<User> {
        return this.http.post<User>(`${environment.baseApiUrl}/registerBusiness`, user);
    }

    // public login(loginDTO: LoginDTO): Observable<string> {
    //     return this.http.post(`${environment.baseApiUrl}/${this.url}/Login`, loginDTO, { responseType: 'text' });
    // }
  
}