package jotframe

import (
	"testing"
	"sync"
)

func Test_NewLine(t *testing.T) {

	line := NewLine(22, &sync.WaitGroup{})

	expectedRow := 22
	actualRow := line.row
	if expectedRow!= actualRow {
		t.Errorf("NewLine(): expected row of %d, but is at row %d", expectedRow, actualRow)
	}
}

func Test_Line_String(t *testing.T) {
	tables := []struct {
		row int
		message string
		expectedStr string
	}{
		{1, "This test will do well..", "<Line row:1 closed:false bufferLen:24>"},
		{2, "...won't it?", "<Line row:2 closed:false bufferLen:12>"},
	}

	for _, table := range tables {
		line := NewLine(table.row, &sync.WaitGroup{})
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
		row int
		moveRows int
	}{
		{1, 12},
		{2, 22},
		{55, -22},
	}

	for _, table := range tables {
		line := NewLine(table.row, &sync.WaitGroup{})
		line.move(table.moveRows)

		expectedResult := table.row + table.moveRows
		actualResult := line.row

		if expectedResult != actualResult {
			t.Errorf("Line.move(): expected row '%d', but got row '%d'", expectedResult, actualResult)
		}

		if !line.stale{
			t.Errorf("Line.move(): expected line %d to be stale, but was not", line.row)
		}
	}

}
