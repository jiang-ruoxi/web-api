package commands

import (
	"api/http/router"

	"github.com/jiang-ruoxi/gopkg/graceful"
	"github.com/jiang-ruoxi/gopkg/server"
	"github.com/urfave/cli/v2"
)

func Serve(c *cli.Context) error {

	// 运行HTTP服务
	graceful.Start(server.NewHttp(server.Addr(":8000"), server.Router(router.All())))

	graceful.Wait()

	return nil
}
