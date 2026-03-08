package bankregistry

// SELoader implements CountryLoader for Sweden Clearing number data format.
type SELoader struct {
	Schema CSVSchema
}

func NewSELoader() *SELoader {
	return &SELoader{
		Schema: CSVSchema{Delimiter: ',', HasHeader: true, BankCodeColumn: 0, BankNameColumn: 2, BICColumn: 1},
	}
}

func (l *SELoader) Load(path string, registry *Registry, meta SourceMetadata) error {
	return registry.LoadCSV(path, "SE", l.Schema, meta)
}
