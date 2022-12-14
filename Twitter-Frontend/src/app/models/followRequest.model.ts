export class FollowRequest {
    sender: string = "";
    receiver: string = "";
    status: string = "";

    FollowRequest(sender: string, receiver: string, status: string) {
        this.sender = sender;
        this.receiver = receiver;
        this.status = status;
    }
}