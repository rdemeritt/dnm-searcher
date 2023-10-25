package sample_newreader

import (
  "log"
  "strings"
  "github.com/PuerkitoBio/goquery"
)

var htmlDocument string = `<center>
	<font color='red'><font size=5>CAPTCHA</font></font><br><br>
	<center>Before accessing the Secret Page, please solve this captcha</center><br>
	<img src="captchas/23/kNOJBjzGAFAVAGAE.jpg"><br>
	<form method="GET">
		<input type="text" name="a" placeholder="answer">
		<input type="submit" value="Submit">
	</form>
</center>`

func main()  {
  htmlReader := strings.NewReader(htmlDocument)
  document, _ := goquery.NewDocumentFromReader(htmlReader)

  // Find and print image URLs
  document.Find("img").Each(func(index int, element *goquery.Selection) {
      imgSrc, exists := element.Attr("src")
      if exists {
          log.Println(imgSrc)
      }
  })
}
