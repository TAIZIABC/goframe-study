package main

import (
	"fmt"
	"sort"
)

// 题目：实现多条件结构体排序
//
// 给定 Employee 结构体 { Name string, Dept string, Salary int }
// 实现 SortEmployees(employees []Employee) []Employee
// 排序规则：
//   1. 按 Dept 升序
//   2. 同部门按 Salary 降序
//   3. 同薪资按 Name 升序

type Employee struct {
	Name   string
	Dept   string
	Salary int
}

func SortEmployees(employees []Employee) []Employee {
	sort.Slice(employees, func(i, j int) bool {
		if employees[i].Dept != employees[j].Dept {
			return employees[i].Dept < employees[j].Dept
		}
		if employees[i].Salary != employees[j].Salary {
			return employees[i].Salary > employees[j].Salary
		}
		return employees[i].Name < employees[j].Name
	})
	return employees
}

func main() {
	emps := []Employee{
		{"Alice", "Engineering", 100},
		{"Bob", "Engineering", 120},
		{"Charlie", "Sales", 90},
		{"Dave", "Engineering", 120},
		{"Eve", "Sales", 95},
	}
	sorted := SortEmployees(emps)
	for _, e := range sorted {
		fmt.Printf("%-10s %-12s %d\n", e.Name, e.Dept, e.Salary)
	}
	// Bob        Engineering  120
	// Dave       Engineering  120
	// Alice      Engineering  100
	// Eve        Sales        95
	// Charlie    Sales        90
}
