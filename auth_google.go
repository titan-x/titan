package devastator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nbusy/neptulon/jsonrpc"
)

// email: profile.emails[0].value,
// name: profile.displayName,
// picture: (yield request.get(profile.image.url, {encoding: 'base64'})).body

// retrieve user info (display name, e-mail, profile pic) using an access token that has 'profile' and 'email' scopes
func googleAuth(ctx *jsonrpc.ReqContext) {
	token := ctx.Req.Params.(map[string]interface{})["accessToken"]
	uri := fmt.Sprintf("https://www.googleapis.com/plus/v1/people/me?access_token=%s", token)

	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var profile interface{}
	json.Unmarshal(b, &profile)

	log.Fatalf("%s", b)

	// if authenticated generate "userid", set it in session, create and send client-certificate as reponse
	ctx.Res = "access granted"
}
