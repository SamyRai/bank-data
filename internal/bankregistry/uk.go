package bankregistry

// UKLoader implements CountryLoader for UK Sort Code data format.
type UKLoader struct {
	Schema CSVSchema
}

func NewUKLoader() *UKLoader {
	return &UKLoader{
		Schema: CSVSchema{Delimiter: ',', HasHeader: true, BankCodeColumn: 0, BankNameColumn: 1, BICColumn: 2},
	}
}

func (l *UKLoader) Load(path string, registry *Registry, meta SourceMetadata) error {
	// The open banking directory data format is mapped via UKLoader for GB country code
	return registry.LoadCSV(path, "GB", l.Schema, meta)
}
