package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	Dev  = 0
	Prod = 1
)

type Reader struct {
	v          *viper.Viper
	configPath string
}

type Entry struct {
	Name      string
	EntryType interface{}
	// this value will be parsed to EntryType, panics if cannot be parsed
	Default  string
	Required bool
}

func NewReader(path string) *Reader {
	return &Reader{
		v:          viper.New(),
		configPath: path,
	}
}

// Loads config using entries slice as a reference.
// Currently only string, int, duration and time types are supported.
// Time must be passed in RFC3339Nano format (2006-01-02T15:04:05.999999999)
func (r Reader) LoadConfig(entries []Entry, target interface{}) error {
	var readFromFile bool
	if len(r.configPath) > 0 {
		readFromFile = true
	}

	if readFromFile {
		if !CheckDirExists(r.configPath) {
			return fmt.Errorf("config dir %s does not exists or cannot be read, "+
				"please provide empty path if this is intended (confguration loaded from env)", r.configPath)
		}
		configBase, configFile := path.Split(r.configPath)
		configFile = strings.TrimSuffix(configFile, filepath.Ext(configFile))
		r.v.AddConfigPath(configBase)
		r.v.SetConfigName(configFile)
	}
	r.v.SetConfigType("env")

	// TODO not reading config values not set in configfiles (form defaults)
	if err := r.setDefaults(entries); err != nil {
		return err
	}
	r.v.AutomaticEnv()

	err := r.v.ReadInConfig()
	if err != nil && len(r.configPath) > 0 {
		return fmt.Errorf("reading config from file failed %w", err)
	}
	if len(r.configPath) == 0 && err != err.(viper.ConfigFileNotFoundError) {
		return fmt.Errorf("reading config from env failed %w", err)
	}

	return r.parseEntries(entries, target)
}

func (r Reader) parseEntries(entries []Entry, target interface{}) error {
	// iterate over config entries again to verify required values and set struct fields

	for _, v := range entries {
		if reflect.ValueOf(v).Kind() == reflect.Ptr {
			return fmt.Errorf("cofnig entries types must be pointers, pass it with for example `new(string)`")
		}

		fieldName, err := CamelCase(v.Name)

		if err != nil {
			return fmt.Errorf("reading config failed, invalid config entry name, error: %w", err)
		}

		if v.Required && !r.v.IsSet(v.Name) {
			return fmt.Errorf("variable %s is required but is not set", v.Name)
		}

		var val interface{}
		switch v.EntryType.(type) {
		case *string:
			val = r.v.GetString(v.Name)
		case *int:
			val = r.v.GetInt(v.Name)
		case *time.Duration:
			val = r.v.GetDuration(v.Name)
		case *time.Time:
			val = r.v.GetTime(v.Name)
		}
		setField(target, fieldName, val)
	}
	return nil
}

func (r Reader) setDefaults(entries []Entry) error {
	// Iterate over config entries to set default values
	for _, v := range entries {
		if len(v.Default) > 0 {
			//var parsed interface{}

			var parsingErr error
			switch v.EntryType.(type) {
			case *string:
				parsed := v.Default
				//viper.SetDefault(v.Name, v.Default)
				r.v.SetDefault(v.Name, parsed)
			case *int:
				parsed, pErr := strconv.Atoi(v.Default)
				r.v.SetDefault(v.Name, parsed)
				parsingErr = pErr
			case *time.Duration:
				parsed, pErr := time.ParseDuration(v.Default)
				r.v.SetDefault(v.Name, parsed)
				parsingErr = pErr
			case *time.Time:
				parsed, pErr := time.Parse(time.RFC3339Nano, v.Default)
				r.v.SetDefault(v.Name, parsed)
				parsingErr = pErr

			}
			if parsingErr != nil {
				return parsingErr
			}
		}
	}
	return nil
}

// if field not exists - do nothing
func setField(target interface{}, fieldName string, value interface{}) {
	v := reflect.ValueOf(target).Elem()

	if target == nil {
		return
	}

	if !v.CanAddr() {
		panic(fmt.Sprintf("cannot assign to the target passed, target must be a pointer in order to assign"))
	}
	fieldNames := map[string]int{}
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		fn := typeField.Name
		fieldNames[fn] = i
	}

	fieldNum, _ := fieldNames[fieldName]
	fieldVal := v.Field(fieldNum)
	fieldVal.Set(reflect.ValueOf(value))
}

func CheckDirExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
