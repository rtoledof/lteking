package internal

import (
	"log/slog"
	"os"
	"runtime/debug"
	"strconv"

	"order.io/pkg/slogctx"
)

func NewAppLogger(name string, level slog.Level) *slog.Logger {
	bi, _ := debug.ReadBuildInfo()

	var ver struct {
		vcs         string
		vcsrevision string
		vcsmodified bool
		cgoenabled  bool
		cgocflags   string
		cgocppflags string
		cgocxxflags string
		cgoldflags  string
		goarch      string
		goos        string
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, v := range info.Settings {
			switch v.Key {
			case "CGO_ENABLED":
				ver.cgoenabled, _ = strconv.ParseBool(v.Value)
			case "CGO_CFLAGS":
				ver.cgocflags = v.Value
			case "CGO_CPPFLAGS":
				ver.cgocppflags = v.Value
			case "CGO_CXXFLAGS":
				ver.cgocxxflags = v.Value
			case "CGO_LDFLAGS":
				ver.cgoldflags = v.Value
			case "vcs.revision":
				ver.vcsrevision = v.Value
			case "vcs.modified":
				ver.vcsmodified, _ = strconv.ParseBool(v.Value)
			case "GOOS":
				ver.goos = v.Value
			case "GOARCH":
				ver.goarch = v.Value
			case "vcs":
				ver.vcs = v.Value
			}
		}
	}
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: level})
	logger := slog.New(slogctx.NewContextAttrsHandler(h))
	return logger.With(
		slog.Group("app",
			slog.String("name", name),
			slog.Int("pid", os.Getpid()),
			slog.String("version", bi.Main.Version),
			slog.Group("vcs",
				slog.String("name", ver.vcs),
				slog.String("revision", ver.vcsrevision),
				slog.Bool("modified", ver.vcsmodified),
			),
			slog.Group("go",
				slog.String("version", bi.GoVersion),
				slog.Bool("cgo_enabled", ver.cgoenabled),
				slog.String("cgo_cflags", ver.cgocflags),
				slog.String("cgo_cppflags", ver.cgocppflags),
				slog.String("cgo_cxxflags", ver.cgocxxflags),
				slog.String("cgo_ldflags", ver.cgoldflags),
				slog.String("goarch", ver.goarch),
				slog.String("goos", ver.goos),
			),
		),
	)
}
