package message

import (
	"fmt"
	"strings"

	"chronicl/internal/config"
	"github.com/DaDevFox/hof"
)

type Message struct {
	CommitType  string
	Scope       string
	Annotations []string
	message     string
}

func Serialize(commitType, scope, scopeAnnotation, commitTypeAnnotations, message string, cfg *config.Config) (string, error) {
	return SerializeMessage(Message{
		commitType,
		scope,
		hof.FilterToArray([]string{scopeAnnotation, commitTypeAnnotations}, func(s string) bool { return strings.Trim(s, " ") != "" }),
		message,
	}, cfg), nil
}

// TODO: make faster
func rightPadTesselation(targ, tesselation string, length int) string {
	result := targ

	for i := len(result); i < length; i++ {
		result += string(tesselation[i%len(tesselation)])
	}

	return result
}

// TODO: validation rules from config?

// TODO: make option to det. whether we want exact divisibility or not, if not skip this sanitization (rename to sanitize for exact divisibilty?)
// TODO: convert to obj output
// TODO: consider calling this config tesselation
func sanitizeTesselation(targetStringPossibilities []string, idealTesselation string) (string, int) {
	maxPadLen := hof.Max(hof.MapToArray(targetStringPossibilities, func(c string) int {
		return len(c)
	}))

	effectivePadLen := min(len(idealTesselation), maxPadLen)
	for maxPadLen%effectivePadLen != 0 {
		effectivePadLen--
	}

	return idealTesselation[0:effectivePadLen], maxPadLen
}

func SerializeMessage(message Message, cfg *config.Config) string {
	commitType :=
		strings.Trim(message.CommitType, " ")

	if cfg.PadCommitType {
		commitTesseleation, maxPadLen := sanitizeTesselation(hof.MapToArray(cfg.CommitTypes, func(c config.MessageFormatterConfig) string {
			return c.Key
		}), cfg.PadCharacter)

		commitType = rightPadTesselation(commitType, commitTesseleation, maxPadLen)
	}
	header := commitType

	scope := strings.Trim(message.Scope, " ")
	if scope != "" {
		if cfg.PadScope {
			commitTesseleation, maxPadLen := sanitizeTesselation(hof.MapToArray(cfg.Scopes, func(c config.MessageFormatterConfig) string {
				return c.Key
			}), cfg.PadCharacter)

			scope = rightPadTesselation(scope, commitTesseleation, maxPadLen)
		}

		header = fmt.Sprintf("%s(%s)", header, scope)
	}

	if len(message.Annotations) > 0 {
		return fmt.Sprintf(
			"%s: [%s] %s",
			header,
			strings.Join(message.Annotations, " "),
			message.message,
		)
	} else {
		return fmt.Sprintf(
			"%s: %s",
			header,
			message.message,
		)
	}
}
