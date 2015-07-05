package devastator

import (
	"encoding/json"
	"testing"
)

// retrieved on 5th of July, 2015
var googleAuthRes = []byte(`{
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

func TestGoogleAuth(t *testing.T) {
	var p googleProfile
	json.Unmarshal(googleAuthRes, &p)

	t.Fatal(p)
}
