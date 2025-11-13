package sheets

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type (
	Service struct {
		*sheets.Service
		logger    *zap.SugaredLogger
		tableData TableData
	}
	Menty struct {
		Name     string
		Telegram string
		Sprint   int
	}

	TableData struct {
		TableName          string
		MentyNameColumn    string
		SprintNumberColumn string
		SpreadsheetId      string
	}
)

func NewSheetsService(ctx context.Context, httpClient *http.Client, logger *zap.SugaredLogger) (*Service, error) {
	sheetsService, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("init sheets service: %w", err)
	}

	tableData := TableData{
		TableName:          os.Getenv("TABLE_NAME"),
		MentyNameColumn:    os.Getenv("MENTY_NAME_COLUMN"),
		SprintNumberColumn: os.Getenv("SPRINT_NUMBER_COLUMN"),
		SpreadsheetId:      os.Getenv("SPREADSHEET_ID"),
	}

	return &Service{
		Service:   sheetsService,
		logger:    logger,
		tableData: tableData,
	}, nil
}

func (s *Service) GetMentyInformation(ctx context.Context) ([]Menty, error) {
	rangeColumns := fmt.Sprintf("%s!%s:%s", s.tableData.TableName, s.tableData.MentyNameColumn, s.tableData.SprintNumberColumn)
	if s.logger != nil {
		s.logger.Infow("fetching menty information", "spreadsheet_id", s.tableData.SpreadsheetId, "range", rangeColumns)
	}
	resp, err := s.Spreadsheets.Values.Get(s.tableData.SpreadsheetId, rangeColumns).
		Context(ctx).
		Do()
	if err != nil {
		if s.logger != nil {
			s.logger.Errorw("failed to fetch menty information", "spreadsheet_id", s.tableData.SpreadsheetId, "range", rangeColumns, "error", err)
		}
		return nil, fmt.Errorf("get menty information: %w", err)
	}

	if len(resp.Values) <= 1 {
		if s.logger != nil {
			s.logger.Infow("menty information is empty", "spreadsheet_id", s.tableData.SpreadsheetId, "range", rangeColumns)
		}
		return nil, nil
	}

	menties := make([]Menty, 0, len(resp.Values)-1)

	for _, row := range resp.Values[1:] {
		if len(row) < 4 {
			if s.logger != nil {
				s.logger.Warnw("skipping incomplete row", "row", row)
			}
			continue
		}

		name := fmt.Sprint(row[0])
		telegram := fmt.Sprint(row[1])
		sprintStr := strings.TrimSpace(fmt.Sprint(row[3]))
		sprintStr = strings.TrimPrefix(sprintStr, "Спринт")
		sprintStr = strings.TrimSpace(sprintStr)

		sprint, err := strconv.Atoi(sprintStr)
		if err != nil {
			if s.logger != nil {
				s.logger.Errorw("failed to parse sprint value", "row", row, "value", sprintStr, "error", err)
			}
			return nil, fmt.Errorf("parse sprint value %q: %w", sprintStr, err)
		}

		menties = append(menties, Menty{
			Name:     name,
			Telegram: telegram,
			Sprint:   sprint,
		})
	}

	if s.logger != nil {
		s.logger.Infow("menty information fetched", "count", len(menties))
	}
	return menties, nil
}
