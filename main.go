/**
 * Copyright 2015-2016, Wothing Co., Ltd.
 * All rights reserved.
 *
 * Created by Elvizlai on 2016/07/18 09:48
 */

package main

import (
	"flag"
	"log"
	"log/syslog"
	"os"
	"os/signal"
	"syscall"

	"github.com/elvizlai/xtunnel/tunnel"
)

var logTo string

var mode string
var laddr string
var raddr string
var cryptoMethod string
var secret string

func init() {
	flag.StringVar(&logTo, "logto", "stdout", "stdout or syslog")
	flag.StringVar(&mode, "mode", "", "run mode: server, client")
	flag.StringVar(&laddr, "listen", "127.0.0.1:9000", "xtunnel local listen")
	flag.StringVar(&raddr, "remote", "127.0.0.1:9001", "xtunnel remote backend")
	flag.StringVar(&cryptoMethod, "crypto", "rc4-md5", "encryption method: blank, rc4, rc4-md5, aes256cfb, chacha20, salsa20")
	flag.StringVar(&secret, "secret", "xtunnel", "password used to encrypt data")
	flag.Parse()
}

type svr interface {
	Run()
}

func wait() {
	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan)
	for sig := range sigChan {
		if sig == syscall.SIGINT || sig == syscall.SIGTERM {
			log.Printf("terminated by signal %v\n", sig)
			return
		}
		log.Printf("received signal: %v, ignore\n", sig)
	}
}

func main() {
	if logTo == "syslog" {
		w, err := syslog.New(syslog.LOG_INFO, "xtunnel")
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(w)
	} else {
		log.SetOutput(os.Stdout)
	}

	var app svr

	switch mode {
	case "server":
		app = tunnel.NewTunnel(laddr, raddr, false, cryptoMethod, secret, 4096)
	case "client":
		app = tunnel.NewTunnel(laddr, raddr, true, cryptoMethod, secret, 4096)
	default:
		log.Fatalf("no such '%s' mode", mode)
	}

	app.Run()
	wait()
}
