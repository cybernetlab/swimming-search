package domain

type Filter struct {
	ShowFullCourses     bool    `json:"showFullCourses"`
	Centre              uint    `json:"centre"`
	CourseGroupCategory []uint  `json:"courseGroupCategory"`
	Offset              int     `json:"offset"`
	Limit               int     `json:"limit"`
	DayOfWeek           []uint8 `json:"dayOfWeek"`
	Region              uint    `json:"region"`
	CourseType          string  `json:"courseType,omitempty"`
}
