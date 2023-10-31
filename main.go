package main

import (
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

	"github.com/rodaine/table"
)

type TicketsResponse struct {
	Hashes []string `json:"hashes"`
}

const dcrctl = "./dcrctl.sh"

// var dcrctlArgs = []string{"--configfile=/home/user/.dcrctl/voter.conf"}
var dcrctlArgs = []string{}

const salt = "DsYYaFKe3nxWJweGmCaVzPqr2qCa7Ve43ed"

func main() {
	fmt.Println()
	if len(os.Args) != 4 && len(os.Args) != 5 {
		fmt.Println("Invalid number of arguments. Expected 4 or 5 but got ", len(os.Args))
		fmt.Println("Usage: ./voters <Treasury Transaction Hash or Policy Key> <Yes Percentage> <No Percentage> <Verbose>")
		os.Exit(1)
	}

	fmt.Println("Starting...")

	// Create a channel to listen for Ctrl+C (interrupt) signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Start a goroutine to handle interrupt signals
	go func() {
		<-interrupt
		fmt.Println("\nClosing...") // Print when Ctrl+C is pressed
		os.Exit(0)
	}()

	tspendOrPolicyKey := os.Args[1]
	yesZone := parsePercentage(os.Args[2])
	noZone := parsePercentage(os.Args[3])
	absZone := 100 - yesZone - noZone

	var verbose bool

	if len(os.Args) == 5 {
		verbose = os.Args[4] == "true"
	}

	if len(tspendOrPolicyKey) == 64 {
		fmt.Printf("Treasury Transaction Hash: %s\n", tspendOrPolicyKey)
	} else {
		fmt.Printf("Policy Key: %s\n", tspendOrPolicyKey)
	}
	fmt.Println()

	assignedTickets := make(map[string]bool)
	const repeatInterval = 5 * time.Minute

	for {
		newTickets := getNewTickets(assignedTickets)
		if len(newTickets) > 0 {
			for _, ticketHash := range newTickets {
				assignedTickets[ticketHash] = true
			}
		} else {
			time.Sleep(repeatInterval)
			continue
		}

		newTickets = shuffleStrings(newTickets)

		totalTickets := len(newTickets)
		var totalYes, totalNo, totalAbstain int

		fmt.Printf("Making random decisions based on  yes 0-%0.f  no %0.f-%0.f  abstain %0.f-100 for %d tickets\n", yesZone, yesZone,
			yesZone+noZone, yesZone+noZone, totalTickets)

		fmt.Println()

		// Calculate policies and keep track of counts in the first loop
		policyCounts := make(map[string]int)
		ticketPolicies := make(map[string]string)

		verbosePolicyTable := table.New("Ticket", "Rand", "Choice")
		for _, ticketHash := range newTickets {
			policy := calculatePolicy(verbosePolicyTable, ticketHash, salt, yesZone, noZone, verbose)
			ticketPolicies[ticketHash] = policy
			policyCounts[policy]++

			setTspendPolicy(tspendOrPolicyKey, ticketHash, policy)
		}

		if verbose {
			verbosePolicyTable.Print()
		}

		// Calculate percentages
		totalYes += policyCounts["yes"]
		totalNo += policyCounts["no"]
		totalAbstain += policyCounts["abstain"]

		summaryTable := table.New("", "yes", "no", "abs", "total")
		if totalTickets > 0 {
			yesPercent := float64(100*totalYes) / float64(totalTickets)
			noPercent := float64(100*totalNo) / float64(totalTickets)
			absPercentage := 100 - (yesPercent + noPercent)

			summaryTable.AddRow("votes", totalYes, totalNo, totalAbstain, totalTickets)
			summaryTable.AddRow("perc", formatPercentage(yesPercent)+"%", formatPercentage(noPercent)+"%",
				formatPercentage(absPercentage)+"%")
			summaryTable.AddRow("targ", formatPercentage(yesZone)+"%", formatPercentage(noZone)+"%",
				formatPercentage(absZone)+"%")
			summaryTable.AddRow(
				"diff",
				formatPercentage(math.Abs(yesPercent-yesZone))+"%",
				formatPercentage(math.Abs(noPercent-noZone))+"%",
				formatPercentage(math.Abs(absPercentage-absZone))+"%",
			)
		}

		if totalTickets > 0 {
			fmt.Println()
			summaryTable.Print()

			fmt.Println()
			fmt.Printf("Checking for tickets every %v\n", repeatInterval)
		}

		time.Sleep(repeatInterval)
	}
}

func getTickets() (*TicketsResponse, error) {
	cmd := exec.Command(dcrctl, "gettickets", "true")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var ticketsResponse TicketsResponse
	err = json.Unmarshal(output, &ticketsResponse)
	if err != nil {
		fmt.Println("Error parsing tickets JSON:", err)
		return nil, err
	}
	return &ticketsResponse, nil
}

func getNewTickets(assignedTickets map[string]bool) []string {
	ticketsResponse, err := getTickets()
	if err != nil {
		fmt.Println("Error fetching tickets:", err)
		return []string{}
	}

	var newTickets []string
	for _, ticketHash := range ticketsResponse.Hashes {
		if !assignedTickets[ticketHash] {
			newTickets = append(newTickets, ticketHash)
		}
	}
	return newTickets
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

func shuffleStrings(strings []string) []string {
	shuffled := make([]string, len(strings))
	copy(shuffled, strings)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}

func calculatePolicy(policyTable table.Table, ticketHash, salt string, yesZone, noZone float64, verbose bool) string {
	data := ticketHash + salt
	hashed := sha256.Sum256([]byte(data))
	seed := new(big.Int).SetBytes(hashed[:]).Uint64()
	r := rand.New(rand.NewSource(int64(seed)))

	// Generate a determinant number between 0 and 100
	determinant := r.Float64() * 100

	// Determine the policy based on random value
	var policy string
	if determinant <= yesZone {
		policy = "yes"
	} else if determinant <= yesZone+noZone {
		policy = "no"
	} else {
		policy = "abstain"
	}

	hashPrnt := ticketHash
	if !verbose {
		fmt.Print(formatPolicy(policy, verbose))
	}

	policyTable.AddRow(hashPrnt, formatPercentage(determinant), formatPolicy(policy, verbose))

	return policy
}

func setTspendPolicy(tspendOrPolicyKey, ticketHash, policy string) {
	if len(tspendOrPolicyKey) == 32 {
		dcrctlArgs = append(dcrctlArgs, "settspendpolicy", tspendOrPolicyKey, policy, ticketHash)
		cmd := exec.Command(dcrctl, dcrctlArgs...)
		cmd.Run()
	} else {
		dcrctlArgs = append(dcrctlArgs, "settreasurypolicy", tspendOrPolicyKey, policy, ticketHash)
		cmd := exec.Command(dcrctl, dcrctlArgs...)
		cmd.Run()
	}
}
