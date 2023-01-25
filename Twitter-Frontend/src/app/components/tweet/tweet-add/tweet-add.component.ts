import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { AddTweetDTO } from 'src/app/dto/addTweetDTO';
import { Tweet } from 'src/app/models/tweet.model';
import { User } from 'src/app/models/user.model';
import { TweetService } from 'src/app/services/tweet.service';
import { UserService } from 'src/app/services/user.service';


@Component({
  selector: 'app-tweet-add',
  templateUrl: './tweet-add.component.html',
  styleUrls: ['./tweet-add.component.css']
})
export class TweetAddComponent implements OnInit {

  constructor(
    private formBuilder: FormBuilder,
    private tweetService: TweetService,
    private userService: UserService,
    private router: Router,
    private http: HttpClient
  ) 
  { }

  tweetFormGroup: FormGroup = new FormGroup({
    text: new FormControl(''),
    image: new FormControl('')
  });

  advertisementFormGroup: FormGroup = new FormGroup({
    residence: new FormControl(''),
    gender: new FormControl(''),
    age_from: new FormControl(''),
    age_to: new FormControl('')
  })

  file!: File;
  formData = new FormData();

  isChecked = false;
  submittedTweet = false;
  submittedAdvertisement = false;
  user: User = new User();

  ngOnInit(): void {
    this.tweetFormGroup = this.formBuilder.group({
      text: ['', [Validators.required, Validators.minLength(5), Validators.maxLength(280)]], // Validators.pattern('[-_a-zA-Z0-9]*')
      image: ['']
    })

    this.advertisementFormGroup = this.formBuilder.group({
      residence: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(35)]],
      gender: ['', [Validators.required]],
      age_from: ['', [Validators.required, Validators.min(5), Validators.max(100)]],
      age_to: ['', [Validators.required, Validators.min(5), Validators.max(100)]]
    })

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

  isBusiness(): boolean {
    if (this.user.userType == "Business") {
      return true;
    } else {
      return false
    }
  }

  get tweetForm(): { [key: string]: AbstractControl } {
    return this.tweetFormGroup.controls;
  }

  get advertisementForm(): { [key: string]: AbstractControl } {
    return this.advertisementFormGroup.controls;
  }

  check() {
    if (this.isChecked == true) {
      this.isChecked = false
    } else {
      this.isChecked = true
    }
  }

  onSubmit() {
    this.submittedTweet = true;
    this.submittedAdvertisement = true;

    if (this.tweetFormGroup.invalid) {
      return;
    }

    if (this.isChecked == true) {
      if (this.advertisementFormGroup.invalid) {
        return;
      }
    }

    // let add advertisementDTO = new AddAdvertisementDTO();
    let addTweet: AddTweetDTO = new AddTweetDTO();

    addTweet.text = this.tweetFormGroup.get("text")?.value;
    addTweet.advertisement = true
    console.log(addTweet)
    console.log(JSON.stringify(addTweet))
    this.formData.append("json", JSON.stringify(addTweet))
    this.tweetService.AddTweet(this.formData)
      .subscribe({
        next: (data: Tweet) => {
          this.router.navigate(['/Main-Page']);
        },
        error: (error) => {
          console.log(error);
        }
      })
  }

  getFile(event: any) {
    console.log("Desio se event")
    this.file = event.target.files[0];
    let fileType = this.file.type.split("/")
    if (fileType[0] === "image") {
        this.formData.append('image', this.file);
    } else {
        console.log('Wrong file type. Only jpeg images are allowed.');
    }
  }

  uploadFile() {
    this.http.post('/api/upload', this.formData).subscribe(response => {
      console.log(response);
    });
  }

}
