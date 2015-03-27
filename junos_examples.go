package junos

// Establishing a session
func Example() {
	jnpr, err := junos.NewSession(host, user, password)
	if err != nil {
		log.Fatal(err)
	}
	defer jnpr.Close()
}
