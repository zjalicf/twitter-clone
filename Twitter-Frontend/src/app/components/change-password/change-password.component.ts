import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { ChangePasswordDTO } from 'src/app/dto/changePasswordDTO';
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
  });

  constructor(private formBuilder: FormBuilder) { }

  submitted = false;

  ngOnInit(): void {
    this.formGroup = this.formBuilder.group({
      currentPassword: ['', [Validators.required, Validators.minLength(8), Validators.maxLength(30), PasswordStrenghtValidator(), Validators.pattern('[-_a-zA-Z0-9]*')]],
      newPassword: ['', [Validators.required, Validators.minLength(8), Validators.maxLength(30), PasswordStrenghtValidator(), Validators.pattern('[-_a-zA-Z0-9]*')]]
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

    changePassword.currentPassword = this.formGroup.get('currentPassword')?.value;
    changePassword.newPassword = this.formGroup.get('newPassword')?.value;

    // this.authService.ChangePassword(changePassword)
    //   .subscribe({

    //   })
  }

}
