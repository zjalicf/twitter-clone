import { Component, Input, OnInit } from '@angular/core';
import { Tweet } from 'src/app/models/tweet.model';

@Component({
  selector: 'app-tweet-list',
  templateUrl: './tweet-list.component.html',
  styleUrls: ['./tweet-list.component.css']
})
export class TweetListComponent implements OnInit {

  constructor() { }

  @Input() tweets: Tweet[] = []

  ngOnInit(): void {
  }

}
