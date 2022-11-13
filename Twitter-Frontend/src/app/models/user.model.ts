export class User {
    userId: number = 0;
    username: string = "";
    password: string = "";
    avatar: string = "";
    email: string = "";
    dateOfRegistration: Date = new Date('2022-05-03');
    description: string = "";
    displayName: string = "";

    User(userId: number, username: string, password: string, avatar: string, email: string,  dateOfRegistration: Date, description: string, displayName: string) {
        this.userId = userId;
        this.username = username;
        this.password = password;
        this.avatar = avatar;
        this.email = email;
        this.dateOfRegistration = dateOfRegistration;
        this.description = description;
        this.displayName = displayName;
    }
}