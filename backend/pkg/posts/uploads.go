package posts

import (
    "io"
    "net/http"
    "os"
    "time"
    "path/filepath"
)

func handleImageUpload(r *http.Request) string {
    return handleFileUpload(r, "image", "uploads/images")
}

func handleGIFUpload(r *http.Request) string {
    return handleFileUpload(r, "gif", "uploads/gifs")
}

func handleFileUpload(r *http.Request, formFieldName string, uploadDir string) string {
    err := os.MkdirAll(uploadDir, os.ModePerm)
    if err != nil {
        return ""
    }

    file, _, err := r.FormFile(formFieldName)
    if err != nil {
        return ""
    }
    defer file.Close()

    fileName := time.Now().Format("20060102150405") + filepath.Ext(formFieldName)
    filePath := filepath.Join(uploadDir, fileName)

    outFile, err := os.Create(filePath)
    if err != nil {
        return ""
    }
    defer outFile.Close()

    _, err = io.Copy(outFile, file)
    if err != nil {
        return ""
    }

    return filePath
}
