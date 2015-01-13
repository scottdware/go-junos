package junos

// rpcCommand lists the commands that will be called.
var rpcCommand = map[string]string{
	"command":                          "<rpc><command format=\"text\">%s</command></rpc>",
	"command-xml":                      "<rpc><command format=\"xml\">%s</command></rpc>",
	"commit":                           "<rpc><commit-configuration/></rpc>",
	"commit-at":                        "<rpc><commit-configuration><at-time>%s</at-time></commit-configuration></rpc>",
	"commit-check":                     "<rpc><commit-configuration><check/></commit-configuration></rpc>",
	"commit-confirm":                   "<rpc><commit-configuration><confirmed/><confirm-timeout>%d</confirm-timeout></commit-configuration></rpc>",
	"facts-re":                         "<rpc><get-route-engine-information/></rpc>",
	"facts-chassis":                    "<rpc><get-chassis-inventory/></rpc>",
	"load-config-local-set":            "<rpc><load-configuration action=\"set\" format=\"text\"><configuration-set>%s</configuration-set></load-configuration></rpc>",
	"load-config-local-text":           "<rpc><load-configuration format=\"text\"><configuration-text>%s</configuration-text></load-configuration></rpc>",
	"load-config-local-xml":            "<rpc><load-configuration format=\"xml\"><configuration>%s</configuration></load-configuration></rpc>",
	"load-config-url-set":              "<rpc><load-configuration action=\"set\" format=\"text\" url=\"%s\"/></rpc>",
	"load-config-url-text":             "<rpc><load-configuration format=\"text\" url=\"%s\"/></rpc>",
	"load-config-url-xml":              "<rpc><load-configuration format=\"xml\" url=\"%s\"/></rpc>",
	"get-rescue-information":           "<rpc><get-rescue-information><format>text</format></get-rescue-information></rpc>",
	"get-rollback-information":         "<rpc><get-rollback-information><rollback>%d</rollback><format>text</format></get-rollback-information></rpc>",
	"get-rollback-information-compare": "<rpc><get-rollback-information><rollback>0</rollback><compare>%d</compare><format>text</format></get-rollback-information></rpc>",
	"lock":            "<rpc><lock><target><candidate/></target></lock></rpc>",
	"rescue-config":   "<rpc><load-configuration rescue=\"rescue\"/></rpc>",
	"rollback-config": "<rpc><load-configuration rollback=\"%d\"/></rpc>",
	"software":        "<rpc><get-software-information/></rpc>",
	"unlock":          "<rpc><unlock><target><candidate/></target></unlock></rpc>",
}
