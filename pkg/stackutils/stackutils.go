package stackutils

import (
	"fmt"
	"runtime"
	"strings"
)

func GetCallingMethodName() string {
	return GetMethodNameAt(3)
}

func GetMethodNameAt(idx int) string {
	pc, _, _, ok := runtime.Caller(idx)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return details.Name()
	}

	return "unknown"
}

// MethodNameFromFullPath get the method name from a full path string
func MethodNameFromFullPath(fullMethod string) string {
	idx := strings.LastIndex(fullMethod, ".")
	if idx > 0 {
		return fullMethod[idx+1:]
	}
	return fullMethod
}

// StructNameFromMethod gets the struct name from a fully qualified method name (or returns a blank if there is no struct
func StructNameFromMethod(methodName string) string {
	idx := strings.LastIndex(methodName, "(")
	if idx > 0 {
		structName := methodName[idx+1:]
		idx = strings.LastIndex(structName, ")")
		if idx > 0 {
			structName = structName[:idx]
			structName = strings.TrimPrefix(structName, "*")
			return structName
		}
	}

	return ""
}

func BasicQualifierFromMethod(fullMethod string) string {
	structName := StructNameFromMethod(fullMethod)
	methodName := MethodNameFromFullPath(fullMethod)

	return fmt.Sprintf("%s.%s", structName, methodName)
}
