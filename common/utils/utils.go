package utils

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/lovelly/leaf/log"
)

func DumpsMapStringInterface(data map[string]interface{}) string {
	bytes, err := json.Marshal(&data)
	if err != nil {
		log.Error("Dumps data:%v, error%v", data, err.Error())
		return ""
	}
	log.Debug("Dumps %v", string(bytes))
	return string(bytes)
}

func SubString(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}

func LoadsMapStringInterface(str string, data map[string]interface{}) error {
	err := json.Unmarshal([]byte(str), &data)
	return err
}

func Load2Obj(filePath string, obj interface{}) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewDecoder(f)
	return enc.Decode(obj)
}

func Dump2File(filePath string, obj interface{}) error {
	os.Remove(filePath)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	return enc.Encode(obj)
}

func DecodeZip(b []byte) ([]byte, error) {
	reader := flate.NewReader(bytes.NewReader(b))
	return ioutil.ReadAll(reader)
}

func EncodeZip(b []byte) ([]byte, error) {
	var buf bytes.Buffer

	writer, err := flate.NewWriter(&buf, flate.BestSpeed)
	if err != nil {
		return b, err
	}
	defer writer.Close()
	_, err = writer.Write(b)
	if err != nil {
		return b, err
	}
	err = writer.Flush()
	if err != nil {
		return b, err
	}
	return buf.Bytes(), nil
}

func DecodeGzip(b []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	defer reader.Close()
	return ioutil.ReadAll(reader)
}

func EncodeGzip(b []byte) ([]byte, error) {
	var buf bytes.Buffer

	writer, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if err != nil {
		return b, err
	}
	defer writer.Close()
	_, err = writer.Write(b)
	if err != nil {
		return b, err
	}
	err = writer.Flush()
	if err != nil {
		return b, err
	}
	return buf.Bytes(), nil
}

func DecodeBase64(b []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return data, err
	}
	return data, nil
}

func EncodeBase64(data []byte) []byte {
	str := base64.StdEncoding.EncodeToString(data)
	return []byte(str)
}

func GetSqlInStr(someList []int) string {

	str := GetStrFromIntList(someList)
	if strings.TrimSpace(str) == "" {
		return "0"
	} else {
		return str
	}
}

func GetStrIntList(str string, split ...string) []int {
	var sper string
	if len(split) == 1 {
		sper = split[0]
	} else {
		sper = ","
	}
	intList := make([]int, 0)
	for _, v := range strings.Split(str, sper) {
		intV, _ := strconv.Atoi(v)
		intList = append(intList, intV)
	}
	return intList
}

func GetStrIntSixteenList(str string, split ...string) []int {
	var sper string
	if len(split) == 1 {
		sper = split[0]
	} else {
		sper = "，"
	}
	intList := make([]int, 0)
	for _, v := range strings.Split(str, sper) {
		intV, _ := strconv.ParseInt(v, 16, 32)
		intList = append(intList, int(intV))
	}
	return intList
}

func GetStrFromIntList(data []int) string {
	strList := make([]string, 0)
	for _, v := range data {
		strList = append(strList, strconv.Itoa(v))
	}
	return strings.Join(strList, ",")
}

func GetStrFloatList(str string) []float64 {
	floatList := make([]float64, 0)
	for _, v := range strings.Split(str, ",") {
		floatV, _ := strconv.ParseFloat(v, 10)
		floatList = append(floatList, floatV)
	}
	return floatList
}

func GetIntListFromStrList(strlist []string) []int {

	intList := make([]int, 0, len(strlist))

	for _, str := range strlist {
		n, err := strconv.Atoi(str)
		if nil == err {
			intList = append(intList, n)
		}

	}

	return intList
}

func ReadFile(fileName string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(fileName)
	return bytes, err

}

func WriteFile(fileName string, content []byte) error {
	return ioutil.WriteFile(fileName, content, os.ModePerm)
}

func IsInIntSlice(s []int, seach int) bool {
	for _, v := range s {
		if v == seach {
			return true
		}
	}
	return false
}

/*
计算组合(无序)
calcComb:4:
 [[4] [3 1] [2 1 1] [1 1 1 1] [2 2]]
calcComb:3:
 [[3] [2 1] [1 1 1]]
*/
func CalcComb(total int) [][]int {
	multRes := calcCombImpl(total)
	for i, cb := range multRes {
		sort.Sort(sort.Reverse(sort.IntSlice(cb)))
		multRes[i] = cb
	}

	var res [][]int
	for _, cb := range multRes {
		find := false
		for _, resCB := range res {
			if len(cb) != len(resCB) ||
				!reflect.DeepEqual(cb, resCB) {
				continue
			}
			find = true
			break
		}

		if !find {
			res = append(res, cb)
		}
	}
	return res
}

func calcCombImpl(total int) [][]int {
	if total <= 0 {
		return nil
	}

	combs := [][]int{[]int{total}}
	for i := 1; i <= total/2; i++ {
		nextValue := total - i
		for _, cb := range calcCombImpl(nextValue) {
			cb = append(cb, i)
			combs = append(combs, cb)
		}
	}
	return combs
}

func IntSliceDelete(s []int, index int) []int {
	if index == 0 {
		return s[1:]
	}

	if index == len(s)-1 {
		return s[:index]
	}

	return append(s[:index], s[index+1:]...)
}

func FilterEmoji(content string) string {
	new_content := ""
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			new_content += string(value)
		}
	}
	return new_content
}

// pascalcase 2 camelcase
func TranslatePascal(name string) (string, error) {
	if len(name) <= 0 {
		return "", errors.New("Param is empty.")
	}
	name = strings.Title(name)
	for strings.Index(name, "_") != -1 {
		idx := strings.Index(name, "_")
		name = name[:idx] + strings.Title(name[idx+1:])
	}

	return name, nil
}
