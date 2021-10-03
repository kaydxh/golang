package profile

import (
	"path/filepath"
	"strings"
	"time"

	time_ "github.com/kaydxh/golang/go/time"

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
	profilePath := os_.GetEnvAsStringOrFallback("PROFILEPATH", "")
	if profilePath != "" {
		profilePath = filepath.Join(profilePath, "profile"+time.Now().Format(time_.TimeMillFormat))
	}

	return start(profileEnv, profilePath)
	//	profiles := strings.Split(profileEnv, ",")
	/*
		for _, pf := range profiles {
			start(pf)
		}
	*/
}

func start(pf, dir string) interface {
	Stop()
} {
	switch strings.ToLower(pf) {
	case "cpu":
		return profile.Start(profile.CPUProfile, profile.NoShutdownHook, profile.ProfilePath(dir))
	case "mem":
		return profile.Start(profile.MemProfile, profile.NoShutdownHook, profile.ProfilePath(dir))
	case "mutex":
		return profile.Start(profile.MutexProfile, profile.NoShutdownHook, profile.ProfilePath(dir))
	case "block":
		return profile.Start(profile.BlockProfile, profile.NoShutdownHook, profile.ProfilePath(dir))
	case "trace":
		return profile.Start(profile.TraceProfile, profile.NoShutdownHook, profile.ProfilePath(dir))
	case "thread_create":
		return profile.Start(profile.ThreadcreationProfile, profile.NoShutdownHook, profile.ProfilePath(dir))
	case "goroutine":
		return profile.Start(profile.GoroutineProfile, profile.NoShutdownHook, profile.ProfilePath(dir))
	default:
		//do nothing
		return new(nop)
	}
}
