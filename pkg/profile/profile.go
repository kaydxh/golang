package profile

import (
	"strings"

	os_ "github.com/kaydxh/golang/go/os"
	"github.com/pkg/profile"
)

//no profile
type nop struct {
}

func (p *nop) Stop() {
}

func StartWithEnv() interface {
	Stop()
} {
	profileEnv := os_.GetEnvAsStringOrFallback("PROFILING", "")
	return start(profileEnv)
	//	profiles := strings.Split(profileEnv, ",")
	/*
		for _, pf := range profiles {
			start(pf)
		}
	*/
}

func start(pf string) interface {
	Stop()
} {
	switch strings.ToLower(pf) {
	case "cpu":
		return profile.Start(profile.CPUProfile, profile.NoShutdownHook)
	case "mem":
		return profile.Start(profile.MemProfile, profile.NoShutdownHook)
	case "mutex":
		return profile.Start(profile.MutexProfile, profile.NoShutdownHook)
	case "block":
		return profile.Start(profile.BlockProfile, profile.NoShutdownHook)
	case "trace":
		return profile.Start(profile.TraceProfile, profile.NoShutdownHook)
	case "thread_create":
		return profile.Start(profile.ThreadcreationProfile, profile.NoShutdownHook)
	case "goroutine":
		return profile.Start(profile.GoroutineProfile, profile.NoShutdownHook)
	default:
		//do nothing
		return new(nop)
	}
}
