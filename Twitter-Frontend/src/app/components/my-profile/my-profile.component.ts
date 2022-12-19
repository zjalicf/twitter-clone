import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Tweet } from 'src/app/models/tweet.model';
import { User } from 'src/app/models/user.model';
import { TweetService } from 'src/app/services/tweet.service';
import { UserService } from 'src/app/services/user.service';

@Component({
  selector: 'app-my-profile',
  templateUrl: './my-profile.component.html',
  styleUrls: ['./my-profile.component.css']
})
export class MyProfileComponent implements OnInit {

  constructor(private userService: UserService,
              private router: Router,
              private tweetService: TweetService) { }

  user: User = new User();
  tweets: Tweet[] = [];
    
  ngOnInit(): void {
    this.userService.GetMe()
      .subscribe({
        next: (data: User) => {
          this.user = data;
        },
        error: (error) => {
          console.log(error);
        },
        complete: () => {
          this.tweetService.GetTweetsForUser(this.user.username)
            .subscribe({
              next: (data: Tweet[]) => {
                this.tweets = data;
              },
              error: (error) => {
                console.log(error);
              }
            });
        }
      });
  }

  updatePassword() {
    this.router.navigateByUrl("Change-Password")
  }

  UpdateVisibility() {
    this.userService.ChangeVisibility().subscribe()
  }

}
