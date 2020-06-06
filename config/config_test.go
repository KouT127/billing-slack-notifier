package config

import "testing"

func TestMustGetEnv(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "should get env",
			args: args{
				"TEST_ENV",
			},
			want: "TEST",
			wantErr: false,
		},
		{
			name: "should not get env",
			args: args{
				"TEST_ENV1",
			},
			want: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MustGetEnv(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("MustGetEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MustGetEnv() got = %v, want %v", got, tt.want)
			}
		})
	}
}
