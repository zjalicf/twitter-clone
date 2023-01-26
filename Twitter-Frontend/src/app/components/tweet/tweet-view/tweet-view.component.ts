import { Component, OnDestroy, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute } from '@angular/router';
import { TimespentDTO } from 'src/app/dto/TimespentDTO';
import { TweetID } from 'src/app/dto/tweetIdDTO';
import { Favorite } from 'src/app/models/favorite.model';
import { Tweet } from 'src/app/models/tweet.model';
import { User } from 'src/app/models/user.model';
import { TweetService } from 'src/app/services/tweet.service';
import { UserService } from 'src/app/services/user.service';
import { TweetLikesDialogComponent } from '../tweet-likes-dialog/tweet-likes-dialog.component';

@Component({
  selector: 'app-tweet-view',
  templateUrl: './tweet-view.component.html',
  styleUrls: ['./tweet-view.component.css']
})
export class TweetViewComponent implements OnInit, OnDestroy {

  constructor(
    private tweetService: TweetService,
    private userService: UserService,
    private route: ActivatedRoute,
    private dialog: MatDialog
  ) { }

  imagePath: string = ""
  likesByTweet: Favorite[] = [];
  tweet: Tweet = new Tweet();
  tweetID: TweetID = new TweetID();
  tweet_id = String(this.route.snapshot.paramMap.get("id"));
  loggedInUser: User = new User();
  totalLikes: number = 0;
  isLiked: boolean = false;
  isRetweeted: boolean = false;
  liked: string = "favorite_border";
  isThatMeLoggedIn: boolean = false;

  startTime: number = 0;
  endTime: number = 0;

  ngOnInit(): void {
    this.startTime = performance.now();
    this.totalLikes = this.tweet.favorite_count;

    this.tweetService.GetOneTweetById(this.tweet_id)
      .subscribe({
        next: (data: Tweet) => {
          this.tweet = data;
          
          if(this.tweet.image) {
            this.tweetService.GetImageByTweet(this.tweet_id).subscribe(response => {
                const fileReader = new FileReader();
                fileReader.readAsDataURL(response);
                fileReader.onload = () => {
                this.imagePath = fileReader.result as string;
            }  
          });
        }
        this.isThatMe()
        }
      });
      
      
  }

  ngOnDestroy(): void {
    this.endTime = performance.now();
    let seconds = Math.round((this.endTime - this.startTime) / 1000);
    
    console.log(seconds);
    
    if (this.tweet.advertisement) {
        let timespent = new TimespentDTO();
        timespent.tweet_id = this.tweet_id;
        timespent.timespent = seconds; 
        this.tweetService.TimespentOnAd(timespent).subscribe();
    }
  }

  isThatMe() {
    this.userService.GetMe()
      .subscribe({
        next: (data: User) => {
          this.loggedInUser = data;
          if (this.tweet.username === this.loggedInUser.username) {
            this.isThatMeLoggedIn = true;
          } else {
            this.isThatMeLoggedIn = false;
          }
          return this.isThatMeLoggedIn
        },
        error: (error) => {
          console.log(error);
          return this.isThatMeLoggedIn

        }
      });
  }

  
  likeTweet(tweet: Tweet) {

    this.tweetID.id = tweet.id
    this.tweetService.LikeTweet(this.tweet).subscribe(
      {
        next: (data) => {
          if (data == 201) {
            this.isLiked = true
            this.tweet.favorite_count++
            this.tweetService.GetLikesByTweet(this.tweet.id)
              .subscribe({
                next: (data) => {
                  this.likesByTweet = data;
                },
                error: (error) => {
                  console.log(error);
                }
              })
          } else {
            this.isLiked = false
            this.tweet.favorite_count--
            this.tweetService.GetLikesByTweet(this.tweet.id).subscribe({
              next: (data) => {
                this.likesByTweet = data;
              },
              error: (error) => {
                console.log(error);
              }
            })
          }
        }
      });
  }

  retweet(tweet: Tweet) {
    alert("retweeted")
  }

  openDialog(): void {
    const dialogRef = this.dialog.open(TweetLikesDialogComponent, {
      data: this.likesByTweet,
    });
    dialogRef.afterClosed().subscribe(result => {
      if (result == "username") {
        this.dialog.closeAll();
      }
    });
  }

  sendCount(): void {
      if (this.tweet.advertisement) {
        let tweetID = new TweetID();
        tweetID.id = this.tweet_id;
        this.tweetService.ViewedProfileFromAd(tweetID).subscribe();
      }
  }

  handleClick() {
    console.log(event)
  }

  isAnAd(): boolean {
    if (this.tweet.advertisement) {
      return true;
    } else {
      return false;
    }
  }

}
