package bankregistry

// PLLoader implements CountryLoader for Poland NRB data format.
type PLLoader struct {
	Schema CSVSchema
}

func NewPLLoader() *PLLoader {
	return &PLLoader{
		Schema: CSVSchema{Delimiter: '|', HasHeader: false, BankCodeColumn: 0, BankNameColumn: 2, BICColumn: 1},
	}
}

func (l *PLLoader) Load(path string, registry *Registry, meta SourceMetadata) error {
	return registry.LoadCSV(path, "PL", l.Schema, meta)
}
