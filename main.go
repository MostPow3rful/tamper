package main

import (
	"flag"
	"fmt"
	"github.com/corpix/uarand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Flag struct {
	IgnoreResponse string
	MatchResponse  string
	Output         string
	Target         string
	Cookie         string
	ExtraMethods   bool
}

func tamper(_httpMethod string, _userFlags Flag) {
	mc := !(_userFlags.MatchResponse == "")
	fc := !(_userFlags.IgnoreResponse == "")

	// Create an instance Of http.Client Struct
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	// Creating Custom Request
	request, err := http.NewRequest(_httpMethod, _userFlags.Target, nil)
	if err != nil {
		fmt.Printf(
			"[%s] Couldn't Make Request [%s]\n",
			_httpMethod, _userFlags.Target,
		)
		return
	}

	// Set HTTP Headers
	if _userFlags.Cookie != "" {
		request.Header.Set("Cookie", _userFlags.Cookie)
	}
	request.Header.Set("User-Agent", uarand.GetRandom())

	// Sending Request
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf(
			"[%s] Couldn't Send Request [%s]\n",
			_httpMethod, _userFlags.Target,
		)
		return
	}

	// Close Response Body [defer]
	defer response.Body.Close()

	if fc {
		if strings.Index(_userFlags.IgnoreResponse, ",") == -1 {
			code, _ := strconv.Atoi(_userFlags.IgnoreResponse)
			if response.StatusCode == code {
				return
			}
		}

		for _, v := range strings.Split(_userFlags.IgnoreResponse, ",") {
			code, _ := strconv.Atoi(v)
			if response.StatusCode == code {
				return
			}
		}
	}

	if mc {
		if strings.Index(_userFlags.MatchResponse, ",") == -1 {
			code, _ := strconv.Atoi(_userFlags.MatchResponse)
			if response.StatusCode != code {
				return
			}
		}

		for _, v := range strings.Split(_userFlags.MatchResponse, ",") {
			code, _ := strconv.Atoi(v)
			if response.StatusCode != code {
				return
			}
		}
	}

	if _userFlags.Output != "" {
		setResultInFile(_userFlags.Output, response.Status, _httpMethod, response.StatusCode)
		return
	}
	fmt.Printf("[%v] - [%d] - [%v]\n", _httpMethod, response.StatusCode, response.Status)
}

func createFile(_fileName string) {
	_, err := os.Create(_fileName)
	if err != nil {
		fmt.Printf("Couldn't Create Output File [%s]\n", _fileName)
		return
	}
}

func setResultInFile(_fileName, _httpMethod, _statusMessage string, _statusCode int) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Couldn't Get Output Of PWD Command")
		return
	}

	outputFile, err := os.OpenFile(fmt.Sprintf("%v/%v", pwd, _fileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0400)
	if err != nil {
		fmt.Printf("Couldn't Open Output File")
		return
	}
	outputFile.WriteString(fmt.Sprintf("[%v] - [%d] - [%v]\n", _httpMethod, _statusCode, _statusMessage))
}

func banner() {
	fmt.Println(" _____      _      __  __   ____    _____   ____")
	fmt.Println("|_   _|    / \\    |  \\/  | |  _ \\  | ____| |  _ \\")
	fmt.Println("  | |     / _ \\   | |\\/| | | |_) | |  _|   | |_) |")
	fmt.Println("  | |    / ___ \\  | |  | | |  __/  | |___  |  _ <")
	fmt.Println("  |_|   /_/   \\_\\ |_|  |_| |_|     |_____| |_| \\_\\")
	fmt.Println("[METHOD] - [STATUS CODE] - [STATUS MESSAGE]")
	fmt.Println()
}

func main() {
	// Creating Variables
	var (
		wg           = sync.WaitGroup{}
		flagInstance = Flag{
			IgnoreResponse: "",
			MatchResponse:  "",
			Output:         "",
			Target:         "",
			Cookie:         "",
			ExtraMethods:   false,
		}
		baseHttpMethods = [9]string{
			"CONNECT",
			"OPTIONS",
			"DELETE",
			"PATCH",
			"TRACE",
			"POST",
			"HEAD",
			"PUT",
			"GET",
		}
		extraHttpMethods = [23]string{
			"VERSION-CONTROL",
			"SHOWMETHOD",
			"TEXTSEARCH",
			"UNCHECKOUT",
			"ORDERPATCH",
			"PROPPATCH",
			"SPACEJUMP",
			"BAMBOOZL",
			"CHECKOUT",
			"PROPFIND",
			"NOEXISTE",
			"CHECKIN",
			"SEARCH",
			"UNLOCK",
			"UNLINK",
			"PURGE",
			"MKCOL",
			"REPOR",
			"TRACK",
			"INDEX",
			"LINK",
			"COPY",
			"LOCK",
			"MOVE",
		}
	)

	// Parse The User Flags
	flag.StringVar(&flagInstance.IgnoreResponse, "fc", "", "Don't Match Response Code [use ',' To Split]")
	flag.StringVar(&flagInstance.MatchResponse, "mc", "", "Match Response Code [use ',' To Split]")
	flag.StringVar(&flagInstance.Output, "o", "", "Name Of File To Set Result in it")
	flag.BoolVar(&flagInstance.ExtraMethods, "x", false, "FUZZ Extra HTTP Methods")
	flag.StringVar(&flagInstance.Target, "d", "", "URL Of Your Target Do You Want To Test")
	flag.StringVar(&flagInstance.Cookie, "c", "", "Set Value Of Cookie Header")
	flag.Parse()

	// Create Output File
	if flagInstance.Output != "" {
		createFile(flagInstance.Output)
	}

	// Check Target
	if flagInstance.Target == "" {
		fmt.Println("Invalid Target")
		os.Exit(0)
	}

	if !strings.HasPrefix(flagInstance.Target, "http://") && !strings.HasPrefix(flagInstance.Target, "https://") {
		fmt.Println("Invalid Prefix For Target")
		os.Exit(0)
	}

	banner()

	// Testing Base HTTP Methods
	for _, httpMethod := range baseHttpMethods {
		wg.Add(1)
		go func(_httpMethod string) {
			tamper(_httpMethod, flagInstance)
			wg.Done()
		}(httpMethod)
	}

	// Testing Extra HTTP Methods
	if flagInstance.ExtraMethods {
		for _, httpMethod := range extraHttpMethods {
			wg.Add(1)
			go func(_httpMethod string) {
				tamper(_httpMethod, flagInstance)
				wg.Done()
			}(httpMethod)
		}
	}

	wg.Wait()
}
