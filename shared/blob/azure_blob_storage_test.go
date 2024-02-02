package blob

import (
	"testing"
)

func TestNewAzureBlobStorage(t *testing.T) {
	cases := map[string]struct {
		inSetenv  func(key, value string)
		inAccount string
		inKey     string
		outStg    Storager
		outErr    error
	}{
		"new repo": {t.Setenv, "devstoreaccount1", "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==", &azure{domain: "https://devstoreaccount1.blob.core.windows.net/"}, nil},
		"missing environment variable `AZURE_STORAGE_ACCOUNT`": {t.Setenv, "", "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==", nil, ErrAccountEnvVar},
		"missing environment variable `AZURE_STORAGE_KEY`":     {t.Setenv, "devstoreaccount1", "", nil, ErrKeyEnvVar},
	}

	for key, testcase := range cases {
		t.Run(key, func(t *testing.T) {
			testcase.inSetenv("AZURE_STORAGE_ACCOUNT", testcase.inAccount)
			testcase.inSetenv("AZURE_STORAGE_KEY", testcase.inKey)

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
