package main

import (
	"bufio"
	"fmt"
	"github.com/qPyth/coingecko-api-task/config"
	cg "github.com/qPyth/coingecko-api-task/internal/coingecko"
	"io"
	"log"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	pathToLog = "logs.txt"
	host      = "https://api.coingecko.com/api/v3/coins/markets"
)

func main() {
	//load cfg
	cfg := config.MustLoad()

	logFile, err := os.OpenFile(pathToLog, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("log file open error: %v", err.Error())
	}

	logger := loggerInit(logFile)

	u, err := url.Parse(host)
	if err != nil {
		logger.Error("url parse error", "err", err)
		return
	}
	client := cg.NewCGClient(u.Host, u.Path, cfg.APIKey)
	service := cg.NewCoinService(client, logger)

	//async update the coin currency in cache
	go func() {
		for {
			time.Sleep(cfg.UpdateTime * time.Second)
			service.AsyncUpdate()
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("enter crypto name: ")
	for scanner.Scan() {
		str := scanner.Text()
		cur, err := service.CoinCurrency(strNormalize(str))
		switch err {
		case nil:
			fmt.Printf("%s = %.03f\n", cur.Name, cur.CurrentPrice)
		default:
			fmt.Println(err)
		}
		fmt.Printf("enter crypto name: ")
	}
}

func strNormalize(str string) string {
	return strings.TrimSpace(strings.ToLower(str))
}

func loggerInit(w io.Writer) *slog.Logger {
	return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}
