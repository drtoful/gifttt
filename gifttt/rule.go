package gifttt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"

	"github.com/drtoful/twik"
	"github.com/drtoful/twik/ast"
)

var (
	_manager  *VariableManager
	varPrefix = "var~"
)

type VariableManager struct {
	cache map[string]*Value
}

type Value struct {
	Value interface{} `json:"value"`
}

func GetManager() *VariableManager {
	if _manager == nil {
		_manager = &VariableManager{
			cache: make(map[string]*Value),
		}
	}
	return _manager
}

func (vm *VariableManager) Get(name string) (interface{}, error) {
	// check cache first
	if v, ok := vm.cache[name]; ok {
		return v.Value, nil
	}

	store := GetStore()
	b, err := store.Get(varPrefix + name)
	if err != nil {
		return nil, fmt.Errorf("undefined symbol: %s", name)
	}

	v := &Value{}
	if err := json.Unmarshal([]byte(b), v); err == nil {
		return v.Value, nil
	} else {
		return nil, err
	}

	panic("never reached")
}

func (vm *VariableManager) Set(name string, value interface{}) error {
	// check if the value has changed since the last time we set it
	old, err := vm.Get(name)
	if err == nil && old == value {
		return nil
	}

	v := &Value{Value: value}
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	store := GetStore()
	vm.cache[name] = v
	return store.Set(varPrefix+name, string(b))
}

// the GlobalScope encapsulated over the DefaultScope of the LISP
// interpreter. Get/Set will be delegated to it, so we can answer
// with the data in the VariableManager
type GlobalScope struct {
	delegate twik.Scope
}

func (s *GlobalScope) Create(symbol string, value interface{}) error {
	panic("never reached")
}

func (s *GlobalScope) Set(symbol string, value interface{}) error {
	manager := GetManager()
	return manager.Set(symbol, value)
}

func (s *GlobalScope) Get(symbol string) (interface{}, error) {
	manager := GetManager()
	return manager.Get(symbol)
}

func (s *GlobalScope) Branch() twik.Scope {
	panic("never reached")
}

func (s *GlobalScope) Eval(node ast.Node) (interface{}, error) {
	return s.delegate.Eval(node)
}

func (s *GlobalScope) Enclose(parent twik.Scope) error {
	panic("never reached")
}

// "run" let's the user execute arbitrary commands
func runFn(args []interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New("run takes at least one argument")
	}

	commands := []string{}
	for _, arg := range args {
		if s, ok := arg.(string); ok {
			commands = append(commands, s)
		} else {
			return nil, errors.New("run only takes string arguments")
		}
	}

	var cmd *exec.Cmd
	if len(commands) == 1 {
		cmd = exec.Command(commands[0])
	} else {
		cmd = exec.Command(commands[0], commands[1:]...)
	}

	if err := cmd.Run(); err == nil {
		cmd.Wait()
	}

	return nil, nil
}

func NewGlobalScope(fset *ast.FileSet) twik.Scope {
	scope := &GlobalScope{
		delegate: twik.NewDefaultScope(fset),
	}
	scope.delegate.Enclose(scope)
	scope.delegate.Create("run", runFn)
	return scope
}

type Rule struct {
	Name    string
	program ast.Node
	scope   twik.Scope
}

func NewRule(name string, r io.Reader) (*Rule, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	fset := twik.NewFileSet()
	scope := NewGlobalScope(fset)

	node, err := twik.Parse(fset, name, data)
	if err != nil {
		return nil, err
	}

	return &Rule{
		Name:    name,
		program: node,
		scope:   scope,
	}, nil
}

func (r *Rule) Run() error {
	_, err := r.scope.Eval(r.program)
	return err
}
