package gifttt

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/drtoful/gifttt/Godeps/_workspace/src/github.com/drtoful/twik"
	"github.com/drtoful/gifttt/Godeps/_workspace/src/github.com/drtoful/twik/ast"
)

var (
	_manager  *VariableManager
	varPrefix = "var~"
)

type VariableManager struct {
	Updates chan *Value
	cache   map[string]*Value
}

type Value struct {
	Name  string      `json:"-"`
	Value interface{} `json:"value"`
}

func GetManager() *VariableManager {
	if _manager == nil {
		_manager = &VariableManager{
			Updates: make(chan *Value),
			cache:   make(map[string]*Value),
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
		return nil, nil
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

	v := &Value{Value: value, Name: name}
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	store := GetStore()
	vm.cache[name] = v
	vm.Updates <- v
	return store.Set(varPrefix+name, string(b))
}

// the GlobalScope encapsulated over the DefaultScope of the LISP
// interpreter. Get/Set will be delegated to it, so we can answer
// with the data in the VariableManager
type GlobalScope struct {
	fset *ast.FileSet
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
	scope := twik.NewDefaultScope(s.fset)
	scope.Enclose(s)
	scope.Create("run", runFn)
	scope.Create("log", logFn)
	return scope.Eval(node)
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

// "log" a message
func logFn(args []interface{}) (interface{}, error) {
	if len(args) == 1 {
		if s, ok := args[0].(string); ok {
			log.Println(s)
			return nil, nil
		}
	}
	return nil, errors.New("log function takes a single string argument")
}

func NewGlobalScope(fset *ast.FileSet) twik.Scope {
	scope := &GlobalScope{
		fset: fset,
	}
	return scope
}

// this scope is trying to find out, which symbols (in gifttt variables)
// there are, that might trigger on value change.
type varScope struct {
	variables chan string
}

func (s *varScope) Create(name string, value interface{}) error {
	panic("not reached")
}

func (s *varScope) Set(name string, value interface{}) error {
	panic("not reached")
}

func (s *varScope) Get(name string) (interface{}, error) {
	panic("not reached")
}

func (s *varScope) Branch() twik.Scope {
	panic("not reached")
}

func (s *varScope) Enclose(scope twik.Scope) error {
	panic("not reached")
}

func (s *varScope) Eval(node ast.Node) (interface{}, error) {
	switch node := node.(type) {
	case *ast.Symbol:
		if node.Name == "log" || node.Name == "run" {
			return nil, nil
		}
		for _, glb := range twik.Globals {
			if node.Name == glb.Name {
				return nil, nil
			}
		}
		s.variables <- node.Name
		return nil, nil
	case *ast.Int:
		return nil, nil
	case *ast.Float:
		return nil, nil
	case *ast.String:
		return nil, nil
	case *ast.List:
		for _, node := range node.Nodes {
			s.Eval(node)
		}
		return nil, nil
	case *ast.Root:
		for _, node := range node.Nodes {
			s.Eval(node)
		}
		close(s.variables)
		return nil, nil
	}
	return nil, nil
}

type Rule struct {
	Name    string
	program ast.Node
	scope   twik.Scope
	lock    *sync.Mutex
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
		lock:    &sync.Mutex{},
	}, nil
}

func (r *Rule) Run() error {
	r.lock.Lock()
	defer r.lock.Unlock()
	_, err := r.scope.Eval(r.program)
	return err
}

type RuleManager struct {
	rules map[string][]*Rule
}

func NewRuleManager(path string) *RuleManager {
	manager := &RuleManager{
		rules: make(map[string][]*Rule),
	}

	files, _ := ioutil.ReadDir(path)
	count := 0
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".rule") {
			filename := filepath.Join(path, f.Name())

			file, err := os.Open(filename)
			if err != nil {
				log.Printf("error opening '%s': %s\n", f.Name(), err.Error())
				continue
			}

			rule, err := NewRule(f.Name(), file)
			file.Close()
			if err != nil {
				log.Println(err)
				continue
			}

			// try to infer which variables are used by this rule, so that we can
			// find out which rules need really to be triggered when a variable
			// changes
			scope := &varScope{variables: make(chan string)}
			go scope.Eval(rule.program)
			for name := range scope.variables {
				if _, ok := manager.rules[name]; !ok {
					manager.rules[name] = []*Rule{}
				}

				rules := manager.rules[name]
				var found bool
				for _, r := range rules {
					found = found || (r.Name == f.Name())
					if found {
						break
					}
				}

				if !found {
					manager.rules[name] = append(rules, rule)
				}
			}
			count += 1
		}
	}
	log.Printf("loaded %d rules\n", count)

	return manager
}

func getSession() string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 16)
	for i := 0; i < len(result); i += 1 {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func (m *RuleManager) Run() {
	vm := GetManager()

	// the rule manager keeps track of time
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			now := time.Now()
			vm.Set("time:second", int64(now.Second()))
			vm.Set("time:minute", int64(now.Minute()))
			vm.Set("time:hour", int64(now.Hour()))
			vm.Set("date:day", int64(now.Day()))
			vm.Set("date:month", int64(now.Month()))
			vm.Set("date:year", int64(now.Year()))
			vm.Set("date:wday", int64(now.Weekday()))
		}
	}()

	for {
		v := <-vm.Updates

		go func() {
			session := getSession()
			if _, ok := m.rules[v.Name]; !ok {
				return
			}

			count := 0
			for _, r := range m.rules[v.Name] {
				go func() {
					err := r.Run()
					if err != nil {
						log.Printf("[%s] error in '%s': %s\n", session, r.Name, err.Error())
					} else {
						log.Printf("[%s] executed rule '%s'\n", session, r.Name)
					}
				}()
				count += 1
			}

			log.Printf("[%s] executed %d rules for change in variable '%s' -> '%#v'\n", session, count, v.Name, v.Value)
		}()
	}
}
