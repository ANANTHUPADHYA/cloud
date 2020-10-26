export interface UserResponse {
        UserID: string;
        FirstName: string;
        LastName: string;
        IsAdmin: boolean;
        files: any;
        EmailAddress: string,
    }


export interface User {
    lastname: string;
    files: 
        {
            downloadUrl: string;
            modifiedDate: string;
            description: string;
            fileName: string;
            uploadedDate: string;
        }[];
    userId: string;
    firstname: string;
    username: string;
    isAdmin: boolean;
    password_SHA512: string;
}

export interface AllUsers {
    success: boolean;
    users: User[];
}

export interface DeleteResponse {
    Success: boolean;
    Message: string;
}