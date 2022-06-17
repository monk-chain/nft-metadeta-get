package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)
type Metadeta struct {
	Name string `json:name`
	Image string `json:image`
	Attributes []Attributes
}
type Attributes struct {
	Trait_type string `json:trait_type`
	Value string `json:value`
}

var (
    result  = "result"
    name  = "Azuki"
    count  = 3
    contract_address = "0xed5af388653567af2f388e6224dc7c4b3241c544"
    url = "https://ikzttp.mypinata.cloud/ipfs/QmQFkLSQysj94s5GvTHPyzTxrawwtjgiiYS2TBLgrvw8CW/"
)

func main() {
	fmt.Println("getting..")

	metadeta_path := result  + "/" + name + "/metadeta/"
	image_path := result + "/" + name + "/image/"
	os.MkdirAll(metadeta_path, 0755)
	os.MkdirAll(image_path, 0755)

	var wg sync.WaitGroup
	wg.Add(1)
	go func () {
		for i := 0; i < count; i++ {
			// metadeta
			fmt.Println(i)
			var data Metadeta
			resp, _ := http.Get(url + "/" + strconv.Itoa(i))
			body, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			json.Unmarshal(body, &data)

			file, err := os.Create(metadeta_path + strconv.Itoa(i) + ".json")
			if err != nil {
				fmt.Println(err)
			}
			defer file.Close()
			json.NewEncoder(file).Encode(data)

			// image
			resp, _ = http.Get(data.Image)
			body , _  = ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()

			img, _, err := image.Decode(bytes.NewReader(body))
			if err != nil {
				continue
			}

			file, err = os.Create(image_path + strconv.Itoa(i) + ".png")
			if err != nil {
				fmt.Println(err)
			}
			defer file.Close()

			err = png.Encode(file, img)
			if err != nil {
				continue
			}
		}
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("done")
}