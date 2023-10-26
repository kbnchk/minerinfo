package minerinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_detectMiner(t *testing.T) {
	type args struct {
		info minerInfo
	}
	tests := []struct {
		name string
		args args
		want Miner
	}{
		{
			name: "antminer",
			args: args{
				info: minerInfo{
					Stats: []stats{
						{Type: "Antminer L7 (Vnish 1.2.0-rc4)"},
					},
				},
			},
			want: antminer{genericMiner{data: minerInfo{
				Stats: []stats{
					{Type: "Antminer L7 (Vnish 1.2.0-rc4)"},
				},
			},
			},
			},
		},
		{
			name: "whatsminer",
			args: args{
				info: minerInfo{
					Devdetails: []devdetails{
						{Name: "SM"},
					},
				},
			},
			want: whatsminer{genericMiner{data: minerInfo{
				Devdetails: []devdetails{
					{Name: "SM"},
				},
			},
			},
			},
		},
		{
			name: "Unknown 1",
			args: args{
				info: minerInfo{},
			},
			want: genericMiner{data: minerInfo{}},
		},
		{
			name: "Unknown 2",
			args: args{
				info: minerInfo{
					Stats: []stats{
						{Type: "Ololo 12"},
					},
				},
			},
			want: genericMiner{data: minerInfo{
				Stats: []stats{
					{Type: "Ololo 12"},
				},
			},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := host{info: tt.args.info}
			got := h.detectMiner()
			assert.Equal(t, tt.want, got)
		})
	}
}
