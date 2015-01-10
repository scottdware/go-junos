package junos

// rpcCommand lists the commands that will be called.
var rpcCommand = map[string]string{
	"command":                          "<rpc><command format=\"text\">%s</command></rpc>",
	"command-xml":                      "<rpc><command format=\"xml\">%s</command></rpc>",
	"commit":                           "<rpc><commit-configuration/></rpc>",
	"get-rescue-information":           "<rpc><get-rescue-information><format>text</format></get-rescue-information></rpc>",
	"get-rollback-information":         "<rpc><get-rollback-information><rollback>%d</rollback><format>text</format></get-rollback-information></rpc>",
	"get-rollback-information-compare": "<rpc><get-rollback-information><rollback>0</rollback><compare>%d</compare><format>text</format></get-rollback-information></rpc>",
	"lock":            "<rpc><lock><target><candidate/></target></lock></rpc>",
	"rescue-config":   "<rpc><load-configuration rescue=\"rescue\"/></rpc>",
	"rollback-config": "<rpc><load-configuration rollback=\"%d\"/></rpc>",
	"unlock":          "<rpc><unlock><target><candidate/></target></unlock></rpc>",
}
