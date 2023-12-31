package access

func CanUser(isAdmin bool, whoOrgID, whoID, userID, userOrgID int) bool {
	if whoID == userID {
		return true
	}
	if isAdmin && whoOrgID == userOrgID {
		return true
	}
	return false

}

func CanAdmin(isAdmin bool, whoID, adminID int) bool {
	if !isAdmin {
		return false
	}
	return whoID == adminID
}

func CanPosition(isAdmin bool, whoOrgID, orgPositionID int) bool {
	if !isAdmin {
		return false
	}
	if whoOrgID != orgPositionID {
		return false
	}

	return true

}

func CanCompany(isAdmin bool, whoOrgID, orgID int) bool {
	if isAdmin && whoOrgID == orgID {
		return true
	}
	return false

}
