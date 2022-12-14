import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { FollowRequest } from 'src/app/models/followRequest.model';
import { User } from 'src/app/models/user.model';
import { UserService } from 'src/app/services/user.service';

@Component({
  selector: 'app-follow-requests',
  templateUrl: './follow-requests.component.html',
  styleUrls: ['./follow-requests.component.css']
})
export class FollowRequestsComponent implements OnInit {

  constructor(private userService: UserService,
              private router: Router) { }

  user: User = new User();
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
      this.firstFollowRequest.sender = "Milan";
      this.firstFollowRequest.receiver = "Petar";
      this.firstFollowRequest.status = "Pending";

      this.secondFollowRequest.sender = "Filip";
      this.secondFollowRequest.receiver = "Petar";
      this.secondFollowRequest.status = "Pending";

      this.followRequests.push(this.firstFollowRequest);
      this.followRequests.push(this.secondFollowRequest);
  }

}
