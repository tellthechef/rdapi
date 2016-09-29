package models

type Restaurant struct {
	Name     string
	Segments []struct {
		RestaurantConsolidated
		Name string
	}
}

type RestaurantConsolidated struct {
	Areas []struct {
		AreaID int
		Name   string
	}
	Tables []struct {
		TableID int
		Number  string
	}
}

// Consolidate merges all segments' tables and areas..
func (restaurant Restaurant) Consolidate() RestaurantConsolidated {
	var r RestaurantConsolidated
	for _, segment := range restaurant.Segments {
		r.Areas = append(r.Areas, segment.Areas...)
		r.Tables = append(r.Tables, segment.Tables...)
	}

	return r
}

func (r RestaurantConsolidated) GetArea(id int) string {
	for _, area := range r.Areas {
		if area.AreaID == id {
			return area.Name
		}
	}

	return ""
}

func (r RestaurantConsolidated) GetTable(id int) string {
	for _, table := range r.Tables {
		if table.TableID == id {
			return table.Number
		}
	}

	return ""
}
