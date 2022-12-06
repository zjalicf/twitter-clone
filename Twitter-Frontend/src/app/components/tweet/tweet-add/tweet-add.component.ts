import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Tweet } from 'src/app/models/tweet.model';
import { TweetService } from 'src/app/services/tweet.service';

@Component({
  selector: 'app-tweet-add',
  templateUrl: './tweet-add.component.html',
  styleUrls: ['./tweet-add.component.css']
})
export class TweetAddComponent implements OnInit {

  formGroup: FormGroup = new FormGroup({
    text: new FormControl('')
  });

  constructor(private formBuilder: FormBuilder,
              private tweetService: TweetService) { }

  submitted = false;

  ngOnInit(): void {
    this.formGroup = this.formBuilder.group({
      text: ['', [Validators.required, Validators.minLength(5), Validators.maxLength(280)]] // Validators.pattern('[-_a-zA-Z0-9]*')
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

    let addTweet: Tweet = new Tweet();

    addTweet.text = this.formGroup.get("text")?.value;

    this.tweetService.AddTweet(addTweet)
      .subscribe(data => {
        alert("Tweet succesfully created!")
      })
  }

}
