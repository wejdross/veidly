package main

import (
	"sport/api"
	"sport/instr"
	"sport/rsv"
	"sport/train"
	"sport/user"
)

type Ctx struct {
	Rsv   *rsv.Ctx
	User  *user.Ctx
	Instr *instr.Ctx
	Api   *api.Ctx
	Train *train.Ctx
}
