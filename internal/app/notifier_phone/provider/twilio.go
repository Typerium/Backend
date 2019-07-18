package provider

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"typerium/internal/pkg/web"
)

const (
	twilioAPI = "https://api.twilio.com"
)

func NewTwilioProvider(clientFactory web.ClientFactory, username, password string) (Provider, error) {
	p := &twilio{
		clientFactory: clientFactory,
		authToken:     web.BasicAuth(username, password),
	}

	err := p.auth(username)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return p, nil
}

type twilio struct {
	clientFactory web.ClientFactory
	authToken     []byte
	authInfo      *twilioAuthInfo
}

type twilioMessage struct {
}

func (p *twilio) Send(number string, body string) error {
	client := p.clientFactory.Acquire()
	defer p.clientFactory.Release(client)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	p.prepareRequest(req, p.authInfo.Resources.Messages, fasthttp.MethodPost)

	data := url.Values{}
	data.Set("Body", body)
	data.Set("To", number)
	data.Set("From", "")
	req.SetBodyString(data.Encode())

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		return errors.WithStack(err)
	}

	out := new(twilioMessage)
	err = p.unmarshalResponse(resp, out)
	if err != nil {
		return err
	}

	return nil
}

func (p *twilio) prepareRequest(req *fasthttp.Request, path, method string) {
	req.SetRequestURI(twilioAPI + path)
	req.Header.SetMethod(method)
	req.Header.SetBytesV(fasthttp.HeaderAuthorization, p.authToken)
}

type twilioError struct {
	Code     int    `json:"code"`
	Detail   string `json:"detail"`
	Message  string `json:"message"`
	MoreInfo string `json:"more_info"`
	Status   int    `json:"status"`
}

func (e *twilioError) Error() string {
	return ""
}

func (p *twilio) unmarshalResponse(resp *fasthttp.Response, out interface{}) (err error) {
	body := resp.Body()

	if resp.StatusCode() != 200 {
		respErr := new(twilioError)
		err = json.Unmarshal(body, respErr)
		if err != nil {
			return errors.WithStack(err)
		}
		return errors.WithStack(respErr)
	}

	err = json.Unmarshal(body, out)
	return errors.WithStack(err)
}

type twilioAuthInfo struct {
	Status          string           `json:"status"`
	AuthToken       string           `json:"auth_token"`
	FriendlyName    string           `json:"friendly_name"`
	OwnerAccountSID string           `json:"owner_account_sid"`
	URI             string           `json:"uri"`
	SID             string           `json:"sid"`
	Type            string           `json:"type"`
	Resources       *twilioResources `json:"subresource_uris"`
	DateCreated     *timeRFC1123Z    `json:"date_created, string"`
	DateUpdated     *timeRFC1123Z    `json:"date_updated, string"`
}

type timeRFC1123Z struct {
	time.Time
}

func (t *timeRFC1123Z) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	t.Time, err = time.Parse(`"`+time.RFC1123Z+`"`, string(data))
	return errors.WithStack(err)
}

type twilioResources struct {
	Addresses             string `json:"addresses"`
	Conferences           string `json:"conferences"`
	SigningKeys           string `json:"signing_keys"`
	Transcriptions        string `json:"transcriptions"`
	ConnectApps           string `json:"connect_apps"`
	SIP                   string `json:"sip"`
	AuthorizedConnectApps string `json:"authorized_connect_apps"`
	Usage                 string `json:"usage"`
	Keys                  string `json:"keys"`
	Applications          string `json:"applications"`
	Recordings            string `json:"recordings"`
	ShortCodes            string `json:"short_codes"`
	Calls                 string `json:"calls"`
	Notifications         string `json:"notifications"`
	IncomingPhoneNumbers  string `json:"incoming_phone_numbers"`
	Queues                string `json:"queues"`
	Messages              string `json:"messages"`
	OutgoingCallerIDs     string `json:"outgoing_caller_ids"`
	AvailablePhoneNumbers string `json:"available_phone_numbers"`
	Balance               string `json:"balance"`
}

func (p *twilio) auth(username string) error {
	client := p.clientFactory.Acquire()
	defer p.clientFactory.Release(client)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	p.prepareRequest(req, fmt.Sprintf("/2010-04-01/Accounts/%s.json", username), fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		return errors.WithStack(err)
	}

	if p.authInfo == nil {
		p.authInfo = new(twilioAuthInfo)
	}

	err = p.unmarshalResponse(resp, p.authInfo)
	if err != nil {
		return err
	}

	return nil
}
