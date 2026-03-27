/* =============================================
   GoFrame 学习之路 - 交互脚本
   ============================================= */

// ==================== 课程数据 ====================
const LESSONS = {
  variables: {
    icon: '📦', title: '变量声明与类型系统',
    desc: 'Go 是静态类型语言，这是与 JavaScript 最大的区别之一。类型在编译时确定，更早发现错误、更好的性能。',
    points: [
      ['静态类型', '变量类型一旦确定不可改变，类似 TypeScript 但更严格'],
      ['类型推断', '使用 <code>:=</code> 短声明可自动推断类型'],
      ['零值机制', '所有变量都有默认零值，不存在 undefined'],
      ['值类型', '基本类型赋值是值拷贝，不是引用'],
    ],
    jsCode: `<span class="comment">// JavaScript - 动态类型</span>
<span class="keyword">let</span> name = <span class="string">"GoFrame"</span>;     <span class="comment">// 可随时改变类型</span>
<span class="keyword">let</span> age = <span class="number">25</span>;
<span class="keyword">let</span> isActive = <span class="keyword">true</span>;
<span class="keyword">let</span> scores = [<span class="number">90</span>, <span class="number">85</span>, <span class="number">92</span>];
<span class="keyword">let</span> user = { name: <span class="string">"张三"</span>, age: <span class="number">25</span> };

<span class="comment">// TypeScript - 静态类型（更接近 Go）</span>
<span class="keyword">let</span> name: <span class="type">string</span> = <span class="string">"GoFrame"</span>;
<span class="keyword">const</span> PI = <span class="number">3.14159</span>;

<span class="comment">// undefined 和 null</span>
<span class="keyword">let</span> data;  <span class="comment">// undefined</span>
<span class="keyword">let</span> empty = <span class="keyword">null</span>;`,
    goCode: `<span class="comment">// Go - 静态类型，编译时检查</span>
<span class="keyword">var</span> name <span class="type">string</span> = <span class="string">"GoFrame"</span>  <span class="comment">// 显式声明</span>
name := <span class="string">"GoFrame"</span>             <span class="comment">// 短声明（推荐）</span>
age := <span class="number">25</span>                      <span class="comment">// 推断为 int</span>
isActive := <span class="keyword">true</span>              <span class="comment">// 推断为 bool</span>
scores := []<span class="type">int</span>{<span class="number">90</span>, <span class="number">85</span>, <span class="number">92</span>}  <span class="comment">// 切片</span>

<span class="comment">// 结构体（替代对象）</span>
<span class="keyword">type</span> <span class="type">User</span> <span class="keyword">struct</span> {
    Name <span class="type">string</span>
    Age  <span class="type">int</span>
}
user := <span class="type">User</span>{Name: <span class="string">"张三"</span>, Age: <span class="number">25</span>}

<span class="comment">// 零值（没有 undefined/null）</span>
<span class="keyword">var</span> s <span class="type">string</span>  <span class="comment">// 零值: ""</span>
<span class="keyword">var</span> n <span class="type">int</span>     <span class="comment">// 零值: 0</span>`,
    editorCode: `package main

import "fmt"

func main() {
    // 试试修改这些变量
    name := "GoFrame 学习者"
    age := 25
    isLearning := true
    
    fmt.Printf("姓名: %s\\n", name)
    fmt.Printf("年龄: %d\\n", age)
    fmt.Printf("正在学习: %v\\n", isLearning)
    
    // 切片（类似 JS 数组）
    skills := []string{"JavaScript", "TypeScript", "Go"}
    fmt.Printf("技能: %v\\n", skills)
    
    // 零值
    var score int
    var title string
    fmt.Printf("默认分数: %d\\n", score)
    fmt.Printf("默认标题: '%s'\\n", title)
}`,
  },
  functions: {
    icon: '⚡', title: '函数与方法',
    desc: 'Go 函数支持多返回值，这是其最强大的特性之一。方法是绑定到特定类型上的函数。',
    points: [
      ['多返回值', 'Go 函数可返回多个值，常用于返回结果+错误'],
      ['命名返回值', '可以给返回值命名，增加可读性'],
      ['方法接收者', '方法通过接收者绑定到类型上，类似 this'],
      ['大小写控制', '大写 = 导出(public)，小写 = 私有(private)'],
    ],
    jsCode: `<span class="comment">// 函数声明</span>
<span class="keyword">function</span> <span class="function">add</span>(a, b) {
  <span class="keyword">return</span> a + b;
}

<span class="comment">// 箭头函数</span>
<span class="keyword">const</span> <span class="function">multiply</span> = (a, b) => a * b;

<span class="comment">// 返回多个值（用对象模拟）</span>
<span class="keyword">function</span> <span class="function">divide</span>(a, b) {
  <span class="keyword">if</span> (b === <span class="number">0</span>) {
    <span class="keyword">return</span> { result: <span class="number">0</span>, error: <span class="string">"除数不能为0"</span> };
  }
  <span class="keyword">return</span> { result: a / b, error: <span class="keyword">null</span> };
}

<span class="comment">// 类方法</span>
<span class="keyword">class</span> <span class="type">Calculator</span> {
  <span class="function">constructor</span>(v) { <span class="keyword">this</span>.value = v; }
  <span class="function">add</span>(n) { <span class="keyword">this</span>.value += n; <span class="keyword">return</span> <span class="keyword">this</span>; }
}`,
    goCode: `<span class="comment">// 函数声明（参数和返回值都有类型）</span>
<span class="keyword">func</span> <span class="function">add</span>(a, b <span class="type">int</span>) <span class="type">int</span> {
    <span class="keyword">return</span> a + b
}

<span class="comment">// 多返回值（Go 的精髓！）</span>
<span class="keyword">func</span> <span class="function">divide</span>(a, b <span class="type">float64</span>) (<span class="type">float64</span>, <span class="type">error</span>) {
    <span class="keyword">if</span> b == <span class="number">0</span> {
        <span class="keyword">return</span> <span class="number">0</span>, fmt.<span class="function">Errorf</span>(<span class="string">"除数不能为0"</span>)
    }
    <span class="keyword">return</span> a / b, <span class="keyword">nil</span>
}

<span class="comment">// 方法（绑定到结构体）</span>
<span class="keyword">type</span> <span class="type">Calculator</span> <span class="keyword">struct</span> { Value <span class="type">float64</span> }

<span class="keyword">func</span> (c *<span class="type">Calculator</span>) <span class="function">Add</span>(n <span class="type">float64</span>) *<span class="type">Calculator</span> {
    c.Value += n
    <span class="keyword">return</span> c
}

<span class="comment">// 闭包</span>
<span class="keyword">func</span> <span class="function">counter</span>() <span class="keyword">func</span>() <span class="type">int</span> {
    count := <span class="number">0</span>
    <span class="keyword">return</span> <span class="keyword">func</span>() <span class="type">int</span> { count++; <span class="keyword">return</span> count }
}`,
    editorCode: `package main

import (
    "errors"
    "fmt"
)

// 多返回值函数
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("除数不能为零")
    }
    return a / b, nil
}

// 可变参数（类似 JS 的 ...args）
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

func main() {
    result, err := divide(10, 3)
    if err != nil {
        fmt.Println("错误:", err)
    } else {
        fmt.Printf("10 / 3 = %.2f\\n", result)
    }
    
    _, err = divide(10, 0)
    if err != nil {
        fmt.Println("错误:", err)
    }
    
    fmt.Printf("求和: %d\\n", sum(1, 2, 3, 4, 5))
}`,
  },
  structs: {
    icon: '🏗️', title: '结构体 vs JavaScript 对象/类',
    desc: 'Go 没有类（class），用结构体（struct）+ 方法来实现面向对象编程。',
    points: [
      ['没有继承', 'Go 用"组合"代替"继承"，更灵活更简洁'],
      ['结构体标签', '可为字段添加 JSON、数据库等标签'],
      ['值类型', '结构体赋值默认是拷贝，用指针共享数据'],
      ['嵌入组合', '通过嵌入其他结构体实现代码复用'],
    ],
    jsCode: `<span class="comment">// JavaScript 类</span>
<span class="keyword">class</span> <span class="type">User</span> {
  <span class="function">constructor</span>(name, email, age) {
    <span class="keyword">this</span>.name = name;
    <span class="keyword">this</span>.email = email;
    <span class="keyword">this</span>.age = age;
  }
  <span class="function">greet</span>() {
    <span class="keyword">return</span> <span class="string">\`你好，我是\${this.name}\`</span>;
  }
}

<span class="comment">// 继承</span>
<span class="keyword">class</span> <span class="type">Admin</span> <span class="keyword">extends</span> <span class="type">User</span> {
  <span class="function">constructor</span>(name, email, age, role) {
    <span class="keyword">super</span>(name, email, age);
    <span class="keyword">this</span>.role = role;
  }
}

<span class="comment">// JSON</span>
JSON.<span class="function">stringify</span>(user);`,
    goCode: `<span class="comment">// Go 结构体（替代 class）</span>
<span class="keyword">type</span> <span class="type">User</span> <span class="keyword">struct</span> {
    Name  <span class="type">string</span> <span class="string">\`json:"name"\`</span>
    Email <span class="type">string</span> <span class="string">\`json:"email"\`</span>
    Age   <span class="type">int</span>    <span class="string">\`json:"age"\`</span>
}

<span class="keyword">func</span> (u *<span class="type">User</span>) <span class="function">Greet</span>() <span class="type">string</span> {
    <span class="keyword">return</span> fmt.<span class="function">Sprintf</span>(<span class="string">"你好，我是%s"</span>, u.Name)
}

<span class="comment">// 组合代替继承</span>
<span class="keyword">type</span> <span class="type">Admin</span> <span class="keyword">struct</span> {
    <span class="type">User</span>           <span class="comment">// 嵌入 User</span>
    Role <span class="type">string</span>
}

<span class="comment">// 构造函数（惯例）</span>
<span class="keyword">func</span> <span class="function">NewUser</span>(n, e <span class="type">string</span>, a <span class="type">int</span>) *<span class="type">User</span> {
    <span class="keyword">return</span> &<span class="type">User</span>{Name: n, Email: e, Age: a}
}`,
    editorCode: `package main

import (
    "encoding/json"
    "fmt"
)

type User struct {
    Name  string \`json:"name"\`
    Email string \`json:"email"\`
    Age   int    \`json:"age"\`
}

func (u *User) Greet() string {
    return fmt.Sprintf("你好！我是 %s，今年 %d 岁", u.Name, u.Age)
}

type Admin struct {
    User
    Role string \`json:"role"\`
}

func main() {
    user := User{Name: "前端开发者", Email: "dev@go.dev", Age: 28}
    fmt.Println(user.Greet())
    
    jsonData, _ := json.MarshalIndent(user, "", "  ")
    fmt.Printf("JSON:\\n%s\\n", jsonData)
    
    admin := Admin{
        User: User{Name: "管理员", Email: "admin@go.dev", Age: 30},
        Role: "super_admin",
    }
    fmt.Println(admin.Greet())
    fmt.Printf("角色: %s\\n", admin.Role)
}`,
  },
  errors: {
    icon: '🛡️', title: '错误处理：告别 try/catch',
    desc: 'Go 的错误处理哲学是"错误是值"，使用显式的错误返回代替异常机制。',
    points: [
      ['错误是值', 'error 是一个接口类型，不是异常'],
      ['显式处理', '每个可能出错的操作都需要检查错误'],
      ['自定义错误', '可创建自己的错误类型，携带更多信息'],
      ['errors 包', 'Go 1.13+ 引入了错误包装和检查机制'],
    ],
    jsCode: `<span class="comment">// try/catch</span>
<span class="keyword">try</span> {
  <span class="keyword">const</span> data = JSON.<span class="function">parse</span>(rawData);
  <span class="keyword">const</span> result = <span class="keyword">await</span> <span class="function">fetchUser</span>(data.id);
} <span class="keyword">catch</span> (err) {
  console.<span class="function">error</span>(<span class="string">"出错了:"</span>, err.message);
}

<span class="comment">// Promise 链</span>
<span class="function">fetch</span>(<span class="string">"/api/users"</span>)
  .<span class="function">then</span>(res => res.<span class="function">json</span>())
  .<span class="function">catch</span>(err => console.<span class="function">error</span>(err));

<span class="comment">// 自定义错误</span>
<span class="keyword">class</span> <span class="type">NotFoundError</span> <span class="keyword">extends</span> <span class="type">Error</span> {
  <span class="function">constructor</span>(id) {
    <span class="keyword">super</span>(<span class="string">\`用户 \${id} 未找到\`</span>);
    <span class="keyword">this</span>.statusCode = <span class="number">404</span>;
  }
}`,
    goCode: `<span class="comment">// 显式错误检查（Go 哲学）</span>
data, err := json.<span class="function">Unmarshal</span>(rawData, &obj)
<span class="keyword">if</span> err != <span class="keyword">nil</span> {
    log.<span class="function">Printf</span>(<span class="string">"解析失败: %v"</span>, err)
    <span class="keyword">return</span>
}

<span class="comment">// 自定义错误类型</span>
<span class="keyword">type</span> <span class="type">NotFoundError</span> <span class="keyword">struct</span> {
    ID         <span class="type">int</span>
    StatusCode <span class="type">int</span>
}
<span class="keyword">func</span> (e *<span class="type">NotFoundError</span>) <span class="function">Error</span>() <span class="type">string</span> {
    <span class="keyword">return</span> fmt.<span class="function">Sprintf</span>(<span class="string">"用户 %d 未找到"</span>, e.ID)
}

<span class="comment">// 错误包装 (Go 1.13+)</span>
err = fmt.<span class="function">Errorf</span>(<span class="string">"查询失败: %w"</span>, origErr)
<span class="keyword">if</span> errors.<span class="function">Is</span>(err, sql.ErrNoRows) {
    <span class="comment">// 处理特定错误</span>
}`,
    editorCode: `package main

import (
    "errors"
    "fmt"
    "strconv"
)

type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("验证失败 [%s]: %s", e.Field, e.Message)
}

func parseAge(input string) (int, error) {
    age, err := strconv.Atoi(input)
    if err != nil {
        return 0, fmt.Errorf("无法解析年龄: %w", err)
    }
    if age < 0 || age > 150 {
        return 0, &ValidationError{Field: "age", Message: fmt.Sprintf("年龄 %d 不合理", age)}
    }
    return age, nil
}

func main() {
    age, err := parseAge("25")
    if err != nil {
        fmt.Println("错误:", err)
    } else {
        fmt.Printf("年龄: %d\\n", age)
    }
    
    _, err = parseAge("abc")
    if err != nil { fmt.Println("错误:", err) }
    
    _, err = parseAge("-5")
    if err != nil {
        var ve *ValidationError
        if errors.As(err, &ve) {
            fmt.Printf("验证错误 - 字段: %s, 消息: %s\\n", ve.Field, ve.Message)
        }
    }
}`,
  },
  concurrency: {
    icon: '🔄', title: '并发编程：Goroutine vs async/await',
    desc: 'Go 的并发模型基于 CSP，使用 goroutine 和 channel 实现。支持真正的并行执行。',
    points: [
      ['Goroutine', '轻量级协程，约 2KB 栈空间，可创建数百万个'],
      ['Channel', 'goroutine 之间的通信管道'],
      ['select', '监听多个 channel，类似 Promise.race'],
      ['sync 包', '提供 WaitGroup、Mutex 等同步原语'],
    ],
    jsCode: `<span class="comment">// async/await</span>
<span class="keyword">async function</span> <span class="function">fetchData</span>() {
  <span class="keyword">const</span> res = <span class="keyword">await</span> <span class="function">fetch</span>(<span class="string">"/api"</span>);
  <span class="keyword">return</span> res.<span class="function">json</span>();
}

<span class="comment">// Promise.all 并发</span>
<span class="keyword">const</span> [users, posts] = <span class="keyword">await</span> Promise.<span class="function">all</span>([
  <span class="function">fetchUsers</span>(),
  <span class="function">fetchPosts</span>()
]);

<span class="comment">// Promise.race 竞争</span>
<span class="keyword">const</span> result = <span class="keyword">await</span> Promise.<span class="function">race</span>([
  <span class="function">fetchData</span>(),
  <span class="function">timeout</span>(<span class="number">5000</span>)
]);

<span class="comment">// 注意：JS 是单线程的！</span>`,
    goCode: `<span class="comment">// goroutine（真正的并行！）</span>
<span class="keyword">go</span> <span class="keyword">func</span>() {
    data := <span class="function">fetchData</span>()
    fmt.<span class="function">Println</span>(data)
}()

<span class="comment">// Channel 通信</span>
ch := <span class="keyword">make</span>(<span class="keyword">chan</span> <span class="type">string</span>)
<span class="keyword">go</span> <span class="keyword">func</span>() { ch <- <span class="function">fetchUsers</span>() }()
users := <-ch

<span class="comment">// WaitGroup（类似 Promise.all）</span>
<span class="keyword">var</span> wg sync.<span class="type">WaitGroup</span>
wg.<span class="function">Add</span>(<span class="number">2</span>)
<span class="keyword">go</span> <span class="keyword">func</span>() { <span class="keyword">defer</span> wg.<span class="function">Done</span>(); <span class="function">task1</span>() }()
<span class="keyword">go</span> <span class="keyword">func</span>() { <span class="keyword">defer</span> wg.<span class="function">Done</span>(); <span class="function">task2</span>() }()
wg.<span class="function">Wait</span>()

<span class="comment">// select（类似 Promise.race）</span>
<span class="keyword">select</span> {
<span class="keyword">case</span> data := <-dataCh:
    fmt.<span class="function">Println</span>(data)
<span class="keyword">case</span> <-time.<span class="function">After</span>(<span class="number">5</span> * time.Second):
    fmt.<span class="function">Println</span>(<span class="string">"超时"</span>)
}`,
    editorCode: `package main

import (
    "fmt"
    "sync"
    "time"
)

func fetchAPI(name string, delay time.Duration) string {
    time.Sleep(delay)
    return fmt.Sprintf("%s 的数据", name)
}

func main() {
    start := time.Now()
    
    var wg sync.WaitGroup
    results := make([]string, 3)
    
    apis := []struct{ name string; delay time.Duration; idx int }{
        {"用户API", 100 * time.Millisecond, 0},
        {"文章API", 150 * time.Millisecond, 1},
        {"评论API", 80 * time.Millisecond, 2},
    }
    
    for _, api := range apis {
        wg.Add(1)
        go func(n string, d time.Duration, i int) {
            defer wg.Done()
            results[i] = fetchAPI(n, d)
            fmt.Printf("✅ %s 完成\\n", n)
        }(api.name, api.delay, api.idx)
    }
    
    wg.Wait()
    fmt.Printf("\\n全部完成！耗时: %v\\n", time.Since(start))
    fmt.Println("串行需要: 330ms")
    for _, r := range results {
        fmt.Printf("  📦 %s\\n", r)
    }
}`,
  },
  interfaces: {
    icon: '🔌', title: '接口：隐式实现的魔力',
    desc: 'Go 的接口是隐式实现的——只要类型拥有接口定义的所有方法，就自动实现了该接口。',
    points: [
      ['隐式实现', '不需要 implements 关键字'],
      ['小接口', 'Go 推崇 1-3 个方法的小接口'],
      ['空接口', '<code>any</code> 类似 TypeScript 的 any'],
      ['接口组合', '大接口可由小接口组合而成'],
    ],
    jsCode: `<span class="comment">// TypeScript 接口</span>
<span class="keyword">interface</span> <span class="type">Shape</span> {
  <span class="function">area</span>(): <span class="type">number</span>;
  <span class="function">perimeter</span>(): <span class="type">number</span>;
}

<span class="keyword">class</span> <span class="type">Circle</span> <span class="keyword">implements</span> <span class="type">Shape</span> {
  <span class="function">constructor</span>(<span class="keyword">private</span> r: <span class="type">number</span>) {}
  <span class="function">area</span>() {
    <span class="keyword">return</span> Math.PI * <span class="keyword">this</span>.r ** <span class="number">2</span>;
  }
  <span class="function">perimeter</span>() {
    <span class="keyword">return</span> <span class="number">2</span> * Math.PI * <span class="keyword">this</span>.r;
  }
}

<span class="comment">// Duck Typing</span>
<span class="keyword">interface</span> <span class="type">Printable</span> {
  <span class="function">toString</span>(): <span class="type">string</span>;
}`,
    goCode: `<span class="comment">// Go 接口（隐式实现！）</span>
<span class="keyword">type</span> <span class="type">Shape</span> <span class="keyword">interface</span> {
    <span class="function">Area</span>() <span class="type">float64</span>
    <span class="function">Perimeter</span>() <span class="type">float64</span>
}

<span class="keyword">type</span> <span class="type">Circle</span> <span class="keyword">struct</span> { Radius <span class="type">float64</span> }

<span class="comment">// 自动满足 Shape 接口</span>
<span class="keyword">func</span> (c <span class="type">Circle</span>) <span class="function">Area</span>() <span class="type">float64</span> {
    <span class="keyword">return</span> math.Pi * c.Radius * c.Radius
}
<span class="keyword">func</span> (c <span class="type">Circle</span>) <span class="function">Perimeter</span>() <span class="type">float64</span> {
    <span class="keyword">return</span> <span class="number">2</span> * math.Pi * c.Radius
}

<span class="comment">// 接口组合</span>
<span class="keyword">type</span> <span class="type">ReadWriter</span> <span class="keyword">interface</span> {
    <span class="type">io.Reader</span>
    <span class="type">io.Writer</span>
}`,
    editorCode: `package main

import (
    "fmt"
    "math"
)

type Shape interface {
    Area() float64
    Perimeter() float64
    String() string
}

type Circle struct{ Radius float64 }
func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }
func (c Circle) String() string     { return fmt.Sprintf("圆(r=%.1f)", c.Radius) }

type Rect struct{ W, H float64 }
func (r Rect) Area() float64      { return r.W * r.H }
func (r Rect) Perimeter() float64 { return 2 * (r.W + r.H) }
func (r Rect) String() string     { return fmt.Sprintf("矩形(%.1f×%.1f)", r.W, r.H) }

func main() {
    shapes := []Shape{
        Circle{5}, Rect{4, 6}, Circle{3}, Rect{10, 2},
    }
    total := 0.0
    for _, s := range shapes {
        fmt.Printf("  %s => 面积: %.2f\\n", s, s.Area())
        total += s.Area()
    }
    fmt.Printf("\\n总面积: %.2f\\n", total)
}`,
  },
  pointers: {
    icon: '🎯', title: '指针：前端没有的概念',
    desc: 'JavaScript 没有指针概念，但 Go 中指针非常重要。指针存储的是变量的内存地址，理解它才能掌握值传递 vs 引用传递。',
    points: [
      ['什么是指针', '<code>*T</code> 表示指向 T 类型的指针，<code>&</code> 取地址，<code>*</code> 解引用'],
      ['值传递', 'Go 函数参数默认值传递（拷贝），JS 对象是引用传递'],
      ['指针接收者', '方法用 <code>*T</code> 接收者才能修改原始数据'],
      ['没有指针运算', 'Go 的指针比 C 安全，不能做算术运算'],
    ],
    jsCode: `<span class="comment">// JS 基本类型 = 值拷贝</span>
<span class="keyword">let</span> a = <span class="number">10</span>;
<span class="keyword">let</span> b = a;   <span class="comment">// b 是 a 的副本</span>
b = <span class="number">20</span>;       <span class="comment">// a 仍然是 10</span>

<span class="comment">// JS 对象 = 引用传递</span>
<span class="keyword">const</span> user = { name: <span class="string">"张三"</span> };
<span class="keyword">const</span> ref = user;  <span class="comment">// ref 和 user 指向同一对象</span>
ref.name = <span class="string">"李四"</span>;  <span class="comment">// user.name 也变了！</span>

<span class="comment">// 函数参数：对象是引用</span>
<span class="keyword">function</span> <span class="function">changeName</span>(u) {
  u.name = <span class="string">"王五"</span>;  <span class="comment">// 会修改原始对象</span>
}
<span class="function">changeName</span>(user);
<span class="comment">// user.name === "王五"</span>

<span class="comment">// 没有办法让函数修改基本类型</span>
<span class="keyword">function</span> <span class="function">tryDouble</span>(n) {
  n = n * <span class="number">2</span>;  <span class="comment">// 不影响外部变量</span>
}`,
    goCode: `<span class="comment">// 指针基础</span>
x := <span class="number">42</span>
p := &x      <span class="comment">// p 是指向 x 的指针</span>
fmt.<span class="function">Println</span>(*p)  <span class="comment">// 42（解引用）</span>
*p = <span class="number">100</span>     <span class="comment">// 通过指针修改 x</span>
<span class="comment">// x 现在是 100</span>

<span class="comment">// 值传递 vs 指针传递</span>
<span class="keyword">func</span> <span class="function">doubleVal</span>(n <span class="type">int</span>) { n *= <span class="number">2</span> }
<span class="keyword">func</span> <span class="function">doublePtr</span>(n *<span class="type">int</span>) { *n *= <span class="number">2</span> }

x := <span class="number">10</span>
<span class="function">doubleVal</span>(x)  <span class="comment">// x 仍是 10（值拷贝）</span>
<span class="function">doublePtr</span>(&x) <span class="comment">// x 变为 20（指针修改）</span>

<span class="comment">// 结构体指针（最常用场景）</span>
<span class="keyword">type</span> <span class="type">User</span> <span class="keyword">struct</span> { Name <span class="type">string</span> }

<span class="comment">// 值接收者：不修改原始数据</span>
<span class="keyword">func</span> (u <span class="type">User</span>) <span class="function">GetName</span>() <span class="type">string</span> { <span class="keyword">return</span> u.Name }

<span class="comment">// 指针接收者：可以修改原始数据</span>
<span class="keyword">func</span> (u *<span class="type">User</span>) <span class="function">SetName</span>(n <span class="type">string</span>) { u.Name = n }`,
    editorCode: `package main

import "fmt"

type User struct {
    Name string
    Age  int
}

// 值接收者 —— 不会修改原始数据（类似 JS 基本类型）
func (u User) Birthday() {
    u.Age++  // 这只修改了副本！
    fmt.Printf("  [值接收者内部] Age = %d\\n", u.Age)
}

// 指针接收者 —— 会修改原始数据（类似 JS 对象引用）
func (u *User) BirthdayPtr() {
    u.Age++  // 修改的是原始数据
    fmt.Printf("  [指针接收者内部] Age = %d\\n", u.Age)
}

// 函数参数对比
func tryModifyValue(n int) {
    n = 999
}

func tryModifyPointer(n *int) {
    *n = 999
}

func main() {
    // 1. 基本类型指针
    x := 42
    p := &x
    fmt.Printf("x = %d, *p = %d\\n", x, *p)
    *p = 100
    fmt.Printf("修改 *p 后: x = %d\\n\\n", x)

    // 2. 值传递 vs 指针传递
    a := 10
    tryModifyValue(a)
    fmt.Printf("值传递后: a = %d（没变）\\n", a)
    tryModifyPointer(&a)
    fmt.Printf("指针传递后: a = %d（变了！）\\n\\n", a)

    // 3. 方法接收者对比
    user := User{Name: "前端开发者", Age: 25}
    fmt.Println("原始 Age:", user.Age)

    user.Birthday()  // 值接收者
    fmt.Println("值接收者调用后:", user.Age, "（没变）")

    user.BirthdayPtr()  // 指针接收者
    fmt.Println("指针接收者调用后:", user.Age, "（变了！）")
}`,
  },
};

// GoFrame 框架课程数据
const GF_LESSONS = [
  { id: 'gf-project', num: '01', title: '项目结构与初始化', content: `<p class="gf-intro-text">GoFrame 项目结构遵循标准 Go 布局，类似 Next.js 约定式目录结构。</p>
<div class="file-tree"><div class="tree-title">📁 GoFrame 项目结构</div><pre class="tree-content">├── api/          <span class="tree-comment"># 接口定义（类似 TS 接口文件）</span>
│   └── v1/       <span class="tree-comment"># API 版本管理</span>
├── internal/     <span class="tree-comment"># 内部代码（私有模块）</span>
│   ├── cmd/      <span class="tree-comment"># 入口命令</span>
│   ├── controller/ <span class="tree-comment"># 控制器（类似 routes/handlers）</span>
│   ├── dao/      <span class="tree-comment"># 数据访问层</span>
│   ├── logic/    <span class="tree-comment"># 业务逻辑</span>
│   ├── model/    <span class="tree-comment"># 数据模型（类似 TS types）</span>
│   └── service/  <span class="tree-comment"># 服务接口</span>
├── manifest/config/ <span class="tree-comment"># 应用配置（类似 .env）</span>
├── resource/     <span class="tree-comment"># 静态资源</span>
├── go.mod        <span class="tree-comment"># 依赖管理（= package.json）</span>
└── main.go       <span class="tree-comment"># 主入口</span></pre></div>
<div class="code-panel full-width"><div class="code-panel-header"><div class="code-lang"><span class="lang-dot go-dot"></span>初始化命令</div></div><pre class="code-block"><code><span class="comment"># 安装 GoFrame CLI（类似 npx create-react-app）</span>
$ go install github.com/gogf/gf/cmd/gf/v2@latest

<span class="comment"># 创建新项目</span>
$ gf init my-project

<span class="comment"># 启动开发服务器（类似 npm run dev）</span>
$ gf run main.go</code></pre></div>` },
  { id: 'gf-router', num: '02', title: '路由与中间件', content: `<p class="gf-intro-text">GoFrame 路由系统类似 Express，支持分组、中间件、RESTful。</p>
<div class="code-comparison"><div class="code-panel"><div class="code-panel-header"><div class="code-lang"><span class="lang-dot js-dot"></span>Express.js</div></div><pre class="code-block"><code><span class="keyword">const</span> app = <span class="function">express</span>();
app.<span class="function">use</span>(express.<span class="function">json</span>());

<span class="keyword">const</span> api = express.<span class="function">Router</span>();
api.<span class="function">get</span>(<span class="string">'/users'</span>, getUsers);
api.<span class="function">post</span>(<span class="string">'/users'</span>, createUser);
app.<span class="function">use</span>(<span class="string">'/api/v1'</span>, api);
app.<span class="function">listen</span>(<span class="number">3000</span>);</code></pre></div><div class="code-vs"><div class="vs-badge">VS</div></div><div class="code-panel"><div class="code-panel-header"><div class="code-lang"><span class="lang-dot go-dot"></span>GoFrame</div></div><pre class="code-block"><code>s := g.<span class="function">Server</span>()
s.<span class="function">Use</span>(ghttp.MiddlewareHandlerResponse)

s.<span class="function">Group</span>(<span class="string">"/api/v1"</span>, <span class="keyword">func</span>(group *ghttp.<span class="type">RouterGroup</span>) {
    group.<span class="function">Middleware</span>(authMiddleware)
    group.<span class="function">GET</span>(<span class="string">"/users"</span>, GetUsers)
    group.<span class="function">POST</span>(<span class="string">"/users"</span>, CreateUser)
})

s.<span class="function">SetPort</span>(<span class="number">8000</span>)
s.<span class="function">Run</span>()</code></pre></div></div>` },
  { id: 'gf-request', num: '03', title: '请求处理与数据校验', content: `<p class="gf-intro-text">GoFrame 提供强大的参数绑定和校验，类似 Joi/Zod 但集成在框架层面。</p>
<div class="code-panel full-width"><div class="code-panel-header"><div class="code-lang"><span class="lang-dot go-dot"></span>请求校验</div></div><pre class="code-block"><code><span class="comment">// 定义请求结构体（类似 Zod schema）</span>
<span class="keyword">type</span> <span class="type">CreateUserReq</span> <span class="keyword">struct</span> {
    g.Meta   <span class="string">\`path:"/users" method:"post"\`</span>
    Name     <span class="type">string</span> <span class="string">\`v:"required|length:2,30#请输入姓名|长度2-30"\`</span>
    Email    <span class="type">string</span> <span class="string">\`v:"required|email#请输入邮箱|格式不正确"\`</span>
    Age      <span class="type">int</span>    <span class="string">\`v:"required|between:1,150"\`</span>
    Password <span class="type">string</span> <span class="string">\`v:"required|length:6,30"\`</span>
}

<span class="comment">// 参数自动绑定 + 校验</span>
<span class="keyword">func</span> (c *<span class="type">Controller</span>) <span class="function">CreateUser</span>(ctx context.Context, req *<span class="type">CreateUserReq</span>) (res *<span class="type">CreateUserRes</span>, err <span class="type">error</span>) {
    <span class="comment">// req 已自动绑定并校验！</span>
    <span class="keyword">return</span> service.User().<span class="function">Create</span>(ctx, req)
}</code></pre></div>
<div class="comparison-table"><h4>对比：参数校验</h4><table><thead><tr><th>功能</th><th>Zod (前端)</th><th>GoFrame</th></tr></thead><tbody><tr><td>必填</td><td><code>z.string()</code></td><td><code>v:"required"</code></td></tr><tr><td>长度</td><td><code>.min(2).max(30)</code></td><td><code>v:"length:2,30"</code></td></tr><tr><td>邮箱</td><td><code>.email()</code></td><td><code>v:"email"</code></td></tr><tr><td>范围</td><td><code>.min(1).max(150)</code></td><td><code>v:"between:1,150"</code></td></tr></tbody></table></div>` },
  { id: 'gf-orm', num: '04', title: '数据库操作 (ORM)', content: `<p class="gf-intro-text">GoFrame 内置强大的 ORM，支持链式操作。类似 Prisma / Sequelize。</p>
<div class="code-comparison"><div class="code-panel"><div class="code-panel-header"><div class="code-lang"><span class="lang-dot js-dot"></span>Prisma (Node.js)</div></div><pre class="code-block"><code><span class="keyword">const</span> users = <span class="keyword">await</span> prisma.user.<span class="function">findMany</span>({
  <span class="function">where</span>: { age: { gte: <span class="number">18</span> } },
  <span class="function">orderBy</span>: { createdAt: <span class="string">'desc'</span> },
  take: <span class="number">10</span>,
});

<span class="keyword">await</span> prisma.user.<span class="function">create</span>({
  data: { name: <span class="string">"张三"</span> }
});

<span class="keyword">await</span> prisma.user.<span class="function">update</span>({
  <span class="function">where</span>: { id: <span class="number">1</span> },
  data: { name: <span class="string">"李四"</span> },
});</code></pre></div><div class="code-vs"><div class="vs-badge">VS</div></div><div class="code-panel"><div class="code-panel-header"><div class="code-lang"><span class="lang-dot go-dot"></span>GoFrame ORM</div></div><pre class="code-block"><code>users, err := g.<span class="function">Model</span>(<span class="string">"user"</span>).
    <span class="function">Where</span>(<span class="string">"age >= ?"</span>, <span class="number">18</span>).
    <span class="function">Order</span>(<span class="string">"created_at DESC"</span>).
    <span class="function">Limit</span>(<span class="number">10</span>).
    <span class="function">All</span>()

_, err := g.<span class="function">Model</span>(<span class="string">"user"</span>).
    <span class="function">Data</span>(g.Map{<span class="string">"name"</span>: <span class="string">"张三"</span>}).
    <span class="function">Insert</span>()

_, err := g.<span class="function">Model</span>(<span class="string">"user"</span>).
    <span class="function">Where</span>(<span class="string">"id"</span>, <span class="number">1</span>).
    <span class="function">Data</span>(g.Map{<span class="string">"name"</span>: <span class="string">"李四"</span>}).
    <span class="function">Update</span>()</code></pre></div></div>` },
  { id: 'gf-config', num: '05', title: '配置管理与日志', content: `<p class="gf-intro-text">GoFrame 配置管理支持 YAML/TOML/JSON 等，支持热加载。类似 dotenv 但更强大。</p>
<div class="code-panel full-width"><div class="code-panel-header"><div class="code-lang"><span class="lang-dot go-dot"></span>config.yaml + 使用</div></div><pre class="code-block"><code><span class="comment"># manifest/config/config.yaml</span>
server:
  address: <span class="string">":8000"</span>
  openapiPath: <span class="string">"/api.json"</span>
database:
  default:
    host: <span class="string">"127.0.0.1"</span>
    port: <span class="string">"3306"</span>
    user: <span class="string">"root"</span>
    name: <span class="string">"mydb"</span>
    type: <span class="string">"mysql"</span>
logger:
  level: <span class="string">"all"</span>
  stdout: <span class="keyword">true</span>
  path: <span class="string">"./logs"</span>

<span class="comment">// 读取配置</span>
port := g.Cfg().<span class="function">MustGet</span>(ctx, <span class="string">"server.address"</span>)

<span class="comment">// 日志</span>
g.Log().<span class="function">Info</span>(ctx, <span class="string">"服务启动"</span>, <span class="string">"端口"</span>, port)
g.Log().<span class="function">Errorf</span>(ctx, <span class="string">"请求失败: %v"</span>, err)</code></pre></div>` },
  { id: 'gf-cache', num: '06', title: '缓存与 Redis', content: `<p class="gf-intro-text">GoFrame 内置缓存组件，支持内存缓存和 Redis 适配器，类似前端的 localStorage + Redis 缓存策略。</p>
<div class="code-comparison"><div class="code-panel"><div class="code-panel-header"><div class="code-lang"><span class="lang-dot js-dot"></span>Node.js + ioredis</div></div><pre class="code-block"><code><span class="keyword">const</span> Redis = <span class="function">require</span>(<span class="string">'ioredis'</span>);
<span class="keyword">const</span> redis = <span class="keyword">new</span> <span class="function">Redis</span>();

<span class="comment">// 设置缓存（带过期时间）</span>
<span class="keyword">await</span> redis.<span class="function">set</span>(<span class="string">'user:1'</span>, JSON.<span class="function">stringify</span>(user), <span class="string">'EX'</span>, <span class="number">3600</span>);

<span class="comment">// 读取缓存</span>
<span class="keyword">const</span> cached = <span class="keyword">await</span> redis.<span class="function">get</span>(<span class="string">'user:1'</span>);
<span class="keyword">if</span> (cached) <span class="keyword">return</span> JSON.<span class="function">parse</span>(cached);

<span class="comment">// 缓存穿透保护</span>
<span class="keyword">async function</span> <span class="function">getUser</span>(id) {
  <span class="keyword">let</span> data = <span class="keyword">await</span> redis.<span class="function">get</span>(<span class="string">\`user:\${id}\`</span>);
  <span class="keyword">if</span> (!data) {
    data = <span class="keyword">await</span> db.<span class="function">findUser</span>(id);
    <span class="keyword">await</span> redis.<span class="function">set</span>(<span class="string">\`user:\${id}\`</span>, data);
  }
  <span class="keyword">return</span> data;
}</code></pre></div><div class="code-vs"><div class="vs-badge">VS</div></div><div class="code-panel"><div class="code-panel-header"><div class="code-lang"><span class="lang-dot go-dot"></span>GoFrame 缓存</div></div><pre class="code-block"><code><span class="comment">// 内存缓存</span>
g.<span class="function">Cache</span>().<span class="function">Set</span>(ctx, <span class="string">"user:1"</span>, user, time.Hour)

val, _ := g.<span class="function">Cache</span>().<span class="function">Get</span>(ctx, <span class="string">"user:1"</span>)

<span class="comment">// Redis 适配器（config.yaml 配置即可）</span>
<span class="comment">// redis:</span>
<span class="comment">//   default:</span>
<span class="comment">//     address: 127.0.0.1:6379</span>
<span class="comment">//     db:      0</span>
g.<span class="function">Redis</span>().<span class="function">Do</span>(ctx, <span class="string">"SET"</span>, <span class="string">"key"</span>, <span class="string">"val"</span>, <span class="string">"EX"</span>, <span class="number">3600</span>)

<span class="comment">// 缓存穿透保护（GetOrSetFunc）</span>
val, _ := g.<span class="function">Cache</span>().<span class="function">GetOrSetFunc</span>(ctx,
    <span class="string">"user:1"</span>,
    <span class="keyword">func</span>(ctx context.Context) (<span class="keyword">interface</span>{}, <span class="type">error</span>) {
        <span class="keyword">return</span> dao.User.<span class="function">FindById</span>(ctx, <span class="number">1</span>)
    },
    time.Hour,
)</code></pre></div></div>
<div class="comparison-table"><h4>缓存策略对比</h4><table><thead><tr><th>场景</th><th>Node.js</th><th>GoFrame</th></tr></thead><tbody><tr><td>内存缓存</td><td><code>node-cache</code></td><td><code>g.Cache()</code> 内置</td></tr><tr><td>Redis</td><td><code>ioredis</code></td><td><code>g.Redis()</code> 内置</td></tr><tr><td>自动回源</td><td>手动实现</td><td><code>GetOrSetFunc</code></td></tr><tr><td>过期策略</td><td><code>'EX', 3600</code></td><td><code>time.Hour</code></td></tr></tbody></table></div>` },
  { id: 'gf-websocket', num: '07', title: 'WebSocket 实时通信', content: `<p class="gf-intro-text">GoFrame 内置 WebSocket 支持，配合 Goroutine 实现高性能实时通信。类似 Socket.io 但性能更强。</p>
<div class="code-comparison"><div class="code-panel"><div class="code-panel-header"><div class="code-lang"><span class="lang-dot js-dot"></span>Socket.io (Node.js)</div></div><pre class="code-block"><code><span class="keyword">const</span> io = <span class="function">require</span>(<span class="string">'socket.io'</span>)(server);

io.<span class="function">on</span>(<span class="string">'connection'</span>, (socket) => {
  console.<span class="function">log</span>(<span class="string">'用户连接'</span>, socket.id);

  socket.<span class="function">on</span>(<span class="string">'message'</span>, (data) => {
    <span class="comment">// 广播给所有人</span>
    io.<span class="function">emit</span>(<span class="string">'message'</span>, {
      user: socket.id,
      text: data,
    });
  });

  socket.<span class="function">on</span>(<span class="string">'disconnect'</span>, () => {
    console.<span class="function">log</span>(<span class="string">'用户断开'</span>);
  });
});</code></pre></div><div class="code-vs"><div class="vs-badge">VS</div></div><div class="code-panel"><div class="code-panel-header"><div class="code-lang"><span class="lang-dot go-dot"></span>GoFrame WebSocket</div></div><pre class="code-block"><code><span class="comment">// 路由注册</span>
s.<span class="function">BindHandler</span>(<span class="string">"/ws"</span>, <span class="keyword">func</span>(r *ghttp.<span class="type">Request</span>) {
    ws, _ := r.<span class="function">WebSocket</span>()
    <span class="keyword">defer</span> ws.<span class="function">Close</span>()

    <span class="keyword">for</span> {
        _, msg, err := ws.<span class="function">ReadMessage</span>()
        <span class="keyword">if</span> err != <span class="keyword">nil</span> { <span class="keyword">break</span> }

        <span class="comment">// 广播给所有连接</span>
        broadcast(msg)
    }
})

<span class="comment">// Goroutine 管理连接池</span>
<span class="keyword">var</span> clients = sync.Map{}

<span class="keyword">func</span> <span class="function">broadcast</span>(msg []<span class="type">byte</span>) {
    clients.<span class="function">Range</span>(<span class="keyword">func</span>(k, v <span class="keyword">interface</span>{}) <span class="type">bool</span> {
        conn := v.(*websocket.<span class="type">Conn</span>)
        conn.<span class="function">WriteMessage</span>(<span class="number">1</span>, msg)
        <span class="keyword">return</span> <span class="keyword">true</span>
    })
}</code></pre></div></div>
<div class="key-points"><h4>💡 Go WebSocket 优势</h4><ul>
<li><strong>Goroutine per Connection</strong>：每个连接一个 goroutine，代码更直观（vs Node.js 回调/事件）</li>
<li><strong>真正的并行处理</strong>：多核 CPU 并行处理消息，不阻塞事件循环</li>
<li><strong>极低内存</strong>：每个连接约 4KB，轻松支撑百万级连接</li>
<li><strong>sync.Map</strong>：并发安全的连接池管理，无需加锁</li>
</ul></div>` },
  { id: 'gf-deploy', num: '08', title: '微服务与部署', content: `<p class="gf-intro-text">Go 编译为单一二进制文件，部署极其简单。GoFrame 支持微服务架构和容器化部署。</p>
<div class="code-panel full-width"><div class="code-panel-header"><div class="code-lang"><span class="lang-dot go-dot"></span>Dockerfile — 多阶段构建</div></div><pre class="code-block"><code><span class="comment"># 构建阶段</span>
<span class="keyword">FROM</span> golang:1.22-alpine <span class="keyword">AS</span> builder
<span class="keyword">WORKDIR</span> /app
<span class="keyword">COPY</span> go.mod go.sum ./
<span class="keyword">RUN</span> go mod download
<span class="keyword">COPY</span> . .
<span class="keyword">RUN</span> CGO_ENABLED=0 go build -o server .

<span class="comment"># 运行阶段（仅 ~10MB 镜像！）</span>
<span class="keyword">FROM</span> alpine:latest
<span class="keyword">COPY</span> --from=builder /app/server /server
<span class="keyword">COPY</span> --from=builder /app/manifest /manifest
<span class="keyword">EXPOSE</span> 8000
<span class="keyword">CMD</span> ["/server"]</code></pre></div>
<div class="code-comparison"><div class="code-panel"><div class="code-panel-header"><div class="code-lang"><span class="lang-dot js-dot"></span>Node.js 部署</div></div><pre class="code-block"><code><span class="comment"># node_modules 动辄几百 MB</span>
<span class="comment"># 需要 Node.js 运行时</span>
<span class="comment"># PM2 进程管理</span>
$ npm install --production
$ pm2 start app.js -i max

<span class="comment"># Docker 镜像 ~300MB+</span>
<span class="keyword">FROM</span> node:20-alpine
<span class="keyword">COPY</span> . .
<span class="keyword">RUN</span> npm ci --production
<span class="keyword">CMD</span> ["node", "app.js"]</code></pre></div><div class="code-vs"><div class="vs-badge">VS</div></div><div class="code-panel"><div class="code-panel-header"><div class="code-lang"><span class="lang-dot go-dot"></span>Go 部署</div></div><pre class="code-block"><code><span class="comment"># 单一二进制，无依赖！</span>
<span class="comment"># 不需要运行时环境</span>
<span class="comment"># 自带高性能 HTTP Server</span>
$ go build -o myapp
$ ./myapp

<span class="comment"># Docker 镜像 ~10MB！</span>
<span class="comment"># 编译后直接丢到 scratch/alpine</span>
<span class="comment"># 启动时间 < 100ms</span>

<span class="comment"># 交叉编译一行搞定</span>
$ GOOS=linux GOARCH=amd64 go build -o myapp</code></pre></div></div>
<div class="comparison-table"><h4>部署对比</h4><table><thead><tr><th>指标</th><th>Node.js</th><th>Go</th></tr></thead><tbody><tr><td>产物大小</td><td>100MB+ (含 node_modules)</td><td>~10MB (单一二进制)</td></tr><tr><td>运行时</td><td>需要 Node.js</td><td>无需任何依赖</td></tr><tr><td>Docker 镜像</td><td>200-400MB</td><td>10-20MB</td></tr><tr><td>启动时间</td><td>1-5 秒</td><td>< 100ms</td></tr><tr><td>内存占用</td><td>50-200MB</td><td>10-30MB</td></tr><tr><td>进程管理</td><td>PM2 / forever</td><td>不需要（内置高可用）</td></tr></tbody></table></div>` },
];

// 实战项目数据
const PROJECTS = [
  { icon: '📝', title: 'Todo List API', diff: '入门', diffClass: 'diff-easy', desc: '用 GoFrame 构建一个完整的 Todo List RESTful API，理解基本的 CRUD 操作。',
    techs: ['GoFrame', 'RESTful', 'SQLite', 'gf cli'],
    features: ['CRUD 接口开发', '请求参数校验', '统一响应格式', '错误处理最佳实践'],
    hasGuide: true,
    steps: [
      { title: 'Step 1：初始化项目', desc: '使用 gf cli 创建项目骨架', code: `<span class="comment"># 创建项目</span>
$ gf init todo-api
$ cd todo-api

<span class="comment"># 项目结构</span>
todo-api/
├── api/v1/todo.go       <span class="comment"># 接口定义</span>
├── internal/
│   ├── controller/      <span class="comment"># 控制器</span>
│   ├── model/           <span class="comment"># 数据模型</span>
│   └── service/         <span class="comment"># 业务逻辑</span>
├── manifest/config/     <span class="comment"># 配置文件</span>
└── main.go` },
      { title: 'Step 2：定义 API 接口', desc: '在 api/v1/ 下定义请求和响应结构体', code: `<span class="comment">// api/v1/todo.go</span>
<span class="keyword">package</span> v1

<span class="keyword">import</span> <span class="string">"github.com/gogf/gf/v2/frame/g"</span>

<span class="comment">// 创建 Todo</span>
<span class="keyword">type</span> <span class="type">TodoCreateReq</span> <span class="keyword">struct</span> {
    g.Meta    <span class="string">\`path:"/todos" method:"post" tags:"Todo"\`</span>
    Title     <span class="type">string</span> <span class="string">\`v:"required|length:1,200" json:"title"\`</span>
    Completed <span class="type">bool</span>   <span class="string">\`json:"completed"\`</span>
}
<span class="keyword">type</span> <span class="type">TodoCreateRes</span> <span class="keyword">struct</span> {
    Id <span class="type">int</span> <span class="string">\`json:"id"\`</span>
}

<span class="comment">// 列表查询</span>
<span class="keyword">type</span> <span class="type">TodoListReq</span> <span class="keyword">struct</span> {
    g.Meta <span class="string">\`path:"/todos" method:"get" tags:"Todo"\`</span>
    Page   <span class="type">int</span> <span class="string">\`d:"1"  v:"min:1"       json:"page"\`</span>
    Size   <span class="type">int</span> <span class="string">\`d:"10" v:"max:100"     json:"size"\`</span>
}
<span class="keyword">type</span> <span class="type">TodoListRes</span> <span class="keyword">struct</span> {
    List  []<span class="type">TodoItem</span> <span class="string">\`json:"list"\`</span>
    Total <span class="type">int</span>        <span class="string">\`json:"total"\`</span>
}

<span class="comment">// 更新 & 删除</span>
<span class="keyword">type</span> <span class="type">TodoUpdateReq</span> <span class="keyword">struct</span> {
    g.Meta    <span class="string">\`path:"/todos/:id" method:"put" tags:"Todo"\`</span>
    Id        <span class="type">int</span>    <span class="string">\`v:"required" in:"path"\`</span>
    Title     <span class="type">string</span> <span class="string">\`v:"length:1,200" json:"title"\`</span>
    Completed *<span class="type">bool</span>  <span class="string">\`json:"completed"\`</span>
}
<span class="keyword">type</span> <span class="type">TodoUpdateRes</span> <span class="keyword">struct</span>{}

<span class="keyword">type</span> <span class="type">TodoDeleteReq</span> <span class="keyword">struct</span> {
    g.Meta <span class="string">\`path:"/todos/:id" method:"delete" tags:"Todo"\`</span>
    Id     <span class="type">int</span> <span class="string">\`v:"required" in:"path"\`</span>
}
<span class="keyword">type</span> <span class="type">TodoDeleteRes</span> <span class="keyword">struct</span>{}` },
      { title: 'Step 3：实现 Controller', desc: '处理请求，调用 Service 层', code: `<span class="comment">// internal/controller/todo.go</span>
<span class="keyword">package</span> controller

<span class="keyword">import</span> (
    <span class="string">"context"</span>
    v1 <span class="string">"todo-api/api/v1"</span>
    <span class="string">"todo-api/internal/service"</span>
)

<span class="keyword">var</span> Todo = &cTodo{}
<span class="keyword">type</span> cTodo <span class="keyword">struct</span>{}

<span class="keyword">func</span> (c *cTodo) <span class="function">Create</span>(ctx context.Context, req *v1.TodoCreateReq) (res *v1.TodoCreateRes, err <span class="type">error</span>) {
    id, err := service.<span class="function">Todo</span>().<span class="function">Create</span>(ctx, req.Title, req.Completed)
    <span class="keyword">if</span> err != <span class="keyword">nil</span> {
        <span class="keyword">return</span> <span class="keyword">nil</span>, err
    }
    <span class="keyword">return</span> &v1.TodoCreateRes{Id: id}, <span class="keyword">nil</span>
}

<span class="keyword">func</span> (c *cTodo) <span class="function">List</span>(ctx context.Context, req *v1.TodoListReq) (res *v1.TodoListRes, err <span class="type">error</span>) {
    list, total, err := service.<span class="function">Todo</span>().<span class="function">List</span>(ctx, req.Page, req.Size)
    <span class="keyword">if</span> err != <span class="keyword">nil</span> {
        <span class="keyword">return</span> <span class="keyword">nil</span>, err
    }
    <span class="keyword">return</span> &v1.TodoListRes{List: list, Total: total}, <span class="keyword">nil</span>
}

<span class="keyword">func</span> (c *cTodo) <span class="function">Update</span>(ctx context.Context, req *v1.TodoUpdateReq) (res *v1.TodoUpdateRes, err <span class="type">error</span>) {
    <span class="keyword">return</span> <span class="keyword">nil</span>, service.<span class="function">Todo</span>().<span class="function">Update</span>(ctx, req.Id, req.Title, req.Completed)
}

<span class="keyword">func</span> (c *cTodo) <span class="function">Delete</span>(ctx context.Context, req *v1.TodoDeleteReq) (res *v1.TodoDeleteRes, err <span class="type">error</span>) {
    <span class="keyword">return</span> <span class="keyword">nil</span>, service.<span class="function">Todo</span>().<span class="function">Delete</span>(ctx, req.Id)
}` },
      { title: 'Step 4：路由注册 & 启动', desc: '将 Controller 绑定到路由，启动服务', code: `<span class="comment">// internal/cmd/cmd.go</span>
<span class="keyword">var</span> Main = gcmd.Command{
    Name: <span class="string">"main"</span>,
    Func: <span class="keyword">func</span>(ctx context.Context, parser *gcmd.Parser) (err <span class="type">error</span>) {
        s := g.<span class="function">Server</span>()
        s.<span class="function">Use</span>(ghttp.MiddlewareHandlerResponse)
        s.<span class="function">Group</span>(<span class="string">"/api/v1"</span>, <span class="keyword">func</span>(group *ghttp.RouterGroup) {
            group.<span class="function">Bind</span>(controller.Todo)
        })
        s.<span class="function">Run</span>()
        <span class="keyword">return</span> <span class="keyword">nil</span>
    },
}

<span class="comment">// 启动 & 测试</span>
$ gf run main.go

<span class="comment"># 创建 Todo</span>
$ curl -X POST http://localhost:8000/api/v1/todos \\
  -H <span class="string">"Content-Type: application/json"</span> \\
  -d <span class="string">'{"title":"学习 GoFrame","completed":false}'</span>

<span class="comment"># 查询列表</span>
$ curl http://localhost:8000/api/v1/todos?page=1&size=10

<span class="comment"># API 文档自动生成</span>
<span class="comment"># 访问 http://localhost:8000/swagger 查看</span>` },
    ],
  },
  { icon: '👤', title: '用户认证系统', diff: '中级', diffClass: 'diff-medium', desc: '实现完整的用户注册、登录、JWT 认证系统，类似前端常见的 Auth 模块。',
    techs: ['GoFrame', 'JWT', 'MySQL', 'bcrypt', '中间件'],
    features: ['用户注册与登录', 'JWT Token 签发与验证', '密码加密存储', '权限中间件'],
    hasGuide: true,
    steps: [
      { title: 'Step 1：用户模型与密码加密', desc: '定义用户表结构，使用 bcrypt 加密密码', code: `<span class="comment">// internal/model/user.go</span>
<span class="keyword">type</span> <span class="type">User</span> <span class="keyword">struct</span> {
    Id        <span class="type">int</span>       <span class="string">\`json:"id"\`</span>
    Username  <span class="type">string</span>    <span class="string">\`json:"username"\`</span>
    Password  <span class="type">string</span>    <span class="string">\`json:"-"\`</span>  <span class="comment">// JSON 不输出密码</span>
    Email     <span class="type">string</span>    <span class="string">\`json:"email"\`</span>
    Role      <span class="type">string</span>    <span class="string">\`json:"role"\`</span>
    CreatedAt *gtime.Time <span class="string">\`json:"created_at"\`</span>
}

<span class="comment">// 密码加密（类似前端 bcrypt.js）</span>
<span class="keyword">import</span> <span class="string">"golang.org/x/crypto/bcrypt"</span>

<span class="keyword">func</span> <span class="function">HashPassword</span>(pwd <span class="type">string</span>) (<span class="type">string</span>, <span class="type">error</span>) {
    bytes, err := bcrypt.<span class="function">GenerateFromPassword</span>(
        []<span class="function">byte</span>(pwd), bcrypt.DefaultCost,
    )
    <span class="keyword">return</span> <span class="function">string</span>(bytes), err
}

<span class="keyword">func</span> <span class="function">CheckPassword</span>(pwd, hash <span class="type">string</span>) <span class="type">bool</span> {
    err := bcrypt.<span class="function">CompareHashAndPassword</span>([]<span class="function">byte</span>(hash), []<span class="function">byte</span>(pwd))
    <span class="keyword">return</span> err == <span class="keyword">nil</span>
}` },
      { title: 'Step 2：JWT Token 签发与验证', desc: '使用 jwt-go 生成和解析 Token', code: `<span class="comment">// internal/service/jwt.go</span>
<span class="keyword">import</span> <span class="string">"github.com/golang-jwt/jwt/v5"</span>

<span class="keyword">var</span> jwtSecret = []<span class="function">byte</span>(g.Cfg().MustGet(ctx, <span class="string">"jwt.secret"</span>).String())

<span class="keyword">type</span> <span class="type">Claims</span> <span class="keyword">struct</span> {
    UserId   <span class="type">int</span>    <span class="string">\`json:"user_id"\`</span>
    Username <span class="type">string</span> <span class="string">\`json:"username"\`</span>
    Role     <span class="type">string</span> <span class="string">\`json:"role"\`</span>
    jwt.RegisteredClaims
}

<span class="comment">// 生成 Token（类似前端 jsonwebtoken.sign）</span>
<span class="keyword">func</span> <span class="function">GenerateToken</span>(userId <span class="type">int</span>, username, role <span class="type">string</span>) (<span class="type">string</span>, <span class="type">error</span>) {
    claims := Claims{
        UserId: userId, Username: username, Role: role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.<span class="function">NewNumericDate</span>(time.Now().Add(<span class="number">24</span> * time.Hour)),
            IssuedAt:  jwt.<span class="function">NewNumericDate</span>(time.Now()),
        },
    }
    token := jwt.<span class="function">NewWithClaims</span>(jwt.SigningMethodHS256, claims)
    <span class="keyword">return</span> token.<span class="function">SignedString</span>(jwtSecret)
}

<span class="comment">// 解析 Token（类似 jsonwebtoken.verify）</span>
<span class="keyword">func</span> <span class="function">ParseToken</span>(tokenStr <span class="type">string</span>) (*Claims, <span class="type">error</span>) {
    token, err := jwt.<span class="function">ParseWithClaims</span>(tokenStr, &Claims{},
        <span class="keyword">func</span>(t *jwt.Token) (<span class="keyword">interface</span>{}, <span class="type">error</span>) {
            <span class="keyword">return</span> jwtSecret, <span class="keyword">nil</span>
        })
    <span class="keyword">if</span> claims, ok := token.Claims.(*Claims); ok && token.Valid {
        <span class="keyword">return</span> claims, <span class="keyword">nil</span>
    }
    <span class="keyword">return</span> <span class="keyword">nil</span>, err
}` },
      { title: 'Step 3：认证中间件', desc: '拦截请求，校验 Token 并注入用户信息', code: `<span class="comment">// internal/middleware/auth.go</span>
<span class="comment">// 类似 Express 的 passport-jwt 中间件</span>
<span class="keyword">func</span> <span class="function">AuthMiddleware</span>(r *ghttp.Request) {
    tokenStr := r.<span class="function">GetHeader</span>(<span class="string">"Authorization"</span>)

    <span class="comment">// 提取 Bearer token</span>
    <span class="keyword">if</span> !strings.<span class="function">HasPrefix</span>(tokenStr, <span class="string">"Bearer "</span>) {
        r.Response.<span class="function">WriteJsonExit</span>(g.Map{
            <span class="string">"code"</span>:    <span class="number">401</span>,
            <span class="string">"message"</span>: <span class="string">"未提供认证令牌"</span>,
        })
    }
    tokenStr = tokenStr[<span class="number">7</span>:]

    <span class="comment">// 解析验证</span>
    claims, err := service.<span class="function">ParseToken</span>(tokenStr)
    <span class="keyword">if</span> err != <span class="keyword">nil</span> {
        r.Response.<span class="function">WriteJsonExit</span>(g.Map{
            <span class="string">"code"</span>:    <span class="number">401</span>,
            <span class="string">"message"</span>: <span class="string">"令牌无效或已过期"</span>,
        })
    }

    <span class="comment">// 注入用户信息到上下文（类似 req.user = decoded）</span>
    r.<span class="function">SetCtxVar</span>(<span class="string">"userId"</span>, claims.UserId)
    r.<span class="function">SetCtxVar</span>(<span class="string">"username"</span>, claims.Username)
    r.<span class="function">SetCtxVar</span>(<span class="string">"role"</span>, claims.Role)
    r.<span class="function">Middleware</span>.<span class="function">Next</span>()
}

<span class="comment">// 路由中应用中间件</span>
s.<span class="function">Group</span>(<span class="string">"/api/v1"</span>, <span class="keyword">func</span>(group *ghttp.RouterGroup) {
    group.<span class="function">POST</span>(<span class="string">"/register"</span>, controller.Auth.Register)
    group.<span class="function">POST</span>(<span class="string">"/login"</span>, controller.Auth.Login)

    <span class="comment">// 需要认证的路由</span>
    group.<span class="function">Group</span>(<span class="string">"/"</span>, <span class="keyword">func</span>(group *ghttp.RouterGroup) {
        group.<span class="function">Middleware</span>(middleware.AuthMiddleware)
        group.<span class="function">GET</span>(<span class="string">"/profile"</span>, controller.Auth.Profile)
        group.<span class="function">Bind</span>(controller.Todo) <span class="comment">// Todo 需要登录</span>
    })
})` },
      { title: 'Step 4：注册与登录接口', desc: '实现完整的注册/登录/获取用户信息', code: `<span class="comment">// internal/controller/auth.go</span>
<span class="keyword">func</span> (c *cAuth) <span class="function">Register</span>(ctx context.Context, req *v1.RegisterReq) (res *v1.RegisterRes, err <span class="type">error</span>) {
    <span class="comment">// 1. 检查用户名是否已存在</span>
    exists, _ := g.<span class="function">Model</span>(<span class="string">"users"</span>).<span class="function">Where</span>(<span class="string">"username"</span>, req.Username).<span class="function">One</span>()
    <span class="keyword">if</span> !exists.IsEmpty() {
        <span class="keyword">return</span> <span class="keyword">nil</span>, gerror.<span class="function">New</span>(<span class="string">"用户名已存在"</span>)
    }

    <span class="comment">// 2. 加密密码</span>
    hash, _ := <span class="function">HashPassword</span>(req.Password)

    <span class="comment">// 3. 插入数据库</span>
    result, err := g.<span class="function">Model</span>(<span class="string">"users"</span>).<span class="function">Data</span>(g.Map{
        <span class="string">"username"</span>: req.Username,
        <span class="string">"password"</span>: hash,
        <span class="string">"email"</span>:    req.Email,
        <span class="string">"role"</span>:     <span class="string">"user"</span>,
    }).<span class="function">Insert</span>()

    id, _ := result.<span class="function">LastInsertId</span>()
    <span class="keyword">return</span> &v1.RegisterRes{Id: <span class="function">int</span>(id)}, err
}

<span class="keyword">func</span> (c *cAuth) <span class="function">Login</span>(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err <span class="type">error</span>) {
    <span class="comment">// 1. 查找用户</span>
    user, _ := g.<span class="function">Model</span>(<span class="string">"users"</span>).<span class="function">Where</span>(<span class="string">"username"</span>, req.Username).<span class="function">One</span>()
    <span class="keyword">if</span> user.IsEmpty() {
        <span class="keyword">return</span> <span class="keyword">nil</span>, gerror.<span class="function">New</span>(<span class="string">"用户不存在"</span>)
    }

    <span class="comment">// 2. 验证密码</span>
    <span class="keyword">if</span> !<span class="function">CheckPassword</span>(req.Password, user[<span class="string">"password"</span>].String()) {
        <span class="keyword">return</span> <span class="keyword">nil</span>, gerror.<span class="function">New</span>(<span class="string">"密码错误"</span>)
    }

    <span class="comment">// 3. 生成 JWT</span>
    token, _ := <span class="function">GenerateToken</span>(
        user[<span class="string">"id"</span>].Int(),
        user[<span class="string">"username"</span>].String(),
        user[<span class="string">"role"</span>].String(),
    )
    <span class="keyword">return</span> &v1.LoginRes{Token: token}, <span class="keyword">nil</span>
}

<span class="comment">// 测试命令</span>
<span class="comment">// 注册</span>
$ curl -X POST :8000/api/v1/register \\
  -d <span class="string">'{"username":"dev","password":"123456","email":"dev@go.dev"}'</span>

<span class="comment">// 登录获取 Token</span>
$ curl -X POST :8000/api/v1/login \\
  -d <span class="string">'{"username":"dev","password":"123456"}'</span>

<span class="comment">// 使用 Token 访问受保护接口</span>
$ curl :8000/api/v1/profile \\
  -H <span class="string">"Authorization: Bearer <your_token>"</span>` },
    ],
  },
  { icon: '📰', title: '博客系统', diff: '中高级', diffClass: 'diff-hard', desc: '全功能博客系统，包含文章管理、分类标签、评论系统，前后端分离架构。',
    techs: ['GoFrame', 'MySQL', 'Redis', 'Swagger', 'CORS'],
    features: ['文章 CRUD + 分页', '分类与标签系统', '评论功能', 'API 文档自动生成'],
    hasGuide: true,
    steps: [
      { title: 'Step 1：数据库设计与模型定义', desc: '设计文章、分类、标签、评论的数据表结构', code: `<span class="comment">-- SQL 建表语句</span>
<span class="keyword">CREATE TABLE</span> articles (
    id          <span class="type">INT</span> PRIMARY KEY AUTO_INCREMENT,
    title       <span class="type">VARCHAR</span>(200) NOT NULL,
    slug        <span class="type">VARCHAR</span>(200) UNIQUE NOT NULL,
    content     <span class="type">TEXT</span> NOT NULL,
    summary     <span class="type">VARCHAR</span>(500),
    cover_image <span class="type">VARCHAR</span>(500),
    category_id <span class="type">INT</span>,
    author_id   <span class="type">INT</span> NOT NULL,
    status      <span class="type">TINYINT</span> DEFAULT 0 <span class="comment">-- 0:草稿 1:已发布</span>,
    view_count  <span class="type">INT</span> DEFAULT 0,
    created_at  <span class="type">DATETIME</span>,
    updated_at  <span class="type">DATETIME</span>
);

<span class="keyword">CREATE TABLE</span> categories (
    id   <span class="type">INT</span> PRIMARY KEY AUTO_INCREMENT,
    name <span class="type">VARCHAR</span>(50) UNIQUE NOT NULL,
    slug <span class="type">VARCHAR</span>(50) UNIQUE NOT NULL,
    sort <span class="type">INT</span> DEFAULT 0
);

<span class="keyword">CREATE TABLE</span> tags (
    id   <span class="type">INT</span> PRIMARY KEY AUTO_INCREMENT,
    name <span class="type">VARCHAR</span>(50) UNIQUE NOT NULL
);

<span class="comment">-- 文章与标签多对多关系</span>
<span class="keyword">CREATE TABLE</span> article_tags (
    article_id <span class="type">INT</span>,
    tag_id     <span class="type">INT</span>,
    PRIMARY KEY (article_id, tag_id)
);

<span class="keyword">CREATE TABLE</span> comments (
    id         <span class="type">INT</span> PRIMARY KEY AUTO_INCREMENT,
    article_id <span class="type">INT</span> NOT NULL,
    user_id    <span class="type">INT</span> NOT NULL,
    content    <span class="type">TEXT</span> NOT NULL,
    parent_id  <span class="type">INT</span> DEFAULT 0 <span class="comment">-- 支持楼中楼回复</span>,
    created_at <span class="type">DATETIME</span>
);

<span class="comment">// internal/model/article.go — Go 模型定义</span>
<span class="keyword">type</span> <span class="type">Article</span> <span class="keyword">struct</span> {
    Id         <span class="type">int</span>         <span class="string">\`json:"id"\`</span>
    Title      <span class="type">string</span>      <span class="string">\`json:"title"\`</span>
    Slug       <span class="type">string</span>      <span class="string">\`json:"slug"\`</span>
    Content    <span class="type">string</span>      <span class="string">\`json:"content"\`</span>
    Summary    <span class="type">string</span>      <span class="string">\`json:"summary"\`</span>
    CoverImage <span class="type">string</span>      <span class="string">\`json:"cover_image"\`</span>
    CategoryId <span class="type">int</span>         <span class="string">\`json:"category_id"\`</span>
    AuthorId   <span class="type">int</span>         <span class="string">\`json:"author_id"\`</span>
    Status     <span class="type">int</span>         <span class="string">\`json:"status"\`</span>
    ViewCount  <span class="type">int</span>         <span class="string">\`json:"view_count"\`</span>
    CreatedAt  *gtime.Time <span class="string">\`json:"created_at"\`</span>
    UpdatedAt  *gtime.Time <span class="string">\`json:"updated_at"\`</span>
    <span class="comment">// 关联字段</span>
    Category   *Category   <span class="string">\`json:"category,omitempty"\`</span>
    Tags       []Tag       <span class="string">\`json:"tags,omitempty"\`</span>
    Author     *User       <span class="string">\`json:"author,omitempty"\`</span>
}` },
      { title: 'Step 2：文章 CRUD + 分页搜索', desc: '实现文章的创建、查询、更新、删除，支持分页与关键词搜索', code: `<span class="comment">// api/v1/article.go — 接口定义</span>
<span class="keyword">type</span> <span class="type">ArticleCreateReq</span> <span class="keyword">struct</span> {
    g.Meta     <span class="string">\`path:"/articles" method:"post" tags:"文章管理"\`</span>
    Title      <span class="type">string</span>   <span class="string">\`v:"required|length:1,200" json:"title"\`</span>
    Content    <span class="type">string</span>   <span class="string">\`v:"required" json:"content"\`</span>
    Summary    <span class="type">string</span>   <span class="string">\`json:"summary"\`</span>
    CoverImage <span class="type">string</span>   <span class="string">\`json:"cover_image"\`</span>
    CategoryId <span class="type">int</span>      <span class="string">\`json:"category_id"\`</span>
    TagIds     []<span class="type">int</span>    <span class="string">\`json:"tag_ids"\`</span>
    Status     <span class="type">int</span>      <span class="string">\`v:"in:0,1" d:"0" json:"status"\`</span>
}
<span class="keyword">type</span> <span class="type">ArticleCreateRes</span> <span class="keyword">struct</span> {
    Id <span class="type">int</span> <span class="string">\`json:"id"\`</span>
}

<span class="keyword">type</span> <span class="type">ArticleListReq</span> <span class="keyword">struct</span> {
    g.Meta     <span class="string">\`path:"/articles" method:"get" tags:"文章管理"\`</span>
    Page       <span class="type">int</span>    <span class="string">\`d:"1" v:"min:1" json:"page"\`</span>
    Size       <span class="type">int</span>    <span class="string">\`d:"10" v:"max:50" json:"size"\`</span>
    Keyword    <span class="type">string</span> <span class="string">\`json:"keyword"\`</span>         <span class="comment">// 搜索关键词</span>
    CategoryId <span class="type">int</span>    <span class="string">\`json:"category_id"\`</span>     <span class="comment">// 按分类筛选</span>
    TagId      <span class="type">int</span>    <span class="string">\`json:"tag_id"\`</span>          <span class="comment">// 按标签筛选</span>
    Status     <span class="type">int</span>    <span class="string">\`d:"-1" json:"status"\`</span>   <span class="comment">// -1:全部 0:草稿 1:已发布</span>
}

<span class="comment">// internal/logic/article.go — 核心业务逻辑</span>
<span class="keyword">func</span> (s *sArticle) <span class="function">List</span>(ctx context.Context, req *v1.ArticleListReq) (list []model.Article, total <span class="type">int</span>, err <span class="type">error</span>) {
    m := g.<span class="function">Model</span>(<span class="string">"articles a"</span>).
        <span class="function">LeftJoin</span>(<span class="string">"categories c"</span>, <span class="string">"a.category_id = c.id"</span>).
        <span class="function">Fields</span>(<span class="string">"a.*, c.name as category_name"</span>)

    <span class="comment">// 条件过滤（类似 Prisma 的 where 条件）</span>
    <span class="keyword">if</span> req.Keyword != <span class="string">""</span> {
        m = m.<span class="function">WhereLike</span>(<span class="string">"a.title"</span>, <span class="string">"%"</span>+req.Keyword+<span class="string">"%"</span>).
            <span class="function">WhereOrLike</span>(<span class="string">"a.content"</span>, <span class="string">"%"</span>+req.Keyword+<span class="string">"%"</span>)
    }
    <span class="keyword">if</span> req.CategoryId > <span class="number">0</span> {
        m = m.<span class="function">Where</span>(<span class="string">"a.category_id"</span>, req.CategoryId)
    }
    <span class="keyword">if</span> req.Status >= <span class="number">0</span> {
        m = m.<span class="function">Where</span>(<span class="string">"a.status"</span>, req.Status)
    }

    <span class="comment">// 分页查询（类似 Prisma 的 skip + take）</span>
    total, err = m.<span class="function">Count</span>()
    <span class="keyword">if</span> err != <span class="keyword">nil</span> { <span class="keyword">return</span> }

    err = m.<span class="function">Page</span>(req.Page, req.Size).
        <span class="function">Order</span>(<span class="string">"a.created_at DESC"</span>).
        <span class="function">Scan</span>(&list)
    <span class="keyword">return</span>
}

<span class="comment">// 创建文章（含标签关联 — 事务处理）</span>
<span class="keyword">func</span> (s *sArticle) <span class="function">Create</span>(ctx context.Context, req *v1.ArticleCreateReq, authorId <span class="type">int</span>) (id <span class="type">int</span>, err <span class="type">error</span>) {
    err = g.<span class="function">DB</span>().<span class="function">Transaction</span>(ctx, <span class="keyword">func</span>(ctx context.Context, tx gdb.TX) <span class="type">error</span> {
        <span class="comment">// 1. 自动生成 slug</span>
        slug := gstr.<span class="function">Replace</span>(strings.<span class="function">ToLower</span>(req.Title), <span class="string">" "</span>, <span class="string">"-"</span>)

        <span class="comment">// 2. 插入文章</span>
        result, err := tx.<span class="function">Model</span>(<span class="string">"articles"</span>).<span class="function">Data</span>(g.Map{
            <span class="string">"title"</span>:       req.Title,
            <span class="string">"slug"</span>:        slug,
            <span class="string">"content"</span>:     req.Content,
            <span class="string">"summary"</span>:     req.Summary,
            <span class="string">"cover_image"</span>: req.CoverImage,
            <span class="string">"category_id"</span>: req.CategoryId,
            <span class="string">"author_id"</span>:   authorId,
            <span class="string">"status"</span>:      req.Status,
        }).<span class="function">Insert</span>()
        <span class="keyword">if</span> err != <span class="keyword">nil</span> { <span class="keyword">return</span> err }

        lastId, _ := result.<span class="function">LastInsertId</span>()
        id = <span class="function">int</span>(lastId)

        <span class="comment">// 3. 批量插入标签关联</span>
        <span class="keyword">if</span> <span class="function">len</span>(req.TagIds) > <span class="number">0</span> {
            batch := make([]g.Map, <span class="number">0</span>, <span class="function">len</span>(req.TagIds))
            <span class="keyword">for</span> _, tagId := <span class="keyword">range</span> req.TagIds {
                batch = <span class="function">append</span>(batch, g.Map{
                    <span class="string">"article_id"</span>: id, <span class="string">"tag_id"</span>: tagId,
                })
            }
            _, err = tx.<span class="function">Model</span>(<span class="string">"article_tags"</span>).<span class="function">Data</span>(batch).<span class="function">Insert</span>()
        }
        <span class="keyword">return</span> err
    })
    <span class="keyword">return</span>
}` },
      { title: 'Step 3：分类与标签系统', desc: '实现分类 CRUD 和标签管理，支持按标签筛选文章', code: `<span class="comment">// api/v1/category.go — 分类接口</span>
<span class="keyword">type</span> <span class="type">CategoryCreateReq</span> <span class="keyword">struct</span> {
    g.Meta <span class="string">\`path:"/categories" method:"post" tags:"分类管理"\`</span>
    Name   <span class="type">string</span> <span class="string">\`v:"required|length:1,50" json:"name"\`</span>
    Sort   <span class="type">int</span>    <span class="string">\`d:"0" json:"sort"\`</span>
}

<span class="keyword">type</span> <span class="type">CategoryListReq</span> <span class="keyword">struct</span> {
    g.Meta <span class="string">\`path:"/categories" method:"get" tags:"分类管理"\`</span>
}
<span class="keyword">type</span> <span class="type">CategoryListRes</span> <span class="keyword">struct</span> {
    List []<span class="type">CategoryWithCount</span> <span class="string">\`json:"list"\`</span>
}

<span class="comment">// 分类带文章计数</span>
<span class="keyword">type</span> <span class="type">CategoryWithCount</span> <span class="keyword">struct</span> {
    Id           <span class="type">int</span>    <span class="string">\`json:"id"\`</span>
    Name         <span class="type">string</span> <span class="string">\`json:"name"\`</span>
    Slug         <span class="type">string</span> <span class="string">\`json:"slug"\`</span>
    ArticleCount <span class="type">int</span>    <span class="string">\`json:"article_count"\`</span>
}

<span class="comment">// internal/logic/category.go</span>
<span class="keyword">func</span> (s *sCategory) <span class="function">ListWithCount</span>(ctx context.Context) ([]model.CategoryWithCount, <span class="type">error</span>) {
    <span class="keyword">var</span> list []model.CategoryWithCount
    err := g.<span class="function">Model</span>(<span class="string">"categories c"</span>).
        <span class="function">LeftJoin</span>(<span class="string">"articles a"</span>, <span class="string">"c.id = a.category_id AND a.status = 1"</span>).
        <span class="function">Fields</span>(<span class="string">"c.id, c.name, c.slug, COUNT(a.id) as article_count"</span>).
        <span class="function">Group</span>(<span class="string">"c.id"</span>).
        <span class="function">Order</span>(<span class="string">"c.sort ASC, c.id ASC"</span>).
        <span class="function">Scan</span>(&list)
    <span class="keyword">return</span> list, err
}

<span class="comment">// api/v1/tag.go — 标签接口</span>
<span class="keyword">type</span> <span class="type">TagListReq</span> <span class="keyword">struct</span> {
    g.Meta <span class="string">\`path:"/tags" method:"get" tags:"标签管理"\`</span>
}

<span class="comment">// 热门标签（按关联文章数排序）</span>
<span class="keyword">func</span> (s *sTag) <span class="function">HotTags</span>(ctx context.Context, limit <span class="type">int</span>) ([]model.TagWithCount, <span class="type">error</span>) {
    <span class="keyword">var</span> list []model.TagWithCount
    err := g.<span class="function">Model</span>(<span class="string">"tags t"</span>).
        <span class="function">LeftJoin</span>(<span class="string">"article_tags at"</span>, <span class="string">"t.id = at.tag_id"</span>).
        <span class="function">Fields</span>(<span class="string">"t.id, t.name, COUNT(at.article_id) as count"</span>).
        <span class="function">Group</span>(<span class="string">"t.id"</span>).
        <span class="function">Order</span>(<span class="string">"count DESC"</span>).
        <span class="function">Limit</span>(limit).
        <span class="function">Scan</span>(&list)
    <span class="keyword">return</span> list, err
}

<span class="comment">// 按标签查文章</span>
<span class="keyword">func</span> (s *sArticle) <span class="function">ListByTag</span>(ctx context.Context, tagId, page, size <span class="type">int</span>) ([]model.Article, <span class="type">int</span>, <span class="type">error</span>) {
    m := g.<span class="function">Model</span>(<span class="string">"articles a"</span>).
        <span class="function">InnerJoin</span>(<span class="string">"article_tags at"</span>, <span class="string">"a.id = at.article_id"</span>).
        <span class="function">Where</span>(<span class="string">"at.tag_id"</span>, tagId).
        <span class="function">Where</span>(<span class="string">"a.status"</span>, <span class="number">1</span>)
    total, _ := m.<span class="function">Count</span>()
    <span class="keyword">var</span> list []model.Article
    err := m.<span class="function">Page</span>(page, size).<span class="function">Order</span>(<span class="string">"a.created_at DESC"</span>).<span class="function">Scan</span>(&list)
    <span class="keyword">return</span> list, total, err
}` },
      { title: 'Step 4：评论系统 + Redis 缓存', desc: '实现树形评论、文章阅读量统计及 Redis 缓存策略', code: `<span class="comment">// api/v1/comment.go — 评论接口</span>
<span class="keyword">type</span> <span class="type">CommentCreateReq</span> <span class="keyword">struct</span> {
    g.Meta    <span class="string">\`path:"/articles/:articleId/comments" method:"post" tags:"评论"\`</span>
    ArticleId <span class="type">int</span>    <span class="string">\`v:"required" in:"path"\`</span>
    Content   <span class="type">string</span> <span class="string">\`v:"required|length:1,1000" json:"content"\`</span>
    ParentId  <span class="type">int</span>    <span class="string">\`d:"0" json:"parent_id"\`</span>  <span class="comment">// 回复某条评论</span>
}

<span class="keyword">type</span> <span class="type">CommentListReq</span> <span class="keyword">struct</span> {
    g.Meta    <span class="string">\`path:"/articles/:articleId/comments" method:"get" tags:"评论"\`</span>
    ArticleId <span class="type">int</span> <span class="string">\`v:"required" in:"path"\`</span>
    Page      <span class="type">int</span> <span class="string">\`d:"1" json:"page"\`</span>
    Size      <span class="type">int</span> <span class="string">\`d:"20" json:"size"\`</span>
}

<span class="comment">// internal/logic/comment.go — 树形评论</span>
<span class="keyword">type</span> <span class="type">CommentTree</span> <span class="keyword">struct</span> {
    model.Comment
    Children []CommentTree <span class="string">\`json:"children,omitempty"\`</span>
    Author   *model.User   <span class="string">\`json:"author"\`</span>
}

<span class="keyword">func</span> (s *sComment) <span class="function">TreeList</span>(ctx context.Context, articleId <span class="type">int</span>) ([]CommentTree, <span class="type">error</span>) {
    <span class="keyword">var</span> comments []model.Comment
    err := g.<span class="function">Model</span>(<span class="string">"comments"</span>).
        <span class="function">Where</span>(<span class="string">"article_id"</span>, articleId).
        <span class="function">Order</span>(<span class="string">"created_at ASC"</span>).
        <span class="function">Scan</span>(&comments)
    <span class="keyword">if</span> err != <span class="keyword">nil</span> { <span class="keyword">return</span> <span class="keyword">nil</span>, err }
    <span class="keyword">return</span> <span class="function">buildTree</span>(comments, <span class="number">0</span>), <span class="keyword">nil</span>
}

<span class="keyword">func</span> <span class="function">buildTree</span>(all []model.Comment, parentId <span class="type">int</span>) []CommentTree {
    <span class="keyword">var</span> tree []CommentTree
    <span class="keyword">for</span> _, c := <span class="keyword">range</span> all {
        <span class="keyword">if</span> c.ParentId == parentId {
            node := CommentTree{Comment: c}
            node.Children = <span class="function">buildTree</span>(all, c.Id)
            tree = <span class="function">append</span>(tree, node)
        }
    }
    <span class="keyword">return</span> tree
}

<span class="comment">// ---- Redis 缓存策略 ----</span>

<span class="comment">// 文章详情缓存（类似前端 SWR / React Query）</span>
<span class="keyword">func</span> (s *sArticle) <span class="function">Detail</span>(ctx context.Context, id <span class="type">int</span>) (*model.Article, <span class="type">error</span>) {
    cacheKey := fmt.<span class="function">Sprintf</span>(<span class="string">"article:%d"</span>, id)

    <span class="comment">// 缓存优先（GetOrSetFunc = 自动回源）</span>
    val, err := g.<span class="function">Cache</span>().<span class="function">GetOrSetFunc</span>(ctx, cacheKey,
        <span class="keyword">func</span>(ctx context.Context) (<span class="keyword">interface</span>{}, <span class="type">error</span>) {
            <span class="keyword">var</span> article model.Article
            err := g.<span class="function">Model</span>(<span class="string">"articles"</span>).<span class="function">Where</span>(<span class="string">"id"</span>, id).<span class="function">Scan</span>(&article)
            <span class="keyword">return</span> &article, err
        }, <span class="number">30</span> * time.Minute)
    <span class="keyword">if</span> err != <span class="keyword">nil</span> { <span class="keyword">return</span> <span class="keyword">nil</span>, err }
    <span class="keyword">return</span> val.Val().(*model.Article), <span class="keyword">nil</span>
}

<span class="comment">// 浏览量统计（Redis INCR，定期同步到 MySQL）</span>
<span class="keyword">func</span> (s *sArticle) <span class="function">IncrViewCount</span>(ctx context.Context, id <span class="type">int</span>) {
    key := fmt.<span class="function">Sprintf</span>(<span class="string">"article:views:%d"</span>, id)
    g.<span class="function">Redis</span>().<span class="function">Do</span>(ctx, <span class="string">"INCR"</span>, key)
}

<span class="comment">// 定时任务：将 Redis 浏览量刷入 MySQL</span>
<span class="keyword">func</span> <span class="function">SyncViewCounts</span>(ctx context.Context) {
    keys, _ := g.<span class="function">Redis</span>().<span class="function">Do</span>(ctx, <span class="string">"KEYS"</span>, <span class="string">"article:views:*"</span>)
    <span class="keyword">for</span> _, key := <span class="keyword">range</span> keys.Strings() {
        id := gstr.<span class="function">TrimLeft</span>(key, <span class="string">"article:views:"</span>)
        count, _ := g.<span class="function">Redis</span>().<span class="function">Do</span>(ctx, <span class="string">"GETDEL"</span>, key)
        g.<span class="function">Model</span>(<span class="string">"articles"</span>).<span class="function">Where</span>(<span class="string">"id"</span>, id).
            <span class="function">Increment</span>(<span class="string">"view_count"</span>, count.Int())
    }
}` },
      { title: 'Step 5：路由注册与 Swagger 文档', desc: '完整路由结构、CORS 配置、自动生成 API 文档', code: `<span class="comment">// internal/cmd/cmd.go — 完整路由结构</span>
<span class="keyword">var</span> Main = gcmd.Command{
    Name: <span class="string">"main"</span>,
    Func: <span class="keyword">func</span>(ctx context.Context, parser *gcmd.Parser) (err <span class="type">error</span>) {
        s := g.<span class="function">Server</span>()

        <span class="comment">// 1. CORS 中间件（前后端分离必备）</span>
        s.<span class="function">Use</span>(ghttp.MiddlewareCORS)
        s.<span class="function">Use</span>(ghttp.MiddlewareHandlerResponse)

        <span class="comment">// 2. API 路由</span>
        s.<span class="function">Group</span>(<span class="string">"/api/v1"</span>, <span class="keyword">func</span>(group *ghttp.RouterGroup) {
            <span class="comment">// 公开接口</span>
            group.<span class="function">GET</span>(<span class="string">"/articles"</span>, controller.Article.List)
            group.<span class="function">GET</span>(<span class="string">"/articles/:id"</span>, controller.Article.Detail)
            group.<span class="function">GET</span>(<span class="string">"/categories"</span>, controller.Category.List)
            group.<span class="function">GET</span>(<span class="string">"/tags"</span>, controller.Tag.List)
            group.<span class="function">GET</span>(<span class="string">"/articles/:articleId/comments"</span>, controller.Comment.List)

            <span class="comment">// 需要认证的接口</span>
            group.<span class="function">Group</span>(<span class="string">"/"</span>, <span class="keyword">func</span>(group *ghttp.RouterGroup) {
                group.<span class="function">Middleware</span>(middleware.Auth)

                <span class="comment">// 文章管理</span>
                group.<span class="function">POST</span>(<span class="string">"/articles"</span>, controller.Article.Create)
                group.<span class="function">PUT</span>(<span class="string">"/articles/:id"</span>, controller.Article.Update)
                group.<span class="function">DELETE</span>(<span class="string">"/articles/:id"</span>, controller.Article.Delete)

                <span class="comment">// 评论</span>
                group.<span class="function">POST</span>(<span class="string">"/articles/:articleId/comments"</span>, controller.Comment.Create)
                group.<span class="function">DELETE</span>(<span class="string">"/comments/:id"</span>, controller.Comment.Delete)

                <span class="comment">// 管理员接口</span>
                group.<span class="function">Group</span>(<span class="string">"/"</span>, <span class="keyword">func</span>(group *ghttp.RouterGroup) {
                    group.<span class="function">Middleware</span>(middleware.AdminOnly)
                    group.<span class="function">Bind</span>(controller.Category) <span class="comment">// 分类 CRUD</span>
                    group.<span class="function">Bind</span>(controller.Tag)      <span class="comment">// 标签 CRUD</span>
                })
            })
        })

        <span class="comment">// 3. Swagger 文档（访问 /swagger 即可查看）</span>
        s.<span class="function">SetOpenApiPath</span>(<span class="string">"/api.json"</span>)
        s.<span class="function">SetSwaggerPath</span>(<span class="string">"/swagger"</span>)

        s.<span class="function">Run</span>()
        <span class="keyword">return</span>
    },
}

<span class="comment">// manifest/config/config.yaml</span>
server:
  address: <span class="string">":8000"</span>
  openapiPath: <span class="string">"/api.json"</span>
  swaggerPath: <span class="string">"/swagger"</span>
database:
  default:
    type:  <span class="string">"mysql"</span>
    host:  <span class="string">"127.0.0.1"</span>
    port:  <span class="string">"3306"</span>
    user:  <span class="string">"root"</span>
    pass:  <span class="string">"your_password"</span>
    name:  <span class="string">"blog_db"</span>
redis:
  default:
    address: <span class="string">"127.0.0.1:6379"</span>
    db:      <span class="number">0</span>

<span class="comment"># 测试接口</span>
$ curl http://localhost:8000/api/v1/articles?page=1&size=10&keyword=Go
$ curl http://localhost:8000/api/v1/categories
$ curl http://localhost:8000/swagger  <span class="comment"># 查看自动生成的 API 文档</span>` },
    ],
  },
  { icon: '🛒', title: '电商 API 平台', diff: '高级', diffClass: 'diff-expert', desc: '企业级电商后端 API，涵盖商品、订单、支付、库存等核心模块。',
    techs: ['GoFrame', 'MySQL', 'Redis', '微服务', 'Docker'],
    features: ['商品管理系统', '订单状态机', '库存并发控制', 'Redis 缓存策略'],
    hasGuide: true,
    steps: [
      { title: 'Step 1：核心数据模型设计', desc: '设计商品、SKU、订单、库存等核心业务表', code: `<span class="comment">-- 商品主表</span>
<span class="keyword">CREATE TABLE</span> products (
    id          <span class="type">BIGINT</span> PRIMARY KEY AUTO_INCREMENT,
    name        <span class="type">VARCHAR</span>(200) NOT NULL,
    description <span class="type">TEXT</span>,
    category_id <span class="type">INT</span>,
    brand       <span class="type">VARCHAR</span>(100),
    status      <span class="type">TINYINT</span> DEFAULT 1 <span class="comment">-- 1:上架 0:下架</span>,
    created_at  <span class="type">DATETIME</span>,
    updated_at  <span class="type">DATETIME</span>
);

<span class="comment">-- SKU 表（商品规格，如颜色/尺码）</span>
<span class="keyword">CREATE TABLE</span> skus (
    id         <span class="type">BIGINT</span> PRIMARY KEY AUTO_INCREMENT,
    product_id <span class="type">BIGINT</span> NOT NULL,
    attrs      <span class="type">JSON</span> <span class="comment">-- {"颜色":"黑色","尺码":"XL"}</span>,
    price      <span class="type">DECIMAL</span>(10,2) NOT NULL,
    stock      <span class="type">INT</span> NOT NULL DEFAULT 0,
    image      <span class="type">VARCHAR</span>(500)
);

<span class="comment">-- 订单主表</span>
<span class="keyword">CREATE TABLE</span> orders (
    id          <span class="type">BIGINT</span> PRIMARY KEY AUTO_INCREMENT,
    order_no    <span class="type">VARCHAR</span>(32) UNIQUE NOT NULL,
    user_id     <span class="type">BIGINT</span> NOT NULL,
    total_amount <span class="type">DECIMAL</span>(10,2) NOT NULL,
    status      <span class="type">TINYINT</span> DEFAULT 0
        <span class="comment">-- 0:待支付 1:已支付 2:已发货 3:已完成 4:已取消</span>,
    address     <span class="type">JSON</span>,
    paid_at     <span class="type">DATETIME</span>,
    created_at  <span class="type">DATETIME</span>
);

<span class="comment">-- 订单明细</span>
<span class="keyword">CREATE TABLE</span> order_items (
    id         <span class="type">BIGINT</span> PRIMARY KEY AUTO_INCREMENT,
    order_id   <span class="type">BIGINT</span> NOT NULL,
    sku_id     <span class="type">BIGINT</span> NOT NULL,
    product_name <span class="type">VARCHAR</span>(200),
    sku_attrs  <span class="type">JSON</span>,
    price      <span class="type">DECIMAL</span>(10,2) NOT NULL,
    quantity   <span class="type">INT</span> NOT NULL
);

<span class="comment">// internal/model/product.go</span>
<span class="keyword">type</span> <span class="type">Product</span> <span class="keyword">struct</span> {
    Id          <span class="type">int64</span>       <span class="string">\`json:"id"\`</span>
    Name        <span class="type">string</span>      <span class="string">\`json:"name"\`</span>
    Description <span class="type">string</span>      <span class="string">\`json:"description"\`</span>
    CategoryId  <span class="type">int</span>         <span class="string">\`json:"category_id"\`</span>
    Brand       <span class="type">string</span>      <span class="string">\`json:"brand"\`</span>
    Status      <span class="type">int</span>         <span class="string">\`json:"status"\`</span>
    Skus        []Sku       <span class="string">\`json:"skus,omitempty"\`</span>
    MinPrice    <span class="type">float64</span>     <span class="string">\`json:"min_price"\`</span>  <span class="comment">// 计算字段</span>
    CreatedAt   *gtime.Time <span class="string">\`json:"created_at"\`</span>
}

<span class="comment">// 订单状态常量</span>
<span class="keyword">const</span> (
    OrderStatusPending   = <span class="number">0</span>  <span class="comment">// 待支付</span>
    OrderStatusPaid      = <span class="number">1</span>  <span class="comment">// 已支付</span>
    OrderStatusShipped   = <span class="number">2</span>  <span class="comment">// 已发货</span>
    OrderStatusCompleted = <span class="number">3</span>  <span class="comment">// 已完成</span>
    OrderStatusCancelled = <span class="number">4</span>  <span class="comment">// 已取消</span>
)` },
      { title: 'Step 2：商品管理与 Redis 缓存', desc: '商品 CRUD、SPU/SKU 管理、多级缓存策略', code: `<span class="comment">// api/v1/product.go — 商品接口定义</span>
<span class="keyword">type</span> <span class="type">ProductListReq</span> <span class="keyword">struct</span> {
    g.Meta     <span class="string">\`path:"/products" method:"get" tags:"商品"\`</span>
    Page       <span class="type">int</span>    <span class="string">\`d:"1" json:"page"\`</span>
    Size       <span class="type">int</span>    <span class="string">\`d:"20" v:"max:100" json:"size"\`</span>
    Keyword    <span class="type">string</span> <span class="string">\`json:"keyword"\`</span>
    CategoryId <span class="type">int</span>    <span class="string">\`json:"category_id"\`</span>
    MinPrice   <span class="type">float64</span> <span class="string">\`json:"min_price"\`</span>
    MaxPrice   <span class="type">float64</span> <span class="string">\`json:"max_price"\`</span>
    SortBy     <span class="type">string</span> <span class="string">\`d:"created_at" v:"in:price,sales,created_at" json:"sort_by"\`</span>
    SortOrder  <span class="type">string</span> <span class="string">\`d:"desc" v:"in:asc,desc" json:"sort_order"\`</span>
}

<span class="comment">// internal/logic/product.go — 商品列表（含缓存）</span>
<span class="keyword">func</span> (s *sProduct) <span class="function">List</span>(ctx context.Context, req *v1.ProductListReq) (*v1.ProductListRes, <span class="type">error</span>) {
    m := g.<span class="function">Model</span>(<span class="string">"products p"</span>).
        <span class="function">LeftJoin</span>(<span class="string">"skus s"</span>, <span class="string">"p.id = s.product_id"</span>).
        <span class="function">Fields</span>(<span class="string">"p.*, MIN(s.price) as min_price"</span>).
        <span class="function">Where</span>(<span class="string">"p.status"</span>, <span class="number">1</span>).
        <span class="function">Group</span>(<span class="string">"p.id"</span>)

    <span class="keyword">if</span> req.Keyword != <span class="string">""</span> {
        m = m.<span class="function">WhereLike</span>(<span class="string">"p.name"</span>, <span class="string">"%"</span>+req.Keyword+<span class="string">"%"</span>)
    }
    <span class="keyword">if</span> req.CategoryId > <span class="number">0</span> {
        m = m.<span class="function">Where</span>(<span class="string">"p.category_id"</span>, req.CategoryId)
    }
    <span class="keyword">if</span> req.MinPrice > <span class="number">0</span> {
        m = m.<span class="function">Having</span>(<span class="string">"min_price >= ?"</span>, req.MinPrice)
    }
    <span class="keyword">if</span> req.MaxPrice > <span class="number">0</span> {
        m = m.<span class="function">Having</span>(<span class="string">"min_price <= ?"</span>, req.MaxPrice)
    }

    total, _ := m.<span class="function">Count</span>()
    <span class="keyword">var</span> products []model.Product
    err := m.<span class="function">Page</span>(req.Page, req.Size).
        <span class="function">Order</span>(req.SortBy + <span class="string">" "</span> + req.SortOrder).
        <span class="function">Scan</span>(&products)

    <span class="keyword">return</span> &v1.ProductListRes{List: products, Total: total}, err
}

<span class="comment">// 商品详情 — 多级缓存策略</span>
<span class="keyword">func</span> (s *sProduct) <span class="function">Detail</span>(ctx context.Context, id <span class="type">int64</span>) (*model.Product, <span class="type">error</span>) {
    cacheKey := fmt.<span class="function">Sprintf</span>(<span class="string">"product:%d"</span>, id)

    <span class="comment">// L1: 本地缓存（进程内，1 分钟）</span>
    <span class="comment">// L2: Redis 缓存（分布式，30 分钟）</span>
    <span class="comment">// L3: 数据库回源</span>
    val, err := g.<span class="function">Cache</span>().<span class="function">GetOrSetFunc</span>(ctx, cacheKey,
        <span class="keyword">func</span>(ctx context.Context) (<span class="keyword">interface</span>{}, <span class="type">error</span>) {
            <span class="keyword">var</span> product model.Product
            err := g.<span class="function">Model</span>(<span class="string">"products"</span>).
                <span class="function">Where</span>(<span class="string">"id"</span>, id).<span class="function">Where</span>(<span class="string">"status"</span>, <span class="number">1</span>).
                <span class="function">Scan</span>(&product)
            <span class="keyword">if</span> err != <span class="keyword">nil</span> { <span class="keyword">return</span> <span class="keyword">nil</span>, err }

            <span class="comment">// 加载 SKU 列表</span>
            g.<span class="function">Model</span>(<span class="string">"skus"</span>).<span class="function">Where</span>(<span class="string">"product_id"</span>, id).
                <span class="function">Order</span>(<span class="string">"price ASC"</span>).<span class="function">Scan</span>(&product.Skus)
            <span class="keyword">return</span> &product, <span class="keyword">nil</span>
        }, <span class="number">30</span> * time.Minute)

    <span class="keyword">if</span> err != <span class="keyword">nil</span> { <span class="keyword">return</span> <span class="keyword">nil</span>, err }
    <span class="keyword">return</span> val.Val().(*model.Product), <span class="keyword">nil</span>
}` },
      { title: 'Step 3：订单状态机与库存扣减', desc: '实现下单流程、订单状态流转、库存并发安全扣减', code: `<span class="comment">// internal/logic/order.go — 下单流程（核心难点）</span>

<span class="comment">// 创建订单（涉及：库存扣减 + 订单创建 + 订单明细 = 事务）</span>
<span class="keyword">func</span> (s *sOrder) <span class="function">Create</span>(ctx context.Context, userId <span class="type">int64</span>, req *v1.OrderCreateReq) (<span class="type">string</span>, <span class="type">error</span>) {
    orderNo := <span class="function">generateOrderNo</span>() <span class="comment">// 生成唯一订单号</span>

    err := g.<span class="function">DB</span>().<span class="function">Transaction</span>(ctx, <span class="keyword">func</span>(ctx context.Context, tx gdb.TX) <span class="type">error</span> {
        totalAmount := <span class="number">0.0</span>
        items := make([]g.Map, <span class="number">0</span>)

        <span class="keyword">for</span> _, item := <span class="keyword">range</span> req.Items {
            <span class="comment">// 1. 查询 SKU 并锁定行（SELECT ... FOR UPDATE）</span>
            sku, err := tx.<span class="function">Model</span>(<span class="string">"skus"</span>).
                <span class="function">Where</span>(<span class="string">"id"</span>, item.SkuId).
                <span class="function">LockUpdate</span>().  <span class="comment">// 悲观锁！防止超卖</span>
                <span class="function">One</span>()
            <span class="keyword">if</span> err != <span class="keyword">nil</span> { <span class="keyword">return</span> err }
            <span class="keyword">if</span> sku.IsEmpty() {
                <span class="keyword">return</span> gerror.<span class="function">Newf</span>(<span class="string">"SKU %d 不存在"</span>, item.SkuId)
            }

            <span class="comment">// 2. 检查库存</span>
            <span class="keyword">if</span> sku[<span class="string">"stock"</span>].Int() < item.Quantity {
                <span class="keyword">return</span> gerror.<span class="function">Newf</span>(<span class="string">"库存不足（剩余 %d）"</span>, sku[<span class="string">"stock"</span>].Int())
            }

            <span class="comment">// 3. 扣减库存（原子操作）</span>
            _, err = tx.<span class="function">Model</span>(<span class="string">"skus"</span>).
                <span class="function">Where</span>(<span class="string">"id"</span>, item.SkuId).
                <span class="function">Where</span>(<span class="string">"stock >= ?"</span>, item.Quantity). <span class="comment">// 双重检查</span>
                <span class="function">Decrement</span>(<span class="string">"stock"</span>, item.Quantity)
            <span class="keyword">if</span> err != <span class="keyword">nil</span> { <span class="keyword">return</span> err }

            <span class="comment">// 4. 构建订单明细</span>
            price := sku[<span class="string">"price"</span>].Float64()
            totalAmount += price * <span class="function">float64</span>(item.Quantity)
            items = <span class="function">append</span>(items, g.Map{
                <span class="string">"sku_id"</span>:       item.SkuId,
                <span class="string">"product_name"</span>: sku[<span class="string">"product_name"</span>].String(),
                <span class="string">"sku_attrs"</span>:    sku[<span class="string">"attrs"</span>].String(),
                <span class="string">"price"</span>:        price,
                <span class="string">"quantity"</span>:     item.Quantity,
            })
        }

        <span class="comment">// 5. 创建订单主表</span>
        orderId, err := tx.<span class="function">Model</span>(<span class="string">"orders"</span>).<span class="function">Data</span>(g.Map{
            <span class="string">"order_no"</span>:     orderNo,
            <span class="string">"user_id"</span>:      userId,
            <span class="string">"total_amount"</span>: totalAmount,
            <span class="string">"status"</span>:       OrderStatusPending,
            <span class="string">"address"</span>:      req.Address,
        }).<span class="function">InsertAndGetId</span>()
        <span class="keyword">if</span> err != <span class="keyword">nil</span> { <span class="keyword">return</span> err }

        <span class="comment">// 6. 批量插入订单明细</span>
        <span class="keyword">for</span> i := <span class="keyword">range</span> items {
            items[i][<span class="string">"order_id"</span>] = orderId
        }
        _, err = tx.<span class="function">Model</span>(<span class="string">"order_items"</span>).<span class="function">Data</span>(items).<span class="function">Insert</span>()
        <span class="keyword">return</span> err
    })

    <span class="keyword">if</span> err != <span class="keyword">nil</span> { <span class="keyword">return</span> <span class="string">""</span>, err }

    <span class="comment">// 7. 设置订单超时自动取消（30 分钟）</span>
    <span class="keyword">go</span> s.<span class="function">scheduleCancel</span>(ctx, orderNo, <span class="number">30</span>*time.Minute)

    <span class="keyword">return</span> orderNo, <span class="keyword">nil</span>
}

<span class="comment">// 生成唯一订单号</span>
<span class="keyword">func</span> <span class="function">generateOrderNo</span>() <span class="type">string</span> {
    <span class="keyword">return</span> time.Now().<span class="function">Format</span>(<span class="string">"20060102150405"</span>) +
        fmt.<span class="function">Sprintf</span>(<span class="string">"%06d"</span>, rand.<span class="function">Intn</span>(<span class="number">999999</span>))
}` },
      { title: 'Step 4：订单状态流转与支付回调', desc: '实现状态机模式管理订单生命周期，处理支付回调', code: `<span class="comment">// internal/logic/order_state.go — 订单状态机</span>

<span class="comment">// 合法的状态转移表（类似前端的状态管理）</span>
<span class="keyword">var</span> validTransitions = <span class="keyword">map</span>[<span class="type">int</span>][]<span class="type">int</span>{
    OrderStatusPending:   {OrderStatusPaid, OrderStatusCancelled},
    OrderStatusPaid:      {OrderStatusShipped, OrderStatusCancelled},
    OrderStatusShipped:   {OrderStatusCompleted},
    OrderStatusCompleted: {},  <span class="comment">// 终态</span>
    OrderStatusCancelled: {},  <span class="comment">// 终态</span>
}

<span class="comment">// 状态转移（带校验）</span>
<span class="keyword">func</span> (s *sOrder) <span class="function">ChangeStatus</span>(ctx context.Context, orderNo <span class="type">string</span>, newStatus <span class="type">int</span>) <span class="type">error</span> {
    order, err := s.<span class="function">GetByOrderNo</span>(ctx, orderNo)
    <span class="keyword">if</span> err != <span class="keyword">nil</span> { <span class="keyword">return</span> err }

    <span class="comment">// 校验状态转移是否合法</span>
    allowed := validTransitions[order.Status]
    isValid := <span class="keyword">false</span>
    <span class="keyword">for</span> _, s := <span class="keyword">range</span> allowed {
        <span class="keyword">if</span> s == newStatus { isValid = <span class="keyword">true</span>; <span class="keyword">break</span> }
    }
    <span class="keyword">if</span> !isValid {
        <span class="keyword">return</span> gerror.<span class="function">Newf</span>(<span class="string">"不允许从 %d 转到 %d"</span>, order.Status, newStatus)
    }

    data := g.Map{<span class="string">"status"</span>: newStatus}
    <span class="keyword">if</span> newStatus == OrderStatusPaid {
        data[<span class="string">"paid_at"</span>] = gtime.<span class="function">Now</span>()
    }

    _, err = g.<span class="function">Model</span>(<span class="string">"orders"</span>).
        <span class="function">Where</span>(<span class="string">"order_no"</span>, orderNo).
        <span class="function">Where</span>(<span class="string">"status"</span>, order.Status). <span class="comment">// 乐观锁：只更新当前状态</span>
        <span class="function">Data</span>(data).<span class="function">Update</span>()
    <span class="keyword">return</span> err
}

<span class="comment">// 支付回调处理</span>
<span class="keyword">func</span> (c *cPayment) <span class="function">Callback</span>(ctx context.Context, req *v1.PayCallbackReq) (res *v1.PayCallbackRes, err <span class="type">error</span>) {
    <span class="comment">// 1. 验证签名（防止伪造回调）</span>
    <span class="keyword">if</span> !<span class="function">verifySign</span>(req) {
        <span class="keyword">return</span> <span class="keyword">nil</span>, gerror.<span class="function">New</span>(<span class="string">"签名验证失败"</span>)
    }

    <span class="comment">// 2. 幂等性检查（防止重复回调）</span>
    lockKey := <span class="string">"pay:lock:"</span> + req.OrderNo
    ok, _ := g.<span class="function">Redis</span>().<span class="function">Do</span>(ctx, <span class="string">"SET"</span>, lockKey, <span class="string">"1"</span>, <span class="string">"NX"</span>, <span class="string">"EX"</span>, <span class="number">60</span>)
    <span class="keyword">if</span> ok.IsNil() {
        g.Log().<span class="function">Info</span>(ctx, <span class="string">"重复回调，忽略"</span>, req.OrderNo)
        <span class="keyword">return</span> &v1.PayCallbackRes{Success: <span class="keyword">true</span>}, <span class="keyword">nil</span>
    }

    <span class="comment">// 3. 更新订单状态</span>
    err = service.<span class="function">Order</span>().<span class="function">ChangeStatus</span>(ctx, req.OrderNo, OrderStatusPaid)
    <span class="keyword">if</span> err != <span class="keyword">nil</span> {
        g.<span class="function">Redis</span>().<span class="function">Do</span>(ctx, <span class="string">"DEL"</span>, lockKey) <span class="comment">// 失败时释放锁</span>
        <span class="keyword">return</span> <span class="keyword">nil</span>, err
    }

    <span class="keyword">return</span> &v1.PayCallbackRes{Success: <span class="keyword">true</span>}, <span class="keyword">nil</span>
}

<span class="comment">// 超时自动取消（goroutine 定时任务）</span>
<span class="keyword">func</span> (s *sOrder) <span class="function">scheduleCancel</span>(ctx context.Context, orderNo <span class="type">string</span>, timeout time.Duration) {
    time.<span class="function">Sleep</span>(timeout)

    order, _ := s.<span class="function">GetByOrderNo</span>(ctx, orderNo)
    <span class="keyword">if</span> order != <span class="keyword">nil</span> && order.Status == OrderStatusPending {
        <span class="comment">// 取消订单 + 恢复库存</span>
        g.<span class="function">DB</span>().<span class="function">Transaction</span>(ctx, <span class="keyword">func</span>(ctx context.Context, tx gdb.TX) <span class="type">error</span> {
            tx.<span class="function">Model</span>(<span class="string">"orders"</span>).<span class="function">Where</span>(<span class="string">"order_no"</span>, orderNo).
                <span class="function">Where</span>(<span class="string">"status"</span>, OrderStatusPending).
                <span class="function">Data</span>(g.Map{<span class="string">"status"</span>: OrderStatusCancelled}).<span class="function">Update</span>()

            <span class="comment">// 恢复库存</span>
            items, _ := tx.<span class="function">Model</span>(<span class="string">"order_items"</span>).
                <span class="function">Where</span>(<span class="string">"order_id"</span>, order.Id).<span class="function">All</span>()
            <span class="keyword">for</span> _, item := <span class="keyword">range</span> items {
                tx.<span class="function">Model</span>(<span class="string">"skus"</span>).<span class="function">Where</span>(<span class="string">"id"</span>, item[<span class="string">"sku_id"</span>]).
                    <span class="function">Increment</span>(<span class="string">"stock"</span>, item[<span class="string">"quantity"</span>].Int())
            }
            <span class="keyword">return</span> <span class="keyword">nil</span>
        })
        g.Log().<span class="function">Info</span>(ctx, <span class="string">"订单超时取消:"</span>, orderNo)
    }
}` },
      { title: 'Step 5：完整路由与 Docker 部署', desc: '路由结构、权限控制、Docker 多阶段构建部署', code: `<span class="comment">// internal/cmd/cmd.go — 电商 API 完整路由</span>
s.<span class="function">Group</span>(<span class="string">"/api/v1"</span>, <span class="keyword">func</span>(group *ghttp.RouterGroup) {
    group.<span class="function">Use</span>(ghttp.MiddlewareHandlerResponse)

    <span class="comment">// --- 公开接口 ---</span>
    group.<span class="function">GET</span>(<span class="string">"/products"</span>, controller.Product.List)
    group.<span class="function">GET</span>(<span class="string">"/products/:id"</span>, controller.Product.Detail)
    group.<span class="function">GET</span>(<span class="string">"/categories"</span>, controller.Category.List)

    <span class="comment">// --- 用户接口（需登录）---</span>
    group.<span class="function">Group</span>(<span class="string">"/"</span>, <span class="keyword">func</span>(group *ghttp.RouterGroup) {
        group.<span class="function">Middleware</span>(middleware.Auth)

        <span class="comment">// 购物车</span>
        group.<span class="function">GET</span>(<span class="string">"/cart"</span>, controller.Cart.List)
        group.<span class="function">POST</span>(<span class="string">"/cart"</span>, controller.Cart.Add)
        group.<span class="function">PUT</span>(<span class="string">"/cart/:id"</span>, controller.Cart.Update)
        group.<span class="function">DELETE</span>(<span class="string">"/cart/:id"</span>, controller.Cart.Remove)

        <span class="comment">// 订单</span>
        group.<span class="function">POST</span>(<span class="string">"/orders"</span>, controller.Order.Create)
        group.<span class="function">GET</span>(<span class="string">"/orders"</span>, controller.Order.MyOrders)
        group.<span class="function">GET</span>(<span class="string">"/orders/:orderNo"</span>, controller.Order.Detail)
        group.<span class="function">POST</span>(<span class="string">"/orders/:orderNo/cancel"</span>, controller.Order.Cancel)
    })

    <span class="comment">// --- 管理员接口 ---</span>
    group.<span class="function">Group</span>(<span class="string">"/admin"</span>, <span class="keyword">func</span>(group *ghttp.RouterGroup) {
        group.<span class="function">Middleware</span>(middleware.Auth, middleware.AdminOnly)

        group.<span class="function">Bind</span>(controller.AdminProduct) <span class="comment">// 商品管理</span>
        group.<span class="function">Bind</span>(controller.AdminOrder)   <span class="comment">// 订单管理</span>
        group.<span class="function">POST</span>(<span class="string">"/orders/:orderNo/ship"</span>, controller.AdminOrder.Ship)
    })

    <span class="comment">// --- 支付回调（第三方调用）---</span>
    group.<span class="function">POST</span>(<span class="string">"/payment/callback"</span>, controller.Payment.Callback)
})

<span class="comment">// ---- Docker 部署 ----</span>

<span class="comment"># Dockerfile（多阶段构建）</span>
<span class="keyword">FROM</span> golang:1.22-alpine AS builder
<span class="keyword">WORKDIR</span> /app
<span class="keyword">COPY</span> go.mod go.sum ./
<span class="keyword">RUN</span> go mod download
<span class="keyword">COPY</span> . .
<span class="keyword">RUN</span> CGO_ENABLED=0 go build -ldflags="-s -w" -o shop .

<span class="keyword">FROM</span> alpine:latest
<span class="keyword">RUN</span> apk --no-cache add ca-certificates tzdata
<span class="keyword">ENV</span> TZ=Asia/Shanghai
<span class="keyword">COPY</span> --from=builder /app/shop /app/shop
<span class="keyword">COPY</span> --from=builder /app/manifest /app/manifest
<span class="keyword">WORKDIR</span> /app
<span class="keyword">EXPOSE</span> 8000
<span class="keyword">CMD</span> ["./shop"]

<span class="comment"># docker-compose.yml</span>
version: <span class="string">"3.8"</span>
services:
  app:
    build: .
    ports: [<span class="string">"8000:8000"</span>]
    depends_on: [mysql, redis]
    environment:
      - GF_DATABASE_DEFAULT_HOST=mysql
      - GF_REDIS_DEFAULT_ADDRESS=redis:6379
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: <span class="string">root123</span>
      MYSQL_DATABASE: <span class="string">shop_db</span>
    volumes: [<span class="string">"mysql_data:/var/lib/mysql"</span>]
  redis:
    image: redis:7-alpine
volumes:
  mysql_data:

<span class="comment"># 一键启动</span>
$ docker-compose up -d
$ curl http://localhost:8000/api/v1/products` },
    ],
  },
  { icon: '💬', title: '实时聊天应用', diff: '高级', diffClass: 'diff-expert', desc: '基于 WebSocket 的实时聊天应用，理解 Go 在实时通信场景中的优势。',
    techs: ['GoFrame', 'WebSocket', 'Redis Pub/Sub', 'Goroutine'],
    features: ['WebSocket 长连接', '消息广播与私聊', 'Goroutine 连接管理', '在线状态追踪'],
    hasGuide: true,
    steps: [
      { title: 'Step 1：连接管理器（Hub）', desc: '设计 WebSocket 连接池，管理所有在线用户连接', code: `<span class="comment">// internal/websocket/hub.go — 连接管理中心</span>
<span class="comment">// 类似前端的 EventEmitter 或 Redux Store</span>

<span class="keyword">package</span> websocket

<span class="keyword">import</span> (
    <span class="string">"sync"</span>
    <span class="string">"github.com/gorilla/websocket"</span>
)

<span class="comment">// 客户端连接</span>
<span class="keyword">type</span> <span class="type">Client</span> <span class="keyword">struct</span> {
    Id       <span class="type">string</span>          <span class="comment">// 用户 ID</span>
    Username <span class="type">string</span>          <span class="comment">// 用户名</span>
    Conn     *websocket.Conn <span class="comment">// WebSocket 连接</span>
    Send     <span class="keyword">chan</span> []<span class="type">byte</span>     <span class="comment">// 发送消息通道</span>
    RoomId   <span class="type">string</span>          <span class="comment">// 当前房间</span>
}

<span class="comment">// 连接管理器（全局单例）</span>
<span class="keyword">type</span> <span class="type">Hub</span> <span class="keyword">struct</span> {
    Clients    <span class="keyword">map</span>[<span class="type">string</span>]*Client    <span class="comment">// userId -> Client</span>
    Rooms      <span class="keyword">map</span>[<span class="type">string</span>]<span class="keyword">map</span>[<span class="type">string</span>]*Client <span class="comment">// roomId -> clients</span>
    Register   <span class="keyword">chan</span> *Client            <span class="comment">// 注册通道</span>
    Unregister <span class="keyword">chan</span> *Client            <span class="comment">// 注销通道</span>
    Broadcast  <span class="keyword">chan</span> *Message            <span class="comment">// 广播通道</span>
    mu         sync.RWMutex
}

<span class="keyword">var</span> hub = &Hub{
    Clients:    <span class="keyword">make</span>(<span class="keyword">map</span>[<span class="type">string</span>]*Client),
    Rooms:      <span class="keyword">make</span>(<span class="keyword">map</span>[<span class="type">string</span>]<span class="keyword">map</span>[<span class="type">string</span>]*Client),
    Register:   <span class="keyword">make</span>(<span class="keyword">chan</span> *Client),
    Unregister: <span class="keyword">make</span>(<span class="keyword">chan</span> *Client),
    Broadcast:  <span class="keyword">make</span>(<span class="keyword">chan</span> *Message, <span class="number">256</span>),
}

<span class="comment">// 启动 Hub（核心事件循环，类似 Node.js 的 Event Loop）</span>
<span class="keyword">func</span> (h *Hub) <span class="function">Run</span>() {
    <span class="keyword">for</span> {
        <span class="keyword">select</span> {
        <span class="keyword">case</span> client := <-h.Register:
            h.mu.<span class="function">Lock</span>()
            h.Clients[client.Id] = client
            <span class="comment">// 加入默认大厅</span>
            h.<span class="function">joinRoom</span>(client, <span class="string">"lobby"</span>)
            h.mu.<span class="function">Unlock</span>()
            <span class="comment">// 广播上线通知</span>
            h.<span class="function">broadcastSystem</span>(<span class="string">"lobby"</span>, client.Username+<span class="string">" 加入了聊天"</span>)

        <span class="keyword">case</span> client := <-h.Unregister:
            h.mu.<span class="function">Lock</span>()
            <span class="keyword">if</span> _, ok := h.Clients[client.Id]; ok {
                h.<span class="function">leaveRoom</span>(client, client.RoomId)
                <span class="keyword">delete</span>(h.Clients, client.Id)
                <span class="keyword">close</span>(client.Send)
            }
            h.mu.<span class="function">Unlock</span>()
            h.<span class="function">broadcastSystem</span>(<span class="string">"lobby"</span>, client.Username+<span class="string">" 离开了聊天"</span>)

        <span class="keyword">case</span> msg := <-h.Broadcast:
            h.<span class="function">handleMessage</span>(msg)
        }
    }
}

<span class="comment">// 房间管理</span>
<span class="keyword">func</span> (h *Hub) <span class="function">joinRoom</span>(client *Client, roomId <span class="type">string</span>) {
    <span class="keyword">if</span> h.Rooms[roomId] == <span class="keyword">nil</span> {
        h.Rooms[roomId] = <span class="keyword">make</span>(<span class="keyword">map</span>[<span class="type">string</span>]*Client)
    }
    h.Rooms[roomId][client.Id] = client
    client.RoomId = roomId
}

<span class="keyword">func</span> (h *Hub) <span class="function">leaveRoom</span>(client *Client, roomId <span class="type">string</span>) {
    <span class="keyword">if</span> room, ok := h.Rooms[roomId]; ok {
        <span class="keyword">delete</span>(room, client.Id)
        <span class="keyword">if</span> <span class="function">len</span>(room) == <span class="number">0</span> {
            <span class="keyword">delete</span>(h.Rooms, roomId)
        }
    }
}` },
      { title: 'Step 2：消息协议与处理', desc: '定义消息格式，实现广播、私聊、房间消息等功能', code: `<span class="comment">// internal/websocket/message.go — 消息协议</span>

<span class="comment">// 消息类型（类似前端 action type）</span>
<span class="keyword">const</span> (
    MsgTypeChat     = <span class="string">"chat"</span>      <span class="comment">// 聊天消息</span>
    MsgTypePrivate  = <span class="string">"private"</span>   <span class="comment">// 私聊消息</span>
    MsgTypeSystem   = <span class="string">"system"</span>    <span class="comment">// 系统通知</span>
    MsgTypeJoin     = <span class="string">"join"</span>      <span class="comment">// 加入房间</span>
    MsgTypeLeave    = <span class="string">"leave"</span>     <span class="comment">// 离开房间</span>
    MsgTypeOnline   = <span class="string">"online"</span>    <span class="comment">// 在线列表</span>
    MsgTypeTyping   = <span class="string">"typing"</span>    <span class="comment">// 正在输入</span>
)

<span class="comment">// 统一消息格式（类似前端 Redux Action）</span>
<span class="keyword">type</span> <span class="type">Message</span> <span class="keyword">struct</span> {
    Type     <span class="type">string</span> <span class="string">\`json:"type"\`</span>
    Content  <span class="type">string</span> <span class="string">\`json:"content"\`</span>
    FromId   <span class="type">string</span> <span class="string">\`json:"from_id"\`</span>
    FromName <span class="type">string</span> <span class="string">\`json:"from_name"\`</span>
    ToId     <span class="type">string</span> <span class="string">\`json:"to_id,omitempty"\`</span>    <span class="comment">// 私聊目标</span>
    RoomId   <span class="type">string</span> <span class="string">\`json:"room_id,omitempty"\`</span>
    Time     <span class="type">int64</span>  <span class="string">\`json:"time"\`</span>
}

<span class="comment">// 消息路由处理（核心分发逻辑）</span>
<span class="keyword">func</span> (h *Hub) <span class="function">handleMessage</span>(msg *Message) {
    msg.Time = time.Now().UnixMilli()

    <span class="keyword">switch</span> msg.Type {
    <span class="keyword">case</span> MsgTypeChat:
        <span class="comment">// 房间广播（发给同一房间的所有人）</span>
        h.<span class="function">sendToRoom</span>(msg.RoomId, msg)

    <span class="keyword">case</span> MsgTypePrivate:
        <span class="comment">// 私聊（只发给目标用户 + 发送者自己）</span>
        h.<span class="function">sendToUser</span>(msg.ToId, msg)
        h.<span class="function">sendToUser</span>(msg.FromId, msg)

    <span class="keyword">case</span> MsgTypeJoin:
        h.mu.<span class="function">Lock</span>()
        <span class="keyword">if</span> client, ok := h.Clients[msg.FromId]; ok {
            h.<span class="function">leaveRoom</span>(client, client.RoomId)
            h.<span class="function">joinRoom</span>(client, msg.RoomId)
        }
        h.mu.<span class="function">Unlock</span>()
        h.<span class="function">broadcastSystem</span>(msg.RoomId, msg.FromName+<span class="string">" 加入了房间"</span>)
        h.<span class="function">sendOnlineList</span>(msg.RoomId)

    <span class="keyword">case</span> MsgTypeLeave:
        h.mu.<span class="function">Lock</span>()
        <span class="keyword">if</span> client, ok := h.Clients[msg.FromId]; ok {
            h.<span class="function">leaveRoom</span>(client, msg.RoomId)
            h.<span class="function">joinRoom</span>(client, <span class="string">"lobby"</span>)
        }
        h.mu.<span class="function">Unlock</span>()
        h.<span class="function">broadcastSystem</span>(msg.RoomId, msg.FromName+<span class="string">" 离开了房间"</span>)

    <span class="keyword">case</span> MsgTypeTyping:
        <span class="comment">// 转发"正在输入"状态给目标用户</span>
        h.<span class="function">sendToUser</span>(msg.ToId, msg)
    }
}

<span class="comment">// 发送给房间内所有用户</span>
<span class="keyword">func</span> (h *Hub) <span class="function">sendToRoom</span>(roomId <span class="type">string</span>, msg *Message) {
    h.mu.<span class="function">RLock</span>()
    <span class="keyword">defer</span> h.mu.<span class="function">RUnlock</span>()

    data, _ := json.<span class="function">Marshal</span>(msg)
    <span class="keyword">if</span> room, ok := h.Rooms[roomId]; ok {
        <span class="keyword">for</span> _, client := <span class="keyword">range</span> room {
            <span class="keyword">select</span> {
            <span class="keyword">case</span> client.Send <- data:
            <span class="keyword">default</span>:
                <span class="comment">// 发送缓冲区满了，关闭连接</span>
                <span class="keyword">close</span>(client.Send)
                <span class="keyword">delete</span>(room, client.Id)
            }
        }
    }
}

<span class="comment">// 发送给指定用户</span>
<span class="keyword">func</span> (h *Hub) <span class="function">sendToUser</span>(userId <span class="type">string</span>, msg *Message) {
    h.mu.<span class="function">RLock</span>()
    <span class="keyword">defer</span> h.mu.<span class="function">RUnlock</span>()

    <span class="keyword">if</span> client, ok := h.Clients[userId]; ok {
        data, _ := json.<span class="function">Marshal</span>(msg)
        client.Send <- data
    }
}

<span class="comment">// 广播在线用户列表</span>
<span class="keyword">func</span> (h *Hub) <span class="function">sendOnlineList</span>(roomId <span class="type">string</span>) {
    h.mu.<span class="function">RLock</span>()
    <span class="keyword">defer</span> h.mu.<span class="function">RUnlock</span>()

    users := make([]g.Map, <span class="number">0</span>)
    <span class="keyword">if</span> room, ok := h.Rooms[roomId]; ok {
        <span class="keyword">for</span> _, c := <span class="keyword">range</span> room {
            users = <span class="function">append</span>(users, g.Map{
                <span class="string">"id"</span>: c.Id, <span class="string">"username"</span>: c.Username,
            })
        }
    }
    data, _ := json.<span class="function">Marshal</span>(g.Map{
        <span class="string">"type"</span>: MsgTypeOnline, <span class="string">"users"</span>: users,
    })
    h.<span class="function">sendToRoom</span>(roomId, &Message{Type: MsgTypeOnline, Content: <span class="function">string</span>(data)})
}` },
      { title: 'Step 3：WebSocket 读写协程', desc: '每个连接启动两个 goroutine 负责读和写，实现高效消息处理', code: `<span class="comment">// internal/websocket/client.go — 客户端读写</span>

<span class="comment">// 读协程：从 WebSocket 读消息 -> 发到 Hub.Broadcast</span>
<span class="comment">// （类似 socket.on('message', callback)）</span>
<span class="keyword">func</span> (c *Client) <span class="function">ReadPump</span>() {
    <span class="keyword">defer</span> <span class="keyword">func</span>() {
        hub.Unregister <- c
        c.Conn.<span class="function">Close</span>()
    }()

    c.Conn.<span class="function">SetReadLimit</span>(<span class="number">4096</span>)
    c.Conn.<span class="function">SetReadDeadline</span>(time.<span class="function">Now</span>().<span class="function">Add</span>(<span class="number">60</span> * time.Second))
    c.Conn.<span class="function">SetPongHandler</span>(<span class="keyword">func</span>(<span class="type">string</span>) <span class="type">error</span> {
        c.Conn.<span class="function">SetReadDeadline</span>(time.<span class="function">Now</span>().<span class="function">Add</span>(<span class="number">60</span> * time.Second))
        <span class="keyword">return</span> <span class="keyword">nil</span>
    })

    <span class="keyword">for</span> {
        _, data, err := c.Conn.<span class="function">ReadMessage</span>()
        <span class="keyword">if</span> err != <span class="keyword">nil</span> {
            <span class="keyword">break</span> <span class="comment">// 连接断开</span>
        }

        <span class="keyword">var</span> msg Message
        <span class="keyword">if</span> err := json.<span class="function">Unmarshal</span>(data, &msg); err != <span class="keyword">nil</span> {
            <span class="keyword">continue</span>
        }
        <span class="comment">// 注入发送者信息</span>
        msg.FromId = c.Id
        msg.FromName = c.Username
        <span class="keyword">if</span> msg.RoomId == <span class="string">""</span> {
            msg.RoomId = c.RoomId
        }
        hub.Broadcast <- &msg
    }
}

<span class="comment">// 写协程：从 Send channel 读消息 -> 写入 WebSocket</span>
<span class="comment">// （类似 socket.emit）</span>
<span class="keyword">func</span> (c *Client) <span class="function">WritePump</span>() {
    ticker := time.<span class="function">NewTicker</span>(<span class="number">30</span> * time.Second)
    <span class="keyword">defer</span> <span class="keyword">func</span>() {
        ticker.<span class="function">Stop</span>()
        c.Conn.<span class="function">Close</span>()
    }()

    <span class="keyword">for</span> {
        <span class="keyword">select</span> {
        <span class="keyword">case</span> msg, ok := <-c.Send:
            <span class="keyword">if</span> !ok {
                <span class="comment">// Hub 关闭了 Send channel</span>
                c.Conn.<span class="function">WriteMessage</span>(websocket.CloseMessage, []<span class="type">byte</span>{})
                <span class="keyword">return</span>
            }
            c.Conn.<span class="function">SetWriteDeadline</span>(time.<span class="function">Now</span>().<span class="function">Add</span>(<span class="number">10</span> * time.Second))
            c.Conn.<span class="function">WriteMessage</span>(websocket.TextMessage, msg)

        <span class="keyword">case</span> <-ticker.C:
            <span class="comment">// 定期心跳保持连接</span>
            c.Conn.<span class="function">SetWriteDeadline</span>(time.<span class="function">Now</span>().<span class="function">Add</span>(<span class="number">10</span> * time.Second))
            <span class="keyword">if</span> err := c.Conn.<span class="function">WriteMessage</span>(websocket.PingMessage, <span class="keyword">nil</span>); err != <span class="keyword">nil</span> {
                <span class="keyword">return</span>
            }
        }
    }
}

<span class="comment">// 每个连接的架构示意：</span>
<span class="comment">//</span>
<span class="comment">// ┌───────────────┐</span>
<span class="comment">// │   Client 浏览器  │</span>
<span class="comment">// └──────┬────────┘</span>
<span class="comment">//        │ WebSocket</span>
<span class="comment">// ┌──────▼────────┐</span>
<span class="comment">// │    Conn 连接    │</span>
<span class="comment">// ├───────────────┤</span>
<span class="comment">// │  ReadPump()   │──▶ Hub.Broadcast channel</span>
<span class="comment">// │  (goroutine)  │</span>
<span class="comment">// ├───────────────┤</span>
<span class="comment">// │  WritePump()  │◀── Client.Send channel</span>
<span class="comment">// │  (goroutine)  │</span>
<span class="comment">// └───────────────┘</span>
<span class="comment">//</span>
<span class="comment">// 优势：每个连接只占 ~4KB 内存</span>
<span class="comment">// 10 万连接 ≈ 400MB，Node.js 同等规模需 2-4GB</span>` },
      { title: 'Step 4：路由注册与 Redis Pub/Sub', desc: '注册 WebSocket 端点，使用 Redis 实现多实例消息同步', code: `<span class="comment">// internal/cmd/cmd.go — WebSocket 路由注册</span>
<span class="keyword">var</span> Main = gcmd.Command{
    Name: <span class="string">"main"</span>,
    Func: <span class="keyword">func</span>(ctx context.Context, parser *gcmd.Parser) (err <span class="type">error</span>) {
        <span class="comment">// 启动 Hub 事件循环</span>
        <span class="keyword">go</span> websocket.GetHub().<span class="function">Run</span>()
        <span class="comment">// 启动 Redis 订阅</span>
        <span class="keyword">go</span> websocket.<span class="function">StartRedisSubscriber</span>(ctx)

        s := g.<span class="function">Server</span>()

        <span class="comment">// REST API</span>
        s.<span class="function">Group</span>(<span class="string">"/api/v1"</span>, <span class="keyword">func</span>(group *ghttp.RouterGroup) {
            group.<span class="function">POST</span>(<span class="string">"/login"</span>, controller.Auth.Login)
            group.<span class="function">GET</span>(<span class="string">"/rooms"</span>, controller.Room.List)
            group.<span class="function">GET</span>(<span class="string">"/messages/:roomId"</span>, controller.Message.History)
        })

        <span class="comment">// WebSocket 端点</span>
        s.<span class="function">BindHandler</span>(<span class="string">"/ws"</span>, <span class="keyword">func</span>(r *ghttp.Request) {
            <span class="comment">// 1. 验证 Token（从 query 参数获取）</span>
            token := r.<span class="function">Get</span>(<span class="string">"token"</span>).String()
            claims, err := service.<span class="function">ParseToken</span>(token)
            <span class="keyword">if</span> err != <span class="keyword">nil</span> {
                r.Response.<span class="function">WriteStatus</span>(<span class="number">401</span>, <span class="string">"认证失败"</span>)
                <span class="keyword">return</span>
            }

            <span class="comment">// 2. 升级为 WebSocket</span>
            ws, err := r.<span class="function">WebSocket</span>()
            <span class="keyword">if</span> err != <span class="keyword">nil</span> {
                g.Log().<span class="function">Error</span>(ctx, <span class="string">"WebSocket 升级失败:"</span>, err)
                <span class="keyword">return</span>
            }

            <span class="comment">// 3. 创建客户端并注册到 Hub</span>
            client := &websocket.Client{
                Id:       fmt.<span class="function">Sprintf</span>(<span class="string">"%d"</span>, claims.UserId),
                Username: claims.Username,
                Conn:     ws.Conn,
                Send:     <span class="keyword">make</span>(<span class="keyword">chan</span> []<span class="type">byte</span>, <span class="number">256</span>),
            }
            websocket.GetHub().Register <- client

            <span class="comment">// 4. 启动读写协程（1 连接 = 2 goroutine）</span>
            <span class="keyword">go</span> client.<span class="function">WritePump</span>()
            client.<span class="function">ReadPump</span>() <span class="comment">// 阻塞直到连接关闭</span>
        })

        s.<span class="function">Run</span>()
        <span class="keyword">return</span>
    },
}

<span class="comment">// ---- Redis Pub/Sub — 多实例消息同步 ----</span>
<span class="comment">// 当部署多个 Go 实例时，消息需要跨实例广播</span>
<span class="comment">// 类似前端 BroadcastChannel API</span>

<span class="keyword">func</span> <span class="function">StartRedisSubscriber</span>(ctx context.Context) {
    conn, _ := g.<span class="function">Redis</span>().<span class="function">Conn</span>(ctx)
    <span class="keyword">defer</span> conn.<span class="function">Close</span>(ctx)

    <span class="comment">// 订阅频道</span>
    _, err := conn.<span class="function">Do</span>(ctx, <span class="string">"SUBSCRIBE"</span>, <span class="string">"chat:broadcast"</span>)
    <span class="keyword">if</span> err != <span class="keyword">nil</span> { <span class="keyword">return</span> }

    <span class="keyword">for</span> {
        reply, err := conn.<span class="function">Receive</span>(ctx)
        <span class="keyword">if</span> err != <span class="keyword">nil</span> { <span class="keyword">break</span> }

        <span class="keyword">var</span> msg Message
        <span class="keyword">if</span> json.<span class="function">Unmarshal</span>(reply.Bytes(), &msg) == <span class="keyword">nil</span> {
            GetHub().<span class="function">handleMessage</span>(&msg)
        }
    }
}

<span class="comment">// 发布消息到 Redis（其他实例会收到）</span>
<span class="keyword">func</span> <span class="function">PublishMessage</span>(ctx context.Context, msg *Message) {
    data, _ := json.<span class="function">Marshal</span>(msg)
    g.<span class="function">Redis</span>().<span class="function">Do</span>(ctx, <span class="string">"PUBLISH"</span>, <span class="string">"chat:broadcast"</span>, <span class="function">string</span>(data))
}

<span class="comment">// 前端连接示例（JavaScript）</span>
<span class="comment">// const ws = new WebSocket('ws://localhost:8000/ws?token=xxx');</span>
<span class="comment">// ws.onmessage = (e) => {</span>
<span class="comment">//   const msg = JSON.parse(e.data);</span>
<span class="comment">//   switch (msg.type) {</span>
<span class="comment">//     case 'chat':    appendMessage(msg); break;</span>
<span class="comment">//     case 'online':  updateUserList(msg.users); break;</span>
<span class="comment">//     case 'typing':  showTyping(msg.from_name); break;</span>
<span class="comment">//   }</span>
<span class="comment">// };</span>
<span class="comment">// ws.send(JSON.stringify({type:'chat', content:'Hello!'}));</span>` },
      { title: 'Step 5：消息持久化与历史记录', desc: '将聊天消息存入数据库，支持历史消息分页查询', code: `<span class="comment">-- 消息存储表</span>
<span class="keyword">CREATE TABLE</span> messages (
    id         <span class="type">BIGINT</span> PRIMARY KEY AUTO_INCREMENT,
    room_id    <span class="type">VARCHAR</span>(50) NOT NULL,
    from_id    <span class="type">VARCHAR</span>(50) NOT NULL,
    from_name  <span class="type">VARCHAR</span>(100),
    type       <span class="type">VARCHAR</span>(20) DEFAULT <span class="string">'chat'</span>,
    content    <span class="type">TEXT</span> NOT NULL,
    created_at <span class="type">DATETIME</span> DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_room_time (room_id, created_at DESC)
);

<span class="keyword">CREATE TABLE</span> rooms (
    id          <span class="type">VARCHAR</span>(50) PRIMARY KEY,
    name        <span class="type">VARCHAR</span>(100) NOT NULL,
    type        <span class="type">TINYINT</span> DEFAULT 0 <span class="comment">-- 0:群聊 1:私聊</span>,
    created_by  <span class="type">VARCHAR</span>(50),
    created_at  <span class="type">DATETIME</span> DEFAULT CURRENT_TIMESTAMP
);

<span class="comment">// internal/logic/message.go — 消息持久化</span>

<span class="comment">// 异步保存消息（不阻塞实时投递）</span>
<span class="keyword">var</span> saveCh = <span class="keyword">make</span>(<span class="keyword">chan</span> *Message, <span class="number">1000</span>) <span class="comment">// 缓冲通道</span>

<span class="keyword">func</span> <span class="function">init</span>() {
    <span class="comment">// 启动多个 worker 消费消息</span>
    <span class="keyword">for</span> i := <span class="number">0</span>; i < <span class="number">4</span>; i++ {
        <span class="keyword">go</span> <span class="function">saveWorker</span>()
    }
}

<span class="keyword">func</span> <span class="function">saveWorker</span>() {
    batch := make([]g.Map, <span class="number">0</span>, <span class="number">50</span>)
    ticker := time.<span class="function">NewTicker</span>(<span class="number">2</span> * time.Second)

    <span class="keyword">for</span> {
        <span class="keyword">select</span> {
        <span class="keyword">case</span> msg := <-saveCh:
            batch = <span class="function">append</span>(batch, g.Map{
                <span class="string">"room_id"</span>:   msg.RoomId,
                <span class="string">"from_id"</span>:   msg.FromId,
                <span class="string">"from_name"</span>: msg.FromName,
                <span class="string">"type"</span>:      msg.Type,
                <span class="string">"content"</span>:   msg.Content,
            })
            <span class="comment">// 攒满 50 条批量插入</span>
            <span class="keyword">if</span> <span class="function">len</span>(batch) >= <span class="number">50</span> {
                <span class="function">flushBatch</span>(&batch)
            }
        <span class="keyword">case</span> <-ticker.C:
            <span class="comment">// 定时刷入（即使不满 50 条）</span>
            <span class="keyword">if</span> <span class="function">len</span>(batch) > <span class="number">0</span> {
                <span class="function">flushBatch</span>(&batch)
            }
        }
    }
}

<span class="keyword">func</span> <span class="function">flushBatch</span>(batch *[]g.Map) {
    ctx := context.<span class="function">Background</span>()
    _, err := g.<span class="function">Model</span>(<span class="string">"messages"</span>).<span class="function">Data</span>(*batch).<span class="function">Insert</span>()
    <span class="keyword">if</span> err != <span class="keyword">nil</span> {
        g.Log().<span class="function">Error</span>(ctx, <span class="string">"批量保存消息失败:"</span>, err)
    }
    *batch = (*batch)[:0] <span class="comment">// 清空 slice</span>
}

<span class="comment">// 查询历史消息（分页，游标方式）</span>
<span class="keyword">func</span> (s *sMessage) <span class="function">History</span>(ctx context.Context, roomId <span class="type">string</span>, beforeId <span class="type">int64</span>, limit <span class="type">int</span>) ([]model.Message, <span class="type">error</span>) {
    m := g.<span class="function">Model</span>(<span class="string">"messages"</span>).
        <span class="function">Where</span>(<span class="string">"room_id"</span>, roomId)

    <span class="keyword">if</span> beforeId > <span class="number">0</span> {
        m = m.<span class="function">Where</span>(<span class="string">"id < ?"</span>, beforeId) <span class="comment">// 游标分页（比 offset 高效）</span>
    }

    <span class="keyword">var</span> msgs []model.Message
    err := m.<span class="function">Order</span>(<span class="string">"id DESC"</span>).<span class="function">Limit</span>(limit).<span class="function">Scan</span>(&msgs)
    <span class="keyword">return</span> msgs, err
}

<span class="comment">// API 端点</span>
<span class="comment">// GET /api/v1/messages/lobby?before_id=100&limit=30</span>
<span class="comment">// 返回：</span>
<span class="comment">// {</span>
<span class="comment">//   "list": [</span>
<span class="comment">//     {"id":99, "from_name":"张三", "content":"你好", "time":"..."},</span>
<span class="comment">//     {"id":98, "from_name":"李四", "content":"在吗", "time":"..."},</span>
<span class="comment">//     ...</span>
<span class="comment">//   ]</span>
<span class="comment">// }</span>` },
    ],
  },
];

// FAQ 数据
const FAQS = [
  { q: '前端转 Go 难度大吗？需要多久？', a: '如果你有 TypeScript 经验，转 Go 会相对顺畅。TS 的静态类型系统和 Go 理念相似。通常 2-3 个月可以达到独立开发后端服务的水平。Go 的语法非常简洁，只有 25 个关键字，学习曲线比你想象的要平缓。' },
  { q: 'GoFrame 和 Gin 应该选哪个？', a: '<strong>Gin</strong> 是轻量级路由框架，需要自行组合各种组件（类似 Express）。<strong>GoFrame</strong> 是全功能企业级框架，内置 ORM、配置管理、日志等组件（类似 NestJS/Spring）。如果你做企业级项目，GoFrame 能大幅提升效率；如果只需简单 API，Gin 也是好选择。' },
  { q: 'Go 的错误处理为什么不用 try/catch？', a: 'Go 的设计哲学是"显式优于隐式"。通过返回值处理错误，让每个可能出错的地方都被明确处理，避免了异常在调用栈中"隐式传播"的问题。虽然代码量多了 <code>if err != nil</code>，但代码的可预测性和可维护性大幅提升。Go 1.13+ 引入的 errors.Is/As 也让错误处理更加灵活。' },
  { q: 'Goroutine 和 JavaScript 的 async/await 有什么本质区别？', a: 'JavaScript 是<strong>单线程 + 事件循环</strong>，async/await 本质是非阻塞 I/O，不能利用多核 CPU。Goroutine 是<strong>真正的并发/并行</strong>，Go 运行时会自动将 goroutine 调度到多个 OS 线程上，充分利用多核 CPU。一个 goroutine 只需约 2KB 栈空间，可以轻松创建数百万个。' },
  { q: 'Go 没有 class，怎么实现面向对象？', a: 'Go 用<strong>结构体 + 方法 + 接口</strong>实现面向对象，用<strong>组合</strong>代替<strong>继承</strong>。这种设计更灵活，避免了深层继承带来的复杂性。在 Go 社区中，"偏好组合，而非继承"是核心设计原则。接口的隐式实现也让代码更加解耦。' },
  { q: 'Go 的包管理和 npm 有什么区别？', a: '<code>go mod</code> 是 Go 的包管理工具，类似 npm。<code>go.mod</code> = <code>package.json</code>，<code>go.sum</code> = <code>package-lock.json</code>。主要区别：Go 模块版本基于 Git tag，直接从代码仓库拉取，没有中心化的 registry（虽然有 proxy.golang.org 加速）。使用 <code>go get</code> 安装依赖，<code>go mod tidy</code> 清理无用依赖。' },
  { q: 'GoFrame 的性能怎么样？', a: 'GoFrame 的底层基于 Go 的 net/http 包，性能非常出色。Go 编译后是原生二进制文件，启动时间通常在毫秒级，内存占用极低。在 TechEmpower 等基准测试中，Go 框架通常比 Node.js 框架快 5-10 倍。对于高并发场景，Go 的 goroutine 模型比 Node.js 的事件循环更有优势。' },
  { q: '前端开发者学 Go 应该重点关注什么？', a: '建议重点关注：<strong>1)</strong> 静态类型系统（类比 TypeScript）；<strong>2)</strong> 错误处理模式；<strong>3)</strong> 接口和组合（替代继承）；<strong>4)</strong> goroutine 和 channel（并发编程）；<strong>5)</strong> 指针的基本概念。不必急于学习底层系统编程，先把 Web 开发这条线学通。' },
  { q: '如何搭建 Go 开发环境？', a: '<strong>1)</strong> 安装 Go：从 <a href="https://go.dev/dl/" target="_blank">go.dev</a> 下载安装；<strong>2)</strong> IDE 推荐 VS Code + Go 扩展（或 GoLand）；<strong>3)</strong> 安装 GoFrame CLI：<code>go install github.com/gogf/gf/cmd/gf/v2@latest</code>；<strong>4)</strong> 设置 Go proxy（国内）：<code>go env -w GOPROXY=https://goproxy.cn,direct</code>。' },
  { q: 'Go 适合做什么？不适合做什么？', a: '<strong>适合</strong>：Web 后端、微服务、CLI 工具、云原生（Docker/K8s 就是 Go 写的）、高并发服务。<strong>不太适合</strong>：GUI 桌面应用、数据科学/机器学习（Python 更好）、前端开发。Go 是后端领域的利器，特别适合需要高性能和高并发的场景。' },
];

// ==================== 渲染函数 ====================

function renderLessons() {
  const container = document.getElementById('lessonContents');
  if (!container) return;
  let html = '';
  for (const [key, l] of Object.entries(LESSONS)) {
    const isActive = key === 'variables' ? ' active' : '';
    html += `<div class="lesson-content${isActive}" id="lesson-${key}">
      <div class="lesson-intro"><div class="lesson-icon">${l.icon}</div><div><h3>${l.title}</h3><p>${l.desc}</p></div></div>
      <div class="key-points"><h4>💡 核心要点</h4><ul>${l.points.map(p => `<li><strong>${p[0]}</strong>：${p[1]}</li>`).join('')}</ul></div>
      <div class="code-comparison">
        <div class="code-panel"><div class="code-panel-header"><div class="code-lang"><span class="lang-dot js-dot"></span>JavaScript / TypeScript</div><button class="copy-btn" onclick="copyCode(this)">复制</button></div><pre class="code-block"><code>${l.jsCode}</code></pre></div>
        <div class="code-vs"><div class="vs-badge">VS</div></div>
        <div class="code-panel"><div class="code-panel-header"><div class="code-lang"><span class="lang-dot go-dot"></span>Go</div><button class="copy-btn" onclick="copyCode(this)">复制</button></div><pre class="code-block"><code>${l.goCode}</code></pre></div>
      </div>
      <div class="try-it-section"><h4>🚀 动手试试</h4><p>在编辑器中修改代码，点击运行查看结果（通过 Go Playground 执行）：</p>
        <div class="interactive-editor"><div class="editor-header"><span class="lang-dot go-dot"></span><span>main.go</span><button class="run-btn" onclick="runCode('${key}')">▶ 运行</button></div>
          <textarea class="code-editor" id="editor-${key}" spellcheck="false">${l.editorCode}</textarea>
          <div class="editor-output" id="output-${key}"><div class="output-placeholder">点击 "运行" 按钮查看输出结果</div></div>
        </div>
      </div>
    </div>`;
  }
  container.innerHTML = html;
}

function renderGfLessons() {
  const container = document.getElementById('gfLessons');
  if (!container) return;
  container.innerHTML = GF_LESSONS.map(l => `
    <div class="gf-lesson-card" id="${l.id}">
      <div class="gf-lesson-header" onclick="toggleGfLesson('${l.id}')">
        <div class="gf-lesson-title"><span class="gf-lesson-num">${l.num}</span><h3>${l.title}</h3></div>
        <span class="gf-chevron">▼</span>
      </div>
      <div class="gf-lesson-body">${l.content}</div>
    </div>`).join('');
}

function renderProjects() {
  const container = document.getElementById('projectsGrid');
  if (!container) return;
  container.innerHTML = PROJECTS.map((p, pi) => `
    <div class="project-card ${p.hasGuide ? 'has-guide' : ''}">
      <div class="project-header"><span class="project-icon">${p.icon}</span><span class="project-difficulty ${p.diffClass}">${p.diff}</span></div>
      <h3>${p.title}</h3>
      <p>${p.desc}</p>
      <div class="project-techs">${p.techs.map(t => `<span>${t}</span>`).join('')}</div>
      <ul class="project-features">${p.features.map(f => `<li>${f}</li>`).join('')}</ul>
      ${p.hasGuide ? `
        <button class="project-guide-btn" onclick="toggleProjectGuide(${pi})">
          <span>📖 查看实现指南</span><span class="guide-chevron" id="guide-chevron-${pi}">▼</span>
        </button>
        <div class="project-guide" id="project-guide-${pi}">
          <div class="guide-steps">
            ${p.steps.map((s, si) => `
              <div class="guide-step">
                <div class="step-header" onclick="toggleProjectStep(${pi}, ${si})">
                  <div class="step-marker">${si + 1}</div>
                  <div class="step-info"><h4>${s.title}</h4><p>${s.desc}</p></div>
                  <span class="step-chevron" id="step-chevron-${pi}-${si}">▶</span>
                </div>
                <div class="step-code" id="step-code-${pi}-${si}">
                  <div class="code-panel full-width"><div class="code-panel-header"><div class="code-lang"><span class="lang-dot go-dot"></span>Go</div><button class="copy-btn" onclick="copyCode(this)">复制</button></div><pre class="code-block"><code>${s.code}</code></pre></div>
                </div>
              </div>
            `).join('')}
          </div>
        </div>
      ` : `<div class="project-coming-soon"><span>🚧</span> 详细实现指南即将推出</div>`}
    </div>`).join('');
}

function toggleProjectGuide(pi) {
  const guide = document.getElementById(`project-guide-${pi}`);
  const chevron = document.getElementById(`guide-chevron-${pi}`);
  if (!guide) return;
  const isOpen = guide.classList.toggle('open');
  if (chevron) chevron.textContent = isOpen ? '▲' : '▼';
}

function toggleProjectStep(pi, si) {
  const code = document.getElementById(`step-code-${pi}-${si}`);
  const chevron = document.getElementById(`step-chevron-${pi}-${si}`);
  if (!code) return;
  const isOpen = code.classList.toggle('open');
  if (chevron) chevron.textContent = isOpen ? '▼' : '▶';
}

function renderFaqs() {
  const container = document.getElementById('faqList');
  if (!container) return;
  container.innerHTML = FAQS.map((f, i) => `
    <div class="faq-item" id="faq-${i}">
      <div class="faq-question" onclick="toggleFaq(${i})"><span class="q-icon">Q</span><span class="q-text">${f.q}</span><span class="faq-chevron">▼</span></div>
      <div class="faq-answer">${f.a}</div>
    </div>`).join('');
}

// ==================== 交互逻辑 ====================

// 课程 Tab 切换
function initLessonTabs() {
  document.querySelectorAll('.lesson-tab').forEach(tab => {
    tab.addEventListener('click', () => {
      document.querySelectorAll('.lesson-tab').forEach(t => t.classList.remove('active'));
      document.querySelectorAll('.lesson-content').forEach(c => c.classList.remove('active'));
      tab.classList.add('active');
      const lesson = tab.getAttribute('data-lesson');
      const content = document.getElementById(`lesson-${lesson}`);
      if (content) content.classList.add('active');
      markCompleted(lesson);
    });
  });
}

// GoFrame 手风琴
function toggleGfLesson(id) {
  const card = document.getElementById(id);
  if (!card) return;
  card.classList.toggle('open');
  markCompleted(id);
}

// FAQ 手风琴
function toggleFaq(index) {
  const item = document.getElementById(`faq-${index}`);
  if (!item) return;
  item.classList.toggle('open');
}

// 学习进度追踪
const completedItems = new Set(JSON.parse(localStorage.getItem('gf-progress') || '[]'));

function markCompleted(id) {
  completedItems.add(id);
  localStorage.setItem('gf-progress', JSON.stringify([...completedItems]));
  updateProgress();
}

function updateProgress() {
  const total = 15;
  const done = Math.min(completedItems.size, total);
  const el = document.getElementById('progressText');
  if (el) el.textContent = `${done} / ${total}`;
}

// 代码复制
function copyCode(btn) {
  const code = btn.closest('.code-panel').querySelector('code');
  if (!code) return;
  const text = code.textContent;
  navigator.clipboard.writeText(text).then(() => {
    btn.textContent = '已复制 ✓';
    setTimeout(() => btn.textContent = '复制', 2000);
  });
}

// 运行代码 (Go Playground)
async function runCode(lesson) {
  const editor = document.getElementById(`editor-${lesson}`);
  const output = document.getElementById(`output-${lesson}`);
  const btn = editor?.closest('.interactive-editor')?.querySelector('.run-btn');
  if (!editor || !output) return;

  if (btn) { btn.classList.add('loading'); btn.textContent = '⏳ 运行中...'; }
  output.innerHTML = '<span style="color:var(--text-muted)">正在编译运行...</span>';

  try {
    const resp = await fetch('https://play.golang.org/compile', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: 'version=2&body=' + encodeURIComponent(editor.value) + '&withVet=true',
    });
    const data = await resp.json();
    if (data.Errors) {
      output.innerHTML = `<span class="output-error">${escapeHtml(data.Errors)}</span>`;
    } else {
      let text = '';
      if (data.Events) {
        data.Events.forEach(e => { text += e.Message; });
      }
      output.innerHTML = `<span style="color:var(--success)">${escapeHtml(text || '(无输出)')}</span>`;
    }
  } catch (err) {
    output.innerHTML = `<span class="output-error">网络错误，请检查网络连接后重试。\n提示：也可以将代码粘贴到 https://go.dev/play/ 运行</span>`;
  }
  if (btn) { btn.classList.remove('loading'); btn.textContent = '▶ 运行'; }
  markCompleted(lesson);
}

function escapeHtml(str) {
  const map = { '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#039;' };
  return str.replace(/[&<>"']/g, c => map[c]);
}

// 主题切换
function initTheme() {
  const toggle = document.getElementById('themeToggle');
  const saved = localStorage.getItem('theme') || 'dark';
  document.documentElement.setAttribute('data-theme', saved);
  updateThemeIcon(saved);

  toggle?.addEventListener('click', () => {
    const current = document.documentElement.getAttribute('data-theme');
    const next = current === 'dark' ? 'light' : 'dark';
    document.documentElement.setAttribute('data-theme', next);
    localStorage.setItem('theme', next);
    updateThemeIcon(next);
  });
}

function updateThemeIcon(theme) {
  const btn = document.getElementById('themeToggle');
  if (btn) btn.textContent = theme === 'dark' ? '☀️' : '🌙';
}

// 移动端菜单
function initMobileMenu() {
  const btn = document.getElementById('mobileMenuBtn');
  const links = document.getElementById('navLinks');
  btn?.addEventListener('click', () => links?.classList.toggle('open'));
  links?.querySelectorAll('.nav-link').forEach(link => {
    link.addEventListener('click', () => links.classList.remove('open'));
  });
}

// 导航高亮
function initNavHighlight() {
  const sections = document.querySelectorAll('.section');
  const links = document.querySelectorAll('.nav-link');
  const observer = new IntersectionObserver(entries => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        const id = entry.target.getAttribute('id');
        links.forEach(l => {
          l.classList.toggle('active', l.getAttribute('href') === `#${id}`);
        });
      }
    });
  }, { rootMargin: '-30% 0px -70% 0px' });
  sections.forEach(s => observer.observe(s));
}

// 回到顶部
function initBackToTop() {
  const btn = document.getElementById('backToTop');
  window.addEventListener('scroll', () => {
    btn?.classList.toggle('visible', window.scrollY > 400);
  });
  btn?.addEventListener('click', () => window.scrollTo({ top: 0, behavior: 'smooth' }));
}

// 滚动动画
function initScrollAnimation() {
  const items = document.querySelectorAll('.roadmap-item');
  const observer = new IntersectionObserver(entries => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        entry.target.classList.add('visible');
      }
    });
  }, { threshold: 0.2 });
  items.forEach(item => observer.observe(item));
}

// Hero 背景粒子
function initHeroCanvas() {
  const canvas = document.getElementById('heroCanvas');
  if (!canvas) return;
  const ctx = canvas.getContext('2d');
  let w, h, particles = [];

  function resize() {
    w = canvas.width = canvas.offsetWidth;
    h = canvas.height = canvas.offsetHeight;
  }
  resize();
  window.addEventListener('resize', resize);

  for (let i = 0; i < 60; i++) {
    particles.push({
      x: Math.random() * w, y: Math.random() * h,
      vx: (Math.random() - 0.5) * 0.5, vy: (Math.random() - 0.5) * 0.5,
      r: Math.random() * 2 + 1, o: Math.random() * 0.5 + 0.1,
    });
  }

  function draw() {
    ctx.clearRect(0, 0, w, h);
    const isDark = document.documentElement.getAttribute('data-theme') !== 'light';
    const color = isDark ? '0, 173, 216' : '0, 130, 170';
    particles.forEach(p => {
      p.x += p.vx; p.y += p.vy;
      if (p.x < 0 || p.x > w) p.vx *= -1;
      if (p.y < 0 || p.y > h) p.vy *= -1;
      ctx.beginPath();
      ctx.arc(p.x, p.y, p.r, 0, Math.PI * 2);
      ctx.fillStyle = `rgba(${color}, ${p.o})`;
      ctx.fill();
    });
    // Draw lines between nearby particles
    for (let i = 0; i < particles.length; i++) {
      for (let j = i + 1; j < particles.length; j++) {
        const dx = particles[i].x - particles[j].x;
        const dy = particles[i].y - particles[j].y;
        const dist = Math.sqrt(dx * dx + dy * dy);
        if (dist < 150) {
          ctx.beginPath();
          ctx.moveTo(particles[i].x, particles[i].y);
          ctx.lineTo(particles[j].x, particles[j].y);
          ctx.strokeStyle = `rgba(${color}, ${0.1 * (1 - dist / 150)})`;
          ctx.stroke();
        }
      }
    }
    requestAnimationFrame(draw);
  }
  draw();
}

// 导航栏滚动效果
function initNavbarScroll() {
  const navbar = document.getElementById('navbar');
  window.addEventListener('scroll', () => {
    if (navbar) {
      navbar.style.boxShadow = window.scrollY > 50 ? 'var(--shadow-md)' : 'none';
    }
  });
}

// Tab 键支持代码编辑器
function initEditorTab() {
  document.querySelectorAll('.code-editor').forEach(editor => {
    editor.addEventListener('keydown', e => {
      if (e.key === 'Tab') {
        e.preventDefault();
        const start = editor.selectionStart;
        const end = editor.selectionEnd;
        editor.value = editor.value.substring(0, start) + '    ' + editor.value.substring(end);
        editor.selectionStart = editor.selectionEnd = start + 4;
      }
    });
  });
}

// ==================== 初始化 ====================
document.addEventListener('DOMContentLoaded', () => {
  renderLessons();
  renderGfLessons();
  renderProjects();
  renderFaqs();
  initLessonTabs();
  initTheme();
  initMobileMenu();
  initNavHighlight();
  initBackToTop();
  initScrollAnimation();
  initHeroCanvas();
  initNavbarScroll();
  initEditorTab();
  updateProgress();
});
