package interp

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"reflect"
	"strconv"
)

// nkind defines the kind of AST, i.e. the grammar category
type nkind uint

// Node kinds for the go language
const (
	undefNode nkind = iota
	addressExpr
	arrayType
	assignStmt
	assignXStmt
	basicLit
	binaryExpr
	blockStmt
	branchStmt
	breakStmt
	callExpr
	caseBody
	caseClause
	chanType
	commClause
	compositeLitExpr
	constDecl
	continueStmt
	declStmt
	deferStmt
	defineStmt
	defineXStmt
	ellipsisExpr
	exprStmt
	fallthroughtStmt
	fieldExpr
	fieldList
	fileStmt
	forStmt0     // for {}
	forStmt1     // for cond {}
	forStmt2     // for init; cond; {}
	forStmt3     // for ; cond; post {}
	forStmt3a    // for init; ; post {}
	forStmt4     // for init; cond; post {}
	forRangeStmt // for range
	funcDecl
	funcLit
	funcType
	goStmt
	gotoStmt
	identExpr
	ifStmt0 // if cond {}
	ifStmt1 // if cond {} else {}
	ifStmt2 // if init; cond {}
	ifStmt3 // if init; cond {} else {}
	importDecl
	importSpec
	incDecStmt
	indexExpr
	interfaceType
	keyValueExpr
	labeledStmt
	landExpr
	lorExpr
	mapType
	parenExpr
	rangeStmt
	returnStmt
	rvalueExpr
	rtypeExpr
	selectStmt
	selectorExpr
	selectorImport
	sendStmt
	sliceExpr
	starExpr
	structType
	switchStmt
	switchIfStmt
	typeAssertExpr
	typeDecl
	typeSpec
	typeSwitch
	unaryExpr
	valueSpec
	varDecl
)

var kinds = [...]string{
	undefNode:        "undefNode",
	addressExpr:      "addressExpr",
	arrayType:        "arrayType",
	assignStmt:       "assignStmt",
	assignXStmt:      "assignXStmt",
	basicLit:         "basicLit",
	binaryExpr:       "binaryExpr",
	blockStmt:        "blockStmt",
	branchStmt:       "branchStmt",
	breakStmt:        "breakStmt",
	callExpr:         "callExpr",
	caseBody:         "caseBody",
	caseClause:       "caseClause",
	chanType:         "chanType",
	commClause:       "commClause",
	compositeLitExpr: "compositeLitExpr",
	constDecl:        "constDecl",
	continueStmt:     "continueStmt",
	declStmt:         "declStmt",
	deferStmt:        "deferStmt",
	defineStmt:       "defineStmt",
	defineXStmt:      "defineXStmt",
	ellipsisExpr:     "ellipsisExpr",
	exprStmt:         "exprStmt",
	fallthroughtStmt: "fallthroughStmt",
	fieldExpr:        "fieldExpr",
	fieldList:        "fieldList",
	fileStmt:         "fileStmt",
	forStmt0:         "forStmt0",
	forStmt1:         "forStmt1",
	forStmt2:         "forStmt2",
	forStmt3:         "forStmt3",
	forStmt3a:        "forStmt3a",
	forStmt4:         "forStmt4",
	forRangeStmt:     "forRangeStmt",
	funcDecl:         "funcDecl",
	funcType:         "funcType",
	funcLit:          "funcLit",
	goStmt:           "goStmt",
	gotoStmt:         "gotoStmt",
	identExpr:        "identExpr",
	ifStmt0:          "ifStmt0",
	ifStmt1:          "ifStmt1",
	ifStmt2:          "ifStmt2",
	ifStmt3:          "ifStmt3",
	importDecl:       "importDecl",
	importSpec:       "importSpec",
	incDecStmt:       "incDecStmt",
	indexExpr:        "indexExpr",
	interfaceType:    "interfaceType",
	keyValueExpr:     "keyValueExpr",
	labeledStmt:      "labeledStmt",
	landExpr:         "landExpr",
	lorExpr:          "lorExpr",
	mapType:          "mapType",
	parenExpr:        "parenExpr",
	rangeStmt:        "rangeStmt",
	returnStmt:       "returnStmt",
	rvalueExpr:       "rvalueExpr",
	rtypeExpr:        "rtypeExpr",
	selectStmt:       "selectStmt",
	selectorExpr:     "selectorExpr",
	selectorImport:   "selectorImport",
	sendStmt:         "sendStmt",
	sliceExpr:        "sliceExpr",
	starExpr:         "starExpr",
	structType:       "structType",
	switchStmt:       "switchStmt",
	switchIfStmt:     "switchIfStmt",
	typeAssertExpr:   "typeAssertExpr",
	typeDecl:         "typeDecl",
	typeSpec:         "typeSpec",
	typeSwitch:       "typeSwitch",
	unaryExpr:        "unaryExpr",
	valueSpec:        "valueSpec",
	varDecl:          "varDecl",
}

func (k nkind) String() string {
	if k < nkind(len(kinds)) {
		return kinds[k]
	}
	return "nKind(" + strconv.Itoa(int(k)) + ")"
}

// astError represents an error during AST build stage
type astError error

// action defines the node action to perform at execution
type action uint

// Node actions for the go language
const (
	aNop action = iota
	aAddr
	aAssign
	aAssignX
	aAdd
	aAddAssign
	aAnd
	aAndAssign
	aAndNot
	aAndNotAssign
	aCall
	aCase
	aCompositeLit
	aDec
	aDefer
	aEqual
	aGreater
	aGreaterEqual
	aGetFunc
	aGetIndex
	aInc
	aLand
	aLor
	aLower
	aLowerEqual
	aMethod
	aMul
	aMulAssign
	aNegate
	aNot
	aNotEqual
	aOr
	aOrAssign
	aQuo
	aQuoAssign
	aRange
	aRecv
	aRem
	aRemAssign
	aReturn
	aSend
	aShl
	aShlAssign
	aShr
	aShrAssign
	aSlice
	aSlice0
	aStar
	aSub
	aSubAssign
	aTypeAssert
	aXor
	aXorAssign
)

var actions = [...]string{
	aNop:          "nop",
	aAddr:         "&",
	aAssign:       "=",
	aAssignX:      "X=",
	aAdd:          "+",
	aAddAssign:    "+=",
	aAnd:          "&",
	aAndAssign:    "&=",
	aAndNot:       "&^",
	aAndNotAssign: "&^=",
	aCall:         "call",
	aCase:         "case",
	aCompositeLit: "compositeLit",
	aDec:          "--",
	aDefer:        "defer",
	aEqual:        "==",
	aGreater:      ">",
	aGetFunc:      "getFunc",
	aGetIndex:     "getIndex",
	aInc:          "++",
	aLand:         "&&",
	aLor:          "||",
	aLower:        "<",
	aMethod:       "Method",
	aMul:          "*",
	aMulAssign:    "*=",
	aNegate:       "-",
	aNot:          "!",
	aNotEqual:     "!=",
	aOr:           "|",
	aOrAssign:     "|=",
	aQuo:          "/",
	aQuoAssign:    "/=",
	aRange:        "range",
	aRecv:         "<-",
	aRem:          "%",
	aRemAssign:    "%=",
	aReturn:       "return",
	aSend:         "<~",
	aShl:          "<<",
	aShlAssign:    "<<=",
	aShr:          ">>",
	aShrAssign:    ">>=",
	aSlice:        "slice",
	aSlice0:       "slice0",
	aStar:         "*",
	aSub:          "-",
	aSubAssign:    "-=",
	aTypeAssert:   "TypeAssert",
	aXor:          "^",
	aXorAssign:    "^=",
}

func (a action) String() string {
	if a < action(len(actions)) {
		return actions[a]
	}
	return "Action(" + strconv.Itoa(int(a)) + ")"
}

func (interp *Interpreter) firstToken(src string) token.Token {
	var s scanner.Scanner
	file := interp.fset.AddFile("", interp.fset.Base(), len(src))
	s.Init(file, []byte(src), nil, 0)

	_, tok, _ := s.Scan()
	return tok
}

// Note: no type analysis is performed at this stage, it is done in pre-order
// processing of CFG, in order to accommodate forward type declarations

// ast parses src string containing Go code and generates the corresponding AST.
// The package name and the AST root node are returned.
func (interp *Interpreter) ast(src, name string) (string, *node, error) {
	var inFunc bool

	// Allow incremental parsing of declarations or statements, by inserting
	// them in a pseudo file package or function. Those statements or
	// declarations will be always evaluated in the global scope
	switch interp.firstToken(src) {
	case token.PACKAGE:
		// nothing to do
	case token.CONST, token.FUNC, token.IMPORT, token.TYPE, token.VAR:
		src = "package main;" + src
	default:
		inFunc = true
		src = "package main; func main() {" + src + "}"
	}

	if !interp.buildOk(interp.context, name, src) {
		return "", nil, nil // skip source not matching build constraints
	}

	f, err := parser.ParseFile(interp.fset, name, src, 0)
	if err != nil {
		return "", nil, err
	}

	var root *node
	var anc astNode
	var st nodestack
	var pkgName string

	addChild := func(root **node, anc astNode, pos token.Pos, kind nkind, act action) *node {
		interp.nindex++
		var i interface{}
		n := &node{anc: anc.node, interp: interp, index: interp.nindex, pos: pos, kind: kind, action: act, val: &i, gen: builtin[act]}
		n.start = n
		if anc.node == nil {
			*root = n
		} else {
			anc.node.child = append(anc.node.child, n)
			if anc.node.action == aCase {
				ancAst := anc.ast.(*ast.CaseClause)
				if len(ancAst.List)+len(ancAst.Body) == len(anc.node.child) {
					// All case clause children are collected.
					// Split children in condition and body nodes to desambiguify the AST.
					interp.nindex++
					body := &node{anc: anc.node, interp: interp, index: interp.nindex, pos: pos, kind: caseBody, action: aNop, val: &i, gen: nop}

					if ts := anc.node.anc.anc; ts.kind == typeSwitch && ts.child[1].action == aAssign {
						// In type switch clause, if a switch guard is assigned, duplicate the switch guard symbol
						// in each clause body, so a different guard type can be set in each clause
						name := ts.child[1].child[0].ident
						interp.nindex++
						gn := &node{anc: body, interp: interp, ident: name, index: interp.nindex, pos: pos, kind: identExpr, action: aNop, val: &i, gen: nop}
						body.child = append(body.child, gn)
					}

					// Add regular body children
					body.child = append(body.child, anc.node.child[len(ancAst.List):]...)
					for i := range body.child {
						body.child[i].anc = body
					}
					anc.node.child = append(anc.node.child[:len(ancAst.List)], body)
				}
			}
		}
		return n
	}

	// Populate our own private AST from Go parser AST.
	// A stack of ancestor nodes is used to keep track of current ancestor for each depth level
	ast.Inspect(f, func(nod ast.Node) bool {
		anc = st.top()
		var pos token.Pos
		if nod != nil {
			pos = nod.Pos()
		}
		switch a := nod.(type) {
		case nil:
			anc = st.pop()

		case *ast.ArrayType:
			st.push(addChild(&root, anc, pos, arrayType, aNop), nod)

		case *ast.AssignStmt:
			var act action
			var kind nkind
			if len(a.Lhs) > 1 && len(a.Rhs) == 1 {
				if a.Tok == token.DEFINE {
					kind = defineXStmt
				} else {
					kind = assignXStmt
				}
				act = aAssignX
			} else {
				kind = assignStmt
				switch a.Tok {
				case token.ASSIGN:
					act = aAssign
				case token.ADD_ASSIGN:
					act = aAddAssign
				case token.AND_ASSIGN:
					act = aAndAssign
				case token.AND_NOT_ASSIGN:
					act = aAndNotAssign
				case token.DEFINE:
					kind = defineStmt
					act = aAssign
				case token.SHL_ASSIGN:
					act = aShlAssign
				case token.SHR_ASSIGN:
					act = aShrAssign
				case token.MUL_ASSIGN:
					act = aMulAssign
				case token.OR_ASSIGN:
					act = aOrAssign
				case token.QUO_ASSIGN:
					act = aQuoAssign
				case token.REM_ASSIGN:
					act = aRemAssign
				case token.SUB_ASSIGN:
					act = aSubAssign
				case token.XOR_ASSIGN:
					act = aXorAssign
				}
			}
			n := addChild(&root, anc, pos, kind, act)
			n.nleft = len(a.Lhs)
			n.nright = len(a.Rhs)
			st.push(n, nod)

		case *ast.BasicLit:
			n := addChild(&root, anc, pos, basicLit, aNop)
			n.ident = a.Value
			switch a.Kind {
			case token.CHAR:
				v, _, _, _ := strconv.UnquoteChar(a.Value[1:len(a.Value)-1], '\'')
				n.rval = reflect.ValueOf(v)
			case token.FLOAT:
				v, _ := strconv.ParseFloat(a.Value, 64)
				n.rval = reflect.ValueOf(v)
			case token.IMAG:
				v, _ := strconv.ParseFloat(a.Value[:len(a.Value)-1], 64)
				n.rval = reflect.ValueOf(complex(0, v))
			case token.INT:
				v, _ := strconv.ParseInt(a.Value, 0, 0)
				n.rval = reflect.ValueOf(int(v))
			case token.STRING:
				v, _ := strconv.Unquote(a.Value)
				n.rval = reflect.ValueOf(v)
			}
			st.push(n, nod)

		case *ast.BinaryExpr:
			kind := binaryExpr
			act := aNop
			switch a.Op {
			case token.ADD:
				act = aAdd
			case token.AND:
				act = aAnd
			case token.AND_NOT:
				act = aAndNot
			case token.EQL:
				act = aEqual
			case token.GEQ:
				act = aGreaterEqual
			case token.GTR:
				act = aGreater
			case token.LAND:
				kind = landExpr
				act = aLand
			case token.LOR:
				kind = lorExpr
				act = aLor
			case token.LEQ:
				act = aLowerEqual
			case token.LSS:
				act = aLower
			case token.MUL:
				act = aMul
			case token.NEQ:
				act = aNotEqual
			case token.OR:
				act = aOr
			case token.REM:
				act = aRem
			case token.SUB:
				act = aSub
			case token.SHL:
				act = aShl
			case token.SHR:
				act = aShr
			case token.QUO:
				act = aQuo
			case token.XOR:
				act = aXor
			}
			st.push(addChild(&root, anc, pos, kind, act), nod)

		case *ast.BlockStmt:
			st.push(addChild(&root, anc, pos, blockStmt, aNop), nod)

		case *ast.BranchStmt:
			var kind nkind
			switch a.Tok {
			case token.BREAK:
				kind = breakStmt
			case token.CONTINUE:
				kind = continueStmt
			case token.FALLTHROUGH:
				kind = fallthroughtStmt
			case token.GOTO:
				kind = gotoStmt
			}
			st.push(addChild(&root, anc, pos, kind, aNop), nod)

		case *ast.CallExpr:
			st.push(addChild(&root, anc, pos, callExpr, aCall), nod)

		case *ast.CaseClause:
			st.push(addChild(&root, anc, pos, caseClause, aCase), nod)

		case *ast.ChanType:
			st.push(addChild(&root, anc, pos, chanType, aNop), nod)

		case *ast.CommClause:
			st.push(addChild(&root, anc, pos, commClause, aNop), nod)

		case *ast.CompositeLit:
			st.push(addChild(&root, anc, pos, compositeLitExpr, aCompositeLit), nod)

		case *ast.DeclStmt:
			st.push(addChild(&root, anc, pos, declStmt, aNop), nod)

		case *ast.DeferStmt:
			st.push(addChild(&root, anc, pos, deferStmt, aDefer), nod)

		case *ast.Ellipsis:
			st.push(addChild(&root, anc, pos, ellipsisExpr, aNop), nod)

		case *ast.ExprStmt:
			st.push(addChild(&root, anc, pos, exprStmt, aNop), nod)

		case *ast.Field:
			st.push(addChild(&root, anc, pos, fieldExpr, aNop), nod)

		case *ast.FieldList:
			st.push(addChild(&root, anc, pos, fieldList, aNop), nod)

		case *ast.File:
			pkgName = a.Name.Name
			st.push(addChild(&root, anc, pos, fileStmt, aNop), nod)

		case *ast.ForStmt:
			// Disambiguate variants of FOR statements with a node kind per variant
			var kind nkind
			if a.Cond == nil {
				if a.Init != nil && a.Post != nil {
					kind = forStmt3a
				} else {
					kind = forStmt0
				}
			} else {
				switch {
				case a.Init == nil && a.Post == nil:
					kind = forStmt1
				case a.Init != nil && a.Post == nil:
					kind = forStmt2
				case a.Init == nil && a.Post != nil:
					kind = forStmt3
				default:
					kind = forStmt4
				}
			}
			st.push(addChild(&root, anc, pos, kind, aNop), nod)

		case *ast.FuncDecl:
			n := addChild(&root, anc, pos, funcDecl, aNop)
			if a.Recv == nil {
				// function is not a method, create an empty receiver list
				addChild(&root, astNode{n, nod}, pos, fieldList, aNop)
			}
			st.push(n, nod)

		case *ast.FuncLit:
			n := addChild(&root, anc, pos, funcLit, aGetFunc)
			addChild(&root, astNode{n, nod}, pos, fieldList, aNop)
			addChild(&root, astNode{n, nod}, pos, undefNode, aNop)
			st.push(n, nod)

		case *ast.FuncType:
			st.push(addChild(&root, anc, pos, funcType, aNop), nod)

		case *ast.GenDecl:
			var kind nkind
			switch a.Tok {
			case token.CONST:
				kind = constDecl
			case token.IMPORT:
				kind = importDecl
			case token.TYPE:
				kind = typeDecl
			case token.VAR:
				kind = varDecl
			}
			st.push(addChild(&root, anc, pos, kind, aNop), nod)

		case *ast.GoStmt:
			st.push(addChild(&root, anc, pos, goStmt, aNop), nod)

		case *ast.Ident:
			n := addChild(&root, anc, pos, identExpr, aNop)
			n.ident = a.Name
			st.push(n, nod)
			if n.anc.kind == defineStmt && n.anc.nright == 0 {
				// Implicit assign expression (in a ConstDecl block).
				// Clone assign source and type from previous
				a := n.anc
				pa := a.anc.child[childPos(a)-1]

				if len(pa.child) > pa.nleft+pa.nright {
					// duplicate previous type spec
					a.child = append(a.child, interp.dup(pa.child[a.nleft], a))
				}

				// duplicate previous assign right hand side
				a.child = append(a.child, interp.dup(pa.lastChild(), a))
				a.nright++
			}

		case *ast.IfStmt:
			// Disambiguate variants of IF statements with a node kind per variant
			var kind nkind
			switch {
			case a.Init == nil && a.Else == nil:
				kind = ifStmt0
			case a.Init == nil && a.Else != nil:
				kind = ifStmt1
			case a.Else == nil:
				kind = ifStmt2
			default:
				kind = ifStmt3
			}
			st.push(addChild(&root, anc, pos, kind, aNop), nod)

		case *ast.ImportSpec:
			st.push(addChild(&root, anc, pos, importSpec, aNop), nod)

		case *ast.IncDecStmt:
			var act action
			switch a.Tok {
			case token.INC:
				act = aInc
			case token.DEC:
				act = aDec
			}
			st.push(addChild(&root, anc, pos, incDecStmt, act), nod)

		case *ast.IndexExpr:
			st.push(addChild(&root, anc, pos, indexExpr, aGetIndex), nod)

		case *ast.InterfaceType:
			st.push(addChild(&root, anc, pos, interfaceType, aNop), nod)

		case *ast.KeyValueExpr:
			st.push(addChild(&root, anc, pos, keyValueExpr, aNop), nod)

		case *ast.LabeledStmt:
			st.push(addChild(&root, anc, pos, labeledStmt, aNop), nod)

		case *ast.MapType:
			st.push(addChild(&root, anc, pos, mapType, aNop), nod)

		case *ast.ParenExpr:
			st.push(addChild(&root, anc, pos, parenExpr, aNop), nod)

		case *ast.RangeStmt:
			// Insert a missing ForRangeStmt for AST correctness
			n := addChild(&root, anc, pos, forRangeStmt, aNop)
			r := addChild(&root, astNode{n, nod}, pos, rangeStmt, aRange)
			st.push(r, nod)
			if a.Key == nil {
				// range not in an assign expression: insert a "_" key variable to store iteration index
				k := addChild(&root, astNode{r, nod}, pos, identExpr, aNop)
				k.ident = "_"
			}

		case *ast.ReturnStmt:
			st.push(addChild(&root, anc, pos, returnStmt, aReturn), nod)

		case *ast.SelectStmt:
			st.push(addChild(&root, anc, pos, selectStmt, aNop), nod)

		case *ast.SelectorExpr:
			st.push(addChild(&root, anc, pos, selectorExpr, aGetIndex), nod)

		case *ast.SendStmt:
			st.push(addChild(&root, anc, pos, sendStmt, aSend), nod)

		case *ast.SliceExpr:
			if a.Low == nil {
				st.push(addChild(&root, anc, pos, sliceExpr, aSlice0), nod)
			} else {
				st.push(addChild(&root, anc, pos, sliceExpr, aSlice), nod)
			}

		case *ast.StarExpr:
			st.push(addChild(&root, anc, pos, starExpr, aStar), nod)

		case *ast.StructType:
			st.push(addChild(&root, anc, pos, structType, aNop), nod)

		case *ast.SwitchStmt:
			if a.Tag == nil {
				st.push(addChild(&root, anc, pos, switchIfStmt, aNop), nod)
			} else {
				st.push(addChild(&root, anc, pos, switchStmt, aNop), nod)
			}

		case *ast.TypeAssertExpr:
			st.push(addChild(&root, anc, pos, typeAssertExpr, aTypeAssert), nod)

		case *ast.TypeSpec:
			st.push(addChild(&root, anc, pos, typeSpec, aNop), nod)

		case *ast.TypeSwitchStmt:
			n := addChild(&root, anc, pos, typeSwitch, aNop)
			st.push(n, nod)
			if a.Init == nil {
				// add an empty init node to disambiguate AST
				addChild(&root, astNode{n, nil}, pos, fieldList, aNop)
			}

		case *ast.UnaryExpr:
			var kind = unaryExpr
			var act action
			switch a.Op {
			case token.AND:
				kind = addressExpr
				act = aAddr
			case token.ARROW:
				act = aRecv
			case token.NOT:
				act = aNot
			case token.SUB:
				act = aNegate
			}
			st.push(addChild(&root, anc, pos, kind, act), nod)

		case *ast.ValueSpec:
			kind := valueSpec
			act := aNop
			if a.Values != nil {
				if len(a.Names) > 1 && len(a.Values) == 1 {
					if anc.node.kind == constDecl || anc.node.kind == varDecl {
						kind = defineXStmt
					} else {
						kind = assignXStmt
					}
					act = aAssignX
				} else {
					if anc.node.kind == constDecl || anc.node.kind == varDecl {
						kind = defineStmt
					} else {
						kind = assignStmt
					}
					act = aAssign
				}
			} else if anc.node.kind == constDecl {
				kind, act = defineStmt, aAssign
			}
			n := addChild(&root, anc, pos, kind, act)
			n.nleft = len(a.Names)
			n.nright = len(a.Values)
			st.push(n, nod)

		default:
			err = astError(fmt.Errorf("ast: %T not implemented, line %s", a, interp.fset.Position(pos)))
			return false
		}
		return true
	})
	if inFunc {
		// Incremental parsing: statements were inserted in a pseudo function.
		// Set root to function body so its statements are evaluated in global scope
		root = root.child[1].child[3]
		root.anc = nil
	}
	return pkgName, root, err
}

type astNode struct {
	node *node
	ast  ast.Node
}

type nodestack []astNode

func (s *nodestack) push(n *node, a ast.Node) {
	*s = append(*s, astNode{n, a})
}

func (s *nodestack) pop() astNode {
	l := len(*s) - 1
	res := (*s)[l]
	*s = (*s)[:l]
	return res
}

func (s *nodestack) top() astNode {
	l := len(*s)
	if l > 0 {
		return (*s)[l-1]
	}
	return astNode{}
}

// dup returns a duplicated node subtree
func (interp *Interpreter) dup(nod, anc *node) *node {
	interp.nindex++
	n := *nod
	n.index = interp.nindex
	n.anc = anc
	n.start = &n
	n.pos = anc.pos
	n.child = nil
	for _, c := range nod.child {
		n.child = append(n.child, interp.dup(c, &n))
	}
	return &n
}
