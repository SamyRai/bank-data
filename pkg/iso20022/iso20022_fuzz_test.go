package iso20022

import "testing"

func FuzzPain001ParseValidate(f *testing.F) {
	f.Add(`<Document><CstmrCdtTrfInitn><PmtInf><Dbtr><Nm>A</Nm><PstlAdr><TwnNm>Berlin</TwnNm><Ctry>DE</Ctry><PstCd>10115</PstCd><AdrLine>Line1</AdrLine></PstlAdr></Dbtr><CdtTrfTxInf><Amt><InstdAmt Ccy="EUR">10</InstdAmt></Amt><CdtrAgt><FinInstnId><BICFI>DEUTDEFF</BICFI></FinInstnId></CdtrAgt><Cdtr><Nm>B</Nm><PstlAdr><TwnNm>Paris</TwnNm><Ctry>FR</Ctry><PstCd>75001</PstCd><AdrLine>Line2</AdrLine></PstlAdr></Cdtr></CdtTrfTxInf></PmtInf></CstmrCdtTrfInitn></Document>`)
	f.Add("<Document></Document>")
	f.Add("not xml")

	f.Fuzz(func(_ *testing.T, input string) {
		doc, err := ParsePain001([]byte(input))
		if err != nil {
			return
		}
		_ = Validate(doc, ModeSCT)
		_ = Validate(doc, ModeSCTInst)
	})
}
