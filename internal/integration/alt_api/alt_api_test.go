package alt_api

import (
	config "img-build-ci-runner/internal/config/viper"

	"testing"
)

// /test for
// /func GetSitePackInfo(name, branch string) (packinfo SiteVersion, err error)
func TestAltapi_GetSitePackInfo(t *testing.T) {
	type args struct {
		branch string
		name   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "noerr_test",
			args:    args{branch: "sisyphus", name: "incus"},
			wantErr: false,
		},
		{
			name:    "noerr_test",
			args:    args{branch: "sisyphus", name: "etcd"},
			wantErr: false,
		},
		{
			name:    "noerr_test",
			args:    args{branch: "sisyphus", name: "etcd"},
			wantErr: false,
		},
		{
			name:    "noerr_test",
			args:    args{branch: "p10", name: "git"},
			wantErr: false,
		},
		{
			name:    "err_test",
			args:    args{branch: "p10", name: "incus"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.New()
			pack, err := a.GetSitePackInfo(tt.args.name, tt.args.branch)
			if (err != nil) && !tt.wantErr {
				t.Errorf("GetSitePackInfo() return error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.args.name != pack.Name {
				t.Errorf("GetSitePackInfo() return info about package = %v, want = %v", pack.Name, tt.args.name)
			}

		})
	}
}
