import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from 'src/app/services/auth.service';
import { RecoveryPasswordService } from 'src/app/services/recoveryPassword.service';

@Component({
  selector: 'app-recovery-new-passwords',
  templateUrl: './recovery-new-passwords.component.html',
  styleUrls: ['./recovery-new-passwords.component.css']
})
export class RecoveryNewPasswordsComponent implements OnInit {

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
      token: ['', [Validators.required, Validators.maxLength(50)]],
    })
  }

  get f(): { [key: string]: AbstractControl } {
    return this.formGroup.controls;
  }

  onSubmit(){
    
  }

}
