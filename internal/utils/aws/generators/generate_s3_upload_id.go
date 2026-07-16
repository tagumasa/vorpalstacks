package generators

// Package generators provides ID generation utilities for vorpalstacks AWS resources.

// GenerateS3UploadID generates a unique S3 multipart upload ID.
//
// Returns:
//   - string: The generated upload ID
//   - error: An error if ID generation fails
func GenerateS3UploadID() (string, error) {
	return generateIDWithPrefix("UPLOAD", 28)
}
