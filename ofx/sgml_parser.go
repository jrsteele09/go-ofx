package ofx

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/civil"
)

// ParseSGML handles SGML-based OFX parsing
func ParseSGML(r io.Reader) ([]Transaction, error) {
	var transactions []Transaction
	var currentTransaction Transaction
	var balance int64

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Parse transaction type
		if strings.Contains(line, "<TRNTYPE>") {
			currentTransaction.Type = extractTagValue(line, "TRNTYPE")
		}

		// Parse date posted
		if strings.Contains(line, "<DTPOSTED>") {
			dateStr := extractTagValue(line, "DTPOSTED")
			if date, err := time.Parse("20060102", dateStr[:8]); err == nil {
				currentTransaction.Date = civil.DateOf(date)
			} else {
				return nil, fmt.Errorf("failed to parse date: %v", err)
			}
		}

		// Parse transaction amount
		if strings.Contains(line, "<TRNAMT>") {
			amountStr := extractTagValue(line, "TRNAMT")
			amount, err := strconv.ParseFloat(amountStr, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse transaction amount: %v", err)
			}

			if amount < 0 {
				currentTransaction.Debit = int64(-amount * 100) // convert to cents
			} else {
				currentTransaction.Credit = int64(amount * 100) // convert to cents
			}
			balance += int64(amount * 100)
			currentTransaction.Balance = balance
		}

		// Parse transaction ID
		if strings.Contains(line, "<FITID>") {
			currentTransaction.ID = extractTagValue(line, "FITID")
		}

		// Parse memo/description
		if strings.Contains(line, "<MEMO>") {
			currentTransaction.Description = extractTagValue(line, "MEMO")
		}

		// End of transaction
		if strings.Contains(line, "</STMTTRN>") {
			transactions = append(transactions, currentTransaction)
			currentTransaction = Transaction{} // Reset for next transaction
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading SGML OFX: %v", err)
	}

	return transactions, nil
}

// Helper function to extract the value between tags
func extractTagValue(line, tag string) string {
	openTag := fmt.Sprintf("<%s>", tag)
	closeTag := fmt.Sprintf("</%s>", tag)

	start := strings.Index(line, openTag) + len(openTag)
	end := strings.Index(line, closeTag)

	if start >= 0 && end >= 0 && start < end {
		return line[start:end]
	}
	return ""
}
