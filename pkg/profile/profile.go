/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
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
