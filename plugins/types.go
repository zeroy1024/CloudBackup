package plugins

import "time"

type UserToken struct {
	UserId       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresTime  string `json:"expire_time"`
	DriveId      string `json:"default_drive_id"`
}

type Config struct {
	BackupPath            []string  `json:"backup_path"`
	CloudBackupFolderName string    `json:"cloud_backup_folder_name"`
	UserToken             UserToken `json:"user_token"`
}

type File struct {
	DriveId         string     `json:"drive_id"`
	DomainId        string     `json:"domain_id"`
	FileId          string     `json:"file_id"`
	Name            string     `json:"name"`
	Type            string     `json:"type"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Hidden          bool       `json:"hidden"`
	Starred         bool       `json:"starred"`
	Status          string     `json:"status"`
	ParentFileId    string     `json:"parent_file_id"`
	EncryptMode     string     `json:"encrypt_mode"`
	FileExtension   string     `json:"file_extension"`
	MimeType        string     `json:"mime_type"`
	MimeExtension   string     `json:"mime_extension"`
	Size            int        `json:"size"`
	UploadId        string     `json:"upload_id"`
	Crc64Hash       string     `json:"crc64_hash"`
	ContentHash     string     `json:"content_hash"`
	ContentHashName string     `json:"content_hash_name"`
	DownloadUrl     string     `json:"download_url"`
	Url             string     `json:"url"`
	Category        string     `json:"category"`
	PunishFlag      int        `json:"punish_flag"`
	Trashed         bool       `json:"trashed"`
	PartInfoList    []PartInfo `json:"part_info_list"`
	Exist           bool       `json:"exist"`
}

type FileList struct {
	Items      []File `json:"items"`
	NextMarker string `json:"next_marker"`
}

type PartInfo struct {
	PartNumber        int    `json:"part_number"`
	UploadUrl         string `json:"upload_url"`
	InternalUploadUrl string `json:"internal_upload_url"`
	ContentType       string `json:"content_type"`
}

type BatchRequest struct {
	ID      string            `json:"id"`
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    map[string]string `json:"body"`
}
