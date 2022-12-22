import { Component, Input, OnInit } from '@angular/core';
import { TweetID } from 'src/app/dto/tweetIdDTO';
import { Tweet } from 'src/app/models/tweet.model';
import { User } from 'src/app/models/user.model';
import { TweetService } from 'src/app/services/tweet.service';
import { UserService } from 'src/app/services/user.service';
import {MatDialog} from '@angular/material/dialog';
import { MatDialogModule } from '@angular/material/dialog';
import { TweetLikesDialogComponent } from '../tweet-likes-dialog/tweet-likes-dialog.component';

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

   @Input() testUsers: User[] = [];

   loggedInUser: User = new User();
   tweetID: TweetID = new TweetID();
   usernames: string[] = ["Milan", "Petar"]
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

      this.testUsers.push(this.loggedInUser);
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
            this.totalLikes =+ 1
            alert("Tweet Liked")

          }else{
            this.totalLikes =- 1
            alert("Tweet Unliked")
          } 
      }});
  }

  openDialog(): void {
    const dialogRef = this.dialog.open(TweetLikesDialogComponent, {
      data: this.loggedInUser,
    });
    dialogRef.afterClosed().subscribe(result => {
      console.log('The dialog was closed');
    });
  }
}
