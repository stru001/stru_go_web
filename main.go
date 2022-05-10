package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"os/signal"
	"stru_web/routes"
	"syscall"
	"time"

	"stru_web/dao/mysql"
	"stru_web/dao/redis"
	"stru_web/logger"
	"stru_web/settings"
)

func main()  {
	//0. 输入命令
	var f string
	flag.StringVar(&f, "f", "./config/config.yml", "配置文件路径")
	flag.Parse()
	//1. 加载配置文件
	if err := settings.Init(f); err != nil {
		return
	}
	//2. 初始化日志
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		return
	}
	zap.L().Sync()
	//3. 初始化Mysql
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		return
	}
	mysql.Close()
	//4. 初始化Redis
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		return
	}
	defer redis.Close()
	//5. 注册路由
	r := routes.Setup()
	//6. 启动服务
	//r.Run()
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d",settings.Conf.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error("listen: %s\n", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)  // 此处不会阻塞
	<-quit  // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}