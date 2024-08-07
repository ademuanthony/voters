package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type TicketsResponse struct {
	Hashes []string `json:"hashes"`
}

// const dcrctl = "/home/user/code/dcrctl/dcrctl"
// const dcrctl = "dcrctl"

const dcrctl = "./dcrctl.sh"

// var dcrctlArgs = []string{"--configfile=/home/user/.dcrctl/voter.conf", "--wallet"}

// var dcrctlArgs = []string{"--wallet", "--testnet"}
var dcrctlArgs = []string{}

var (
	yesZone float64
	noZone  float64
)

const (
	salt              = "DsYYaFKe3nxWJweGmCaVzPqr2qCa7Ve43ed"
	tspendOrPolicyKey = "03f6e7041f1cf51ee10e0a01cd2b0385ce3cd9debaabb2296f7e9dee9329da946c"
	verbose           = true
	repeatInterval    = 1 * time.Minute
)

var (
	totalYes     = 0
	totalNo      = 0
	totalAbstain = 0
	totalTickets = 0
)

func main() {
	fmt.Println()

	fmt.Printf("- Treasury Voter v1.0.0 (2023-11-01), will run every %v.\n\n", repeatInterval)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interrupt
		fmt.Println("\nClosing...")
		os.Exit(0)
	}()

	var err error

	yesZoneInput := getCliInput("Enter yes zone (for computing the number of yes votes):")
	yesZone, err = strconv.ParseFloat(yesZoneInput, 64)
	if err != nil {
		fmt.Println("Please enter a valid number for yes zone")
		return
	}

	noZoneInput := getCliInput("Enter no zone (for computing the number of no votes):")
	noZone, err = strconv.ParseFloat(noZoneInput, 64)
	if err != nil {
		fmt.Println("Please enter a valid number for no zone")
		return
	}

	absZone := 100 - yesZone - noZone

	assignedTickets := make(map[string]bool)

	round := 1

	for {
		if len(tspendOrPolicyKey) == 32 {
			fmt.Printf("*** loop %d ***  tspend %s\n", round, tspendOrPolicyKey)
		} else {
			fmt.Printf("*** loop %d ***  politeiakey %s\n", round, tspendOrPolicyKey)
		}

		tspendInMempool, err := getTspendInMempool()
		if err != nil {
			fmt.Println("error in getting tspend from mempool", err.Error())
		} else if len(tspendInMempool) > 0 {
			fmt.Printf("- found %d new tspend in mempool\n", len(tspendInMempool))
			for _, hash := range tspendInMempool {
				fmt.Println(hash)
			}
		}

		fmt.Printf(
			"- targets: yes %s%%  no %s%%  abstain %s%%, randzones: yes 0-%s  no %s-%s  abstain %s-100\n",
			formatPercentage(yesZone),
			formatPercentage(noZone),
			formatPercentage(absZone),
			formatPercentage(yesZone),
			formatPercentage(yesZone),
			formatPercentage(yesZone+noZone),
			formatPercentage(yesZone+noZone),
		)

		startGetTicketTime := time.Now()
		fmt.Printf("- get tickets... ")
		newTickets, removedTickets := getNewTickets(assignedTickets)

		if round > 1 {
			fmt.Printf("got %d new tickets  %d removed tickets  completed in %s.\n", len(newTickets), len(removedTickets), formatDuration(time.Since(startGetTicketTime)))
		} else {
			fmt.Printf("got %d tickets completed in %s.\n", len(newTickets), formatDuration(time.Since(startGetTicketTime)))
		}

		if round > 1 {
			if len(removedTickets) > 0 {
				fmt.Println("removed tickets")
				fmt.Println("Count \tTicket \t\t\t\t\tRand \tChoice \tSymbol")
				for i, hash := range removedTickets {
					determinant, policy := calculatePolicy(i+1, hash, salt, yesZone, noZone, verbose)
					fmt.Printf("%d \t%s \t%s \t%s \t%s\n", i+1, hash, formatPercentage(determinant), formatPolicy(policy, verbose), formatPolicy(policy, false))
				}
			}
		}

		round++

		if len(newTickets) > 0 {
			for _, ticketHash := range newTickets {
				assignedTickets[ticketHash] = true
			}

			newTicketsCount := len(newTickets)
			totalTickets += newTicketsCount

			fmt.Printf("Making random decisions for %d tickets\n", newTicketsCount)

			policyCounts := make(map[string]int)
			ticketPolicies := make(map[string]string)

			fmt.Println("new tickets")
			if verbose {
				fmt.Println("Count \tTicket \t\t\t\t\tRand \tChoice \tSymbol")
			}
			for i, ticketHash := range newTickets {
				determinant, policy := calculatePolicy(i+1, ticketHash, salt, yesZone, noZone, verbose)
				ticketPolicies[ticketHash] = policy
				policyCounts[policy]++

				setTspendPolicy(tspendOrPolicyKey, ticketHash, policy)

				if verbose {
					fmt.Printf("%d \t%s \t%s \t%s \t%s\n", i+1, ticketHash, formatPercentage(determinant), formatPolicy(policy, verbose), formatPolicy(policy, false))
				} else {
					fmt.Print(formatPolicy(policy, verbose))
				}
			}

			totalYes += policyCounts["yes"]
			totalNo += policyCounts["no"]
			totalAbstain += policyCounts["abstain"]
		}

		fmt.Println("\tyes \tno \tabs \ttotal")
		yesPercent := float64(100*totalYes) / float64(totalTickets)
		noPercent := float64(100*totalNo) / float64(totalTickets)
		absPercentage := 100 - (yesPercent + noPercent)

		fmt.Printf("votes \t%d \t%d \t%d \t%d\n", totalYes, totalNo, totalAbstain, totalTickets)
		fmt.Printf("perc \t%s \t%s \t%s\n", formatPercentage(yesPercent)+"%", formatPercentage(noPercent)+"%",
			formatPercentage(absPercentage)+"%")

		fmt.Printf("targ \t%s \t%s \t%s\n", formatPercentage(yesZone)+"%", formatPercentage(noZone)+"%",
			formatPercentage(absZone)+"%")
		fmt.Printf("diff \t%s \t%s \t%s\n",
			formatPercentage(math.Abs(yesPercent-yesZone))+"%",
			formatPercentage(math.Abs(noPercent-noZone))+"%",
			formatPercentage(math.Abs(absPercentage-absZone))+"%")

		nextRun := time.Now().Add(repeatInterval)
		fmt.Printf("- sleeping for %v, next run at %v...\n", repeatInterval, nextRun.Format("2006-01-02 15h-04m-05s"))
		fmt.Println()

		time.Sleep(repeatInterval)
	}
}

func getCliInput(prompt string) (input string) {
	fmt.Println(prompt)
	_, err := fmt.Scan(&input)
	if err != nil {
		fmt.Println("error in reading input", err)
	}
	return
}

func getTickets() (*TicketsResponse, error) {
	args := append(dcrctlArgs, "gettickets", "true")
	cmd := exec.Command(dcrctl, args...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprint(err) + ": " + stderr.String())
	}

	var ticketsResponse TicketsResponse
	if err := json.Unmarshal(out.Bytes(), &ticketsResponse); err != nil {
		return nil, fmt.Errorf("error parsing tickets JSON: %w", err)
	}
	return &ticketsResponse, nil
}

var previousTickts []string

func getNewTickets(assignedTickets map[string]bool) ([]string, []string) {
	ticketsResponse, err := getTickets()
	if err != nil {
		fmt.Println("Error fetching tickets:", err)
		return []string{}, []string{}
	}

	var allTickets = make(map[string]bool)
	for _, hash := range ticketsResponse.Hashes {
		allTickets[hash] = true
	}

	var removedTickets []string
	for _, hash := range previousTickts {
		if _, f := allTickets[hash]; !f {
			removedTickets = append(removedTickets, hash)
		}
	}

	previousTickts = ticketsResponse.Hashes

	var newTickets []string
	for _, ticketHash := range ticketsResponse.Hashes {
		if !assignedTickets[ticketHash] {
			newTickets = append(newTickets, ticketHash)
		}
	}
	return newTickets, removedTickets
}

var seenMempoolTspend = make(map[string]bool)

func getTspendInMempool() ([]string, error) {
	args := append(dcrctlArgs, "getrawmempool", "false", "tspend")
	cmd := exec.Command(dcrctl, args...)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var hashStrings []string
	err = json.Unmarshal(output, &hashStrings)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, hash := range hashStrings {
		if seenMempoolTspend[hash] {
			continue
		}

		seenMempoolTspend[hash] = true
		result = append(result, hash)
	}

	return result, nil
}

func parsePercentage(percentageStr string) float64 {
	percentage, err := strconv.ParseFloat(percentageStr, 64)
	if err != nil {
		fmt.Println("Invalid percentage:", err)
		os.Exit(1)
	}
	return percentage
}

func formatPercentage(percentage float64) string {
	if percentage == float64(int(percentage)) {
		return fmt.Sprintf("%.0f", percentage)
	}
	return fmt.Sprintf("%.2f", percentage)
}

func formatPolicy(policy string, verbose bool) string {
	if !verbose {
		switch policy {
		case "yes":
			return "+"

		case "no":
			return "-"

		default:
			return ""
		}

	}
	return policy
}

func formatDuration(duration time.Duration) string {
	units := []struct {
		unit  string
		value int64
	}{
		{"h", int64(time.Hour)},
		{"m", int64(time.Minute)},
		{"s", int64(time.Second)},
		{"ms", int64(time.Millisecond)},
	}

	formattedDuration := ""

	for _, unit := range units {
		if duration >= time.Duration(unit.value) {
			count := int64(duration / time.Duration(unit.value))

			formattedDuration += fmt.Sprintf("%d%s", count, unit.unit)

			duration -= time.Duration(count) * time.Duration(unit.value)
		}
	}

	return formattedDuration
}

func calculatePolicy(no int, ticketHash, salt string, yesZone, noZone float64, verbose bool) (float64, string) {
	data := ticketHash + salt
	hashed := sha256.Sum256([]byte(data))
	seed := new(big.Int).SetBytes(hashed[:]).Uint64()
	r := rand.New(rand.NewSource(int64(seed)))

	determinant := r.Float64() * 100

	var policy string
	if determinant <= yesZone {
		policy = "yes"
	} else if determinant <= yesZone+noZone {
		policy = "no"
	} else {
		policy = "abstain"
	}

	return determinant, policy
}

func setTspendPolicy(tspendOrPolicyKey, ticketHash, policy string) {
	if len(tspendOrPolicyKey) == 32 {
		args := append(dcrctlArgs, "settspendpolicy", tspendOrPolicyKey, policy, ticketHash)
		cmd := exec.Command(dcrctl, args...)
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error in setting policy", err.Error())
		}
	} else {
		args := append(dcrctlArgs, "settreasurypolicy", tspendOrPolicyKey, policy, ticketHash)
		cmd := exec.Command(dcrctl, args...)
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error in setting policy", err.Error())
		}
	}
}
