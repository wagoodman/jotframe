package frame

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func suppressOutput(f func()) {
	originalStdOut := os.Stdout
	var err error
	os.Stdout, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	f()
	os.Stdout = originalStdOut
}

type TestEventHandler struct {
	t      *testing.T
	events []*ScreenEvent
}

func NewTestEventHandler(t *testing.T) *TestEventHandler {
	return &TestEventHandler{
		t: t,
		events: make([]*ScreenEvent, 0),
	}
}

func (handler *TestEventHandler) onEvent(event *ScreenEvent) {
	handler.events = append(handler.events, event)
}

type drawTestParams struct {
	rows           int
	hasHeader      bool
	hasFooter      bool
	startRow       int
	terminalHeight int
	events         []ScreenEvent
	errors         []string
}


func validateEvents(t *testing.T, test string, table drawTestParams, errs []error, frame *logicalFrame, handler *TestEventHandler) {

	if len(frame.activeLines) != table.rows {
		t.Errorf("[case=%s] expected %d rows, got %d", test, table.rows, len(frame.activeLines))
	}

	if len(table.errors) != len(errs) {
		t.Errorf("[case=%s] expected %d errors, got %d", test, len(table.errors), len(errs))
	} else {
		for idx, err := range errs {
			if err.Error() != table.errors[idx] {
				t.Errorf("[case=%s] unexpected error: %v", test, err)
			}
		}
	}

	if len(table.events) != len(handler.events) {
		t.Errorf("[case=%s] expected %d events, got %d", test, len(table.events), len(handler.events))
	} else {
		for idx, event := range table.events {

			if bytes.Compare(event.value, handler.events[idx].value) != 0 {
				t.Errorf("[case=%s] event=%d: expected value='%v', got '%v'", test, idx, string(event.value), string(handler.events[idx].value))
			}

			if event.row != handler.events[idx].row {
				t.Errorf("[case=%s] event=%d: expected row='%v', got '%v'", test, idx, event.row, handler.events[idx].row)
			}
		}
	}

	if t.Failed() {
		t.Logf("[case=%s] actual events", test)
		for idx, event := range handler.events {
			t.Log(fmt.Sprintf("   event=%d: row=%d value='%s'", idx, event.row, string(event.value)))
		}
		t.Logf("[case=%s] actual errors", test)
		for idx, err := range errs {
			t.Log(fmt.Sprintf("   error=%d: %v", idx, err.Error()))
		}
		t.Fatal("Stopping test")
	}
}


