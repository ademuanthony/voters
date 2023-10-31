## gettickets: `gettickets includeimmature`

Returning the hashes of the tickets currently owned by wallet.

Arguments:
1. includeimmature (boolean, required) If true include immature tickets in the results.

Result:
{
   "hashes": ["value",...], (array of string) Hashes of the tickets owned by the wallet encoded as strings
}      


## ticketsforaddress: `ticketsforaddress "address"`

Request all the tickets for an address.

Arguments:
1. address (string, required) Address to look for.

Result:
true|false (boolean) Tickets owned by the specified address.


## ticketinfo: `ticketinfo (startheight=0)`

Returns details of each wallet ticket transaction

Arguments:
1. startheight (numeric, optional, default=0) Specify the starting block height to scan from

Result:
[{
  "hash": "value",               (string)          Transaction hash of the ticket
 "cost": n.nnn,                 (numeric)         Amount paid to purchase the ticket; this may be greater than the ticket price at time of purchase
 "votingaddress": "value",      (string)          Address of 0th output, which describes the requirements to spend the ticket
 "status": "value",             (string)          Description of ticket status (unknown, unmined, immature, mature, live, voted, missed, expired, unspent, revoked)
 "blockhash": "value",          (string)          Hash of block ticket is mined in
 "blockheight": n,              (numeric)         Height of block ticket is mined in
 "vote": "value",               (string)          Transaction hash of vote which spends the ticket
 "revocation": "value",         (string)          Transaction hash of revocation which spends the ticket
 "choices": [{                  (array of object) Vote preferences set for the ticket
  "agendaid": "value",          (string)          The ID for the agenda the choice concerns
  "agendadescription": "value", (string)          A description of the agenda the choice concerns
  "choiceid": "value",          (string)          The ID of the current choice for this agenda
  "choicedescription": "value", (string)          A description of the current choice for this agenda
 },...],                                          
 "vsphost": "value",            (string)          VSP Host associated with the ticket (if any)
},...]


## setvotechoice: `setvotechoice "agendaid" "choiceid" ("tickethash")`

Sets choices for defined agendas in the latest stake version supported by this software

Arguments:
1. agendaid   (string, required) The ID for the agenda to modify
2. choiceid   (string, required) The ID for the choice to choose
3. tickethash (string, optional) The hash of the ticket to set choices for

Result:
Nothing


## settspendpolicy: `settspendpolicy "hash" "policy" ("ticket")`

Set a voting policy for a treasury spend transaction

Arguments:
1. hash   (string, required) Hash of treasury spend transaction to set policy for
2. policy (string, required) Voting policy for a tspend transaction (invalid/abstain, yes, or no)
3. ticket (string, optional) Ticket hash to set a per-ticket tspend approval policy

Result:
Nothing,


## settreasurypolicy:  `settreasurypolicy "key" "policy" ("ticket")`

Set a voting policy for treasury spends by a particular key

Arguments:
1. key    (string, required) Treasury key to set policy for
2. policy (string, required) Voting policy for a treasury key (invalid/abstain, yes, or no)
3. ticket (string, optional) Ticket hash to set a per-ticket treasury key policy

Result:
Nothing,


## setdisapprovepercent: `setdisapprovepercent percent`

Sets the wallet's block disapprove percent per vote. The wallet will randomly disapprove blocks with this percent of votes. Only used for testing purposes and will fail on mainnet.
Arguments:
1. percent (numeric, required) The percent of votes to disapprove blocks. i.e. 100 means that all votes disapprove the block they are called on. Must be between zero and one hundred.

Result:
nNothing

## disapprovepercent: `disapprovepercent`
Returns the wallet's current block disapprove percent per vote. i.e. 100 means that all votes disapprove the block they are called on. Only used for testing purposes.

Arguments:
None

Result:
n (numeric) The disapprove percent. When voting, this percent of votes will randomly disapprove the block they are called on


sendtotreasury: `sendtotreasury amount`

Send decred to treasury

Arguments:
1. amount (numeric, required) Amount to send to treasury

Result:
"value" (string) The transaction hash of the sent transaction,


## sendfromtreasury: `sendfromtreasury "key" amounts`

Send from treasury balance to multiple recipients.

Arguments:
1. key     (string, required) Politeia public key
2. amounts (object, required) Pairs of payment addresses and the output amount to pay each
{
  "Address to pay": Amount to send to the payment address valued in decred, (object) JSON object using payment addresses as keys and output amounts valued in decred to send to each address
...
}

Result:
"value" (string) The transaction hash of the sent transaction,



## treasurypolicy: `treasurypolicy ("key" "ticket")`

Return voting policies for treasury spend transactions by key

Arguments:
1. key    (string, optional) Return the policy for a particular key
2. ticket (string, optional) Return policies used by a specific ticket hash

Result (no key provided):
[{
 "key": "value",    (string) Treasury key associated with a policy
 "policy": "value", (string) Voting policy description (abstain, yes, or no)
 "ticket": "value", (string) Ticket hash of a per-ticket treasury key approval policy
},...]

Result (key specified):
{
 "key": "value",    (string) Treasury key associated with a policy
 "policy": "value", (string) Voting policy description (abstain, yes, or no)
 "ticket": "value", (string) Ticket hash of a per-ticket treasury key approval policy
}

## tspendpolicy: `tspendpolicy ("hash" "ticket")`

Return voting policies for treasury spend transactions

Arguments:
1. hash   (string, optional) Return the policy for a particular tspend hash
2. ticket (string, optional) Return policies used by a specific ticket hash

Result (no tspend hash provided):
[{
   "hash": "value",   (string) Treasury spend transaction hash
 "policy": "value", (string) Voting policy description (abstain, yes, or no)
 "ticket": "value", (string) Ticket hash of a per-ticket tspend approval policy
},...]

Result (tspend hash specified):
{
   "hash": "value",   (string) Treasury spend transaction hash
 "policy": "value", (string) Voting policy description (abstain, yes, or no)
 "ticket": "value", (string) Ticket hash of a per-ticket tspend approval policy
}


## getvotechoices: `getvotechoices ("tickethash")`

Retrieve the currently configured default vote choices for the latest supported stake agendas

Arguments:
1. tickethash (string, optional) The hash of the ticket to return vote choices for. If the ticket has no choices set, the default vote choices are returned

Result:
{
   "version": n,                  (numeric)         The latest stake version supported by the software and the version of the included agendas
 "choices": [{                  (array of object) The currently configured agenda vote choices, including abstaining votes
  "agendaid": "value",          (string)          The ID for the agenda the choice concerns
  "agendadescription": "value", (string)          A description of the agenda the choice concerns
  "choiceid": "value",          (string)          The ID of the current choice for this agenda
  "choicedescription": "value", (string)          A description of the current choice for this agenda
 },...],                                          
}


## purchaseticket: `purchaseticket "fromaccount" spendlimit (minconf=1 "ticketaddress" numtickets=1 "pooladdress" poolfees expiry "comment" dontsigntx)`

Purchase ticket using available funds.

Arguments:
1.  fromaccount   (string, required)             The account to use for purchase (default="default")
2.  spendlimit    (numeric, required)            Limit on the amount to spend on ticket
3.  minconf       (numeric, optional, default=1) Minimum number of block confirmations required
4.  ticketaddress (string, optional)             Override the ticket address to which voting rights are given
5.  numtickets    (numeric, optional, default=1) The number of tickets to purchase
6.  pooladdress   (string, optional)             The address to pay stake pool fees to
7.  poolfees      (numeric, optional)            The amount of fees to pay to the stake pool
8.  expiry        (numeric, optional)            Height at which the purchase tickets expire
9.  comment       (string, optional)             Unused
10. dontsigntx    (boolean, optional)            Return unsigned split and ticket transactions instead of signing and publishing

Result:
"value" (string) Hash of the resulting ticket


## processunmanagedticket: `processunmanagedticket ("tickethash")`

Processes tickets for vsp client based on ticket hash.

Arguments:
1. tickethash (string, optional) The ticket hash of ticket to be processed by the vsp client.

Result:
Nothing

