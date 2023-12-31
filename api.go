package mapdf

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func PostFileWithHeader[T any](targetURL, tokenKey, tokenValue, filePath, formdataName string) (result T, err error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	// Create a buffer to store the form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create a part for the file
	filePart, err := writer.CreateFormFile(formdataName, filePath)
	if err != nil {
		return
	}

	// Copy the file content to the part
	_, err = io.Copy(filePart, file)
	if err != nil {
		return
	}

	// Add the token as a form field
	writer.WriteField(tokenKey, tokenValue)

	// Close the form data writer
	writer.Close()

	// Create a new HTTP request with the form data
	request, err := http.NewRequest("POST", targetURL, &requestBody)
	if err != nil {
		return
	}

	// Set the content type header
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// Set the token header
	request.Header.Set(tokenKey, tokenValue)

	// Create an HTTP client
	client := &http.Client{}

	// Perform the request
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	if err = json.Unmarshal(responseBody, &result); err != nil {
		return
	}

	return
}
