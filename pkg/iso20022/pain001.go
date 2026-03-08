// Package iso20022 provides minimal, deterministic pain.001 parsing and validation.
package iso20022

import (
	"encoding/xml"
	"errors"
	"fmt"
)

// ValidationMode controls rule strictness.
type ValidationMode string

const (
	ModeSCT     ValidationMode = "SCT"
	ModeSCTInst ValidationMode = "SCT_INST"
)

// Document models the subset of pain.001 used by this library.
type Document struct {
	XMLName          xml.Name          `xml:"Document"`
	CustomerTransfer CustomerTransfer  `xml:"CstmrCdtTrfInitn"`
}

type CustomerTransfer struct {
	Payments []PaymentInfo `xml:"PmtInf"`
}

type PaymentInfo struct {
	Debtor Debtor `xml:"Dbtr"`
	Txs    []Tx   `xml:"CdtTrfTxInf"`
}

type Debtor struct {
	Name    string      `xml:"Nm"`
	Address PostalAddress `xml:"PstlAdr"`
}

type Tx struct {
	Amount Amount `xml:"Amt>InstdAmt"`
	Agent  Agent  `xml:"CdtrAgt"`
	Creditor Creditor `xml:"Cdtr"`
}

type Amount struct {
	Currency string  `xml:"Ccy,attr"`
	Value    float64 `xml:",chardata"`
}

type Agent struct {
	BIC string `xml:"FinInstnId>BICFI"`
}

type Creditor struct {
	Name    string       `xml:"Nm"`
	Address PostalAddress `xml:"PstlAdr"`
}

type PostalAddress struct {
	TownName  string   `xml:"TwnNm"`
	Country   string   `xml:"Ctry"`
	PostCode  string   `xml:"PstCd"`
	AddrLines []string `xml:"AdrLine"`
}

func ParsePain001(data []byte) (*Document, error) {
	var doc Document
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	if len(doc.CustomerTransfer.Payments) == 0 {
		return nil, errors.New("pain.001 contains no PmtInf blocks")
	}
	return &doc, nil
}

func Validate(doc *Document, mode ValidationMode) error {
	if doc == nil {
		return errors.New("document is nil")
	}
	if len(doc.CustomerTransfer.Payments) == 0 {
		return errors.New("no payment info blocks")
	}

	for pi, p := range doc.CustomerTransfer.Payments {
		if err := validateStructuredAddress(p.Debtor.Address, fmt.Sprintf("PmtInf[%d].Dbtr", pi)); err != nil {
			return err
		}
		if len(p.Txs) == 0 {
			return fmt.Errorf("PmtInf[%d] has no transactions", pi)
		}
		for ti, tx := range p.Txs {
			if err := validateStructuredAddress(tx.Creditor.Address, fmt.Sprintf("PmtInf[%d].Tx[%d].Cdtr", pi, ti)); err != nil {
				return err
			}
			if mode == ModeSCTInst {
				if tx.Amount.Value > 100000 {
					return fmt.Errorf("PmtInf[%d].Tx[%d]: SCT Inst amount must be <= 100000", pi, ti)
				}
				if tx.Agent.BIC == "" {
					return fmt.Errorf("PmtInf[%d].Tx[%d]: SCT Inst requires creditor BIC", pi, ti)
				}
				if tx.Amount.Currency != "EUR" {
					return fmt.Errorf("PmtInf[%d].Tx[%d]: SCT Inst currency must be EUR", pi, ti)
				}
			}
		}
	}

	return nil
}

func validateStructuredAddress(addr PostalAddress, owner string) error {
	if addr.TownName == "" || addr.Country == "" || addr.PostCode == "" {
		return fmt.Errorf("%s: structured address requires TwnNm/Ctry/PstCd", owner)
	}
	if len(addr.AddrLines) > 2 {
		return fmt.Errorf("%s: structured address allows max 2 AdrLine entries", owner)
	}
	return nil
}
