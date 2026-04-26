// package structs 演示Go结构体、接口、方法、嵌入
// 公共类型定义，供各个handler文件使用
package structs

import "fmt"

// ── 基础结构体 ────────────────────────────────────────────────
// Java: public class Person { private String name; ... }
// Go: 没有class，用struct；首字母大小写控制访问权限

type Person struct {
	Name  string // 大写 = public (exported)
	Age   int    // 大写 = public
	email string // 小写 = package-private (unexported)
}

func NewPerson(name string, age int, email string) *Person {
	return &Person{Name: name, Age: age, email: email}
}

func (p Person) Greet() string {
	return fmt.Sprintf("Hi, I'm %s, age %d", p.Name, p.Age)
}

// 指针接收者：可修改结构体字段
func (p *Person) SetAge(age int) {
	p.Age = age
}

func (p Person) String() string {
	return fmt.Sprintf("Person{Name:%s, Age:%d}", p.Name, p.Age)
}

// ── 形状接口（用于接口演示）──────────────────────────────────
// Java: public interface Shape { double area(); }
// Go: 隐式实现，不需要 implements 关键字！
type Shape interface {
	Area() float64
	Perimeter() float64
	Describe() string
}

// ── 嵌入演示用类型 ────────────────────────────────────────────
type Animal struct {
	Name string
}

func (a Animal) Breathe() string { return fmt.Sprintf("%s: 呼吸中", a.Name) }
func (a Animal) Eat(food string) string { return fmt.Sprintf("%s: 吃%s", a.Name, food) }

// Dog嵌入Animal，获得其所有方法（类似继承，但Go称为组合）
type Dog struct {
	Animal
	Breed string
}

func (d Dog) Breathe() string { // 覆盖Animal的Breathe
	return fmt.Sprintf("%s(狗/%s): 用鼻子呼吸", d.Name, d.Breed)
}
func (d Dog) Fetch() string { return fmt.Sprintf("%s: 取回物品！", d.Name) }

type Cat struct {
	Animal
	Indoor bool
}

func (c Cat) Purr() string {
	label := "室外猫"
	if c.Indoor {
		label = "室内猫"
	}
	return fmt.Sprintf("%s: 呼噜呼噜(%s)", c.Name, label)
}

// Swimmer接口
type Swimmer interface {
	Swim() string
}

type Duck struct {
	Animal
}

func (d Duck) Swim() string  { return fmt.Sprintf("%s: 游泳中~", d.Name) }
func (d Duck) Quack() string { return fmt.Sprintf("%s: 嘎嘎嘎!", d.Name) }

// Logger 演示多嵌入
type Logger struct {
	Prefix string
}

func (l Logger) Log(msg string) string {
	return fmt.Sprintf("[%s] %s", l.Prefix, msg)
}

// Service 演示同时嵌入多个类型（类似Java多接口+代码复用）
type Service struct {
	Animal
	Logger
	Version string
}

// ── 编译时接口检查（可选但推荐）─────────────────────────────
// Java: 编译器自动检查implements
// Go: 需要显式加这行才能在编译时发现类型未满足接口
var _ Swimmer = Duck{}
