package minerinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_whatsminer_Type(t *testing.T) {
	a := whatsminer{}
	got := a.Type()
	assert.Equal(t, "Whatsminer", got)
}

func Test_whatsminer_Model_Version(t *testing.T) {
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
			name: "Whatsminer M30S++VG30",
			fields: fields{
				info: minerInfo{
					Devdetails: []devdetails{{Model: "M30S++VG30"}},
				},
			},
			wantModel:   "M30S++",
			wantVersion: "VG30",
		},
		{
			name: "Whatsminer M31S+V40",
			fields: fields{
				info: minerInfo{
					Devdetails: []devdetails{{Model: "M31S+V40"}},
				},
			},
			wantModel:   "M31S+",
			wantVersion: "V40",
		},
		{
			name: "Whatsminer M31S++V90",
			fields: fields{
				info: minerInfo{
					Devdetails: []devdetails{{Model: "M31S++V90"}},
				},
			},
			wantModel:   "M31S++",
			wantVersion: "V90",
		},
		{
			name: "Whatsminer M31PV90",
			fields: fields{
				info: minerInfo{
					Devdetails: []devdetails{{Model: "M31PV90"}},
				},
			},
			wantModel:   "M31P",
			wantVersion: "V90",
		},
		{
			name: "Whatsminer M31PV9V$V##V0",
			fields: fields{
				info: minerInfo{
					Devdetails: []devdetails{{Model: "M31PV9V$V##V0"}},
				},
			},
			wantModel:   "M31P",
			wantVersion: "V9V$V##V0",
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
			a := whatsminer{
				genericMiner: genericMiner{tt.fields.info},
			}
			gotModel := a.Model()
			assert.Equal(t, tt.wantModel, gotModel)
			gotVersion := a.Version()
			assert.Equal(t, tt.wantVersion, gotVersion)
		})
	}
}

func Test_whatsminer_Version(t *testing.T) {
	w := whatsminer{}
	got := w.Version()
	assert.Equal(t, "", got)
}
