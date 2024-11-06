package healthchecker

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

type CustomError struct {
	Expected int `json:"expected,omitempty"`
	Got      int `json:"got,omitempty"`
}

func (c *CustomError) Error() string {
	return fmt.Sprintf("Expected: %d got: %d", c.Expected, c.Got)
}

var (
	configSetup1 = Config{Name: "test 1"}
	configSetup2 = Config{Name: "test 1", Version: "v1"}
	configSetup3 = Config{
		Name:    "test 1",
		Version: "v1",
		Integrations: []Check{
			{
				Name:   "func 1",
				Handle: func() CheckResponse { return CheckResponse{} },
			},
		},
	}
)

func TestNew(t *testing.T) {
	type args struct {
		conf Config
	}

	tests := []struct {
		name string
		args args
		want *HealthCheck
	}{
		{
			name: "should create a new instance of HealthCheck with no arguments",
			args: args{Config{}},
			want: &HealthCheck{},
		},
		{
			name: "should create with only name",
			args: args{configSetup1},
			want: &HealthCheck{configSetup1},
		},
		{
			name: "should create with name and version",
			args: args{configSetup2},
			want: &HealthCheck{configSetup2},
		},
		{
			name: "should create with name version, and checks",
			args: args{configSetup3},
			want: &HealthCheck{configSetup3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.conf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHealthCheck_Liveness(t *testing.T) {
	type fields struct {
		config Config
	}
	tests := []struct {
		name   string
		fields fields
		want   Liveness
	}{
		{
			name:   "should return fully functional",
			fields: fields{config: configSetup2},
			want: Liveness{
				Status:  fullyFunctional,
				Version: "v1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HealthCheck{
				config: tt.fields.config,
			}
			if got := h.Liveness(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HealthCheck.Liveness() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHealthCheck_Readiness(t *testing.T) {
	type fields struct {
		config Config
	}
	tests := []struct {
		name   string
		fields fields
		want   Readiness
	}{
		{
			name:   "shoud fake return false",
			fields: fields{config: configSetup3},
			want: Readiness{
				Name:    configSetup3.Name,
				Version: configSetup3.Version,
				Status:  true,
				Integrations: []Integration{
					{
						Name: configSetup3.Integrations[0].Name,
					},
				},
			},
		},
		{
			name: "shoud fake return false",
			fields: fields{config: Config{
				Integrations: []Check{
					{
						Handle: func() CheckResponse {
							return CheckResponse{
								Error: errors.New("test 1"),
								URL:   "test v1",
							}
						},
					},
				},
			}},
			want: Readiness{
				Status: false,
			},
		},
		{
			name: "shoud perform a real fake integration",
			fields: fields{config: Config{
				Integrations: []Check{
					{
						Handle: func() CheckResponse {
							result := CheckResponse{
								URL: "https://github.com/statusss",
							}
							client := http.Client{}
							request, _ := http.NewRequest("GET", result.URL, nil)
							response, err := client.Do(request)
							if err != nil {
								fmt.Println("foi aqui dfp")
								result.Error = err
								return result
							}
							if response.StatusCode != http.StatusOK {
								result.Error = &CustomError{Expected: http.StatusOK, Got: response.StatusCode}
							}
							return result
						},
					},
				},
			}},
			want: Readiness{
				Status: false,
			},
		},
		{
			name: "shoud perform a real integration",
			fields: fields{config: Config{
				Integrations: []Check{
					{
						Handle: func() CheckResponse {
							result := CheckResponse{
								URL: "https://github.com/status",
							}
							client := http.Client{}
							request, _ := http.NewRequest("GET", result.URL, nil)
							response, err := client.Do(request)
							result.Error = err
							if err != nil {
								if response.StatusCode != http.StatusOK {
									result.Error = &CustomError{Expected: http.StatusOK, Got: response.StatusCode}
								}
							}
							return result
						},
					},
				},
			}},
			want: Readiness{
				Status: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HealthCheck{
				config: tt.fields.config,
			}
			got := h.Readiness()
			if got.Status != tt.want.Status {
				t.Errorf("Test Readiness() fail want: %v got: %v", tt.want.Status, got.Status)
			}
		})
	}
}
