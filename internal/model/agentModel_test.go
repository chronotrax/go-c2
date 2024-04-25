package model

import (
	"github.com/google/uuid"
	"net"
	"reflect"
	"testing"
)

func TestParseAgentModel(t *testing.T) {
	type args struct {
		id string
		ip string
	}
	tests := []struct {
		name    string
		args    args
		want    *Agent
		wantErr bool
	}{
		{name: "valid agent",
			args: args{
				id: "6834d229-3c16-45af-b492-ef26ca5d6770",
				ip: "192.168.1.1"},
			want: &Agent{
				ID: uuid.MustParse("6834d229-3c16-45af-b492-ef26ca5d6770"),
				IP: net.IPv4(192, 168, 1, 1),
			}},
		{name: "valid agent: IPv6",
			args: args{
				id: "6834d229-3c16-45af-b492-ef26ca5d6770",
				ip: "0102:0304:0506:0708:090a:0b0c:0d0e:0fff"},
			want: &Agent{
				ID: uuid.MustParse("6834d229-3c16-45af-b492-ef26ca5d6770"),
				IP: net.IP{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0xff},
			}},
		{name: "invalid ip",
			args: args{
				id: "6834d229-3c16-45af-b492-ef26ca5d6770",
				ip: "192..1"},
			want:    nil,
			wantErr: true},
		{name: "invalid ip: IPv6",
			args: args{
				id: "6834d229-3c16-45af-b492-ef26ca5d6770",
				ip: "0102:0fff"},
			want:    nil,
			wantErr: true},
		{name: "invalid uuid: too short",
			args: args{
				id: "6834d229-3c16-45af-b492-ef26ca5d677",
				ip: "192.168.1.1"},
			want:    nil,
			wantErr: true},
		{name: "invalid uuid: invalid char",
			args: args{
				id: "6834d229-3c16-45af-b492-ef26ca5d677y",
				ip: "192.168.1.1"},
			want:    nil,
			wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAgentModel(tt.args.id, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAgentModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAgentModel() got = %v, want %v", got, tt.want)
			}
		})
	}
}
