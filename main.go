package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Fanfou status structure
type Status struct {
	CreatedAt string `json:"created_at"`
	Id        string `json:"id"`
	RawId     int64  `json:"rawid"`
	Text      string `json:"text"`
	Source    string `json:"source"`
}

// Check error information
func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// Check if file exist
func isFileExist(name string) bool {
	_, err := os.Open(name)

	// If err is not nil, file must exist
	if err == nil {
		return true
	}

	return os.IsExist(err)
}

func main() {

	var key, secret string
	var configFile, logFile string
	var DBFileName string
	flag.StringVar(&key, "key", "", "设置你的OAuth consumer key")
	flag.StringVar(&secret, "secret", "", "设置你的OAuth consumer secret")
	flag.StringVar(&configFile, "config", "config.json", "指定你的配置文件，非必须")
	flag.StringVar(&logFile, "log", "fanfou.log", "指定你的log文件，非必须")
	flag.StringVar(&DBFileName, "db", "fanfou.sqlite", "指定你的数据库文件")
	flag.Parse()

	var fanfou *Fanfou
	if isFileExist(configFile) {
		fanfou = &Fanfou{}
		fanfou.LoadConfig(configFile)
	} else {
		if key == "" || secret == "" {
			flag.Usage()
			return
		}
		fanfou = &Fanfou{
			Key:    key,
			Secret: secret,
		}

		fanfouClient, err := fanfou.NewFanfouClient()
		if err != nil {
			log.Fatalln(err)
		}

		rtoken, cbUrl, err := fanfouClient.GetRequestTokenAndUrl("oob")
		if err != nil {
			log.Fatalln(err)
		}

		var verifyCode string
		fmt.Printf("请点击下面的链接来为应用授权。\n%s\n", cbUrl)
		fmt.Printf("授权完成后请按Enter键继续。\n")
		fmt.Scanln(&verifyCode)

		atoken, err := fanfouClient.AuthorizeToken(rtoken, verifyCode)
		if err != nil {
			log.Fatalln(err)
		}
		fanfou.AccessToken = atoken
		fanfou.WriteConfig(configFile)
	}

	fanfouClient, err := fanfou.NewFanfouClient()
	checkError(err)

	err = CreateFanfouDB(DBFileName)
	checkError(err)

	db, err := OpenDB(DBFileName)
	checkError(err)
	defer db.Close()

	// Set output log file
	file, err := os.Create(logFile)
	checkError(err)
	log.SetOutput(file)

	var sinceId string
	for {
		resp, err := fanfouClient.Get("http://api.fanfou.com/statuses/user_timeline.json", map[string]string{
			"mode": "lite", "count": "60"}, fanfouClient.Token)
		checkError(err)
		defer resp.Body.Close()

		bytes, err := ioutil.ReadAll(resp.Body)
		checkError(err)

		// Unmarshal all json data
		statuses := []Status{}
		json.Unmarshal(bytes, &statuses)
		// If we already read all statuses, then quit
		if len(statuses) == 0 {
			break
		}

		// Insert all statuses into table
		err = InsertStatuses(db, &statuses)
		checkError(err)

		// Get last status's id
		sinceId = statuses[len(statuses)-1].Id
		log.Printf("id: %s OK\n", sinceId)
	}
}
