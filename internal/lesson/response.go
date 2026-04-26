// package lesson 提供统一的RESTful响应格式
// Java对比: 类似SpringBoot的统一返回体 Result<T>
package lesson

import (
	"encoding/json"
	"net/http"
)

// Response 是所有学习端点的统一响应结构
// Java对比: public class Result<T> { private int code; private String message; private T data; }
type Response struct {
	Code    int         `json:"code"`              // HTTP状态码
	Topic   string      `json:"topic"`             // 本节主题
	Java    string      `json:"java_equivalent"`   // Java对比
	Summary string      `json:"summary"`           // 概念说明
	Points  []string    `json:"key_points"`        // 要点列表
	Data    interface{} `json:"data"`              // 演示数据
	Tips    []string    `json:"tips,omitempty"`    // 额外提示
}

// ErrorResponse 统一错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

// WriteJSON 写出JSON响应 —— Go惯例: 函数而非方法
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v) //nolint
}

// WriteError 写出错误响应
func WriteError(w http.ResponseWriter, status int, errType, message string) {
	WriteJSON(w, status, ErrorResponse{
		Code:    status,
		Error:   errType,
		Message: message,
	})
}
