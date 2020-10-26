export interface FileResponse {
    success: boolean;
    data? : {
        message: string;
    }
}

export  interface FileParams {
    UserID: string;
    file: File;
    description: string;
    filename?: string;
}


export  interface EditFileParams {
    UserID: string;
    description: string;
    filename?: string;
}