package dataset

// DataSource defines a generic interface for loading and decoding data assets.
type DataSource interface {
	Load() error
	Get() interface{}
}
