package util

import (
	"encoding/xml"
)

type RoutingRules []RoutingRule

type RoutingRule struct {
	Condition Condition `xml:"Condition"`
}

func main() {
	fmt.Println("vim-go")
}
