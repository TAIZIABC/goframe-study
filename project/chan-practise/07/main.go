package main

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

func test() {
	var s = "test"

	for i, r := range s {
		fmt.Println(i, r)
	}

	a := strings.Split(s, "")

	a[len(a)-1] = "123"
	b := append(a, "123")
	var c []string
	d := copy(b, c)
	// for _, v := range b {
	// 	fmt.Println(v)
	// }

	fmt.Println(slices.Equal(a, b))
	fmt.Println(slices.Equal(c, b))
	fmt.Println(d)

	fmt.Println("Hello, World!")
	fmt.Println(a)
}

func isPalindrome(s string) bool {
	arr := strings.Split(s, "")
	fmt.Println(arr)
	slices.Reverse(arr)
	return s == strings.Join(arr, "")
}

func Compress(s string) string {
	// 在此编写你的代码
	if len(s) == 0 {
		return ""
	}
	var b strings.Builder
	count := 1
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			count++
		} else {
			b.WriteByte(s[i-1])
			if count > 1 {
				b.WriteString(strconv.Itoa(count))
			}
			count = 1
		}
	}
	b.WriteByte(s[len(s)-1])
	if count > 1 {
		b.WriteString(strconv.Itoa(count))
	}
	return b.String()
}

func CountVowels(s string) int {
	// 在此编写你的代码
	vowels := "aeiouAEIOU"
	count := 0
	for _, c := range s {
		if strings.ContainsRune(vowels, c) {
			count++
		}
	}
	return count
}

func Unique(s []int) []int {
	// 在此编写你的代码
	var result []int
	m := make(map[int]bool)
	for _, v := range s {
		if _, ok := m[v]; !ok {
			m[v] = true
			result = append(result, v)
		}
	}
	return result
}

func Zip(a, b []int) [][2]int {
	// 在此编写你的代码
	result := make([][2]int, 0)
	maxLen := max(len(b), len(a))
	for i := range maxLen {
		var va, vb int
		if i < len(a) {
			va = a[i]
		}
		if i < len(b) {
			vb = b[i]
		}
		result = append(result, [2]int{va, vb})
	}

	return result
}

// GroupBy 按字符串长度分组
func GroupBy(words []string) map[int][]string {
	// 在此编写你的代码
	result := make(map[int][]string)
	for _, v := range words {
		result[len(v)] = append(result[len(v)], v)
	}

	return result
}

// InvertMap 键值互换
func InvertMap(m map[string]int) map[int]string {
	// 在此编写你的代码
	result := make(map[int]string)
	for k, v := range m {
		result[v] = k
	}

	return result
}

func MergeMaps(maps ...map[string]int) map[string]int {
	// 在此编写你的代码
	result := make(map[string]int)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// FrequencySort 按频率排序（高频在前，频率相同则数值小的在前）
func FrequencySort(s []int) []int {
	// 在此编写你的代码
	result := make([]int, len(s))
	m := make(map[int]int)
	for _, v := range s {
		m[v]++
	}

	copy(result, s)
	slices.SortFunc(result, func(a, b int) int {
		fmt.Println(a, b)
		if m[a] == m[b] {
			return a - b
		}
		return m[b] - m[a]
	})
	fmt.Println(m)
	return result
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum
}

func testChan() {
	c := make(chan int)
	go sum([]int{1, 2, 3, 4}, c)
	go sum([]int{1, 2, 3, 4}, c)
	x, y := <-c, <-c
	fmt.Println(x, y, x+y)
}

// TitleCase 每个单词首字母大写，其余小写
func TitleCase(s string) string {
	// 在此编写你的代码
	result := ""
	temp := strings.Fields(s)
	for _, v := range temp {
		result += strings.ToUpper(v[:1]) + strings.ToLower(v[1:]) + " "
	}
	result = result[:len(result)-1]
	return result
}

func MaskEmail(email string) string {
	// 在此编写你的代码
	if len(email) == 0 {
		return email
	}
	index := strings.Index(email, "@")
	if index == -1 {
		return email
	}
	preStr := email[:index]
	if len(preStr) <= 2 {
		return preStr[:1] + "***" + email[index:]
	} else {
		return preStr[:1] + "***" + preStr[len(preStr)-1:] + email[index:]
	}
}

func SafeAtoi(s string, defaultVal int) int {
	// 在此编写你的代码
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}

type Student struct {
	Name  string
	Score int
}

func SortByScore(students []Student) []Student {
	// 在此编写你的代码
	slices.SortFunc(students, func(a, b Student) int {
		if a.Score == b.Score {
			return strings.Compare(a.Name, b.Name)
		}
		return b.Score - a.Score
	})

	return students
}

// TopK 返回前K大的元素（降序）
func TopK(s []int, k int) []int {
	// 在此编写你的代码
	if k <= 0 || len(s) == 0 {
		return []int{}
	}

	copys := slices.Clone(s)
	slices.SortFunc(copys, func(a, b int) int {
		return b - a
	})
	if k > len(s) {
		k = len(s)
	}
	return copys[:k]
}

// ToJSON 序列化为 JSON 字符串
func ToJSON(v any) string {
	// 在此编写你的代码
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// FromJSON 反序列化 JSON 字符串
func FromJSON(data string, v any) error {
	// 在此编写你的代码
	return json.Unmarshal([]byte(data), v)
}

// FormatDuration 人性化时间描述
// 规则：
// - < 1秒 → "刚刚"
// - < 1分钟 → "X秒前"
// - < 1小时 → "X分钟前"
// - < 24小时 → "X小时前"
// - 否则 → "X天前"
func FormatDuration(d time.Duration) string {
	// 在此编写你的代码
	if d < time.Second {
		return "刚刚"
	}
	if d < time.Minute {
		return fmt.Sprintf("%ds前", d/time.Second)
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm前", d/time.Minute)
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh前", d/time.Hour)
	}
	return ""
}

func ToIntSlice(s []string) ([]int, error) {
	// 在此编写你的代码
	result := make([]int, len(s))
	var err error
	for i, v := range s {
		result[i], err = strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type Students struct {
	Name  string
	Age   int
	Score float64
}

func ParseCSVLine(line string) (Students, error) {
	// 在此编写你的代码
	result := strings.Split(line, ",")

	if len(result) != 3 {
		return Students{}, fmt.Errorf("invalid line")
	}
	age, err := strconv.Atoi(result[1])
	if err != nil {
		return Students{}, fmt.Errorf("invalid age")
	}
	score, err := strconv.ParseFloat(result[2], 64)
	if err != nil {
		return Students{}, fmt.Errorf("invalid score")
	}
	return Students{
		Name:  result[0],
		Age:   age,
		Score: score,
	}, nil
}

func Generate(nums ...int) <-chan int {
	// 在此编写
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func Square(in <-chan int) <-chan int {
	// 在此编写
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func Sum(in <-chan int) int {
	// 在此编写
	total := 0
	for n := range in {
		total += n
	}
	return total
}

func main() {

	// var s = "aaaccssssgd"
	// var m map[string]string
	// fmt.Println(len(s))
	// fmt.Println(isPalindrome(s))
	// fmt.Println(Compress(s))
	// fmt.Println(Zip([]int{1, 2, 3}, []int{4, 5}))
	// fmt.Println(FrequencySort([]int{1, 1, 2, 2, 2, 3}))
	// testChan()
	// fmt.Println(TitleCase("hello WORLD go"))
	// fmt.Println(MaskEmail("alice@example.com"))
	// fmt.Println(SortByScore([]Student{{"Bob", 90}, {"Alice", 90}, {"Charlie", 85}}))
	// fmt.Println(ToJSON(map[string]int{"a": 1}))
	// fmt.Println(FromJSON(`{"name":"Go"}`, &m), m)
	// fmt.Println(Sum(Square(Generate(1, 2, 3)))) // 14
	for i, v := range []int{1, 2, 3} {
		go func() { fmt.Println(i, v) }()
		time.Sleep(time.Second)
	}

}
