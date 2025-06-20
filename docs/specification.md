# OAuth 2.0 認可サーバー仕様書

## 1. 概要

本仕様書は、OAuth 2.0 の認可コードフローを学習目的で実装する簡易的な認可サーバーの設計仕様である。基本的な認可エンドポイントおよびトークンエンドポイントを備える。

## 2. 機能要件

### 2.1 認可エンドポイント `/authorize`
- クライアントからの認可リクエストを受け付ける。
- ユーザーに対してログインと同意画面を表示する。
- 認可コードを生成し、リダイレクト URI に付与してリダイレクトする。

### 2.2 トークンエンドポイント `/token`
- 認可コードを受け取り、アクセストークンを発行する。
- クライアント認証は不要（パブリッククライアントのみ対応）。

### 2.3 ユーザー認証
- 単一ユーザーの固定アカウント（例: user/password）によるログイン処理。
- セッションを利用したログイン状態の管理。

### 2.4 クライアント管理
- クライアント情報（client_id, client_name, redirect_uri）をコード内にハードコード。
- 単一クライアントのみ対応。

## 3. 非機能要件

### 3.1 セキュリティ
- PKCE（Proof Key for Code Exchange）未対応。
- CSRF対策未実装。
- トークンやコードのランダム性は学習用途として最小限。

### 3.2 可用性・保守性
- リクエストログなどの監査機能なし。
- トークンやコードの期限管理や削除処理なし。
- リフレッシュトークン未対応。

### 3.3 拡張性
- Authorization Code Flow のみ対応。
- 複数クライアントやスコープ対応なし。
- クライアント秘密鍵を使った機密クライアント対応なし。

## 4. インターフェース仕様

### 4.1 認可エンドポイント `GET /authorize`
**パラメータ**:

| No. | フィールド名     | フィールドの説明               | フィールドの型 | フィールドの制約         | 備考                             |
|-----|------------------|-------------------------------|----------------|---------------------------|----------------------------------|
| 1   | response_type    | レスポンスタイプの指定        | 文字列         | 必須、固定値 `code`       | 認可コードフローのみ対応         |
| 2   | client_id        | クライアントの識別子          | 文字列         | 必須                      | ハードコードされたクライアントID |
| 3   | redirect_uri     | 認可後のリダイレクト先 URI     | 文字列（URI）  | 必須                      | クライアントに事前登録されている |
| 4   | state            | CSRF 対策用トークン           | 文字列         | 任意

**レスポンス**:
```json
{
	"access_token": "xxxxxxxxxxxxx",
	"token_type": "bearer",
	"expires_in": 3600
}
```

### 4.2 トークンエンドポイント `POST /token`
- **パラメータ**（x-www-form-urlencoded）:

| No. | フィールド名     | フィールドの説明                    | フィールドの型 | フィールドの制約                     | 備考                                     |
|-----|------------------|-------------------------------------|----------------|--------------------------------------|------------------------------------------|
| 1   | grant_type       | グラントタイプの指定               | 文字列         | 必須、固定値 `authorization_code`    | 認可コードフローのみ対応                 |
| 2   | code             | 認可コード                         | 文字列         | 必須                                 | 認可エンドポイントで発行された値         |
| 3   | redirect_uri     | リダイレクト URI                   | 文字列（URI）  | 必須                                 | 認可リクエスト時と同一である必要がある   |
| 4   | client_id        | クライアントの識別子               | 文字列         | 必須                                 | 照合のみ、クライアント認証は不要         |

**レスポンス**（JSON形式）:
  ```json
  {
    "access_token": "xxxxxxxxxxxxx",
    "token_type": "bearer",
    "expires_in": 3600
  }
```
