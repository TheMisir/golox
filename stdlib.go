package main

import (
	"io/ioutil"
	"time"
)

func InitializeStdLib(environment *Environment) {
	environment.define("clock", MakeLoxCallable(0, lox_clock))
	environment.define("readfile", MakeLoxCallable(0, lox_readfile))
	environment.define("writefile", MakeLoxCallable(2, lox_writefile))
}

func lox_clock(interpreter *Interpreter, arguments []Any) Any {
	return float64(time.Now().Unix())
}

func lox_readfile(interpreter *Interpreter, arguments []Any) Any {
	content, err := ioutil.ReadFile(arguments[0].(string))
	if err != nil {
		return nil
	}

	return string(content)
}

func lox_writefile(interpreter *Interpreter, arguments []Any) Any {
	err := ioutil.WriteFile(arguments[0].(string), []byte(arguments[1].(string)), 0644)
	return err == nil
}
