package model

// DictionaryItem — одна запись справочника в JSON-представлении (camelCase).
type DictionaryItem = map[string]any

// DictionaryList — результат постраничной выборки.
type DictionaryList struct {
	Items []DictionaryItem
	Total int
}
