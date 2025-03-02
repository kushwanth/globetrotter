package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var errorBytes []byte

const badRequest = "Bad Request"

func serverListenerAddress() string {
	addrHost := strings.TrimSpace(os.Getenv("HOST"))
	addrPort := strings.TrimSpace(os.Getenv("PORT"))
	addrPortNum, addrPortNumErr := strconv.Atoi(addrPort)
	if net.ParseIP(addrHost) == nil {
		addrHost = "127.0.0.1"
	}
	if addrPortNumErr != nil || !(addrPortNum > 1024 && addrPortNum < 65536) {
		addrPort = "8080"
	}
	return net.JoinHostPort(addrHost, addrPort)
}

func toJson(data any) []byte {
	responseJson, err := json.Marshal(data)
	if err != nil {
		return errorBytes
	} else {
		return responseJson
	}
}

func MD5Hash(s string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(s))
	return hasher.Sum(nil)
}

func getMD5Hex(s string) string {
	return hex.EncodeToString(MD5Hash(s))
}

func shuffleOptions(arr [4]QuestionOption) {
	rSrc := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rSrc)
	r.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
}

func getHttpRateLimit() int {
	customRateLimit, customRateLimitErr := strconv.Atoi(strings.TrimSpace(os.Getenv("HTTP_RATE_LIMIT")))
	if customRateLimitErr == nil {
		return customRateLimit
	} else {
		return 9
	}
}

func isSelectedOptionCorrect(city City, selectedOption QuestionOption) bool {
	isLabelEqual := strings.EqualFold(city.City, selectedOption.Label)
	userSelectedHash, hashErr := hex.DecodeString(selectedOption.Id)
	if hashErr != nil {
		fmt.Println(hashErr.Error())
		return false
	}
	cityHash := MD5Hash(city.City)
	isHashesEqual := bytes.Equal(cityHash, userSelectedHash)
	return isHashesEqual && isLabelEqual
}
