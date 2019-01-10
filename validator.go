package packman

import (
	"errors"
	"regexp"
	"strings"
)

// PkgRequest encapsulates the three main fields of a client's request to the packman service
type PkgRequest struct {
	Cmd     string
	Pkg     string
	DepList []string
}

// ParsePkgRequest validates a package request by decomposing it and
// passing those components to helper functions.
func ParsePkgRequest(msg string) (PkgRequest, error) {
	var pkgReq PkgRequest

	// Requests are expected to be "<command>|<pkg>|<dependency>\n"
	msgParts := strings.Split(msg, "|")
	if len(msgParts) < 3 {
		return pkgReq, errors.New("invalid request format")
	}

	err := parseCmd(msgParts[0])
	if err != nil {
		return pkgReq, err
	}

	dListString := strings.TrimSuffix(msgParts[2], "\n")
	pkgDepList := strings.Split(dListString, ",")
	err = parsePkgNames(msgParts[1], pkgDepList)
	if err != nil {
		return pkgReq, err
	}

	// Otherwise, return a request struct
	pkgReq = PkgRequest{
		Cmd:     msgParts[0],
		Pkg:     msgParts[1],
		DepList: pkgDepList,
	}

	return pkgReq, nil
}

// parseCmd Validates request type is one of INDEX, QUERY, or REMOVE
func parseCmd(cmd string) error {
	validCmds := map[string]struct{}{"INDEX": {}, "QUERY": {}, "REMOVE": {}}

	_, cmdOk := validCmds[cmd]
	if !cmdOk {
		return errors.New("invalid command: not INDEX, QUERY, or REMOVE")
	}

	return nil
}

// ParsePkgNames validates the format of package and dependency names
func parsePkgNames(pkgName string, pkgDepList []string) error {

	fullPkgList := make([]string, 1)
	fullPkgList[0] = pkgName
	if pkgDepList[0] != "" {
		fullPkgList = append(fullPkgList, pkgDepList...)
	}

	validName := false
	validator := regexp.MustCompile(`^[a-zA-Z0-9]+[a-zA-Z0-9\-\_\+]*[a-zA-Z0-9\+]*$`)
	for _, name := range fullPkgList {
		validName = validator.MatchString(name)
		if !validName {
			return errors.New("invalid package name")
		}
	}

	return nil
}
