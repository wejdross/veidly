package config

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Ctx struct {
	defines   map[string]string
	config    []byte
	configMap map[interface{}]interface{}
	Path      string
}

func (ctx *Ctx) Unmarshal(val interface{}) error {
	return yaml.Unmarshal(ctx.config, val)
}

func (ctx *Ctx) GetKey(key string) ([]byte, error) {
	v, e := ctx.configMap[key]
	if !e {
		return nil, fmt.Errorf("cannot find key: %s in configMap", key)
	}
	return yaml.Marshal(v)
}

func (ctx *Ctx) UnmarshalKey(
	key string, value interface{}, valFunc func() error,
) error {
	c, err := ctx.GetKey(key)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(c, value); err != nil {
		return err
	}

	if valFunc == nil {
		return nil
	}

	return valFunc()
}

func (ctx *Ctx) UnmarshalKeyPanic(key string, value interface{}, valFunc func() error) {
	if err := ctx.UnmarshalKey(key, value, valFunc); err != nil {
		panic(err)
	}
}

func Decrypt(val string, key []byte) (string, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	parts := strings.Split(val, ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid ecnrypted value")
	}

	nonce, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return "", err
	}
	etext, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	decrypted, err := gcm.Open(nil, nonce, etext, nil)
	return string(decrypted), err
}

func (ctx *Ctx) setDefines(definesFragment []byte, group, keypath string) error {

	var d struct {
		Defines map[string]interface{} `yaml:"defines"`
	}
	err := yaml.Unmarshal(definesFragment, &d)
	if err != nil {
		return err
	}

	var key []byte

	if keypath != "" {
		key, err = ioutil.ReadFile(keypath)
		if err != nil {
			fmt.Println("warn: no decryption key specified")
		}
		key, err = base64.StdEncoding.DecodeString(string(key))
		if err != nil {
			return err
		}
	}

	var defines = make(map[string]string)

	const encryptedFieldPrefix = "ENCRYPTED:"
	const extFieldPrefix = "EXT:"

	var extFiles = make(map[string]map[string]string)

	for k := range d.Defines {
		switch d.Defines[k].(type) {
		case string:
			defines[k] = d.Defines[k].(string)
		case map[interface{}]interface{}:
			v := d.Defines[k].(map[interface{}]interface{})
			groupVal, e := v[group]
			if !e {
				return fmt.Errorf("value for group: %s no found in key: %s", group, k)
			}
			gv := strings.TrimSpace(groupVal.(string))
			if strings.HasPrefix(gv, encryptedFieldPrefix) {
				gv = strings.Replace(gv, encryptedFieldPrefix, "", 1)
				if len(key) == 0 {
					return fmt.Errorf("%s requires decryption, but no key was provided", k)
				}
				defines[k], err = Decrypt(gv, key)
				if err != nil {
					return err
				}
			} else if strings.HasPrefix(gv, extFieldPrefix) {
				gv = strings.Replace(gv, extFieldPrefix, "", 1)
				parts := strings.Split(gv, ":")
				if len(parts) != 2 {
					return fmt.Errorf("invalid %s field value: %s", extFieldPrefix, gv)
				}
				extPath := parts[0]
				extField := parts[1]
				var extFc map[string]string

				if extFc = extFiles[extPath]; extFc == nil {
					c, err := ioutil.ReadFile(extPath)
					if err != nil {
						return err
					}
					if err := yaml.Unmarshal(c, &extFc); err != nil {
						return err
					}
					extFiles[extPath] = extFc
				}

				v := extFc[extField]
				if v == "" {
					return fmt.Errorf("secret value for %s not defined", extField)
				}

				defines[k] = v

			} else {
				defines[k] = gv
			}
		default:
			return fmt.Errorf("unrecognized type")
		}
	}

	replacerArr := make([]string, 0, len(defines)*2)

	for k := range defines {
		replacerArr = append(replacerArr, "${"+k+"}", defines[k])
	}

	replacer := strings.NewReplacer(replacerArr...)

	for k1 := range defines {
		defines[k1] = replacer.Replace(defines[k1])
	}

	ctx.defines = defines

	return nil
}

func preprocess(fragment []byte, defines map[string]string) []byte {

	for k := range defines {
		fragment = bytes.Replace(
			fragment,
			[]byte("${"+k+"}"),
			[]byte(defines[k]),
			-1)
	}

	return fragment
}

/*
	decode yaml segment at 0-based index {index}
	from file located at {path}, unmarshall it into {val}
	switch file values based on {group} decrypt values if needed using {key}
*/
func NewCtx(path string, ver, keypath string) *Ctx {

	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	ctx := new(Ctx)

	if err := ctx.setDefines(content, ver, keypath); err != nil {
		panic(err)
	}

	content = preprocess(content, ctx.defines)

	if err := yaml.Unmarshal(content, &ctx.configMap); err != nil {
		panic(err)
	}

	ctx.config = content
	ctx.Path = path

	return ctx
}

const DefaultConfigPath = "../config.yml"

func NewLocalCtx() *Ctx {
	return NewCtx(DefaultConfigPath, "local", "")
}
