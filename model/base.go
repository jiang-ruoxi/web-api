package model

import (
	"github.com/jiang-ruoxi/gopkg/db"
	"gorm.io/gorm"
)

// Default 默认数据库
func Default() *gorm.DB { return db.MustGet("web") }
