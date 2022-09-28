package hcnettoml

import "log"

// ExampleGetTOML gets the hcnet.toml file for coins.asia
func ExampleClient_GetHcnetToml() {
	_, err := DefaultClient.GetHcnetToml("coins.asia")
	if err != nil {
		log.Fatal(err)
	}
}
