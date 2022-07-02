package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs/v2"
)

func main() {

	source_file := flag.String("source_file", "", "path of source file")
	source_directory := flag.String("source_directory", "", "path of source directory file")
	output_directory := flag.String("output_directory", "output", "specify the output directory")
	flag.Parse()

	if len(os.Args) < 3 {
		fmt.Println(`Usage of ./migrator:
		-output_directory string
			  specify the output directory
		-source_directory string
			  path of source directory file
		-source_file string
			  path of source file`)
	}

	if *source_directory != "" {
		source_directoryValue := *source_directory
		output_directoryValue := *output_directory
		files, err := ioutil.ReadDir(source_directoryValue)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			modification(source_directoryValue+file.Name(), output_directoryValue)
		}
	}

	if *source_file != "" {
		source_fileValue := *source_file
		output_directoryValue := *output_directory
		modification(source_fileValue, output_directoryValue)

	}

}

/*
func test(filename string, output_directoryValue string) {
	//fmt.Print(filename)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Error when opening trest file: ", err)
	}
	if _, err := os.Stat(output_directoryValue + "/" + filename); os.IsNotExist(err) {
		os.MkdirAll(output_directoryValue, 0777) // Create your file
	}
	err = ioutil.WriteFile(output_directoryValue+"/"+filepath.Base(filename), content, 0777)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
}
*/

func modification(filename string, output_directoryValue string) {

	jsonParsed, err := gabs.ParseJSON([]byte(`{
		"templating": [
			{
			  "current": {
				"selected": true,
				"tags": [],
				"text": [
				  "Prometheus"
				],
				"value": [
				  "Prometheus"
				]
			  },
			  "description": null,
			  "error": null,
			  "hide": 0,
			  "includeAll": false,
			  "label": null,
			  "multi": true,
			  "name": "prom_datasource",
			  "options": [],
			  "query": "prometheus",
			  "queryValue": "",
			  "refresh": 1,
			  "regex": "",
			  "skipUrlSync": false,
			  "type": "datasource"
			},
			{
			  "current": {
				"selected": true,
				"tags": [],
				"text": [
				  "Aries-Functional-Tests"
				],
				"value": [
				  "Aries-Functional-Tests"
				]
			  },
			  "description": null,
			  "error": null,
			  "hide": 0,
			  "includeAll": false,
			  "label": null,
			  "multi": true,
			  "name": "es_datasource",
			  "options": [],
			  "query": "elasticsearch",
			  "queryValue": "",
			  "refresh": 1,
			  "regex": "",
			  "skipUrlSync": false,
			  "type": "datasource"
			},
			{
			  "current": {
				"selected": true,
				"tags": [],
				"text": [
				  "InfluxDB"
				],
				"value": [
				  "InfluxDB"
				]
			  },
			  "description": null,
			  "error": null,
			  "hide": 0,
			  "includeAll": false,
			  "label": null,
			  "multi": true,
			  "name": "influx_datasource",
			  "options": [],
			  "query": "influxdb",
			  "queryValue": "",
			  "refresh": 1,
			  "regex": "",
			  "skipUrlSync": false,
			  "type": "datasource"
			},
			{
			  "current": {
				"selected": true,
				"tags": [],
				"text": [
				  "CloudWatch"
				],
				"value": [
				  "CloudWatch"
				]
			  },
			  "description": null,
			  "error": null,
			  "hide": 0,
			  "includeAll": false,
			  "label": null,
			  "multi": true,
			  "name": "cloudwatch_datasource",
			  "options": [],
			  "query": "cloudwatch",
			  "queryValue": "",
			  "refresh": 1,
			  "regex": "",
			  "skipUrlSync": false,
			  "type": "datasource"
			},
			{
			  "current": {
				"selected": true,
				"tags": [],
				"text": [
				  "pgwatch2_pg_metrics"
				],
				"value": [
				  "pgwatch2_pg_metrics"
				]
			  },
			  "description": null,
			  "error": null,
			  "hide": 0,
			  "includeAll": false,
			  "label": null,
			  "multi": true,
			  "name": "postgres_datasource",
			  "options": [],
			  "query": "postgres",
			  "queryValue": "",
			  "refresh": 1,
			  "regex": "",
			  "skipUrlSync": false,
			  "type": "datasource"
			}
		]
	  }`))

	if err != nil {
		panic(err)
	}
	// Let's first read the `config.json` file
	//output: fmt.Println(listname, jsonOutput)
	// change data : payload.SetP("ritesh", "panels.0.datasource.type")
	// fetech data : listname := payload.Path("panels.0.datasource.type").Data()
	//jsonOutput := payload.String()
	fmt.Println(filename)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	if _, err := os.Stat(output_directoryValue + "/" + filename); os.IsNotExist(err) {
		os.MkdirAll(output_directoryValue, 0777) // Create your file
	}
	// Now let's unmarshall the data into `payload`
	payload, _ := gabs.ParseJSON(content)

	//fmt.Println(jsonParsed)
	log.Println("Adding datasource variable")

	for _, child := range jsonParsed.S("templating").Children() {
		payload.ArrayAppend(child.Data(), "templating", "list")
	}
	log.Println("Adding datasource done")

	log.Println("Starting to  modify  panel in the dashboard json")
	for key, child := range payload.S("panels").Children() {
		modi := string("panels.") + strconv.Itoa(key) + string(".datasource")

		if strings.HasPrefix(strings.ToLower(child.Path("datasource").String()), "\"prometheus") {

			payload.SetP("${prom_datasource}", modi)
		}

		if strings.HasPrefix(strings.ToLower(child.Path("datasource").String()), "\"elasticsearch") || strings.HasPrefix(strings.ToLower(child.Path("datasource").String()), "\"es") {
			payload.SetP("${es_datasource}", modi)
		}

		if strings.HasPrefix(strings.ToLower(child.Path("datasource").String()), "\"postgresql") {
			payload.SetP("${postgres_datasource}", modi)
		}

		if strings.HasPrefix(strings.ToLower(child.Path("datasource").String()), "\"influxdb") {
			payload.SetP("${influx_datasource}", modi)
		}

		if strings.HasPrefix(strings.ToLower(child.Path("datasource").String()), "\"cloudwatch") {
			payload.SetP("${cloudwatch_datasource}", modi)
		}

	}
	log.Println("Modify panel json done ")
	//fmt.Println(payload.String())
	output_json, _ := json.MarshalIndent(payload.Data(), "", "\t")
	err = ioutil.WriteFile(output_directoryValue+"/"+filepath.Base(filename), output_json, 0777)
	if err != nil {
		log.Fatal("Error Writing to  file: ", err)
	}
}
