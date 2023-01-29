import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { Tweet } from 'src/app/models/tweet.model';
import { User } from 'src/app/models/user.model';
import { FollowService } from 'src/app/services/follow.service';
import { TweetService } from 'src/app/services/tweet.service';
import { UserService } from 'src/app/services/user.service';
import { FollowComponentDialogComponent } from '../follow-component-dialog/follow-component-dialog.component';
import { FollowingComponentDialogComponent } from '../following-component-dialog/following-component-dialog.component';

@Component({
  selector: 'app-my-profile',
  templateUrl: './my-profile.component.html',
  styleUrls: ['./my-profile.component.css']
})
export class MyProfileComponent implements OnInit {

  constructor(private userService: UserService,
              private router: Router,
              private tweetService: TweetService,
              private followService: FollowService,
              public dialog: MatDialog,
              ) { }

  user: User = new User();
  tweets: Tweet[] = [];
  isBusinessBool: boolean = false
  followings: User[] = []
  followers: User[] = []
    
  ngOnInit(): void {


    this.userService.GetMe()
      .subscribe({
        next: (data: User) => {
          this.user = data;
          if (this.user.userType == "Business"){
            this.isBusinessBool = true
          }
        },
        error: (error) => {
          console.log(error);
        },
        complete: () => {
          this.tweetService.GetTweetsForUser(this.user.username)
            .subscribe({
              next: (data: Tweet[]) => {
                this.tweets = data;
                this.followService.GetFollowingsForMe().subscribe(
                  data => {
                    this.followings = data
                    console.log(this.followings)
                  }
                )

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

  openDialogFollowers(): void {
    const dialogRef = this.dialog.open(FollowComponentDialogComponent, {
      data: this.followers,
    });
    dialogRef.afterClosed().subscribe(result => {
      if (result == "username") {
        this.dialog.closeAll();
      }
    });
  }

  openDialogFollowings(): void {
    const dialogRef = this.dialog.open(FollowingComponentDialogComponent, {
      data: this.followings,
    });
    dialogRef.afterClosed().subscribe(result => {
      if (result == "username") {
        this.dialog.closeAll();
      }
    });
  }

}
