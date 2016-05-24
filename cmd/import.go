// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/estesp/elastistack/goroutine"
	"github.com/maruel/panicparse/stack"
	es "github.com/mattbaird/elastigo/lib"
	"github.com/spf13/cobra"
)

const (
	esDefaultHost     = "localhost"
	esDefaultPort     = 9200
	esBulkConnections = 5
)

var (
	defaultIndex      = "stacktrace"
	defaultType       = "goroutine"
	esHost, inputFile string
	esPort            int
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a stack trace into Elasticsearch",
	Long: `Given a textual Golang stack trace, the import
command will parse the input file and insert the stack
trace data into Elasticsearch for further analysis.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		esConn := es.NewConn()
		esConn.SetFromUrl(fmt.Sprintf("http://%s:%d", esHost, esPort))

		bulkIndexer := esConn.NewBulkIndexer(esBulkConnections)

		if inputFile == "" {
			log.Error("Must specify a filename for --input")
			return fmt.Errorf("Must specify a filename for --input")
		}

		inputBytes, err := ioutil.ReadFile(inputFile)
		if err != nil {
			log.Errorf("Error reading from %q: %v", inputFile, err)
			return err
		}
		fReader := bytes.NewReader(inputBytes)
		log.Info("parsing stack dump")
		routines, err := stack.ParseDump(fReader, nil)
		if err != nil {
			log.Errorf("Error trying to parse dump: %v", err)
			return err
		}
		bulkIndexer.Start()
		defer esConn.Close()
		defer bulkIndexer.Stop()

		// base time from which offsets will be calculated
		timeNow := time.Now()

		log.Infof("Loading routine data in elastic search..%d\n", len(routines))
		for idx, routine := range routines {
			log.Debugf("[%03d] routine #%d\n", idx, routine.ID)
			rTime := timeNow.Add(time.Minute * time.Duration(-routine.SleepMin))
			routineEntry := goroutine.NewGoroutineTrace(routine, rTime)
			bulkIndexer.Index(defaultIndex, defaultType, fmt.Sprintf("%d", idx), "", "", &rTime, routineEntry)
		}

		done := make(chan struct{}, 1)
		go func() {
			for bulkIndexer.PendingDocuments() > 0 {
				time.Sleep(2 * time.Second)
			}
			done <- struct{}{}
		}()

		<-done
		return nil
	},
}

func init() {
	RootCmd.AddCommand(importCmd)

	importCmd.PersistentFlags().StringVarP(&esHost, "host", "e", esDefaultHost, "Hostname for Elasticstack endpoint (localhost)")
	importCmd.PersistentFlags().IntVarP(&esPort, "port", "p", esDefaultPort, "Port for Elasticstack endpoint (9200)")

	importCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input filename with Golang stack trace data")
}
