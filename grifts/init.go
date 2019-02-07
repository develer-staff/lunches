package grifts

import (
	"github.com/develersrl/lunches/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
