// Copyright 2022 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Serves the `resultFile` as a response to queries. This is useful
// for assembling a custom healthcheck as a sidecar.
//
// resultFile is of format <status code> SPACE <response text>
//
// You need a separate container that does the custom set of
// healthchecking on your apps and write a resulting status file. See
// `healtcheck-example.sh`.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var f = struct {
	bindAddr string
	resultFile string
}{
	bindAddr: ":8080",
	resultFile: "/result",
}


func main() {
	flag.StringVar(&f.bindAddr, "bindAddr", f.bindAddr, "Address to bind healthcheck response from")
	flag.StringVar(&f.resultFile, "resultFile", f.resultFile, `Location of the result file. File should be a single line of the format <status code> SPACE <response>; example: "500 server is unhealthy"`)
	flag.Parse()

	http.HandleFunc("/", handle)
	log.Printf("Listening on %q", f.bindAddr)
	log.Fatal(http.ListenAndServe(f.bindAddr, nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	resultText, err := os.ReadFile(f.resultFile)
	if err != nil {
		log.Fatalf("os.ReadFile(%q) = _, %v", err)
	}
	result := strings.SplitN(string(resultText), " ", 2)
	if len(result) != 2 {
		w.WriteHeader(500)
		err := fmt.Errorf("Invalid resultsText format (%q)", resultText)
		fmt.Fprintf(w, "%v\n", err)
		log.Print(err)
		return
	}
	code, err := strconv.Atoi(result[0])
	if err != nil {
		w.WriteHeader(500)
		err := fmt.Errorf("Invalid resultText %q: invalid status code", resultText)
		fmt.Fprintf(w, "%v\n", err)
		log.Print(err)
		return
	}
	if code < 200 || code >= 600 {
		w.WriteHeader(500)
		err := fmt.Errorf("Invalid resultText %q: invalid status code value", resultText)
		fmt.Fprintf(w, "%v\n", err)
		log.Print(err)
		return
	}
	w.WriteHeader(int(code))
	w.Write([]byte(result[1]))
}
