package core

import (
	"context"
	"ego/src/boot/global"
	"ego/src/boot/initialize"
	"ego/src/utils/backrun"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func RunServer() {
	Router := initialize.Routers()
	sys := global.C_CONFIG.System
	s := &http.Server{
		Addr:           ":" + sys.Addr,
		Handler:        Router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	time.Sleep(10 * time.Microsecond)
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Printf(`
    欢迎使用 `+sys.Webname+`
    当前版本:`+sys.Webversion+`
    访问地址:http://localhost%s/
`, ":"+sys.Addr)

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("接收到关闭信号，正在处理剩余请求并关闭服务...")

	// ==================== 1. 优雅关闭 HTTP 服务 ====================
	// 给 HTTP 服务 10 秒的独立关闭时间
	httpCtx, httpCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer httpCancel()

	if err := s.Shutdown(httpCtx); err != nil {
		log.Fatalf("服务器强制关闭异常: %s\n", err)
	}

	// ==================== 2. 优雅关闭定时任务 ====================
	if backrun.CronInstance != nil {
		fmt.Println("正在停止后台定时任务...")

		// 停止接收新任务，并拿到 cron 内部的结束通知 context
		cronCtx := backrun.CronInstance.Stop()

		// 给定时任务单独分配 5 秒的独立关闭时间（可根据实际任务长短调整）
		cronTimeoutCtx, cronCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cronCancel()

		// 优雅合并：要么定时任务自己结束，要么单独给它的 5 秒超时到了
		select {
		case <-cronCtx.Done():
			fmt.Println("后台定时任务已全部安全停止")
		case <-cronTimeoutCtx.Done():
			log.Println("警告：定时任务关闭超时，强制退出")
		}
	}

	log.Println("服务器已安全退出")
}
