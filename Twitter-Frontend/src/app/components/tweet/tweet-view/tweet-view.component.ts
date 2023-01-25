import { Component, OnDestroy, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Router } from '@angular/router';
import { TweetID } from 'src/app/dto/tweetIdDTO';
import { Favorite } from 'src/app/models/favorite.model';
import { Tweet } from 'src/app/models/tweet.model';
import { User } from 'src/app/models/user.model';
import { TweetService } from 'src/app/services/tweet.service';
import { TweetLikesDialogComponent } from '../tweet-likes-dialog/tweet-likes-dialog.component';

@Component({
  selector: 'app-tweet-view',
  templateUrl: './tweet-view.component.html',
  styleUrls: ['./tweet-view.component.css']
})
export class TweetViewComponent implements OnInit, OnDestroy {

  constructor(
    private tweetService: TweetService,
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

  ngOnInit(): void {
    this.totalLikes = this.tweet.favorite_count;

    this.tweetService.GetOneTweetById(this.tweet_id)
      .subscribe({
        next: (data: Tweet) => {
          this.tweet = data;
        }
      });

      if(this.tweet.image) {
        this.tweetService.GetImageByTweet(this.tweet.id).subscribe(response => {
          const fileReader = new FileReader();
          fileReader.readAsDataURL(response);
          fileReader.onload = () => {
            this.imagePath = fileReader.result as string;
          }
        });
      }
  }

  date = new Date();
  ngOnDestroy(): void {
    
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

  handleClick() {
    console.log(event)
  }

}
