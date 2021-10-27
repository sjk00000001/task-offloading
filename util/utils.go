package util

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"unsafe"
)

func DecimalFormat(val float64) float64 {
	val, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", val), 64)
	return val
}

func String2Bytes(s string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: stringHeader.Data,
		Len:  stringHeader.Len,
		Cap:  stringHeader.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Bytes2String(b []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}
	return *(*string)(unsafe.Pointer(&sh))
}

/* make a two-dimensional array slice */
func MakeSlicesFloat64(r int, c int) [][]float64 {
	arr := make([][]float64, r)
	for i := range arr {
		arr[i] = make([]float64, c)
	}
	return arr
}

/* generate random numbers in the specified interval */
func IntervalRandGenerator(arr [2]float64) float64 {
	return arr[0] + rand.Float64()*(arr[1]-arr[0])
}

/* determine whether the string is an IP address plus a port */
func IsPortedIP(s []string) bool {
	portedIP := "(\\d{1,2}|1\\d\\d|2[0-4]\\d|25[0-5])\\.(\\d{1,2}|1\\d\\d|2[0-4]\\d|25[0-5])\\.(\\d{1,2}" +
		"|1\\d\\d|2[0-4]\\d|25[0-5])\\.(\\d{1,2}|1\\d\\d|2[0-4]\\d|25[0-5])\\:([0-9]|[1-9]\\d{1,3}|[1-5]\\d{4}|6[0-5]{2}[0-3][0-5])"
	return judge(s, portedIP)
}

/* determine whether the string is an IP address */
func IsIP(s []string) bool {
	IP := "((?:(?:25[0-5]|2[0-4]\\d|[01]?\\d?\\d)\\.){3}(?:25[0-5]|2[0-4]\\d|[01]?\\d?\\d))"
	return judge(s, IP)
}

/* determine whether each string in s fits the regexp */
func judge(s []string, reg string) bool {
	match := true
	for i := range s {
		b, _ := regexp.MatchString(reg, s[i])
		match = match && b
	}
	return match
}

/* print the result of task offloading decision */
func PrintResult(m map[int]int, f float64, i int, d int) {
	fmt.Print("[TO Evaluation] Task offloading decision map: ")
	PrintMap(m)
	fmt.Print("[TO Evaluation] Total cost: ", f)
	fmt.Println()
	fmt.Print("[TO Evaluation] Total iterations: ", i, " times")
	fmt.Println()
	fmt.Print("[TO Evaluation] Execution duration: ", d, "ms")
	fmt.Println()
}

/* print key-value pairs */
func PrintMap(m map[int]int) {
	for val := range m {
		fmt.Print("(", val, ",", m[val], ") ")
	}
	fmt.Println()
}

/* print new error */
func PrintError(s string) {
	fmt.Println(errors.New(s))
}
