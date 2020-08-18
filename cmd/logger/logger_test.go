package logger

import (
	"reflect"
	"testing"
)

func TestNewLogger(t *testing.T) {
	type args struct {
		options   []Option
		skipDepth int
	}
	tests := []struct {
		name    string
		args    args
		want    *Logger
		wantErr bool
	}{
		{
			name: "format不合法",
			args: args{
				options: []Option{
					{OutputFormat: 0},
				},
				skipDepth: 0,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "writeToType不合法",
			args: args{
				options: []Option{
					{
						OutputFormat: FormatJSON,
						WriteTo: WriteTo{
							Type: 5,
							Path: "",
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "writeToFile缺少路径",
			args: args{
				options: []Option{
					{
						OutputFormat: FormatJSON,
						WriteTo: WriteTo{
							Type: WriteToFile,
							Path: "",
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "level不合法",
			args: args{
				options: []Option{
					{
						OutputFormat: FormatJSON,
						WriteTo: WriteTo{
							Type: WriteToStdout,
						},
						MinLevel: -2,
						MaxLevel: 0,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "skipDepth不合法",
			args: args{
				options: []Option{
					{
						OutputFormat: FormatJSON,
						WriteTo: WriteTo{
							Type: WriteToStdout,
						},
						MinLevel: DebugLevel,
						MaxLevel: ErrorLevel,
					},
				},
				skipDepth: -1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLogger(tt.args.options, tt.args.skipDepth)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLogger() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogger_Info(t *testing.T) {
	type args struct {
		args []interface{}
	}
	options := []Option{
		{
			OutputFormat: FormatPlainText,
			WriteTo: WriteTo{
				Type: WriteToStdout,
				Path: "",
			},
			MinLevel: InfoLevel,
			MaxLevel: WarnLevel,
		},
		{
			OutputFormat: FormatPlainText,
			WriteTo: WriteTo{
				Type: WriteToFile,
				Path: "/var/log/test/access.log",
			},
			MinLevel: InfoLevel,
			MaxLevel: WarnLevel,
		},
		{
			OutputFormat: FormatJSON,
			WriteTo: WriteTo{
				Type: WriteToStderr,
				Path: "",
			},
			MinLevel: ErrorLevel,
			MaxLevel: FatalLevel,
		},
		{
			OutputFormat: FormatJSON,
			WriteTo: WriteTo{
				Type: WriteToFile,
				Path: "/var/log/test/err.log",
			},
			MinLevel: ErrorLevel,
			MaxLevel: FatalLevel,
		},
		{
			OutputFormat: FormatPlainText,
			WriteTo: WriteTo{
				Type: WriteToQYWeiXinBot,
				Path: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=",
			},
			MinLevel: ErrorLevel,
			MaxLevel: FatalLevel,
		},
	}
	logger, err := NewLogger(options, 0)
	if err != nil {
		t.Error(err)
		return
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "人肉判断",
			args: args{args: []interface{}{"1", "2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.Info(tt.args)
			logger.Error(tt.args)
		})
	}
}
