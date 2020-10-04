package models

// TripQuery is the struct for making trip-related queries based on the indexed key field
type TripQuery struct {
	Value          string
	FilterFinished bool
	Finished       bool
}
