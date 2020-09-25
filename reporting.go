package shure

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"
	"time"
)

var reportRegex = regexp.MustCompile("< REP ([0-9])? ?([A-Z,_]*) (.*)? ?>")

const (
	_fullReportIndex = 0
	_channelIndex    = 1
	_typeIndex       = 2
	_valueIndex      = 3
)
const ErrorReportType = "ERROR"

type Report struct {
	Type       string
	Value      string
	Channel    int
	FullReport string
}

func (u *ULXDReceiver) StartReporting(ctx context.Context) (chan Report, error) {

	// Connect to the receiver

	timeout := 30 * time.Second // default timeout
	// Set timeout according to context if a deadline exists
	if d, ok := ctx.Deadline(); ok {
		timeout = time.Until(d)
	}

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:2202", u.address), timeout)
	if err != nil {
		return nil, fmt.Errorf("Error while connecting to the receiver: %s", err)
	}

	// Create reader
	r := bufio.NewReader(conn)
	c := make(chan Report, 10) // Make slightly buffered channel

	// Start monitoring
	go monitorReporting(r, c)
	return c, nil
}

func monitorReporting(r *bufio.Reader, c chan Report) {
	for {

		// Read the Report
		data, err := r.ReadString('>')
		if err != nil {
			// If the connection was closed
			if errors.Is(err, io.EOF) {
				c <- Report{
					Type:       ErrorReportType,
					Value:      "ConnectionClosedError",
					Channel:    -1,
					FullReport: "The connection to the receiver was closed",
				}
				close(c)
				return
			}

			c <- Report{
				Type:       ErrorReportType,
				Value:      "ReadError",
				Channel:    -1,
				FullReport: fmt.Sprintf("Error while reading from receiver: %s", err),
			}
			continue
		}

		c <- parseReport(data)

	}
}

func parseReport(data string) Report {
	// Parse Raw Data
	parts := reportRegex.FindStringSubmatch(data)

	// Handle matching error
	if len(parts) == 0 {
		return Report{
			Type:       ErrorReportType,
			Channel:    -1,
			Value:      "ParseError",
			FullReport: fmt.Sprintf("Report did not match expected format: %s", data),
		}
	}

	// Translate channel to int
	var err error
	channel := 0
	// -1 if there is no channel
	if len(parts[_channelIndex]) == 0 {
		channel = -1
	} else { // otherwise convert to int
		channel, err = strconv.Atoi(parts[_channelIndex])
		// Error if we can't convert channel to an int
		if err != nil {
			return Report{
				Type:       ErrorReportType,
				Channel:    -1,
				Value:      "ParseError",
				FullReport: fmt.Sprintf("Error while parsing channel: %s", err),
			}
		}

	}

	// Populate Report
	return Report{
		Type:       parts[_typeIndex],
		Channel:    channel,
		Value:      parts[_valueIndex],
		FullReport: parts[_fullReportIndex],
	}
}
