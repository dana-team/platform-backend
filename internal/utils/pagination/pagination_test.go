package pagination

import (
	"github.com/dana-team/platform-backend/internal/types"
	"github.com/dana-team/platform-backend/internal/utils"
	"github.com/dana-team/platform-backend/internal/utils/testutils/mocks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	strForPagination = "string-for-pagination"
)

type TestPaginator struct {
	str string
}

func (p *TestPaginator) FetchList(listOptions metav1.ListOptions) (*types.List[string], error) {
	limit := listOptions.Limit
	startIndex := 0

	items := []string{
		p.str + "-1",
		p.str + "-2",
		p.str + "-3",
		p.str + "-4",
		p.str + "-5",
	}

	endIndex := startIndex + int(limit)
	if endIndex > len(items) {
		endIndex = len(items)
	}

	return &types.List[string]{
		Items: items[startIndex:endIndex],
	}, nil
}

func Test_buildListOptions(t *testing.T) {
	type args struct {
		limit         int64
		continueToken string
	}

	type want struct {
		listOptions metav1.ListOptions
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldBuildListOptions": {
			args: args{
				limit:         int64(20),
				continueToken: "token",
			},
			want: want{
				listOptions: metav1.ListOptions{
					Limit:         int64(20),
					Continue:      "token",
					LabelSelector: utils.ManagedLabelSelector,
				},
			},
		},
		"ShouldBuildEmptyListOptions": {
			args: args{
				limit:         0,
				continueToken: "",
			},
			want: want{
				listOptions: metav1.ListOptions{
					Limit:         0,
					Continue:      "",
					LabelSelector: utils.ManagedLabelSelector,
				},
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			listOptions := buildListOptions(test.args.limit, test.args.continueToken)
			assert.Equal(t, test.want.listOptions, listOptions)
		})
	}
}

func Test_FetchPage(t *testing.T) {
	type args struct {
		page      int
		limit     int
		paginator Paginator[types.List[string]]
	}

	type want struct {
		strings []string
		err     error
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldReturnOneStringOnFirstPage": {
			args: args{
				limit: 1,
				page:  1,
				paginator: &TestPaginator{
					str: strForPagination,
				},
			},
			want: want{
				strings: []string{
					strForPagination + "-1",
				},
				err: nil,
			},
		},
		"ShouldReturnAllStrings": {
			args: args{
				limit: 5,
				page:  1,
				paginator: &TestPaginator{
					str: strForPagination,
				},
			},
			want: want{
				strings: []string{
					strForPagination + "-1",
					strForPagination + "-2",
					strForPagination + "-3",
					strForPagination + "-4",
					strForPagination + "-5",
				},
				err: nil,
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			stringsList, err := FetchPage(test.args.limit, test.args.page, test.args.paginator)
			if test.want.err != nil {
				assert.Equal(t, test.want.err, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.want.strings, stringsList)

		})
	}
}

func Test_extractLimitFromCtx(t *testing.T) {
	const defaultPaginationLimitStr = "100"
	const defaultPaginationLimitInt = 100

	_ = os.Setenv(envDefaultPaginationLimit, defaultPaginationLimitStr)

	type args struct {
		limit int
	}

	type want struct {
		limit int
		err   error
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldReturnLimitFromCtx": {
			args: args{
				limit: 5,
			},
			want: want{
				limit: 5,
				err:   nil,
			},
		},
		"ShouldReturnDefaultLimitFromCtx": {
			want: want{
				limit: defaultPaginationLimitInt,
				err:   nil,
			},
		},

		"ShouldReturnInvalidLimit": {
			args: args{
				limit: -1,
			},
			want: want{
				limit: -1,
				err:   nil,
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			c := mocks.GinContext()
			mocks.SetPaginationValues(c, test.args.limit, 0)

			limit, err := extractLimitFromCtx(c)
			if test.want.err != nil {
				assert.Equal(t, test.want.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.want.limit, limit)
		})
	}
}

func Test_extractPageFromCtx(t *testing.T) {
	type args struct {
		page int
	}

	type want struct {
		page int
		err  error
	}

	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldReturnPage": {
			args: args{
				page: 5,
			},
			want: want{
				page: 5,
				err:  nil,
			},
		},
		"ShouldReturnDefaultPage": {
			want: want{
				page: 1,
				err:  nil,
			},
		},

		"ShouldReturnInvalidLimit": {
			args: args{
				page: -1,
			},
			want: want{
				page: -1,
				err:  nil,
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			c := mocks.GinContext()
			mocks.SetPaginationValues(c, 0, test.args.page)

			page, err := extractPageFromCtx(c)
			if test.want.err != nil {
				assert.Equal(t, test.want.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.want.page, page)
		})
	}
}
