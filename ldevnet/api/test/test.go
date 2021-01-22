package test

import (
	"os"

	"github.com/filecoin-project/lotus/build"
)

func init() {
	os.Setenv("BELLMAN_NO_GPU", "1")

	build.InsecurePoStValidation = true
}
