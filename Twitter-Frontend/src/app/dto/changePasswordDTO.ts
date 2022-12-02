export class ChangePasswordDTO {
    currentPassword: string = "";
    newPassword: string = "";

    ChangePasswordDTO(currentPassword: string, newPassword: string) {
        this.currentPassword = currentPassword;
        this.newPassword = newPassword;
    }
}