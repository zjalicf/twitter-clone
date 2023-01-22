import { Component, Input, OnInit } from '@angular/core';
import { TweetID } from 'src/app/dto/tweetIdDTO';
import { Tweet } from 'src/app/models/tweet.model';
import { User } from 'src/app/models/user.model';
import { TweetService } from 'src/app/services/tweet.service';
import { UserService } from 'src/app/services/user.service';
import { MatDialog } from '@angular/material/dialog';
import { MatDialogModule } from '@angular/material/dialog';
import { TweetLikesDialogComponent } from '../tweet-likes-dialog/tweet-likes-dialog.component';
import { Favorite } from 'src/app/models/favorite.model';
import { Router } from '@angular/router';
import { HttpHeaderResponse } from '@angular/common/http';

@Component({
  selector: 'app-tweet-item',
  templateUrl: './tweet-item.component.html',
  styleUrls: ['./tweet-item.component.css']
})
export class TweetItemComponent implements OnInit {

  constructor(private userService: UserService,
    private tweetService: TweetService,
    public dialog: MatDialog) { }

  @Input() tweet: Tweet = new Tweet();

  imagePath: string = ""

  likesByTweet: Favorite[] = [];

  loggedInUser: User = new User();
  tweetID: TweetID = new TweetID();
  totalLikes: number = 0
  isLiked: boolean = false;
  isRetweeted: boolean = false;
  liked: string = "favorite_border";

  ngOnInit(): void {
    this.totalLikes = this.tweet.favorite_count;

    if(this.tweet.image) {
      this.tweetService.GetImageByTweet(this.tweet.id).subscribe(response => {
        const fileReader = new FileReader();
        fileReader.readAsDataURL(response);
        fileReader.onload = () => {
          this.imagePath = fileReader.result as string;
        }
      });
  
    }

    this.userService.GetMe()
      .subscribe({
        next: (data: User) => {
          this.loggedInUser = data;
        },
        error: (error) => {
          console.log(error);
        }
      });

    this.tweetService.GetLikesByTweet(this.tweet.id)
      .subscribe({
        next: (data) => {
          this.likesByTweet = data;
          if (data != null) {
            this.likesByTweet.forEach(like => {
              if (like.username == this.loggedInUser.username) {
                this.isLiked = true;
              } else {
                this.isLiked = false;
              }
            });
          }
        },
        error: (error) => {
          console.log(error);
        }
      })
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
    this.tweetService.LikeTweet(this.tweetID).subscribe(
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
