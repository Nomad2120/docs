package crypto

// CryptoAPI -
type CryptoAPI interface {
	LoadKeyStore(storage int, container, password, alias string) error
	SignData(alias string, flags int, inData string, inSign string) ([]byte, error)
	ExtractCertFromCMS(cms string, signID int, flags int) ([]byte, error)
	SignWSSE(alias string, flags int, inData, signNodeID string) ([]byte, error)
	Close() error
}
