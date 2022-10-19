package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/Pandahoro/kutego-api/pkg/swagger/server/models"
	"github.com/Pandahoro/kutego-api/pkg/swagger/server/restapi"
	"github.com/Pandahoro/kutego-api/pkg/swagger/server/restapi/operations"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/google/go-github/v37/github"
)

func main() {

	// Initialize Swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewKutegoAPIAPI(swaggerSpec)
	// Use swaggerUI instead of reDoc on /docs
	api.UseSwaggerUI()

	server := restapi.NewServer(api)

	defer func() {
		if err := server.Shutdown(); err != nil {
			// error handle
			log.Fatalln(err)
		}
	}()

	server.Port = 8080

	api.CheckHealthHandler = operations.CheckHealthHandlerFunc(Health)

	api.GetCatNameHandler = operations.GetCatNameHandlerFunc(GetCatName)

	api.GetCatsHandler = operations.GetCatsHandlerFunc(GetCats)

	api.GetCatRandomHandler = operations.GetCatRandomHandlerFunc(GetCatRandom)

	// Start server which listening
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

}

//Health route returns OK
func Health(operations.CheckHealthParams) middleware.Responder {
	return operations.NewCheckHealthOK().WithPayload("OK")
}

//GetCatName returns Cat image (png)
func GetCatName(cat operations.GetCatNameParams) middleware.Responder {

	var URL string
	if cat.Name != "" {
		URL = "https://github.com/Pandahoro/cats/raw/main/" + cat.Name + ".gif"
	} else {
		//by default we return Gandalf cat
		URL = "https://github.com/Pandahoro/cats/raw/main/SadCatto.gif"
	}

	response, err := http.Get(URL)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		srcImage, _ := getFireCatError("Ooops, Error")
		return operations.NewGetCatNameOK().WithPayload(convertImgToIoCloser(srcImage))
	}

	srcImage, _, err := image.Decode(response.Body)
	if err != nil {
		log.Fatalf("failed to Decode image: %v", err)
	}

	if cat.Size != nil {
		srcImage = resizeImage(srcImage, *cat.Size)
	}

	return operations.NewGetCatNameOK().WithPayload(convertImgToIoCloser(srcImage))
}

/**
Display Cat list with optional filter
*/
func GetCats(cat operations.GetCatsParams) middleware.Responder {

	catList := GetCatsList()

	if cat.Name != nil {
		var arr []*models.Cat
		for key, value := range catList {
			if value.Name == *cat.Name {
				arr = append(arr, catList[key])
				return operations.NewGetCatsOK().WithPayload(arr)
			}
		}
	}

	return operations.NewGetCatsOK().WithPayload(catList)
}

/**
Display a random Cat Image
*/
func GetCatRandom(cat operations.GetCatRandomParams) middleware.Responder {
	var URL string

	// Get Cats List
	arr := GetCatsList()

	// Get a Random Index
	rand.Seed(time.Now().UnixNano())
	var index int
	index = rand.Intn(len(arr) - 1)

	URL = "https://github.com/Pandahoro/cats/raw/main/" + arr[index].Name + ".gif"

	response, err := http.Get(URL)
	if err != nil {
		fmt.Println("error")
		srcImage, _ := getFireCatError("Ooops, Error")
		return operations.NewGetCatNameOK().WithPayload(convertImgToIoCloser(srcImage))
	}
	defer response.Body.Close()

	srcImage, _, err := image.Decode(response.Body)
	if err != nil {
		log.Fatalf("failed to Decode image: %v", err)
	}

	if cat.Size != nil {
		srcImage = resizeImage(srcImage, *cat.Size)
	}

	return operations.NewGetCatNameOK().WithPayload(convertImgToIoCloser(srcImage))
}

/**
Display Fire Cat with a message (error)
*/
func getFireCatError(message string) (image.Image, error) {

	// Open local file
	file, err := os.Open("./assets/fire-cat.gif")
	if err != nil {
		log.Fatalf("failed to Open fire-cat image: %v", err)
	}

	srcImage, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("failed to Decode image: %v", err)
		return srcImage, err
	}
	// Add Text on Cat
	srcImage, err = TextOnCat(srcImage, "Ooops, Error! It's on fire!")

	// Resize Image
	srcImage = resizeImage(srcImage, "medium")
	if err != nil {
		log.Fatalf("failed to put Text on Cat: %v", err)
		return srcImage, err
	}

	return srcImage, nil
}

/**
Get Cats List from Scraly repository
*/
func GetCatsList() []*models.Cat {

	client := github.NewClient(nil)

	// list public repositories for org "github"
	ctx := context.Background()
	// list all repositories for the authenticated user
	_, directoryContent, _, err := client.Repositories.GetContents(ctx, "Pandahoro", "cats", "/", nil)
	if err != nil {
		fmt.Println(err)
	}

	var arr []*models.Cat

	for _, c := range directoryContent {
		if *c.Name == ".gitignore" || *c.Name == "README.md" {
			continue
		}

		var name string = strings.Split(*c.Name, ".")[0]

		arr = append(arr, &models.Cat{name, *c.Path, *c.DownloadURL})

	}

	return arr
}

/**
Resize Image
*/
func resizeImage(srcImage image.Image, size string) image.Image {

	var height int
	switch size {
	case "x-small":
		height = 50
	case "small":
		height = 100
	case "medium":
		height = 300
	default:
		// Mouhouhahaha!
		height = 1000
	}

	// Resize the cropped image to width = 200px preserving the aspect ratio.
	srcImage = imaging.Resize(srcImage, 0, height, imaging.Lanczos)

	return srcImage

}

/**
Convert Image to io.close (for reply format)
*/
func convertImgToIoCloser(srcImage image.Image) io.ReadCloser {
	encoded := &bytes.Buffer{}
	png.Encode(encoded, srcImage)

	return ioutil.NopCloser(encoded)
}

/**
Add text on Image
*/
func TextOnCat(bgImage image.Image, text string) (image.Image, error) {

	imgWidth := bgImage.Bounds().Dx()
	imgHeight := bgImage.Bounds().Dy()

	dc := gg.NewContext(imgWidth, imgHeight)
	dc.DrawImage(bgImage, 0, 0)

	if err := dc.LoadFontFace("assets/FiraSans-Light.ttf", 50); err != nil {
		return nil, err
	}

	x := float64((imgWidth / 2))
	y := float64((imgHeight / 12))
	maxWidth := float64(imgWidth) - 60.0
	dc.SetColor(color.Black)
	dc.DrawStringWrapped(text, x, y, 0.5, 0.5, maxWidth, 1.5, gg.AlignRight)

	return dc.Image(), nil
}
