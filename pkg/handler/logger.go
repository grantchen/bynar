/**
    @author: dongjs
    @date: 2023/9/14
    @description:
**/

package handler

import (
	"database/sql"
	goErrors "errors"

	"github.com/sirupsen/logrus"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/errors"
)

// LogInternalError log internal error messages
func LogInternalError(err error) {
	fromError := errors.FromError(err)
	if fromError != nil && fromError.IsInternal() {
		if !goErrors.Is(err, sql.ErrNoRows) {
			logrus.Errorf("%s, stack: %s", fromError.FullError(), fromError.Stack())
		}

		if !config.IsProductionEnv() {
			logrus.Printf("%s, stack: %s", fromError.FullError(), fromError.Stack())
		}
	}
}
