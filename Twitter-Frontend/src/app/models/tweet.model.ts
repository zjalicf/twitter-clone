export class Tweet {
    id: string = "";
    text: string = "";
    username: string = "";
    favorite_count: number = 0;
    image: boolean = false;

    Tweet(id: string, text: string, username: string, favorite_count: number, image: boolean) {
        this.id = id
        this.text = text;
        this.username = username;
        this.favorite_count = favorite_count;
        this.image = image;
    }
}