package crs

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantCode  int
		wantYX    bool
		wantEmpty bool
	}{
		{
			name:     "EPSG short form",
			input:    "EPSG:25832",
			wantCode: 25832,
		},
		{
			name:     "URN form no version",
			input:    "urn:ogc:def:crs:EPSG::25832",
			wantCode: 25832,
		},
		{
			name:     "URN form with version",
			input:    "urn:ogc:def:crs:EPSG:6.12:25832",
			wantCode: 25832,
		},
		{
			name:     "HTTP form",
			input:    "http://www.opengis.net/def/crs/EPSG/0/25832",
			wantCode: 25832,
		},
		{
			name:     "HTTPS form",
			input:    "https://www.opengis.net/def/crs/EPSG/0/4326",
			wantCode: 4326,
			wantYX:   true,
		},
		{
			name:     "WGS84 Y,X order",
			input:    "EPSG:4326",
			wantCode: 4326,
			wantYX:   true,
		},
		{
			name:     "ETRS89 Y,X order",
			input:    "EPSG:4258",
			wantCode: 4258,
			wantYX:   true,
		},
		{
			name:     "UTM zone 32N X,Y order",
			input:    "EPSG:32632",
			wantCode: 32632,
			wantYX:   false,
		},
		{
			name:     "ADV UTM32 compound CRS",
			input:    "urn:adv:crs:ETRS89_UTM32*DE_DHHN92_NH",
			wantCode: 25832,
		},
		{
			name:     "ADV UTM33 compound CRS",
			input:    "urn:adv:crs:ETRS89_UTM33*DE_DHHN92_NH",
			wantCode: 25833,
		},
		{
			name:      "unrecognized format",
			input:     "some-unknown-crs",
			wantCode:  0,
			wantEmpty: true,
		},
		{
			name:      "empty string",
			input:     "",
			wantCode:  0,
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Parse(tt.input)
			if c.Raw != tt.input {
				t.Errorf("Raw = %q, want %q", c.Raw, tt.input)
			}

			if c.Code != tt.wantCode {
				t.Errorf("Code = %d, want %d", c.Code, tt.wantCode)
			}

			if c.IsYXOrder != tt.wantYX {
				t.Errorf("IsYXOrder = %v, want %v", c.IsYXOrder, tt.wantYX)
			}
		})
	}
}
