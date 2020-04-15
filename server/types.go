package server

type (
	Data struct {
		UserKey     string                 `json:"userKey"`
		Data        string                 `json:"data"`
		ContentType string                 `json:"contentType"`
		Meta        map[string]interface{} `json:"meta"`
		CreatedAt   int64                  `json:"-"`
	}
	DataInfo struct {
		UserKey     string                 `json:"-"`
		Data        string                 `json:"data"`
		ContentType string                 `json:"-"`
		Meta        map[string]interface{} `json:"meta"`
		CreatedAt   int64                  `json:"created_at"`
	}
)
