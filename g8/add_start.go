package g8

import (
	"e8vm.io/e8vm/arch8"
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/types"
)

const (
	initName      = ":init"
	startName     = ":start"
	testStartName = ":test"
)

func findFunc(b *builder, name string, t types.T) *objFunc {
	s := b.scope.Query(name)
	if s == nil {
		return nil
	}
	f, isFunc := s.Item.(*objFunc)
	if !isFunc {
		return nil
	}
	if !types.SameType(f.ref.Type(), t) {
		return nil
	}
	return f
}

func wrapFunc(b *builder, name, wrapName string) {
	f := findFunc(b, name, types.VoidFunc)
	if f == nil {
		return
	}

	b.f = b.p.NewFunc(wrapName, nil, ir.VoidFuncSig)
	b.b = b.f.NewBlock(nil)
	b.b.Call(nil, f.IR())
}

func addStart(b *builder) { wrapFunc(b, "main", startName) }

func addInit(b *builder) { wrapFunc(b, "init", initName) }

var testMainFuncType = types.NewVoidFunc(types.VoidFunc)

func addTestStart(b *builder, testList ir.Ref, n int) {
	b.f = b.p.NewFunc(testStartName, nil, ir.VoidFuncSig)
	b.b = b.f.NewBlock(nil)

	argAddr := ir.Num(arch8.AddrBootArg) // the arg
	index := b.newTempIR(types.Uint)     // to save the index
	b.b.Assign(index, ir.NewAddrRef(argAddr, arch8.RegSize, 0, false, true))

	size := ir.Num(uint32(n))
	checkInRange(b, index, size, "u<")

	base := b.newPtr()
	b.b.Arith(base, nil, "&", testList)
	addr := b.newPtr()
	b.b.Arith(addr, index, "*", ir.Num(arch8.RegSize))
	b.b.Arith(addr, base, "+", addr)

	f := ir.NewFuncPtr(
		ir.VoidFuncSig,
		ir.NewAddrRef(addr, arch8.RegSize, 0, false, true),
	)

	testMain := findFunc(b, "testMain", testMainFuncType)
	if testMain == nil {
		b.b.Call(nil, f)
	} else {
		b.b.Call(nil, testMain.ref.IR(), f)
	}
}
