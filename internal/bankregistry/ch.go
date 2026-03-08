package bankregistry

// CHLoader implements CountryLoader for Switzerland BC-Nr data format.
type CHLoader struct {
	Schema CSVSchema
}

func NewCHLoader() *CHLoader {
	return &CHLoader{
		Schema: CSVSchema{Delimiter: '\t', HasHeader: true, BankCodeColumn: 0, BankNameColumn: 1, BICColumn: 2},
	}
}

func (l *CHLoader) Load(path string, registry *Registry, meta SourceMetadata) error {
	return registry.LoadCSV(path, "CH", l.Schema, meta)
}
