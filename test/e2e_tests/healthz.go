package e2e_tests

import (
	"fmt"
	"github.com/dana-team/platform-backend/src/utils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("Validate Healthz route", func() {
	It("Should return correct status code", func() {
		uri := fmt.Sprintf("%s/healthz", platformURL)
		status, response := performHTTPRequest(httpClient, nil, http.MethodGet, uri, "", "", userToken)

		expectedResponse := map[string]interface{}{
			testutils.MessageKey: "ok",
		}

		Expect(status).Should(Equal(http.StatusOK))
		compareResponses(response, expectedResponse)
	})
})
