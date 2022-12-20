import { Component, Input, OnInit, ChangeDetectorRef } from '@angular/core';
import { TweetID } from 'src/app/dto/tweetIdDTO';
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

  constructor(private userService: UserService,
              private tweetService: TweetService,
              private crf: ChangeDetectorRef) { }

   @Input() tweet: Tweet = new Tweet();

   loggedInUser: User = new User();
   tweetID: TweetID = new TweetID();
   totalLikes: number = 0

  ngOnInit(): void {

    this.totalLikes = this.tweet.favorite_count

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

  likeTweet(tweet: Tweet) {
    this.tweetID.id = tweet.id
    console.log(this.tweetID)
    console.log(tweet)
    this.tweetService.LikeTweet(this.tweetID).subscribe(
      {next : (data) => {
        if (data == 201) {
          this.totalLikes = this.tweet.favorite_count + 1
          alert("Tweet Liked")

        }else{
          this.totalLikes = this.tweet.favorite_count - 1
          alert("Tweet Unliked")
        }
          
      }, complete: () => {
        
      }})
      
  }
}
