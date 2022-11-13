import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-register-regular',
  templateUrl: './register-regular.component.html',
  styleUrls: ['./register-regular.component.css']
})
export class RegisterRegularComponent implements OnInit {

  genders: string[] = [
    'Male',
    'Female',
    'Other'
  ];

  constructor() { }

  ngOnInit(): void {
  }

}
