package integration

func (t *TCP) Check() error {
	conn, err := net.DialTimeout("tcp", t.URL, 5*time.Second)
	if err != nil {
		return fmt.Errorf("%s TCP ist nicht erreichbar", t.URL)
	}
	defer conn.Close()
	return nil
}

func (t *TCP) GetUrl() string {
	return t.URL
}
