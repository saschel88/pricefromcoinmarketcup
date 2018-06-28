package main

import (
	"strconv"
	"strings"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"flag"
)


// Уровень цены в BTC
type priceBTC struct {
	Price float64 `json:"price"`
}
// Уровень цены в USD
type priceUSD struct {
	Price float64 `json:"price"`
}
//  Уровень данных по USD
type quotesStruct struct {
	USD priceUSD `json:"USD"`
	BTC priceBTC `json:"BTC"`

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

// Верхний уровень структуры Test
type jsonStructTest struct {
	Data coinsIDNameSymbol `json:"data"`

}
//Структура для заполнения данных Symbol and Price
type coinsSymbolPrcie struct {
	Symbol string
	Price float64
}

var urlApiV2StartLimit ="https://api.coinmarketcap.com/v2/ticker/?start="
//https://api.coinmarketcap.com/v2/ticker/?start=101&limit=10
var startUrlApiV2ConvertToBTC="https://api.coinmarketcap.com/v2/ticker/"
var endStrToUrlApiV2ConvertToBTC ="/?convert=BTC"
const osargsGetID  = "ID"
// Флаг определяющий какой запрос отрабатывать
var flagID  = flag.Bool("ID",false,"Получение значения ID по значению Symbol")

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
func getPriceByID (currencyListInFile []string) ( string, error) {
	// Контейнер для распарсенных данных
	coinsDataTest:=jsonStructTest{}
	// Строка для формирования конечных данных
	stringIDSymbolPriceUSDAndBTC:=""

	currecyIDArray :=make([]int,len(currencyListInFile)) //Создаем слайс int для заполенения его нашими ID
	for k, v := range currencyListInFile {
		for i := 0; i < len(v); i++ {

			if v[i] == 32 || v[i] == 9 { // ID является значение от начала до первого пробела, либо TAB
				currecyIDArray[k],_ = strconv.Atoi(v[:i]) //формируем массив строк ID
				break

			}
		}
	}
	//Формирование URL запросов на получение данных ID Symbol price USD и BTC
	urlApiV2ConvertToBTC:=make([]string,len(currecyIDArray))
	for k,v :=range currecyIDArray {
		urlApiV2ConvertToBTC[k]=startUrlApiV2ConvertToBTC+strconv.Itoa(v)+endStrToUrlApiV2ConvertToBTC
	}
	//Получение  данных  Symbol price USD и BTC по ID
	for _,v:=range urlApiV2ConvertToBTC{

		resp, err := http.Get(v) //получаем данные с coinmarcetcap c помощью Get запроса по URL = urlApiV2ConvertToBTC
		if err != nil {
			log.Fatal("Получить данные по GET запросу не удалось: ", err)
		}
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body) // записываем полученную с запроса информацию в []byte  respBody
		if err != nil {
			log.Fatal("Считать данные с тела resp в []byte  не удалось: ", err)
		}
		err = json.Unmarshal(respBody, &coinsDataTest) // Распарсиваем данные с respBody в массив структур coinsDataTest
		if err != nil {
			log.Fatal("Not UnMarshaling:", err)
		}
		// Формирование данных для записи в выходной файл
		stringIDSymbolPriceUSDAndBTC+=strconv.Itoa(coinsDataTest.Data.ID)+"\t"+coinsDataTest.Data.Symbol+"\t"+strconv.FormatFloat(coinsDataTest.Data.Qoutes.BTC.Price, 'f', -1, 64)+"\t"+strconv.FormatFloat(coinsDataTest.Data.Qoutes.USD.Price, 'f', -1, 64)+"\n" //fmt 'f' означает, что представление без экспоненты

	}


return strings.Replace(stringIDSymbolPriceUSDAndBTC,".",",",-1),nil

}

func main() {
	tmpURLApiV2StartLimit:=""
	tmpStart:=0
	tmpLimit:="100"
	coinsData  := jsonStruct{}

	bs, err := ioutil.ReadFile("currencyold.txt") //Считываем с файла перечень валлют, которыми торгуем

	if err != nil {
		log.Fatal("НЕ удалось считать данные из файла ",err)
	}


	currencyListInFile := strings.Split(string(bs), "\n")  // получаем массив наименований валют из строки string(bs) по раделителю "\n"


	flag.Parse()
	// Получени ID валюты по SYMBOL
	if *flagID {
	//if true{
		for i:=1;i<10 ;i++  {
			//формируем URL
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

			currencyStringToFileNew, _ := getPriceByID(currencyListInFile)
			// Запись данных в файл currencynew.txt
			ioutil.WriteFile("currencynew.txt", []byte(currencyStringToFileNew), 0777)
		}

}
