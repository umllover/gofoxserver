package msg

import "time"

type AuthInfo struct {
	RetCode      int //非0位失败
	UserID       int //
	SpreaderID   int //推广人id
	RegisterDate *time.Time
}
