package auditd

import (
	"fmt"
	"testing"
)

/*
@Author: OvO
@Date: 2024/8/21 14:38
*/

func TestAudit(t *testing.T) {
	mapA := parseAudit("type=CRED_ACQ msg=audit(1722813841.098:2810352): pid=110205 uid=0 auid=4294967295 ses=4294967295 msg='op=PAM:setcred grantors=pam_env,pam_unix acct=\"root\" exe=\"/usr/sbin/crond\" hostname=? addr=? terminal=cron res=success'")
	for s, s2 := range mapA {
		fmt.Println(s, s2)
	}
}
