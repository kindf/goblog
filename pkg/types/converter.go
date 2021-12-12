package types

import (
    "strconv" // 字符串和其他类型转换
    "goblog/pkg/logger"
)

func Int64ToString(num int64) string {
    return strconv.FormatInt(num, 10)
}

func StringToUint64(str string) uint64 {
    i, err := strconv.ParseUint(str, 10, 64)
    if err != nil {
        logger.LogError(err)
    }
    return i
}

func Uint64ToString(num uint64) string {
    return strconv.FormatUint(num, 10)
}

