package minerinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_antminer_Type(t *testing.T) {
	a := antminer{}
	got := a.Type()
	assert.Equal(t, "Antminer", got)
}

func Test_antminer_Model_Version(t *testing.T) {
	type fields struct {
		info minerInfo
	}
	tests := []struct {
		name        string
		fields      fields
		wantModel   string
		wantVersion string
	}{
		{
			name: "Antminer L7",
			fields: fields{
				info: minerInfo{
					Stats: []stats{{Type: "Antminer L7"}},
				},
			},
			wantModel:   "L7",
			wantVersion: "",
		},
		{
			name: "Antminer L7 (Vnish 1.2.0-rc4)",
			fields: fields{
				info: minerInfo{
					Stats: []stats{{Type: "Antminer L7 (Vnish 1.2.0-rc4)"}},
				},
			},
			wantModel:   "L7",
			wantVersion: "Vnish 1.2.0-rc4",
		},
		{
			name: "Antminer S19 XP",
			fields: fields{
				info: minerInfo{
					Stats: []stats{{Type: "Antminer S19 XP"}},
				},
			},
			wantModel: "S19 XP",
		},
		{
			name: "Antminer S19 XP (Vnish 1.2.0-rc1)",
			fields: fields{
				info: minerInfo{
					Stats: []stats{{Type: "Antminer S19 XP (Vnish 1.2.0-rc1)"}},
				},
			},
			wantModel:   "S19 XP",
			wantVersion: "Vnish 1.2.0-rc1",
		},
		{
			name: "Antminer S19 XP (33432947623()()((())))",
			fields: fields{
				info: minerInfo{
					Stats: []stats{{Type: "Antminer S19 XP (33432947623()()((())))"}},
				},
			},
			wantModel:   "S19 XP",
			wantVersion: "33432947623()()((()))",
		},
		{
			name: "unknown model",
			fields: fields{
				info: minerInfo{
					Devdetails: []devdetails{},
				},
			},
			wantModel: "Unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := antminer{
				genericMiner: genericMiner{tt.fields.info},
			}
			gotModel := a.Model()
			assert.Equal(t, tt.wantModel, gotModel)
			gotVersion := a.Version()
			assert.Equal(t, tt.wantVersion, gotVersion)
		})
	}
}

func Test_antminer_Hashrate(t *testing.T) {
	type fields struct {
		info minerInfo
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "Antminer L7",
			fields: fields{
				info: minerInfo{
					Stats:   []stats{{Type: "Antminer L7"}},
					Summary: []summary{{GHS: 9807.629}},
				},
			},
			want: float64(9807.629),
		},
		{
			name: "Antminer S19 XP (Vnish 1.2.0-rc1)",
			fields: fields{
				info: minerInfo{
					Stats:   []stats{{Type: "Antminer S19 XP"}},
					Summary: []summary{{GHS: 146776.72}},
				},
			},
			want: float64(146776720.00),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := antminer{
				genericMiner: genericMiner{tt.fields.info},
			}
			got := a.Hashrate()
			assert.Equal(t, tt.want, got)
		})
	}
}
