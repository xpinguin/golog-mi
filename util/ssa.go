package util

import (
	"bytes"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"regexp"
	"strings"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

/////////
/////////
func LoadProgram(mainPkg string, debugInfo bool) *ssa.Program {
	ldrCfg := loader.Config{
		Build:       &build.Default,
		AllowErrors: true,
	}
	ldrCfg.Import(mainPkg)

	ldr, err := ldrCfg.Load()
	if err != nil {
		log.Fatal("Unable to load pkg: ", err)
		return nil
	}

	buildMode := ssa.BuilderMode(0)
	if debugInfo {
		buildMode &= ssa.GlobalDebug
	}
	prog := ssautil.CreateProgram(ldr, buildMode)
	prog.Build()

	return prog
}

///////////
//////////
func BuildPkg(name, src string) (*ssa.Package, *types.Package) {
	// There is a more substantial test of BuildPackage and the
	// SSA program it builds in ../ssa/builder_test.go.

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, name+".go", src, 0)
	if err != nil {
		log.Fatal(err)
	}

	pkg := types.NewPackage(name, "")
	ssapkg, _, err := ssautil.BuildPackage(&types.Config{Importer: importer.Default()}, fset, pkg, []*ast.File{f}, 0)
	if err != nil {
		log.Fatal(err)
	}
	return ssapkg, pkg
}

///////////
//////////
func FindPackage(prog *ssa.Program, name string) *ssa.Package {
	for _, p := range prog.AllPackages() {
		switch name {
		case p.Pkg.Name(), p.Pkg.Path():
			return p
		}
	}
	return nil
}

func PkgMethod(pkg *ssa.Package, tyMemb ssa.Member, name string) *ssa.Function {
	var ty *ssa.Type
	if ty, _ = tyMemb.(*ssa.Type); ty == nil {
		return nil
	}

	prog := pkg.Prog
	sel := prog.MethodSets.MethodSet(types.NewPointer(ty.Type())).Lookup(pkg.Pkg, name)
	if sel == nil {
		return nil
	}

	return prog.MethodValue(sel)
}

func FindMethod(typeName, funcName string, pkgs ...*ssa.Package) <-chan *ssa.Function {
	allTypes := len(typeName) == 0
	matchingFuncs := make(chan *ssa.Function)

	go func() {
		defer close(matchingFuncs)

		yieldFunc := func(f *ssa.Function) {
			if f != nil {
				matchingFuncs <- f
			}
		}
		for _, pkg := range pkgs {
			if pkg == nil {
				continue
			}
			/*if !isSelfPkgPath(pkg.Pkg.Path()) {
				continue
			}*/
			if allTypes {
				for _, memb := range pkg.Members {
					yieldFunc(PkgMethod(pkg, memb, funcName))
				}
				// package's top-level functions too!
				yieldFunc(pkg.Func(funcName))
			} else {
				yieldFunc(PkgMethod(pkg, pkg.Members[typeName], funcName))
			}
		}
	}()
	return matchingFuncs
}

//////////
func ShortMethodName(methStr string, omitPkg bool) string {
	re := regexp.MustCompile(`(?:(?:[(*]+)(?P<path>[/.\w\d]+)[/](?P<pkg>[\w\d]+)[.](?P<type>[\w\d]+)(?:[)])[.])?(?P<name>[\w\d]+)$`)

	ms := re.FindStringSubmatch(methStr)
	if len(ms) == 0 {
		return ""
	}

	var name, pkg string
	for i, m := range ms {
		switch re.SubexpNames()[i] {
		case "name":
			name = name + m
		case "type":
			name = m + "::" + name
		case "pkg":
			pkg = pkg + m
		}
	}
	if omitPkg {
		return strings.Trim(strings.Replace(name, "::", ".", -1), ". ")
	}
	if pkg != "" {
		pkg += "."
	}
	return strings.Trim(pkg+name, ". ")
}

//////////
func IndirectType(t types.Type) types.Type {
	switch tt := t.(type) {
	case *types.Pointer:
		return tt.Elem()
	}
	return t
}

//////////
func RelTypeName(t types.Type, ssapkg *ssa.Package) string {
	var pkg *types.Package
	if ssapkg != nil {
		pkg = ssapkg.Pkg
	}
	return types.TypeString(t, types.RelativeTo(pkg))
}

func RelMethodName(m *ssa.Function, pkg *ssa.Package) string {
	var recvName string

	if recv := m.Signature.Recv(); recv != nil {
		// Method (declared or wrapper)?
		recvName = RelTypeName(IndirectType(recv.Type().Underlying()), pkg)

	} else if len(m.FreeVars) == 1 && strings.HasSuffix(m.Name(), "$bound") {
		// Bound?
		recvName = RelTypeName(IndirectType(m.FreeVars[0].Type().Underlying()), pkg)
	} else {
		// default
		return m.RelString(pkg.Pkg)
	}

	return strings.TrimLeft(recvName+"."+m.Name(), ". ")
}

///////////
///////////
func FuncStr(f *ssa.Function) string {
	buf := bytes.Buffer{}
	ssa.WriteFunction(&buf, f)
	return buf.String()
}

func FuncStrFields(f *ssa.Function) (flds [][]string) {
	for _, line := range strings.Split(FuncStr(f), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		line = strings.TrimRight(line, "\r\n \t")
		strings.ReplaceAll(line, "\t", " ")
		flds = append(flds, strings.Fields(line))
	}
	return
}

// WriteFunction writes to buf a human-readable "disassembly" of f.
// func WriteFunction(buf *bytes.Buffer, f *Function) {
// 	fmt.Fprintf(buf, "# Name: %s\n", f.String())
// 	if f.Pkg != nil {
// 		fmt.Fprintf(buf, "# Package: %s\n", f.Pkg.Pkg.Path())
// 	}
// 	if syn := f.Synthetic; syn != "" {
// 		fmt.Fprintln(buf, "# Synthetic:", syn)
// 	}
// 	if pos := f.Pos(); pos.IsValid() {
// 		fmt.Fprintf(buf, "# Location: %s\n", f.Prog.Fset.Position(pos))
// 	}

// 	if f.parent != nil {
// 		fmt.Fprintf(buf, "# Parent: %s\n", f.parent.Name())
// 	}

// 	if f.Recover != nil {
// 		fmt.Fprintf(buf, "# Recover: %s\n", f.Recover)
// 	}

// 	from := f.pkg()

// 	if f.FreeVars != nil {
// 		buf.WriteString("# Free variables:\n")
// 		for i, fv := range f.FreeVars {
// 			fmt.Fprintf(buf, "# % 3d:\t%s %s\n", i, fv.Name(), relType(fv.Type(), from))
// 		}
// 	}

// 	if len(f.Locals) > 0 {
// 		buf.WriteString("# Locals:\n")
// 		for i, l := range f.Locals {
// 			fmt.Fprintf(buf, "# % 3d:\t%s %s\n", i, l.Name(), relType(deref(l.Type()), from))
// 		}
// 	}
// 	writeSignature(buf, from, f.Name(), f.Signature, f.Params)
// 	buf.WriteString(":\n")

// 	if f.Blocks == nil {
// 		buf.WriteString("\t(external)\n")
// 	}

// 	// NB. column calculations are confused by non-ASCII
// 	// characters and assume 8-space tabs.
// 	const punchcard = 80 // for old time's sake.
// 	const tabwidth = 8
// 	for _, b := range f.Blocks {
// 		if b == nil {
// 			// Corrupt CFG.
// 			fmt.Fprintf(buf, ".nil:\n")
// 			continue
// 		}
// 		n, _ := fmt.Fprintf(buf, "%d:", b.Index)
// 		bmsg := fmt.Sprintf("%s P:%d S:%d", b.Comment, len(b.Preds), len(b.Succs))
// 		fmt.Fprintf(buf, "%*s%s\n", punchcard-1-n-len(bmsg), "", bmsg)

// 		if false { // CFG debugging
// 			fmt.Fprintf(buf, "\t# CFG: %s --> %s --> %s\n", b.Preds, b, b.Succs)
// 		}
// 		for _, instr := range b.Instrs {
// 			buf.WriteString("\t")
// 			switch v := instr.(type) {
// 			case Value:
// 				l := punchcard - tabwidth
// 				// Left-align the instruction.
// 				if name := v.Name(); name != "" {
// 					n, _ := fmt.Fprintf(buf, "%s = ", name)
// 					l -= n
// 				}
// 				n, _ := buf.WriteString(instr.String())
// 				l -= n
// 				// Right-align the type if there's space.
// 				if t := v.Type(); t != nil {
// 					buf.WriteByte(' ')
// 					ts := relType(t, from)
// 					l -= len(ts) + len("  ") // (spaces before and after type)
// 					if l > 0 {
// 						fmt.Fprintf(buf, "%*s", l, "")
// 					}
// 					buf.WriteString(ts)
// 				}
// 			case nil:
// 				// Be robust against bad transforms.
// 				buf.WriteString("<deleted>")
// 			default:
// 				buf.WriteString(instr.String())
// 			}
// 			buf.WriteString("\n")
// 		}
// 	}
// 	fmt.Fprintf(buf, "\n")
// }
