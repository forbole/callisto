package keybase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetAvatarURL returns the avatar URL from the given identity.
// If no identity is found, it returns an empty string instead.
func GetAvatarURL(identity string) (string, error) {
	if len(identity) < 16 {
		return "", nil
	}

	var response IdentityQueryResponse
	endpoint := fmt.Sprintf("/user/lookup.json?key_suffix=%[1]s&fields=basics&fields=pictures", identity)
	err := queryKeyBase(endpoint, &response)
	if err != nil {
		return "", fmt.Errorf("error while querying keybase: %s", err)
	}

	// The server responded with an error
	if response.Status.Code != 0 {
		return "", fmt.Errorf("response code not valid: %s", response.Status.ErrDesc)
	}

	// No images found
	if len(response.Objects) == 0 {
		return "", nil
	}

	// Either the pictures do not exist, or the primary one does not exist, or the URL is empty
	data := response.Objects[0]
	if data.Pictures == nil || data.Pictures.Primary == nil || len(data.Pictures.Primary.URL) == 0 {
		return "", nil
	}

	// The picture URL is found
	return data.Pictures.Primary.URL, nil
}

// queryKeyBase queries the Keybase APIs for the given endpoint, and de-serializes
// the response as a JSON object inside the given ptr
func queryKeyBase(endpoint string, ptr interface{}) error {
	resp, err := http.Get("https://keybase.io/_/api/1.0" + endpoint)
	if err != nil {
		return fmt.Errorf("error while querying keybase APIs: %s", err)
	}

	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error while reading response body: %s", err)
	}

	err = json.Unmarshal(bz, &ptr)
	if err != nil {
		return fmt.Errorf("error while unmarshaling response body: %s", err)
	}

	return nil
}
