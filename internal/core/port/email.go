package port

import (
	"context"
	"stori-challenge/internal/core/domain"
)

type EmailClient interface {
	SendEmailReport(ctx context.Context, reportDocument domain.ReportDocument)
}
