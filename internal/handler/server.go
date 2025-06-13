package handler

import (
	"log"
	"net/http"
	"oauth-tutorial/internal/dto"
)

func StartServer() {
	http.HandleFunc("/authorize", authorizeHandler)
	http.HandleFunc("/token", tokenHandler)
	log.Println("[INFO] OAuth2.0 Authorization Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func authorizeHandler(w http.ResponseWriter, r *http.Request) {
	req, err := dto.NewAuthorizeRequest(r.URL.Query())
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// 本来はここでユーザー認証・同意画面を表示する
	// 今回は簡易的に認可コードを発行
	code := "dummy-auth-code" // 実際はランダムなコードを生成・保存する

	// 認可コードをリダイレクトURIに付与してリダイレクト
	http.Redirect(w, r, req.RedirectURI()+"?code="+code+"&state="+req.State(), http.StatusFound)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: トークン発行処理の実装
	w.Write([]byte("token endpoint"))
}
