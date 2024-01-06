package dl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

func GetLastFromSplit(s string, sep string) string {
	split := strings.Split(s, sep)
	if len(split) == 0 {
		return ""
	}
	return split[len(split)-1]
}

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

	fmt.Println("url: ", chunkUrl)
	// If it's breaking, 99% this is where the problem is
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

func getSVGSheetUrl(sheetId string, pageNumber int, headers map[string]string) (string, error) {
	urlPath := fmt.Sprintf(`https://musescore.com/api/jmuse?id=6102579&index=%d&type=img&v2=1`, pageNumber)
	fmt.Println("url: ", urlPath)

	// Creating new http request object
	request, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return "", err
	}

	// Adding all custom headers to request
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	// Creating a client to execute the request
	client := http.Client{}
	// Executing
	res, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// reading svg url data
	var data struct {
		Result string
		Status int
		Error  string
		Info   struct {
			Url string
		}
	}

	json.NewDecoder(res.Body).Decode(&data)

	fmt.Println(data.Status)

	if len(data.Info.Url) == 0 {
		return "", errors.New("Failed to get SVG url")
	}

	// fmt.Println("SVG url: ", data.Info.Url)
	return data.Info.Url, nil
}

func downloadFile(fileUrl string, fileName string) error {
	res, err := http.Get(fileUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		return errors.New("Request failed")
	}

	if len(fileName) == 0 {
		fileName = "DummyFilename"
	}
	err = os.MkdirAll("downloads", 0777)
	if err != nil {
		fmt.Println("Failed to create folder")
		return err
	}
	fileName = "downloads/" + fileName
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Failed to create file")
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func downloadSvgsTillFailure(url string, headers map[string]string) {
	var pageIndex int = 0
	for {
		svgUrl, err := getSVGSheetUrl(url, pageIndex, headers)
		if err != nil {
			panic(err)
		}
		fileName := fmt.Sprintf("%s_%d.svg", GetLastFromSplit(url, "/"), pageIndex)
		err = downloadFile(svgUrl, fileName)
		if err != nil {
			if pageIndex == 0 {
				panic(err)
			} else {
				break
			}
		}
		pageIndex++
	}
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
		downloadSvgsTillFailure(url, map[string]string{"authorization": token})
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
