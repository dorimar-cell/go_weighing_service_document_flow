package dto

// DictionaryListResponse — контракт GET /spa/api/weighing/dictionaries/{code}.
type DictionaryListResponse struct {
	Items    []map[string]any `json:"items"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
	Total    int              `json:"total"`
}

// DictionaryItemResponse — контракт POST/PATCH: одна запись.
type DictionaryItemResponse = map[string]any
