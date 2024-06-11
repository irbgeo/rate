package storage

type StartOpts struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

type token struct {
	Token string `db:"token"`
}
