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
	"strings"
	"time"
)

type ReportType int

const (
	ERROR = iota + 1
	BATTERY_CYCLES
	BATTERY_CHARGE_MINUTES
	BATTERY_TYPE
	INTERFERENCE
	POWER
)

var channelRegex = regexp.MustCompile("REP [\\d]")

type Report struct {
	Type    ReportType
	Value   string
	Channel int
	Message string
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
					Type:    ERROR,
					Value:   "ConnectionClosedError",
					Channel: -1,
					Message: "The connection to the receiver was closed",
				}
				close(c)
				return
			}

			c <- Report{
				Type:    ERROR,
				Value:   "ReadError",
				Channel: -1,
				Message: fmt.Sprintf("Error while reading from receiver: %s", err),
			}
			continue
		}

		// If we have a report then send it, else continue
		if report, ok := parseReport(data); ok {
			c <- report
		}

	}
}

func parseReport(data string) (Report, bool) {
	report := Report{}

	// Get Channel
	channel := channelRegex.FindString(data)
	if len(channel) == 0 {
		// TODO: Confirm that below assumption is valid?
		// No relevant data so skip
		return report, false
	}

	c, err := strconv.Atoi(channel[len(channel)-1:])
	if err != nil {
		// Report error
		report.Channel = -1
		report.Type = ERROR
		report.Value = "ParseError"
		report.Message = fmt.Sprintf("Error while parsing channel: %s", err)
		return report, true
	}

	report.Channel = c

	// Determine event type and populate data
	if strings.Contains(data, "RF_INT_DET") {
		populateInterferenceReport(data, &report)
	} else if strings.Contains(data, "TX_TYPE") {
		populatePowerReport(data, &report)
	} else if strings.Contains(data, "BATT_CYCLE") {
		populateBatteryCyclesReport(data, &report)
	} else if strings.Contains(data, "BATT_RUN_TIME") {
		populateBatteryChargeMinutesReport(data, &report)
	} else if strings.Contains(data, "BATT_TYPE") {
		populateBatteryTypeReport(data, &report)
	} else {
		// Unrecognized event
		report.Channel = -1
		report.Type = ERROR
		report.Value = "UnrecognizedReport"
		report.Message = fmt.Sprintf("Encountered an unrecognized report: %s", data)
	}

	return report, true
}

func populateInterferenceReport(data string, r *Report) {
	r.Type = INTERFERENCE

	// Determine interference
	if strings.Contains(data, "NONE") {
		r.Value = "NONE"
		r.Message = fmt.Sprintf("No interference on channel %d", r.Channel)
	} else if strings.Contains(data, "CRITICAL") {
		r.Value = "CRITICAL"
		r.Message = fmt.Sprintf("Interference on channel %d", r.Channel)
	} else {
		// Invalid data
		r.Type = ERROR
		r.Channel = -1
		r.Value = "ParseError"
		r.Message = "Invalid value for interference report"
		return
	}
}

func populatePowerReport(data string, r *Report) {
	r.Type = POWER

	// Determine power state
	if strings.Contains(data, "UNKN") {
		r.Value = "STANDBY"
		// TODO: Is this by channel?
		r.Message = "Power is in standby mode"
	} else {
		r.Value = "ON"
		r.Message = "Power is on"
	}
}

func populateBatteryCyclesReport(data string, r *Report) {
	r.Type = BATTERY_CYCLES

	re := regexp.MustCompile("[1-9][0-9]*")
	cycles := re.FindString(data)

	switch cycles {
	case "65535":
		r.Value = "UNKNOWN"
		r.Message = fmt.Sprintf("Channel %d has an unknown number of battery cycles", r.Channel)
	case "":
		r.Value = "0"
		r.Message = fmt.Sprintf("Channel %d has 0 battery cycles", r.Channel)
	default:
		r.Value = cycles
		r.Message = fmt.Sprintf("Channel %d has %s battery cyles", r.Channel, cycles)
	}
}

func populateBatteryChargeMinutesReport(data string, r *Report) {
	r.Type = BATTERY_CHARGE_MINUTES

	re := regexp.MustCompile("[1-9][0-9]*")
	time := re.FindString(data)

	// TODO: Other values?
	switch time {
	case "65535":
		r.Value = "UNKNOWN"
		r.Message = fmt.Sprintf("Channel %d has an unknown number of minutes left", r.Channel)
	case "65534":
		r.Value = "CALCULATING"
		r.Message = fmt.Sprintf("Channel %d is calculating the number of minutes left", r.Channel)
	case "":
		r.Value = "0"
		r.Message = fmt.Sprintf("Channel %d has 0 minutes left", r.Channel)
	default:
		r.Value = time
		r.Message = fmt.Sprintf("Channel %d has %s minutes left", r.Channel, time)
	}
}

func populateBatteryTypeReport(data string, r *Report) {
	r.Type = BATTERY_TYPE

	re := regexp.MustCompile("[\\s][A-Z]{4}[\\s]")
	batteryType := re.FindString(data)

	switch batteryType {
	case " UNKN ":
		r.Value = "UNKNOWN"
		r.Message = fmt.Sprintf("Channel %d has an unknown battery type", r.Channel)
	default:
		r.Value = strings.TrimSpace(batteryType)
		r.Message = fmt.Sprintf("Channel %d has a %s battery type", r.Channel, r.Value)
	}
}
