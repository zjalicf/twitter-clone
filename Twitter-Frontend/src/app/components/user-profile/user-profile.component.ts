import { Component, OnInit } from '@angular/core';
import { User } from 'src/app/models/user.model';
import { UserServiceService } from 'src/app/services/user.service.service';

@Component({
  selector: 'app-user-profile',
  templateUrl: './user-profile.component.html',
  styleUrls: ['./user-profile.component.css']
})
export class UserProfileComponent implements OnInit {

  user?: User
  
  constructor(private UserService: UserServiceService) { }

  ngOnInit(): void {

    this.UserService.GetOneUserByUsername("nani13051411").subscribe(
      data => {
        this.user = data
      }
    )

  }

}
