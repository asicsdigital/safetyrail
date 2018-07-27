package util

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"

	jww "github.com/spf13/jwalterweatherman"
)

type XMLRoutingRules struct {
	XMLName      xml.Name         `xml:"RoutingRules"`
	RoutingRules []XMLRoutingRule `xml:"RoutingRule"`
}

type XMLRoutingRule struct {
	XMLName   xml.Name     `xml:"RoutingRule"`
	Condition XMLCondition `xml:"Condition"`
	Redirect  XMLRedirect  `xml:"Redirect"`
}

type XMLCondition struct {
	XMLName         xml.Name `xml:"Condition"`
	KeyPrefixEquals string   `xml:"KeyPrefixEquals"`
}

type XMLRedirect struct {
	XMLName              xml.Name `xml:"Redirect"`
	Protocol             string   `xml:"Protocol"`
	HostName             string   `xml:"HostName"`
	ReplaceKeyPrefixWith string   `xml:"ReplaceKeyPrefixWith"`
}

func (r *XMLRoutingRules) Parse(infile string) error {
	b, err := ioutil.ReadFile(infile)

	if err != nil {
		jww.ERROR.Println("unable to open rules file %s: %s", infile, err)
	} else {
		err = xml.Unmarshal(b, r)
	}

	return err
}

func (r *XMLRoutingRule) Dump() string {
	b, err := xml.MarshalIndent(r, "", "  ")

	if err != nil {
		jww.ERROR.Panic(err)
	}

	return string(b)
}

func (r *XMLRoutingRule) ModRewrite() string {
	sourcePath := r.Condition.KeyPrefixEquals
	protocol := r.Redirect.Protocol
	hostname := r.Redirect.HostName
	destPath := r.Redirect.ReplaceKeyPrefixWith

	return fmt.Sprintf("^/%s(.*) %s://%s/%s$1 [R,L]", sourcePath, protocol, hostname, destPath)
}
