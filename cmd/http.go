package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

type httpCmdFlags struct {
	Url     string
	Timeout string
}

type httpClientData struct {
	id      int
	total   int64
	average int64
	maxrtt  time.Duration
	minrtt  time.Duration
	count   int64
	errors  int
}

var httpFlags httpCmdFlags

func init() {
	rootCmd.AddCommand(httpCmd)
	httpCmd.Flags().StringVarP(&httpFlags.Url, "url", "u", "https://www.google.com", "Url to HTTP server to test.")
	httpCmd.Flags().StringVarP(&httpFlags.Timeout, "timeout", "t", "10s", "HTTP client time out duration.")
}

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Execute HTTP Tests",
	Long:  "aapep http - Used for stress (load) testing HTTP servers.",
	Run: func(cmd *cobra.Command, args []string) {
		if httpFlags.Url == "" {
			showHelp(cmd, "Error - Url can not be blank.")
		}

		timer, err := time.ParseDuration(rootFlags.Timer)
		if err != nil {
			showHelp(cmd, "Error - can not parse timer duration.")
		}

		fmt.Println("Starting Clients...")

		wg.Add(rootFlags.ClientCount)
		clientData := []*httpClientData{}
		for i := 0; i < rootFlags.ClientCount; i++ {
			d := httpClientData{id: i}
			clientData = append(clientData, &d)
			go httpClient(&d, rootFlags.Interval, timer)
		}

		fmt.Println("Waiting for clients to complete...")
		wg.Wait()
		fmt.Println("Clients completed")

		totals := httpClientData{}
		for _, x := range clientData {
			totals.average = totals.average + x.average
			totals.errors = totals.errors + x.errors
			if totals.maxrtt < x.maxrtt {
				totals.maxrtt = x.maxrtt
			}

			if totals.minrtt == 0 {
				totals.minrtt = x.minrtt
			}

			if x.minrtt < totals.minrtt {
				totals.minrtt = x.minrtt
			}
		}

		fmt.Printf("\nResults:\n")
		fmt.Printf("Average: %v\n", time.Duration(totals.average/int64(len(clientData))))
		fmt.Printf("Max Time: %v\n", totals.maxrtt)
		fmt.Printf("Min Time: %v\n", totals.minrtt)
		fmt.Printf("Errors: %v\n", totals.errors)
	},
}

func httpClient(data *httpClientData, interval int, timer time.Duration) {
	defer wg.Done()
	countMax := int64(0)
	s := time.Duration(interval) * time.Millisecond
	countMax = (timer * time.Millisecond).Milliseconds() / int64(s)

	t, _ := time.ParseDuration(httpFlags.Timeout)

	for {
		data.count++
		start := time.Now()
		transCfg := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: transCfg, Timeout: t}
		response, err := client.Get(httpFlags.Url)
		if err != nil {
			data.errors++
		} else {
			response.Body.Close()
			if response.StatusCode != 200 {
				data.errors++
			} else {
				rtt := time.Duration(time.Since(start))
				data.total = data.total + int64(rtt)
				data.average = data.total / data.count

				if rtt > data.maxrtt {
					data.maxrtt = rtt
				}

				if data.minrtt == 0 {
					data.minrtt = rtt
				}

				if rtt < data.minrtt {
					data.minrtt = rtt
				}
			}
		}

		if data.count >= countMax {
			break
		}
		time.Sleep(s)
	}
}
