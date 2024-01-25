package blob

import (
	"testing"
)

func TestNewAzureBlobStorage(t *testing.T) {
	cases := map[string]struct {
		inDomain string
		inSetenv func(key, value string)
		outStg   Storager
		outErr   error
	}{
		"new repo":                     {"custom-domain", t.Setenv, &azure{domain: "https://custom-domain.blob.core.windows.net/"}, nil},
		"missing environment variable": {"", t.Setenv, nil, ErrDomainEnv},
	}

	for key, testcase := range cases {
		t.Run(key, func(t *testing.T) {
			testcase.inSetenv("AZURE_STORAGE_ACCOUNT", testcase.inDomain)
			blobstg, err := NewAzureBlobStorage()
			if err != nil {
				if err != testcase.outErr {
					t.Errorf("blob: azure_blob_storage: test_new_azure_blob_storage: error interface mismatch (result = %q, expected = %q)\n", err, testcase.outErr)
				}
				return
			}

			switch {
			case blobstg.String() != testcase.outStg.String():
				t.Errorf("blob: azure_blob_storage: test_new_azure_blob_storage: `Storager` interface mismatch (result = %q, expected = %q)\n", blobstg, testcase.outStg)
			}
		})
	}
}
