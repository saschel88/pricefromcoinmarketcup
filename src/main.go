package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"github.com/astaxie/beego/logs"
)



type coinsData struct {
	//	Price_usd string `json:"price_usd"`
	ID int `json:"id"`
	Name string `json:"name"`
	Symbol string `json:"symbol"`

}

type jsonStruct struct {
	Data map[string]coinsData `json:"data"`

}

const countCoins  = 50
var urlApiv2 = "https://api.coinmarketcap.com/v2/ticker/?limit=1"
var  urlApiv2Listing ="https://api.coinmarketcap.com/v2/listings/"

//Функция получения ID по значению Symbol валюты

//func getIDbySymbol (sliceOfCoinData []coinsData) (map[string]string, string) {
//
//}

func main() {

	coinsData  := jsonStruct{}


	resp, err := http.Get(urlApiv2) //получаем данные с coinmarcetcap c помощью Get запроса по URL = urlApiv2
	if err != nil {               //проверка на ошибку
		log.Fatal("Получить данные по GET запросу не удалось: ", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body) // записываем полученную с запроса информацию в []byte  respBody
	if err != nil {               //проверка на ошибку
		log.Fatal("Считать данные с тела resp в []byte  не удалось: ", err)
	}
	logs.Warning(string(respBody))
	logs.Info(coinsData)
	err = json.Unmarshal(respBody, &coinsData) // Распарсиваем данные с respBody в массив структур sliceOfCoinData
	if err != nil {                                       //Проверяем на ошибки
		log.Fatal("Not UnMarshaling:", err)
	}

	logs.Info(coinsData)
}
