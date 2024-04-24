package main

import (
	"fmt"
	"net/url"
)

const myurl string = "https://www.google.com/webhp?hl=en&sa=X&ved=0ahUKEwiUrNeRzdmFAxVpvokEHSvDAcsQPAgJ"

func main() {

	content, _ := url.Parse(myurl)

	fmt.Println(content.Scheme)
	fmt.Println(content.Host)
	fmt.Println(content.Path)
	fmt.Println(content.Port())
	fmt.Println(content.RawQuery)

	qparam := content.Query()
	fmt.Println(qparam["h1"])

	for key, val := range qparam {
		fmt.Println("query data: ", key, val)
	}

	partsOfUrl := &url.URL{
		Scheme: "https",
		Host:   "www.google.com",
	}
	anotherUrl := partsOfUrl.String()
	fmt.Println(anotherUrl)
}
