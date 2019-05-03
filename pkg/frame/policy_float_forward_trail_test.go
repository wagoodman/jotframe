package frame

import (
	"fmt"
	"sort"
	"testing"
)

var floatForwardDrawTestCases = map[string]drawTestParams{
	"FloatForward_goCase": {3, false, false, 10, PolicyFloatForwardTrail, 40,
		[]ScreenEvent{
			// create the frame (pave a blank spot)
			{row: 10, value: []byte("")},
			{row: 11, value: []byte("")},
			{row: 12, value: []byte("")},
			// draw the first update
			{row: 10, value: []byte("LineIdx:0")},
			{row: 11, value: []byte("LineIdx:1")},
			{row: 12, value: []byte("LineIdx:2")},
		},
		[]string{},
	},
	"FloatForward_Header": {3, true, false, 10, PolicyFloatForwardTrail, 40,
		[]ScreenEvent{
			// create the frame (pave a blank spot)
			{row: 10, value: []byte("")},
			{row: 11, value: []byte("")},
			{row: 12, value: []byte("")},
			{row: 13, value: []byte("")},
			// draw the first update
			{row: 10, value: []byte("theHeader")},
			{row: 11, value: []byte("LineIdx:0")},
			{row: 12, value: []byte("LineIdx:1")},
			{row: 13, value: []byte("LineIdx:2")},
		},
		[]string{},
	},
	"FloatForward_Footer": {3, false, true, 10, PolicyFloatForwardTrail, 40,
		[]ScreenEvent{
			// create the frame (pave a blank spot)
			{row: 10, value: []byte("")},
			{row: 11, value: []byte("")},
			{row: 12, value: []byte("")},
			{row: 13, value: []byte("")},
			// draw the first update
			{row: 10, value: []byte("LineIdx:0")},
			{row: 11, value: []byte("LineIdx:1")},
			{row: 12, value: []byte("LineIdx:2")},
			{row: 13, value: []byte("theFooter")},
		},
		[]string{},
	},
	"FloatForward_HeaderFooter": {3, true, true, 10, PolicyFloatForwardTrail, 40,
		[]ScreenEvent{
			// create the frame (pave a blank spot)
			{row: 10, value: []byte("")},
			{row: 11, value: []byte("")},
			{row: 12, value: []byte("")},
			{row: 13, value: []byte("")},
			{row: 14, value: []byte("")},
			// draw the first update
			{row: 10, value: []byte("theHeader")},
			{row: 11, value: []byte("LineIdx:0")},
			{row: 12, value: []byte("LineIdx:1")},
			{row: 13, value: []byte("LineIdx:2")},
			{row: 14, value: []byte("theFooter")},
		},
		[]string{},
	},
	"FloatForward_TermHeightSmall_AtTop": {3, false, false, 1, PolicyFloatForwardTrail, 2,
		[]ScreenEvent{
			// create the frame (pave a blank spot)
			{row: 1, value: []byte("")},
			{row: 2, value: []byte("")},
			// draw the first update
			{row: 1, value: []byte("LineIdx:0")},
			{row: 2, value: []byte("LineIdx:1")},
		},
		[]string{
			"line is out of bounds (row=3)",
		},
	},
	"FloatForward_TermHeightSmall_AtTop_Header": {3, true, false, 1, PolicyFloatForwardTrail, 2,
		[]ScreenEvent{
			// create the frame (pave a blank spot)
			{row: 1, value: []byte("")},
			{row: 2, value: []byte("")},
			// draw the first update
			{row: 1, value: []byte("theHeader")},
			{row: 2, value: []byte("LineIdx:0")},
		},
		[]string{
			"line is out of bounds (row=3)",
			"line is out of bounds (row=4)",
		},
	},
	"FloatForward_TermHeightSmall_AtTop_Footer": {3, false, true, 1, PolicyFloatForwardTrail, 2,
		[]ScreenEvent{
			// create the frame (pave a blank spot)
			{row: 1, value: []byte("")},
			{row: 2, value: []byte("")},
			// draw the first update
			{row: 1, value: []byte("LineIdx:0")},
			{row: 2, value: []byte("LineIdx:1")},
		},
		[]string{
			"line is out of bounds (row=3)",
			"line is out of bounds (row=4)",
		},
	},
	"FloatForward_TermHeightSmall_AtTop_HeaderFooter": {3, true, true, 1, PolicyFloatForwardTrail, 2,
		[]ScreenEvent{
			// create the frame (pave a blank spot)
			{row: 1, value: []byte("")},
			{row: 2, value: []byte("")},
			// draw the first update
			{row: 1, value: []byte("theHeader")},
			{row: 2, value: []byte("LineIdx:0")},
		},
		[]string{
			"line is out of bounds (row=3)",
			"line is out of bounds (row=4)",
			"line is out of bounds (row=5)",
		},
	},
	"FloatForward_TermHeightSmall_AtBottom": {3, false, false, 49, PolicyFloatForwardTrail, 50,
		[]ScreenEvent{
			// create the frame (pave a blank spot)
			{row: 48, value: []byte("")},
			{row: 49, value: []byte("")},
			{row: 50, value: []byte("")},
			// draw the first update
			{row: 48, value: []byte("LineIdx:0")},
			{row: 49, value: []byte("LineIdx:1")},
			{row: 50, value: []byte("LineIdx:2")},
		},
		[]string{},
	},
	"FloatForward_TermHeightSmall_AtBottom_Header": {3, true, false, 49, PolicyFloatForwardTrail, 50,
		[]ScreenEvent{
			// create the frame (pave a blank spot)
			{row: 47, value: []byte("")},
			{row: 48, value: []byte("")},
			{row: 49, value: []byte("")},
			{row: 50, value: []byte("")},
			// draw the first update
			{row: 47, value: []byte("theHeader")},
			{row: 48, value: []byte("LineIdx:0")},
			{row: 49, value: []byte("LineIdx:1")},
			{row: 50, value: []byte("LineIdx:2")},
		},
		[]string{},
	},
	"FloatForward_termHeightSmall_AtBottom_Footer": {3, false, true, 49, PolicyFloatForwardTrail, 50,
		[]ScreenEvent{
			// create the frame (pave a blank spot)
			{row: 47, value: []byte("")},
			{row: 48, value: []byte("")},
			{row: 49, value: []byte("")},
			{row: 50, value: []byte("")},
			// draw the first update
			{row: 47, value: []byte("LineIdx:0")},
			{row: 48, value: []byte("LineIdx:1")},
			{row: 49, value: []byte("LineIdx:2")},
			{row: 50, value: []byte("theFooter")},
		},
		[]string{},
	},
	"FloatForward_TermHeightSmall_AtBottom_HeaderFooter": {3, true, true, 49, PolicyFloatForwardTrail, 50,
		[]ScreenEvent{
			// create the frame (pave a blank spot)
			{row: 46, value: []byte("")},
			{row: 47, value: []byte("")},
			{row: 48, value: []byte("")},
			{row: 49, value: []byte("")},
			{row: 50, value: []byte("")},
			// draw the first update
			{row: 46, value: []byte("theHeader")},
			{row: 47, value: []byte("LineIdx:0")},
			{row: 48, value: []byte("LineIdx:1")},
			{row: 49, value: []byte("LineIdx:2")},
			{row: 50, value: []byte("theFooter")},
		},
		[]string{},
	},
}

func Test_FloatForwardPolicy_Frame_Draw(t *testing.T) {

	names := make([]string, 0, len(floatForwardDrawTestCases))
	for name := range floatForwardDrawTestCases {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, test := range names {
		table := floatForwardDrawTestCases[test]
		suppressOutput(func() {
			// setup...
			terminalHeight = table.terminalHeight
			handler := NewTestEventHandler(t)
			screenHandlers = make([]ScreenEventHandler, 0)
			addScreenHandler(handler)

			// run test...
			var errs []error
			frame := New(Config{
				Lines:          table.rows,
				HasHeader:      table.hasHeader,
				HasFooter:      table.hasFooter,
				startRow:       table.startRow,
				PositionPolicy: table.policy,
			})
			if table.hasHeader {
				frame.Header.buffer = []byte("theHeader")
			}
			for idx, line := range frame.Lines {
				line.buffer = []byte(fmt.Sprintf("LineIdx:%d", idx))
			}
			if table.hasFooter {
				frame.Footer.buffer = []byte("theFooter")
			}
			errs = frame.Draw()

			// assert results...
			validateEvents(t, test, table, errs, frame, handler)

		})
	}

}

func Test_FloatForwardPolicy_Frame_AdhocDraw(t *testing.T) {

	names := make([]string, 0, len(floatForwardDrawTestCases))
	for name := range floatForwardDrawTestCases {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, test := range names {
		table := floatForwardDrawTestCases[test]
		suppressOutput(func() {
			// setup...
			terminalHeight = table.terminalHeight
			handler := NewTestEventHandler(t)
			screenHandlers = make([]ScreenEventHandler, 0)
			addScreenHandler(handler)

			// run test...
			var err error
			var errs = make([]error, 0)
			frame := New(Config{
				Lines:          table.rows,
				HasHeader:      table.hasHeader,
				HasFooter:      table.hasFooter,
				startRow:       table.startRow,
				PositionPolicy: table.policy,
			})
			if table.hasHeader {
				err = frame.Header.WriteString("theHeader")
				if err != nil {
					errs = append(errs, err)
				}
			}
			for idx, line := range frame.Lines {
				err = line.WriteString(fmt.Sprintf("LineIdx:%d", idx))
				if err != nil {
					errs = append(errs, err)
				}
			}
			if table.hasFooter {
				err = frame.Footer.WriteString("theFooter")
				if err != nil {
					errs = append(errs, err)
				}
			}

			// assert results...
			validateEvents(t, test, table, errs, frame, handler)

		})
	}

}
