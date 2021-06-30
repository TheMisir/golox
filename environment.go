package main

type Environment struct {
	context   *LoxContext
	values    map[string]Any
	enclosing *Environment
}

func MakeEnvironment(context *LoxContext, enclosing *Environment) *Environment {
	return &Environment{
		context:   context,
		values:    make(map[string]Any, 0),
		enclosing: enclosing,
	}
}

func (e *Environment) define(name string, value Any) {
	e.values[name] = value
}

func (e *Environment) get(name *Token) Any {
	if val, ok := e.values[name.lexme]; ok {
		return val
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	e.context.runtimeError(name, "Undefined variable '%s'.", name.lexme)
	return nil // will not be executed
}

func (e *Environment) assign(name *Token, value Any) {
	if _, ok := e.values[name.lexme]; ok {
		e.values[name.lexme] = value
		return
	}

	if e.enclosing != nil {
		e.enclosing.assign(name, value)
		return
	}

	e.context.runtimeError(name, "Undefined variable '%s'.", name.lexme)
}
