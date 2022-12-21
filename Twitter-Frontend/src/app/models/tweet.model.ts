export class Tweet {
    id: string = "";
    text: string = "";
    username: string = "";
    favorite_count: number = 0;

    Tweet(id: string, text: string, username: string, favorite_count: number) {
        this.id = id
        this.text = text;
        this.username = username;
        this.favorite_count = favorite_count
    }
}