package ofx

import (
	"bytes"
	"fmt"
	"io"

	"cloud.google.com/go/civil"
)

type OFXType string

const (
	XMLType  OFXType = "XML"
	SGMLType OFXType = "SGML"
)

type Transaction struct {
	ID             string     `json:"id"`
	ImportIDs      []string   `json:"importIds"`
	Type           string     `json:"type"`
	Description    string     `json:"description"`
	CurrencySymbol string     `json:"symbol"`
	Debit          int64      `json:"debit"`
	Credit         int64      `json:"credit"`
	Balance        int64      `json:"balance"`
	Date           civil.Date `json:"date"`
}

// detectFormat detects whether the OFX file is XML or SGML
func detectFormat(r io.Reader) (OFXType, error) {
	// Read a few bytes to detect if it's XML (has XML declaration)
	buffer := make([]byte, 100)
	_, err := r.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("could not read input: %v", err)
	}

	if bytes.Contains(buffer, []byte("<?xml")) {
		return XMLType, nil
	}

	return SGMLType, nil
}

// Parse is the main function that selects the appropriate parser
func Parse(r io.Reader) ([]Transaction, error) {
	allBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %v", err)
	}

	format := SGMLType
	if bytes.Contains(allBytes, []byte("<?xml")) {
		format = XMLType
	}

	buf := bytes.NewReader(allBytes)
	switch format {
	case XMLType:
		return ParseXML(buf)
	case SGMLType:
		return ParseSGML(buf)
	default:
		return nil, fmt.Errorf("unknown format")
	}
}
