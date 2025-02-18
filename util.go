package dbpg

import (
	"fmt"
	"regexp"
	"strings"
)

var reCreateTable = regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?([^\s(]+)`)

type Pair[T, U any] struct {
	A T
	B U
}

type KeyValPair[T, U any] struct {
	Key T
	Val U
}

type Map[K comparable, V any] map[K]V

func (dict Map[K, V]) Keys() []K {
	var keys = make([]K, 0, len(dict))
	for k := range dict {
		keys = append(keys, k)
	}
	return keys
}

func (dict Map[K, V]) Values() []V {
	var vals = make([]V, 0, len(dict))
	for _, v := range dict {
		vals = append(vals, v)
	}
	return vals
}

func (dict Map[K, V]) Flatten() []KeyValPair[K, V] {
	var pairs = make([]KeyValPair[K, V], 0, len(dict))
	for k, v := range dict {
		pairs = append(pairs, KeyValPair[K, V]{k, v})
	}
	return pairs
}

func MapFn[T any, U any](slice []T, fn func(T) U) []U {
	var result = make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

func KeysToMap[K comparable, V any](keys []K, val V) Map[K, V] {
	var dict = make(map[K]V, len(keys))
	for _, k := range keys {
		dict[k] = val
	}
	return dict
}

func ColumnNames(cols []string) string {
	return strings.Join(cols, ",")
}

func ColumnPlaceholders(cols []string) string {
	placeholders := make([]string, len(cols))
	for i := range cols {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	return strings.Join(placeholders, ", ")
}

func UpdatePlaceholders(cols []string) string {
	placeholders := make([]string, len(cols))
	for i, col := range cols {
		placeholders[i] = fmt.Sprintf("%s=$%d", col, i+1)
	}
	return strings.Join(placeholders, ", ")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func TableNameFromCreateSql(sql string) (string, error) {
	var matches = reCreateTable.FindStringSubmatch(sql)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("table name not found")
}
