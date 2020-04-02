package utils

import (
	"github.com/go-resty/resty/v2"
	"github.com/subosito/gotenv"
	"go.uber.org/zap"
	"os"
)

func LogRestyTrace(log *zap.SugaredLogger, resp *resty.Response) {
	log.Debug("Request Trace Info:")
	ti := resp.Request.TraceInfo()
	log.Debug("DNSLookup    :", ti.DNSLookup)
	log.Debug("ConnTime     :", ti.ConnTime)
	log.Debug("TLSHandshake :", ti.TLSHandshake)
	log.Debug("ServerTime   :", ti.ServerTime)
	log.Debug("ResponseTime :", ti.ResponseTime)
	log.Debug("TotalTime    :", ti.TotalTime)
	log.Debug("IsConnReused :", ti.IsConnReused)
	log.Debug("IsConnWasIdle:", ti.IsConnWasIdle)
	log.Debug("ConnIdleTime :", ti.ConnIdleTime)
}

// file name like ../../.env
func MustLoadDevEnv(filenames ...string) {
	if os.Getenv("SERVICE_ENV") == "dev" {
		gotenv.Must(gotenv.Load, filenames...)
	}
}
