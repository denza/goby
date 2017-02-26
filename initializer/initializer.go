package initializer

import (
	"fmt"
	"github.com/st0012/rooby/ast"
	"github.com/st0012/rooby/object"
)

var (
	ObjectClass = InitializeObjectClass()
	ClassClass  = InitializeClassClass()
)

var BuiltinGlobalMethods = map[string]*object.BuiltInMethod{
	"puts": {
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}

				return object.NULL
			}
		},
		Des:  "Print arguments",
		Name: "puts",
	},
	"class": {
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				switch r := receiver.(type) {
				case *object.BaseObject:
					return r.Class
				case *object.Class:
					return r.Class
				}

				fmt.Print(receiver.Inspect())
				return receiver
			}
		},
		Des:  "return receiver's class",
		Name: "class",
	},
}

var BuiltinClassMethods = map[string]*object.BuiltInMethod{
	"new": {
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				class := receiver.(*object.Class)
				instance := InitializeInstance(class)
				initMethod := class.LookupInstanceMethod("initialize")

				if initMethod != nil {
					instance.InitializeMethod = initMethod.(*object.Method)
				}

				return instance
			}
		},
		Des:  "Initialize class's instance",
		Name: "new",
	},
	"name": {
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				name := receiver.(*object.Class).Name
				nameString := &object.String{Value: name.Value}
				return nameString
			}
		},
		Des:  "return self's name",
		Name: "name",
	},
}

func InitializeMainObject() *object.BaseObject {
	obj := &object.BaseObject{Class: ObjectClass, InstanceVariables: object.NewEnvironment()}
	scope := &object.Scope{Self: obj, Env: object.NewEnvironment()}
	obj.Scope = scope
	return obj
}

func InitializeObjectClass() *object.Class {
	name := &ast.Constant{Value: "Object"}
	globalMethods := object.NewEnvironment()

	for key, value := range BuiltinGlobalMethods {
		globalMethods.Set(key, value)
	}

	class := &object.Class{Name: name, Class: ClassClass, Methods: globalMethods}

	return class
}

func InitializeClassClass() *object.Class {
	methods := object.NewEnvironment()

	for key, value := range BuiltinClassMethods {
		methods.Set(key, value)
	}

	name := &ast.Constant{Value: "Class"}
	class := &object.Class{Name: name, Methods: methods}

	return class
}

func InitializeClass(name *ast.Constant, scope *object.Scope) *object.Class {
	class := &object.Class{Name: name, Methods: object.NewEnvironment(), Class: ClassClass, SuperClass: ObjectClass}
	classScope := &object.Scope{Self: class, Env: object.NewClosedEnvironment(scope.Env)}
	class.Scope = classScope

	return class
}

func InitializeInstance(c *object.Class) *object.BaseObject {
	instance := &object.BaseObject{Class: c, InstanceVariables: object.NewEnvironment()}

	return instance
}