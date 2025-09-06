# OAuth 2.0 認可サーバー仕様書

## 1. 概要

本仕様書は、OAuth 2.0 の認可コードフローを学習目的で実装する簡易的な認可サーバーの設計仕様である。基本的な認可エンドポイントおよびトークンエンドポイントを備える。

## 2. 機能要件

### 2.1 認可エンドポイント `/authorize`
- エンドユーザーからの認可リクエストを受け付ける。
- ユーザーに対してログインと同意画面を表示する。

### 2.2 認可コード発行エンドポイント `/decision`
- エンドユーザーからの認可コード発行リクエストを受け付ける。
- 認可コードを生成し、リダイレクト URI に付与してリダイレクトする。

### 2.3 トークンエンドポイント `/token`
- 認可コードを受け取り、アクセストークンを発行する。
- クライアント認証は必要(コンフィデンシャルクライアントのみ対応)

### 2.3 ユーザー認証
- 単一ユーザーの固定アカウント（例: user/password）によるログイン処理。
- セッションを利用したログイン状態の管理。

### 2.4 クライアント管理
- クライアント情報（client_id, client_name, redirect_uri）をインメモリで保管。
- 単一クライアントのみ対応。

## 3. 非機能要件

### 3.1 セキュリティ

### 3.2 可用性・保守性

### 3.3 拡張性
- Authorization Code Flow のみ対応。

## 4. インターフェース仕様
簡易実装なのでRFCと違う部分あり
### 4.1 認可エンドポイント `GET /authorize`
**クエリパラメータ**:

| No. | フィールド名     | フィールドの説明               | フィールドの型 | フィールドの制約         | 備考                             |
|-----|------------------|-------------------------------|----------------|---------------------------|----------------------------------|
| 1   | response_type    | レスポンスタイプの指定        | string | 必須、固定値 `code`       | 認可コードフローのみ対応         |
| 2   | client_id        | クライアントの識別子          | string | 必須                      | 事前登録されている想定 |
| 3   | redirect_uri     | 認可後のリダイレクト先 URI     | string(URL形式) | 必須                      | 事前登録されている想定 |
| 4   | scope           | 認可する操作の範囲    | read, write    | 必須 |  |
| 5   | state            | CSRF 対策用トークン           | string | 必須(PKCEサポート次第任意)

**成功レスポンス**:
```json
// 簡易実装なので画面ではなく、OKを返すのみとする。
{
	"message": "OK",
}
```

**エラーレスポンス** (JSON):
| フィールド | 型 | 説明 |
|---|---|---|
| error | string | 例: invalid_request, unsupported_response_type, server_error |
| error_description | string | エラー詳細メッセージ |
| state | string | 入力stateを返却 (存在する場合) |

HTTP ステータス:
- 400: invalid_request / unsupported_response_type / invalid_redirect_uri
- 500: server_error

### 4.2 認可コード発行エンドポイント `POST /decision`
**Content-Type**:
application/x-www-form-urlencoded

**ボディ**:
| No. | フィールド名     | フィールドの説明               | フィールドの型 | フィールドの制約         | 備考                             |
|-----|------------------|-------------------------------|----------------|---------------------------|----------------------------------|
| 1   | login_id    | ユーザーのログインID        | string | 必須 |  |
| 2   | password    | パスワード                 | string | 必須 |  |
| 3   | approved    | 認可フラグ                 | boolean | 必須   |  |

**ヘッダー**:
- Cookie: `session_id` (サーバが `/authorize` 応答時に付与)

**成功時**:
- HTTP 303 See Other
- Location: `<redirect_uri>?code=<authorization_code>&state=<state>`

**エラー時**:
- セッション不在/取得エラー: JSON で返却 (400)
  - `{ "message": "..." }`
- ユーザーが拒否: リダイレクト (303)
  - `<redirect_uri>?error=access_denied&error_description=...&state=...`
- 資格情報誤り: JSON で返却 (401)
  - `{ "message": "invalid login credentials" }`

### 4.3 トークンエンドポイント `POST /token`
**Content-Type**:
application/x-www-form-urlencoded

**ヘッダー**:
- Authorization: `Basic <base64(client_id:client_secret)>` （Basic認証のみサポート）
- 例: `Authorization: Basic Y2xpZW50SWQ6c2VjcmV0`

**ボディ**:
| No. | フィールド名     | フィールドの説明                    | フィールドの型 | フィールドの制約                     | 備考                                     |
|-----|------------------|-------------------------------------|----------------|--------------------------------------|------------------------------------------|
| 1   | grant_type       | グラントタイプの指定               | 文字列         | 必須、固定値 `authorization_code`    | 認可コードフローのみ対応                 |
| 2   | code             | 認可コード                         | 文字列         | 必須                                 | 認可エンドポイントで発行された値         |
| 3   | redirect_uri     | リダイレクト URI                   | 文字列（URI）  | 必須                                 | 認可リクエスト時と同一である必要がある   |

**レスポンス**（JSON形式）
  ```json
  {
    "access_token": "xxxxxxxxxxxxx",
    "refresh_token": "xxxxxxxxxxxxx",
    "token_type": "bearer",
    "expires_in": 3600
  }
```

**エラーレスポンス**

- 認証エラー（Basic 認証不備/不正）
  - HTTP ステータス: 401 Unauthorized
  - レスポンスヘッダー: `WWW-Authenticate: Basic realm="token", charset="UTF-8"`
  - ボディ例:
    ```json
    { "error": "invalid_client", "error_description": "client authentication failed" }
    ```

- バリデーション/業務エラー（JSON ボディ、WWW-Authenticate は付与しない）
  - 400 Bad Request: `invalid_request`（必須欠落/形式不正）
  - 400 Bad Request: `unsupported_grant_type`（grant_type が authorization_code 以外）
  - 400 Bad Request: `invalid_grant`（code 不正/期限切れ、redirect_uri 不一致）
  - 400 Bad Request: `unauthorized_client`（クライアントに許可されていない）
  - 500 Internal Server Error: `server_error`
  - ボディ例:
    ```json
    { "error": "invalid_grant", "error_description": "authorization code is invalid or expired" }
    ```
