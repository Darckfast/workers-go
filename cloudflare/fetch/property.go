//go:build js && wasm

package fetch

import (
	"context"
	"errors"
	"syscall/js"

	jsconv "github.com/Darckfast/workers-go/internal/conv"
	jsruntime "github.com/Darckfast/workers-go/internal/runtime"
)

type IncomingBotManagementJsDetection struct {
	Passed bool
}

func NewIncomingBotManagementJsDetection(cf js.Value) *IncomingBotManagementJsDetection {
	if cf.IsUndefined() {
		return nil
	}
	return &IncomingBotManagementJsDetection{
		Passed: cf.Get("passed").Bool(),
	}
}

type IncomingBotManagement struct {
	CorporateProxy bool
	VerifiedBot    bool
	JsDetection    *IncomingBotManagementJsDetection
	StaticResource bool
	Score          int
}

func NewIncomingBotManagement(cf js.Value) *IncomingBotManagement {
	if cf.IsUndefined() {
		return nil
	}
	return &IncomingBotManagement{
		CorporateProxy: cf.Get("corporateProxy").Bool(),
		VerifiedBot:    cf.Get("verifiedBot").Bool(),
		JsDetection:    NewIncomingBotManagementJsDetection(cf.Get("jsDetection")),
		StaticResource: cf.Get("staticResource").Bool(),
		Score:          cf.Get("score").Int(),
	}
}

type IncomingTLSClientAuth struct {
	CertIssuerDNLegacy    string
	CertIssuerSKI         string
	CertSubjectDNRFC2253  string
	CertSubjectDNLegacy   string
	CertFingerprintSHA256 string
	CertNotBefore         string
	CertSKI               string
	CertSerial            string
	CertIssuerDN          string
	CertVerified          string
	CertNotAfter          string
	CertSubjectDN         string
	CertPresented         string
	CertRevoked           string
	CertIssuerSerial      string
	CertIssuerDNRFC2253   string
	CertFingerprintSHA1   string
}

func NewIncomingTLSClientAuth(cf js.Value) *IncomingTLSClientAuth {
	if cf.IsUndefined() {
		return nil
	}
	return &IncomingTLSClientAuth{
		CertIssuerDNLegacy:    jsconv.MaybeString(cf.Get("certIssuerDNLegacy")),
		CertIssuerSKI:         jsconv.MaybeString(cf.Get("certIssuerSKI")),
		CertSubjectDNRFC2253:  jsconv.MaybeString(cf.Get("certSubjectDNRFC2253")),
		CertSubjectDNLegacy:   jsconv.MaybeString(cf.Get("certSubjectDNLegacy")),
		CertFingerprintSHA256: jsconv.MaybeString(cf.Get("certFingerprintSHA256")),
		CertNotBefore:         jsconv.MaybeString(cf.Get("certNotBefore")),
		CertSKI:               jsconv.MaybeString(cf.Get("certSKI")),
		CertSerial:            jsconv.MaybeString(cf.Get("certSerial")),
		CertIssuerDN:          jsconv.MaybeString(cf.Get("certIssuerDN")),
		CertVerified:          jsconv.MaybeString(cf.Get("certVerified")),
		CertNotAfter:          jsconv.MaybeString(cf.Get("certNotAfter")),
		CertSubjectDN:         jsconv.MaybeString(cf.Get("certSubjectDN")),
		CertPresented:         jsconv.MaybeString(cf.Get("certPresented")),
		CertRevoked:           jsconv.MaybeString(cf.Get("certRevoked")),
		CertIssuerSerial:      jsconv.MaybeString(cf.Get("certIssuerSerial")),
		CertIssuerDNRFC2253:   jsconv.MaybeString(cf.Get("certIssuerDNRFC2253")),
		CertFingerprintSHA1:   jsconv.MaybeString(cf.Get("certFingerprintSHA1")),
	}
}

type IncomingTLSExportedAuthenticator struct {
	ClientFinished  string
	ClientHandshake string
	ServerHandshake string
	ServerFinished  string
}

func NewIncomingTLSExportedAuthenticator(cf js.Value) *IncomingTLSExportedAuthenticator {
	if cf.IsUndefined() {
		return nil
	}
	return &IncomingTLSExportedAuthenticator{
		ClientFinished:  jsconv.MaybeString(cf.Get("clientFinished")),
		ClientHandshake: jsconv.MaybeString(cf.Get("clientHandshake")),
		ServerHandshake: jsconv.MaybeString(cf.Get("serverHandshake")),
		ServerFinished:  jsconv.MaybeString(cf.Get("serverFinished")),
	}
}

type IncomingProperties struct {
	Longitude                string
	Latitude                 string
	TLSCipher                string
	Continent                string
	Asn                      int
	ClientAcceptEncoding     string
	Country                  string
	TLSClientAuth            *IncomingTLSClientAuth
	TLSExportedAuthenticator *IncomingTLSExportedAuthenticator
	TLSVersion               string
	Colo                     string
	Timezone                 string
	City                     string
	VerifiedBotCategory      string
	// EdgeRequestKeepAliveStatus int
	RequestPriority string
	HttpProtocol    string
	Region          string
	RegionCode      string
	AsOrganization  string
	PostalCode      string
	BotManagement   *IncomingBotManagement
}

func NewIncomingProperties(ctx context.Context) (*IncomingProperties, error) {
	obj := jsruntime.MustExtractTriggerObj(ctx)
	cf := obj.Get("cf")
	if cf.IsUndefined() {
		return nil, errors.New("runtime is not cloudflare")
	}

	return &IncomingProperties{
		Longitude:                jsconv.MaybeString(cf.Get("longitude")),
		Latitude:                 jsconv.MaybeString(cf.Get("latitude")),
		TLSCipher:                jsconv.MaybeString(cf.Get("tlsCipher")),
		Continent:                jsconv.MaybeString(cf.Get("continent")),
		Asn:                      jsconv.MaybeInt(cf.Get("asn")),
		ClientAcceptEncoding:     jsconv.MaybeString(cf.Get("clientAcceptEncoding")),
		Country:                  jsconv.MaybeString(cf.Get("country")),
		TLSClientAuth:            NewIncomingTLSClientAuth(cf.Get("tlsClientAuth")),
		TLSExportedAuthenticator: NewIncomingTLSExportedAuthenticator(cf.Get("tlsExportedAuthenticator")),
		TLSVersion:               cf.Get("tlsVersion").String(),
		Colo:                     cf.Get("colo").String(),
		Timezone:                 cf.Get("timezone").String(),
		City:                     jsconv.MaybeString(cf.Get("city")),
		VerifiedBotCategory:      jsconv.MaybeString(cf.Get("verifiedBotCategory")),
		RequestPriority:          jsconv.MaybeString(cf.Get("requestPriority")),
		HttpProtocol:             cf.Get("httpProtocol").String(),
		Region:                   jsconv.MaybeString(cf.Get("region")),
		RegionCode:               jsconv.MaybeString(cf.Get("regionCode")),
		AsOrganization:           cf.Get("asOrganization").String(),
		PostalCode:               jsconv.MaybeString(cf.Get("postalCode")),
		BotManagement:            NewIncomingBotManagement(cf.Get("botManagement")),
	}, nil
}
