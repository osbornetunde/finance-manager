package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	appErrors "finance-manager/internal/errors"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

const transactionDedupTTL = 60 * time.Second

type CreateTransactionRequest struct {
	UserID      int64           `json:"user_id"`
	Amount      int64           `json:"amount"`
	CategoryID  int64           `json:"category_id"`
	Description string          `json:"description"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	Tags        []string        `json:"tags,omitempty"`
}

func (a *API) GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := a.service.GetTransactions(ctx)
	if err != nil {
		a.logger.Error("Failed to fetch transactions", "error", err)
		a.httpError(w, http.StatusInternalServerError, "Failed to fetch transactions")
		return
	}

	a.jsonResponse(w, http.StatusOK, res)
}

func (a *API) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var body CreateTransactionRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		a.httpError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		a.httpError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	requestKey := strings.TrimSpace(r.Header.Get("X-Deduplication-Key"))
	if requestKey == "" {
		requestKey = "fp:" + buildTransactionFingerprint(body)
	} else {
		requestKey = "dedup:" + requestKey
	}
	if a.dedup.seenRecently(requestKey, transactionDedupTTL) {
		w.Header().Set("Retry-After", strconv.FormatInt(int64(transactionDedupTTL/time.Second), 10))
		a.logger.Debug("Duplicate transaction request blocked", "key", requestKey)
		a.httpError(w, http.StatusConflict, "Duplicate transaction request")
		return
	}

	res, err := a.service.CreateTransaction(ctx, body.UserID, body.Amount, body.CategoryID, body.Description, body.Metadata, body.Tags)
	if err != nil {
		if appErrors.IsValidationError(err) {
			a.dedup.forget(requestKey)
			a.logger.Debug("Validation failed", "error", err)
			a.httpError(w, http.StatusBadRequest, err.Error())
			return
		}
		if appErrors.IsInvalidUserReference(err) || appErrors.IsInvalidCategoryReference(err) {
			a.dedup.forget(requestKey)
			a.logger.Debug("Invalid transaction reference", "error", err)
			a.httpError(w, http.StatusBadRequest, err.Error())
			return
		}
		a.logger.Error("Failed to create transaction", "error", err)
		a.httpError(w, http.StatusInternalServerError, "Failed to create transaction")
		return
	}
	a.jsonResponse(w, http.StatusCreated, res)
}

func buildTransactionFingerprint(body CreateTransactionRequest) string {
	var normalizedMetadata any
	if len(body.Metadata) > 0 {
		if err := json.Unmarshal(body.Metadata, &normalizedMetadata); err != nil {
			normalizedMetadata = strings.TrimSpace(string(body.Metadata))
		}
	}

	tags := make([]string, len(body.Tags))
	copy(tags, body.Tags)
	slices.Sort(tags)

	payload := struct {
		UserID      int64    `json:"user_id"`
		Amount      int64    `json:"amount"`
		CategoryID  int64    `json:"category_id"`
		Description string   `json:"description"`
		Metadata    any      `json:"metadata,omitempty"`
		Tags        []string `json:"tags,omitempty"`
	}{
		UserID:      body.UserID,
		Amount:      body.Amount,
		CategoryID:  body.CategoryID,
		Description: strings.TrimSpace(body.Description),
		Metadata:    normalizedMetadata,
		Tags:        tags,
	}

	b, _ := json.Marshal(payload)
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
