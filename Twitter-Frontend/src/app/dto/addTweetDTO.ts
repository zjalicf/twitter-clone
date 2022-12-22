export class AddTweetDTO {
    text: string = "";

    Tweet(text: string) {
        this.text = text;
    }
}