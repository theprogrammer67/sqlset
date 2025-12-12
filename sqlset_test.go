package sqlset_test

import (
	"embed"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theprogrammer67/sqlset"
)

//go:embed testdata/valid/*.sql
var testdataValid embed.FS

//go:embed testdata/invalid/meta1.sql
var testdataInvalidMeta1 embed.FS

//go:embed testdata/invalid/meta2.sql
var testdataInvalidMeta2 embed.FS

//go:embed testdata/invalid/syntax1.sql
var testdataInvalidSyntax1 embed.FS

//go:embed testdata/invalid/syntax2.sql
var testdataInvalidSyntax2 embed.FS

//go:embed testdata/invalid/long-lines.sql
var testdataInvalidLongLines embed.FS

//nolint:funlen,lll
func TestSQLSet(t *testing.T) {
	set, err := sqlset.New(testdataValid)
	require.NoError(t, err)
	require.NotNil(t, set)

	queryTests := []struct {
		setID         string
		queryID       string
		expected      string
		expectedFound bool
	}{
		{
			setID:         "test-id-override-1",
			queryID:       "GetData1",
			expected:      "SELECT '515bbf3c-93c5-476a-8dbc-4a6db4fe3c0c' AS id, 'Igor' AS name, 'en' AS language, 'igor@example.com' AS email, ARRAY['token1','token2'] AS tokens;",
			expectedFound: true,
		},
		{
			setID:         "test-id-override-1",
			queryID:       "GetData2",
			expected:      "SELECT 'ef84af8f-bb55-4f74-9d7c-3db30e740d20' AS id, 'Alexey' AS name, 'en' AS language, 'alex@example.com' AS email, '{}'::varchar[] as tokens;",
			expectedFound: true,
		},
		{
			setID:         "test-id-override-1",
			queryID:       "GetData3",
			expected:      "SELECT 'e192f9e5-5e5c-4bba-b13e-0f9de32ec6bd' AS id, 'Denis' AS name, 'en' AS language, 'denis@example.com' AS email, ARRAY['token3','token4'] AS tokens;",
			expectedFound: true,
		},
		{
			setID:         "test-id-override-1",
			queryID:       "unknown",
			expectedFound: false,
		},
	}

	for _, test := range queryTests {
		t.Run("GetQuery "+test.setID+":"+test.queryID, func(t *testing.T) {
			t.Parallel()

			query, err := set.GetQuery(test.setID, test.queryID)

			if test.expectedFound {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, sqlset.ErrNotFound)
			}

			assert.Equal(t, test.expected, query)
		})

		t.Run("MustGetQuery "+test.setID+":"+test.queryID, func(t *testing.T) {
			t.Parallel()

			var query string

			fn := func() {
				query = set.MustGetQuery(test.setID, test.queryID)
			}

			if test.expectedFound {
				assert.NotPanics(t, fn)
			} else {
				assert.Panics(t, fn)
			}

			assert.Equal(t, test.expected, query)
		})
	}

	t.Run("GetAllMetas", func(t *testing.T) {
		t.Parallel()

		metas := set.GetAllMetas()

		require.Len(t, metas, 2)
		assert.Contains(t, metas, sqlset.QuerySetMeta{
			ID:          "test-id-override-1",
			Name:        "Test 1",
			Description: "Test description 1",
		})
		assert.Contains(t, metas, sqlset.QuerySetMeta{
			ID:          "test2",
			Name:        "test2",
			Description: "Test description 2",
		})
	})
}

func TestNew_WhenInvalid_ExpectError(t *testing.T) {
	tests := []struct {
		name        string
		fs          fs.FS
		expectedErr error
	}{
		{
			name:        "invalid meta 1",
			fs:          testdataInvalidMeta1,
			expectedErr: sqlset.ErrInvalidSyntax,
		},
		{
			name:        "invalid meta 2",
			fs:          testdataInvalidMeta2,
			expectedErr: sqlset.ErrInvalidSyntax,
		},
		{
			name:        "invalid syntax 1",
			fs:          testdataInvalidSyntax1,
			expectedErr: sqlset.ErrInvalidSyntax,
		},
		{
			name:        "invalid syntax 2",
			fs:          testdataInvalidSyntax2,
			expectedErr: sqlset.ErrInvalidSyntax,
		},
		{
			name:        "long lines",
			fs:          testdataInvalidLongLines,
			expectedErr: sqlset.ErrMaxLineLenExceeded,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			set, err := sqlset.New(test.fs)

			//nolint:testifylint
			assert.ErrorIs(t, err, test.expectedErr)
			assert.Nil(t, set)
		})
	}
}
