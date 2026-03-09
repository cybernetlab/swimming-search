package domain

import (
	"strings"
	"time"
)

type CourseSchedule struct {
	Type      string
	DateDescr string `json:"dateDescription"`
	TimeDescr string `json:"timeDescription"`
	Time      struct {
		Start string
		End   string
	}
	DayOfWeek string
}

type CourseGroup struct {
	ID   int `json:"id"`
	Name string
	Icon string `json:"homeportal_icon"`
}

type Course struct {
	ID           int `json:"courseId"`
	Type         string
	Name         string
	Size         int
	Schedule     CourseSchedule
	Availability struct {
		StartDate CustomDate
		Spaces    struct {
			All      int
			Occupied int
			Reserved int
			Free     int
		}
		EnablePortalSpaces int `json:"enable_portal_spaces"`
	}
	Restrictions    []string
	PriceIndication struct {
		Price string
		Per   string
	}
	Level struct {
		Description string
	}
	CourseGroup         CourseGroup
	CourseGroupCategory CourseGroup
	Centre              Centre
	PaymentGatewayValid bool
}

type CustomDate struct {
	time.Time
}

func (c *CustomDate) UnmarshalJSON(b []byte) error {
	var err error

	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return nil
	}
	c.Time, err = time.Parse("2006-01-02 15:04:05", s)
	return err
}

func (c Course) IsMatch(search Search) bool {
	if c.Schedule.DayOfWeek == "" {
		return false
	}
	if !strings.Contains(strings.ToLower(c.Name), strings.ToLower(search.NameQuery)) {
		return false
	}
	day := dayToInt(c.Schedule.DayOfWeek)
	for _, x := range search.Days {
		if x == day {
			return true
		}
	}
	return false
}

func dayToInt(day string) uint8 {
	switch day {
	case "Monday":
		return 1
	case "Tuesday":
		return 2
	case "Wednesday":
		return 3
	case "Thursday":
		return 4
	case "Friday":
		return 5
	case "Saturday":
		return 6
	case "Sunday":
		return 7
	}
	return 0
}
