import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { AddTweetDTO } from 'src/app/dto/addTweetDTO';
import { Tweet } from 'src/app/models/tweet.model';
import { TweetService } from 'src/app/services/tweet.service';


@Component({
  selector: 'app-tweet-add',
  templateUrl: './tweet-add.component.html',
  styleUrls: ['./tweet-add.component.css']
})
export class TweetAddComponent implements OnInit {

  formGroup: FormGroup = new FormGroup({
    text: new FormControl(''),
    image: new FormControl('')
  });

  file!: File;
  formData = new FormData();

  constructor(private formBuilder: FormBuilder,
              private tweetService: TweetService,
              private router: Router,
              private http: HttpClient) { }

  submitted = false;

  ngOnInit(): void {
    this.formGroup = this.formBuilder.group({
      text: ['', [Validators.required, Validators.minLength(5), Validators.maxLength(280)]], // Validators.pattern('[-_a-zA-Z0-9]*')
      image: ['']
    })
  }

  get tweetForm(): { [key: string]: AbstractControl } {
    return this.formGroup.controls;
  }

  onSubmit() {
    this.submitted = true;

    if (this.formGroup.invalid) {
      return;
    }

    let addTweet: AddTweetDTO = new AddTweetDTO();

    addTweet.text = this.formGroup.get("text")?.value;

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
    if (this.file.type === 'image/jpeg' || this.file.type === 'image/jpg') {
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
