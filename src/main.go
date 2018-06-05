package main

import (
	"encoding/json"


	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"os"

)

type coinSymbolAndPrice struct {
	Id        string `json:"id"`
	Symbol    string `json: "symbol"`
	Price_usd string `json:"price_usd"`
	//Price_usd float64 `json:"price_usd,string"` //парсим строковое значение в float64
}

var urlApi = "https://api.coinmarketcap.com/v1/ticker/?limit=1000"

//Данная функция применяется, когда нам известен только Symbol валюты
func dublicateCurrencyList(coinFromUnmarshalingFunc []coinSymbolAndPrice, lenArgs int) (map[string]string, string) { //Функция для переноса распарсенных данных в карту, а так же для  поиска повторений по параметру Symbol
	strListRepeatCurrency := ""
	currencyMap := make(map[string]string)
	if lenArgs > 1 {
		for i := 0; i < len(coinFromUnmarshalingFunc); i++ { //заполняем карту значениями из массива структур coinFromUnmarshaling
			_, ok := currencyMap[coinFromUnmarshalingFunc[i].Symbol]

			if ok {
				strListRepeatCurrency += coinFromUnmarshalingFunc[i].Symbol + ","
			}
			currencyMap[coinFromUnmarshalingFunc[i].Symbol] = coinFromUnmarshalingFunc[i].Id + "\t" + coinFromUnmarshalingFunc[i].Price_usd
		}

	} else {
		for i := 0; i < len(coinFromUnmarshalingFunc); i++ { //заполняем карту значениями из массива структур coinFromUnmarshaling
			currencyMap[coinFromUnmarshalingFunc[i].Id] = coinFromUnmarshalingFunc[i].Symbol + "\t" + coinFromUnmarshalingFunc[i].Price_usd
		}
	}

	return currencyMap, strListRepeatCurrency
}

func main() {
	lenArgs := len(os.Args)


	coinFromUnmarshaling := []coinSymbolAndPrice{} // объявляем массив структур coinSymbolAndPrice

	resp, err := http.Get(urlApi) //получаем данные с coinmarcetcap c помощью Get запроса по URL = urlApi
	if err != nil {               //проверка на ошибку
		log.Fatal("Возникала ошибка: ", err)
	}
	defer resp.Body.Close()                  //По окончанию работы закрываем resp.Body
	respBody, _ := ioutil.ReadAll(resp.Body) // записываем полученную с запроса информацию в []byte  respBody

	err = json.Unmarshal(respBody, &coinFromUnmarshaling) // Распарсиваем данные с respBody в массив структур coinFromUnmarshaling
	if err != nil {                                       //Проверяем на ошибки
		log.Fatal("Not UnMarshaling:", err)
	}


	bs, err := ioutil.ReadFile("currencyold.txt") //Считываем с файла перечень валлют, которыми торгуем
	if err != nil {
		log.Fatal(err)
	}
	currencyListInFile := strings.Split(string(bs), "\n")                   // получаем массив наименований валют из строки string(bs) по раделителю "\n"

	mapCoinSymbolAndPrice, strListRepeatCurrency := dublicateCurrencyList(coinFromUnmarshaling, lenArgs) //вызов функции
	currencyListInFile = currencyListInFile[:len(currencyListInFile)-1]
	stringCurrencyAndPrice := "" //пустая строка для формирования списка валют с ценами

	if lenArgs > 1 { //формируем список валюты по ID либо по SYMBOL в зависимости от lenArgs (количество аргументов при вызове программы) если 0 (по умолчанию)то режим работы с ID, else то режим работы с SYMBOL
		strListRepeatCurrency = strListRepeatCurrency[:len(strListRepeatCurrency)-1]
		arrayStrListRenameCurrency := strings.Split(strListRepeatCurrency, ",") //получаем массив повторяющихся валют из строки strListRenameCurrency по разделителю ","
		strListRepeatCurrency=""
		for k, v := range currencyListInFile {

			for i := 0; i < len(v); i++ {

				if v[i] == 32 || v[i] == 9 {
					currencyListInFile[k] = v[:i] //формируем массив строк SYMBOL
					i = 100

				}
			}
		}

		for _, v := range currencyListInFile { //проверяем на наличие в нашем списке валют (symbol), повторяющихся из arrayStrListRenameCurrency
			for j := 0; j < len(arrayStrListRenameCurrency)-1; j++ {
				if v == arrayStrListRenameCurrency[j] {
					strListRepeatCurrency += v + ", "
				}
			}
		}
		strListRepeatCurrency = strListRepeatCurrency[:len(strListRepeatCurrency)-2]

	} else {
		for k, v := range currencyListInFile {

			for i := 0; i < len(v); i++ {

				if v[i] == 32 || v[i] == 9 {
					currencyListInFile[k] = v[i+1:] //формируем массив строк ID
					i = 100

				}
			}
		}

	}
	textMessageAboutRepeatCurrency := "Требуется ручная проверка в связи с тем, что данная валюта имеет повторения \n"
	for _, v := range currencyListInFile { //формирования списка валют с ценами
		stringCurrencyAndPrice += v + "\t" + mapCoinSymbolAndPrice[v] + "\n"

	}
	stringCurrencyAndPriceNew := strings.Replace(stringCurrencyAndPrice, ".", ",", -1) // замена точек на запятые
	stringCurrencyAndPriceNew += textMessageAboutRepeatCurrency+strListRepeatCurrency
	ioutil.WriteFile("currencynew.txt", []byte(stringCurrencyAndPriceNew), 0777)

}
