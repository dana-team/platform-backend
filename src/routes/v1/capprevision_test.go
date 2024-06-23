package v1_test

import (
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const (
	cappRevisionNamespace = testName + "-capp-revision-ns"
	cappRevisionName      = testName + "-capp-revision"
	labelSelectorKey      = "labelSelector"
	labelKey              = "key"
	labelValue            = "value"
)

func setupCappRevisions() {
	createTestNamespace(cappRevisionNamespace)
	createTestCappRevision(cappRevisionName+"-1", cappRevisionNamespace, map[string]string{labelKey + "-1": labelValue + "-1"}, map[string]string{})
	createTestCappRevision(cappRevisionName+"-2", cappRevisionNamespace, map[string]string{labelKey + "-2": labelValue + "-2"}, map[string]string{})
}

func TestGetCappRevisions(t *testing.T) {
	type selector struct {
		keys   []string
		values []string
	}

	type requestParams struct {
		namespace     string
		labelSelector selector
	}

	type want struct {
		statusCode int
		response   types.CappRevisionList
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedGettingCappRevisions": {
			requestParams: requestParams{
				namespace: cappRevisionNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response:   types.CappRevisionList{Count: 2, CappRevisions: []string{cappRevisionName + "-1", cappRevisionName + "-2"}},
			},
		},
		"ShouldFailWithBadRequestInvalidURI": {
			requestParams: requestParams{
				namespace: "",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   types.CappRevisionList{},
			},
		},
		"ShouldSucceedGettingCappRevisionsWithLabelSelector": {
			requestParams: requestParams{
				namespace: cappRevisionNamespace,
				labelSelector: selector{
					keys:   []string{labelKey + "-1"},
					values: []string{labelValue + "-1"},
				},
			},
			want: want{
				statusCode: http.StatusOK,
				response:   types.CappRevisionList{Count: 1, CappRevisions: []string{cappRevisionName + "-1"}},
			},
		},
		"ShouldFailGettingCappRevisionsWithInvalidLabelSelector": {
			requestParams: requestParams{
				namespace: cappRevisionNamespace,
				labelSelector: selector{
					keys:   []string{labelKey + "-1"},
					values: []string{labelValue + " 1"},
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   types.CappRevisionList{},
			},
		},
		"ShouldSucceedGettingNoCappRevisionsWithLabelSelector": {
			requestParams: requestParams{
				namespace: cappRevisionNamespace,
				labelSelector: selector{
					keys:   []string{labelKey + "-3"},
					values: []string{labelValue + "-3"},
				},
			},
			want: want{
				statusCode: http.StatusOK,
				response:   types.CappRevisionList{Count: 0, CappRevisions: nil},
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}

			for i, key := range test.requestParams.labelSelector.keys {
				params.Add(labelSelectorKey, key+"="+test.requestParams.labelSelector.values[i])
			}

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capprevisions/", test.requestParams.namespace)

			request, _ := http.NewRequest("GET", fmt.Sprintf("%s?%s", baseURI, params.Encode()), nil)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			if writer.Code == http.StatusOK {
				var response types.CappRevisionList
				if err := json.Unmarshal(writer.Body.Bytes(), &response); err != nil {
					panic(err)
				}

				assert.Equal(t, test.want.response.Count, response.Count)
				assert.Equal(t, test.want.response.CappRevisions, response.CappRevisions)
			}
		})
	}
}

func TestGetCappRevision(t *testing.T) {
	type requestParams struct {
		name      string
		namespace string
	}

	type want struct {
		statusCode int
		response   types.CappRevision
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedGettingCappRevision": {
			requestParams: requestParams{
				namespace: cappRevisionNamespace,
				name:      cappRevisionName + "-1",
			},
			want: want{
				statusCode: http.StatusOK,
				response: utils.GetBareCappRevisionType(cappRevisionName+"-1", cappRevisionNamespace,
					[]types.KeyValue{{Key: "key1", Value: "value-1"}}, []types.KeyValue{}),
			},
		},
		"ShouldFailWithBadRequestInvalidURI": {
			requestParams: requestParams{
				namespace: "",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   types.CappRevision{},
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/capprevisions/%s", test.requestParams.namespace, test.requestParams.name)
			request, _ := http.NewRequest("GET", baseURI, nil)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			if writer.Code == http.StatusOK {
				var response types.CappRevision
				if err := json.Unmarshal(writer.Body.Bytes(), &response); err != nil {
					panic(err)
				}

				assert.Equal(t, test.want.response.Metadata.Name, response.Metadata.Name)
				assert.Equal(t, test.want.response.Metadata.Namespace, response.Metadata.Namespace)
			}
		})
	}
}
