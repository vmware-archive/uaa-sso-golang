# THIS CODEBASE IS DEPRECATED. IT IS NO LONGER SUPPORTED.

### Example usage

#### creating a session for a user

When a user does not have a session you need to redirect to a UAA authorize url.  To do this, make use of the LogingURL method in uaa.UAAInterface.

First you have to new up a UAA object. We provide constructors for this:

	loginHost := "http://login.10.244.0.34.xip.io" // bosh-lite example
	uaaHost := "http://uaa.10.244.0.34.xip.io"     // bosh-lite example
	
	uaaObject := uaa.NewUAA(loginHost, uaaHost, "your-uaa-client-id", "your-uaa-client-secret")

A few other things need set before you begin:

	uaaObject.RedirectURL = "your/session/create/url"
	uaaObject.Scope = "scope settings for your app"
	uaaObject.AccessType = "offline" // or online

Now you can use the uaa object we set up to get the login url you need to redirect to.  An example ServeHTTP method is shown:

	func (handler SessionsNew) ServeHTTP(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, uaaObject.LoginURL(), http.StatusFound)
	}
	
The login url contains a number of pieces of information for UAA. ClientID and RedirectURI amongst them.  The redirect uri is what UAA will redirect too on successful login.

UAA will send along a code. You send the code to uaa.Exchange which will ping the cf uaa and get access and refresh tokens for you. It will return them as members on an uaa.Token.  An example call to Exchange:

	func (handler SessionsCreate) ServeHTTP(w http.ResponseWriter, req *http.Request) {
		code := req.URL.Query().Get("code")
		token, err := uaaObject.Exchange(code)
		
		// you probably want to store the tokens in your session,
		// they are large, so maybe in different sessions if session is in cookie
		token.Access
		token.Refresh
	}
	
Tokens when converted to json will give you the appropriate keys:

	access_token
	refresh_token

	

#### verifing a users session

To verify a user that has a session you need to create a uaa.Token and populate the members from your session:

	token := uaa.Token{
		Access: "your access token",
		Refresh: "your refresh token",
	}
	
After that you can call token.IsPresent() to determine if both access and refresh tokens are set. And then use the IsExpired method on the token to see if it has expired:

	expired, err := token.IsExpired()
		
If your token has expired you can use the Refresh method to get a new one:
	
	if expired {
		token, err : = uaaObject.Refresh(token.Refresh)
		if err == uaa.InvalidRefreshToken {
			// handle appropriately, re-login?
		}
		// set new access / refresh tokens in your session and carry on
	}
	

### godoc

The documentation can be found [here](http://godoc.org/github.com/pivotal-cf/uaa-sso-golang/uaa).
  
To view documentation locally:
* Checkout the repo into your `$GOPATH`.
* Run `godoc -http=:6060` at your terminal.
* Then navigate to [this](http://localhost:6060/pkg/github.com/pivotal-cf/uaa-sso-golang/uaa/) in your browser.
