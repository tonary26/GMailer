package worker

import (
	"bytes"
	"context"
	"gmailer/internal/database"
	"gmailer/internal/mailer"
	"gmailer/internal/models"
	"log"
	"text/template"
	"time"

	"golang.org/x/time/rate"
)

type Worker struct {
	store  *database.MailingStore
	sender mailer.Sender
	limit  *rate.Limiter
}

func NewWorker(store *database.MailingStore, sender mailer.Sender) *Worker {
	return &Worker{
		store:  store,
		sender: sender,
		limit:  rate.NewLimiter(rate.Every(time.Second), 1),
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Воркер остановлен")
			return
		case <-ticker.C:
			w.proccesBatch(ctx)
		}
	}
}

func (w *Worker) proccesBatch(ctx context.Context) {
	messages, err := w.store.FetchPending(ctx, 10)
	if err != nil {
		log.Printf("Ошибка получения письма: %v", err)
		return
	}

	for _, msg := range messages {
		if err := w.limit.Wait(ctx); err != nil {
			return
		}
		w.sendOne(ctx, msg)
	}
}

func (w *Worker) sendOne(ctx context.Context, msg models.PendingMessage) {
	body, err := w.renderTemplate(msg.BodyTemplate, msg)
	if err != nil {
		if storeErr := w.store.MarkFailed(ctx, msg.ID, 3, time.Minute); storeErr != nil {
			log.Printf("Ошибка обновления письма %d: %v", msg.ID, storeErr)
		}
		w.finishMailing(ctx, msg.MailingID)
		return
	}

	err = w.sender.Send(&mailer.Message{
		To:      msg.Email,
		Subject: msg.Subject,
		Body:    body,
	})
	if err != nil {
		log.Printf("Ошибка отправки на %s: %v", msg.Email, err)
		if storeErr := w.store.MarkFailed(ctx, msg.ID, 3, time.Minute); storeErr != nil {
			log.Printf("Ошибка обновления письма %d: %v", msg.ID, storeErr)
		}
		w.finishMailing(ctx, msg.MailingID)
		return
	}

	if err := w.store.MarkSent(ctx, msg.ID); err != nil {
		log.Printf("Ошибка обновления письма %d: %v", msg.ID, err)
		return
	}
	w.finishMailing(ctx, msg.MailingID)
}

func (w *Worker) finishMailing(ctx context.Context, mailingID int) {
	if err := w.store.MarkDoneIfComplete(ctx, mailingID); err != nil {
		log.Printf("Ошибка завершения рассылки %d: %v", mailingID, err)
	}
}

func (w *Worker) renderTemplate(tmplStr string, data models.PendingMessage) (string, error) {
	tmpl, err := template.New("email").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer
	if err := tmpl.Execute(&buff, data); err != nil {
		return "", err
	}

	return buff.String(), nil
}
