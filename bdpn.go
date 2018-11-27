package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	ConnectionError = "BDPN недоступен."
	DataError       = "Ошибка чтения данных."
	XMLError        = "Ответ BDPN не XML."
	OperatorError   = "Оператор определен не верно."
	CheckOK         = "OK"
	Unknown         = "Статус проверки неизвестен."
	crit            = 2
	warning         = 1
	normal          = 0
)

type BDPN struct {
	IsMegaFon bool   `xml:"Is_MegaFon"`
	Operator  string `xml:"Operator>Name"`
}

func StopCheck(m string, e error, c int) {
	fmt.Printf("%s\n", m)
	if e != nil {
		fmt.Printf("%s\n", e)
	}
	os.Exit(c)
}

func main() {
	host := flag.String("host", "10.236.26.169", "a host")
	msisdn := flag.String("msisdn", "79262631598", "a msisdn")
	operator := flag.String("operator", "MegaFon", "a operator")
	flag.Parse()
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	s := fmt.Sprintf("http://%s:80/mnp/msisdn/%s/xml", *host, *msisdn)
	u, err := url.Parse(s)
	if err != nil {
		StopCheck(DataError, err, crit)
	}
	resp, err := client.Get(u.String())
	if err != nil {
		StopCheck(ConnectionError, err, crit)
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		StopCheck(DataError, err, crit)
	}

	var bdpn BDPN
	err = xml.Unmarshal(contents, &bdpn)
	if err != nil {
		StopCheck(XMLError, err, crit)
	}

	if bdpn.Operator != *operator || !bdpn.IsMegaFon {
		fmt.Printf("%s Оператор: %s .IsMegaFon: %t\n", OperatorError, bdpn.Operator, bdpn.IsMegaFon)
	}
	StopCheck(CheckOK, nil, normal)
}
