package integration

type (
	TCP struct {
		URL string
	}
	HTTP struct {
		URL string
	}
)

type Checker interface {
	Check() error
	GetURL() string
}
