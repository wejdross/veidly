package dal

/*
	DAL stands for data access layer
	that is global service which is accessed by all modules
	willing to perform operations on vieewery database
*/

import (
	"database/sql"
	"fmt"
	"path"
	"sport/config"

	"github.com/kzaag/dp/cmn"
	"github.com/kzaag/dp/pgsql"
	"github.com/kzaag/dp/target"
	_ "github.com/kzaag/dp/target"
	_ "github.com/lib/pq"
)

type Ctx struct {
	Db *sql.DB
	Cs string
}

/*
	initialize global database context from target taken from dp-config
*/
func newDal(configCtx *config.Ctx, overrideDb string) (*Ctx, error) {
	var uargv *target.Args
	if overrideDb != "" {
		uargv = &target.Args{
			Set: map[string]string{
				"db": overrideDb,
			},
		}
	}

	dpConfig, err := configCtx.GetKey("pg")
	if err != nil {
		return nil, err
	}

	config, err := target.NewConfigFromBytes(dpConfig, configCtx.Path, uargv)
	if err != nil {
		return nil, err
	}
	const targetName = "ddl"
	var t *target.Target
	for _, _t := range config.Targets {
		if _t.Name == targetName {
			t = _t
			break
		}
	}
	if t == nil {
		return nil, fmt.Errorf("dp target '%s' not found", targetName)
	}

	conn, err := pgsql.TargetGetDB(t)
	if err != nil {
		return nil, err
	}

	cs, err := pgsql.TargetGetCS(t)
	if err != nil {
		return nil, err
	}

	ret := new(Ctx)
	ret.Db = conn.(*sql.DB)
	ret.Cs = cs

	return ret, err
}

func NewDal(configCtx *config.Ctx, overrideDb string) *Ctx {
	d, err := newDal(configCtx, overrideDb)
	if err != nil {
		panic(err)
	}
	return d
}

// completely reset and redeploy DAL based on config path
func deployDb(configCtx *config.Ctx, verbose bool, overrideDb string, reset bool) error {
	targetCtx := pgsql.TargetCtxNew()
	var targetConfig *target.Config
	var err error
	targetArgs := target.Args{
		ConfigPath: configCtx.Path,
		Verbose:    verbose,
		Execute:    true,
		Raw:        false,
		Demand:     target.MapFlags{},
		Set:        map[string]string{},
	}
	if overrideDb != "" {
		targetArgs.Set["db"] = overrideDb
	}
	if reset {
		targetArgs.Demand["reset"] = struct{}{}
	}
	b, err := configCtx.GetKey("pg")
	if err != nil {
		return err
	}

	if targetConfig, err = target.NewConfigFromBytes(b, targetArgs.ConfigPath, &targetArgs); err != nil {
		return err
	}
	return targetCtx.ExecConfig(targetConfig, &targetArgs)
}

// same as ResetDal but panics on error
func DeployDb(configCtx *config.Ctx, verbose bool, overrideDb string, reset bool) {
	err := deployDb(configCtx, verbose, overrideDb, reset)
	if err != nil {
		panic(err)
	}
}

func fixAbsPaths(basepath string, paths []string) {
	var i int
	var p *string
	for i = 0; i < len(paths); i++ {
		p = &paths[i]
		if !path.IsAbs(*p) {
			*p = path.Join(basepath, *p)
		}
	}
}

// verify that ddl is up to date
func ValidateNoPendingChanges(configCtx *config.Ctx, verbose bool, overrideDb string, noFix bool, forceFix bool) error {
	targetCtx := pgsql.TargetCtxNew()
	var targetConfig *target.Config
	var err error
	checkTargetName := "ddl"
	targetArgs := target.Args{
		ConfigPath: configCtx.Path,
		Verbose:    verbose,
		Execute:    false,
		Raw:        false,
		Demand: target.MapFlags{
			checkTargetName: struct{}{},
		},
		Set: map[string]string{},
	}
	if overrideDb != "" {
		targetArgs.Set["db"] = overrideDb
	}
	b, err := configCtx.GetKey("pg")
	if err != nil {
		return err
	}
	if targetConfig, err = target.NewConfigFromBytes(b, configCtx.Path, &targetArgs); err != nil {
		return err
	}
	var t *target.Target
	for i := range targetConfig.Targets {
		if targetConfig.Targets[i].Name == checkTargetName {
			t = targetConfig.Targets[i]
			break
		}
	}
	if t == nil {
		return fmt.Errorf("couldnt perform db verification, '%s' target is not found", checkTargetName)
	}
	// remove scipt execs
	t.Exec = t.Exec[0:1]
	fixAbsPaths(targetConfig.Base, t.Exec[0].Args)
	if len(t.Exec) != 1 {
		return fmt.Errorf("couldnt perform db verification, ambigious merge exec found")
	}
	db, err := targetCtx.DbNew(t)
	if err != nil {
		return err
	}
	s, err := targetCtx.GetMergeScript(targetConfig, db, t, t.Exec[0].Args)
	if err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	if noFix {
		return fmt.Errorf("ddl not up to date")
	}
	if forceFix {
		if err := deployDb(configCtx, true, overrideDb, false); err != nil {
			return err
		}
		s, err := targetCtx.GetMergeScript(targetConfig, db, t, t.Exec[0].Args)
		if err != nil {
			return err
		}
		if s == "" {
			fmt.Printf("%sYour DDL is up to date now!%s\n", cmn.ForeGreen, cmn.AttrOff)
			return nil
		} else {
			return fmt.Errorf("there are still pending changes:\n\n%s\n\nContact Kamil. This may be bug within dp", s)
		}
	} else {
		fmt.Printf("%sYour DDL is not up to date! Pending operations are:%s\n\n", cmn.ForeYellow, cmn.AttrOff)
		fmt.Println(s)
		fmt.Printf("%sDo you want to execute them? ( [y]es / [n]o / [i]gnore , default no)%s\n", cmn.ForeYellow, cmn.AttrOff)
		res := ""
		fmt.Scanln(&res)
		switch res {
		case "y":
			if err := deployDb(configCtx, true, overrideDb, false); err != nil {
				return err
			}
			s, err := targetCtx.GetMergeScript(targetConfig, db, t, t.Exec[0].Args)
			if err != nil {
				return err
			}
			if s == "" {
				fmt.Printf("%sYour DDL is up to date now!%s\n", cmn.ForeGreen, cmn.AttrOff)
				return nil
			} else {
				return fmt.Errorf("there are still pending changes:\n\n%s\n\nContact Kamil. This may be bug within dp", s)
			}
		case "i":
			return nil
		default:
			return fmt.Errorf("ddl not up to date")
		}
	}
}
