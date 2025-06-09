/**
 * Description：
 * FileName：const.go
 * Author：CJiaの用心
 * Create：2025/6/8 22:32:19
 * Remark：
 */

package menu

type MethodConst int

const (
	MethodConstGET    MethodConst = iota + 1 // GET
	MethodConstPOST                          // POST
	MethodConstPUT                           // PUT
	MethodConstDELETE                        // DELETE
)

// MethodMapping 接口请求方法映射
var MethodMapping = map[MethodConst]string{
	MethodConstGET:    "GET",
	MethodConstPOST:   "POST",
	MethodConstPUT:    "PUT",
	MethodConstDELETE: "DELETE",
}

// MethodImportMapping 接口请求方法映射
var MethodImportMapping = map[string]MethodConst{
	"GET":    MethodConstGET,
	"POST":   MethodConstPOST,
	"PUT":    MethodConstPUT,
	"DELETE": MethodConstDELETE,
}
