import { Component, Input, OnInit } from '@angular/core';
import { Tweet } from 'src/app/models/tweet.model';
import { TweetService } from 'src/app/services/tweet.service';

@Component({
  selector: 'app-tweet-item',
  templateUrl: './tweet-item.component.html',
  styleUrls: ['./tweet-item.component.css']
})
export class TweetItemComponent implements OnInit {

  constructor(private tweetService :TweetService) { }

   @Input() tweet: Tweet = new Tweet();

  ngOnInit(): void {
    
  }

}
