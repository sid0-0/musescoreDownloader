package dl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

func getAuthTokenFromChunk(chunkUrl string) (string, error) {
	res, err := http.Get(chunkUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", errors.New("GET Request failed")
	}
	bodyContentInBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var bodyString string = string(bodyContentInBytes)

	// .mainServer,g\),".*?"
	fmt.Println("url: ", chunkUrl)
	compiledRegex, err := regexp.Compile(`[a-z,A-Z],.\.mainServer,.\),\s*"(.*?)"`)
	if err != nil {
		return "", errors.New(fmt.Sprint("Faulty regex", err))
	}
	foundStringMap := compiledRegex.FindAllStringSubmatch(bodyString, 1)
	if len(foundStringMap) == 0 {
		return "", errors.New("Token not found")
	}

	token := foundStringMap[0][1]
	fmt.Println("Token found: ", token)

	return token, nil
}

func downloadFile(sheetId string, pageNumber int, headers map[string]string) error {
	fmt.Println("url: ", fmt.Sprintf(`https://musescore.com/api/jmuse?id=6102579&index=%d&type=img&v2`, pageNumber))
	res, err := http.Get(fmt.Sprintf(`https://musescore.com/api/jmuse?id=6102579&index=%d&type=img&v2`, pageNumber))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err != nil {
		return err
	}

	// bodyAsString, _ := io.ReadAll(res.Body)

	var data struct {
		Result string `json:"result"`
		Status int    `json:"status"`
		Error  string `json:"error"`
		Info   struct {
			Url string `json:"url"`
		}
	}
	json.NewDecoder(res.Body).Decode(&data)
	// fmt.Println("data: ", data)
	// fmt.Println("data: ", string(bodyAsString))
	fmt.Println("data: ", data)
	// fmt.Println("./downloads/" + splits[len(splits)-1])
	// if err := r.Save("./downloads/test"); err != nil {
	// 	fmt.Println("Failed to download file ", err)
	// }
	return nil
}

func DownloadFromUrl(url string) error {
	fmt.Println(url)
	c := colly.NewCollector()
	c.OnHTML("link[rel=\"preload\"][as=\"script\"]", func(e *colly.HTMLElement) {
		var href string = e.Attr("href")
		var withoutParams string = strings.Split(href, "?")[0]
		// filtering out non svg url
		if strings.Count(href, "/") < 8 {
			return
		}
		token, err := getAuthTokenFromChunk(withoutParams)
		if err != nil {
			panic(fmt.Sprintln("Failed to get token. ", err))
		}
		downloadFile(url, 0, map[string]string{"Authorization": token})
	})

	//	c.OnHTML("body", func(e *colly.HTMLElement) {
	//		// Extract data from HTML elements
	//		quote ,_:= e.DOM.Html()
	//
	//		fmt.Printf("Quote: %s", quote)
	//	})
	if err := c.Visit(url); err != nil {
		fmt.Println("Scraping failed. ", err)
		return err
	}

	return nil
}
