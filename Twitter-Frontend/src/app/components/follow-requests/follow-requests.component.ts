import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { FollowRequest } from 'src/app/models/followRequest.model';
import { User } from 'src/app/models/user.model';
import { FollowService } from 'src/app/services/follow.service';
import { UserService } from 'src/app/services/user.service';

@Component({
  selector: 'app-follow-requests',
  templateUrl: './follow-requests.component.html',
  styleUrls: ['./follow-requests.component.css']
})
export class FollowRequestsComponent implements OnInit {

  constructor(private userService: UserService,
              private router: Router,
              private followService: FollowService) { }

  user: User = new User();
  requests: FollowRequest[] = []
  firstFollowRequest: FollowRequest = new FollowRequest();
  secondFollowRequest: FollowRequest = new FollowRequest();

  followRequests: FollowRequest[] = []

  ngOnInit(): void {
    this.userService.GetMe()
      .subscribe({
        next: (data: User) => {
          this.user = data;
        },
        error: (error) => {
          console.log(error);
        } 
      });

      this.followService.GetRequestsForUser().subscribe(
        data => {
          this.requests = data
        }
      )
  }


AcceptRequest(id: string){
  this.followService.AcceptRequest(id).subscribe(
    data => {
      alert("Follow request accepted!")
    }
  )
}

DeclineRequest(id: string){
  this.followService.DeclineRequest(id).subscribe(
    data => {
      alert("Follow request rejected!")
    }
  )
}



}
