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

// Retrieve user info (display name, e-mail, profile pic) using an access token that has 'profile' and 'email' scopes.
func googleAuth(ctx *jsonrpc.ReqContext) {
	token := ctx.Req.Params.(map[string]interface{})["accessToken"]
	uri := fmt.Sprintf("https://www.googleapis.com/plus/v1/people/me?access_token=%s", token)

	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var profile gProfile
	if err := json.Unmarshal(resBody, &profile); err != nil {
		log.Fatal(err)
	}

	// if authenticated generate "userid", set it in session, create, store in database, and send client-certificate as reponse
	ctx.Res = "access granted"
}
