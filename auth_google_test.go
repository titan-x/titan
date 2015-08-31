package devastator

import (
	"encoding/json"
	"testing"
)

// Response from: GET https://www.googleapis.com/plus/v1/people/me?access_token=...
// (retrieved on 5th of July, 2015)
var googleAuthResStr = []byte(`{
 "kind": "plus#person",
 "etag": "\"abcd123\"",
 "gender": "male",
 "emails": [
  {
   "value": "chuck@devastator.com",
   "type": "account"
  }
 ],
 "objectType": "person",
 "id": "1234567890",
 "displayName": "Chuck Norris",
 "name": {
  "familyName": "Chuck",
  "givenName": "Norris"
 },
 "url": "https://plus.google.com/1234567890",
 "image": {
  "url": "https://lh3.googleusercontent.com/asdfgsgdsfs/afsdggg/sdafdsff/1234dfgg/photo.jpg?sz=50",
  "isDefault": false
 },
 "isPlusUser": true,
 "language": "en",
 "circledByCount": 999,
 "verified": false
}`)

var (
	displayName = "Chuck Norris"
	email       = "chuck@devastator.com"
	imgURL      = "https://lh3.googleusercontent.com/asdfgsgdsfs/afsdggg/sdafdsff/1234dfgg/photo.jpg?sz=50"
)

func TestGoogleAuth(t *testing.T) {
	var profile gProfile
	if err := json.Unmarshal(googleAuthResStr, &profile); err != nil {
		t.Fatal(err)
	}

	if profile.DisplayName != displayName ||
		profile.Emails[0].Value != email ||
		profile.Image.URL != imgURL {
		t.Fatal("Cannot deserialize Google 'plus/v1/people/me' response.")
	}
}
