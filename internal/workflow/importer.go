package workflow

import (
	"emergency-claim-code/internal/model"
	"fmt"
	"strings"
)

type ImportRow struct{ BatchID, Applicant, Quantity, Note string }

func ParseRow(line string) (ImportRow, error) {
	parts := strings.Split(line, "|")
	if len(parts) != 4 {
		return ImportRow{}, fmt.Errorf("row needs four fields")
	}
	return ImportRow{BatchID: strings.TrimSpace(parts[0]), Applicant: strings.TrimSpace(parts[1]), Quantity: strings.TrimSpace(parts[2]), Note: strings.TrimSpace(parts[3])}, nil
}

func ParseRows(input string) ([]ImportRow, []error) {
	rows := make([]ImportRow, 0)
	errs := make([]error, 0)
	for _, line := range strings.Split(input, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		row, err := ParseRow(line)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		rows = append(rows, row)
	}
	return rows, errs
}

func ToRecord(row ImportRow, id string) (model.Record, error) {
	var quantity int
	if _, err := fmt.Sscanf(row.Quantity, "%d", &quantity); err != nil {
		return model.Record{}, err
	}
	return model.Record{ID: id, BatchID: row.BatchID, Applicant: row.Applicant, Quantity: quantity, Note: row.Note}, nil
}
