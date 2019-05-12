package frame

import (
	"github.com/google/uuid"
	"testing"
)

func Test_NewLine(t *testing.T) {
	events := make(chan ScreenEvent, 100)
	line := NewLine(22, events)

	expectedRow := 22
	actualRow := line.row
	if expectedRow != actualRow {
		t.Errorf("NewLine(): expected row of %d, but is at row %d", expectedRow, actualRow)
	}
}

func Test_Line_String(t *testing.T) {
	tables := []struct {
		row         int
		message     string
		expectedStr string
	}{
		{1, "This test will do well..", "<Line row:1 buff:24 id:f47ac10b-58cc-0372-8567-0e02b2c3d479>"},
		{2, "...won't it?", "<Line row:2 buff:12 id:f47ac10b-58cc-0372-8567-0e02b2c3d479>"},
	}

	for _, table := range tables {
		events := make(chan ScreenEvent, 100)
		line := NewLine(table.row, events)
		u, _ := uuid.Parse("f47ac10b-58cc-0372-8567-0e02b2c3d479")
		line.id = u
		line.buffer = []byte(table.message)

		expectedResult := table.expectedStr
		actualResult := line.String()

		if expectedResult != actualResult {
			t.Errorf("Line.String(): expected '%s', but got '%s'", expectedResult, actualResult)
		}
	}
}

func Test_Line_move(t *testing.T) {
	tables := []struct {
		row      int
		moveRows int
	}{
		{1, 12},
		{2, 22},
		{55, -22},
	}

	for _, table := range tables {
		events := make(chan ScreenEvent, 100)
		line := NewLine(table.row, events)
		line.move(table.moveRows)

		expectedResult := table.row + table.moveRows
		actualResult := line.row

		if expectedResult != actualResult {
			t.Errorf("Line.Move(): expected row '%d', but got row '%d'", expectedResult, actualResult)
		}

		// if !line.stale {
		// 	t.Errorf("Line.Move(): expected line %d to be stale, but was not", line.row)
		// }
	}

}
