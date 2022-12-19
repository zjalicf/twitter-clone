import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { FollowRequest } from 'src/app/models/followRequest.model';
import { FollowService } from 'src/app/services/follow.service';

@Component({
  selector: 'app-follow-request-item',
  templateUrl: './follow-request-item.component.html',
  styleUrls: ['./follow-request-item.component.css']
})
export class FollowRequestItemComponent implements OnInit {

  @Input() followRequest: FollowRequest = new FollowRequest();
  @Output() answerFollowRequest = new EventEmitter<any>()

  constructor(private followService: FollowService) { }

  ngOnInit(): void {

  }

  AcceptRequest(id: string){
    this.followService.AcceptRequest(id)
      .subscribe({
        next: (data) => {
          this.answerFollowRequest.emit();
          alert("Request Accepted");
        },
        error: (error) => {
          console.log(error);
        }
      }
    )
  }
  
  DeclineRequest(id: string){
    this.followService.DeclineRequest(id)
      .subscribe({
        next: (data) => {
          alert("Request Denied");
        },
        error: (error) => {
          console.log(error);
        }
      }
    )
  }

}
