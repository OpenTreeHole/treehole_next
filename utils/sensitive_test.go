package utils

import (
	"testing"

	"github.com/importcjj/sensitive"
	"github.com/stretchr/testify/assert"

	"treehole_next/config"
	"treehole_next/data"
)

func TestIsSensitive(t *testing.T) {
	data.SensitiveWordFilter = sensitive.New()
	data.SensitiveWordFilter.AddWord("天气", "小明", "123")
	config.Config.OpenSensitiveCheck = true
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{content: "今天天气真好"},
			want: true,
		},
		{
			name: "test2",
			args: args{content: "昨天吃了什么"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, IsSensitive(tt.args.content), "IsSensitive(%v)", tt.args.content)
		})
	}
}
