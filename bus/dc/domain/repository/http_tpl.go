package repository

var loginTpl string = `mutation login($pwd: String) {
login(
role: TravelAgent 
login: "HKGKLOOK-GQL" 
password: $pwd
sourceCode:"OPENTRAVEL"
) {
token,
agency {
name,
code,
id
}
role,
expiration {
timeout,
inactivityTimeout,
expires,
autoExtend
}
versionInfo
}
}
`

type LoginRequest struct {
	Pwd string `json:"pwd"`
}
type LoginResponse struct {
	Data struct {
		Login struct {
			Token string `json:"token"`
		} `json:"login"`
	} `json:"data"`
}

var scheduleTpl string = `query AvailableVoyages($datefrom: Date,$dateto: Date,$ports: [String],$ships: [String],$minDur: Period,$maxDur: Period) {
availableVoyages(
params: {
availability: OK
startDateRange: { 
from: $datefrom, 
to: $dateto
}
departurePortKeys: $ports
shipKeys: $ships
lengthRange: {
from: $minDur
to: $maxDur
}
} 
) {
reference
pkg {
code
name
id
key
description
sailDays
}
sail {
ship {
code
name
}
from {
dateTime
port {
code
name
}
sailRefID
}
to {
dateTime
port {
code
name
}
sailRefID
}}
sailActivities {
dateTime
portCode
port {
name
}
mayEmbark
mayDisembark
}
availableCategories {
cabinCategory {
description
id
}
}
}
}
`

type DcError struct {
	Errors []struct {
		Message   string `json:"message"`
		Locations []struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"locations"`
		Extensions struct {
			Classification string `json:"classification"`
		} `json:"extensions"`
	} `json:"errors"`
}

type ScheduleReqest struct {
	Datefrom string   `json:"datefrom"`
	Dateto   string   `json:"dateto"`
	Ports    []string `json:"ports"`
	Ships    []string `json:"ships"`
	MinDur   int      `json:"minDur"`
	MaxDur   int      `json:"maxDur"`
}

type AvailableVoyageRow struct {
	Reference string `json:"reference"`
	Pkg       struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		Id          string `json:"id"`
		Key         string `json:"key"`
		Description string `json:"description"`
		SailDays    int    `json:"sailDays"`
	} `json:"pkg"`
	Sail struct {
		Ship struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"ship"`
		From struct {
			DateTime string `json:"dateTime"`
			Port     struct {
				Code string `json:"code"`
				Name string `json:"name"`
			} `json:"port"`
			SailRefID int `json:"sailRefID"`
		} `json:"from"`
		To struct {
			DateTime string `json:"dateTime"`
			Port     struct {
				Code string `json:"code"`
				Name string `json:"name"`
			} `json:"port"`
			SailRefID int `json:"sailRefID"`
		} `json:"to"`
	} `json:"sail"`
	SailActivities []struct {
		DateTime string  `json:"dateTime"`
		PortCode *string `json:"portCode"`
		Port     *struct {
			Name string `json:"name"`
		} `json:"port"`
		MayEmbark    bool `json:"mayEmbark"`
		MayDisembark bool `json:"mayDisembark"`
	} `json:"sailActivities"`
	AvailableCategories []struct {
		CabinCategory struct {
			Description string `json:"description"`
			Id          string `json:"id"`
		} `json:"cabinCategory"`
	} `json:"availableCategories"`
}

type ScheduleResponse struct {
	DcError
	Data struct {
		AvailableVoyages []AvailableVoyageRow `json:"availableVoyages"`
	} `json:"data"`
}