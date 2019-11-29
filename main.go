package sm

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	filename := flag.String("conf", "./conf.yaml", "配置文件")
	flag.Parse()

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	server := NewServer()
	server.InitConfig(*filename)
	server.Serve()
}
