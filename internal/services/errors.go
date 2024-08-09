package services

import "fmt"

var (
	UserNotFound = fmt.Errorf("there is no user")
)
