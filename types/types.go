package types

import (
	"database/sql"

	"github.com/Moonyongjung/xpla.go/client"
	"github.com/spf13/viper"
)

const (
	GenesisBlockNum = "1"
)

var (
	Db           *sql.DB
	ErrLogTable  = "errLog"
	InfoLogTable = "infoLog"
)

type App struct {
	Viper       *viper.Viper
	PubClient   *client.XplaClient
	PrivClient  *client.XplaClient
	Channels    Channels
	HomePath    string
	AppFilePath string
}
