
export  interface LoginResponse {
  EmailAddress: string;
  Password: string;
  UserID: string;
  IsAdmin: boolean;
}

export interface UserParams {
    emailAddress: string;
    firstName: string;
    lastName: string;
    password: string;
    isAdmin: boolean;
}

export interface RegisterResponse {
    FirstName: string;
    LastName: string;
    EmailAddress: string;
    IsAdmin: string;
}

