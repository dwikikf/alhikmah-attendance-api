package repository

import (
	"database/sql"
	"encoding/json"

	"alhikmah-attendance-api/core/domain"
)

type reportCachePostgres struct {
	db *sql.DB
}

func NewReportCacheRepository(db *sql.DB) domain.ReportCacheRepository {
	return &reportCachePostgres{
		db: db,
	}
}

func (r *reportCachePostgres) Get(reportType, classID, periodStart, periodEnd string) (json.RawMessage, error) {
	var data json.RawMessage
	query := `
		SELECT report_data 
		FROM reports 
		WHERE report_type = $1 AND class_id = $2 AND period_start = $3 AND period_end = $4 
		  AND generated_at > NOW() - INTERVAL '1 hour'
	`
	err := r.db.QueryRow(query, reportType, classID, periodStart, periodEnd).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // cache miss
		}
		return nil, err
	}
	return data, nil
}

func (r *reportCachePostgres) Set(reportType, classID, periodStart, periodEnd, generatedBy string, data json.RawMessage) error {
	query := `
		INSERT INTO reports (report_type, class_id, period_start, period_end, generated_by, report_data) 
		VALUES ($1::report_type_enum, $2, $3, $4, $5, $6) 
		ON CONFLICT (report_type, class_id, period_start, period_end) 
		DO UPDATE SET 
			report_data = EXCLUDED.report_data, 
			generated_by = EXCLUDED.generated_by, 
			generated_at = NOW()
	`
	var genBy interface{} = generatedBy
	if generatedBy == "" {
		// fallback to any valid user UUID because generated_by is NOT NULL
		var fallbackUser string
		err := r.db.QueryRow("SELECT id FROM users LIMIT 1").Scan(&fallbackUser)
		if err == nil {
			genBy = fallbackUser
		} else {
			// fallback if somehow users table is empty, though this will likely fail FK constraint
			genBy = "f88dad22-46d3-4c7a-80eb-7a9272254d23"
		}
	}

	_, err := r.db.Exec(query, reportType, classID, periodStart, periodEnd, genBy, data)
	return err
}

func (r *reportCachePostgres) Delete(reportType, classID, periodStart, periodEnd string) error {
	query := `
		DELETE FROM reports 
		WHERE report_type = $1 AND class_id = $2 AND period_start = $3 AND period_end = $4
	`
	_, err := r.db.Exec(query, reportType, classID, periodStart, periodEnd)
	return err
}

