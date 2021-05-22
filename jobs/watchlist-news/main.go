package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strings"
	"time"
	"tradingbot/lib/newsapi"
	"tradingbot/lib/traderdb"
	"tradingbot/lib/utils"

	"github.com/jackc/pgx/v4/pgxpool"
)

type EmailParams struct {
	GeneralNewsByCompany map[string][]newsapi.Article
}

const userID = 1

func main() {
	recipientsFlag := flag.String("recipients", "", "Email addresses delimited by commas")
	flag.Parse()

	recipients := strings.TrimSpace(*recipientsFlag)
	if recipients == "" {
		fmt.Println("At least one recipient must be specified")
		os.Exit(1)
	}

	databaseUrl := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseUrl == "" {
		fmt.Println("DATABASE_URL environment variable is required")
		os.Exit(1)
	}
	newsAPIKey := strings.TrimSpace(os.Getenv("NEWS_API_KEY"))
	if newsAPIKey == "" {
		fmt.Println("NEWS_API_KEY environment variable is required")
		os.Exit(1)
	}
	alpacaAPIKey := strings.TrimSpace(os.Getenv("ALPACA_API_KEY"))
	if alpacaAPIKey == "" {
		fmt.Println("ALPACA_API_KEY environment variable is required")
		os.Exit(1)
	}
	email := strings.TrimSpace(os.Getenv("GMAIL_USERNAME"))
	if email == "" {
		fmt.Println("GMAIL_USERNAME environment variable is required")
		os.Exit(1)
	}
	password := strings.TrimSpace(os.Getenv("GMAIL_PASSWORD"))
	if password == "" {
		fmt.Println("GMAIL_PASSWORD environment variable is required")
		os.Exit(1)
	}

	emailTemplate, err := template.ParseFiles("email-template.html")
	if err != nil {
		fmt.Println("Failed to parse email template:", err)
		os.Exit(1)
	}

	now := time.Now()
	httpClient := utils.NewHttpClient()
	newsClient := newsapi.NewClient(httpClient, newsAPIKey)

	pool, err := pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		fmt.Println("Failed to connect to DB:", err)
		os.Exit(1)
	}

	emailParams, err := getEmailParams(pool, newsClient, now)
	if err != nil {
		fmt.Println("Failed to fetch news items:", err)
		os.Exit(1)
	}

	var body bytes.Buffer
	subject := "Subject: Watchlist News " + now.Format("Jan 02") + "\n"
	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.WriteString(subject + headers)
	err = emailTemplate.Execute(&body, emailParams)
	if err != nil {
		fmt.Println("Failed to populate email template:", err)
		os.Exit(1)
	}

	auth := smtp.PlainAuth("", email, password, "smtp.gmail.com")
	err = smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		email,
		strings.Split(recipients, ","),
		body.Bytes(),
	)
	if err != nil {
		fmt.Println("Failed to send email:", err)
		os.Exit(1)
	}
}

func getEmailParams(
	db *pgxpool.Pool,
	newsClient *newsapi.Client,
	now time.Time,
) (params EmailParams, err error) {
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		return params, err
	}

	stocks, err := traderdb.GetWatchlistStocksWithUserID(db, userID)
	if err != nil {
		return params, err
	}

	startTime := utils.GetLastWeekday(now)

	generalNewsByCompany := make(map[string][]newsapi.Article)
	for _, stock := range stocks {
		// Using domains because some news sources are missing from /sources
		// (e.g. seekingalpha.com, yahoo.finance.com)
		response, err := newsClient.GetAllHeadlinesWithSources(newsapi.AllArticlesQueryParams{
			Query:     utils.TrimCompanyName(stock.Company) + " OR " + stock.Symbol + " Stock",
			StartTime: startTime,
			Domains:   utils.NewsDomains,
			PageSize:  20,
		})
		if err != nil {
			return params, err
		}
		for index := range response.Articles {
			response.Articles[index].PublishedAt = response.Articles[index].PublishedAt.In(location)
		}
		generalNewsByCompany[stock.Symbol] = response.Articles
	}

	params = EmailParams{
		GeneralNewsByCompany: generalNewsByCompany,
	}
	return params, nil
}
