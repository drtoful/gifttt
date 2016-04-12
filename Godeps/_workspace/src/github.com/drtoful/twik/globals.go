package twik

import (
	"errors"
	"fmt"
	"strings"

	"github.com/drtoful/gifttt/Godeps/_workspace/src/github.com/drtoful/twik/ast"
)

var Globals = []struct {
	Name  string
	value interface{}
}{
	{"true", true},
	{"false", false},
	{"nil", nil},
	{"error", errorFn},
	{"==", eqFn},
	{"!=", neFn},
	{"+", plusFn},
	{"-", minusFn},
	{"*", mulFn},
	{"/", divFn},
	{">", gtFn},
	{">=", gteFn},
	{"<", ltFn},
	{"<=", lteFn},
	{"or", orFn},
	{"and", andFn},
	{"if", ifFn},
	{"when", whenFn},
	{"unless", unlessFn},
	{"var", varFn},
	{"set", setFn},
	{"do", doFn},
	{"func", funcFn},
	{"for", forFn},
	{"range", rangeFn},
	{"split", splitFn},
	{"nth", nthFn},
	{"length", lengthFn},
}

func errorFn(args []interface{}) (value interface{}, err error) {
	if len(args) == 1 {
		if s, ok := args[0].(string); ok {
			return nil, errors.New(s)
		}
	}
	return nil, errors.New("error function takes a single string argument")
}

func eqFn(args []interface{}) (value interface{}, err error) {
	if len(args) != 2 {
		return nil, errors.New("== takes two values")
	}
	return args[0] == args[1], nil
}

func neFn(args []interface{}) (value interface{}, err error) {
	if len(args) != 2 {
		return nil, errors.New("!= takes two values")
	}
	return args[0] != args[1], nil
}

func cmpFn(fn string, args []interface{}) (value interface{}, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(`function "%s" takes two arguments`, fn)
	}

	var resf float64
	factor := 1.0
	for _, arg := range args {
		switch arg := arg.(type) {
		case int64:
			resf += float64(arg) * factor
		case float64:
			resf += arg * factor
		default:
			return nil, fmt.Errorf("cannot compare %#v", arg)
		}
		factor = -1
	}

	switch fn {
	case ">":
		return resf > 0, nil
	case ">=":
		return resf >= 0, nil
	case "<":
		return resf < 0, nil
	case "<=":
		return resf <= 0, nil
	}

	panic("illegal compare function")
}

func gtFn(args []interface{}) (value interface{}, err error) {
	return cmpFn(">", args)
}

func gteFn(args []interface{}) (value interface{}, err error) {
	return cmpFn(">=", args)
}

func ltFn(args []interface{}) (value interface{}, err error) {
	return cmpFn("<", args)
}

func lteFn(args []interface{}) (value interface{}, err error) {
	return cmpFn("<=", args)
}

func splitFn(args []interface{}) (interface{}, error) {
	if len(args) >= 2 {
		text, ok1 := args[0].(string)
		sep, ok2 := args[1].(string)
		if ok1 && ok2 {
			slice := strings.Split(text, sep)
			result := make([]interface{}, len(slice))
			for i := range slice {
				result[i] = slice[i]
			}
			return result, nil
		}
	}
	return nil, errors.New("split function takes two string arguments")
}

func nthFn(args []interface{}) (interface{}, error) {
	if len(args) == 2 {
		list, ok1 := args[0].([]interface{})
		n, ok2 := args[1].(int64)

		if ok1 && ok2 {
			if int(n) >= len(list) {
				return nil, errors.New("index out of bounds")
			}
			return list[n], nil
		}
	}
	return nil, errors.New("nth function takes a list and integer argument")
}

func lengthFn(args []interface{}) (interface{}, error) {
	if len(args) == 1 {
		if list, ok := args[0].([]interface{}); ok {
			return int64(len(list)), nil
		}
	}
	return nil, errors.New("length function takes a list as argument")
}

func plusFn(args []interface{}) (value interface{}, err error) {
	var resi int64
	var resf float64
	var f bool
	for _, arg := range args {
		switch arg := arg.(type) {
		case int64:
			resi += arg
			resf += float64(arg)
		case float64:
			resf += arg
			f = true
		default:
			return nil, fmt.Errorf("cannot sum %#v", arg)
		}
	}
	if f {
		return resf, nil
	}
	return resi, nil
}

func minusFn(args []interface{}) (value interface{}, err error) {
	if len(args) == 0 {
		return nil, fmt.Errorf(`function "-" takes one or more arguments`)
	}
	var resi int64
	var resf float64
	var f bool
	for i, arg := range args {
		switch arg := arg.(type) {
		case int64:
			if i == 0 && len(args) > 1 {
				resi = arg
				resf = float64(arg)
			} else {
				resi -= arg
				resf -= float64(arg)
			}
		case float64:
			if i == 0 && len(args) > 1 {
				resf = arg
			} else {
				resf -= arg
			}
			f = true
		default:
			return nil, fmt.Errorf("cannot subtract %#v", arg)
		}
	}
	if f {
		return resf, nil
	}
	return resi, nil
}

func mulFn(args []interface{}) (value interface{}, err error) {
	var resi = int64(1)
	var resf = float64(1)
	var f bool
	for _, arg := range args {
		switch arg := arg.(type) {
		case int64:
			resi *= arg
			resf *= float64(arg)
		case float64:
			resf *= arg
			f = true
		default:
			return nil, fmt.Errorf("cannot multiply %#v", arg)
		}
	}
	if f {
		return resf, nil
	}
	return resi, nil
}

func divFn(args []interface{}) (value interface{}, err error) {
	if len(args) < 2 {
		return nil, fmt.Errorf(`function "/" takes two or more arguments`)
	}
	var resi int64
	var resf float64
	var f bool
	for i, arg := range args {
		switch arg := arg.(type) {
		case int64:
			if i == 0 && len(args) > 1 {
				resi = arg
				resf = float64(arg)
			} else {
				resi /= arg
				resf /= float64(arg)
			}
		case float64:
			if i == 0 && len(args) > 1 {
				resf = float64(arg)
			} else {
				resf /= arg
			}
			f = true
		default:
			return nil, fmt.Errorf("cannot divide with %#v", arg)
		}
	}
	if f {
		return resf, nil
	}
	return resi, nil
}

func andFn(scope Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) == 0 {
		return true, nil
	}
	for _, arg := range args {
		value, err = scope.Eval(arg)
		if err != nil {
			return nil, err
		}
		if value == false {
			return false, nil
		}
	}
	return value, err
}

func orFn(scope Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) == 0 {
		return false, nil
	}
	for _, arg := range args {
		value, err = scope.Eval(arg)
		if err != nil {
			return nil, err
		}
		if value != false {
			return value, nil
		}
	}
	return value, err
}

func ifFn(scope Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) != 3 {
		return nil, errors.New(`function "if" takes three arguments`)
	}
	value, err = scope.Eval(args[0])
	if err != nil {
		return nil, err
	}
	if value == false {
		if len(args) == 3 {
			return scope.Eval(args[2])
		}
		return false, nil
	}
	return scope.Eval(args[1])
}

func whenFn(scope Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) != 2 {
		return nil, errors.New(`function "when" takes two arguments`)
	}

	value, err = scope.Eval(args[0])
	if err != nil {
		return nil, err
	}

	if value == false {
		return false, nil
	}
	return scope.Eval(args[1])
}

func unlessFn(scope Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) != 2 {
		return nil, errors.New(`function "unless" takes two arguments`)
	}

	value, err = scope.Eval(args[0])
	if err != nil {
		return nil, err
	}

	if value == false {
		return scope.Eval(args[1])
	}
	return false, nil
}

func varFn(scope Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) == 0 || len(args) > 2 {
		return nil, errors.New("var takes one or two arguments")
	}
	symbol, ok := args[0].(*ast.Symbol)
	if !ok {
		return nil, errors.New("var takes a symbol as first argument")
	}
	if len(args) == 1 {
		value = nil
	} else {
		value, err = scope.Eval(args[1])
		if err != nil {
			return nil, err
		}
	}
	return nil, scope.Create(symbol.Name, value)
}

func setFn(scope Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) != 2 {
		return nil, errors.New(`function "set" takes two arguments`)
	}
	symbol, ok := args[0].(*ast.Symbol)
	if !ok {
		return nil, errors.New(`function "set" takes a symbol as first argument`)
	}
	value, err = scope.Eval(args[1])
	if err != nil {
		return nil, err
	}
	return nil, scope.Set(symbol.Name, value)
}

func doFn(scope Scope, args []ast.Node) (value interface{}, err error) {
	scope = scope.Branch()
	for _, arg := range args {
		value, err = scope.Eval(arg)
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

func funcFn(scope Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) < 2 {
		return nil, errors.New(`func takes three or more arguments`)
	}
	i := 0
	var name string
	if symbol, ok := args[0].(*ast.Symbol); ok {
		name = symbol.Name
		i++
	}
	list, ok := args[i].(*ast.List)
	if !ok {
		return nil, errors.New(`func takes a list of parameters`)
	}
	params := list.Nodes
	for _, param := range params {
		if _, ok := param.(*ast.Symbol); !ok {
			return nil, errors.New("func's list of parameters must be a list of symbols")
		}
	}
	body := args[i+1:]
	if len(body) == 0 {
		return nil, fmt.Errorf("func takes a body sequence")
	}
	fn := func(args []interface{}) (value interface{}, err error) {
		if len(args) != len(params) {
			nameInfo := "anonymous function"
			if name != "" {
				nameInfo = fmt.Sprintf("function %q", name)
			}
			switch len(params) {
			case 0:
				return nil, fmt.Errorf("%s takes no arguments", nameInfo)
			case 1:
				return nil, fmt.Errorf("%s takes one argument", nameInfo)
			default:
				return nil, fmt.Errorf("%s takes %d arguments", nameInfo, len(params))
			}
		}
		scope = scope.Branch()
		for i, arg := range args {
			err := scope.Create(params[i].(*ast.Symbol).Name, arg)
			if err != nil {
				panic("must not happen: " + err.Error())
			}
		}
		for _, node := range body {
			value, err = scope.Eval(node)
			if err != nil {
				return nil, err
			}
		}
		return value, nil
	}
	if name != "" {
		if err = scope.Create(name, fn); err != nil {
			return nil, err
		}
	}
	return fn, nil
}

func forFn(scope Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) < 4 {
		return nil, errors.New(`for takes four or more arguments`)
	}
	init, test, step, code := args[0], args[1], args[2], args[3:]
	scope = scope.Branch()
	_, err = scope.Eval(init)
	if err != nil {
		return nil, err
	}
	for {
		more, err := scope.Eval(test)
		if err != nil {
			return nil, err
		}
		if more == false {
			return value, nil
		}

		for _, c := range code {
			value, err = scope.Eval(c)
			if err != nil {
				return nil, err
			}
		}

		_, err = scope.Eval(step)
		if err != nil {
			return nil, err
		}
	}
	panic("unreachable")
}

func rangeFn(scope Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) < 3 {
		return nil, errors.New(`range takes three or more arguments`)
	}
	var iname, ename string
	if symbol, ok := args[0].(*ast.Symbol); ok {
		iname = symbol.Name
	} else if list, ok := args[0].(*ast.List); ok && len(list.Nodes) == 2 {
		symbol1, ok1 := list.Nodes[0].(*ast.Symbol)
		symbol2, ok2 := list.Nodes[1].(*ast.Symbol)
		if ok1 && ok2 {
			iname = symbol1.Name
			ename = symbol2.Name
		}
	}
	if iname == "" {
		return nil, errors.New(`range takes var name or (i elem) var name pair as first argument`)
	}
	scope = scope.Branch()
	value, err = scope.Eval(args[1])
	if err != nil {
		return nil, err
	}
	code := args[2:]
	if n, ok := value.(int64); ok {
		scope.Create(iname, 0)
		for i := int64(0); i < n; i++ {
			scope.Set(iname, i)
			for _, c := range code {
				value, err = scope.Eval(c)
				if err != nil {
					return nil, err
				}
			}
		}
		return value, nil
	}
	if list, ok := value.([]interface{}); ok {
		scope.Create(iname, 0)
		scope.Create(ename, nil)
		for i, e := range list {
			scope.Set(iname, i)
			scope.Set(ename, e)
			for _, c := range code {
				value, err = scope.Eval(c)
				if err != nil {
					return nil, err
				}
			}
		}
		return value, nil
	}
	return nil, errors.New(`range takes an integer or a list as second argument`)
}
