package bankregistry

// ITLoader implements CountryLoader for Italy ABI mapping data format.
type ITLoader struct {
	Schema CSVSchema
}

func NewITLoader() *ITLoader {
	return &ITLoader{
		Schema: CSVSchema{Delimiter: ',', HasHeader: true, BankCodeColumn: 1, BankNameColumn: 2, BICColumn: 0},
	}
}

func (l *ITLoader) Load(path string, registry *Registry, meta SourceMetadata) error {
	return registry.LoadCSV(path, "IT", l.Schema, meta)
}
