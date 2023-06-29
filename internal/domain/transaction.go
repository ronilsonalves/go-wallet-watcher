package domain

type Transaction struct {
	BlockNumber       string `json:"blockNumber,omitempty"`
	TimeStamp         string `json:"timeStamp,omitempty"`
	Hash              string `json:"hash,omitempty"`
	Nonce             string `json:"nonce,omitempty"`
	BlockHash         string `json:"blockHash,omitempty"`
	TransactionIndex  string `json:"transactionIndex,omitempty"`
	From              string `json:"from,omitempty"`
	To                string `json:"to,omitempty"`
	Value             string `json:"value,omitempty"`
	Gas               string `json:"gas,omitempty"`
	GasPrice          string `json:"gasPrice,omitempty"`
	IsError           string `json:"isError,omitempty"`
	TxreceiptStatus   string `json:"txreceipt_status,omitempty"`
	Input             string `json:"input,omitempty"`
	ContractAddress   string `json:"contractAddress,omitempty"`
	CumulativeGasUsed string `json:"cumulativeGasUsed,omitempty"`
	GasUsed           string `json:"gasUsed,omitempty"`
	Confirmations     string `json:"confirmations,omitempty"`
	MethodId          string `json:"methodId,omitempty"`
	FunctionName      string `json:"functionName,omitempty"`
}
