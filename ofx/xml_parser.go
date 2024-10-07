package ofx

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"time"

	"cloud.google.com/go/civil"
)

// OFX is the top-level structure for parsing the XML OFX file (same as before)
type OFX struct {
	BankMsgs BankMsgs `xml:"BANKMSGSRSV1"`
}

type BankMsgs struct {
	StatementResponse StatementResponse `xml:"STMTTRNRS"`
}

type StatementResponse struct {
	Statement Statement `xml:"STMTRS"`
}

type Statement struct {
	TransactionList TransactionList `xml:"BANKTRANLIST"`
}

type TransactionList struct {
	Transactions []OFXTransaction `xml:"STMTTRN"`
}

type OFXTransaction struct {
	Type        string `xml:"TRNTYPE"`
	DatePosted  string `xml:"DTPOSTED"`
	Amount      string `xml:"TRNAMT"`
	ID          string `xml:"FITID"`
	Description string `xml:"MEMO"`
	CheckNum    string `xml:"CHECKNUM,omitempty"`
}

// ParseXML handles XML-based OFX parsing
func ParseXML(r io.Reader) ([]Transaction, error) {
	var ofx OFX
	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(&ofx); err != nil {
		return nil, fmt.Errorf("error decoding XML OFX: %v", err)
	}
	return ConvertOFXToTransactions(ofx)
}

// ConvertOFXToTransactions converts parsed OFX data to Transaction struct
func ConvertOFXToTransactions(ofx OFX) ([]Transaction, error) {
	var transactions []Transaction
	for _, ofxTran := range ofx.BankMsgs.StatementResponse.Statement.TransactionList.Transactions {
		transaction := Transaction{
			ID:          ofxTran.ID,
			Type:        ofxTran.Type,
			Description: ofxTran.Description,
		}

		if date, err := time.Parse("20060102", ofxTran.DatePosted[:8]); err == nil {
			transaction.Date = civil.DateOf(date)
		} else {
			return nil, fmt.Errorf("failed to parse date: %v", err)
		}

		amount, err := strconv.ParseFloat(ofxTran.Amount, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse amount: %v", err)
		}

		if amount < 0 {
			transaction.Debit = int64(-amount * 100)
		} else {
			transaction.Credit = int64(amount * 100)
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
