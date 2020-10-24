package models

// File is the structure to represent a file
type FileInfo struct {
	FileName    string `json:"file_name"`
	Description string `json:"description"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}

// UpdateFileInfo is the file info for file update
type UpdateFileInfo struct {
	Description string `json:"description"`
}

// DownloadFileInfo is the download file Info
type DownloadFileInfo struct {
	PresignedURL string `json:"presignedURL,omitempty"`
}
