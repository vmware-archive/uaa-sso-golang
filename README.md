### Example usage

#### creating a session for a user

When a user does not have a session you need to redirect to a UAA authorize url.  To do this, make use of the LogingURL method in uaa.UAAInterface.

First you have to new up a UAA struct. We provide constructors for this:

	loginHost := "http://login.10.244.0.34.xip.io" // bosh-lite example
	uaaHost := "http://uaa.10.244.0.34.xip.io"     // bosh-lite example
	
	uaa := uaa.NewUAA(loginHost, uaaHost, "your-uaa-client-id", "your-uaa-client-secret")

A few other things are need set before you begin:

	uaa.RedirectURL = "your/session/create/url"
	uaa.Scope = "scope settings for your app"
	uaa.AccessType = "offline" // or online

Now you can use the uaa we set up to get the loginUrl you need to redirect to.  An example ServeHTTP method is shown:

	func (handler SessionsNew) ServeHTTP(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, uaa.LoginURL(), http.StatusFound)
	}
	
The login url contains a number of pieces of information for UAA. ClientID and RedirectURI amongst them.  The redirect uri is what UAA will redirect too on successful login.

UAA will send along a "code". You send the code to uaa.Exchange(code) which will ping the cf uaa and get access and refresh tokens for you. It will return them as members on the uaa.Token struct.  An example call to Exchange:

	func (handler SessionsCreate) ServeHTTP(w http.ResponseWriter, req *http.Request) {
		code := req.URL.Query().Get("code")
		token, err := uaa.Exchange(code)
		
		// you probably want to store the tokens in your session,
		// they are large, so maybe in different sessions
		token.Access
		token.Refresh
	}
	

#### verifing a users session

For a request you need to verify you need to create a uaa.Token and populate the members from your session:

	token := uaa.Token{
		Access: "your access token",
		Refresh: "your refresh token",
	}
	
After that you can call token.IsPresent() to determine if both are set. And then use the IsExpired method to see if your token has expired:

	expired, err := token.IsExpired()
		
If your token has expried you can use the Refresh method to get a new one:
	
	if expired {
		token, err : = uaa.Refresh(token.Refresh)
		if err == uaa.InvalidRefreshToken {
			// handle appropriately, re-login?
		}
		// set new access / refresh tokens in your session and carry on
	}
	

### godoc
  
Check out the project into your go workspace then run the the following:

  godoc -http=:6060
  
Then navigate to the following in your favorite browser:
       
       http://localhost:6060/pkg/github.com/pivotal-cf/uaa-sso-golang/uaa/
