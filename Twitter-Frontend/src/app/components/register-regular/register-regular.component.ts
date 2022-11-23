import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormControl, FormGroup, ValidationErrors, ValidatorFn, Validators } from '@angular/forms';
import { User } from 'src/app/models/user.model';
import { AuthService } from 'src/app/services/auth.service';
import { PasswordSpecialCharacterValidator, PasswordStrenghtValidator } from 'src/app/services/customValidators';

@Component({
  selector: 'app-register-regular',
  templateUrl: './register-regular.component.html',
  styleUrls: ['./register-regular.component.css']
})
export class RegisterRegularComponent implements OnInit {

  formGroup: FormGroup = new FormGroup({
    firstName: new FormControl(''),
    lastName: new FormControl(''),
    gender: new FormControl(''), 
    age: new FormControl(''),
    residence: new FormControl(''),
    email: new FormControl(''),
    username: new FormControl(''),
    password: new FormControl('')
  });

  genders: string[] = [
    'Male',
    'Female'
  ];

  constructor(private authService: AuthService,
              private formBuilder: FormBuilder) { }

  // @ts-ignore
  formGroup: FormGroup;
  submitted = false;

  ngOnInit(): void {
    this.formGroup = this.formBuilder.group({
      firstName: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(20)]],
      lastName: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(20)]],
      gender: ['', [Validators.required]],
      age: ['', [Validators.required, Validators.min(1), Validators.max(100)]],
      residence: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(35)]],
      email: ['', [Validators.required, Validators.email, Validators.minLength(3), Validators.maxLength(35)]],
      username: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(30)]],
      password: ['', [Validators.required, Validators.minLength(8), Validators.maxLength(30), PasswordStrenghtValidator(), PasswordSpecialCharacterValidator()]],
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

    let registerUser: User = new User();

    registerUser.firstName = this.formGroup.get("firstName")?.value;
    registerUser.lastName = this.formGroup.get("lastName")?.value;
    registerUser.gender = this.formGroup.get("gender")?.value;
    registerUser.age = this.formGroup.get("age")?.value;
    registerUser.residence = this.formGroup.get("residence")?.value;
    registerUser.email = this.formGroup.get("email")?.value;
    registerUser.username = this.formGroup.get("username")?.value;
    registerUser.password = this.formGroup.get("password")?.value;

    this.authService.Register(registerUser)
      .subscribe({
        next: (data: User) => {
          console.log(data);
          alert("You have been successfully registered to Twitter");
        },
        error: (error) => {
          console.log(error)
        }
      });
  }

}

// export function createPasswordStrenghtValidator(): ValidatorFn {
//   return (control: AbstractControl) : ValidationErrors | null => {
//     const value = control.value;

//     if (!value) {
//       return null;
//     }

//     const hasUpperCase = /[A-Z]+/.test(value)

//     const hasLowerCase = /[a-z]+/.test(value);

//     const hasNumeric = /[0-9]+/.test(value);

//     const passwordValid = hasUpperCase && hasLowerCase && hasNumeric;

//     return !passwordValid ? {passwordStrength:true}: null;
//   }
// }