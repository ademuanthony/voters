#!/bin/bash

# This Bash script sets the tspend policy for owned tickets in Decred based on specified Yes, No, and Abstain percentages.

# Check if the correct number of arguments is provided
if [ "$#" -ne 4 ]; then
    echo "Usage: $0 <Salt> <Yes Percentage> <No Percentage> <Abstain Percentage>"
    exit 1
fi

# Extract Salt, Yes, No, and Abstain percentages from command-line arguments
salt=$1
yes_percentage=$2
no_percentage=$3
abstain_percentage=$4

# Get a list of owned tickets and their hashes from the dcrctl.sh script
tickets=$(./dcrctl.sh gettickets)

shuffled_tickets=$(echo "$tickets" | jq -r '.hashes | .[]' | sort --random-sort)
echo "${#shuffled_tickets[@]}"

# Calculate the number of tickets for each policy
total_tickets=${#shuffled_tickets[@]}
scaling_factor=10000  # Use 10000 as the scaling factor
yes_count=$((total_tickets * yes_percentage * scaling_factor / 100))
no_count=$((total_tickets * no_percentage * scaling_factor / 100))
abstain_count=$((total_tickets * scaling_factor - yes_count - no_count))

# Cycle through shuffled owned tickets and set tspend policy for each ticket based on the calculated counts
for ticket_hash in "${shuffled_tickets[@]}"; do
    # Calculate the policy for this ticket based on the seed and percentages
    seed=$(echo -n "$ticket_hash$salt" | openssl dgst -sha256 | cut -d ' ' -f2)
    
    # Convert the seed to an integer between 0 and scaling_factor
    seed_integer=$((0x$seed))  # Convert from hexadecimal
    seed_normalized=$((seed_integer % scaling_factor))
    
    # Calculate the policy based on the seed and percentages
    if [ "$seed_normalized" -le "$yes_count" ]; then
        policy="yes"
    elif [ "$seed_normalized" -le "$((yes_count + no_count))" ]; then
        policy="no"
    else
        policy="abstain"
    fi

    # Set the tspend policy for the current ticket using the dcrctl.sh script
    ./dcrctl.sh settspendpolicy "$ticket_hash" "$policy"

    echo "Treasury spend policy set for ticket $ticket_hash: $policy"
    echo "___________________________________________________________"
done
