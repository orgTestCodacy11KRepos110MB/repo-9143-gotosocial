/*
   GoToSocial
   Copyright (C) 2021-2022 GoToSocial Authors admin@gotosocial.org

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package testrig

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/superseriousbusiness/activity/pub"
	"github.com/superseriousbusiness/activity/streams"
	apimodel "github.com/superseriousbusiness/gotosocial/internal/api/model"
	"github.com/superseriousbusiness/gotosocial/internal/concurrency"
	"github.com/superseriousbusiness/gotosocial/internal/db"
	"github.com/superseriousbusiness/gotosocial/internal/federation"
	"github.com/superseriousbusiness/gotosocial/internal/messages"
	"github.com/superseriousbusiness/gotosocial/internal/transport"
)

// NewTestTransportController returns a test transport controller with the given http client.
//
// Obviously for testing purposes you should not be making actual http calls to other servers.
// To obviate this, use the function NewMockHTTPClient in this package to return a mock http
// client that doesn't make any remote calls but just returns whatever you tell it to.
//
// Unlike the other test interfaces provided in this package, you'll probably want to call this function
// PER TEST rather than per suite, so that the do function can be set on a test by test (or even more granular)
// basis.
func NewTestTransportController(client pub.HttpClient, db db.DB, fedWorker *concurrency.WorkerPool[messages.FromFederator]) transport.Controller {
	return transport.NewController(db, NewTestFederatingDB(db, fedWorker), &federation.Clock{}, client)
}

// NewMockHTTPClient returns a client that conforms to the pub.HttpClient interface.
//
// If do is nil, then a standard response set will be mocked out, which includes models stored in the
// testrig, and webfinger responses as well.
//
// If do is not nil, then the given do function will always be used, which allows callers
// to customize how the client is mocked.
//
// Note that you should never ever make ACTUAL http calls with this thing.
func NewMockHTTPClient(do func(req *http.Request) (*http.Response, error)) pub.HttpClient {
	if do == nil {
		do = func(req *http.Request) (*http.Response, error) {
			responseCode := http.StatusNotFound
			responseBytes := []byte(`{"error":"404 not found"}`)
			responseContentType := "application/json"
			responseContentLength := len(responseBytes)

			if strings.Contains(req.URL.String(), ".well-known/webfinger") {
				responseCode, responseBytes, responseContentType, responseContentLength = WebfingerResponse(req)
			} else {
				if note, ok := suite.testRemoteStatuses[req.URL.String()]; ok {
					// the request is for a note that we have stored
					noteI, err := streams.Serialize(note)
					if err != nil {
						panic(err)
					}
					noteJson, err := json.Marshal(noteI)
					if err != nil {
						panic(err)
					}
					responseBytes = noteJson
					responseContentType = "application/activity+json"
				}

				if person, ok := suite.testRemotePeople[req.URL.String()]; ok {
					// the request is for a person that we have stored
					personI, err := streams.Serialize(person)
					if err != nil {
						panic(err)
					}
					personJson, err := json.Marshal(personI)
					if err != nil {
						panic(err)
					}
					responseBytes = personJson
					responseContentType = "application/activity+json"
				}

				if group, ok := suite.testRemoteGroups[req.URL.String()]; ok {
					// the request is for a person that we have stored
					groupI, err := streams.Serialize(group)
					if err != nil {
						panic(err)
					}
					groupJson, err := json.Marshal(groupI)
					if err != nil {
						panic(err)
					}
					responseBytes = groupJson
					responseContentType = "application/activity+json"
				}

				if service, ok := suite.testRemoteServices[req.URL.String()]; ok {
					serviceI, err := streams.Serialize(service)
					if err != nil {
						panic(err)
					}
					serviceJson, err := json.Marshal(serviceI)
					if err != nil {
						panic(err)
					}
					responseBytes = serviceJson
					responseContentType = "application/activity+json"
				}

				if attachment, ok := suite.testRemoteAttachments[req.URL.String()]; ok {
					responseBytes = attachment.Data
					responseContentType = attachment.ContentType
				}
			}

			logrus.Debugf("returning response %s", string(responseBytes))
			reader := bytes.NewReader(responseBytes)
			readCloser := io.NopCloser(reader)
			return &http.Response{
				StatusCode:    responseCode,
				Body:          readCloser,
				ContentLength: int64(responseContentLength),
				Header: http.Header{
					"content-type": {responseContentType},
				},
			}, nil
		}
	}

	return &mockHTTPClient{
		do: do,
	}
}

type mockHTTPClient struct {
	do func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.do(req)
}

func WebfingerResponse(req *http.Request) (responseCode int, responseBytes []byte, responseContentType string, responseContentLength int) {
	var wfr *apimodel.WellKnownResponse

	switch req.URL.String() {
	case "https://unknown-instance.com/.well-known/webfinger?resource=acct:some_group@unknown-instance.com":
		wfr = &apimodel.WellKnownResponse{
			Subject: "acct:some_group@unknown-instance.com",
			Links: []apimodel.Link{
				{
					Rel:  "self",
					Type: "application/activity+json",
					Href: "https://unknown-instance.com/groups/some_group",
				},
			},
		}
	case "https://owncast.example.org/.well-known/webfinger?resource=acct:rgh@owncast.example.org":
		wfr = &apimodel.WellKnownResponse{
			Subject: "acct:rgh@example.org",
			Links: []apimodel.Link{
				{
					Rel:  "self",
					Type: "application/activity+json",
					Href: "https://owncast.example.org/federation/user/rgh",
				},
			},
		}
	case "https://unknown-instance.com/.well-known/webfinger?resource=acct:brand_new_person@unknown-instance.com":
		wfr = &apimodel.WellKnownResponse{
			Subject: "acct:brand_new_person@unknown-instance.com",
			Links: []apimodel.Link{
				{
					Rel:  "self",
					Type: "application/activity+json",
					Href: "https://unknown-instance.com/users/brand_new_person",
				},
			},
		}
	case "https://turnip.farm/.well-known/webfinger?resource=acct:turniplover6969@turnip.farm":
		wfr = &apimodel.WellKnownResponse{
			Subject: "acct:turniplover6969@turnip.farm",
			Links: []apimodel.Link{
				{
					Rel:  "self",
					Type: "application/activity+json",
					Href: "https://turnip.farm/users/turniplover6969",
				},
			},
		}
	case "https://fossbros-anonymous.io/.well-known/webfinger?resource=acct:foss_satan@fossbros-anonymous.io":
		wfr = &apimodel.WellKnownResponse{
			Subject: "acct:foss_satan@fossbros-anonymous.io",
			Links: []apimodel.Link{
				{
					Rel:  "self",
					Type: "application/activity+json",
					Href: "https://fossbros-anonymous.io/users/foss_satan",
				},
			},
		}
	case "https://example.org/.well-known/webfinger?resource=acct:some_user@example.org":
		wfr = &apimodel.WellKnownResponse{
			Subject: "acct:some_user@example.org",
			Links: []apimodel.Link{
				{
					Rel:  "self",
					Type: "application/activity+json",
					Href: "https://example.org/users/some_user",
				},
			},
		}
	}

	if wfr == nil {
		logrus.Debugf("webfinger response not available for %s", req.URL)
		responseCode = http.StatusNotFound
		responseBytes = []byte(`{"error":"not found"}`)
		responseContentType = "application/json"
		responseContentLength = len(responseBytes)
		return
	}

	wfrJson, err := json.Marshal(wfr)
	if err != nil {
		panic(err)
	}
	responseCode = http.StatusOK
	responseBytes = wfrJson
	responseContentType = "application/json"
	responseContentLength = len(wfrJson)
	return
}
