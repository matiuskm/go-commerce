package helpers

import (
	"context"
	"fmt"
	"math/rand"
	"mime/multipart"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/resend/resend-go/v2"
)

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func generateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateOrderNumber() string {
	datePart := time.Now().Format("20060102")
	randomPart := generateRandomString(6)
	return fmt.Sprintf("ORD-%s-%s", datePart, randomPart)
}

func UploadToCloudinary(file multipart.File, filename string) (string, error) {
    cld, err := cloudinary.NewFromParams(
        os.Getenv("CLOUDINARY_CLOUD_NAME"),
        os.Getenv("CLOUDINARY_API_KEY"),
        os.Getenv("CLOUDINARY_API_SECRET"),
    )
    if err != nil {
        return "", err
    }

    uploadRes, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
        PublicID: "product_" + filename,
        Folder:   "go-commerce",
    })
    if err != nil {
        return "", err
    }

    return uploadRes.SecureURL, nil
}

func SendEmail(to string, subject string, body string) error {
    client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

	params := &resend.SendEmailRequest{
        From:    "GoCommerce <no-reply@arunikadigital.com>",
        To:      []string{},
        Html:    body,
        Subject: subject,
    }

	sent, err := client.Emails.Send(params)
	if err!= nil {
		return err
	}
	fmt.Println(sent)
    return nil
}