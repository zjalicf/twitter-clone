import { HttpErrorResponse } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { MatLegacySnackBar as MatSnackBar } from '@angular/material/legacy-snack-bar';
import { Router } from '@angular/router';
import { AuthService } from 'src/app/services/auth.service';
import { RecoveryPasswordService } from 'src/app/services/recoveryPassword.service';

@Component({
  selector: 'app-recovery-enter-mail',
  templateUrl: './recovery-enter-mail.component.html',
  styleUrls: ['./recovery-enter-mail.component.css']
})
export class RecoveryEnterMailComponent implements OnInit {

  formGroup: FormGroup = new FormGroup({
    email: new FormControl(''),
  });
  submitted = false;

  constructor(
    private authService: AuthService,
    private formBuilder: FormBuilder,
    private router: Router,
    private recoveryService: RecoveryPasswordService,
    private _snackBar: MatSnackBar
  ) { }

  ngOnInit(): void {
    this.formGroup = this.formBuilder.group({
      email: ['', [Validators.required, Validators.email]],
    })
  }

  get form(): { [key: string]: AbstractControl } {
    return this.formGroup.controls;
  }

  onSubmit(){
    this.submitted = true;

    if (this.formGroup.invalid) {
      return;
    }

    let email = this.formGroup.get("email")?.value

    this.authService.RequestRecoverPassword(email).subscribe({
      next: (token: string) => {
        this.recoveryService.updateToken(token)
        this.openSnackBar("Recovery token has sended to your e-mail. Please type token from e-mail", "")
        this.router.navigate(['/Recovery-Token'])
      },
      error: (error: HttpErrorResponse) => {
        if(error.status == 404){
          this.formGroup.setErrors({userNotExist:true})
        this.openSnackBar("User with that e-mail not exists in system.", "")

        }
      }
    })
  }

  openSnackBar(message: string, action: string) {
    this._snackBar.open(message, action,  {
      duration: 3500
    });
  }

}
