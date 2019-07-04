package mi

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/mndrix/golog"
	"github.com/mndrix/golog/term"
	"github.com/mndrix/ps"

	"github.com/xpinguin/golog-mi/util"
)

////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////
type (
	Variable = term.Variable
	Var      = term.Variable
	Bindings = term.Bindings

	Machine    = golog.Machine
	ForeignRet = golog.ForeignReturn
)

////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////
// Bindings (from the original pkg by "mndrix")
func NewPrologEnv() (env *PrologEnv) {
	env = &PrologEnv{
		bindings: ps.NewMap(),
		names:    ps.NewMap(),
		Ext:      ps.NewMap(),
	}
	return env
}

type PrologEnv struct {
	bindings ps.Map // v.Indicator() => term.Term
	names    ps.Map // v.Name => *Variable
	Ext      ps.Map
}

func (pl *PrologEnv) Bind(v *Variable, val term.Term) (Bindings, error) {
	_, ok := pl.bindings.Lookup(v.Indicator())
	if ok {
		// binding already exists for this variable
		return pl, term.AlreadyBound
	}

	// at this point, we know that v is a free variable

	// create a new environment with the binding in place
	newEnv := pl.clone()
	newEnv.bindings = pl.bindings.Set(v.Indicator(), val)

	// error slot in return is for attributed variables someday
	return newEnv, nil
}
func (pl *PrologEnv) Resolve_(v *Variable) term.Term {
	r, err := pl.Resolve(v)
	util.Exn(err)
	return r
}

func (pl *PrologEnv) Resolve(v *Variable) (term.Term, error) {
	for {
		t, err := pl.Value(v)
		if err == term.NotBound {
			return v, nil
		}
		if err != nil {
			return nil, err
		}
		if term.IsVariable(t) {
			v = t.(*Variable)
		} else {
			return t.ReplaceVariables(pl), nil
		}
	}
}
func (pl *PrologEnv) Size() int {
	return pl.bindings.Size()
}
func (pl *PrologEnv) Value(v *Variable) (term.Term, error) {
	name := v.Indicator()
	value, ok := pl.bindings.Lookup(name)
	if !ok {
		return nil, term.NotBound
	}
	return value.(term.Term), nil
}
func (pl *PrologEnv) clone() *PrologEnv {
	newEnv := *pl
	return &newEnv
}

func (pl *PrologEnv) ByName(name string) (term.Term, error) {
	v, ok := pl.names.Lookup(name)
	if !ok {
		return nil, term.NotBound
	}
	return pl.Resolve(v.(*Variable))
}

func (pl *PrologEnv) ByName_(name string) term.Term {
	x, err := pl.ByName(name)
	util.Exn(err)
	return x
}

func (pl *PrologEnv) String() string {
	return pl.bindings.String()
}

func (pl *PrologEnv) WithNames(names ps.Map) Bindings {
	if !pl.names.IsNil() {
		panic("Can't set names when names have already been set")
	}

	b := pl.clone()
	b.names = names
	return b
}

////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////
func AtomStr(t Term) (string, bool) {
	a, ok := t.(*term.Atom)
	if !ok {
		return "", false
	}
	return a.Name(), true
}

////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////
func env(m Machine) *PrologEnv {
	return m.Bindings().(*PrologEnv)
}

func ext(m Machine) ps.Map {
	return env(m).Ext
}

////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////
/// ext_call/5 (+Pfx, +Name, ?Args, ?Vargs, ?Ret)
/// ext_call/4 (+Pfx, +Name, ?Args, ?Ret)
/// ext_call/3 (+Pfx, +Name, -Ret) % TODO: ?Ret
func ExtCall(m Machine, args []Term) ForeignRet {
	var pfx, fn string
	var fargs []Term

	var ret Term
	ret, args = args[len(args)-1], args[:len(args)-1]

	switch len(args) {
	case 4:
		//unify varargs with map
		fallthrough
	case 3:
		//unify args with list
		fargs = term.ListToSlice(args[2])
		fallthrough
	case 2:
		pfx, _ = AtomStr(args[0])
		fn, _ = AtomStr(args[1])

	default:
		return golog.ForeignFail()
	}

	switch strings.TrimSpace(pfx) {
	case "":
	case "go":
		log.Fatal("go-call is not implemented")
	case "defer":
		log.Fatal("defer is not implemented")
	default:
		log.Fatal("Unknown ext_call prefix: ", pfx)
	}

	f, ok := ext(m).Lookup(fn)
	if !ok {
		return golog.ForeignFail()
	}

	call := reflect.ValueOf(f)
	var callargs []reflect.Value
	for _, farg := range fargs {
		//t, _ := read.Term(farg.String())
		callargs = append(callargs, reflect.ValueOf(farg.String()))
	}
	if call.Kind() != reflect.Func {
		panic(fmt.Errorf("IS NOT FUNC: %#v", call))
	}

	callret := call.Call(callargs)
	switch len(callret) {
	case 0:
		return golog.ForeignUnify(term.SliceToList([]Term{}), ret)
	case 1:
		return golog.ForeignUnify(
			term.SliceToList([]Term{
				term.NewAtom(fmt.Sprintf("%v", callret[0].Interface())),
			}),
			ret)
	default:
		rterms := []Term{}
		for _, cret := range callret {
			rterms = append(rterms,
				term.NewAtom(fmt.Sprintf("%v", cret.Interface())))
		}
		return golog.ForeignUnify(term.SliceToList(rterms), ret)
	}

	/*switch ret.Type() {
	case VariableType:
	case FloatType:
	case	 IntegerType:
	case AtomType:
	}*/
	fmt.Println("WTF?!")
	return golog.ForeignFail()

}
