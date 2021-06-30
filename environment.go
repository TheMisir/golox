package main

type Environment struct {
	context *LoxContext
	values  map[string]Any
}

func MakeEnvironment(context *LoxContext) *Environment {
	return &Environment{
		values: make(map[string]Any, 0),
	}
}

func (e *Environment) define(name string, value Any) {
	e.values[name] = value
}

func (e *Environment) get(name *Token) Any {
	if val, ok := e.values[name.lexme]; ok {
		return val
	}

	panic(MakeRuntimeError(name, "Undefined variable '%s'.", name.lexme))
}

func (e *Environment) assign(name *Token, value Any) {
	if _, ok := e.values[name.lexme]; ok {
		e.values[name.lexme] = value
		return
	}

	panic(MakeRuntimeError(name, "Undefined variable '%s'.", name.lexme))
}
