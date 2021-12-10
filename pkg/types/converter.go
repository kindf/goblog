package types

import (
    "strconv" // 字符串和其他类型转换
)

func Int64ToString(num int64) string {
    return strconv.FormatInt(num, 10)
}
