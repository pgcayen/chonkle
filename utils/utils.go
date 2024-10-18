package utils

import "golang.org/x/exp/rand"

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckOK(ok bool, msg string) {
	if !ok {
		panic(msg)
	}
}

func CoalesceString(s string, ifEmpty string) string {
	if s == "" {
		return ifEmpty
	}

	return s
}

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
