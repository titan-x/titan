package devastator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nbusy/neptulon/jsonrpc"
)

// Response from GET https://www.googleapis.com/plus/v1/people/me?access_token=... (with scope 'profile' and 'email')
// has the following structure with denoted fields of interest (rest is left out):
type gProfile struct {
	Emails      []gEmail
	DisplayName string
	Image       gImage
}

type gEmail struct {
	Value string
}

type gImage struct {
	URL string
}

// CertResponse is the success response returned after a successful Google authentication.
type CertResponse struct {
	Cert, Key []byte
}

// googleAuth authenticates a user with Google+ using provided OAuth 2.0 access token.
// If authenticated successfully, user profile is retrieved from Google+ and user is given a TLS client-certificate in return.
func googleAuth(ctx *jsonrpc.ReqContext, db DB, certMgr *CertMgr) {
	t := ctx.Req.Params.(map[string]interface{})["accessToken"]
	p, i, err := getGProfile(t.(string))
	if err != nil {
		ctx.ResErr = &jsonrpc.ResError{Code: 666, Message: "Failed to authenticated user with Google+ OAuth access token."}
		log.Printf("Errored during Google+ profile call using provided access token: %v with error: %v", t, err)
	}

	// retrieve user information
	user, ok := db.GetByMail(p.Emails[0].Value)
	if !ok {
		// this is a first-time registration so create user profile via Google+ profile info
		user = &User{Email: p.Emails[0].Value, Name: p.DisplayName, Picture: i, Cert: make([]byte, 555)}
	}

	// re-create user client certificate as we don't store private key
	cert, key, err := certMgr.GenClientCert(string(user.ID))
	if err != nil {
		log.Fatal("Failed to generate client certificate for user:", err)
	}

	user.Cert = cert
	if err := db.SaveUser(user); err != nil {
		log.Fatal("Failed to persist user information:", err)
	}

	ctx.Conn.Data.Set("userid", user.ID)
	ctx.Res = CertResponse{Cert: user.Cert, Key: key}
	return
}

// getGProfile retrieves user info (display name, e-mail, profile pic) using an access token that has 'profile' and 'email' scopes.
// Also retrieves user profile image via profile image URL provided the response.
func getGProfile(token string) (profile *gProfile, profilePic []byte, err error) {
	// retrieve profile info from Google
	uri := fmt.Sprintf("https://www.googleapis.com/plus/v1/people/me?access_token=%s", token)
	res, err := http.Get(uri)
	if err != nil {
		return
	}

	resBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return
	}

	var p gProfile
	if err = json.Unmarshal(resBody, &p); err != nil {
		return
	}
	profile = &p

	// retrieve profile image
	uri = profile.Image.URL
	res, err = http.Get(uri)
	if err != nil {
		return
	}

	profilePic, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return
	}

	return
}
