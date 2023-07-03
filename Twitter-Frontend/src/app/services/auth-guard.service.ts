import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, Router, RouterStateSnapshot } from '@angular/router';
import { User } from '../models/user.model';
import { UserService } from './user.service';

@Injectable({
  providedIn: 'root'
})
export class AuthGuard {

  constructor(private router: Router,
              private userService: UserService) { }

  user: User = new User();
    
  ngOnInit(): void {
    this.userService.GetMe()
      .subscribe({
        next: (data: User) => {
          this.user = data;
        },
        error: (error) => {
          console.log(error);
        }
      })
    
  }

  canActivate(next: ActivatedRouteSnapshot, state: RouterStateSnapshot) {  
      if (localStorage.getItem('authToken')) {  
          return true;  
      } else { 
      this.router.navigate(['/Login']);  
      return false
    }  
  } 
  
}
