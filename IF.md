# IF設計書（API仕様）

## 認可エンドポイント
- URL: `/authorize`
- メソッド: GET
- パラメータ:
  - response_type: code
  - client_id: クライアントID
  - redirect_uri: リダイレクトURI
  - scope: スコープ
  - state: CSRF対策用ランダム文字列
- レスポンス: 認可コードをリダイレクトURIに付与してリダイレクト

## トークンエンドポイント
- URL: `/token`
- メソッド: POST
- パラメータ:
  - grant_type: authorization_code
  - code: 認可コード
  - redirect_uri: リダイレクトURI
  - client_id: クライアントID
  - client_secret: クライアントシークレット
- レスポンス:
  - access_token: アクセストークン
  - token_type: bearer
  - expires_in: 有効期限（秒）
  - refresh_token: リフレッシュトークン
