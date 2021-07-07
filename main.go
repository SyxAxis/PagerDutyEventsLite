package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var relversion = "1.2"

// ====================================================================================================
//
//  Auth: G. Johnson
//  Date: 14-NOV-2020
//  Desc: Lite version of the master PD util that will work on lower Windows versions like Win2008
//
// ====================================================================================================
//
// BUILDING - Assumes PowerShell command line in use on Windows for building Win and Linux Versions
//
// Windows build:
//
//  $env:GOOS="windows"
//  go build pagerdutylite.go
//
// LINUX build
//
//  $env:GOOS="linux"
//  go build pagerdutylite.go
//
// If bulding Win/Lin on Linux use "export" command to set the same params for each platform
//
// ====================================================================================================
//
// original flags from main PD util
//
// usage: PagerDuty events [-h] --routing_key ROUTING_KEY [--msg MSG]
//                         [--source SOURCE] --keyname KEYNAME
//                         [--event {TRIGGER,RESOLVE,ACKNOWLEDGE}]
//                         [--details DETAILS]
//                         [--jdetails JDETAILS]
//                         [--severity {info,critical,warning,error}]
//                         [--retries RETRIES] [--retry_interval {30,60,120,300}]
//                         [--proxy_server PROXY_SERVER]
//
// limited options implemented in this util
//
// usage: PagerDuty events [-h]
//        					--routing_key ROUTING_KEY
//  						--msg MSG
//                         	--source SOURCE
//                          --keyname KEYNAME
//                          --event {TRIGGER,RESOLVE,ACKNOWLEDGE}
//                          --details DETAILS
//                          --jdetails JSON_LOG_DETAILS
//                          --severity {info,critical,warning,error}
//                          --proxy_server PROXY_SERVER
//
// ====================================================================================================

var showlog bool
var rtnJSONrslt bool

// type PyLdContext struct {
// 	CtxType   string `json:"type"` // link or image
// 	CtxHref   string `json:"href"` // URL
// 	CtxText   string `json:"text"` // Alt text
// 	CtxSource string `json:"src"`  // for image this would URI to the online image
// }

type PayLoad struct {
	PyLdSummary   string           `json:"summary"`        // main incident title
	PyLdSource    string           `json:"source"`         // source of the incident ( hostname, appname, etc)
	PyLdSeverity  string           `json:"severity"`       // info warning error critical
	CustomDetails *json.RawMessage `json:"custom_details"` // takes raw JSON to be used to make a simple info table in the PD alert page. JSON is simply sets of K:V pairs in an array

	// removed
	// PyLdComponent string `json:"component"` // system sub-component type ( db, restsvr, etc )
	// PyLdGroup     string `json:"group"`     // grouping if several incidents get raised ( pricingApp, warehouseheating, etc )
	// PyLdClass     string `json:"class"`     // type of error in the incident ( highCPU, Lowtemp, diskSpace, etc )
	// PyLdClient    string `json:"client"`   		//
	// PyLDTimestamp string `json:"timestamp"`		// override PD timestamp, useful if there's a delay in delivering to PD
}

type PDEventTrigger struct {
	RoutingKey  string  `json:"routing_key"`  // the key used from within a service for an API call
	EventAction string  `json:"event_action"` // trigger acknowledge resolve
	DeDupeKey   string  `json:"dedup_key"`    // unique key you use to log and update on, you set it
	PyLd        PayLoad `json:"payload"`      // payload struct
}

type PDEventResolve struct {
	RoutingKey  string `json:"routing_key"`  // the key used from within a service for an API call
	EventAction string `json:"event_action"` // trigger acknowledge resolve
	DeDupeKey   string `json:"dedup_key"`    // unique key you use to log and update on, you set it
}

type PDEventAcknowledge struct {
	RoutingKey  string `json:"routing_key"`  // the key used from within a service for an API call
	EventAction string `json:"event_action"` // trigger acknowledge resolve
	DeDupeKey   string `json:"dedup_key"`    // unique key you use to log and update on, you set it
}

//
func FlagUsage() {
	fmt.Println("\n")
	fmt.Printf("PagerDuty Util Lite - %s", relversion)
	fmt.Println("\n")
	fmt.Println("--routing_key  	", "<string> - The primary routing key for the PD event rule or service")
	fmt.Println("--keyname      	", "<string> - Unique user defined key.")
	fmt.Println("--event        	", "<string> - Must be one of { trigger | acknowledge | resolve}")
	fmt.Println("--severity     	", "<string> - {info | critical | warning | error} ")
	fmt.Println("--msg          	", "<string> - Primary message alert title.")
	fmt.Println("--source       	", "<string> - Source of the alert, advise use of hostname.")
	fmt.Println("--details      	", "<string> - Simple logging details for the alert.")
	fmt.Println("--jdetails     	", "<string> - JSON formatted structure declaring sets of key:value pairs with log information sets. ( OPT )")
	fmt.Println("--proxy_server 	", "<string> - Force specific proxy server ( default use HTTP_PROXY/HTTPS_PROXY from environment).")
	fmt.Println("--JSONresult      	", "Return result in JSON format.")
	fmt.Println("--saveJSONresponse	", "Save the JSON result to a file names <keyname>.json.")
	fmt.Println("\n")
	fmt.Println("Note : ")
	fmt.Println("  - If you need to use a proxy then set HTTP_PROXY or HTTPS_PROXY in the environment first ( set on Windows cmd line or export on Unix ).")
	fmt.Println("  - When supplying <jdetails> param, be sure to escape the double-quotes. PowerShell : '{\\\"key01\\\":\\\"value01\\\"}'")
	fmt.Println("  - When supplying <jdetails> param it will override <details> param, although both will be sent to PD if supplied.")
	fmt.Println("\n")

}

// primary outtput but only shows if the "--showlog" param is invoked
func FuncOutputMsg(outMsg string) {
	if showlog {
		log.Println(outMsg)
	}
}

func PDRequest(jsonReq []byte) (retVal int, bodyResult string) {

	FuncOutputMsg("Proceeding       : Submitting request to PagerDuty")

	// new HTTP request
	req, err := http.NewRequest(http.MethodPost, "https://events.pagerduty.com/v2/enqueue", bytes.NewBuffer(jsonReq))
	// basic headers
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("accept", "application/vnd.pagerduty+json;version=2")
	// client object up
	client := &http.Client{}
	// send the request and check it for errors
	// note errors during the call, data issues will not be known unless scanned
	resp, err := client.Do(req)
	if err != nil {
		// this will cause an exit
		log.Fatalln(err)
	}

	// leave the connection open until we're done in this func
	defer resp.Body.Close()
	// collect the result from the call
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// Convert response body to string
	bodyString := string(bodyBytes)

	FuncOutputMsg("Result           : " + bodyString)

	if rtnJSONrslt {
		fmt.Println(bodyString)
	}

	if strings.Contains(bodyString, "\"status\":\"success\"") {
		return 0, bodyString
	} else {
		return 1, bodyString
	}

}

func DumpJSONResultToFile(PDCallJSONDumPFile string, PDCallJSONResponseBody string) {
	FuncOutputMsg("Write response   : " + PDCallJSONDumPFile + ".json")
	datOut := []byte(PDCallJSONResponseBody)
	err := ioutil.WriteFile(PDCallJSONDumPFile+".json", datOut, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {

	clfRoutingKey := flag.String("routing_key", "", "<string> - The primary routing key for the PD event rule or service")
	clfMessageSummary := flag.String("msg", "", "<string> - Primary message alert title.")
	clfAlertSource := flag.String("source", "", "<string> - Source of the alert, advise use of hostname.")
	clfUniqueKey := flag.String("keyname", "", "<string> - Unique user defined key.")
	clfAction := flag.String("event", "trigger", "<string> - Must be one of { trigger | acknowledge | resolve}")
	clfAlertDetails := flag.String("details", "", "<string> - JSON formatted details for the alert")
	clfCustomJSONDetails := flag.String("jdetails", "", "<string> - JSON formatted details for the alert")
	clfSeverity := flag.String("severity", "", "<string> - {info | critical | warning | error} ")
	clfRtnJSONRslt := flag.Bool("JSONresult", false, "Return result in JSON format.")
	clfShowOpsLog := flag.Bool("showlog", false, "Display the output of the operations.")
	clfDumpJSONResultToFile := flag.Bool("saveJSONresponse", false, "Save the JSON result to a file names <keyname>.json ")

	// clfJSONLogDetails := flag.String("jdetails", "", "JSON formatted log details.")
	clfProxyServer := flag.String("proxy_server", "", "Force use of specific proxy server.")

	// despite the above defaults, there is a custom help message defined in a func
	flag.Usage = FlagUsage

	// parse the inbound flags
	flag.Parse()

	// fmt.Println(*clfCustomJSONDetails)
	// fmt.Println(*clfUniqueKey)

	// optional flag to output the log
	showlog = *clfShowOpsLog
	rtnJSONrslt = *clfRtnJSONRslt

	FuncOutputMsg("Routing Key      : " + *clfRoutingKey)
	FuncOutputMsg("Message Key      : " + *clfUniqueKey)

	// start building the struct for JSON delivery
	jsonReq := []byte{}

	var takeAction string
	if strings.ToUpper(*clfAction) == "TRIGGER" || strings.ToUpper(*clfAction) == "T" {
		takeAction = "trigger"
	} else if strings.ToUpper(*clfAction) == "ACKNOWLEDGE" || strings.ToUpper(*clfAction) == "A" {
		takeAction = "acknowledge"
	} else if strings.ToUpper(*clfAction) == "RESOLVE" || strings.ToUpper(*clfAction) == "R" {
		takeAction = "resolve"
	}

	FuncOutputMsg("Action           : " + takeAction)

	if takeAction == "trigger" {

		FuncOutputMsg("Severity         : " + *clfSeverity)
		FuncOutputMsg("Message Source   : " + *clfAlertSource)
		FuncOutputMsg("Message Summary  : " + *clfMessageSummary)

		if len(*clfProxyServer) > 0 {
			FuncOutputMsg("Proxy Server     : " + *clfProxyServer)
		}

		// define type, start with nothing
		customJSONDetails := json.RawMessage("{}")
		// if it needs to be updated with a real value off the cmd line then swap that in
		if len(*clfCustomJSONDetails) > 0 && len(*clfAlertDetails) > 0 {
			FuncOutputMsg("Custom Details   : " + *clfCustomJSONDetails)
			FuncOutputMsg("Alert Details    : " + *clfAlertDetails)
			tmpJSONStr := *clfCustomJSONDetails
			if strings.HasPrefix(tmpJSONStr, "{") {
				// tmpJSONStr = tmpJSONStr[:len(tmpJSONStr)-len("}")]
				tmpJSONStr = tmpJSONStr[1:]
				newJSONStr := "{\"details\":\"" + *clfAlertDetails + "\"," + tmpJSONStr
				customJSONDetails = json.RawMessage(newJSONStr)
			} else {
				FuncOutputMsg("Unable to append 'details' into 'jdetails'. 'jdetails' param must start with valid '{' character.")
				customJSONDetails = json.RawMessage(*clfCustomJSONDetails)
			}
		} else if len(*clfAlertDetails) > 0 {
			FuncOutputMsg("Alert Details    : " + *clfAlertDetails)
			customJSONDetails = json.RawMessage("{\"details\":\"" + *clfAlertDetails + "\"}")
		} else if len(*clfCustomJSONDetails) > 0 {
			FuncOutputMsg("Custom Details   : " + *clfCustomJSONDetails)
			customJSONDetails = json.RawMessage(*clfCustomJSONDetails)
		}

		// now fill in the structs that will be used to draw up the JSON to send off
		pyld := PayLoad{*clfMessageSummary,
			*clfAlertSource,
			*clfSeverity,
			&customJSONDetails}

		mainPyLd := PDEventTrigger{
			*clfRoutingKey,
			takeAction,
			*clfUniqueKey,
			pyld}

		jsonReq, _ = json.Marshal(mainPyLd)
	}

	if takeAction == "acknowledge" {
		mainPyLd := PDEventAcknowledge{
			*clfRoutingKey,
			takeAction,
			*clfUniqueKey}

		// marshal up the struct into a JSON formatted string
		jsonReq, _ = json.Marshal(mainPyLd)
	}

	if takeAction == "resolve" {
		mainPyLd := PDEventResolve{
			*clfRoutingKey,
			takeAction,
			*clfUniqueKey}

		// marshal up the struct into a JSON formatted string
		jsonReq, _ = json.Marshal(mainPyLd)
	}

	FuncOutputMsg("JSON request     : " + string(jsonReq))

	// post it to PD and see if they accept it
	PDCallResponseStatus, PDCallJSONResponseBody := PDRequest(jsonReq)

	if *clfDumpJSONResultToFile {
		DumpJSONResultToFile(*clfUniqueKey, PDCallJSONResponseBody)
	}

	os.Exit(PDCallResponseStatus)

}
