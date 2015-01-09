package junos

var RPCCommand = map[string]string{
	"lock":                             "<rpc><lock><target><candidate/></target></lock></rpc>",
	"unlock":                           "<rpc><unlock><target><candidate/></target></unlock></rpc>",
    "command": "<rpc><command format=\"text\">%s</command></rpc>",
	"get-rescue-information":           "<rpc><get-rescue-information><format>text</format></get-rescue-information></rpc>",
	"get-rollback-information":         "<rpc><get-rollback-information><rollback>%d</rollback><format>text</format></get-rollback-information></rpc>",
	"get-rollback-information-compare": "<rpc><get-rollback-information><rollback>0</rollback><compare>%d</compare><format>text</format></get-rollback-information></rpc>",
}
