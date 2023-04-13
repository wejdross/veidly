package api

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func (ctx *Ctx) RunTls() error {

	if !ctx.Request.WithTls {
		return fmt.Errorf("requested TLS without specifying tls params in config")
	}

	c := tls.Config{
		CipherSuites: []uint16{
			// 1.3
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_AES_128_GCM_SHA256,
			// 1.2 but required
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS13,
		MaxVersion:               tls.VersionTLS13,
	}

	s := http.Server{
		TLSConfig: &c,
		Handler:   ctx.engine,
		Addr:      ctx.Request.Addr,
	}

	return s.ListenAndServeTLS(ctx.Request.CertPath, ctx.Request.KeyPath)
}
