package decode

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/yanmengfei/poc-engine-xray/library/proto"
	"github.com/yanmengfei/poc-engine-xray/library/utils"
	exp "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

var containsStringFunc = &functions.Overload{
	Operator: "contains_string",
	Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
		v1, ok := lhs.(types.String)
		if !ok {
			return types.ValOrErr(lhs, "unexpected type '%v' passed to contains", lhs.Type())
		}
		v2, ok := rhs.(types.String)
		if !ok {
			return types.ValOrErr(rhs, "unexpected type '%v' passed to contains", rhs.Type())
		}
		return types.Bool(strings.Contains(string(v1), string(v2)))
	},
}

var stringIContainsStringDecl = decls.NewFunction("icontains", decls.NewInstanceOverload(
	"string_icontains_string", []*exp.Type{decls.String, decls.String}, decls.Bool))
var stringIContainsStringFunc = &functions.Overload{
	Operator: "string_icontains_string",
	Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
		v1, ok := lhs.(types.String)
		if !ok {
			return types.ValOrErr(lhs, "unexpected type '%v' passed to icontains", lhs.Type())
		}
		v2, ok := rhs.(types.String)
		if !ok {
			return types.ValOrErr(rhs, "unexpected type '%v' passed to icontains", rhs.Type())
		}
		return types.Bool(strings.Contains(strings.ToLower(string(v1)), strings.ToLower(string(v2))))
	},
}

var bytesBContainsBytesDecl = decls.NewFunction("bcontains", decls.NewInstanceOverload(
	"bytes_bcontains_bytes", []*exp.Type{decls.Bytes, decls.Bytes}, decls.Bool))
var bytesBContainsBytesFunc = &functions.Overload{
	Operator: "bytes_bcontains_bytes",
	Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
		v1, ok := lhs.(types.Bytes)
		if !ok {
			return types.ValOrErr(lhs, "unexpected type '%v' passed to bcontains", lhs.Type())
		}
		v2, ok := rhs.(types.Bytes)
		if !ok {
			return types.ValOrErr(rhs, "unexpected type '%v' passed to bcontains", rhs.Type())
		}
		return types.Bool(bytes.Contains(v1, v2))
	},
}

var matchesStringFunc = &functions.Overload{
	Operator: "matches_string",
	Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
		v1, ok := lhs.(types.String)
		if !ok {
			return types.ValOrErr(lhs, "unexpected type '%v' passed to match", lhs.Type())
		}
		v2, ok := rhs.(types.String)
		if !ok {
			return types.ValOrErr(rhs, "unexpected type '%v' passed to match", rhs.Type())
		}
		ok, err := regexp.Match(string(v1), []byte(v2))
		if err != nil {
			return types.NewErr("%v", err)
		}
		return types.Bool(ok)
	},
}

var stringBmatchBytesDecl = decls.NewFunction("bmatches",
	decls.NewInstanceOverload("string_bmatch_bytes",
		[]*exp.Type{decls.String, decls.Bytes},
		decls.Bool))
var stringBmatchBytesFunc = &functions.Overload{
	Operator: "string_bmatch_bytes",
	Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
		v1, ok := lhs.(types.String)
		if !ok {
			return types.ValOrErr(lhs, "unexpected type '%v' passed to bmatch", lhs.Type())
		}
		v2, ok := rhs.(types.Bytes)
		if !ok {
			return types.ValOrErr(rhs, "unexpected type '%v' passed to bmatch", rhs.Type())
		}
		ok, err := regexp.Match(string(v1), v2)
		if err != nil {
			return types.NewErr("%v", err)
		}
		return types.Bool(ok)
	},
}

var stringInMapKeyDecl = decls.NewFunction("in", decls.NewInstanceOverload("string_in_map_key",
	[]*exp.Type{decls.String, decls.NewMapType(decls.String, decls.String)}, decls.Bool))
var stringInMapKeyFunc = &functions.Overload{
	Operator: "string_in_map_key",
	Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
		v1, ok := lhs.(types.String)
		if !ok {
			return types.ValOrErr(lhs, "unexpected type '%v' passed to in", lhs.Type())
		}
		v2, ok := rhs.(types.Bytes)
		if !ok {
			return types.ValOrErr(rhs, "unexpected type '%v' passed to in", lhs.Type())
		}
		return types.Bool(bytes.Contains(v2, []byte(v1)))
	},
}

var md5StringDecl = decls.NewFunction("md5", decls.NewOverload(
	"md5_string", []*exp.Type{decls.String}, decls.String))
var md5StringFunc = &functions.Overload{
	Operator: "md5_string",
	Unary: func(value ref.Val) ref.Val {
		v, ok := value.(types.String)
		if !ok {
			return types.ValOrErr(value, "unexpected type '%v' passed to md5_string", value.Type())
		}
		return types.String(fmt.Sprintf("%x", md5.Sum([]byte(v))))
	},
}

var randomIntDecl = decls.NewFunction("randomInt", decls.NewOverload(
	"randomInt_int_int", []*exp.Type{decls.Int, decls.Int}, decls.Int))
var randomIntFunc = &functions.Overload{
	Operator: "randomInt_int_int",
	Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
		from, ok := lhs.(types.Int)
		if !ok {
			return types.ValOrErr(lhs, "unexpected type '%v' passed to randomInt", lhs.Type())
		}
		to, ok := rhs.(types.Int)
		if !ok {
			return types.ValOrErr(rhs, "unexpected type '%v' passed to randomInt", rhs.Type())
		}
		min, max := int(from), int(to)
		return types.Int(rand.Intn(max-min) + min)
	},
}

//	指定长度的小写字母组成的随机字符串
var randomLowercaseDecl = decls.NewFunction("randomLowercase", decls.NewOverload(
	"randomLowercase_int", []*exp.Type{decls.Int}, decls.String))
var randomLowercaseFunc = &functions.Overload{
	Operator: "randomLowercase_int",
	Unary: func(value ref.Val) ref.Val {
		n, ok := value.(types.Int)
		if !ok {
			return types.ValOrErr(value, "unexpected type '%v' passed to randomLowercase", value.Type())
		}
		return types.String(utils.RandLetters(int(n)))
	},
}

//	将字符串进行 base64 编码
var base64StringDecl = decls.NewFunction("base64", decls.NewOverload(
	"base64_string", []*exp.Type{decls.String}, decls.String))
var base64StringFunc = &functions.Overload{
	Operator: "base64_string",
	Unary: func(value ref.Val) ref.Val {
		v, ok := value.(types.String)
		if !ok {
			return types.ValOrErr(value, "unexpected type '%v' passed to base64_string", value.Type())
		}
		return types.String(base64.StdEncoding.EncodeToString([]byte(v)))
	},
}

//	将bytes进行 base64 编码
var base64BytesDecl = decls.NewFunction("base64", decls.NewOverload(
	"base64_bytes", []*exp.Type{decls.Bytes}, decls.String))
var base64BytesFunc = &functions.Overload{
	Operator: "base64_bytes",
	Unary: func(value ref.Val) ref.Val {
		v, ok := value.(types.Bytes)
		if !ok {
			return types.ValOrErr(value, "unexpected type '%v' passed to base64_bytes", value.Type())
		}
		return types.String(base64.StdEncoding.EncodeToString(v))
	},
}

//	将字符串进行 base64 解码
var base64DecodeStringDecl = decls.NewFunction("base64Decode", decls.NewOverload(
	"base64Decode_string", []*exp.Type{decls.String}, decls.String))
var base64DecodeStringFunc = &functions.Overload{
	Operator: "base64Decode_string",
	Unary: func(value ref.Val) ref.Val {
		v, ok := value.(types.String)
		if !ok {
			return types.ValOrErr(value, "unexpected type '%v' passed to base64Decode_string", value.Type())
		}
		decodeBytes, err := base64.StdEncoding.DecodeString(string(v))
		if err != nil {
			return types.NewErr("%v", err)
		}
		return types.String(decodeBytes)
	},
}

//	将bytes进行 base64 编码
var base64DecodeBytesDecl = decls.NewFunction("base64Decode", decls.NewOverload(
	"base64Decode_bytes", []*exp.Type{decls.Bytes}, decls.String))
var base64DecodeBytesFunc = &functions.Overload{
	Operator: "base64Decode_bytes",
	Unary: func(value ref.Val) ref.Val {
		v, ok := value.(types.Bytes)
		if !ok {
			return types.ValOrErr(value, "unexpected type '%v' passed to base64Decode_bytes", value.Type())
		}
		decodeBytes, err := base64.StdEncoding.DecodeString(string(v))
		if err != nil {
			return types.NewErr("%v", err)
		}
		return types.String(decodeBytes)
	},
}

//	将字符串进行 urlencode 编码
var urlencodeStringDecl = decls.NewFunction("urlencode", decls.NewOverload(
	"urlencode_string", []*exp.Type{decls.String}, decls.String))
var urlencodeStringFunc = &functions.Overload{
	Operator: "urlencode_string",
	Unary: func(value ref.Val) ref.Val {
		v, ok := value.(types.String)
		if !ok {
			return types.ValOrErr(value, "unexpected type '%v' passed to urlencode_string", value.Type())
		}
		return types.String(url.QueryEscape(string(v)))
	},
}

//	将bytes进行 urlencode 编码
var urlencodeBytesDecl = decls.NewFunction("urlencode", decls.NewOverload(
	"urlencode_bytes", []*exp.Type{decls.Bytes}, decls.String))
var urlencodeBytesFunc = &functions.Overload{
	Operator: "urlencode_bytes",
	Unary: func(value ref.Val) ref.Val {
		v, ok := value.(types.Bytes)
		if !ok {
			return types.ValOrErr(value, "unexpected type '%v' passed to urlencode_bytes", value.Type())
		}
		return types.String(url.QueryEscape(string(v)))
	},
}

//	将字符串进行 urldecode 解码
var urldecodeStringDecl = decls.NewFunction("urldecode", decls.NewOverload(
	"urldecode_string", []*exp.Type{decls.String}, decls.String))
var urldecodeStringFunc = &functions.Overload{
	Operator: "urldecode_string",
	Unary: func(value ref.Val) ref.Val {
		v, ok := value.(types.String)
		if !ok {
			return types.ValOrErr(value, "unexpected type '%v' passed to urldecode_string", value.Type())
		}
		decodeString, err := url.QueryUnescape(string(v))
		if err != nil {
			return types.NewErr("%v", err)
		}
		return types.String(decodeString)
	},
}

//	将 bytes 进行 urldecode 解码
var urldecodeBytesDecl = decls.NewFunction("urldecode", decls.NewOverload(
	"urldecode_bytes", []*exp.Type{decls.Bytes}, decls.String))
var urldecodeBytesFunc = &functions.Overload{
	Operator: "urldecode_bytes",
	Unary: func(value ref.Val) ref.Val {
		v, ok := value.(types.Bytes)
		if !ok {
			return types.ValOrErr(value, "unexpected type '%v' passed to urldecode_bytes", value.Type())
		}
		decodeString, err := url.QueryUnescape(string(v))
		if err != nil {
			return types.NewErr("%v", err)
		}
		return types.String(decodeString)
	},
}

//	截取字符串
var substrDecl = decls.NewFunction("substr", decls.NewOverload(
	"substr_string_int_int", []*exp.Type{decls.String, decls.Int, decls.Int}, decls.String))
var substrFunc = &functions.Overload{
	Operator: "substr_string_int_int",
	Function: func(values ...ref.Val) ref.Val {
		if len(values) == 3 {
			str, ok := values[0].(types.String)
			if !ok {
				return types.NewErr("invalid string to 'substr'")
			}
			start, ok := values[1].(types.Int)
			if !ok {
				return types.NewErr("invalid start to 'substr'")
			}
			length, ok := values[2].(types.Int)
			if !ok {
				return types.NewErr("invalid length to 'substr'")
			}
			runes := []rune(str)
			if start < 0 || length < 0 || int(start+length) > len(runes) {
				return types.NewErr("invalid start or length to 'substr'")
			}
			return types.String(runes[start : start+length])
		} else {
			return types.NewErr("too many arguments to 'substr'")
		}
	},
}

//	暂停执行等待指定的秒数
var sleepDecl = decls.NewFunction("sleep", decls.NewOverload(
	"sleep_int", []*exp.Type{decls.Int}, decls.Null))
var sleepFunc = &functions.Overload{
	Operator: "sleep_int",
	Unary: func(value ref.Val) ref.Val {
		v, ok := value.(types.Int)
		if !ok {
			return types.ValOrErr(value, "unexpected type '%v' passed to sleep", value.Type())
		}
		time.Sleep(time.Duration(v) * time.Second)
		return nil
	},
}

//	反连平台结果
var reverseWaitDecl = decls.NewFunction("wait", decls.NewInstanceOverload(
	"reverse_wait_int", []*exp.Type{decls.Any, decls.Int}, decls.Bool))
var reverseWaitFunc = &functions.Overload{
	Operator: "reverse_wait_int",
	Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
		rev, ok := lhs.Value().(*proto.Reverse)
		if !ok {
			return types.ValOrErr(lhs, "unexpected type '%v' passed to 'wait'", lhs.Type())
		}
		timeout, ok := rhs.Value().(int64)
		if !ok {
			return types.ValOrErr(rhs, "unexpected type '%v' passed to 'wait'", rhs.Type())
		}
		return types.Bool(proto.VerifyReverse(rev, timeout))
	},
}
