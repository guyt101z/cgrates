package engine

import (
	"testing"

	"github.com/cgrates/cgrates/utils"
)

var testMap = UserMap{
	table: map[string]map[string]string{
		"test:user":   map[string]string{"t": "v"},
		":user":       map[string]string{"t": "v"},
		"test:":       map[string]string{"t": "v"},
		"test1:user1": map[string]string{"t": "v", "x": "y"},
	},
	index: make(map[string][]string),
}

func TestUsersAdd(t *testing.T) {
	tm := NewUserMap(ratingStorage)
	var r string
	up := UserProfile{
		Tenant:   "test",
		UserName: "user",
		Profile: map[string]string{
			"t": "v",
		},
	}
	tm.SetUser(up, &r)
	p, found := tm.table[up.GetId()]
	if r != utils.OK ||
		!found ||
		p["t"] != "v" ||
		len(tm.table) != 1 ||
		len(p) != 1 {
		t.Error("Error setting user: ", tm)
	}
}

func TestUsersUpdate(t *testing.T) {
	tm := NewUserMap(ratingStorage)
	var r string
	up := UserProfile{
		Tenant:   "test",
		UserName: "user",
		Profile: map[string]string{
			"t": "v",
		},
	}
	tm.SetUser(up, &r)
	p, found := tm.table[up.GetId()]
	if r != utils.OK ||
		!found ||
		p["t"] != "v" ||
		len(tm.table) != 1 ||
		len(p) != 1 {
		t.Error("Error setting user: ", tm)
	}
	up.Profile["x"] = "y"
	tm.UpdateUser(up, &r)
	p, found = tm.table[up.GetId()]
	if r != utils.OK ||
		!found ||
		p["x"] != "y" ||
		len(tm.table) != 1 ||
		len(p) != 2 {
		t.Error("Error updating user: ", tm)
	}
}

func TestUsersUpdateNotFound(t *testing.T) {
	tm := NewUserMap(ratingStorage)
	var r string
	up := UserProfile{
		Tenant:   "test",
		UserName: "user",
		Profile: map[string]string{
			"t": "v",
		},
	}
	tm.SetUser(up, &r)
	up.UserName = "test1"
	err := tm.UpdateUser(up, &r)
	if err != utils.ErrNotFound {
		t.Error("Error detecting user not found on update: ", err)
	}
}

func TestUsersUpdateInit(t *testing.T) {
	tm := NewUserMap(ratingStorage)
	var r string
	up := UserProfile{
		Tenant:   "test",
		UserName: "user",
	}
	tm.SetUser(up, &r)
	up = UserProfile{
		Tenant:   "test",
		UserName: "user",
		Profile: map[string]string{
			"t": "v",
		},
	}
	tm.UpdateUser(up, &r)
	p, found := tm.table[up.GetId()]
	if r != utils.OK ||
		!found ||
		p["t"] != "v" ||
		len(tm.table) != 1 ||
		len(p) != 1 {
		t.Error("Error updating user: ", tm)
	}
}

func TestUsersRemove(t *testing.T) {
	tm := NewUserMap(ratingStorage)
	var r string
	up := UserProfile{
		Tenant:   "test",
		UserName: "user",
		Profile: map[string]string{
			"t": "v",
		},
	}
	tm.SetUser(up, &r)
	p, found := tm.table[up.GetId()]
	if r != utils.OK ||
		!found ||
		p["t"] != "v" ||
		len(tm.table) != 1 ||
		len(p) != 1 {
		t.Error("Error setting user: ", tm)
	}
	tm.RemoveUser(up, &r)
	p, found = tm.table[up.GetId()]
	if r != utils.OK ||
		found ||
		len(tm.table) != 0 {
		t.Error("Error removing user: ", tm)
	}
}

func TestUsersGetFull(t *testing.T) {
	up := UserProfile{
		Tenant:   "test",
		UserName: "user",
		Profile: map[string]string{
			"t": "v",
		},
	}
	results := make([]*UserProfile, 0)
	testMap.GetUsers(up, &results)
	if len(results) != 1 {
		t.Error("error getting users: ", results)
	}
}

func TestUsersGetTenant(t *testing.T) {
	up := UserProfile{
		Tenant:   "testX",
		UserName: "user",
		Profile: map[string]string{
			"t": "v",
		},
	}
	results := make([]*UserProfile, 0)
	testMap.GetUsers(up, &results)
	if len(results) != 0 {
		t.Error("error getting users: ", results)
	}
}

func TestUsersGetUserName(t *testing.T) {
	up := UserProfile{
		Tenant:   "test",
		UserName: "userX",
		Profile: map[string]string{
			"t": "v",
		},
	}
	results := make([]*UserProfile, 0)
	testMap.GetUsers(up, &results)
	if len(results) != 0 {
		t.Error("error getting users: ", results)
	}
}

func TestUsersGetNotFoundProfile(t *testing.T) {
	up := UserProfile{
		Tenant:   "test",
		UserName: "user",
		Profile: map[string]string{
			"o": "p",
		},
	}
	results := make([]*UserProfile, 0)
	testMap.GetUsers(up, &results)
	if len(results) != 0 {
		t.Error("error getting users: ", results)
	}
}

func TestUsersGetMissingTenant(t *testing.T) {
	up := UserProfile{
		UserName: "user",
		Profile: map[string]string{
			"t": "v",
		},
	}
	results := make([]*UserProfile, 0)
	testMap.GetUsers(up, &results)
	if len(results) != 2 {
		t.Error("error getting users: ", results)
	}
}

func TestUsersGetMissingUserName(t *testing.T) {
	up := UserProfile{
		Tenant: "test",
		Profile: map[string]string{
			"t": "v",
		},
	}
	results := make([]*UserProfile, 0)
	testMap.GetUsers(up, &results)
	if len(results) != 2 {
		t.Error("error getting users: ", results)
	}
}

func TestUsersGetMissingId(t *testing.T) {
	up := UserProfile{
		Profile: map[string]string{
			"t": "v",
		},
	}
	results := make([]*UserProfile, 0)
	testMap.GetUsers(up, &results)
	if len(results) != 4 {
		t.Error("error getting users: ", results)
	}
}

func TestUsersGetMissingIdTwo(t *testing.T) {
	up := UserProfile{
		Profile: map[string]string{
			"t": "v",
			"x": "y",
		},
	}
	results := make([]*UserProfile, 0)
	testMap.GetUsers(up, &results)
	if len(results) != 1 {
		t.Error("error getting users: ", results)
	}
}

func TestUsersAddIndex(t *testing.T) {
	var r string
	testMap.AddIndex([]string{"t"}, &r)
	if r != utils.OK ||
		len(testMap.index) != 1 ||
		len(testMap.index[utils.ConcatenatedKey("t", "v")]) != 4 {
		t.Error("error adding index: ", testMap.index)
	}
}

func TestUsersAddIndexFull(t *testing.T) {
	var r string
	testMap.index = make(map[string][]string) // reset index
	testMap.AddIndex([]string{"t", "x", "UserName", "Tenant"}, &r)
	if r != utils.OK ||
		len(testMap.index) != 6 ||
		len(testMap.index[utils.ConcatenatedKey("t", "v")]) != 4 {
		t.Error("error adding index: ", testMap.index)
	}
}

func TestUsersAddIndexNone(t *testing.T) {
	var r string
	testMap.index = make(map[string][]string) // reset index
	testMap.AddIndex([]string{"test"}, &r)
	if r != utils.OK ||
		len(testMap.index) != 0 {
		t.Error("error adding index: ", testMap.index)
	}
}

func TestUsersGetFullindex(t *testing.T) {
	var r string
	testMap.index = make(map[string][]string) // reset index
	testMap.AddIndex([]string{"t", "x", "UserName", "Tenant"}, &r)
	up := UserProfile{
		Tenant:   "test",
		UserName: "user",
		Profile: map[string]string{
			"t": "v",
		},
	}
	results := make([]*UserProfile, 0)
	testMap.GetUsers(up, &results)
	if len(results) != 1 {
		t.Error("error getting users: ", results)
	}
}

func TestUsersGetTenantindex(t *testing.T) {
	var r string
	testMap.index = make(map[string][]string) // reset index
	testMap.AddIndex([]string{"t", "x", "UserName", "Tenant"}, &r)
	up := UserProfile{
		Tenant:   "testX",
		UserName: "user",
		Profile: map[string]string{
			"t": "v",
		},
	}
	results := make([]*UserProfile, 0)
	testMap.GetUsers(up, &results)
	if len(results) != 0 {
		t.Error("error getting users: ", results)
	}
}

func TestUsersGetUserNameindex(t *testing.T) {
	var r string
	testMap.index = make(map[string][]string) // reset index
	testMap.AddIndex([]string{"t", "x", "UserName", "Tenant"}, &r)
	up := UserProfile{
		Tenant:   "test",
		UserName: "userX",
		Profile: map[string]string{
			"t": "v",
		},
	}
	results := make([]*UserProfile, 0)
	testMap.GetUsers(up, &results)
	if len(results) != 0 {
		t.Error("error getting users: ", results)
	}
}

func TestUsersGetNotFoundProfileindex(t *testing.T) {
	var r string
	testMap.index = make(map[string][]string) // reset index
	testMap.AddIndex([]string{"t", "x", "UserName", "Tenant"}, &r)
	up := UserProfile{
		Tenant:   "test",
		UserName: "user",
		Profile: map[string]string{
			"o": "p",
		},
	}
	results := make([]*UserProfile, 0)
	testMap.GetUsers(up, &results)
	if len(results) != 0 {
		t.Error("error getting users: ", results)
	}
}

func TestUsersGetMissingTenantindex(t *testing.T) {
	var r string
	testMap.index = make(map[string][]string) // reset index
	testMap.AddIndex([]string{"t", "x", "UserName", "Tenant"}, &r)
	up := UserProfile{
		UserName: "user",
		Profile: map[string]string{
			"t": "v",
		},
	}
	results := make([]*UserProfile, 0)
	testMap.GetUsers(up, &results)
	if len(results) != 2 {
		t.Error("error getting users: ", results)
	}
}

func TestUsersGetMissingUserNameindex(t *testing.T) {
	var r string
	testMap.index = make(map[string][]string) // reset index
	testMap.AddIndex([]string{"t", "x", "UserName", "Tenant"}, &r)
	up := UserProfile{
		Tenant: "test",
		Profile: map[string]string{
			"t": "v",
		},
	}
	results := make([]*UserProfile, 0)
	testMap.GetUsers(up, &results)
	if len(results) != 2 {
		t.Error("error getting users: ", results)
	}
}

func TestUsersGetMissingIdindex(t *testing.T) {
	var r string
	testMap.index = make(map[string][]string) // reset index
	testMap.AddIndex([]string{"t", "x", "UserName", "Tenant"}, &r)
	up := UserProfile{
		Profile: map[string]string{
			"t": "v",
		},
	}
	results := make([]*UserProfile, 0)
	testMap.GetUsers(up, &results)
	if len(results) != 4 {
		t.Error("error getting users: ", results)
	}
}

func TestUsersGetMissingIdTwoINdex(t *testing.T) {
	var r string
	testMap.index = make(map[string][]string) // reset index
	testMap.AddIndex([]string{"t", "x", "UserName", "Tenant"}, &r)
	up := UserProfile{
		Profile: map[string]string{
			"t": "v",
			"x": "y",
		},
	}
	results := make([]*UserProfile, 0)
	testMap.GetUsers(up, &results)
	if len(results) != 1 {
		t.Error("error getting users: ", results)
	}
}