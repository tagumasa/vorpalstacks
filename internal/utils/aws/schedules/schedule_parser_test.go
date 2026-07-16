package schedules

import (
	"testing"
)

func TestGetScheduleType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ScheduleType
		wantErr bool
	}{
		{"valid rate expression", "rate(1 minute)", ScheduleTypeRate, false},
		{"valid cron expression", "cron(0 12 * * ? *)", ScheduleTypeCron, false},
		{"rate with spaces", "rate( 1 hour )", ScheduleTypeRate, false},
		{"cron with spaces", "cron( 0 12 * * ? * )", ScheduleTypeCron, false},
		{"empty string", "", "", true},
		{"invalid - no prefix", "every 5 minutes", "", true},
		{"invalid - wrong prefix", "interval(5 minutes)", "", true},
		{"case sensitive rate", "RATE(1 minute)", "", true},
		{"case sensitive cron", "CRON(0 12 * * ? *)", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetScheduleType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetScheduleType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetScheduleType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRateSchedule(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *RateSchedule
		wantErr bool
	}{
		{"valid rate 1 minute", "rate(1 minute)", &RateSchedule{Value: 1, Unit: "minute"}, false},
		{"valid rate 1 hour", "rate(1 hour)", &RateSchedule{Value: 1, Unit: "hour"}, false},
		{"valid rate 1 day", "rate(1 day)", &RateSchedule{Value: 1, Unit: "day"}, false},
		{"valid rate 6 hours", "rate(6 hours)", &RateSchedule{Value: 6, Unit: "hours"}, false},
		{"valid rate 7 days", "rate(7 days)", &RateSchedule{Value: 7, Unit: "days"}, false},
		{"valid rate plural minutes", "rate(1 minutes)", &RateSchedule{Value: 1, Unit: "minutes"}, false},
		{"valid rate plural hours", "rate(6 hours)", &RateSchedule{Value: 6, Unit: "hours"}, false},
		{"valid rate plural days", "rate(7 days)", &RateSchedule{Value: 7, Unit: "days"}, false},
		{"invalid - not rate", "cron(0 12 * * ? *)", nil, true},
		{"invalid - zero value", "rate(0 minute)", nil, true},
		{"invalid - negative value", "rate(-1 hour)", nil, true},
		{"valid rate 2 minutes", "rate(2 minutes)", &RateSchedule{Value: 2, Unit: "minutes"}, false},
		{"valid rate 3 hours", "rate(3 hours)", &RateSchedule{Value: 3, Unit: "hours"}, false},
		{"valid rate 5 days", "rate(5 days)", &RateSchedule{Value: 5, Unit: "days"}, false},
		{"invalid - missing unit", "rate(1)", nil, true},
		{"invalid - missing value", "rate(minute)", nil, true},
		{"invalid - invalid unit", "rate(1 week)", nil, true},
		{"empty string", "", nil, true},
		{"rate with extra spaces", "rate(  1  hour  )", &RateSchedule{Value: 1, Unit: "hour"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRateSchedule(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRateSchedule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && got != nil {
				if got.Value != tt.want.Value || got.Unit != tt.want.Unit {
					t.Errorf("GetRateSchedule() = %+v, want %+v", got, tt.want)
				}
			} else if tt.want != nil || got != nil {
				t.Errorf("GetRateSchedule() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
