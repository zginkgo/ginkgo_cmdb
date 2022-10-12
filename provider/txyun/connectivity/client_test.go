package connectivity

import (
	"testing"
)

func TestTencentCloudClient(t *testing.T) {
	conn := C()
	if err := conn.Check(); err != nil {
		t.Fatal(err)
	}
	t.Log(conn.AccountID())
	t.Log(conn.cvmConn.)
}

func init() {
	err := LoadClientFromEnv()
	if err != nil {
		panic(err)
	}
}
