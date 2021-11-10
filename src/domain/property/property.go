package property

import "github.com/superbkibbles/bookstore_utils-go/rest_errors"

const (
	STATUS_ACTIVE   = "active"
	STATUS_DEACTIVE = "deactive"
)

type Property struct {
	ID        string `json:"id"`
	AgencyID  string `json:"agency_id"`
	ComplexID string `json:"complex_id"`

	ComplexName    string `json:"complex_name"`
	Description    string `json:"description"`
	Title          string `json:"title"`
	FloorNumber    int64  `json:"floor_nummber"`
	BuildingNumber int64  `json:"building_number"`
	DirectionFace  string `json:"direction_face"`
	PropertyType   string `json:"property_type"`
	BuiltYear      int64  `json:"built_year"`
	Price          int64  `json:"price"`
	Rooms          int64  `json:"rooms"`
	Bathrooms      int64  `json:"bathrooms"`
	Bedrooms       int64  `json:"bedrooms"`
	LivingRoom     int64  `json:"living_rooms"`
	Hall           int64  `json:"hall"`
	Balcony        int64  `json:"balcony"`
	Kitchen        int64  `json:"kitchen"`
	PropertyKind   string `json:"property_kind"`

	Promoted bool `json:"promoted"`

	Space        float64 `json:"space"`
	BuildingSize float64 `json:"building_size"`
	Area         float64 `json:"area"`

	Location    string      `json:"location"`
	City        string      `json:"city"`
	GPS         coordinates `json:"gps"`
	NearSchools []school    `json:"near_schools"`

	Visuals []Visual `json:"visuals"`

	ForRent     bool   `json:"for_rent"`
	PropertyNo  int64  `json:"property_no"`
	Viewers     int64  `json:"Viewers"`
	Status      string `json:"status"`
	DateCreated string `json:"date_created"`
	IsSold      bool   `json:"is_sold"`
	SoldDate    string `json:"sold_date"`
}

type Visual struct {
	Url string `json:"url"`
}

type school struct {
	Name string `json:"name"`
}

type coordinates struct {
	Lat  string `json:"lat"`
	Long string `json:"long"`
}

type Properties []Property

func (p *Property) Validate() rest_errors.RestErr {
	return nil
}

// Create another struct to get the PIC or photo request
// Get array of pics and photos
