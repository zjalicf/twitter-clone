import { Component, OnInit } from '@angular/core';
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
  tweets: Tweet[] = []
  profileUsername = String(this.route.snapshot.paramMap.get("username"));
  
  constructor(private UserService: UserService,
              private route: ActivatedRoute,
              private router: Router,
              private TweetService: TweetService,
              private followService: FollowService) { }

  ngOnInit(): void {
    this.UserService.GetOneUserByUsername(this.profileUsername)
      .subscribe({
        next: (data: User) => {
          this.user = data;
        },
        error: (error) => {
          console.log(error);
        }
      })
    this.TweetService.GetTweetsForUser(this.profileUsername)
      .subscribe({
        next: (data: Tweet[]) => {
          this.tweets = data;
        },
        error: (error) => {
          console.log(error);
        }
      })
  
  }

  SendRequest(user: User){
    var followReq = new FollowRequest()
    followReq.receiver = user.username
    if (user.visibility){
      this.followService.SendRequest("private", followReq).subscribe()
    }else {
      this.followService.SendRequest("public", followReq).subscribe()
    }
  }

}
