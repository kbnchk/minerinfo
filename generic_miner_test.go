package minerinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_genericMiner_Hashrate(t *testing.T) {
	type fields struct {
		data minerInfo
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "MHS float",
			fields: fields{
				data: minerInfo{
					Summary: []summary{{MHS: 123456.78}},
				},
			},
			want: float64(123456.78),
		},
		{
			name: "MHS text",
			fields: fields{
				data: minerInfo{
					Summary: []summary{{MHS: "123456.78"}},
				},
			},
			want: float64(123456.78),
		},
		{
			name: "GHs text",
			fields: fields{
				data: minerInfo{
					Summary: []summary{{GHS: "123456.78"}},
				},
			},
			want: float64(123456780),
		},
		{
			name: "Empty values",
			fields: fields{
				data: minerInfo{},
			},
			want: float64(0),
		},
		{
			name: "Wrong values",
			fields: fields{
				data: minerInfo{
					Summary: []summary{{GHS: false}},
				},
			},
			want: float64(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := genericMiner{
				data: tt.fields.data,
			}
			got := g.Hashrate()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_genericMiner_Pools(t *testing.T) {
	type fields struct {
		data minerInfo
	}
	tests := []struct {
		name   string
		fields fields
		want   []Pools
	}{
		{
			name: "default",
			fields: fields{
				data: minerInfo{
					Pools: []Pools{
						{
							Url:    "testPool 1",
							User:   "testUser 1",
							Status: "Alive",
						},
						{
							Url:    "testPool 2",
							User:   "testUser 2",
							Status: "Dead",
						},
					},
				},
			},
			want: []Pools{
				{
					Url:    "testPool 1",
					User:   "testUser 1",
					Status: "Alive",
				},
				{
					Url:    "testPool 2",
					User:   "testUser 2",
					Status: "Dead",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := genericMiner{
				data: tt.fields.data,
			}
			got := g.Pools()
			assert.Equal(t, tt.want, got)
		})
	}
}
