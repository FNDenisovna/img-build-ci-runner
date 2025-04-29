package storage

import (
	model "img-build-ci-runner/internal/model"
	"img-build-ci-runner/internal/storage/sqllite"
	"testing"
	"time"
)

func TestGetPackages(t *testing.T) {
	type args struct {
		branch string
		limit  int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Ok1",
			args: args{
				branch: "sisyphus", limit: 5,
			},
			wantErr: false,
		},
		{
			name: "Ok2",
			args: args{
				branch: "sisyphus", limit: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := sqllite.New()
			if (err != nil) != tt.wantErr {
				t.Errorf("sqllite.New return error during creating db = %v, wantErr %v", err, tt.wantErr)
				return
			}

			c := New(db)
			defer c.Close()
			packs, err := c.GetPackages(tt.args.branch, tt.args.limit)

			if len(packs) > tt.args.limit {
				t.Errorf("GetPackages() return %v packs, limit was %v", len(packs), tt.args.limit)
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPackages() return error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetPackage(t *testing.T) {
	type args struct {
		name    string
		release string
		version string
		changed time.Time
		branch  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Ok1",
			args: args{
				name:    "test",
				version: "1.0.0",
				release: "alt1",
				branch:  "sisyphus",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := sqllite.New()
			if (err != nil) != tt.wantErr {
				t.Errorf("sqllite.New return error during creating db = %v, wantErr %v", err, tt.wantErr)
				return
			}

			c := New(db)
			defer c.Close()

			pack := &model.SqlPack{
				Name:    tt.args.name,
				Version: tt.args.version,
				Release: tt.args.release,
				Branch:  tt.args.branch,
				Changed: time.Now(),
			}
			id_res, err := c.InsertPackage(pack)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertPackage() return error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			pack_res, err := c.GetPackage(pack.Name, tt.args.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPackages() return error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if id_res != pack_res.Id {
				t.Errorf("InsertPackage() and GetPackages() return different results")
				return
			}

			if err := c.DeletePackageById(pack_res.Id); err != nil {
				t.Errorf("DeletePackageById() can't delete info by id %v from GetPackages(). Error: %v", pack_res.Id, err)
				return
			}
		})
	}
}

func TestInsertPackage(t *testing.T) {
	type args struct {
		name    string
		release string
		version string
		changed time.Time
		branch  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Ok1",
			args: args{
				name:    "test",
				release: "alt1",
				version: "1.0.0",
				branch:  "sisyphus",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := sqllite.New()
			if (err != nil) != tt.wantErr {
				t.Errorf("sqllite.New return error during creating db: %v, wantErr %v", err, tt.wantErr)
				return
			}

			c := New(db)
			defer c.Close()
			pack := &model.SqlPack{
				Name:    tt.args.name,
				Version: tt.args.version,
				Release: tt.args.release,
				Branch:  tt.args.branch,
				Changed: time.Now(),
			}
			_, err = c.InsertPackage(pack)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertPackage() return error: %v, wantErr %v", err, tt.wantErr)
				return
			}

			if id, err := c.packExists(tt.args.name, tt.args.branch); err == nil && id > 0 {
				err = c.DeletePackageById(id)
				if err != nil {
					t.Errorf("Deleting package error: %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else {
				t.Errorf("Can't find package in db after inserting: %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUpdatePackage(t *testing.T) {
	type args struct {
		name    string
		release string
		version string
		changed time.Time
		branch  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Ok1",
			args: args{
				name:    "test",
				release: "alt1",
				version: "1.0.0",
				branch:  "sisyphus",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := sqllite.New()
			if (err != nil) != tt.wantErr {
				t.Errorf("sqllite.New return error during creating db: %v, wantErr %v", err, tt.wantErr)
				return
			}

			c := New(db)
			defer c.Close()
			pack := &model.SqlPack{
				Name:    tt.args.name,
				Version: tt.args.version,
				Release: tt.args.release,
				Branch:  tt.args.branch,
				Changed: time.Now(),
			}
			_, err = c.InsertPackage(pack)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertPackage() return error: %v, wantErr %v", err, tt.wantErr)
				return
			}

			if id, err := c.packExists(tt.args.name, tt.args.branch); err == nil && id > 0 {
				pack.Version = "2.0.0"
				id, err = c.UpdatePackage(pack, id)
				if err != nil {
					t.Errorf("Updating package error: %v, wantErr %v", err, tt.wantErr)
					return
				}

				err = c.DeletePackageById(id)
				if err != nil {
					t.Errorf("Deleting package error: %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else {
				t.Errorf("Can't find package in db after inserting: %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
