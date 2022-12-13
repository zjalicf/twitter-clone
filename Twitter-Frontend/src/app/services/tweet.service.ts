import { HttpClient, HttpHeaders } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { environment } from "src/environments/environment";
import { Tweet } from "../models/tweet.model";

@Injectable({
    providedIn: 'root'
    })

    export class TweetService {
        private url = "tweets";
        constructor(private http: HttpClient) { }
    

        public AddTweet(tweet: Tweet): Observable<Tweet> {
            let headers = new HttpHeaders({
                "Content-Type":"application/json",
                "Authorization": "" + localStorage.getItem("authToken")
            })
            let options = {headers: headers}
            console.log(localStorage.getItem("authToken"))

            return this.http.post<Tweet>(`${environment.baseApiUrl}/${this.url}/`, tweet, options);
        }

        public GetAllTweets(): Observable<any> {
            let headers = new HttpHeaders({
                "Content-Type":"application/json",
                "Authorization": "" + localStorage.getItem("authToken")
            })
            let options = {headers: headers}
            console.log(localStorage.getItem("authToken"))
            return this.http.get<any>(`${environment.baseApiUrl}/${this.url}/`, options);
        }
    
        public GetTweetsForUser(username: string): Observable<any> {
            return this.http.get<any>(`${environment.baseApiUrl}/${this.url}/user/` + username)
        }

    }