package memoryDS

import (
	"MyProject/apiSchema/userSchema"
	userDataModel "MyProject/models/user/dataModel"
	"context"
	"sync"
	"sync/atomic"
)

type UserDBDS struct {
	idCounter int64
	users     []userDataModel.User
	lock      sync.RWMutex
}

func NewTaskDBDS(stastID int64) *UserDBDS {
	return &UserDBDS{
		idCounter: stastID,
		users:     []userDataModel.User{},
	}
}

func (ds *UserDBDS) CreateUser(ctx context.Context, req userSchema.LoginRequest) (userDataModel.User, error) {
	_ = ctx
	user := userDataModel.User{
		ID:     atomic.AddInt64(&ds.idCounter, 1),
		Code:   req.Code,
		Name:   req.Name,
		Family: req.Family,
	}

	ds.lock.Lock()
	ds.users = append(ds.users, user)
	ds.lock.Unlock()

	return user, nil
}
