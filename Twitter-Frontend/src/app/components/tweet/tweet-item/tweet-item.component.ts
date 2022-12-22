import { Component, Input, OnInit } from '@angular/core';
import { TweetID } from 'src/app/dto/tweetIdDTO';
import { Tweet } from 'src/app/models/tweet.model';
import { User } from 'src/app/models/user.model';
import { TweetService } from 'src/app/services/tweet.service';
import { UserService } from 'src/app/services/user.service';
import {MatDialog} from '@angular/material/dialog';
import { MatDialogModule } from '@angular/material/dialog';
import { TweetLikesDialogComponent } from '../tweet-likes-dialog/tweet-likes-dialog.component';
import { Favorite } from 'src/app/models/favorite.model';
import { Router } from '@angular/router';

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

   likesByTweet: Favorite[] = [];

   loggedInUser: User = new User();
   tweetID: TweetID = new TweetID();
   totalLikes: number = 0

  ngOnInit(): void {
    this.totalLikes = this.tweet.favorite_count;

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
          },
          error: (error) => {
            console.log(error);
          }
        })
  }

  // isLiked(): boolean {
  //   for (let favorite of this.likesByTweet) {
  //     if (favorite.username == this.loggedInUser.username) {
  //       return true;
  //     }
  //   }
  //   return false;
  // }

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
      {next : (data) => {
        console.log(data)
          if (data == 201) {
            this.tweet.favorite_count++

          }else{
            this.tweet.favorite_count--
          } 
      }});
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
}
