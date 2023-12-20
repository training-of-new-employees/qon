package impl

import (
	"context"
	"reflect"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/pkg/doar"
	"github.com/training-of-new-employees/qon/internal/pkg/jwttoken"
	"github.com/training-of-new-employees/qon/internal/store"
	"github.com/training-of-new-employees/qon/internal/store/cache"
	mock_store "github.com/training-of-new-employees/qon/mocks/store"
)

func Test_newUserService(t *testing.T) {
	type args struct {
		db         store.Storages
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		cache      cache.Cache
		jwtGen     jwttoken.JWTGenerator
		jwtVal     jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	tests := []struct {
		name string
		args args
		want *uService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newUserService(tt.args.db, tt.args.secretKey, tt.args.aTokenTime, tt.args.rTokenTime, tt.args.cache, tt.args.jwtGen, tt.args.jwtVal, tt.args.sender); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uService_WriteAdminToCache(t *testing.T) {
	type fields struct {
		db         store.Storages
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx context.Context
		val model.CreateAdmin
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.CreateAdmin
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &uService{
				db:         tt.fields.db,
				cache:      tt.fields.cache,
				secretKey:  tt.fields.secretKey,
				aTokenTime: tt.fields.aTokenTime,
				rTokenTime: tt.fields.rTokenTime,
				tokenGen:   tt.fields.tokenGen,
				tokenVal:   tt.fields.tokenVal,
				sender:     tt.fields.sender,
			}
			got, err := u.WriteAdminToCache(tt.args.ctx, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.WriteAdminToCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uService.WriteAdminToCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uService_GetUserByEmail(t *testing.T) {
	type fields struct {
		db         store.Storages
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &uService{
				db:         tt.fields.db,
				cache:      tt.fields.cache,
				secretKey:  tt.fields.secretKey,
				aTokenTime: tt.fields.aTokenTime,
				rTokenTime: tt.fields.rTokenTime,
				tokenGen:   tt.fields.tokenGen,
				tokenVal:   tt.fields.tokenVal,
				sender:     tt.fields.sender,
			}
			got, err := u.GetUserByEmail(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uService.GetUserByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uService_GetUserByID(t *testing.T) {
	type fields struct {
		db         store.Storages
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.UserInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &uService{
				db:         tt.fields.db,
				cache:      tt.fields.cache,
				secretKey:  tt.fields.secretKey,
				aTokenTime: tt.fields.aTokenTime,
				rTokenTime: tt.fields.rTokenTime,
				tokenGen:   tt.fields.tokenGen,
				tokenVal:   tt.fields.tokenVal,
				sender:     tt.fields.sender,
			}
			got, err := u.GetUserByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uService.GetUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uService_ArchiveUser(t *testing.T) {
	type fields struct {
		db         store.Storages
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx             context.Context
		id              int
		editorCompanyID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &uService{
				db:         tt.fields.db,
				cache:      tt.fields.cache,
				secretKey:  tt.fields.secretKey,
				aTokenTime: tt.fields.aTokenTime,
				rTokenTime: tt.fields.rTokenTime,
				tokenGen:   tt.fields.tokenGen,
				tokenVal:   tt.fields.tokenVal,
				sender:     tt.fields.sender,
			}
			if err := u.ArchiveUser(tt.args.ctx, tt.args.id, tt.args.editorCompanyID); (err != nil) != tt.wantErr {
				t.Errorf("uService.ArchiveUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_uService_EditUser(t *testing.T) {
	type fields struct {
		db         store.Storages
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx             context.Context
		val             *model.UserEdit
		editorCompanyID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.UserEdit
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &uService{
				db:         tt.fields.db,
				cache:      tt.fields.cache,
				secretKey:  tt.fields.secretKey,
				aTokenTime: tt.fields.aTokenTime,
				rTokenTime: tt.fields.rTokenTime,
				tokenGen:   tt.fields.tokenGen,
				tokenVal:   tt.fields.tokenVal,
				sender:     tt.fields.sender,
			}
			got, err := u.EditUser(tt.args.ctx, tt.args.val, tt.args.editorCompanyID)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.EditUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uService.EditUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uService_GetUsersByCompany(t *testing.T) {
	type fields struct {
		db         store.Storages
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx       context.Context
		companyID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &uService{
				db:         tt.fields.db,
				cache:      tt.fields.cache,
				secretKey:  tt.fields.secretKey,
				aTokenTime: tt.fields.aTokenTime,
				rTokenTime: tt.fields.rTokenTime,
				tokenGen:   tt.fields.tokenGen,
				tokenVal:   tt.fields.tokenVal,
				sender:     tt.fields.sender,
			}
			got, err := u.GetUsersByCompany(tt.args.ctx, tt.args.companyID)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.GetUsersByCompany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uService.GetUsersByCompany() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uService_GenerateTokenPair(t *testing.T) {
	type fields struct {
		db         store.Storages
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx       context.Context
		userId    int
		isAdmin   bool
		companyID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Tokens
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &uService{
				db:         tt.fields.db,
				cache:      tt.fields.cache,
				secretKey:  tt.fields.secretKey,
				aTokenTime: tt.fields.aTokenTime,
				rTokenTime: tt.fields.rTokenTime,
				tokenGen:   tt.fields.tokenGen,
				tokenVal:   tt.fields.tokenVal,
				sender:     tt.fields.sender,
			}
			got, err := u.GenerateTokenPair(tt.args.ctx, tt.args.userId, tt.args.isAdmin, tt.args.companyID)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.GenerateTokenPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uService.GenerateTokenPair() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uService_CreateUser(t *testing.T) {
	type fields struct {
		db         store.Storages
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx context.Context
		val model.UserCreate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &uService{
				db:         tt.fields.db,
				cache:      tt.fields.cache,
				secretKey:  tt.fields.secretKey,
				aTokenTime: tt.fields.aTokenTime,
				rTokenTime: tt.fields.rTokenTime,
				tokenGen:   tt.fields.tokenGen,
				tokenVal:   tt.fields.tokenVal,
				sender:     tt.fields.sender,
			}
			got, err := u.CreateUser(tt.args.ctx, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uService.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uService_GetAdminFromCache(t *testing.T) {
	type fields struct {
		db         store.Storages
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx  context.Context
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.CreateAdmin
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &uService{
				db:         tt.fields.db,
				cache:      tt.fields.cache,
				secretKey:  tt.fields.secretKey,
				aTokenTime: tt.fields.aTokenTime,
				rTokenTime: tt.fields.rTokenTime,
				tokenGen:   tt.fields.tokenGen,
				tokenVal:   tt.fields.tokenVal,
				sender:     tt.fields.sender,
			}
			got, err := u.GetAdminFromCache(tt.args.ctx, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.GetAdminFromCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uService.GetAdminFromCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uService_DeleteAdminFromCache(t *testing.T) {
	type fields struct {
		db         store.Storages
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &uService{
				db:         tt.fields.db,
				cache:      tt.fields.cache,
				secretKey:  tt.fields.secretKey,
				aTokenTime: tt.fields.aTokenTime,
				rTokenTime: tt.fields.rTokenTime,
				tokenGen:   tt.fields.tokenGen,
				tokenVal:   tt.fields.tokenVal,
				sender:     tt.fields.sender,
			}
			if err := u.DeleteAdminFromCache(tt.args.ctx, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("uService.DeleteAdminFromCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_uService_CreateAdmin(t *testing.T) {
	type fields struct {
		db         store.Storages
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx context.Context
		val *model.CreateAdmin
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &uService{
				db:         tt.fields.db,
				cache:      tt.fields.cache,
				secretKey:  tt.fields.secretKey,
				aTokenTime: tt.fields.aTokenTime,
				rTokenTime: tt.fields.rTokenTime,
				tokenGen:   tt.fields.tokenGen,
				tokenVal:   tt.fields.tokenVal,
				sender:     tt.fields.sender,
			}
			got, err := u.CreateAdmin(tt.args.ctx, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.CreateAdmin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uService.CreateAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uService_UpdatePasswordAndActivateUser(t *testing.T) {
	type fields struct {
		db         store.Storages
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx      context.Context
		email    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &uService{
				db:         tt.fields.db,
				cache:      tt.fields.cache,
				secretKey:  tt.fields.secretKey,
				aTokenTime: tt.fields.aTokenTime,
				rTokenTime: tt.fields.rTokenTime,
				tokenGen:   tt.fields.tokenGen,
				tokenVal:   tt.fields.tokenVal,
				sender:     tt.fields.sender,
			}
			if err := u.UpdatePasswordAndActivateUser(tt.args.ctx, tt.args.email, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("uService.UpdatePasswordAndActivateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_uService_ResetPassword(t *testing.T) {
	type fields struct {
		db         store.Storages
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &uService{
				db:         tt.fields.db,
				cache:      tt.fields.cache,
				secretKey:  tt.fields.secretKey,
				aTokenTime: tt.fields.aTokenTime,
				rTokenTime: tt.fields.rTokenTime,
				tokenGen:   tt.fields.tokenGen,
				tokenVal:   tt.fields.tokenVal,
				sender:     tt.fields.sender,
			}
			if err := u.ResetPassword(tt.args.ctx, tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("uService.ResetPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_uService_EditAdmin(t *testing.T) {
	type fields struct {
		userdb     *mock_store.MockRepositoryUser
		cache      cache.Cache
		secretKey  string
		aTokenTime time.Duration
		rTokenTime time.Duration
		tokenGen   jwttoken.JWTGenerator
		tokenVal   jwttoken.JWTValidator
		sender     doar.EmailSender
	}
	type args struct {
		ctx context.Context
		val model.AdminEdit
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		setargs func() args
		want    *model.AdminEdit
		wantErr bool
	}{
		{
			"Invalid Data",
			nil,
			func() args {
				email := "a"
				aedit := model.AdminEdit{
					ID:    1,
					Email: &email,
				}
				return args{
					nil,
					aedit,
				}
			},
			nil,
			true,
		},
		{
			"Valid data, db problem",
			func(f *fields) {
				email := "valid@mail.com"
				aedit := model.AdminEdit{
					ID:    1,
					Email: &email,
				}
				f.userdb.EXPECT().EditAdmin(gomock.Any(), aedit).Return(nil, errs.ErrNotFound)
			},
			func() args {
				email := "valid@mail.com"
				aedit := model.AdminEdit{
					ID:    1,
					Email: &email,
				}
				return args{
					nil,
					aedit,
				}
			},
			nil,
			true,
		},
		{
			"Valid Data",
			func(f *fields) {
				email := "valid@mail.com"
				aedit := model.AdminEdit{
					ID:    1,
					Email: &email,
				}
				f.userdb.EXPECT().EditAdmin(gomock.Any(), aedit).Return(&aedit, nil)

			},
			func() args {
				email := "valid@mail.com"
				aedit := model.AdminEdit{
					ID:    1,
					Email: &email,
				}
				return args{
					nil,
					aedit,
				}
			},
			&model.AdminEdit{
				ID: 1,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := &fields{}
			f.userdb = mock_store.NewMockRepositoryUser(ctrl)
			if tt.prepare != nil {
				tt.prepare(f)
			}
			storages := mockUserStorage(ctrl, f.userdb)
			u := &uService{
				db:         storages,
				cache:      f.cache,
				secretKey:  f.secretKey,
				aTokenTime: f.aTokenTime,
				rTokenTime: f.rTokenTime,
				tokenGen:   f.tokenGen,
				tokenVal:   f.tokenVal,
				sender:     f.sender,
			}
			args := tt.setargs()
			got, err := u.EditAdmin(args.ctx, args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("uService.EditAdmin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("uService.EditAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockUserStorage(ctrl *gomock.Controller, uStore *mock_store.MockRepositoryUser) *mock_store.MockStorages {
	posStore := mock_store.NewMockRepositoryUser(ctrl)
	storages := mock_store.NewMockStorages(ctrl)
	storages.EXPECT().UserStorage().Return(uStore).AnyTimes()
	storages.EXPECT().UserStorage().Return(posStore).AnyTimes()
	return storages
}
