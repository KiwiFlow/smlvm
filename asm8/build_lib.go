package asm8

import (
	"e8vm.io/e8vm/sym8"
)

func declareSymbol(b *builder, sym *sym8.Symbol) bool {
	// declare in the scope
	exists := b.scope.Declare(sym)
	if exists != nil {
		b.Errorf(sym.Pos, "%q already declared", sym.Name())
		b.Errorf(exists.Pos, "  previously declared here")
		return false
	}
	return true
}

func declareFile(b *builder, pkg *pkg, file *file) {
	// declare functions
	for _, fn := range file.funcs {
		t := fn.Name
		sym := sym8.Make(b.symPkg, t.Lit, SymFunc, fn, t.Pos)
		if !declareSymbol(b, sym) {
			continue
		}

		b.curPkg.declare(sym)
		// b.index(t.Lit, b.curPkg.Declare(sym))
	}

	// declare variables
	for _, v := range file.vars {
		t := v.Name
		sym := sym8.Make(b.symPkg, t.Lit, SymVar, v, t.Pos)
		if !declareSymbol(b, sym) {
			continue
		}

		b.curPkg.declare(sym)
		// b.index(t.Lit, b.curPkg.Declare(sym))
	}
}

func buildPkgScope(b *builder, pkg *pkg) {
	if pkg.imports != nil {
		for as, stmt := range pkg.imports.stmts {
			sym := sym8.Make(b.symPkg, as, SymImport, stmt, stmt.Path.Pos)
			if !declareSymbol(b, sym) {
				continue
			}

			b.curPkg.Import(stmt.pkg.Lib)
			b.importPkg(stmt.path, as)
			// b.index(as, b.curPkg.Require(stmt.lib))
		}
	}

	for _, file := range pkg.files {
		declareFile(b, pkg, file)
	}
}

func checkUnusedImport(b *builder, pkg *pkg) {
	if pkg.imports == nil {
		return
	}

	for as, stmt := range pkg.imports.stmts {
		if _, used := b.pkgUsed[as]; !used {
			b.Errorf(stmt.Path.Pos, "package %q imported but not used", as)
		}
	}
}

func buildLib(b *builder, pkg *pkg) *lib {
	ret := newLib(pkg.path)
	b.curPkg = ret

	b.scope.Push()
	defer b.scope.Pop()

	buildPkgScope(b, pkg)
	if b.Errs() != nil {
		return nil // error on declaring, so just return
	}

	for _, file := range pkg.files {
		buildFile(b, file)
	}

	checkUnusedImport(b, pkg)

	if b.Errs() != nil {
		return nil // error on building
	}

	return ret
}
