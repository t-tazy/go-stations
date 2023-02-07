package server

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

// wrap http.Server
type Server struct {
	srv *http.Server
}

func NewServer(addr string, mux http.Handler) *Server {
	return &Server{&http.Server{
		Addr:    addr,
		Handler: mux,
	}}
}

func (s Server) Run(ctx context.Context) error {
	// シグナルをハンドリング
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// error伝達用
	errCh := make(chan error)

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("failed to terminate server")
			errCh <- err
		}
		errCh <- nil
	}()

	// コンテキストを通じて処理の中断を検知したとき
	// ShutdownメソッドでHTTPサーバーの機能を終了する
	select {
	case <-ctx.Done():
		if err := s.srv.Shutdown(context.Background()); err != nil {
			log.Printf("failed to shutdown: %+v", err)
			return err
		}
		return <-errCh
	case err := <-errCh: // serverの起動に失敗
		return err
	}
}
