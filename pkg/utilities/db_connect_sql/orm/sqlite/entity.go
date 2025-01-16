package sqlite

import (
	"xorm.io/xorm"
)

type Service interface {
	Init() *xorm.Engine
}
