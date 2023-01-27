import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ActivatedRoute, Router } from '@angular/router';
import { FollowRequest } from 'src/app/models/followRequest.model';
import { Tweet } from 'src/app/models/tweet.model';
import { User } from 'src/app/models/user.model';
import { FollowService } from 'src/app/services/follow.service';
import { TweetService } from 'src/app/services/tweet.service';
import { UserService } from 'src/app/services/user.service';

@Component({
  selector: 'app-user-profile',
  templateUrl: './user-profile.component.html',
  styleUrls: ['./user-profile.component.css']
})
export class UserProfileComponent implements OnInit {

  user: User = new User();
  loggedInUser = new User();
  tweets: Tweet[] = []
  profileUsername = String(this.route.snapshot.paramMap.get("username"));
  isFollowing: boolean = false;

  constructor(
    private UserService: UserService,
    private route: ActivatedRoute,
    private router: Router,
    private TweetService: TweetService,
    private followService: FollowService,
    private _snackBar: MatSnackBar
  ) { }

  ngOnInit(): void {
      this.followService.IsFollowExist(this.profileUsername).subscribe(response => {
        this.isFollowing = response;
      })
    

    this.UserService.GetOneUserByUsername(this.profileUsername)
      .subscribe({
        next: (data: User) => {
          this.user = data;
        },
        error: (error) => {
          console.log(error);
          this.router.navigate(["/404"])
        }
      });

    this.TweetService.GetTweetsForUser(this.profileUsername)
      .subscribe({
        next: (data: Tweet[]) => {
          this.tweets = data;
        },
        error: (error) => {
          console.log(error);
        }
      });

    this.UserService.GetMe()
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
    if (this.user.username == this.loggedInUser.username) {
      return true;
    } else {
      return false;
    }
  }

  // isPrivate(): boolean {
  //   console.log("Privacy is " + this.user.privacy)
  //   return this.user.privacy
  // }

  SendRequest(user: User){
    var followReq = new FollowRequest()
    followReq.receiver = user.username
    if (user.privacy){
      this.followService.SendRequest("private", followReq).subscribe(
        data => {
          console.log(data.status)
          this.openSnackBar("Request sended!", "")
        },
        error => {
          if (error.status == 400) {
            alert("You already follow this user")

          }
        }
        )
    }else {
      this.followService.SendRequest("public", followReq).subscribe(
        data => {
          console.log(data.status)
        },
        error => {
          if (error.status == 400) {
            alert("You already follow this user")

          }
        }
      )
    }
  }

  openSnackBar(message: string, action: string) {
    this._snackBar.open(message, action,  {
      duration: 3500
    });
  }

}
