package bankregistry

// ESLoader implements CountryLoader for Spain CCC data format.
type ESLoader struct {
	Schema CSVSchema
}

func NewESLoader() *ESLoader {
	return &ESLoader{
		Schema: CSVSchema{Delimiter: ';', HasHeader: true, BankCodeColumn: 0, BankNameColumn: 1, BICColumn: 2},
	}
}

func (l *ESLoader) Load(path string, registry *Registry, meta SourceMetadata) error {
	return registry.LoadCSV(path, "ES", l.Schema, meta)
}
