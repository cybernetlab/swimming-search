package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cybernetlab/swimming-search/internal/domain"

	. "github.com/cybernetlab/swimming-search/test/support"
)

func withName(c domain.Course, name string) domain.Course {
	c.Name = name
	return c
}

func withQuery(s domain.Search, query string) domain.Search {
	s.NameQuery = query
	return s
}

func withDay(c domain.Course, day string) domain.Course {
	c.Schedule.DayOfWeek = day
	return c
}

func withDays(s domain.Search, days ...uint8) domain.Search {
	s.Days = days
	return s
}

func Test_IsMatch_ByName(t *testing.T) {
	c := domain.Course{Schedule: domain.CourseSchedule{DayOfWeek: "Monday"}}
	s := domain.Search{Days: []uint8{1, 2, 3, 4, 5, 6, 7}}

	require.True(t, withName(c, "some test").IsMatch(withQuery(s, "test")))
	require.True(t, withName(c, "some TEST").IsMatch(withQuery(s, "test")))
	require.True(t, withName(c, "some test").IsMatch(withQuery(s, "TEST")))
	require.False(t, withName(c, "some test").IsMatch(withQuery(s, "other")))
}

func Test_IsMatch_ByDays(t *testing.T) {
	c := domain.Course{Name: "test"}
	s := domain.Search{NameQuery: "test"}

	require.True(t, withDay(c, "Monday").IsMatch(withDays(s, 1)))
	require.True(t, withDay(c, "Tuesday").IsMatch(withDays(s, 2)))
	require.True(t, withDay(c, "Wednesday").IsMatch(withDays(s, 3)))
	require.True(t, withDay(c, "Thursday").IsMatch(withDays(s, 4)))
	require.True(t, withDay(c, "Friday").IsMatch(withDays(s, 5)))
	require.True(t, withDay(c, "Saturday").IsMatch(withDays(s, 6)))
	require.True(t, withDay(c, "Sunday").IsMatch(withDays(s, 7)))
	require.False(t, withDay(c, "Monday").IsMatch(withDays(s, 2, 3, 4, 5, 6, 7)))
	require.False(t, withDay(c, "").IsMatch(withDays(s, 1, 2, 3, 4, 5, 6, 7)))
	require.False(t, withDay(c, "Invalid").IsMatch(withDays(s, 1, 2, 3, 4, 5, 6, 7)))
}

func Test_Unmarshal(t *testing.T) {
	c := ReadFixture[domain.Course](t, "course")

	require.Equal(t, 84, c.ID)
	require.Equal(t, "course", c.Type)
	require.Equal(t, "Learn2Swim - Preschool 1 (3-4yrs)", c.Name)
	require.Equal(t, 10, c.Size)
	require.Equal(t, []string{}, c.Restrictions)
	require.Equal(t, true, c.PaymentGatewayValid)
	require.Equal(t, "dayOfWeek", c.Schedule.Type)
	require.Equal(t, "Weekly from 16 March 2026", c.Schedule.DateDescr)
	require.Equal(t, "15:35 - 16:05", c.Schedule.TimeDescr)
	require.Equal(t, "15:35:00", c.Schedule.Time.Start)
	require.Equal(t, "16:05:00", c.Schedule.Time.End)
	require.Equal(t, "Monday", c.Schedule.DayOfWeek)
	require.Equal(t, 1, c.CourseGroup.ID)
	require.Equal(t, "Learn2Swim", c.CourseGroup.Name)
	require.Equal(t, "", c.CourseGroup.Icon)
	require.Equal(t, 1, c.CourseGroupCategory.ID)
	require.Equal(t, "Swimming", c.CourseGroupCategory.Name)
	require.Equal(t, "swimming", c.CourseGroupCategory.Icon)
	require.Equal(t, 3, c.Centre.ID)
	require.Equal(t, "Scott Hall", c.Centre.Name)
	require.Equal(t, "24.50", c.PriceIndication.Price)
	require.Equal(t, "month", c.PriceIndication.Per)
	require.Equal(t, "This is the description\nfor a level", c.Level.Description)
	require.Equal(t, time.Date(2026, 3, 16, 15, 35, 0, 0, time.UTC).UTC(), c.Availability.StartDate.Time)
	require.Equal(t, 1, c.Availability.EnablePortalSpaces)
	require.Equal(t, 10, c.Availability.Spaces.All)
	require.Equal(t, 10, c.Availability.Spaces.Occupied)
	require.Equal(t, 0, c.Availability.Spaces.Reserved)
	require.Equal(t, 0, c.Availability.Spaces.Free)
}

func Test_Unmarshal_NoDate(t *testing.T) {
	var zeroTime time.Time
	c := ReadFixture[domain.Course](t, "course", Subst{Key: "startDate", Value: `null`})
	require.Equal(t, domain.CustomDate{zeroTime}, c.Availability.StartDate)
}
