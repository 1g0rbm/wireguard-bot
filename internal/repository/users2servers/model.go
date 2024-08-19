package users2servers

type Users2Servers struct {
	UserId   int64  `db:"user_id"`
	ServerId int    `db:"server_id"`
	Address  string `db:"address"`
}
