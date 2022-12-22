import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { User } from 'src/app/models/user.model';

export interface DialogData {
  testUsers: User[]
}

@Component({
  selector: 'app-tweet-likes-dialog',
  templateUrl: './tweet-likes-dialog.component.html',
  styleUrls: ['./tweet-likes-dialog.component.css']
})
export class TweetLikesDialogComponent {

  constructor(
    public dialogRef: MatDialogRef<TweetLikesDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: User,
  ) {}

  onOkClick(): void {
    this.dialogRef.close();
    console.log(this.data)
  }

}
