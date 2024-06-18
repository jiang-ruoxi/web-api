package commands

import (
	"api/model"
	"github.com/jiang-ruoxi/gopkg/db"
	"github.com/urfave/cli/v2"
)

func Migrate() *cli.Command {

	return &cli.Command{
		Name:  "migrate",
		Usage: "数据库迁移",
		Subcommands: []*cli.Command{
			{
				Name:        "up",
				Usage:       "自动迁移数据库",
				Description: "自动迁移数据库",
				Action: func(ctx *cli.Context) error {
					tx := db.MustGet("web").Debug()
					tables := []interface{}{
						&model.Lexicon{},
					}
					return tx.AutoMigrate(tables...)
				},
			},
		},
	}
}
