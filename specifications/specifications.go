package specifications

//MethodPost struct to the POST verb of an endpoint
type MethodPost struct {
	POST EndPoint
}

//MethodGet struct to the GET verb of an endpoint
type MethodGet struct {
	GET EndPoint
}

//MethodPut struct to the PUT verb of an endpoint
type MethodPut struct {
	PUT EndPoint
}

//MethodOptions struct to the OPTIONS verb of an endpoint
type MethodOptions struct {
	OPTIONS EndPoint
}

//EndPoint struct to Store information about a prticular endpoint
type EndPoint struct {
	Description string  `json:"description"`
	Parameters  []Param `json:"parameters"`
}

//Param struct to store information about the parameter of the Query string of of an endpoint.
type Param struct {
	Name             string `json:"name"`
	ParameterDetails Detail `json:"details"`
}

//Detail stuct to store detailed information about the params a http verb Documentation.
type Detail struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

//UserOptions is a variable the that contains user information about the OPTIONS verb of this api endpint
var UserOptions = MethodOptions{OPTIONS: EndPoint{Description: "This page"}}

//UserPostParameters is a variable that contains information about the parameters of the POST verd this api endpoint
var UserPostParameters = []Param{{Name: "email", ParameterDetails: Detail{Type: "string",
	Description: "A new user's Email adderess", Required: false}}}

//UserPost is a variable that Contains information about POST verb of this api endpoint
var UserPost = MethodPost{POST: EndPoint{Description: "Create a new user", Parameters: UserPostParameters}}

//UserGet is a variable that contains information about the GET verb of This api endpoint
var UserGet = MethodGet{GET: EndPoint{Description: "Access a user"}}
