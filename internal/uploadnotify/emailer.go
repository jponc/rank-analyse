package uploadnotify

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jponc/rank-analyse/internal/repository/dbrepository"
	"github.com/jponc/rank-analyse/pkg/ses"
)

const sender = "ponce.julianalfonso@gmail.com"

type Emailer struct {
	dbrepository *dbrepository.Repository
	sesClient    *ses.Client
	apiBaseURL   string
}

func NewEmailer(dbrepository *dbrepository.Repository, sesClient *ses.Client, apiBaseURL string) *Emailer {
	e := &Emailer{
		dbrepository: dbrepository,
		sesClient:    sesClient,
		apiBaseURL:   apiBaseURL,
	}

	return e
}

func (e *Emailer) Send(ctx context.Context, crawlID uuid.UUID) error {
	crawl, err := e.dbrepository.GetCrawl(ctx, crawlID)
	if err != nil {
		return fmt.Errorf("failed to get crawl: %v", err)
	}

	toAddresses := []string{crawl.Email}
	ccAddress := []string{}
	subject := fmt.Sprintf("Crawl finished for keyword [%s]", crawl.Keyword)

	link := fmt.Sprintf("%s/crawls/%s", e.apiBaseURL, crawlID)

	htmlBody := fmt.Sprintf(`
		<h1>Finished crawling for keyword "%s"</h1>
		Download it here: <a href="%s">%s</a>
	`, crawl.Keyword, link, link)

	textBody := fmt.Sprintf(`
		Finished crawling for keyword %s, download it here: %s
	`, crawl.Keyword, link)

	err = e.sesClient.SendEmail(
		ctx,
		toAddresses,
		ccAddress,
		sender,
		subject,
		htmlBody,
		textBody,
	)

	if err != nil {
		return fmt.Errorf("failed to send crawl result email: %v", err)
	}

	return nil
}
