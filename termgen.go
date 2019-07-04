package mi

import (
	"fmt"
	"go/types"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ssa"

	"github.com/davecgh/go-spew/spew"
	"github.com/iancoleman/strcase"
	"github.com/mndrix/golog/term"

	"github.com/xpinguin/golog-mi/util"
)

///////////////////////////////////
//////////////////////////////////
/////////////////////////////////
type (
	Type = reflect.Type

	Term     = term.Term
	Callable = term.Callable
)

func Var_(name string) Term {
	defer func() {
		if ex := recover(); ex != nil {
			fmt.Printf("[%s]: %v\n", name, ex)
		}
	}()
	name = strcase.ToCamel(name)
	return term.NewVar(name)
}

func Fntr_(name string, args ...interface{}) Callable {
	subts := []Term{}
	for _, arg := range args {
		var targ Term
		switch t := arg.(type) {
		case Term:
			targ = t
		default:
			if arg == nil {
				targ = term.NewTermList([]Term{})
			} else {
				sarg := fmt.Sprintf("%v", arg)
				if r, e := strconv.ParseFloat(sarg, 64); e == nil {
					targ = term.NewFloat64(r)
				} else if num, e := strconv.Atoi(sarg); e == nil {
					targ = term.NewInt64(int64(num))
				} else {
					targ = term.NewAtom(sarg)
				}
			}
		}
		subts = append(subts, targ)
	}
	return term.NewCallable(name, subts...)
}

/////////////////////////////////
////////////////////////////////
///////////////////////////////
type FuncConvData struct {
	f                *ssa.Function
	blockLengthTerms map[int][][]Term
	blkIdxToDom      map[int]int
}

func (c *FuncConvData) BlockDomIdx(blk *ssa.BasicBlock) int {
	return c.blkIdxToDom[blk.Index]
}

/////////////////////////////////

/////////////////////////////////
////////////////////////////////
var blockLengthTerms = map[int][][]Term{}

func collapseBlockTerms(ts []Term) (r []Term) {
	defer func() {
		var ltrms [][]Term
		if ltrms, _ = blockLengthTerms[len(r)]; ltrms == nil {
			ltrms = [][]Term{}
		}
		blockLengthTerms[len(r)] = append(ltrms, r)
	}()

	if len(ts) < 2 {
		return ts
	}

	r = append(r, ts[len(ts)-1])
	for i := 0; i < len(ts)-1; i++ {
		t := ts[len(ts)-2-i]
		for _, rt := range r {
			if strings.Contains(rt.String(), t.String()) {
				goto _cont
			}
		}
		r = append([]Term{t}, r...)
	_cont:
	}
	return r
}

/////////////////////////////////
/////////////////////////////////
var blkToDomIdx = map[*ssa.BasicBlock]int{}

func blockDomIndex(blk *ssa.BasicBlock) int {
	idx, ok := blkToDomIdx[blk]
	if !ok {
		panic("ssa.Block still unregistred!\n" + spew.Sdump(blk))
	}
	return idx
}

/////////////////////////////////
func FunctionTerm(f *ssa.Function, nameFunc func(*ssa.Function) string) Term {
	if f == nil {
		return nil
	}
	if nameFunc == nil {
		nameFunc = func(f *ssa.Function) string {
			//return util.ShortMethodName(f.String(), true)
			// TODO: anon-funcs
			return util.RelMethodName(f, f.Package())
		}
	}

	nameTerms := []Term{}
	for _, subName := range strings.Split(nameFunc(f), ".") {
		subName = strcase.ToSnake(subName)
		nameTerms = append(nameTerms, term.NewAtom(subName))
	}
	if f.Pkg != nil && f.Pkg.Pkg != nil {
		if path := f.Pkg.Pkg.Path(); path != "" && path != string(*nameTerms[0].(*term.Atom)) {
			nameTerms = append([]Term{term.NewAtom(path)}, nameTerms...)
		}
	}

	paramsTerms := []Term{}
	for _, p := range f.Params {
		pt := p.Type()
		for pt != nil {
			if _, ok := pt.(*types.Struct); ok {
				break
			}
			if ptr, ok := pt.(*types.Pointer); ok {
				pt = ptr.Elem()
			} else if named, ok := pt.(*types.Named); ok {
				pt = named.Underlying()
			} else {
				break
			}
		}

		var pTerms []Term
		var param Term

		if ps, ok := pt.(*types.Struct); ok {
			for i := 0; i < ps.NumFields(); i++ {
				pTerms = append(pTerms, term.NewVar("_"))
			}
		}
		if len(pTerms) > 0 {
			param = term.NewCallable(strcase.ToSnake(p.Type().String()), pTerms...)
		} else {
			param = term.NewVar("_")
		}

		paramsTerms = append(paramsTerms, param)
	}

	//blkIdxToDom := map[int]int{}
	for idxDom, blk := range f.DomPreorder() {
		blkToDomIdx[blk] = idxDom
	}

	fts := []Term{}
	for idxDom, blk := range f.DomPreorder() {
		if t := BlockTerm(blk); t != nil {
			btsList := t.(Callable).Arguments()[1]
			if term.IsList(btsList) {
				bts := term.ListToSlice(btsList)
				///
				// fmt.Printf("%%%% --------- {%d} -------------\n", idxDom)
				// fmt.Println("%% ---------------------------")
				for j, bt := range bts {
					btStr := strings.TrimRight(bt.String(), " \n\r")
					if /*!strings.HasPrefix(btStr, ":=(")*/ j == 0 {
						btStr = fmt.Sprintf("'%.3d': ", idxDom) + btStr
					}
					btStr += " ,"
					/*if j+1 < len(bts) {
						btStr += ","
					} else {
						btStr += "."
					}*/

					fmt.Println(btStr)
				}
				// fmt.Println("%% ---------------------------.")
			}
			///
			fts = append(fts, t)
		}
	}

	return term.NewCallable("func", term.SliceToList(nameTerms), term.SliceToList(paramsTerms), term.SliceToList(fts))
}

/////////////////////////////////
func BlockTerm(blk *ssa.BasicBlock) Term {
	if blk == nil {
		return nil
	}
	bts := []Term{}
	for _, instr := range blk.Instrs {
		//if _, ok := instr.(ssa.Value); !ok {
		if t := InstrTerm(instr); t != nil {
			bts = append(bts, t)
		}
		//}
	}
	if len(bts) == 0 {
		return nil
	}

	return term.NewCallable("block",
		term.NewInt64(int64(blockDomIndex(blk))),
		term.SliceToList(collapseBlockTerms(bts)))
}

/////////////////////////////////
func InstrTerm(instr ssa.Instruction) (t Term) {
	if instr == nil {
		return nil //meta.TypeFntr_("nil")
	} else if v, ok := instr.(ssa.Value); ok {
		return ValueInstrTerm(v)
	}
	///

	///
	blocksFntrs_ := func(blks []*ssa.BasicBlock) (ids []Callable) {
		for _, blk := range blks {
			//ids = append(ids, Fntr_("block", blockDomIndex(blk)))
			ids = append(ids, term.NewAtom(fmt.Sprintf("%.3d", blockDomIndex(blk))))
		}
		return
	}
	jumpFntrs_ := func(jmpToBlocks []*ssa.BasicBlock) (jmps []interface{}) {
		for _, blkFntr := range blocksFntrs_(jmpToBlocks) {
			jmps = append(jmps, Fntr_("jmp", blkFntr))
		}
		return
	}
	///
	switch x := instr.(type) {
	case *ssa.Jump: // block-end
		var jmpFntr Callable
		for _i, _jmp := range jumpFntrs_(x.Block().Succs) {
			switch _i {
			case 0:
				jmpFntr = _jmp.(Callable)
			default:
				panic(fmt.Sprintf("Too many blocks for a single jump: %#v\n\n;;;\n\n%#v", x.Block(), x))
			}
		}
		if jmpFntr == nil {
			panic(fmt.Sprintf("No blocks to jump: %#v\n\n;;;\n\n%#v", x.Block(), x))
		}
		return jmpFntr

	case *ssa.If: // block-end
		return Fntr_("if", append(
			[]interface{}{ValueInstrTerm(x.Cond)},
			jumpFntrs_(x.Block().Succs)...)...)

	case *ssa.Return: // func-end
		rts := []Term{}
		for _, rt := range x.Results {
			rts = append(rts, ValueInstrTerm(rt))
		}
		return Fntr_("return", term.SliceToList(rts))

	case *ssa.MapUpdate:
		return Fntr_("update", ValueInstrTerm(x.Map),
			ValueInstrTerm(x.Key), ValueInstrTerm(x.Value)) //, term.NewVar("_"))

	case *ssa.Store:
		return Fntr_("store", ValueInstrTerm(x.Addr), ValueInstrTerm(x.Val)) //, term.NewVar("_"))

	case *ssa.Send: // side-effect (treat like return?)
		return Fntr_("send", ValueInstrTerm(x.Chan), ValueInstrTerm(x.X)) //, term.NewVar("_"))
	case *ssa.RunDefers: // block|func-post-end
		return Fntr_("run", "defers")
	case *ssa.Panic: // dynamic-block-end
		return Fntr_("panic", ValueInstrTerm(x.X))
	case *ssa.Defer:
		return Fntr_("defer", CallInstrTerm(x.Call))
	case *ssa.Go:
		return Fntr_("go", CallInstrTerm(x.Call))
	default:
		operands := []interface{}{}
		if instr, _ := x.(ssa.Instruction); instr != nil {
			var vrnd interface{}
			for _, rnd := range instr.Operands([]*ssa.Value{}) {
				if rnd == nil {
					vrnd = term.NewAtom("nil")
				} else if (*rnd) != nil {
					vrnd = ValueInstrTerm(*rnd)
				}
				operands = append(operands, vrnd)
			}
		}
		//operands = append(operands, Var_(x.Name())) // return
		return Fntr_(
			strcase.ToSnake(reflect.Indirect(reflect.ValueOf(x)).Type().Name()),
			operands...,
		)
	}
	return
}

/////////////////////////////////
func CallInstrTerm(c ssa.CallCommon) (t Term) {
	meth := []Term{}
	targs := []Term{}
	//rts := []Term{}
	rts_len := 0

	for _, arg := range c.Args {
		targs = append(targs, ValueInstrTerm(arg))
	}

	switch {
	case c.IsInvoke():
		targs = append([]Term{ValueInstrTerm(c.Value)}, targs...)
	case c.Method != nil:
		meth = append(meth, term.NewAtom(util.ShortMethodName(c.Method.String(), true)))
		/*case c.Value != nil:
		//// FIXME: c.Value is a function itself!
		rts = append(rts, ValueInstrTerm(c.Value))*/
	}

	if sgn := c.Signature(); sgn != nil {
		if rs := sgn.Results(); rs != nil {
			rts_len = rs.Len()
		}
	}

	if len(meth) == 0 {
		if f := c.StaticCallee(); f != nil {
			meth = append(meth, term.NewAtom(util.RelMethodName(f, f.Pkg)))
		}
	}

	return Fntr_("call",
		Fntr_("p_r", len(targs), rts_len),
		term.SliceToList(targs),
		term.SliceToList(meth))
}

/////////////////////////////////
func ValueInstrTerm(v ssa.Value) (t Term) {
	if v == nil {
		return nil //meta.TypeFntr_("nil")
	}
	///
	defer func() {
		if t == nil {
			return
		}
		t0 := t
		vreg := v.Name()
		t = Fntr_("=:", t0, vreg)
	}()
	///
	switch x := v.(type) {
	case *ssa.UnOp:
		/// TODO: op/3
		/// eg. op(600, fy, '*').
		if x.CommaOk {
			return Fntr_(x.Op.String(), ValueInstrTerm(x.X),
				term.NewVar("_")) //, Var_(x.Name()))
		}
		return Fntr_(x.Op.String(), ValueInstrTerm(x.X)) //Var_(x.Name()))
	case *ssa.BinOp:
		/// TODO: op/3
		return Fntr_(x.Op.String(),
			ValueInstrTerm(x.X), ValueInstrTerm(x.Y)) //, Var_(x.Name()))
	case *ssa.Parameter:
		for idx, p := range x.Parent().Params {
			if p == x {
				return Fntr_("p", idx)
			}
		}
		return Fntr_("p", term.NewAtom(strcase.ToCamel(x.Name())))
	case *ssa.Const:
		if x.Value != nil {
			return
			return Fntr_("v", x.Value)
		} else {
			return Fntr_("n", x.Name())
		}
	case *ssa.FieldAddr:
		// op(600, fy, '&').
		// op(500, xfy, ':')
		// return Fntr_("&", Fntr_(":", ValueInstrTerm(x.X), x.Field))
		return Fntr_(":", ValueInstrTerm(x.X), x.Field)
	case *ssa.Phi:
		edgeTerms := []Term{}
		for i, edge := range x.Edges {
			var edgeVal interface{}
			switch ee := edge.(type) {
			case *ssa.Const:
				edgeVal = ValueInstrTerm(ee)
			default:
				edgeVal = Fntr_("load", ee.Name())
			}
			blkIdx := blockDomIndex(v.(ssa.Instruction).Block().Preds[i])
			edgeTerm := Fntr_("block", blkIdx, edgeVal)
			edgeTerms = append(edgeTerms, edgeTerm)
		}
		return Fntr_("phi", term.SliceToList(edgeTerms)) //, Fntr_("comment", x.Comment)) //, Var_(x.Name()))
	case *ssa.Call:
		return CallInstrTerm(x.Call)
	case *ssa.Extract:
		return Fntr_("extract", x.Index, Fntr_("tuple", ValueInstrTerm(x.Tuple)))
	case *ssa.Alloc:
		return Fntr_("alloc", Fntr_("heap", x.Heap) /*, Fntr_("comment", x.Comment)*/)
	default:
		operands := []interface{}{}
		if instr, _ := x.(ssa.Instruction); instr != nil {
			var vrnd interface{}
			for _, rnd := range instr.Operands([]*ssa.Value{}) {
				if rnd == nil {
					vrnd = term.NewAtom("nil")
				} else if (*rnd) != nil {
					vrnd = ValueInstrTerm(*rnd)
				}
				operands = append(operands, vrnd)
			}
		}
		//operands = append(operands, Var_(x.Name())) // return
		return Fntr_(
			strcase.ToSnake(reflect.Indirect(reflect.ValueOf(x)).Type().Name()),
			operands...,
		)
	}
	return
}
