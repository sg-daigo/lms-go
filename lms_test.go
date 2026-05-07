package lms

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func getLogger() *slog.Logger {
	level := new(slog.LevelVar)
	level.Set(slog.LevelDebug)
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(handler)

}

var logger = getLogger()

// テスト用サーバーの起動
func NewTestServer(t *testing.T, mockResponse string) string {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, mockResponse)
	}))

	t.Cleanup(func() {
		ts.Close()
	})

	// サーバーのURLをモックサーバーに向ける
	return ts.URL
}

func TestPlayers(t *testing.T) {
	server := "http://taro-iot:9000"
	req := NewPlayersRequest()
	result, err := SendRequest[PlayersResult](server, req, logger)
	if err != nil {
		t.Error("Unexpected error: %w", err)
	}

	fmt.Println("count =", result.Count)
	for i := 0; i < result.Count; i++ {
		player := result.Players[i]
		fmt.Println("ID =", player.PlayerID)
		fmt.Println("Name =", player.Name)
		fmt.Println("Address =", player.IP)
	}
}

func TestPlayersMock(t *testing.T) {
	// 期待するレスポンスを定義
	mockResponse := `{
		"id": "test-id",
		"method": "slim.request",
		"params": ["", ["players", "0", "100"]],
		"result": {"count": 1, "players_loop": [{"playerid": "00:00:00:00:00:00", "name": "test", "ip": "0.0.0.0"}]}
	}`

	// テスト用サーバーの起動
	server := NewTestServer(t, mockResponse)

	// 実行
	request := NewPlayersRequest()
	request.ID = "test-id" // IDを固定して比較しやすくする

	res, err := SendRequest[PlayersResult](server, request, nil)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}

	// 検証
	if res.Count != 1 {
		t.Errorf("Expected count 1, got %d", res.Count)
	}
	if res.Players[0].Name != "test" {
		t.Errorf("Expected name test, got %s", res.Players[0].Name)
	}
}
