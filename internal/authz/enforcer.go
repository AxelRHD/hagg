package authz

import (
	"log"
	"os"
	"path/filepath"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

func MustNewEnforcer() *casbin.Enforcer {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	modelPath := filepath.Join(wd, "model.conf")
	policyPath := filepath.Join(wd, "policy.csv")

	m, err := model.NewModelFromFile(modelPath)
	if err != nil {
		log.Fatalf("casbin: cannot load model.conf: %v", err)
	}

	a := fileadapter.NewAdapter(policyPath)

	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		log.Fatalf("casbin: cannot create enforcer: %v", err)
	}

	if err := e.LoadPolicy(); err != nil {
		log.Fatalf("casbin: cannot load policy.csv: %v", err)
	}

	return e
}
