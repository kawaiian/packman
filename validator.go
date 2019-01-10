package packman

import (
	"errors"
	"regexp"
	"strings"
)

type PkgRequest struct {
	req     string
	pkg     string
	depList []string
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

	err := parseRequest(msgParts[0])
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
		req:     msgParts[0],
		pkg:     msgParts[1],
		depList: pkgDepList,
	}

	return pkgReq, nil
}

// ParseRequest Validates request type is one of INDEX, QUERY, or REMOVE
func parseRequest(req string) error {
	validReqs := map[string]struct{}{"INDEX": {}, "QUERY": {}, "REMOVE": {}}

	_, reqOk := validReqs[req]
	if !reqOk {
		return errors.New("invalid request type: not INDEX, QUERY, or REMOVE")
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
