import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import {  Router } from '@angular/router';
import { Tweet } from 'src/app/models/tweet.model';
import { User } from 'src/app/models/user.model';
import { FollowService } from 'src/app/services/follow.service';
import { TweetService } from 'src/app/services/tweet.service';
import { UserService } from 'src/app/services/user.service';

@Component({
  selector: 'app-main-page',
  templateUrl: './main-page.component.html',
  styleUrls: ['./main-page.component.css']
})
export class MainPageComponent implements OnInit {

  tweets: Tweet[] = [];
  user: User = new User();
  recommendations: string[] = [];

  constructor(private tweetService: TweetService,
    private userService: UserService,
    private followService: FollowService,
    private _snackBar: MatSnackBar,
    private router: Router
    ) { }


  //treba napraviti da se prikazu samo tvitovi usera koje pratimo i tvitovi ulogovanog usera

  ngOnInit(): void {
    this.userService.GetMe()
      .subscribe({
        next: (data) => {
          this.user = data;
        },
        error: (error) => {
          console.log(error);
        },
        complete: () => {

          this.followService.Recommendations()
            .subscribe({
              next: (data) => {
                this.recommendations = data;
                console.log(this.recommendations);
              },
              error: (error) => {
                this.openSnackBar("The service is currently unavailable. Try again later.", "")
                console.log(error);
              }
            });
            
          this.tweetService.GetHomeFeed()
            .subscribe({
              next: (data) => {
                this.tweets = data;
              },
              error: (error) => {
                this.openSnackBar("The service is currently unavailable. Try again later.", "")
                console.log(error);
              }
            });
          
          
        }
      });
  }

  openSnackBar(message: string, action: string) {
    this._snackBar.open(message, action,  {
      duration: 3500
    });
  }

  OpenProfile(name: string){
    this.router.navigate(["View-Profile/" + name])

  }
}
