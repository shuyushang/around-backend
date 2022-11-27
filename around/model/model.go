package model

type Post struct {
	//struct ~~ class
	//public or private? fields start with UpperCase -> public; lowercase -> private
	//fields
	Id string `json:"id"`
	//json label is for converting to json; in json, how to match for the key value pair
	//`` -> if there is a special symbol in the string, just treat it like a regular symbol
	User    string `json:"user"`
	Message string `json:"message"`
	Url     string `json:"url"`
	//picture save in drive, and return an url；
	//if multiple pictures in one post, we can change the type to array
	Type string `json:"type"`
	//save the post type is a pic or video for frontend to make a difference
}

// user struct
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Age      int64  `json:"age"`
	Gender   string `json:"gender"`
}
