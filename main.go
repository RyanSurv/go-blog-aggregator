package main

import (
	"github.com/ryansurv/go-blog-aggregator/internal/config"
)

func main() {
	cfg := config.Read()
	cfg.SetUser("ryan")
	cfg.PrettyPrint()
}