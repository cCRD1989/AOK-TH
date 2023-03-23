package aokmodel

type Userlogin struct {
	//ID              string
	Username string
	// Password        string
	// Gold            int
	// Cash            int
	// Email           string
	// IsEmailVerified int
	// AuthType        int
	// AccessToken     string
	// UserLevel       int
	// UnbanTime       int
	// CreateAt        time.Time
	// UpdateAt        time.Time
}

func (n *Userlogin) TableName() string {
	return "userlogin"
}
