package main

import (
	"bufio"
	"encoding/csv"
	"net/http"
	"net/url"
	"encoding/json"
	"bytes"
	"strings"
	"os"
	"fmt"
	"io"
	"io/ioutil"
	"time"
)

type Token struct {
	AccessToken string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
	TokenType string `json:"token_type" bson:"token_type"`
	ExpiresIn int `json:"expires_in" bson:"expires_in"`
	Scope string `json:"scope" bson:"scope"`
	UserId int `json:"user_id" bson:"user_id"`
}

type Question struct {
	Text string `json:"text" bson:"text"`
	ItemId string `json:"item_id" bson:"item_id"`
}

func main() {
	// Load a CSV file.
	fmt.Println("Lendo arquivo CSV...")
	f, err := os.Open("leads.csv")
	if err != nil {
	 	fmt.Println("Não encontrei o arquivo leads.csv")
	 	fmt.Println("O dois arquivos estão na mesma pasta?")
	 	return
	 }

	fmt.Println("Autenticando no MercadoLivre...")
	resp, err := http.PostForm("https://api.mercadolibre.com/oauth/token", url.Values{"grant_type": {"client_credentials"}, "client_id": {"APP_ID"}, "client_secret": {"APP_SECRET"}})
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	var t Token   
	err = decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(f))
	r.Comma = ';'
	
	for {

		record, err := r.Read()

		// Stop at EOF.
		if err == io.EOF {
			break
		}

		data := strings.Split(record[3], " - ")

		question := Question{ 
			"Nome: " + record[12] + " | Telefone: " + record[13] + " | Email: " + record[11],
			data[0],
		}
		
		ret := makeQuestion(data[0], t.AccessToken, question)

		fmt.Println(ret)

		fmt.Println("Aguardando 2 minutos para a proxima pergunta...")
		time.Sleep(2 * time.Minute)
	}
}

func makeQuestion(item string, token string, question Question) string{

	url := "https://api.mercadolibre.com/questions/" + item + "?access_token=" + token

	out, err := json.Marshal(question)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(out))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	retorno := item + " => " + resp.Status
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(ioutil.NopCloser(bytes.NewBuffer(body)))

	return retorno
}