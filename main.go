package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/quavious/GoStories/confidential"
	"github.com/quavious/GoStories/controller"
	"github.com/quavious/GoStories/model"
)

func main() {
	splashKEY := confidential.SplashKEY
	file, err := os.Open("urls.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	item := model.Storage{}
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.Unmarshal(buf.Bytes(), &item)

	for _, url := range item.Item {
		container := controller.ScrapeURL(url)
		fmt.Println(url)
		if title := container[0]; len(title) == 0 {
			continue
		}
		translated := controller.Translate(container)
		title := translated[0]
		textContent := strings.Join(translated[1:], "  \n\n")
		tags := strings.Split(strings.ReplaceAll(title, ", ", ","), ",")

		var imageURL string
		for _, tag := range tags {
			splashURL := "https://api.unsplash.com/photos/random?client_id=" + splashKEY + "&query=" + tag + "&orientation=landscape"
			imageURL = controller.FetchImage(splashURL)
			if len(imageURL) > 0 {
				break
			}
		}
		if len(imageURL) < 10 {
			splashURL := "https://api.unsplash.com/photos/random?client_id=" + splashKEY + "&query=" + "life" + "&orientation=landscape"
			imageURL = controller.FetchImage(splashURL)
		}

		var postContent string
		postContent += "---\n"
		postContent += "title : " + title + "\n"
		postContent += "subtitle : " + "Story#" + time.Now().Format("200601021504") + "\n"
		postContent += "draft : false\n"
		postContent += "tags :\n"
		postContent += " - life\n"
		postContent += " - daily\n"
		for _, tag := range tags {
			postContent += " - " + tag + "\n"
		}
		postContent += "date : " + time.Now().Format("2006-01-02T15:04:05") + "+0900" + "\n"
		postContent += "toc : false\n"
		postContent += "images : \n"
		postContent += "thumbnail : " + imageURL + "\n"
		postContent += "---\n"
		postContent += textContent

		err = os.Chdir("C:/Hugo/hugostory")
		if err != nil {
			log.Fatalln(err)
			return
		}

		fmt.Println(os.Getwd())
		path := time.Now().Format("2006\\01\\02")
		filename := time.Now().Format("20060102150405") + ".md"
		cmd := exec.Command("hugo", "new", "posts\\"+path+"\\"+filename)
		cmd.Dir = "C:\\Hugo\\hugostory"

		output, err := cmd.Output()

		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(output))

		err = ioutil.WriteFile("C:\\Hugo\\hugostory\\content\\posts\\"+path+"\\"+filename, []byte(postContent), 0644)

		if err != nil {
			log.Println(err)
			return
		}
		interval := rand.Intn(120) + 120

		fmt.Println("Make File Successful")
		time.Sleep(time.Second * time.Duration(interval))
	}
}
