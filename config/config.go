package config

// RolePermission defines what each role is allowed to do.
// Public things like browsing products don't need to come through here.
func RolePermission(userRole, action string) bool {
	// Admins can do everything.
	if userRole == "admin" {
		return true
	}

	// Customers can do a limited set of things.
	if userRole == "customer" {
		switch action {
		case
			"manage:cart",  // add/update/remove items, view own cart
			"checkout",     // place an order
			"write:review", // create or update review
			"read:orders":  // view own order history
			return true
		}
	}

	// Default: not allowed.
	return false
}
