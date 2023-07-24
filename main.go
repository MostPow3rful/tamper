package main

import (
	"bufio"
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
	CustomHeaders  string
	Output         string
	Target         string
	Cookie         string
	ExtraMethods   bool
}

var (
	customHeaders = make(map[string]string, 0)
	flagInstance  = Flag{
		IgnoreResponse: "",
		MatchResponse:  "",
		CustomHeaders:  "",
		Output:         "",
		Target:         "",
		Cookie:         "",
		ExtraMethods:   false,
	}
)

func tamper(_httpMethod string) {
	mc := !(flagInstance.MatchResponse == "")
	fc := !(flagInstance.IgnoreResponse == "")

	// Create an instance Of http.Client Struct
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	// Creating Custom Request
	request, err := http.NewRequest(_httpMethod, flagInstance.Target, nil)
	if err != nil {
		fmt.Printf(
			"[%s] Couldn't Make Request [%s]\n",
			_httpMethod, flagInstance.Target,
		)
		return
	}

	// Set HTTP Headers
	for header, value := range customHeaders {
		request.Header.Set(header, value)
	}
	if flagInstance.Cookie != "" {
		request.Header.Set("Cookie", flagInstance.Cookie)
	}
	request.Header.Set("User-Agent", uarand.GetRandom())

	// Sending Request
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf(
			"[%s] Couldn't Send Request [%s]\n",
			_httpMethod, flagInstance.Target,
		)
		return
	}

	// Close Response Body [defer]
	defer response.Body.Close()

	// if fc switch defined
	if fc {
		if strings.Index(flagInstance.IgnoreResponse, ",") == -1 {
			code, _ := strconv.Atoi(flagInstance.IgnoreResponse)
			if response.StatusCode == code {
				return
			}
		}

		for _, v := range strings.Split(flagInstance.IgnoreResponse, ",") {
			code, _ := strconv.Atoi(v)
			if response.StatusCode == code {
				return
			}
		}
	}

	//if mc switch defined
	if mc {
		if strings.Index(flagInstance.MatchResponse, ",") == -1 {
			code, _ := strconv.Atoi(flagInstance.MatchResponse)
			if response.StatusCode != code {
				return
			}
		}

		for _, v := range strings.Split(flagInstance.MatchResponse, ",") {
			code, _ := strconv.Atoi(v)
			if response.StatusCode != code {
				return
			}
		}
	}

	// if defined output file
	if flagInstance.Output != "" {
		setResultInFile(flagInstance.Output, response.Status, _httpMethod, response.StatusCode)
		return
	}
	fmt.Printf("[%v] - [%d] - [%v]\n", _httpMethod, response.StatusCode, response.Status)
}

func parseCustomHeaders(_fileName string) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Couldn't Get Output Of PWD Command")
		return
	}

	file, err := os.Open(fmt.Sprintf("%v/%v", pwd, _fileName))
	if err != nil {
		fmt.Printf("Couldn't Open Output File")
		return
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, " ") {
			line = line[1:]
		}

		// Empty Headers Only ?
		if strings.Index(line, ":") == -1 {
			if strings.HasSuffix(line, " ") {
				customHeaders[line[0:strings.Index(line, " ")]] = "127.0.0.1"
				continue
			}

			customHeaders[line] = "127.0.0.1"
			continue
		}

		// Header And Defined Value
		if strings.Index(line, ":") != -1 {
			if strings.HasSuffix(line, " ") {
				line = line[0:strings.Index(line, " ")]
			}

			data := strings.Split(line, ":")
			if data[0] == "" {
				continue
			}
			if data[1] == "" {
				data[1] = "127.0.0.1"
			}

			if strings.HasSuffix(data[0], " ") {
				data[0] = data[0][0:strings.LastIndex(data[0], " ")]
			}
			if strings.HasPrefix(data[1], " ") {
				data[1] = data[1][strings.LastIndex(data[1], " ")+1:]
			}
			customHeaders[data[0]] = data[1]
		}
	}

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
	defer outputFile.Close()

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
		wg              = sync.WaitGroup{}
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
		extraHttpMethods = [41]string{
			"UPDATEREDIRECTREF",
			"BASELINE-CONTROL",
			"VERSION-CONTROL",
			"MKREDIRECTREF",
			"MKWORKSPACE",
			"MKACTIVITY",
			"SHOWMETHOD",
			"TEXTSEARCH",
			"ORDERPATCH",
			"UNCHECKOUT",
			"MKCALENDAR",
			"PROPPATCH",
			"ARBITRARY",
			"SPACEJUMP",
			"BAMBOOZL",
			"NOEXISTE",
			"PROPFIND",
			"CHECKOUT",
			"CHECKIN",
			"SEARCH",
			"UNLOCK",
			"REPORT",
			"REBIND",
			"UNBIND",
			"UNLINK",
			"UPDATE",
			"PURGE",
			"REPOR",
			"TRACK",
			"INDEX",
			"QUERY",
			"LABEL",
			"MERGE",
			"MKCOL",
			"BIND",
			"LINK",
			"COPY",
			"MOVE",
			"LOCK",
			"ACL",
			"PRI",
		}
	)

	// Parse The User Flags
	flag.StringVar(&flagInstance.IgnoreResponse, "fc", "", "Don't Match Response Code [use ',' To Split]")
	flag.StringVar(&flagInstance.MatchResponse, "mc", "", "Match Response Code [use ',' To Split]")
	flag.StringVar(&flagInstance.Target, "d", "", "URL Of Your Target Do You Want To Test")
	flag.StringVar(&flagInstance.CustomHeaders, "h", "", "Set Custom Headers To Test")
	flag.StringVar(&flagInstance.Output, "o", "", "Name Of File To Set Result in it")
	flag.StringVar(&flagInstance.Cookie, "c", "", "Set Value Of Cookie Header")
	flag.BoolVar(&flagInstance.ExtraMethods, "x", false, "FUZZ Extra HTTP Methods")
	flag.Parse()

	// Check Target
	if !strings.HasPrefix(flagInstance.Target, "http://") && !strings.HasPrefix(flagInstance.Target, "https://") {
		fmt.Println("Invalid Prefix For Target")
		os.Exit(0)
	}
	if flagInstance.Target == "" {
		fmt.Println("Invalid Target")
		os.Exit(0)
	}

	// Set Custom Headers
	if flagInstance.CustomHeaders != "" {
		parseCustomHeaders(flagInstance.CustomHeaders)
	}

	// Create Output File
	if flagInstance.Output != "" {
		createFile(flagInstance.Output)
	}

	banner()

	// Testing Base HTTP Methods
	for _, httpMethod := range baseHttpMethods {
		wg.Add(1)
		go func(_httpMethod string) {
			tamper(_httpMethod)
			wg.Done()
		}(httpMethod)
	}

	// Testing Extra HTTP Methods
	if flagInstance.ExtraMethods {
		for _, httpMethod := range extraHttpMethods {
			wg.Add(1)
			go func(_httpMethod string) {
				tamper(_httpMethod)
				wg.Done()
			}(httpMethod)
		}
	}

	wg.Wait()
}
