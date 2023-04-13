package static

/*
NOTE THAT CURRENTLY STATIC SERVER IS NOT SECURE
AND IS MEANT TO BE USED __ONLY__ IN DEV ENV
*/
import (
	"fmt"
	"net/url"
	"sport/api"
)

type Config struct {
	Enabled      bool   `yaml:"enabled"`
	Basepath     string `yaml:"basepath"`
	BaseUrl      string `yaml:"baseurl"`
	HttpBasePath string `yaml:"http_base_path"`
}

func (c *Config) Validate() error {
	const hdr = "Validate CtxRequest:"
	if c.Basepath == "" {
		return fmt.Errorf("%s invalid basepath", hdr)
	}
	if c.BaseUrl == "" {
		return fmt.Errorf("%s invalid baseurl", hdr)
	}
	if c.HttpBasePath == "" {
		return fmt.Errorf("%s invalid http_base_path", hdr)
	}
	if _, err := url.Parse(c.BaseUrl); err != nil {
		return fmt.Errorf("%s %v", hdr, err)
	}
	return nil
}

type Ctx struct {
	Config *Config
	Api    *api.Ctx
}

func NewCtx(a *api.Ctx) *Ctx {
	c := new(Ctx)
	c.Config = new(Config)
	a.Config.UnmarshalKeyPanic("static", c.Config, c.Config.Validate)
	c.Api = a
	c.RegisterAllHandlers()
	return c
}

func (ctx *Ctx) RegisterAllHandlers() {
	if !ctx.Config.Enabled {
		return
	}

	ctx.Api.EmptyGroup.GET(ctx.Config.HttpBasePath+"/:dir/:file", ctx.GetFile())
}

// func InitializeModule(configPath string, api *api.Api, enableStatic bool) {
// 	req, err := newCtxRequestFromConfig(configPath)
// 	if err != nil {
// 		panic(err)
// 	}
// 	c := new(Ctx)
// 	c.Config = req
// 	if _, err = url.Parse(req.BaseUrl); err != nil {
// 		panic(err)
// 	}
// 	if req.Enabled || enableStatic {
// 		// if err := createNecessaryDirectories(req); err != nil {
// 		// 	log.Fatal(err)
// 		// }
// 		api.EmptyGroup.GET(req.HttpBasePath+"/:dir/:file", GetFile)
// 	}
// }

// helper func to make sure directory structure is initialized
// func createNecessaryDirectories(req *CtxRequest) error {
// 	fd, err := os.Stat(req.Basepath)
// 	if err != nil {
// 		if len(req.Dirs_to_create) > 0 {
// 			for _, val := range req.Dirs_to_create {
// 				dirTODO := fmt.Sprintf("%s/%s", req.Basepath, val)
// 				if err := os.MkdirAll(dirTODO, 0700); err != nil {
// 					return fmt.Errorf("Problems while creating directory: %s,\nErrror is: %s", dirTODO, err.Error())
// 				}
// 			}
// 		}
// 	} else {
// 		if !fd.IsDir() {
// 			return fmt.Errorf("Path: %s is not directory!\nErrror is: %s", fd.Name(), err.Error())
// 		}
// 	}
// 	return nil
// }
