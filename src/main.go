package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

// Уровень цены в USD
type priceUSD struct {
	Price float64 `json:"price"`
}
//  Уровень данных по USD
type quotesStruct struct {
	USD priceUSD `json:"USD"`
}
// Уровень данных Валюты
type coinsIDNameSymbol struct {

	ID int `json:"id"`
	Name string `json:"name"`
	Symbol string `json:"symbol"`
	Qoutes quotesStruct `json:"quotes"`
}
// Верхний уровень структуры
type jsonStruct struct {
	Data map[string]coinsIDNameSymbol `json:"data"`

}

var urlApiv2 = "https://api.coinmarketcap.com/v2/ticker/?limit=2"


//Функция получения ID по значению Symbol валюты

//func getIDbySymbol (sliceOfCoinData []jsonStruct) (map[string]string, string) {
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


	err = json.Unmarshal(respBody, &coinsData) // Распарсиваем данные с respBody в массив структур sliceOfCoinData
	if err != nil {                                       //Проверяем на ошибки
		log.Fatal("Not UnMarshaling:", err)
	}

	logs.Info(coinsData)
}
