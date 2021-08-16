package decode

import (
	"errors"
	"regexp"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types/ref"
	"github.com/yanmengfei/spoce/library/proto"
	exp "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type CelLibrary struct {
	EnvOptions []cel.EnvOption
	ProOptions []cel.ProgramOption
}

func (c *CelLibrary) CompileOptions() []cel.EnvOption {
	return c.EnvOptions
}

func (c *CelLibrary) ProgramOptions() []cel.ProgramOption {
	return c.ProOptions
}

func (c *CelLibrary) UpdateCompileOptions(args map[string]string) {
	for k, v := range args {
		var d *exp.Decl
		if strings.HasPrefix(v, "randomInt") {
			d = decls.NewVar(k, decls.Int)
		} else if strings.HasPrefix(v, "newReverse") {
			d = decls.NewVar(k, decls.NewObjectType("proto.Reverse"))
		} else {
			d = decls.NewVar(k, decls.String)
		}
		c.EnvOptions = append(c.EnvOptions, cel.Declarations(d))
	}
}

func NewCelOption() (c CelLibrary) {
	c.EnvOptions = []cel.EnvOption{
		cel.Container("proto"),
		// 对象类型注入
		cel.Types(
			&proto.UrlType{},
			&proto.Request{},
			&proto.Response{},
			&proto.Reverse{},
		),
		// 定义对象
		cel.Declarations(
			decls.NewVar("request", decls.NewObjectType("proto.Request")),
			decls.NewVar("response", decls.NewObjectType("proto.Response")),
		),
		// 定义运算符
		cel.Declarations(
			bytesBContainsBytesDecl, stringIContainsStringDecl, stringBmatchBytesDecl, md5StringDecl,
			stringInMapKeyDecl, randomIntDecl, randomLowercaseDecl, base64StringDecl,
			base64BytesDecl, base64DecodeStringDecl, base64DecodeBytesDecl, urlencodeStringDecl,
			urlencodeBytesDecl, urldecodeStringDecl, urldecodeBytesDecl, substrDecl, sleepDecl, reverseWaitDecl,
		),
	}
	// 定义运算逻辑
	c.ProOptions = []cel.ProgramOption{cel.Functions(
		containsStringFunc, stringIContainsStringFunc, bytesBContainsBytesFunc, matchesStringFunc, md5StringFunc,
		stringInMapKeyFunc, randomIntFunc, randomLowercaseFunc, stringBmatchBytesFunc, base64StringFunc,
		base64BytesFunc, base64DecodeStringFunc, base64DecodeBytesFunc, urlencodeStringFunc, urlencodeBytesFunc,
		urldecodeStringFunc, urldecodeBytesFunc, substrFunc, sleepFunc, reverseWaitFunc,
	)}
	return c
}

// NewCelEnv new env from set
func NewCelEnv(set map[string]string) (*cel.Env, error) {
	option := NewCelOption()
	if set != nil {
		option.UpdateCompileOptions(set)
	}
	return cel.NewEnv(cel.Lib(&option))
}

// Evaluate 执行运算
func Evaluate(env *cel.Env, expression string, params map[string]interface{}) (ref.Val, error) {
	ast, iss := env.Compile(expression)
	if iss.Err() != nil {
		return nil, iss.Err()
	}
	prg, err := env.Program(ast)
	if err != nil {
		return nil, err
	}
	out, _, err := prg.Eval(params)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Search 完成正则匹配
func Search(re string, body string, set map[string]interface{}) error {
	r, err := regexp.Compile(re)
	if err != nil {
		return err
	}
	result := r.FindStringSubmatch(body)
	names := r.SubexpNames()
	if len(result) > 1 && len(names) > 1 {
		for i, name := range names {
			if i > 0 && i <= len(result) {
				set[name] = result[i]
			}
		}
		return nil
	}
	return errors.New("not matched")
}
