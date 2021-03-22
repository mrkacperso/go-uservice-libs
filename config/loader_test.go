package config

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	type SingleString struct {
		DbHost string
	}

	type AllTypes struct {
		ValString   string
		ValInt      int
		ValDuration time.Duration
		ValTime     time.Time
	}

	type args struct {
		configPath string
		entries    []Entry
		//target           interface{}
		target           AllTypes
		want             AllTypes
		configVarsString map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Plain 1 env variable - no defaults, all required",
			args: args{
				configPath: "",
				entries: []Entry{
					{
						Name:      "VAL_STRING",
						EntryType: new(string),
						Default:   "some string",
						Required:  true,
					},
				},
				target: AllTypes{},
				configVarsString: map[string]string{
					"VAL_STRING": "localhost",
				},
				want: AllTypes{ValString: "localhost"},
			},
			wantErr: false,
		},
		{
			name: "All variable types - no defaults, all required",
			args: args{
				configPath: "",
				entries: []Entry{
					{
						Name:      "VAL_STRING",
						EntryType: new(string),
						Default:   "",
						Required:  true,
					},
					{
						Name:      "VAL_INT",
						EntryType: new(int),
						Default:   "",
						Required:  true,
					},
					{
						Name:      "VAL_DURATION",
						EntryType: new(time.Duration),
						Default:   "",
						Required:  true,
					},
					{
						Name:      "VAL_TIME",
						EntryType: new(time.Time),
						Default:   "",
						Required:  true,
					},
				},
				target: AllTypes{},
				configVarsString: map[string]string{
					"VAL_STRING":   "ValStringContent",
					"VAL_INT":      "2",
					"VAL_DURATION": "10s",
					"VAL_TIME":     "2006-01-02T15:04:05.999999999",
				},
				want: AllTypes{
					ValString:   "ValStringContent",
					ValInt:      2,
					ValDuration: 10 * time.Second,
					ValTime:     time.Date(2006, time.January, 02, 15, 04, 05, 999999999, time.UTC),
				},
			},
			wantErr: false,
		},
		{
			name: "All variable types - unset but required",
			args: args{
				configPath: "",
				entries: []Entry{
					{
						Name:      "VAL_STRING",
						EntryType: new(string),
						Default:   "",
						Required:  true,
					},
					{
						Name:      "VAL_INT",
						EntryType: new(int),
						Default:   "",
						Required:  true,
					},
					{
						Name:      "VAL_DURATION",
						EntryType: new(time.Duration),
						Default:   "",
						Required:  true,
					},
					{
						Name:      "VAL_TIME",
						EntryType: new(time.Time),
						Default:   "",
						Required:  true,
					},
				},
				target: AllTypes{},
				configVarsString: map[string]string{
					//"VAL_STRING":   "ValStringContent",
					"VAL_INT":      "2",
					"VAL_DURATION": "10s",
					"VAL_TIME":     "2006-01-02T15:04:05.999999999",
				},
				want: AllTypes{
					ValString:   "ValStringContent",
					ValInt:      2,
					ValDuration: 10 * time.Second,
					ValTime:     time.Date(2006, time.January, 02, 15, 04, 05, 999999999, time.UTC),
				},
			},
			wantErr: true,
		},
		{
			name: "All variable types - unset but NOT required",
			args: args{
				configPath: "",
				entries: []Entry{
					{
						Name:      "VAL_STRING",
						EntryType: new(string),
						Default:   "",
						Required:  false,
					},
					{
						Name:      "VAL_INT",
						EntryType: new(int),
						Default:   "",
						Required:  true,
					},
					{
						Name:      "VAL_DURATION",
						EntryType: new(time.Duration),
						Default:   "",
						Required:  true,
					},
					{
						Name:      "VAL_TIME",
						EntryType: new(time.Time),
						Default:   "",
						Required:  true,
					},
				},
				target: AllTypes{},
				configVarsString: map[string]string{
					"VAL_INT":      "2",
					"VAL_DURATION": "10s",
					"VAL_TIME":     "2006-01-02T15:04:05.999999999",
				},
				want: AllTypes{
					ValString:   "",
					ValInt:      2,
					ValDuration: 10 * time.Second,
					ValTime:     time.Date(2006, time.January, 02, 15, 04, 05, 999999999, time.UTC),
				},
			},
			wantErr: false,
		},
		{
			name: "All variable types - unset but in default",
			args: args{
				configPath: "",
				entries: []Entry{
					{
						Name:      "VAL_STRING",
						EntryType: new(string),
						Default:   "ValStringContent",
						Required:  true,
					},
					{
						Name:      "VAL_INT",
						EntryType: new(int),
						Default:   "",
						Required:  true,
					},
					{
						Name:      "VAL_DURATION",
						EntryType: new(time.Duration),
						Default:   "",
						Required:  true,
					},
					{
						Name:      "VAL_TIME",
						EntryType: new(time.Time),
						Default:   "",
						Required:  true,
					},
				},
				target: AllTypes{},
				configVarsString: map[string]string{
					"VAL_INT":      "2",
					"VAL_DURATION": "10s",
					"VAL_TIME":     "2006-01-02T15:04:05.999999999",
				},
				want: AllTypes{
					ValString:   "ValStringContent",
					ValInt:      2,
					ValDuration: 10 * time.Second,
					ValTime:     time.Date(2006, time.January, 02, 15, 04, 05, 999999999, time.UTC),
				},
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				for k := range tt.args.configVarsString {
					_ = os.Unsetenv(k)
				}
			})

			for k, v := range tt.args.configVarsString {
				_ = os.Setenv(k, v)
			}

			c := NewReader(tt.args.configPath)

			err := c.LoadConfig(tt.args.entries, &tt.args.target)

			if tt.wantErr && err != nil {
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("----- TEST %s: LoadConfig() error = \"%v\", wantErr %v\n", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.target, tt.args.want) {
				t.Errorf("----- TEST %s: failed LoadConfig() \nwant:\n%+v\ngot:\n%+v\n", tt.name, tt.args.want, tt.args.target)
			}

		})
	}
}
