package webserver

type FileUpload struct {
	Filename string
	Size     int64
	Content  []byte
}
