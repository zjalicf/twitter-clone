import { Component, OnInit } from '@angular/core';
import { Tweet } from 'src/app/models/tweet.model';
import { User } from 'src/app/models/user.model';
import { TweetService } from 'src/app/services/tweet.service';
import { UserService } from 'src/app/services/user.service';

@Component({
  selector: 'app-main-page',
  templateUrl: './main-page.component.html',
  styleUrls: ['./main-page.component.css']
})
export class MainPageComponent implements OnInit {

  tweets: Tweet[] = []
  user: User = new User()

  constructor(private tweetService: TweetService,
    private userService: UserService) { }


  //treba napraviti da se prikazu samo tvitovi usera koje pratimo

  ngOnInit(): void {
    
    this.userService.GetMe().subscribe(response => {
      this.user = response
      this.tweetService.GetTweetsForUser(this.user.username).subscribe( data => {
        this.tweets = data
      })
    })
    
  }

}
