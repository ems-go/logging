package logging

import (
	"fmt"

	flag "github.com/spf13/pflag"
)

const (
	LOG_DIR_NAME string = "loggingDir"
)

func AddFlags(fs *flag.FlagSet) {
	fs.AddFlag(flag.Lookup(LOG_DIR_NAME))
}

type DbgValue string

// const (
//
//	DbgFalse DbgValue = "."
//
// )
// type
var (
	DefaultLogPath = LoggingDir(LOG_DIR_NAME, ".", "appcation config logging file dir")
)

func (v *DbgValue) IsBoolFlag() bool {
	return true
}

func (v *DbgValue) Get() interface{} {
	return DbgValue(*v)
}

func (v *DbgValue) Set(s string) error {

	*v = DbgValue(s)

	return nil

}

func (v *DbgValue) String() string {

	return fmt.Sprintf("%v", *v)
}

func (v *DbgValue) Type() string {
	return "string"
}

func DbgVar(p *DbgValue, name string, value DbgValue, usage string) {
	*p = value
	flag.Var(p, name, usage)
	flag.Lookup(name).NoOptDefVal = "."
}

func LoggingDir(name string, value DbgValue, usage string) *DbgValue {
	p := new(DbgValue)
	DbgVar(p, name, value, usage)
	return p
}

const (
	LOG_WRITE_FILE_NAME string = "loggingFile"
)

func AddFileFlags(fs *flag.FlagSet) {
	fs.AddFlag(flag.Lookup(LOG_WRITE_FILE_NAME))
}

type DbgFileValue bool

// const (
// 	DbgFalse DbgValue = "."
// )

var (
	dbgFileValue = LoggingFile(LOG_WRITE_FILE_NAME, "true", "appcation config logging write to file ")
)

func (v *DbgFileValue) IsBoolFlag() bool {
	return true
}

func (v *DbgFileValue) Get() interface{} {
	return DbgFileValue(*v)
}

func (v *DbgFileValue) Set(s string) error {
	if s == "true" {
		*v = DbgFileValue(true)
	} else {
		*v = DbgFileValue(false)
	}

	return nil

}

func (v *DbgFileValue) String() string {

	return fmt.Sprintf("%v", *v)
}

func (v *DbgFileValue) Type() string {
	return "bool"
}

func LoggingFile(name string, value DbgValue, usage string) *DbgFileValue {
	p := new(DbgFileValue)
	flag.Var(p, name, usage)
	flag.Lookup(name).NoOptDefVal = "true"
	// DbgVar(p, name, value, usage)
	return p
}
