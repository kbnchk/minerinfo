package minerinfo

import (
	"bufio"
	"bytes"
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
			got := detectMiner(tt.args.info)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_sendCommand(t *testing.T) {
	type args struct {
		req        minerRequest
		respString string
	}
	tests := []struct {
		name     string
		args     args
		wantResp minerInfo
		wanterr  bool
	}{
		{
			name: "S19 XP stats",
			args: args{
				req: minerRequest{
					Command: "stats",
				},
				respString: `{"STATUS": [{"STATUS":"S","When":1696258519,"Code":70,"Msg":"CGMiner stats","Description":"cgminer 4.11.1"}],"STATS": [{"Cgminer":"4.11.1","Miner":"uart_trans.1.3","CompileTime":"Sat Jun 24 11:17:59 UTC 2023","Type":"Antminer S19 XP (Vnish 1.2.0-rc1)"},{"STATS":0,"ID":"BTM_S19XP","Elapsed":331311,"Calls":0,"Wait":0.000000,"Max":0.000000,"Min":99999999.000000,"GHS 5s":"140258.58","GHS av":142562.48,"miner_count":3,"total_acn":330,"frequency":495,"total_freqavg":489.46,"total_rateideal":144399.77,"total_rate":140258.58,"state":"mining","fan_num":4,"fan_mode":0,"fan_pwm":25,"fan1":2790,"fan2":2850,"fan3":2940,"fan4":2910,"temp_num":3,"temp1":60,"temp2":57,"temp3":59,"temp2_1":75,"temp2_2":72,"temp2_3":74,"temp3_1":75,"temp3_2":72,"temp3_3":74,"temp_sens1":"m-m","temp_sens2":"m-m","temp_sens3":"m-m","temp_pcb1":"37-60","temp_pcb2":"35-57","temp_pcb3":"37-59","temp_chip1":"52-75","temp_chip2":"50-72","temp_chip3":"52-74","chain_state1":"mining","chain_state2":"mining","chain_state3":"mining","chain_acn1":110,"chain_acn2":110,"chain_acn3":110,"chain_vol1":12625,"chain_vol2":12625,"chain_vol3":12625,"freq_avg1":490.82,"freq_avg2":487.43,"freq_avg3":490.13,"chain_rateideal1":48267.06,"chain_rateideal2":47933.60,"chain_rateideal3":48199.12,"chain_rate1":47016.68,"chain_rate2":47262.92,"chain_rate3":45978.98,"chain_hw1":0,"chain_hw2":7994,"chain_hw3":0,"chain_acs1":"oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo","chain_acs2":"oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo","chain_acs3":"oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo oooooooooo","chain_consumption1":0,"chain_consumption2":0,"chain_consumption3":0,"miner_version":"uart_trans.1.3","build_version":"1.2.0-rc1"}],"id":1}`,
			},
			wantResp: minerInfo{
				Stats: []stats{
					{
						Type: "Antminer S19 XP (Vnish 1.2.0-rc1)",
					},
					{},
				},
			},
			wanterr: false,
		},
		{
			name: "S19 XP summary",
			args: args{
				req: minerRequest{
					Command: "summary",
				},
				respString: `{"STATUS": [{"STATUS":"S","When":1696258988,"Code":11,"Msg":"Summary","Description":"cgminer 4.11.1"}],"SUMMARY": [{"Elapsed":331780,"GHS 5s":"140258.58","GHS av":142562.73,"Found Blocks":0,"Getworks":14434,"Accepted":226146,"Rejected":94,"Hardware Errors":8004,"Utility":40.90,"Discarded":175963,"Stale":177,"Get Failures":17,"Local Work":86289179,"Remote Failures":0,"Network Blocks":549,"Total MH":47299461638452.0000,"Work Utility":1991845.62,"Difficulty Accepted":11019960320.00000000,"Difficulty Rejected":7772160.00000000,"Difficulty Stale":681984.00000000,"Best Share":11482220369,"Device Hardware%":0.5332,"Device Rejected%":0.0706,"Pool Rejected%":0.0705,"Pool Stale%":0.0062,"Last getwork":1696258987,"Fee Percent":3.8449}],"id":1}`,
			},
			wantResp: minerInfo{
				Summary: []summary{
					{
						GHS: "140258.58",
					},
				},
			},
			wanterr: false,
		},
		{
			name: "S19 XP pools",
			args: args{
				req: minerRequest{
					Command: "pools",
				},
				respString: `{"STATUS": [{"STATUS":"S","When":1696260268,"Code":7,"Msg":"6 Pool(s)","Description":"cgminer 4.11.1"}],"POOLS": [{"POOL":0,"URL":"stratum+tcp://eu.ss.btc.com:1800","Status":"Alive","Priority":0,"Long Poll":"N","Getworks":11871,"Accepted":19335,"Rejected":12,"Stale":54,"Works":83106431,"Discarded":169184,"Get Failures":14,"Remote Failures":0,"User":"divcr1","Type":0,"Last Share Time":"0:00:27","Diff":"524K","Diff1 Shares":0,"Proxy Type":"","Proxy":"","Difficulty Accepted":10636623872.00000000,"Difficulty Rejected":8126464.00000000,"Difficulty Stale":524288.00000000,"Last Share Difficulty":524288.00000000,"Work Difficulty":524288.00000000,"Has Stratum":true,"Stratum Active":true,"Stratum URL":"eu.ss.btc.com","Stratum Difficulty":524288.00000000,"Has Vmask":true,"Best Share":11482220369,"Pool Rejected%":0.0763,"Pool Stale%":0.0049,"Bad Work":0,"Current Block Height":810322,"Current Block Version":536870912},{"POOL":1,"URL":"stratum+tcp://eu.ss.btc.com:443","Status":"Alive","Priority":1,"Long Poll":"N","Getworks":0,"Accepted":0,"Rejected":0,"Stale":0,"Works":0,"Discarded":0,"Get Failures":0,"Remote Failures":0,"User":"divcr1","Type":0,"Last Share Time":"0","Diff":"131K","Diff1 Shares":0,"Proxy Type":"","Proxy":"","Difficulty Accepted":0.00000000,"Difficulty Rejected":0.00000000,"Difficulty Stale":0.00000000,"Last Share Difficulty":0.00000000,"Work Difficulty":131072.00000000,"Has Stratum":true,"Stratum Active":false,"Stratum URL":"","Stratum Difficulty":0.00000000,"Has Vmask":true,"Best Share":0,"Pool Rejected%":0.0000,"Pool Stale%":0.0000,"Bad Work":0,"Current Block Height":809772,"Current Block Version":536870912},{"POOL":2,"URL":"stratum+tcp://eu.ss.btc.com:25","Status":"Dead","Priority":2,"Long Poll":"N","Getworks":0,"Accepted":0,"Rejected":0,"Stale":0,"Works":0,"Discarded":0,"Get Failures":0,"Remote Failures":0,"User":"divcr1","Type":0,"Last Share Time":"0","Diff":"","Diff1 Shares":0,"Proxy Type":"","Proxy":"","Difficulty Accepted":0.00000000,"Difficulty Rejected":0.00000000,"Difficulty Stale":0.00000000,"Last Share Difficulty":0.00000000,"Work Difficulty":0.00000000,"Has Stratum":true,"Stratum Active":false,"Stratum URL":"","Stratum Difficulty":0.00000000,"Has Vmask":false,"Best Share":0,"Pool Rejected%":0.0000,"Pool Stale%":0.0000,"Bad Work":0,"Current Block Height":0,"Current Block Version":0},{"POOL":3,"URL":"DevFee","Status":"Alive","Priority":3,"Long Poll":"N","Getworks":2617,"Accepted":207793,"Rejected":83,"Stale":123,"Works":3337415,"Discarded":7459,"Get Failures":3,"Remote Failures":0,"User":"DevFee","Type":1,"Last Share Time":"0:06:36","Diff":"2.05K","Diff1 Shares":0,"Proxy Type":"","Proxy":"","Difficulty Accepted":425560064.00000000,"Difficulty Rejected":169984.00000000,"Difficulty Stale":157696.00000000,"Last Share Difficulty":2048.00000000,"Work Difficulty":2048.00000000,"Has Stratum":true,"Stratum Active":false,"Stratum URL":"","Stratum Difficulty":0.00000000,"Has Vmask":true,"Best Share":397281270,"Pool Rejected%":0.0399,"Pool Stale%":0.0370,"Bad Work":0,"Current Block Height":810321,"Current Block Version":536870912}],"id":1}`,
			},
			wantResp: minerInfo{
				Pools: []Pools{
					{
						Url:    "stratum+tcp://eu.ss.btc.com:1800",
						User:   "divcr1",
						Status: "Alive",
					},
					{
						Url:    "stratum+tcp://eu.ss.btc.com:443",
						User:   "divcr1",
						Status: "Alive",
					},
					{
						Url:    "stratum+tcp://eu.ss.btc.com:25",
						User:   "divcr1",
						Status: "Dead",
					},
					{
						Url:    "DevFee",
						User:   "DevFee",
						Status: "Alive",
					},
				},
			},
			wanterr: false,
		},
		{
			name: "Cheetah devdetails",
			args: args{
				req: minerRequest{
					Command: "pools",
				},
				respString: `{"STATUS":[{"STATUS":"S","When":1696262507,"Code":69,"Msg":"Device Details","Description":"cgminer 4.10.0"}],"DEVDETAILS":[{"DEVDETAILS":0,"Name":"C3012","ID":0,"Driver":"C3012","Kernel":"","Model":"0","Device Path":""},{"DEVDETAILS":1,"Name":"C3012","ID":0,"Driver":"C3012","Kernel":"","Model":"1","Device Path":""},{"DEVDETAILS":2,"Name":"C3012","ID":0,"Driver":"C3012","Kernel":"","Model":"2","Device Path":""}],"id":1}`,
			},
			wantResp: minerInfo{
				Devdetails: []devdetails{

					{
						Name:  "C3012",
						Model: "0",
					},
					{
						Name:  "C3012",
						Model: "1",
					},
					{
						Name:  "C3012",
						Model: "2",
					},
				},
			},
			wanterr: false,
		},
		{
			name: "invalid response",
			args: args{
				req: minerRequest{
					Command: "summary",
				},
				respString: "asjdasldjakldjakle34124dqdqf",
			},
			wantResp: minerInfo{},
			wanterr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			connreader := bufio.NewReader(bytes.NewBufferString(tt.args.respString))
			connwriter := bufio.NewWriter(bytes.NewBuffer(nil))
			conn := bufio.NewReadWriter(connreader, connwriter)

			var got minerInfo
			err := sendCommand(conn, tt.args.req, &got)
			if tt.wanterr == true {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResp, got)
			}
		})
	}
}
