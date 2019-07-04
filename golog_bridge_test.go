package mi

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/mndrix/golog"
)

func TestExtCall(t *testing.T) {
	m := golog.NewInteractiveMachine().SetBindings(NewPrologEnv())
	m = m.RegisterForeign(map[string]golog.ForeignPredicate{
		"ext_call/3": ExtCall,
		"ext_call/4": ExtCall,
		"ext_call/5": ExtCall,
	})

	//////////////////
	env(m).Ext = ext(m).Set("hello", func(x interface{}) string {
		return "tagil11" + ": " + spew.Sdump(x)
	})
	m = m.Consult(`
		hello(X, Ret) :- ext_call('', 'hello', [X], Ret).
	`)
	m.ProveAll(`hello(R0, Rs), Rs = [R0], printf('= %v~n', Rs), printf('= %v~n', R0).`)
}
