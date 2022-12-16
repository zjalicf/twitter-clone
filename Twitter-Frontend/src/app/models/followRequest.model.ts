export class FollowRequest {
    id: string = "";
    requester: string = "";
    receiver: string = "";
    status: string = "";

    FollowRequest(id: string, requester: string, receiver: string, status: string) {
        this.id = id
        this.requester = requester;
        this.receiver = receiver;
        this.status = status;
    }
}