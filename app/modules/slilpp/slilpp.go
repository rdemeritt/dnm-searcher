package slilpp

import (
	dnmtools "code.8labs.io/rdemeritt/dnm-searcher/app"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var siteUrl = "https://slilpp.top"
var loginUrl = siteUrl + "/" + "log.php"
var shopsUrl = siteUrl + "/" + "shops"

var siteUsername = "obayoshinori"
var sitePassword = "w8kdEfTVN7gyPTjX"

func Start(client *http.Client) bool {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// initiate login
	success := loginSequence(client)
	if !success {
		log.Println("login failed")
		return false
	}

	pages := getPages(client, shopsUrl)
	log.Println(pages)

	return true
}


func loginSequence(client *http.Client) bool {
	// get first captcha image for login
	captchaUrl, success := getCaptchaImageUrl(client, loginUrl)
	if  !success {
		log.Println("unable to get login captcha:", captchaUrl)
		return false
	}
	log.Println(captchaUrl, success)

	// now that we have our image url let's submit the image for remote solution
	solution, success := dnmtools.SolveCaptcha(client, captchaUrl)
	if  !success {
		log.Println("remote captcha solution failed:", captchaUrl)
		return false
	}
	log.Println(solution, success)

	// let's login
	success = login(client, siteUsername, sitePassword, solution)

	return true
}


func login(client *http.Client, username string, password string, captcha string) bool {
// login parameters
//	submitted=1&username=obayoshinori&password=w8kdEfTVN7gyPTjX&keystring=544271&login=
	loginFormData := url.Values{
		"username": {username},
		"password": {password},
		"keystring": {captcha},
		"submitted": {"1"},
		"login": {""},
	}

	log.Println("loginUrl:", loginUrl)
	log.Println("loginFormData:", loginFormData.Encode())

	response, err := client.PostForm(loginUrl, loginFormData)

	if err != nil {
		log.Println(err)
		return false
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(response.StatusCode, string(body))

	if response.StatusCode == 200 {
		log.Println("login successful")
		return true
	}

	return false
}


// we have to solve our captcha before we can submit the login page
func getCaptchaImageUrl(client *http.Client, url string) (string, bool) {
	// now that our session has been initialized, lets get past this first captcha
	response, err := client.Get(url)
	if err != nil {
		log.Println("Error making GET request:", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	// for some reason it does not like it when i give it
	// response.Body...  i have to create a new io.Reader
	// from the string of the html...  not sure why
	htmlReader := strings.NewReader(string(body))

	// Create a goquery document from the HTTP response
	// document, err := goquery.NewDocumentFromReader(response.Body)
	document, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		log.Println("Error loading HTTP response body. ", err)
	}

	// Find and print image URLs
	var imgSrcUrl string
	document.Find("img").EachWithBreak(func(index int, element *goquery.Selection) bool {
		imgSrc, exists := element.Attr("src")
		if exists {
			log.Println(imgSrc)
			if strings.Contains(imgSrc, "kcaptcha") {
				imgSrcUrl = imgSrc[2:]
				return false
			}
		}
		return true
	})

	if imgSrcUrl == "" {
		return imgSrcUrl, false
	}
	imgSrcUrl = siteUrl + "/" + imgSrcUrl
	return imgSrcUrl, true
}


func getPages(client *http.Client, pageUrl string) string {
	// fetch "otherShoPPs"

	var otherShoPPs = url.Values{
		"submitted": {"1"},
	}

	log.Println("getting:", pageUrl + "/search.php?" + otherShoPPs.Encode())
	response, err := client.Get(pageUrl + "/search.php?" + otherShoPPs.Encode())
	if err != nil {
		log.Println(err)
		return ""
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	return string(body)
}