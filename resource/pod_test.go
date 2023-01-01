package resource

import "testing"

func TestRuncmd(t *testing.T) {
	pf := NewPodFactor(kubectlPodCmd)
	scan, _ := pf.runCmd()
	for scan.Scan() {
		t.Log(scan.Text())
	}
}
