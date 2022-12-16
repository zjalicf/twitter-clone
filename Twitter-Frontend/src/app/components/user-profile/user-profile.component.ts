import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Tweet } from 'src/app/models/tweet.model';
import { User } from 'src/app/models/user.model';
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
              private TweetService: TweetService) { }

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

  isPrivate(): boolean {
    if (this.user.visibility == true) {
      return true;
    } else {
      return false;
    }
  }

}
