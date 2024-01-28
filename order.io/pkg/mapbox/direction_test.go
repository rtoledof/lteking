package mapbox

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"order.io/pkg/order"
)

func TestDirectionService_GetRoute(t *testing.T) {
	var tests = []struct {
		name    string
		server  *httptest.Server
		request order.DirectionRequest
		want    *order.DirectionResponse
		wantErr bool
	}{
		{
			name: "valid request",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				data, _ := os.ReadFile("testdata/direction_ok.json")
				w.Write(data)
			})),
			request: order.DirectionRequest{
				Points: []*order.Point{
					{
						Lat: 13.426579,
						Lng: 52.508068,
					},
					{
						Lat: 13.427292,
						Lng: 52.506902,
					},
				},
			},
			want: &order.DirectionResponse{
				Routes: []*order.Route{
					{
						Geometry: "mnn_Ick}pAfBiF`CzA",
						Duration: 26.2,
						Distance: 176.7,
						Waitpoints: []*order.WaitPoint{
							{
								Location: []float64{13.426579, 52.508068},
								Name:     "Köpenicker Straße",
							},
							{
								Location: []float64{13.427292, 52.506902},
								Name:     "Engeldamm",
							},
						},
						Legs: []*order.Legs{
							{
								Summary:  "Köpenicker Straße, Engeldamm",
								Duration: 26.2,
								Distance: 176.7,
								Weight:   44.4,
								Steps: []order.Step{
									{
										Intersections: []*order.Intersection{
											{
												Location: []float64{13.426579, 52.508068},
												Bearings: []int{125},
												Entry:    []bool{true},
											},
											{
												Location: []float64{13.426688, 52.508022},
												Bearings: []int{30, 120, 300},
												Entry:    []bool{true, true, false},
												In:       2,
												Out:      1,
											},
										},
										DrivingSide: "right",
										Geometry:    "mnn_Ick}pAHUlAqDNa@",
										Mode:        "driving",
										Maneuver: &order.Maneuver{
											Location:      []float64{13.426579, 52.508068},
											BearingBefore: 0,
											BearingAfter:  125,
											Type:          "depart",
											Modifier:      "right",
											Instruction:   "Head southeast on Köpenicker Straße (L 1066)",
										},
										Ref:      "L 1066",
										Weight:   35.9,
										Duration: 17.7,
										Name:     "Köpenicker Straße (L 1066)",
										Distance: 98.1,
										VoiceInstructions: []*order.VoiceInstructions{
											{
												Announcement:          "Head southeast on Köpenicker Straße (L 1066), then turn right onto Engeldamm",
												DistanceAlongGeometry: 98.1,
												SsmlAnnouncement:      "<speak><amazon:effect name=\"drc\"><prosody rate=\"1.08\">Head southeast on Köpenicker Straße (L <say-as interpret-as=\"address\">1066</say-as>), then turn right onto Engeldamm</prosody></amazon:effect></speak>",
											},
											{
												Announcement:          "Turn right onto Engeldamm, then you will arrive at your destination",
												DistanceAlongGeometry: 83.1,
												SsmlAnnouncement:      "<speak><amazon:effect name=\"drc\"><prosody rate=\"1.08\">Turn right onto Engeldamm, then you will arrive at your destination</prosody></amazon:effect></speak>",
											},
										},
										BannerInstructions: []*order.BannerInstructions{
											{
												DistanceAlongGeometry: 98.1,
												Primary: &order.Primary{
													Text:       "Engeldamm",
													Type:       "turn",
													Modifier:   "right",
													Components: []*order.Component{{Text: "Engeldamm"}},
												},
											},
										},
									},
									{
										Intersections: []*order.Intersection{
											{
												Location: []float64{13.427752, 52.50755},
												Bearings: []int{30, 120, 210, 300},
												Entry:    []bool{false, true, true, false},
												In:       3,
												Out:      2,
											},
										},
										DrivingSide: "right",
										Geometry:    "ekn_Imr}pARL\\T^RHDd@\\",
										Mode:        "driving",
										Maneuver: &order.Maneuver{
											Location:      []float64{13.427752, 52.50755},
											BearingBefore: 125,
											BearingAfter:  202,
											Type:          "turn",
											Modifier:      "right",
											Instruction:   "Turn right onto Engeldamm",
										},
										Weight:   8.5,
										Duration: 8.5,
										Name:     "Engeldamm",
										Distance: 78.6,
										VoiceInstructions: []*order.VoiceInstructions{
											{
												Announcement:          "You have arrived at your destination",
												DistanceAlongGeometry: 27.7,
												SsmlAnnouncement:      "<speak><amazon:effect name=\"drc\"><prosody rate=\"1.08\">You have arrived at your destination</prosody></amazon:effect></speak>",
											},
										},
										BannerInstructions: []*order.BannerInstructions{
											{
												DistanceAlongGeometry: 78.6,
												Primary: &order.Primary{
													Text:       "You will arrive at your destination",
													Type:       "arrive",
													Modifier:   "straight",
													Components: []*order.Component{{Text: "You will arrive at your destination"}},
												},
												Secondary: &order.Primary{
													Text:       "Engeldamm",
													Type:       "arrive",
													Modifier:   "straight",
													Components: []*order.Component{{Text: "Engeldamm"}},
												},
											},
											{
												DistanceAlongGeometry: 15,
												Primary: &order.Primary{
													Text:       "You have arrived at your destination",
													Type:       "arrive",
													Modifier:   "straight",
													Components: []*order.Component{{Text: "You have arrived at your destination"}},
												},
												Secondary: &order.Primary{
													Text:       "Engeldamm",
													Components: []*order.Component{{Text: "Engeldamm"}},
												},
											},
										},
									},
									{
										Intersections: []*order.Intersection{
											{
												Location: []float64{13.427292, 52.506902},
												Bearings: []int{25},
												Entry:    []bool{true},
											},
										},
										DrivingSide: "right",
										Geometry:    "cgn_Iqo}pA",
										Mode:        "driving",
										Maneuver: &order.Maneuver{
											Location:      []float64{13.427292, 52.506902},
											BearingBefore: 205,
											BearingAfter:  0,
											Type:          "arrive",
											Instruction:   "You have arrived at your destination",
										},
										Name:               "Engeldamm",
										VoiceInstructions:  []*order.VoiceInstructions{},
										BannerInstructions: []*order.BannerInstructions{},
									},
								},
							},
						},
						WeightName: "auto",
						Weight:     44.4,
					},
				},
			},
		},
	}

	c := NewClient("pk.eyJ1IjoiY3ViYXdoZWVsZXIiLCJhIjoiY2s5Z2Z6Z2Z4MDJ6ZjNscWx6Z2Z")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.BaseURL, _ = url.Parse(tt.server.URL)
			c.client = tt.server.Client()
			got, _, err := c.Directions.GetRoute(context.Background(), tt.request)
			if err != nil && !tt.wantErr {
				t.Errorf("DirectionService.GetRoute() error = %v", err)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("DirectionService.GetRoute() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
