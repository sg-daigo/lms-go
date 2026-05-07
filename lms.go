/*
Lyrion Music Server(LMS)をJSON/RPCで操作するライブラリ
*/

package lms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// CLIリクエスト
type LMSRequest struct {
	ID     string `json:"id"`
	Method string `json:"method"`
	Params []any  `json:"params"`
}

// CLIレスポンス（共通）
type LMSResponse struct {
	ID     string          `json:"id"`
	Method string          `json:"method"`
	Params []any           `json:"params"`
	Result json.RawMessage `json:"result"`
}

// CLIレスポンス（Players）
type PlayersResult struct {
	Count   int `json:"count"`
	Players []struct {
		PlayerID string `json:"playerid"`
		Name     string `json:"name"`
		IP       string `json:"ip"`
	} `json:"players_loop"`
}

// CLIリクエストを生成する
func NewRequest(playerID string, command ...string) LMSRequest {
	return LMSRequest{
		ID:     uuid.NewString(),
		Method: "slim.request",
		Params: []any{
			playerID,
			command,
		},
	}
}

// Playersリクエストを生成する
func NewPlayersRequest() LMSRequest {
	return NewRequest("", "players", "0", "100")
}

// LMSにリクエストを送信し、レスポンスを受ける
func SendRequest[T any](server string, req LMSRequest, logger *slog.Logger) (result T, err error) {
	b, err := json.Marshal(req)
	if err != nil {
		return result, fmt.Errorf("Json marshal faild: %w", err)
	}

	url := server + "/jsonrpc.js"
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return result, fmt.Errorf("Call LMS API faild: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("Unexpected error: %w", err)
	}

	if logger != nil {
		logger.Debug(string(body))
	}

	var res LMSResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return result, fmt.Errorf("Json unmarshal faild: %w", err)
	}

	if logger != nil {
		logger.Debug("Response: " + string(res.Result))
	}

	if res.ID != req.ID {
		return result, fmt.Errorf("Expected ID to %s, but received %s", req.ID, res.ID)
	}

	if err = json.Unmarshal(res.Result, &result); err != nil {
		return result, fmt.Errorf("Json unmarshal faild: %w", err)
	}

	return result, nil
}
