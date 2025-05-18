package apicalls

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type APIResponse struct {
	Data []Artwork `json:"data"`
}

type Artwork struct {
	ID      int    `json:"id"`
	APILink string `json:"api_link"`
	Title   string `json:"title"`
}

type sec_APIResponse struct {
	Data ImageData `json:"data"`
}

type ImageData struct {
	ID          int    `json:"id"`
	ImageID     string `json:"image_id"`
	CreditLine  string `json:"credit_line"`
	ArtistTitle string `json:"artist_title"`
	Dimensions  string `json:"dimensions"` //размеры картины в реальности
	//publication_history тут длинная история публикации, ее по идее надо бы в отдельную кнопку вынести
	Сlassification_title string `json:"classification_title"` //классификация, к примеру скульптура
	Date_display         string `json:"date_display"`
}

func Full_text_search(text string, chatID int64) [50]ImageData {
	path := fmt.Sprintf("https://api.artic.edu/api/v1/artworks/search?q=%s", text)
	resp, err := http.Get(path)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var apiResponse APIResponse
	if err := json.Unmarshal(result, &apiResponse); err != nil {
		log.Fatal("JSON parse error:", err)
	}

	var data_array [50]ImageData

	for count, artwork := range apiResponse.Data {
		count += 1
		path := artwork.APILink
		resp, err := http.Get(path)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		full_responce, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var api_response sec_APIResponse
		if err := json.Unmarshal(full_responce, &api_response); err != nil {
			log.Fatal("JSON parse error:", err)
		}

		data_array[count] = api_response.Data

		Get_image(api_response.Data.ImageID, chatID)

		if count == 50 {
			break
		}
	}

	return data_array
}

func Get_image(image_id_api string, user_id int64) {
	path := "https://www.artic.edu/iiif/2/" + image_id_api + "/full/843,/0/default.jpg"
	resp, err := http.Get(path)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	str_user_id := strconv.FormatInt(user_id, 10)

	dirPath := filepath.Join(str_user_id)
	filePath := filepath.Join(dirPath, image_id_api+".jpg") //TODO сохранять в /tmp
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		fmt.Println("ошибка создания директорий: ", err)
	}
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("ошибка создания файла: ", err)
	}
	defer file.Close()

	if _, err := file.Write([]byte(body)); err != nil {
		fmt.Println("ошибка записи в файл: ", err)
	}

}
