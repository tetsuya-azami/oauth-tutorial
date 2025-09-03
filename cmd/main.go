package main

import (
	"log"
	"net/http"
	"oauth-tutorial/internal/infrastructure"
	pAuthorize "oauth-tutorial/internal/presentation/authorize"
	pDecision "oauth-tutorial/internal/presentation/decision"
	pToken "oauth-tutorial/internal/presentation/token"
	"oauth-tutorial/internal/session"
	uAuthorize "oauth-tutorial/internal/usecase/authorize"
	uDecision "oauth-tutorial/internal/usecase/decision"
	"oauth-tutorial/pkg/mycrypto"
	"oauth-tutorial/pkg/mylogger"
)

func main() {
	// ロガー構築
	logger := mylogger.NewLogger()
	logger.Info("start server")

	// 認可リクエストのためのコンポーネントを初期化
	cr := infrastructure.NewClientRepository()
	sig := session.NewSessionIDGenerator()
	ss := infrastructure.NewSessionStorage()
	acf := uAuthorize.NewAuthorizationCodeFlow(logger, cr, sig, ss)

	// 認可コード発行のためのコンポーネントを初期化
	rg := &mycrypto.RandomGenerator{}
	ur := infrastructure.NewUserRepository()
	ar := infrastructure.NewAuthCodeRepository()
	pac := uDecision.NewPublishAuthorizationCodeUseCase(logger, rg, ss, ur, ar)

	// ハンドラーの登録
	http.Handle("GET /authorize", pAuthorize.NewAuthorizeHandler(logger, acf))
	http.Handle("POST /decision", pDecision.NewDecisionHandler(logger, pac))
	http.Handle("POST /token", pToken.NewTokenHandler(logger))

	// サーバーの起動
	logger.Info("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
