package api

// // if noFixing is provided this function will fail without trying to redownload config
// func ValidateApiConfigVersion(path string, noFixing bool) error {
// 	var x struct {
// 		Version string
// 	}
// 	if err := helpers.YmlParseFile(path, &x); err != nil {
// 		return err
// 	}

// 	cvs, err := ioutil.ReadFile("api_conf_ver")
// 	if err != nil {
// 		return err
// 	}
// 	ver := string(cvs)

// 	if x.Version != ver {
// 		if noFixing {
// 			return fmt.Errorf(
// 				`Invalid config version.
// 			This revision expects config version to be %s, and you got: "%s"
// 			update your config by executing "make conf" `,
// 				ver, x.Version)
// 		}
// 		fmt.Printf("%vYour API config (api.yml) is out of date! Got: %s, expected: %s%v\n", cmn.ForeYellow, x.Version, ver, cmn.AttrOff)
// 		fmt.Printf("%sDo you want to download it? ( [y]es / [n]o / [i]gnore , default no)%s\n", cmn.ForeYellow, cmn.AttrOff)
// 		fmt.Printf("%sNote that your local api.yml will be overriden if you input 'y'%s\n", cmn.ForeYellow, cmn.AttrOff)
// 		res := ""
// 		fmt.Scanln(&res)
// 		switch res {
// 		case "y":
// 			c := exec.Command("make", "-C", "..", "conf", "ver="+ver)
// 			c.Stderr = os.Stderr
// 			c.Stdout = os.Stdout
// 			err := c.Run()
// 			if err != nil {
// 				return err
// 			}
// 			if err := helpers.YmlParseFile(path, &x); err != nil {
// 				return err
// 			}
// 			if x.Version != ver {
// 				return fmt.Errorf("%sKamil fucked up.\n He did not upload proper version of config to the server.\n Contact him to solve it\n%s", cmn.ForeRed, cmn.AttrOff)
// 			}
// 			return nil
// 		case "i":
// 			return nil
// 		default:
// 			return fmt.Errorf(
// 				`Invalid config version.
// 				This revision expects config version to be %s, and you got: "%s"
// 				update your config by executing "make conf" `,
// 				ver, x.Version)
// 		}
// 	}

// 	return nil
// }
