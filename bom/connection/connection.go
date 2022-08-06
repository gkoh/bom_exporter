package connection

// Retriever is the interface that wraps a data connection.
//
// Identifier returns the connection identity.
// Retrieve obtains the data from the underlying connection.
type Retriever interface {
	Identifier() string
	Retrieve() ([]byte, error)
}
