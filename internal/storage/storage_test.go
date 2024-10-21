package storage

import (
	"altpack-vers-checker/internal/storage/sqllite"
	"testing"
	"time"
)

func TestGetPackages(t *testing.T) {
	type args struct {
		offset int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Ok1",
			args:    args{offset: 5},
			wantErr: false,
		},
		{
			name:    "Ok2",
			args:    args{offset: 1},
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
			_, err = c.GetPackages(tt.args.offset)
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
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Ok1",
			args:    args{name: "test"},
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

			pack := &SqlPack{
				Name:    tt.args.name,
				Version: tt.args.version,
				Release: tt.args.release,
				Changed: time.Now(),
			}
			id_res, err := c.InsertPackage(pack)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertPackage() return error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			pack_res, err := c.GetPackage(pack.Name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPackages() return error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if id_res != pack_res.Id {
				t.Errorf("InsertPackage() and GetPackages() return different results")
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
			pack := &SqlPack{
				Name:    tt.args.name,
				Version: tt.args.version,
				Release: tt.args.release,
				Changed: time.Now(),
			}
			_, err = c.InsertPackage(pack)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertPackage() return error: %v, wantErr %v", err, tt.wantErr)
				return
			}

			if id, err := c.packExists(tt.args.name); err == nil && id > 0 {
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
