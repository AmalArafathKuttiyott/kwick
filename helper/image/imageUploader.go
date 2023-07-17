package image

import (
	"fmt"
	"io"
	"kwick/helper/database"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// func ImageUploader(ctx *gin.Context, i uint) bool {
// 	err := ctx.Request.ParseMultipartForm(64 << 20) // 64MB
// 	if err != nil {
// 		return false
// 	}
// 	form := ctx.Request.MultipartForm
// 	files := form.File["images"]
// 	number := uint(i)
// 	productId := strconv.FormatUint(uint64(number), 10)

// 	//Get the productId

//		// Iterate through each uploaded file
//		for _, file := range files {
//			// Save the file to disk
//			fileName := "public/images/" + productId + file.Filename
//			_, err := os.Create(fileName)
//			if err != nil {
//				return false
//			}
//			num := fmt.Sprint(i)
//			added := database.AddImages(i, "/static/"+num+file.Filename)
//			if !added {
//				return false
//			}
//		}
//		return true
//	}
func ImageUploader(ctx *gin.Context, i uint) bool {
	// Parse the multipart form data
	err := ctx.Request.ParseMultipartForm(32 << 20) // 32MB maximum file size
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return false
	}
	for _, fileHeaders := range ctx.Request.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			// Open the uploaded file
			file, err := fileHeader.Open()
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return false
			}
			defer file.Close()
			id := database.GetProductId(ctx.Request.FormValue("productName"))
			num := strconv.Itoa(int(id))
			fileName := "public/images/" + num + fileHeader.Filename

			// Create a new file in the server's local storage
			destinationFile, err := os.Create(fileName)
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return false
			}
			// add image with id of product as image name in database
			defer destinationFile.Close()

			// Copy the uploaded file's data to the destination file
			_, err = io.Copy(destinationFile, file)
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return false
			}

			database.AddImages(id, "/static/"+num+fileHeader.Filename)
			fmt.Printf("Image '%s' uploaded successfully!\n", fileHeader.Filename)
		}
	}
	return true
}
