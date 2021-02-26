package config

import "testing"

func TestCamelCase1(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "TEST_TEST",args: args{str: "TEST_TEST"},want: "TestTest",wantErr: false},
		{name: "TEST-TEST",args: args{str: "TEST-TEST"},want: "",wantErr: true},
		{name: "test_Test",args: args{str: "test_Test"},want: "TestTest",wantErr: false},
		{name: "test",args: args{str: "test"},want: "Test",wantErr: false},
		{name: "test_Test_TEST",args: args{str: "test_Test_TEST"},want: "TestTestTest",wantErr: false},
		{name: "Test",args: args{str: "Test"},want: "Test",wantErr: false},
		{name: "Test123",args: args{str: "Test123"},want: "Test123",wantErr: false},
		{name: "123", args: args{str: "123"}, want: "123", wantErr: false},
		{name: "123", args: args{str: "123"}, want: "123", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CamelCase(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("CamelCase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CamelCase() got = %v, want %v", got, tt.want)
			}
		})
	}
}