package main

import (
	"flag"
	"github.com/caddyserver/certmagic"
	"github.com/julienschmidt/httprouter"
	"github.com/t73liu/trading-bot/trader/news"
	"github.com/t73liu/trading-bot/trader/stock"
	"github.com/t73liu/trading-bot/trader/utils"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	prodFlag := flag.Bool("prod", false, "Run in production mode")
	httpsFlag := flag.Bool("https", false, "Run with HTTPS")
	emailFlag := flag.String(
		"email",
		"",
		"Email to receive expiration alerts for certificates (Optional)",
	)
	domainsFlag := flag.String(
		"domains",
		"",
		"Comma-delimited domains (Required for HTTPS)",
	)
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	logger.Println("Initializing ...")

	client := utils.NewHttpClient()

	handler := initApp(logger, client)

	if *httpsFlag {
		// https://github.com/caddyserver/certmagic#requirements
		logger.Println("Starting service with HTTPS")
		domainNames := strings.Split(strings.TrimSpace(*domainsFlag), ",")
		if len(domainNames) == 0 {
			logger.Fatalln("domains are required for HTTPS")
		}
		certmagic.DefaultACME.Agreed = true
		email := strings.TrimSpace(*emailFlag)
		if email != "" {
			certmagic.DefaultACME.Email = email
		}
		if !*prodFlag {
			certmagic.DefaultACME.CA = certmagic.LetsEncryptStagingCA
		}
		logger.Fatalln(certmagic.HTTPS(domainNames, handler))
	}

	logger.Printf("Starting service with HTTP at port %s\n", ":8080")
	server := utils.NewHttpServer(&handler)
	logger.Fatalln(server.ListenAndServe())
}

func initApp(logger *log.Logger, client *http.Client) http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/index.html")
	})
	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))

	newsClient := news.NewClient(client, os.Getenv("NEWS_API_KEY"))
	newsHandlers := news.NewHandlers(logger, newsClient)
	newsHandlers.AddRoutes(router)

	stockHandlers := stock.NewHandlers(logger)
	stockHandlers.AddRoutes(router)

	return router
}
