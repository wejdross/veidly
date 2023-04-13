package helpers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func WrapErr(prefix string, err error) error {
	return fmt.Errorf("%s: %s", prefix, err)
}

type ElementNotFoundErr struct {
	subject string
}

func (e ElementNotFoundErr) Error() string {
	return fmt.Sprintf("%s not found.", e.subject)
}

/*
	Note for maintainers (and future me):
	ENF (aka element not found) error used to be separate error
	but i noticed that this is too complex and unneccessary
	so i decided to use it as synonym to sql.NoRows error.

	for future reference:
	you shouldnt use those wrappers. Just use plain sql.NoRows in your code
*/

func NewElementNotFoundErr(subj interface{}) error {
	// return ElementNotFoundErr{
	// 	subject: fmt.Sprintf("%v", subj),
	// }
	return sql.ErrNoRows
}

// checks if error is "element not found"
func IsENF(err error) bool {
	return err == sql.ErrNoRows
	// _, ok := err.(ElementNotFoundErr)
	// return ok
}

type ModuleDisabledError struct {
	module string
}

func (m ModuleDisabledError) Error() string {
	return fmt.Sprintf("Module %s is disabled", m.module)
}

func NewModuleDisabledError(moduleName string) error {
	return ModuleDisabledError{module: moduleName}
}

type MuxError struct {
	title  string
	errors []error
}

func (m *MuxError) HasErrors() bool {
	return len(m.errors) != 0
}

func NewMuxError(title string) MuxError {
	ret := MuxError{}
	ret.title = title
	ret.errors = make([]error, 0, 10)
	return ret
}

func (m *MuxError) Add(err error) {
	m.errors = append(m.errors, err)
}

func (m MuxError) Error() string {
	ret := m.title + "\n"
	for i, err := range m.errors {
		ret += "\t" + strconv.Itoa(i+1) + ": " + err.Error() + "\n"
	}
	return ret
}

type HttpError struct {
	err    error
	retMsg string
	code   int
}

func (e *HttpError) Error() string {
	return e.err.Error()
}

func NewHttpError(status int, publicMsg string, err error) error {
	if err == nil {
		if publicMsg == "" {
			panic("err cant be nil")
		}
		err = errors.New(publicMsg)
	}
	return &HttpError{
		err:    err,
		retMsg: publicMsg,
		code:   status,
	}
}

type HttpErrorBody struct {
	Err string
}

// wont write if e.msg is empty
func (e HttpError) Write(w io.Writer) error {
	if e.retMsg == "" {
		return nil
	}
	var b HttpErrorBody
	b.Err = e.retMsg
	bt, err := json.Marshal(b)
	if err != nil {
		return err
	}
	_, err = w.Write(bt)
	return err
}

func (e HttpError) WriteAndAbort(g *gin.Context) {
	if !g.IsAborted() {
		g.AbortWithError(e.code, e.err)
	}
	err := e.Write(g.Writer)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// return nil if isnt
func HttpErr(err error) *HttpError {
	if e, ok := err.(*HttpError); ok {
		return e
	} else {
		return nil
	}
}
