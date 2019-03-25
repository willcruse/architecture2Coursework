package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "regexp"
    "strings"
)

const ipReg = "^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$"

type address struct {
    Ip            string  `json:ip`
    City          string  `json:city`
    Region        string  `json:region`
    RegionCode    string  `json:region_code`
    Country       string  `json:country`
    CountryName   string  `json:country_name`
    ContinentCode string  `json:continent_code`
    Postal        string  `json:postal`
    Latitude      float64 `json:latitude`
    Longitude     float64 `json:longituide`
    Timezone      string  `json:timezone`
    Org           string  `json:org`
}

func main() {
    var ma map[string]address
    ma = make(map[string]address)
    dataB, err := ioutil.ReadFile("routes.txt")
    if err != nil {
        fmt.Println(err)
    }
    data := strings.Split(string(dataB), " ")
    for _, k := range data {
        k = strings.ReplaceAll(k, "(", "")
        k = strings.ReplaceAll(k, ")", "")
        com, err := regexp.Compile(ipReg)
        if err != nil {
            fmt.Println(err)
        }
        match := com.Match([]byte(k))
        if err != nil {
            //fmt.Println(err)
        }

        if match && len(k) >= 8 && len(k) <= 15 {
            _, ok := ma[k]
            if !ok {
                ma[k] = address{Ip: k}
            }
        }
    }
    for k, addres := range ma {
        url := fmt.Sprintf("https://ipapi.co/%s/json/", addres.Ip)
        req, err := http.Get(url)
        if err != nil {
            fmt.Println(err)
        }
        decoder := json.NewDecoder(req.Body)
        err = decoder.Decode(&addres)
        ma[k] = addres
        if err != nil {
            fmt.Println(err)
        }
    }
    var toWrite []byte
    for _, a := range ma {
        b, e := json.Marshal(a)
        if e == nil {
            for _, l := range b {
                toWrite = append(toWrite, l)
            }
            toWrite = append(toWrite, '\n')
        } else {
            fmt.Println(e)
        }
    }
    fmt.Println(len(ma))
    fmt.Println(string(toWrite))
    ioutil.WriteFile("res.txt", toWrite, 0644)

}

