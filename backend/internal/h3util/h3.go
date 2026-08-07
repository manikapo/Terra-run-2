package h3util

import (
	"fmt"

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
		cell, err := h3.LatLngToCell(h3.NewLatLng(p.Lat, p.Lon), resolution)
		if err != nil {
			return nil, fmt.Errorf("h3 cell: %w", err)
		}
		hex := h3.IndexToString(uint64(cell))
		if !seen[hex] {
			seen[hex] = true
			cells = append(cells, hex)
		}
	}
	return cells, nil
}

// ParentCell returns parent H3 index at target resolution.
func ParentCell(cellHex string, parentRes int) (string, error) {
	idx, err := h3.StringToIndex(cellHex)
	if err != nil {
		return "", err
	}
	parent, err := h3.CellToParent(h3.Cell(idx), parentRes)
	if err != nil {
		return "", err
	}
	return h3.IndexToString(uint64(parent)), nil
}

// TileCell parses tile path param (h3 index hex at tile resolution).
func ParseCellHex(hex string) (h3.Cell, error) {
	idx, err := h3.StringToIndex(hex)
	if err != nil {
		return 0, err
	}
	return h3.Cell(idx), nil
}

// CellBoundaryGeoJSON returns GeoJSON polygon coordinates for a cell.
func CellBoundaryGeoJSON(cellHex string) (map[string]interface{}, error) {
	cell, err := ParseCellHex(cellHex)
	if err != nil {
		return nil, err
	}
	boundary, err := h3.CellToBoundary(cell)
	if err != nil {
		return nil, err
	}

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
	children, err := h3.CellToChildren(cell, childRes)
	if err != nil {
		return nil, err
	}
	out := make([]string, len(children))
	for i, c := range children {
		out[i] = h3.IndexToString(uint64(c))
	}
	return out, nil
}
