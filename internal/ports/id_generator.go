package ports

// IDGenerator abstracts ID generation to make services testable.
type IDGenerator interface {
	NewID() string
}
