package domain

type Centre struct {
	Type   string
	ID     uint `json:"value"`
	Name   string
	Region struct {
		ID   uint `json:"value"`
		Name string
	}
	CourseGroupCategoryIds []uint
}
