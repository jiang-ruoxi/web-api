package model

import (
	"github.com/jiang-ruoxi/gopkg/db"
	"gorm.io/gorm"
)

func DefaultWeb() *gorm.DB { return db.MustGet("web") }

func DefaultMarket() *gorm.DB { return db.MustGet("market") }
