package profile

import (
	"strings"

	os_ "github.com/kaydxh/golang/go/os"
	"github.com/pkg/profile"
)

//const ()

//type noprofie struct{}

//func (p *noprofie) Stop() {}

//type Profile struct {
//}

func StartWithEnv() {
	profileEnv := os_.GetEnvAsStringOrFallback("PROFILING", "")
	profiles := strings.Split(profileEnv, ",")
	for _, pf := range profiles {
		start(pf)
	}
}

func start(pf string) {
	switch strings.ToLower(pf) {
	case "cpu":
		defer profile.Start(profile.CPUProfile, profile.NoShutdownHook).Stop()
	case "mem":
		defer profile.Start(profile.MemProfile, profile.NoShutdownHook).Stop()
	case "mutex":
		defer profile.Start(profile.MutexProfile, profile.NoShutdownHook).Stop()
	case "block":
		defer profile.Start(profile.BlockProfile, profile.NoShutdownHook).Stop()
	case "trace":
		defer profile.Start(profile.TraceProfile, profile.NoShutdownHook).Stop()
	case "thread_create":
		defer profile.Start(profile.ThreadcreationProfile, profile.NoShutdownHook).Stop()
	case "goroutine":
		defer profile.Start(profile.GoroutineProfile, profile.NoShutdownHook).Stop()
	default:
		//do nothing
	}
}
