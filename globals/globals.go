// Globals: Secret, UserKey, SuscriptionID, ResourceGroupName
package globals

// We could add many others like the port of the server.
// We could also import these from a JSON, YAML, or local system envs.

// These variables start with an uppercase letter so they can be exported to other packages.
var Secret = []byte("secret")

const (
	// ...
	Userkey = "user"
	// Should from the body of a request
	SubscriptionID = "d88a3473-bb31-4a61-80a7-1e614fa1c2cc"
	// Should from the body of a request
	ResourceGroupName = "play-storage"
)
