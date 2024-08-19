package users2servers

type Users2Servers struct {
	UserID   int64  `db:"user_id"`
	ServerID int    `db:"server_id"`
	Address  string `db:"address"`
}
