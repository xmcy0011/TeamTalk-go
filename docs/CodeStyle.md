# 命名规范

## 通用命名规则

```go
int priceCountReader     // 无缩写  
int numErrors            // "num" 是一个常见的写法  
int numDnsConnections    // 人人都知道 "DNS" 是什么  
int n                    // 毫无意义.  
int nerr                 // 含糊不清的缩写.  
int nCompConns           // 含糊不清的缩写.  
int wgcConnections       // 只有贵团队知道是什么意思.  
int pcReader             // "pc" 有太多可能的解释了.  
int cstmrId              // 删减了若干字母.  
```


## 模块(文件夹)命名

```go
-io           // 全小写
 |--ioutil
-strconv      // 可缩写
```

## 文件命名

```go
my_useful_class.go     // 推荐
my-useful-class.go     // 可接受
myusefulclass.go       // 可接受
myusefulclass_test.go  // 可接受
```

## 类型命名

```go
// 类和结构体
type UrlTable struct { ...
type UrlTableTester struct { ...
type UrlTableProperties struct { ...

// 类型定义
type PropertiesMap map[string]string
```

## 变量命名

### 普通变量命名

```go
var tableName       // 好 - 首字母小写.
var TableName       // 好 - 首字母大写.
var tablename       // 差 - 全小写.
var TABLENAME       // 差 - 全大写.

var localVariable   // 外部模块不可访问全局变量
var PublicVariable  // 可访问全局变量

type Student struct {
    localVariable string  // 外部模块不可访问成员变量
    PublicVariable string // 可访问成员变量
}
```

## 常量命名

```go
const KDAYS_IN_A_WEEK = 7 // 命名时以 “K” 开头, 大写加"_"混合
```

## 函数命名

一般来说, 函数名的每个单词首字母大写 (即 “驼峰变量名” 或 “帕斯卡变量名”), 没有下划线. 对于首字母缩写的单词, 更倾向于将它们视作一个单词进行首字母大写 (例如, 写作 StartRpc() 而非 StartRPC()).

```go
func AddTableEntry(){}
func DeleteUrl(){}

// 对于首字母缩写的单词, 更倾向于将它们视作一个单词进行首字母大写 (例如, 写作 StartRpc() 而非 StartRPC()).
func StartRpc(){} 

// 同样的命名规则同时适用于类作用域与命名空间作用域的常量, 因为它们是作为 API 的一部分暴露对外的, 因此应当让它们看起来像是一个函数, 因为在这时, 它们实际上是一个对象而非函数的这一事实对外不过是一个无关紧要的实现细节.
func OpenFileOrDie(){}  
```

## 枚举命名

```go

// 先声明类型
type PolicyType int32

const (
    POLICY_MIN      PolicyType = 0  // 全大写
    POLICY_MAX      PolicyType = 1
    POLICY_MID      PolicyType = 2
    POLICY_AVG      PolicyType = 3
)
```