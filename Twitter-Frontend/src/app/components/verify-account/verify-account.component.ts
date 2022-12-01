import { HttpErrorResponse } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { VerificationRequest } from 'src/app/dto/verificationRequest';
import { AuthService } from 'src/app/services/auth.service';
import { VerificationService } from 'src/app/services/verify.service';

@Component({
  selector: 'app-verify-account',
  templateUrl: './verify-account.component.html',
  styleUrls: ['./verify-account.component.css']
})
export class VerifyAccountComponent implements OnInit {

  formGroup: FormGroup = new FormGroup({
    verificationToken: new FormControl(''),
  });
  submitted = false;

  constructor(private authService: AuthService,
              private formBuilder: FormBuilder,
              private router: Router,
              private verificationService: VerificationService) { }


  ngOnInit(): void {
    this.formGroup = this.formBuilder.group({
      verificationToken: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(36)]],
    })
  }

  get f(): { [key: string]: AbstractControl } {
    return this.formGroup.controls;
  }

  onSubmit() {
    this.submitted = true;

    if (this.formGroup.invalid) {
      return;
    }
    
    let userToken = "";
    this.verificationService.currentVerificationToken.subscribe(token => userToken = token);
    let mailToken: string = this.formGroup.get("verificationToken")?.value;
    let request = new VerificationRequest();
    request.user_token = userToken;
    request.mail_token = mailToken;
    console.log(request)
    this.authService.VerifyAccount(request)
      .subscribe({
          next: (response: void) => {
            alert("You have been successfully registered to Twitter");
            this.router.navigate(['/Login'])
          },
          error: (error: HttpErrorResponse) => {
            if (error.status == 406 || error.status == 400) {
              this.formGroup.setErrors({invalidToken:true})                
            }
            else if(error.status == 404){
              this.formGroup.setErrors({expiredToken:true})
            }
          }
      })
  }

}
