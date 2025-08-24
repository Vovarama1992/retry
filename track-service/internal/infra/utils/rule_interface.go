package utils

import (
	"github.com/Vovarama1992/retry/pkg/domain"
)

type ActionRule interface {
	ShouldSkip(a domain.Action) bool
}
