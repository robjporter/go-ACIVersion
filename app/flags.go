package app

import (
	"strings"

	"github.com/robjporter/go-functions/as"
	"github.com/robjporter/go-functions/kingpin"
)

var (
	add       = kingpin.Command("add", "Register a new ACI APIC.")
	update    = kingpin.Command("update", "Update an ACI APIC.")
	delete    = kingpin.Command("delete", "Remove an ACI APIC.")
	show      = kingpin.Command("show", "Show an ACI APIC.")
	run       = kingpin.Command("run", "Run the main application.")
	addUCS    = add.Command("aci", "Add an APIC")
	updateUCS = update.Command("aci", "Update an APIC")
	deleteUCS = delete.Command("aci", "Delete an APIC")
	showUCS   = show.Command("aci", "Show an APIC")

	showAll = show.Command("all", "Show all")

	addACIIP       = addUCS.Flag("ip", "IP Address or DNS name for ACI APIC, without http(s).").Required().IP()
	addACIUsername = addUCS.Flag("username", "Name of user.").Required().String()
	addACIPassword = addUCS.Flag("password", "Password for user in plain text.").Required().String()

	updateACIIP       = updateUCS.Flag("ip", "IP Address or DNS name for ACI APIC, without http(s).").Required().IP()
	updateACIUsername = updateUCS.Flag("username", "Name of user.").Required().String()
	updateACIPassword = updateUCS.Flag("password", "Password for user in plain text.").Required().String()

	deleteACIIP = deleteUCS.Flag("ip", "IP Address or DNS name for ACI APIC, without http(s).").Required().IP()

	showACIIP = showUCS.Flag("ip", "IP Address or DNS name for ACI APIC, without http(s).").Required().IP()
)

func ProcessCommandLineArguments() string {
	switch kingpin.Parse() {
	case "run":
		return "RUN"
	case "add aci":
		return "ADDACI|" + as.ToString(*addACIIP) + "|" + *addACIUsername + "|" + *addACIPassword
	case "update aci":
		return "UPDATEACI|" + as.ToString(*updateACIIP) + "|" + *updateACIUsername + "|" + *updateACIPassword
	case "delete aci":
		return "DELETEACI|" + as.ToString(*deleteACIIP)
	case "show aci":
		return "SHOWACI|" + as.ToString(*showACIIP)
	case "show all":
		return "SHOWALL"
	}
	return ""
}

func (a *Application) processResponse(response string) {
	a.Log("Processing command line options.", map[string]interface{}{"args": response}, true)
	splits := strings.Split(response, "|")
	switch splits[0] {
	case "RUN":
		a.runAll()
	case "ADDACI":
		a.addACISystem(splits[1], splits[2], splits[3])
	case "UPDATEACI":
		a.updateACISystem(splits[1], splits[2], splits[3])
	case "DELETEACI":
		a.deleteACISystem(splits[1])
	case "SHOWACI":
		a.showACISystem(splits[1])
	case "SHOWALL":
		a.showACISystems()

	}
}

func (a *Application) addACISystem(ip, username, password string) {
	if !a.checkACIExists(ip) {
		if a.addACI(ip, username, password) {
			a.saveConfig()
			a.LogInfo("New ACI APIC has been added successfully.", map[string]interface{}{"IP": ip, "Username": username}, false)
		} else {
			a.LogInfo("ACI APIC could not be added.", map[string]interface{}{"IP": ip, "Username": username}, false)
		}
	} else {
		a.LogInfo("An ACI APIC already exsists in the config file.", map[string]interface{}{"IP": ip, "Username": username}, false)
	}
}

func (a *Application) addACI(ip, username, password string) bool {
	if ip != "" {
		if username != "" {
			if password != "" {
				tmp := ACISystemInfo{}
				tmp.ip = ip
				tmp.username = username
				tmp.password = a.EncryptPassword(password)
				a.ACI = append(a.ACI, tmp)
				return true
			} else {
				a.Log("The password for the ACI APIC cannot be blank.", nil, false)
			}
		} else {
			a.Log("The username for the ACI APIC cannot be blank.", nil, false)
		}
	} else {
		a.Log("The URL for the ACI APIC cannot be blank.", nil, false)
	}
	return false
}

func (a *Application) deleteACISystem(ip string) {
	if a.checkACIExists(ip) {
		if a.deleteACI(ip) {
			a.saveConfig()
			a.LogInfo("ACI APIC has been deleted successfully.", map[string]interface{}{"IP": ip}, true)
		} else {
			a.Log("ACI APIC could not be deleted.", map[string]interface{}{"IP": ip}, false)
		}
	} else {
		a.LogInfo("ACI APIC does not exsists and so cannot be deleted.", map[string]interface{}{"IP": ip}, false)
	}
}

func (a *Application) deleteACI(ip string) bool {
	for i := 0; i < len(a.ACI); i++ {
		if a.ACI[i].ip == as.ToString(ip) {
			a.ACI = append(a.ACI[:i], a.ACI[i+1:]...)
		}
	}
	return true
}

func (a *Application) showACI(ip string) {
	for i := 0; i < len(a.ACI); i++ {
		if a.ACI[i].ip == as.ToString(ip) {
			a.LogInfo("ACI APIC", map[string]interface{}{"URL": a.ACI[i].ip}, false)
			a.LogInfo("ACI APIC", map[string]interface{}{"Username": a.ACI[i].username}, false)
			a.LogInfo("ACI APIC", map[string]interface{}{"Password": a.ACI[i].password}, false)
		}
	}
}

func (a *Application) showACISystem(ip string) {
	if a.checkACIExists(ip) {
		a.showACI(ip)
	} else {
		a.Log("The ACI APIC does not exist and so cannot be displayed.", map[string]interface{}{"URL": ip}, false)
	}
}

func (a *Application) showACISystems() {
	a.getAllSystems()
	for i := 0; i < len(a.ACI); i++ {
		a.LogInfo("ACI APIC", map[string]interface{}{"URL": a.ACI[i].ip}, false)
		a.LogInfo("ACI APIC", map[string]interface{}{"Username": a.ACI[i].username}, false)
		a.LogInfo("ACI APIC", map[string]interface{}{"Password": a.ACI[i].password}, false)
	}
}

func (a *Application) updateACI(ip, username, password string) bool {
	for i := 0; i < len(a.ACI); i++ {
		if a.ACI[i].ip == as.ToString(ip) {
			a.ACI[i].username = username
			a.ACI[i].password = a.EncryptPassword(password)
		}
	}
	return true
}

func (a *Application) updateACISystem(ip, username, password string) {
	if a.checkACIExists(ip) {
		if a.updateACI(ip, username, password) {
			a.saveConfig()
			a.LogInfo("Update to ACI APIC has been completed successfully.", map[string]interface{}{"IP": ip, "Username": username}, false)
		} else {
			a.LogInfo("ACI APIC could not be updated.", map[string]interface{}{"IP": ip, "Username": username}, false)
		}
	} else {
		a.LogInfo("ACI APIC does not exsist and can therefore not be updated.", map[string]interface{}{"IP": ip, "Username": username}, false)
	}
}

func (a *Application) checkACIExists(ip string) bool {
	a.Log("Searching for ACI APIC in config file", map[string]interface{}{"IP": ip}, true)
	if a.Config.IsSet("aci.systems") {
		a.getAllSystems()
		for i := 0; i < len(a.ACI); i++ {
			if strings.TrimSpace(a.ACI[i].ip) == strings.TrimSpace(ip) {
				return true
			}
		}
		return false
	}
	return false
}

func (a *Application) getAllSystems() {
	tmp := as.ToSlice(a.Config.Get("aci.systems"))
	a.Log("Located ACI APIC in the config file", map[string]interface{}{"Systems": len(tmp)}, true)
	a.readSystems(tmp)
}

func (a *Application) readSystems(acii []interface{}) bool {
	a.ACI = nil
	for i := 0; i < len(acii); i++ {
		var newlist map[string]string
		newlist = as.ToStringMapString(acii[i])
		tmp := ACISystemInfo{}
		tmp.ip = newlist["url"]
		tmp.username = newlist["username"]
		tmp.password = newlist["password"]
		a.ACI = append(a.ACI, tmp)
	}
	return true
}
