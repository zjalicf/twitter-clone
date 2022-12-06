import { HttpBackend, HttpErrorResponse } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { ChangePasswordDTO } from 'src/app/dto/changePasswordDTO';
import { AuthService } from 'src/app/services/auth.service';
import { PasswordStrenghtValidator } from 'src/app/services/customValidators';

@Component({
  selector: 'app-change-password',
  templateUrl: './change-password.component.html',
  styleUrls: ['./change-password.component.css']
})
export class ChangePasswordComponent implements OnInit {
  formGroup: FormGroup = new FormGroup({
    currentPassword: new FormGroup(''),
    newPassword: new FormGroup(''),
    newPasswordConfirm: new FormGroup('')
  });


  constructor(private router: Router,
              private formBuilder: FormBuilder,
              private authService: AuthService
  ) { }

  submitted = false;

  ngOnInit(): void {
    this.formGroup = this.formBuilder.group({
      currentPassword: ['', [Validators.required, Validators.minLength(8), Validators.maxLength(30), PasswordStrenghtValidator(), Validators.pattern('[-_a-zA-Z0-9]*')]],
      newPassword: ['', [Validators.required, Validators.minLength(8), Validators.maxLength(30), PasswordStrenghtValidator(), Validators.pattern('[-_a-zA-Z0-9]*')]],
      newPasswordConfirm: ['', [Validators.required, Validators.minLength(8), Validators.maxLength(30), PasswordStrenghtValidator(), Validators.pattern('[-_a-zA-Z0-9]*')]]
    });
  }

  get changePasswordGroup(): { [key: string]: AbstractControl } {
    return this.formGroup.controls;
  }

  onSubmit() {
    this.submitted = true;

    if (this.formGroup.invalid) {
      return;
    }

    let changePassword: ChangePasswordDTO = new ChangePasswordDTO();

    changePassword.old_password = this.formGroup.get('currentPassword')?.value
    changePassword.new_password = this.formGroup.get('newPassword')?.value;
    changePassword.new_password_confirm = this.formGroup.get('newPasswordConfirm')?.value;

    console.log(changePassword.new_password + changePassword.new_password_confirm + changePassword.old_password)

    this.authService.ChangePassword(changePassword)
      .subscribe({
        next: (data: string) => {
          localStorage.clear
          this.router.navigate(["/Login"])
        },
        error: (err: HttpErrorResponse) => {
          if(err.status == 409){
            alert("Old password not match!")
          } else if (err.status == 406){
              alert("New passwrod not match!")
          } else if (err.status == 200){
            alert("Password changed successfully!")
          }
          // localStorage.setItem("authToken", "")
          // this.router.navigate(["/Login"])
        }
      }
      )
  }

}
