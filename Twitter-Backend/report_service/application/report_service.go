package application

import (
	"go.opentelemetry.io/otel/trace"
	"report_service/domain"
)

type ReportService struct {
	store  domain.ReportStore
	tracer trace.Tracer
}

func NewReportService(store domain.ReportStore, tracer trace.Tracer) *ReportService {
	return &ReportService{
		store: store,

		tracer: tracer,
	}
}
