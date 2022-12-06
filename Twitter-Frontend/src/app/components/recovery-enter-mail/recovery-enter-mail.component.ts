import { HttpErrorResponse } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
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
    private recoveryService: RecoveryPasswordService
  ) { }

  

  ngOnInit(): void {
    this.formGroup = this.formBuilder.group({
      email: ['', [Validators.required, Validators.email]],
    })
  }

  get f(): { [key: string]: AbstractControl } {
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
        this.router.navigate(['/Recovery-Token'])
      },
      error: (error: HttpErrorResponse) => {
        if(error.status == 404){
          this.formGroup.setErrors({userNotExist:true})
        }
      }
    })
  }

}
