package main

import (
	"flag"
	"fmt"
	"github.com/mrjones/oauth"
	"log"
	"os"
)

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

// Request user's permission and return an access token.
func Authorize(client *FanfouClient) (*oauth.AccessToken, error) {
	rtoken, cbUrl, err := client.GetRequestTokenAndUrl("oob")
	if err != nil {
		return nil, err
	}

	var verifyCode string
	fmt.Printf("请点击下面的链接来为应用授权。\n%s\n", cbUrl)
	fmt.Printf("授权完成后请按Enter键继续。\n")
	fmt.Scanln(&verifyCode)

	atoken, err := client.AuthorizeToken(rtoken, verifyCode)
	if err != nil {
		return nil, err
	}

	return atoken, nil
}

func main() {

	var key, secret string
	var configFile, logFile string
	var DBFileName string

	// parse commandline options 
	flag.StringVar(&key, "key", "", "设置你的OAuth consumer key")
	flag.StringVar(&secret, "secret", "", "设置你的OAuth consumer secret")
	flag.StringVar(&configFile, "config", "config.json", "指定你的配置文件，非必须")
	flag.StringVar(&logFile, "log", "fanfou.log", "指定你的log文件，非必须")
	flag.StringVar(&DBFileName, "db", "fanfou.sqlite", "指定你的数据库文件")
	flag.Parse()

	// predefined these variables to make sure the scope
	var err error
	var fanfouClient *FanfouClient

	// If configuration file exist, just use it.
	// Or the user must manually gives a key and a secret.
	if isFileExist(configFile) {
		fanfou := &Fanfou{}
		fanfou.LoadConfig(configFile)

		fanfouClient, err = fanfou.NewFanfouClient()
		checkError(err)
	} else {
		// Check key and secret
		if key == "" || secret == "" {
			flag.Usage()
			return
		}
		fanfou := &Fanfou{
			Key:    key,
			Secret: secret,
		}

		fanfouClient, err = fanfou.NewFanfouClient()
		checkError(err)

		// Request user's permission
		atoken, err := Authorize(fanfouClient)
		checkError(err)

		fanfou.AccessToken = atoken
		fanfou.WriteConfig(configFile)
	}

	err = CreateFanfouDB(DBFileName)
	checkError(err)

	db, err := OpenDB(DBFileName)
	checkError(err)
	defer db.Close()

	// Set output log file
	file, err := os.Create(logFile)
	checkError(err)
	log.SetOutput(file)

	err = GetAndStoreAllStatuses(fanfouClient, db)
	checkError(err)
}
