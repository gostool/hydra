package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/vars/rlog"
)

func init() {
	hydra.OnReady(func() {
		hydra.Conf.RPC(":8092")
		hydra.Conf.Vars().RLog("/rpc/log", rlog.WithDisable())
	})
}
