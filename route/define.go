package route

type ROOT struct {
	ANY
}

func (r ROOT) Path() string {
	return "/"
}

type ANY struct{}

func (ANY) Method() string { return "ANY" }

type GET struct{}

func (GET) Method() string { return "GET" }

type POST struct{}

func (POST) Method() string { return "POST" }

type DELETE struct{}

func (DELETE) Method() string { return "DELETE" }

type PATCH struct{}

func (PATCH) Method() string { return "PATCH" }

type PUT struct{}

func (PUT) Method() string { return "PUT" }

type OPTIONS struct{}

func (OPTIONS) Method() string { return "OPTIONS" }

type HEAD struct{}

func (HEAD) Method() string { return "HEAD" }
