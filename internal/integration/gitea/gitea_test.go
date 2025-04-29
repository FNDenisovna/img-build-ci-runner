package gitea

import (
	model "img-build-ci-runner/internal/model"
	"testing"
)

// /test for
// /func (g *GiteaApi) RunBuildImage(tag *model.GiteaTag, token string) error
func TestGitea_RunBuildImage(t *testing.T) {
	type args struct {
		branch  string
		version string
		image   string
		target  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "noerr",
			args:    args{branch: "p10", image: "incus", version: "6.0.0", target: "master"},
			wantErr: false,
		},
		{
			name:    "targeterr",
			args:    args{branch: "p10", image: "incus", version: "6.0.0", target: "fix_master"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New("")
			tag := &model.GiteaTag{
				Message: "building",
				Target:  tt.args.target,
				Image:   tt.args.image,
				Version: tt.args.version,
				Branch:  tt.args.branch,
			}
			token := "..."
			err := a.RunBuildImage(tag, token)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTag() return error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
