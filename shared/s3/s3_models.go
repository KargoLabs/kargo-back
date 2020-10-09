package s3

// StructWithUploadURL is the struct handler to return upload URL
type StructWithUploadURL struct {
	Struct                interface{} `json:"struct"`
	UploadProfilePhotoURL string      `json:"upload_profile_photo_url"`
}

// StructWithGetURL is the struct handler to return get URL
type StructWithGetURL struct {
	Struct             interface{} `json:"struct"`
	GetProfilePhotoURL string      `json:"get_profile_photo_url"`
}
