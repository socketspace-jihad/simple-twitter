package db

type DB interface {
	Connect() error
	Disconnect() error
}
