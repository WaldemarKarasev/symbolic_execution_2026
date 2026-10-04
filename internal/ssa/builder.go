// Package ssa предоставляет функции для построения SSA представления
package ssa

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// Builder отвечает за построение SSA из исходного кода Go
type Builder struct {
	fset *token.FileSet
}

// NewBuilder создаёт новый экземпляр Builder
func NewBuilder() *Builder {
	return &Builder{
		fset: token.NewFileSet(),
	}
}

// TODO: Реализуйте следующие методы в рамках домашнего задания

// ParseAndBuildSSA парсит исходный код Go и создаёт SSA представление
// Возвращает SSA программу и функцию по имени
func (b *Builder) ParseAndBuildSSA(source string, funcName string) (*ssa.Function, error) {
	// TODO: Реализовать
	// Шаги:
	file, err := parser.ParseFile(b.fset, "virt-file.go", source, 0)

	if err != nil {
		return nil, err
	}

	pkg := types.NewPackage("input", file.Name.Name)

	conf := &types.Config{
		Importer: importer.Default(),
	}
	// // 2. Создание SSA программы
	ssaPkg, _, err := ssautil.BuildPackage(conf, b.fset, pkg, []*ast.File{file}, ssa.InstantiateGenerics)

	if err != nil {
		return nil, err
	}

	// // 3. Поиск нужной функции по имени
	ssaFunc := ssaPkg.Func(funcName)

	if ssaFunc == nil {
		return nil, fmt.Errorf("Failed to find funcName: %s", funcName)
	}

	return ssaFunc, nil

	// Подсказки:
	// - Используйте parser.ParseFile для парсинга
	// - Создайте packages.Config и загрузите пакет
	// - Используйте ssautil.CreateProgram для создания SSA
	// - Найдите функцию в SSA программе
}
