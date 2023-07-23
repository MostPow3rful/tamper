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

type Flags struct {
	forgetResponse string
	matchResponse  string
	output         string
	target         string
	cookie         string
	extraMethods   bool
}

func tamper(_httpMethod string, _userFlags Flags) {
	mc := !(_userFlags.matchResponse == "")
	fc := !(_userFlags.forgetResponse == "")

	// Create an instance Of http.Client Struct
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	// Creating Custom Request
	request, err := http.NewRequest(_httpMethod, _userFlags.target, nil)
	if err != nil {
		fmt.Printf(
			"[%s] Couldn't Make Request [%s]\n",
			_httpMethod, _userFlags.target,
		)
		return
	}

	// Set HTTP Headers
	if _userFlags.cookie != "" {
		request.Header.Set("Cookie", _userFlags.cookie)
	}
	request.Header.Set("User-Agent", uarand.GetRandom())

	// Sending Request
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf(
			"[%s] Couldn't Send Request [%s]\n",
			_httpMethod, _userFlags.target,
		)
		return
	}

	// Close Response Body [defer]
	defer response.Body.Close()

	if fc {
		if strings.Index(_userFlags.forgetResponse, ",") == -1 {
			code, _ := strconv.Atoi(_userFlags.forgetResponse)
			if response.StatusCode == code {
				return
			}
		}

		for _, v := range strings.Split(_userFlags.forgetResponse, ",") {
			code, _ := strconv.Atoi(v)
			if response.StatusCode == code {
				return
			}
		}
	}

	if mc {
		if strings.Index(_userFlags.matchResponse, ",") == -1 {
			code, _ := strconv.Atoi(_userFlags.matchResponse)
			if response.StatusCode != code {
				return
			}
		}

		for _, v := range strings.Split(_userFlags.matchResponse, ",") {
			code, _ := strconv.Atoi(v)
			if response.StatusCode != code {
				return
			}
		}
	}

	if _userFlags.output != "" {
		setResultInFile(_userFlags.output, response.Status, _httpMethod, response.StatusCode)
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
		wg        = sync.WaitGroup{}
		userFlags = Flags{
			extraMethods: false,
			target:       "",
			output:       "",
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
	flag.StringVar(&userFlags.output, "o", "", "Name Of File To Send Result in it")
	flag.StringVar(&userFlags.cookie, "c", "", "Set Cookie")
	flag.BoolVar(&userFlags.extraMethods, "x", false, "FUZZ Extra HTTP Methods")
	flag.StringVar(&userFlags.target, "d", "", "URL Of Your Target")
	flag.StringVar(&userFlags.matchResponse, "mc", "", "Match Response Code [use ',' To Split]")
	flag.StringVar(&userFlags.forgetResponse, "fc", "", "Don't Match Response Code [use ',' To Split]")
	flag.Parse()

	// Create Output File
	if userFlags.output != "" {
		createFile(userFlags.output)
	}

	// Check Target
	if userFlags.target == "" {
		fmt.Println("Invalid Target")
		os.Exit(0)
	}

	if !strings.HasPrefix(userFlags.target, "http://") && !strings.HasPrefix(userFlags.target, "https://") {
		fmt.Println("Invalid Prefix For Target")
		os.Exit(0)
	}

	banner()

	// Testing Base HTTP Methods
	for _, httpMethod := range baseHttpMethods {
		wg.Add(1)
		go func(_httpMethod string) {
			tamper(_httpMethod, userFlags)
			wg.Done()
		}(httpMethod)
	}

	// Testing Extra HTTP Methods
	if userFlags.extraMethods {
		for _, httpMethod := range extraHttpMethods {
			wg.Add(1)
			go func(_httpMethod string) {
				tamper(_httpMethod, userFlags)
				wg.Done()
			}(httpMethod)
		}
	}

	wg.Wait()
}
