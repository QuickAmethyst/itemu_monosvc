package sql

import "fmt"

func BuildPostgresUri(conf ConnectionOptions) string {
	ssl := `disable`

	if conf.SSL {
		ssl = `require`
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		conf.Host,
		conf.Port,
		conf.User,
		conf.Password,
		conf.DatabaseName,
		ssl,
	)
}
