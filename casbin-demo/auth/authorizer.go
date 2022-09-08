package auth

import (
	"errors"
	"fmt"

	"github.com/casbin/casbin/v2"
	_ "github.com/casbin/casbin/v2/model"
	// xormadapter "github.com/casbin/xorm-adapter/v2"
	//"google.golang.org/grpc/codes"
	//_ "google.golang.org/grpc/internal/status"
)

type Authorizer struct {
	enforcer *casbin.Enforcer
}

func New(model, policy string) *Authorizer {
	enforcer, _ := casbin.NewEnforcer(model, policy)
	return &Authorizer{
		enforcer: enforcer,
	}
}

func (a *Authorizer) Authorizer(subject, object, action string) error {
	if ok, _ := a.enforcer.Enforce(subject, object, action); !ok {
		msg := fmt.Sprintf(
			"%s not permitted to %s %s",
			subject,
			action,
			object,
		)
		// st := status.New(codes.PermissionDenied, msg)
		// return st.Err()
		return errors.New("PermissionDenied" + msg)
	}
	return nil
}
