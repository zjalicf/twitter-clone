import { Component, Input, OnInit } from '@angular/core';
import { Tweet } from 'src/app/models/tweet.model';
import { User } from 'src/app/models/user.model';
import { TweetService } from 'src/app/services/tweet.service';
import { UserService } from 'src/app/services/user.service';

@Component({
  selector: 'app-tweet-item',
  templateUrl: './tweet-item.component.html',
  styleUrls: ['./tweet-item.component.css']
})
export class TweetItemComponent implements OnInit {

  constructor(private userService: UserService) { }

   @Input() tweet: Tweet = new Tweet();

   loggedInUser: User = new User();

  ngOnInit(): void {
    this.userService.GetMe()
      .subscribe({
        next: (data: User) => {
          this.loggedInUser = data;
        },
        error: (error) => {
          console.log(error);
        }
      });
  }

  isThatMe(): boolean {
    if (this.tweet.username == this.loggedInUser.username) {
      return true;
    } else {
      return false;
    }
  }

  likeTweet() {
    alert("Tweet Liked")
  }

}
