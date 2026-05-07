# lms-go

Lyrion Music Server (LMS) を JSON/RPC で操作するための Go ライブラリです。

## 概要

このライブラリは、[Lyrion Music Server](https://lyrion.org/) が提供する JSON/RPC API をシンプルに呼び出すためのユーティリティを提供します。リクエストの構築・送信・レスポンスのデコードを汎用的に扱えます。

## インストール

```bash
>go get github.com/sg-daigo/lms-go
```

## 使い方

### プレイヤー一覧の取得

```go
package main

import (
    "fmt"
    "log/slog"
    "os"

    "github.com/sg-daigo/lms-go"
)

func main() {
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))

    req := lms.NewPlayersRequest()
    result, err := lms.SendRequest[lms.PlayersResult]("http://localhost:9000", req, logger)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Players (%d):\n", result.Count)
    for _, p := range result.Players {
        fmt.Printf("  - %s (%s) @ %s\n", p.Name, p.PlayerID, p.IP)
    }
}
```

### 任意のコマンドを送信する

```go
// 特定のプレイヤーで再生を開始する
req := lms.NewRequest("aa:bb:cc:dd:ee:ff", "play")
// SendRequest の型パラメータにレスポンスの構造体を指定する
result, err := lms.SendRequest[map[string]any]("http://localhost:9000", req, nil)
```

## API リファレンス

### 型

#### `LMSRequest`

JSON/RPC リクエストを表す構造体です。`NewRequest` または `NewPlayersRequest` で生成します。

| フィールド | 型       | 説明                          |
|------------|----------|-------------------------------|
| `ID`       | `string` | リクエストを識別する UUID      |
| `Method`   | `string` | 常に `"slim.request"`         |
| `Params`   | `[]any`  | `[playerID, [command...]]`    |

#### `LMSResponse`

サーバーからのレスポンスを表す構造体です。`Result` フィールドに JSON の生データが格納されます。

#### `PlayersResult`

`players` コマンドのレスポンス結果を表す構造体です。

| フィールド | 型     | 説明           |
|------------|--------|----------------|
| `Count`    | `int`  | プレイヤー数   |
| `Players`  | 配列   | プレイヤーの詳細 |

各プレイヤーは以下のフィールドを持ちます。

| フィールド | 型       | 説明                    |
|------------|----------|-------------------------|
| `PlayerID` | `string` | プレイヤーの MAC アドレス |
| `Name`     | `string` | プレイヤーの表示名       |
| `IP`       | `string` | プレイヤーの IP アドレス  |

---

### 関数

#### `NewRequest(playerID string, command ...string) LMSRequest`

指定したプレイヤー ID とコマンドを持つリクエストを生成します。ID には UUID が自動付与されます。

```go
req := lms.NewRequest("aa:bb:cc:dd:ee:ff", "pause", "1")
```

#### `NewPlayersRequest() LMSRequest`

接続中のプレイヤー一覧を取得するリクエストを生成します。`NewRequest` のショートカットです。

```go
req := lms.NewPlayersRequest()
```

#### `SendRequest[T any](server string, req LMSRequest, logger *slog.Logger) (T, error)`

LMS にリクエストを送信し、レスポンスの `result` フィールドを型 `T` にデコードして返します。

| 引数     | 説明                                                      |
|----------|-----------------------------------------------------------|
| `server` | LMS サーバーのベース URL（例: `http://192.168.1.10:9000`）|
| `req`    | 送信するリクエスト                                        |
| `logger` | デバッグログ用の `*slog.Logger`。不要な場合は `nil`       |

レスポンスの ID がリクエストの ID と一致しない場合はエラーを返します。

## 依存ライブラリ

| パッケージ | 用途 |
|---|---|
| [`github.com/google/uuid`](https://github.com/google/uuid) | リクエスト ID の UUID 生成 |

## ライセンス

MIT
