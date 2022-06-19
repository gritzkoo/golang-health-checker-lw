package healthchecker

// Readiness response
type Readiness struct {
	Name         string        `json:"name,omitempty"`
	Status       bool          `json:"status"`
	Version      string        `json:"version,omitempty"`
	Date         string        `json:"date"`
	Duration     float64       `json:"duration"`
	Integrations []Integration `json:"integrations"`
}

const fullyFunctional = "fully functional"

// Liveness response
type Liveness struct {
	Status  string `json:"status"`
	Version string `json:"version,omitempty"`
}

// Integration default response contract
type Integration struct {
	// The name you give in the Check.Name struct when your create one
	Name string `json:"name"`
	// The status if returned or not an error in the CheckResponse.Error
	Status bool `json:"status"`
	// The amount of time expend to execute your Check.Handle function
	ResponseTime float64 `json:"response_time"`
	// The value you placed in CheckResponse.URL
	URL string `json:"url"`
	// The error passed in the CheckResponse.Error

	Error error `json:"error,omitempty"`
}

// CheckResponse is the main struct to be used outside this package
// here you will inform this package
type CheckResponse struct {
	/*
		Here comes the error when you execute your test.

		When this field is not nil, then the package will
		assume that your test fails, and returns a status false in this integration
		and place in the main status the same value

		If you are getting a empty object in the JSON response,
		try create a struct with a Error() string function in it
		to make pretty responses like:

		type MyCustomError struct {
			Message string `json:"message,omitempty"`
			Code int `json:code,omitempty`
		}

		func (e *MyCustomError) Error() string {
			return fmt.Printf("A error occur! message: %s, code: %d", e.Message, e.Code)
		}

		Then:

		func MyTest() CheckResponse {
			return CheckResponse{
				Error: &MyCustomError{Message: "something got wrong!", Code: 9999},
			}
		}
	*/
	Error error `json:"error,omitempty"`
	// Use this URL field to expose the host of your test to make simple
	// troubleshooting.
	//
	// It's optional but can help you a lot!
	URL string `json:"url,omitempty"`
}

// Config the setup of this package
type Config struct {
	// The name of your application, can be empty
	Name string
	/*
		Version is the identifier of you app current runs on
		You can create a file lien revision.txt with the command
		git show -s --format="%ai %H %s %aN" HEAD > revision.txt
		then place its content in this field

		One example of this command out put is:

		2022-06-19 10:17:22 -0300 2b99fb28d4596fca9d782d7c582bfd71d2a592b4 docs: include new infos Gritzko Daniel Kleiner

		Using in the healthchecker.New example:

		version, err := ioutil.ReadFile("revision.txt")
		if err != nil {
			version = []byte{}
		}
		var check = healthchecker.New(healthchecker.Config{
			Version: string(version)
		})
	*/
	Version string
	// Is the list of checks you need to execute.
	//
	// You can place as many functions as you like, then this package will
	// generate a channel and execute all of the asynchronously
	Integrations []Check
}

// Check used to inform each integration config
type Check struct {
	// Use this name to identify you integration to make easy to
	// create querys in your log interface to show how many times your
	// integration fails
	Name string
	// The Handle function will execute wherever you need
	// just create a function to test someting and return a
	// CheckResponse interface and done, you can test anything
	Handle func() CheckResponse
}
