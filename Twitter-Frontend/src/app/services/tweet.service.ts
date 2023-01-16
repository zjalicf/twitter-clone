import { HttpClient, HttpHeaders } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { environment } from "src/environments/environment";
import { AddTweetDTO } from "../dto/addTweetDTO";
import { TweetID } from "../dto/tweetIdDTO";
import { Favorite } from "../models/favorite.model";
import { Tweet } from "../models/tweet.model";

@Injectable({
    providedIn: 'root'
    })

    export class TweetService {
        private url = "tweets";
        constructor(private http: HttpClient) { }
    

        public AddTweet(tweet: AddTweetDTO): Observable<Tweet> {
            return this.http.post<Tweet>(`${environment.baseApiUrl}/${this.url}/`, tweet);
        }

        public GetHomeFeed(): Observable<any> {
            return this.http.get<any>(`${environment.baseApiUrl}/${this.url}/feed`);
        }
    
        public GetTweetsForUser(username: string): Observable<any> {
            return this.http.get<any>(`${environment.baseApiUrl}/${this.url}/user/` + username)
        }

        public LikeTweet(tweetID: TweetID): Observable<any> {
            return this.http.post<any>(`${environment.baseApiUrl}/${this.url}/favorite`, tweetID)
        }

        public GetLikesByTweet(tweetID: string): Observable<Favorite[]> {
            return this.http.get<Favorite[]>(`${environment.baseApiUrl}/${this.url}/whoLiked/` + tweetID)
        }

        public Retweet(tweetID: TweetID): Observable<void> {
            return this.http.post<void>(`${environment.baseApiUrl}/${this.url}/retweet/`,tweetID)
        }

    }