package dnmtools

import (
	"encoding/base64"
	"github.com/nuveo/anticaptcha"
	"io/ioutil"
	"log"
	"net/http"
	)

var antiCaptchaKey = "11473f9181c240c403f68f1cbd47aed5"

func b64EncodeImgAtUrl(client *http.Client, url string) string {
	response, err := client.Get(url)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer response.Body.Close()

	captcha, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return ""
	}

	imgBase64Str := base64.StdEncoding.EncodeToString(captcha)

	log.Println(imgBase64Str)
	return imgBase64Str
}


func SolveCaptcha(client *http.Client, url string) (string, bool) {
	// base64 encode the image @ url
	encImg := b64EncodeImgAtUrl(client, url)
	if encImg == "" {
		return encImg, false
	}

	// submit encoded image to solve captcha
	antiClient := &anticaptcha.Client{APIKey: antiCaptchaKey}

	log.Println("submitting captcha")
	solution, err := antiClient.SendImage(encImg)
	if err != nil {
		log.Println(err)
		return solution, false
	}
	log.Println(solution)
	// process response
	return solution, true
}