package mongodb

import (
	"errors"

	"github.com/spf13/cobra"
)

var ErrCredentialsNotFound = errors.New("error mongodb credentials not found")

const (
	hostFlag     = "db-host"
	portFlag     = "db-port"
	passwordFlag = "db-password"
	userFlag     = "db-user"
	dbFlag       = "db-db"
)

func SetFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String(hostFlag, "", "mongodb host")
	cmd.PersistentFlags().String(portFlag, "", "mongodb port")
	cmd.PersistentFlags().String(passwordFlag, "", "mongodb password")
	cmd.PersistentFlags().String(userFlag, "", "mongodb user")
	cmd.PersistentFlags().String(dbFlag, "", "mongodb db")
}
