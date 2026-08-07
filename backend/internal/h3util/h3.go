package h3util

import (
	"github.com/uber/h3-go/v4"
)

// CellsFromRoute returns unique H3 cell indexes (as hex strings) along a GPS route.
func CellsFromRoute(points []struct{ Lat, Lon float64 }, resolution int) ([]string, error) {
	if len(points) == 0 {
		return nil, nil
	}

	seen := make(map[string]bool)
	var cells []string

	for _, p := range points {
		cell := h3.LatLngToCell(h3.NewLatLng(p.Lat, p.Lon), resolution)
		hex := cell.String()
		if !seen[hex] {
			seen[hex] = true
			cells = append(cells, hex)
		}
	}
	return cells, nil
}

// ParentCell returns parent H3 index at target resolution.
func ParentCell(cellHex string, parentRes int) (string, error) {
	cell, err := h3.IndexFromString(cellHex)
	if err != nil {
		return "", err
	}
	parent := cell.Parent(parentRes)
	return parent.String(), nil
}

// ParseCellHex parses an H3 index hex string.
func ParseCellHex(hex string) (h3.Cell, error) {
	return h3.IndexFromString(hex)
}

// CellBoundaryGeoJSON returns GeoJSON polygon coordinates for a cell.
func CellBoundaryGeoJSON(cellHex string) (map[string]interface{}, error) {
	cell, err := ParseCellHex(cellHex)
	if err != nil {
		return nil, err
	}
	boundary := cell.Boundary()

	coords := make([][]float64, 0, len(boundary)+1)
	for _, latLng := range boundary {
		coords = append(coords, []float64{latLng.Lng, latLng.Lat})
	}
	if len(coords) > 0 {
		coords = append(coords, coords[0]) // close ring
	}

	return map[string]interface{}{
		"type":        "Polygon",
		"coordinates": [][][]float64{coords},
	}, nil
}

// ChildrenAtResolution returns child cell hex strings for a tile cell.
func ChildrenAtResolution(tileHex string, childRes int) ([]string, error) {
	cell, err := ParseCellHex(tileHex)
	if err != nil {
		return nil, err
	}
	children := cell.Children(childRes)
	out := make([]string, len(children))
	for i, c := range children {
		out[i] = c.String()
	}
	return out, nil
}

// CellToInt64 converts hex to int64 for DB storage.
func CellToInt64(cellHex string) (int64, error) {
	cell, err := ParseCellHex(cellHex)
	if err != nil {
		return 0, err
	}
	return int64(cell), nil
}
