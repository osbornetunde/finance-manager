package core

import (
	"encoding/json"
	"time"
)

type Transaction struct {
	ID          int64           `json:"id"`
	UserID      int64           `json:"user_id"`
	Amount      int64           `json:"amount"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date"`
	CategoryID  int64           `json:"category_id"`
	Metadata    json.RawMessage `json:"metadata"`
	Tags        []string        `json:"tags"`
}
