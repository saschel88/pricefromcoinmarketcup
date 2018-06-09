package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"strings"
	"strconv"
	"os"
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
//Структура для заполнения данных Symbol and Price
type coinsSymbolPrcie struct {
	Symbol string
	Price float64
}

var urlApiV2 = "https://api.coinmarketcap.com/v2/ticker/?limit=200"
var urlApiV2StartLimit ="https://api.coinmarketcap.com/v2/ticker/?start="
//https://api.coinmarketcap.com/v2/ticker/?start=101&limit=10

const osargsGetID  = "ID"



//Функция получения ID по значению Symbol валютыё
func getIDbySymbol (coinsData jsonStruct) (map[string]int, string, error) {
	currencyMap := make(map[string]int)
	strListRepeatCurrency:=""

	for _,v :=range coinsData.Data{
		// Формирование списка валют с одинаковыми Symbol
		_,ok := currencyMap[v.Symbol]
		if ok {
			strListRepeatCurrency += v.Symbol + ","
		}
		currencyMap[v.Symbol]=v.ID // Заполнение карты currencyMap:Symbol=ID
	}
	if len(strListRepeatCurrency)>1 {
		strListRepeatCurrency=strListRepeatCurrency[:len(strListRepeatCurrency)-1] //Если есть хоть один элемент в строке, то удалем последний (",")
	}
return currencyMap, strListRepeatCurrency, nil // Необходимо отработать условия, при которых возвраащется err

}
//Функция получения Price по значению ID валютыё
func getPriceByID (coinsData jsonStruct, currencyListInFile []string) ( string, error) {
	stringIDSymbolAndPrice:=""
	currecyIDArray :=make([]int,len(currencyListInFile)) //Создаем слайс int для заполенения его нашими ID
	for k, v := range currencyListInFile {
		for i := 0; i < len(v); i++ {

			if v[i] == 32 || v[i] == 9 { // ID является значение от начала до первого пробела, либо TAB
				currecyIDArray[k],_ = strconv.Atoi(v[:i]) //формируем массив строк ID
				break

			}
		}
	}

// формируем строковый массив с данными по ID, Symbol и Price из coinsData
	for k, arrayV := range currecyIDArray {
		for _,v:=range coinsData.Data{
			if v.ID==arrayV {
				stringIDSymbolAndPrice+=currencyListInFile[k]+"\t"+strconv.FormatFloat(v.Qoutes.USD.Price, 'f', -1, 64)+"\n" //fmt 'f' означает, что представление без экспоненты

			}
		}
	}
return stringIDSymbolAndPrice,nil

}

func main() {
	tmpURLApiV2StartLimit:=""
	tmpStart:=0
	tmpLimit:="100"
	coinsData  := jsonStruct{}

	for i:=1;i<4 ;i++  {

		tmpURLApiV2StartLimit=urlApiV2StartLimit+strconv.Itoa(tmpStart+1)+"&limit="+tmpLimit

		resp, err := http.Get(tmpURLApiV2StartLimit) //получаем данные с coinmarcetcap c помощью Get запроса по URL = urlApiv2
	if err != nil {               //проверка на ошибку
		log.Fatal("Получить данные по GET запросу не удалось: ", err)
	}
		tmpStart+=100
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body) // записываем полученную с запроса информацию в []byte  respBody
	if err != nil {               //проверка на ошибку
		log.Fatal("Считать данные с тела resp в []byte  не удалось: ", err)
	}
	err = json.Unmarshal(respBody, &coinsData) // Распарсиваем данные с respBody в массив структур sliceOfCoinData
	if err != nil {                                       //Проверяем на ошибки
		log.Fatal("Not UnMarshaling:", err)
	}


		tmpURLApiV2StartLimit=""
	}

	bs, err := ioutil.ReadFile("currencyold.txt") //Считываем с файла перечень валлют, которыми торгуем

	if err != nil {
		log.Fatal("НЕ удалось считать данные из файла ",err)
	}


	currencyListInFile := strings.Split(string(bs), "\n")  // получаем массив наименований валют из строки string(bs) по раделителю "\n"

	// Получени ID валюты по SYMBOL
	if len(os.Args)>1 && os.Args[1]==osargsGetID {
	//if true{
		currencyMapBySymbol, strListRepeatCurrency, _ := getIDbySymbol(coinsData) //Получение  карты currencyMap:Symbol=ID и списка валют с одинаковыми Symbol

		stringIDAndSymbol := ""                                                   //пустая строка для формирования списка валют с ID
		for _, v := range currencyListInFile {
			_, ok := currencyMapBySymbol[v]
			if ok {
				stringIDAndSymbol += strconv.Itoa(currencyMapBySymbol[v]) + "\t" + v + "\n"
			}
		}


		textMessageAboutRepeatCurrency := "\n Требуется ручная проверка в связи с тем, что данная валюта имеет повторения \n"
		// Осуществляем проверку наличия валют из файла в списке повторяющихся валют

		if len(strListRepeatCurrency) > 0 {
			arrayStrListRepeatCurrency := strings.Split(strListRepeatCurrency, ",") //получаем массив повторяющихся валют из строки strListRenameCurrency по разделителю ","

			strListRepeatCurrency = ""

			for _, v := range currencyListInFile {
				for _, val := range arrayStrListRepeatCurrency {
					if v == val {
						strListRepeatCurrency += v + ","
					}
				}
			}

			if len(strListRepeatCurrency) > 1 {
				strListRepeatCurrency = strListRepeatCurrency[:len(strListRepeatCurrency)-1] //Если есть хоть один элемент в строке, то удалем последний (",")
				stringIDAndSymbol += textMessageAboutRepeatCurrency + strListRepeatCurrency
				}

		}

		ioutil.WriteFile("currencynew.txt", []byte(stringIDAndSymbol), 0777)

	}else {

			//Формирование данных для получения ID Symbol Price

			currencyStringToFileNew, _ := getPriceByID(coinsData, currencyListInFile)

			ioutil.WriteFile("currencynew.txt", []byte(currencyStringToFileNew), 0777)
		}
}
