package silkroad4

import (
  "io/ioutil"
  "log"
  "net/http"
  "net/url"
  "strings"

  "code.8labs.io/rdemeritt/dnm-searcher/app"
  "github.com/PuerkitoBio/goquery"
)

var siteUrl = "http://silkroadkaxmspva.onion"
var authUrl = siteUrl + "/auth.php"
var loginUrl = siteUrl + "/?road"

var sr4Username = "akiyaenri"
var sr4Password = "LmfkUYgeiucUCCr7"


func Start(client *http.Client) bool {
  log.SetFlags(log.LstdFlags | log.Lshortfile)

	// initiate login loginSequence
	success := loginSequence(client)
	if !success {
		log.Println("login failed")
    return false
	}

  pages := getPages(client, siteUrl)
  if pages == nil {
		log.Println("get pages failed")
    return false
	}

  body, err := ioutil.ReadAll(pages.Body)
  if err != nil {
    log.Println(err)
  }
  defer pages.Body.Close()
  log.Println("getPages:", string(body))

  return true
}


func initializeSession(client *http.Client) (string, bool) {
  // silkroad4 requires us to go through an extra step before
  // we can login

  // populate our cookies
  response, err := client.Get(authUrl)
  if err != nil {
    log.Println("Error making GET request:", err)
  }

  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    log.Println(err)
  }
  defer response.Body.Close()
  log.Println(string(body))

  // lets make sure that the response makes it clear the session is initialized
  if response.StatusCode >= 200 && strings.Contains(string(body), "complete a captcha") {
    return string(body), true
  }
  return string(body), false
}


// we have to solve our first captcha before we can get the login page
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
          if strings.Contains(imgSrc, "captchas") {
            imgSrcUrl = imgSrc
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


func submitInitialCaptchaSolution(client *http.Client, solutionUrl string, solution string) bool {
  // query string {"Query string":{"a":"hrphc"}}
  initialCaptchaSolution := url.Values{
    "a": {solution},
  }

  // submit get w/ solution
  response, err := client.Get(solutionUrl + "?" + initialCaptchaSolution.Encode())
  if err != nil {
      log.Println(err)
      return false
  }
	defer response.Body.Close()

  if response.StatusCode == 200 {
    return true
  }
  log.Println("response status:", response.StatusCode)
  return false
}


func login(client *http.Client, username string, password string, captcha string) bool {
  // login parameters
  // {"Query string":{"road":""},"Form data":{"username":"akiyaenri","password":"LmfkUYgeiucUCCr7","attempt":"1","captcha":"aegadnzs","login":"Login"}}
  loginFormData := url.Values{
    "username": {username},
    "password": {password},
    "captcha": {captcha},
    "attempt": {"1"},
    "login": {"Login"},
  }

	log.Println("loginUrl:", loginUrl)
	log.Println("loginFormData:", loginFormData.Encode())
  response, err := client.PostForm(loginUrl, loginFormData)
  if err != nil {
    log.Println(err)
    return false
  }
	defer response.Body.Close()

  if response.StatusCode == 200 {
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
      log.Println(err)
    }
    log.Println("login successful")
    log.Println(string(body))
    for name, values := range response.Header {
      // Loop over all values for the name.
      for _, value := range values {
        log.Println(name, value)
        }
      }
    log.Println("list cookies")
		for _, cookie := range response.Cookies() {
	  	log.Println("response cookie:", cookie.Name, cookie.Value)
		}
    return true
  }
  return false
}


func loginSequence(client *http.Client) bool {
  // initilize our session
  body, success := initializeSession(client)
  if  !success {
    log.Println("unable to initilize session:", body)
    return false
  }

  // we have to solve our first captcha before we can get to our login page
  captchaUrl, success := getCaptchaImageUrl(client, authUrl)
  if  !success {
    log.Println("unable to get initial captcha:", captchaUrl)
    return false
  }

  // now that we have our image url let's submit the image for remote solution
  solution, success := dnmtools.SolveCaptcha(client, captchaUrl)
  if  !success {
    log.Println("remote captcha solution failed:", captchaUrl)
    return false
  }
  log.Println(solution, success)

  // we have our solution... let's submit it and move to logging in
  success = submitInitialCaptchaSolution(client, authUrl, solution)
  if !success {
    log.Println("unable to solve initial captcha")
    return false
  }

  // if we've made it this far, we can now get the login page, solve the
  // captcha, and login
  captchaUrl, success = getCaptchaImageUrl(client, loginUrl)
  if  !success {
    log.Println("unable to get login captcha:", captchaUrl)
    return false
  }

  // now that we have our image url let's submit the image for remote solution
  solution, success = dnmtools.SolveCaptcha(client, captchaUrl)
  if  !success {
    log.Println("remote captcha solution failed:", captchaUrl)
    return false
  }
  log.Println(solution, success)

  // success = login(client, sr4Username, sr4Password, solution)
	success = login(client, sr4Username, sr4Password, solution)
  if  !success {
    log.Println("login failed:", sr4Username, sr4Password, solution)
    return false
  }
  return true
}


func getPages(client *http.Client, pageUrl string) *http.Response {
	// fetch drops
	// {"Query string":{"road":"","cat":"97"}}

	var query99 = url.Values{
		"road": {""},
		"cat": {"99"},
	}

	log.Println("getting:", pageUrl + "/?" + query99.Encode())
	response, err := client.Get(pageUrl + "/?" + query99.Encode())
	if err != nil {
      log.Println(err)
      return nil
  }
	// defer response.Body.Close()

	return response
}
