package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

var (

	letters        = []string{}
)

type photo struct {
	Photo string `json:"img"`
}

func addthis() {
	pos := " .:!|l1G0@"
	for _, x := range pos {
		for y := 0; y < 257/(len(pos)); y++ {
			letters = append(letters, string(x))
		}
	}
}

func openThis(f io.Reader) {
	//this open the image and print the pixels
	img, err := png.Decode(f)
	if err != nil {
		log.Println(err)
		return
	}

	division := 7
	limitY, limitX := img.Bounds().Max.Y/division, img.Bounds().Max.X/division
	for y := img.Bounds().Min.Y; y < limitY; y++ {
		for x := img.Bounds().Min.X; x < limitX; x++ {

			r, g, b, _ := img.At(x*division, y*division).RGBA()
			fmt.Print(letters[int(((r/257)+(g/257)+(b/257))/3)%len(letters)])
		}
		fmt.Println()

	}
}

func imagePNG(input string) io.Reader {
	return base64.NewDecoder(base64.StdEncoding, strings.NewReader(input))
}
func saycheese(_ http.ResponseWriter, r *http.Request) {

	// decode the bodyrequest
	var conf photo

	json.NewDecoder(r.Body).Decode(&conf)

	imageData := imagePNG(conf.Photo[31:])

	openThis(imageData)

}


func main() {

	// clear the console
	if err := exec.Command("clear").Run(); err != nil {
		exec.Command("cls").Run()
	}
	// start the interface
	
	fmt.Println("\033[34mstarting  server \033[0m")


	addthis()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "view"+r.URL.Path)
	})
	http.HandleFunc("/photo", saycheese)
	http.ListenAndServe(":8000", nil)

}
