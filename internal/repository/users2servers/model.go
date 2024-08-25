package users2servers

type Users2Servers struct {
	UserID   int64  `db:"user_id"`
	ServerID int    `db:"server_id"`
	Address  string `db:"address"`
}

type UsersServers struct {
	UserID           int64  `db:"user_id"`
	Username         string `db:"username"`
	FirstName        string `db:"first_name"`
	LastName         string `db:"last_name"`
	Role             int8   `db:"role"`
	UserPublicKey    string `db:"user_public_key"`
	UserPrivateKey   string `db:"user_private_key"`
	State            string `db:"state"`
	UserAddress      string `db:"user_address"`
	ServerID         int    `db:"server_id"`
	ServerName       string `db:"server_name"`
	ServerAddress    string `db:"server_address"`
	ServerPublicKey  string `db:"server_public_key"`
	ServerPrivateKey string `db:"server_private_key"`
}
