package cmd

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
	"github.com/spf13/cobra"
)

type dnsCmdFlags struct {
	Server  string
	Name    string
	Type    string
	dnsType uint16
}

// type dnsClientMessage struct {
// 	average time.Duration
// 	err     bool
// }

type dnsClientData struct {
	id      int
	total   int64
	average int64
	maxrtt  time.Duration
	minrtt  time.Duration
	count   int64
	errors  int
}

var dnsFlags dnsCmdFlags

func init() {
	rootCmd.AddCommand(dnsCmd)
	dnsCmd.Flags().StringVarP(&dnsFlags.Server, "server", "s", "127.0.0.1:53", "Name or IP address of the name server to query.")
	dnsCmd.Flags().StringVarP(&dnsFlags.Name, "name", "n", "www.google.com.", "Name of the resource record that is to be looked up.")
	dnsCmd.Flags().StringVarP(&dnsFlags.Type, "type", "t", "A", "Indicates what type of query is required.")
}

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Execute DNS Tests",
	Long:  "aapep dns - Used for stress (load) testing DNS servers/forwarders.",
	Run: func(cmd *cobra.Command, args []string) {
		if dnsFlags.Server == "" {
			showHelp(cmd, "Error - name server can not be blank.")
		}

		if dnsFlags.Name == "" {
			showHelp(cmd, "Error - name of the resource record can not be blank.")
		}

		switch dnsFlags.Type {
		case "A":
			dnsFlags.dnsType = dns.TypeA
		case "CNAME":
			dnsFlags.dnsType = dns.TypeCNAME
		default:
			showHelp(cmd, fmt.Sprintf("Error - unspported query type (%v).", dnsFlags.Type))
		}

		timer, err := time.ParseDuration(rootFlags.Timer)
		if err != nil {
			showHelp(cmd, "Error - can not parse timer duration.")
		}

		// c1 := make(chan dnsClientMessage)
		// c2 := make(chan ui.DnsCmdResMsg)

		// go chanListener(c1, c2)

		fmt.Println("Starting Clients...")

		wg.Add(rootFlags.ClientCount)
		clientData := []*dnsClientData{}
		for i := 0; i < rootFlags.ClientCount; i++ {
			d := dnsClientData{id: i}
			clientData = append(clientData, &d)
			go dnsClient(&d, rootFlags.Interval, timer)
		}

		fmt.Println("Waiting for clients to complete...")
		wg.Wait()
		fmt.Println("Clients completed")

		totals := dnsClientData{}
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

		// m := ui.DnsCmdModel{
		// 	Res: c2,
		// }
		// m.CMD.Version = rootCmd.Version
		// m.CMD.Clients = rootFlags.ClientCount
		// m.CMD.Interval = rootFlags.Interval

		// p := tea.NewProgram(m)
		// if err := p.Start(); err != nil {
		// 	fmt.Println("could not start program:", err)
		// 	os.Exit(1)
		// }
	},
}

func dnsClient(data *dnsClientData, interval int, timer time.Duration) {
	defer wg.Done()
	countMax := int64(0)
	s := time.Duration(interval) * time.Millisecond
	countMax = (timer * time.Millisecond).Milliseconds() / int64(s)

	for {
		data.count++
		m := dns.Msg{}
		m.SetQuestion(dnsFlags.Name, dnsFlags.dnsType)

		d := dns.Client{}
		r, rtt, err := d.Exchange(&m, dnsFlags.Server)
		if err != nil {
			data.errors++
		} else {
			if len(r.Answer) < 1 {
				data.errors++
			} else {
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

// func chanListener(c1 chan dnsClientMessage, c2 chan ui.DnsCmdResMsg) {
// 	var errors int = 0
// 	var average time.Duration = 0
// 	for c := range c1 {
// 		if c.err {
// 			errors++
// 		} else {
// 			if average != 0 {
// 				average = (average + c.average) / 2
// 			} else {
// 				average = c.average
// 			}
// 		}
// 		c2 <- ui.DnsCmdResMsg{Average: average, Errors: errors}
// 	}
// }

// func dnsClient(c chan dnsClientMessage, interval int) {
// 	var avg time.Duration = 0
// 	var count int = 0
// 	var total time.Duration = 0
// 	for {
// 		count++
// 		m := dns.Msg{}
// 		m.SetQuestion(dnsFlags.Name, dnsFlags.dnsType)

// 		d := dns.Client{}
// 		r, rtt, err := d.Exchange(&m, dnsFlags.Server)
// 		if err != nil {
// 			c <- dnsClientMessage{average: time.Duration(0), err: true}
// 		} else {
// 			if len(r.Answer) < 1 {
// 				c <- dnsClientMessage{average: time.Duration(0), err: true}
// 			} else {
// 				total = total + rtt
// 				avg = total / time.Duration(count)
// 				c <- dnsClientMessage{average: avg, err: false}
// 			}
// 		}
// 		time.Sleep(time.Duration(interval) * time.Millisecond)
// 	}
// }
