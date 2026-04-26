// Go学习后台系统入口
// 对比Java: public static void main(String[] args) in a class
// Go的main()必须在package main中，且函数名就是main
package main

import (
	"fmt"
	"log"
	"net/http"

	"gogo/api"
)

func main() {
	router := api.NewRouter()

	addr := ":8080"
	fmt.Println("================================================")
	fmt.Println("  Go语言学习后台系统")
	fmt.Println("  专为Java开发者设计")
	fmt.Println("================================================")
	fmt.Printf("  服务地址: http://localhost%s\n", addr)
	fmt.Printf("  API文档:  http://localhost%s/\n", addr)
	fmt.Printf("  健康检查: http://localhost%s/health\n", addr)
	fmt.Println("  访问 GET / 查看所有学习端点")
	fmt.Println("================================================")

	// log.Fatal = 打印错误后调用os.Exit(1)
	// Java: SpringApplication.run(App.class, args) 类比
	log.Fatal(http.ListenAndServe(addr, router))
}
