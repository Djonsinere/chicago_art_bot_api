package apicalls

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func Full_text_search(text string, chatID int64) string {

	path := fmt.Sprintf("https://api.artic.edu/api/v1/artworks/search?q=%s", text)
	resp, err := http.Get(path)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)

}

func Get_image(image_id string) { //достать из https://api.artic.edu/api/v1/artworks/656
	path := "https://www.artic.edu/iiif/2/" + image_id + "/full/843,/0/default.jpg"
	resp, err := http.Get(path)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("filename.jpg", []byte(body), 0755)
	if err != nil {
		fmt.Println("unable to write file: %w", err)
	}
}

//"data:image/gif;base64,R0lGODlhCAAFAPUAADY5NTQ/PkI4IT1COjtEQkNHQk5YWFteWWtkTmtoX291WnJ1WX1xZ355cYZ3Zn+Cf42KaoWEeJOKfpeWcoCEgI2OiJaQgp+XgpWRiJWTjKKbjqqhkK+olbGml7q0q7e2sMG5rcbAs8fBuMrGv8/HvsnJud3d3evr7gAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACH5BAAAAAAALAAAAAAIAAUAAAYlQJPnIwKRQpPTyPDgQDaC0gUQIAwsDISjcUgUMB2JpkLJRBSLIAA7"
