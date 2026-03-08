package iso20022

import "testing"

const validPain001 = `
<Document>
  <CstmrCdtTrfInitn>
    <PmtInf>
      <Dbtr>
        <Nm>Debtor GmbH</Nm>
        <PstlAdr>
          <TwnNm>Berlin</TwnNm>
          <Ctry>DE</Ctry>
          <PstCd>10115</PstCd>
          <AdrLine>Street 1</AdrLine>
        </PstlAdr>
      </Dbtr>
      <CdtTrfTxInf>
        <Amt><InstdAmt Ccy="EUR">100.00</InstdAmt></Amt>
        <CdtrAgt><FinInstnId><BICFI>DEUTDEFF</BICFI></FinInstnId></CdtrAgt>
        <Cdtr>
          <Nm>Creditor AG</Nm>
          <PstlAdr>
            <TwnNm>Paris</TwnNm>
            <Ctry>FR</Ctry>
            <PstCd>75001</PstCd>
            <AdrLine>Street 2</AdrLine>
          </PstlAdr>
        </Cdtr>
      </CdtTrfTxInf>
    </PmtInf>
  </CstmrCdtTrfInitn>
</Document>`

func TestParseAndValidate(t *testing.T) {
	doc, err := ParsePain001([]byte(validPain001))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if err := Validate(doc, ModeSCT); err != nil {
		t.Fatalf("SCT validate failed: %v", err)
	}
	if err := Validate(doc, ModeSCTInst); err != nil {
		t.Fatalf("SCT Inst validate failed: %v", err)
	}
}

func TestValidateFailsForSCTInstRules(t *testing.T) {
	doc, err := ParsePain001([]byte(validPain001))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	doc.CustomerTransfer.Payments[0].Txs[0].Amount.Value = 100001
	if err := Validate(doc, ModeSCTInst); err == nil {
		t.Fatalf("expected SCT Inst amount validation error")
	}
}
