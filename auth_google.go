package devastator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nbusy/neptulon/jsonrpc"
)

// Response from GET https://www.googleapis.com/plus/v1/people/me?access_token=...
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

func googleAuth(ctx *jsonrpc.ReqContext) {
	token := ctx.Req.Params.(map[string]interface{})["accessToken"]
	_, _, err := getGProfile(token.(string))
	if err != nil {
		ctx.ResErr = &jsonrpc.ResError{Code: 666, Message: "Failed to authenticated user with Google+ OAuth access token."}
		log.Printf("Errored during Google+ profile call using provided access token: %v with error: %v", token, err)
	}

	// user is authenticated at this point so check if this is a first-time registration
	// if user, ok := users.GetByMail("..."); ok { ... }

	// if authenticated generate "userid", set it in session, create, store in database, and send client-certificate as reponse

	ctx.Res = "access granted"
}

// Retrieve user info (display name, e-mail, profile pic) using an access token that has 'profile' and 'email' scopes.
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

	if err = json.Unmarshal(resBody, profile); err != nil {
		return
	}

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
