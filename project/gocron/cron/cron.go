// cron/cron.go
// 6 字段 cron 表达式解析器 (秒 分 时 日 月 周)
package cron

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Schedule struct {
	Second []int
	Minute []int
	Hour   []int
	Day    []int
	Month  []int
	Weekday []int
}

// Parse 解析 6 字段 cron 表达式: 秒 分 时 日 月 周
func Parse(expr string) (*Schedule, error) {
	fields := strings.Fields(expr)
	if len(fields) != 6 {
		return nil, fmt.Errorf("cron 需要 6 个字段 (秒 分 时 日 月 周), 实际 %d 个", len(fields))
	}

	s := &Schedule{}
	var err error

	if s.Second, err = parseField(fields[0], 0, 59); err != nil {
		return nil, fmt.Errorf("秒: %w", err)
	}
	if s.Minute, err = parseField(fields[1], 0, 59); err != nil {
		return nil, fmt.Errorf("分: %w", err)
	}
	if s.Hour, err = parseField(fields[2], 0, 23); err != nil {
		return nil, fmt.Errorf("时: %w", err)
	}
	if s.Day, err = parseField(fields[3], 1, 31); err != nil {
		return nil, fmt.Errorf("日: %w", err)
	}
	if s.Month, err = parseField(fields[4], 1, 12); err != nil {
		return nil, fmt.Errorf("月: %w", err)
	}
	if s.Weekday, err = parseField(fields[5], 0, 6); err != nil {
		return nil, fmt.Errorf("周: %w", err)
	}

	return s, nil
}

// Next 计算从 now 之后的下一个触发时间
func (s *Schedule) Next(now time.Time) time.Time {
	t := now.Add(time.Second).Truncate(time.Second)

	for i := 0; i < 366*24*60*60; i++ {
		if s.match(t) {
			return t
		}
		t = t.Add(time.Second)
	}
	return time.Time{}
}

func (s *Schedule) match(t time.Time) bool {
	return contains(s.Month, int(t.Month())) &&
		contains(s.Day, t.Day()) &&
		contains(s.Weekday, int(t.Weekday())) &&
		contains(s.Hour, t.Hour()) &&
		contains(s.Minute, t.Minute()) &&
		contains(s.Second, t.Second())
}

func contains(vals []int, v int) bool {
	for _, x := range vals {
		if x == v {
			return true
		}
	}
	return false
}

// parseField 解析单个 cron 字段，支持 *, */n, n-m, n,m,k
func parseField(field string, min, max int) ([]int, error) {
	var result []int

	for _, part := range strings.Split(field, ",") {
		vals, err := parsePart(part, min, max)
		if err != nil {
			return nil, err
		}
		result = append(result, vals...)
	}

	// 去重
	seen := map[int]bool{}
	var unique []int
	for _, v := range result {
		if !seen[v] {
			seen[v] = true
			unique = append(unique, v)
		}
	}
	return unique, nil
}

func parsePart(part string, min, max int) ([]int, error) {
	// */n
	if strings.HasPrefix(part, "*/") {
		step, err := strconv.Atoi(part[2:])
		if err != nil || step <= 0 {
			return nil, fmt.Errorf("无效步长: %s", part)
		}
		return rangeVals(min, max, step), nil
	}

	// *
	if part == "*" {
		return rangeVals(min, max, 1), nil
	}

	// n-m or n-m/step
	if strings.Contains(part, "-") {
		rangeParts := strings.SplitN(part, "/", 2)
		bounds := strings.SplitN(rangeParts[0], "-", 2)
		lo, err1 := strconv.Atoi(bounds[0])
		hi, err2 := strconv.Atoi(bounds[1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("无效范围: %s", part)
		}
		step := 1
		if len(rangeParts) > 1 {
			step, _ = strconv.Atoi(rangeParts[1])
			if step <= 0 {
				step = 1
			}
		}
		return rangeVals(lo, hi, step), nil
	}

	// 单个值
	v, err := strconv.Atoi(part)
	if err != nil {
		return nil, fmt.Errorf("无效值: %s", part)
	}
	if v < min || v > max {
		return nil, fmt.Errorf("值 %d 超出范围 [%d, %d]", v, min, max)
	}
	return []int{v}, nil
}

func rangeVals(min, max, step int) []int {
	var r []int
	for i := min; i <= max; i += step {
		r = append(r, i)
	}
	return r
}
