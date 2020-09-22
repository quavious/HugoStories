package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/quavious/GoStories/model"
)

var agent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36 Edg/85.0.564.51"

//ScrapeURL returns scraped texts.
func ScrapeURL(url string) []string {
	elements := []string{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []string{}
	}
	req.Header.Add("User-Agent", agent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	var title string

	title, _ = doc.Find("tiara-page").Attr("data-tiara-tags")
	elements = append(elements, title)
	doc.Find(".item_type_text").Each(func(i int, s *goquery.Selection) {
		sample := strings.TrimSpace(s.Text())
		if len(sample) != 0 {
			elements = append(elements, sample)
		}
	})

	return elements
}

//Translate returns text translated to us.
func Translate(elements []string) []string {
	translated := []string{}
	options := []chromedp.ExecAllocatorOption{
		chromedp.UserAgent(agent),
		chromedp.Flag("Headless", true),
	}
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	ctxt, cancel := context.WithTimeout(ctx, time.Minute*30)
	defer cancel()

	for _, element := range elements {
		var source string

		var apiURL string = "https://translate.google.co.kr/?hl=ko#view=home&op=translate&sl=ko&tl=en&text=" + element

		if err := chromedp.Run(ctxt,
			chromedp.Navigate(apiURL),
			chromedp.Sleep(time.Second*4),
			//chromedp.OuterHTML("html", &source, chromedp.ByQuery),
			chromedp.ActionFunc(func(ctx context.Context) error {
				node, err := dom.GetDocument().Do(ctx)
				if err != nil {
					return err
				}
				source, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
				return err
			}),
		); err != nil {
			fmt.Println(err)
			chromedp.Cancel(ctxt)
			return []string{}
		}
		time.Sleep(time.Second * 2)
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(source))
		sample := doc.Find(".tlid-translation").Text()
		//fmt.Println(sample)
		translated = append(translated, strings.TrimSpace(sample))
		timeInterval := rand.Intn(8) + 8
		time.Sleep(time.Second * time.Duration(timeInterval))
	}
	return translated
}

//FetchImage returns a splash image from apiKEY.
func FetchImage(apiURL string) string {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return ""
	}
	req.Header.Add("User-Agent", agent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	bytes := bytes.NewBuffer(nil)
	_, err = io.Copy(bytes, resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	thumbnail := model.ImageURL{}
	err = json.Unmarshal(bytes.Bytes(), &thumbnail)

	return thumbnail.Urls.Thumbnail
}
