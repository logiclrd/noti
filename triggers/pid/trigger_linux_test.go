package pid

import "testing"

func TestProcName(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Name:	nvme\nUmask:	0000\nState:	S (sleeping)\nTgid:	123\nNgid:	0\nPid:	123\nPPid:	2\nTracerPid:	0\nUid:	0	0	0	0\nGid:	0	0	0	0\nFDSize:	64\nGroups:	\nNStgid:	123\nNSpid:	123\nNSpgid:	0\nNSsid:	0\nThreads:	1\nSigQ:	0/62396\nSigPnd:	0000000000000000\nShdPnd:	0000000000000000\nSigBlk:	0000000000000000\nSigIgn:	ffffffffffffffff\nSigCgt:	0000000000000000\nCapInh:	0000000000000000\nCapPrm:	0000003fffffffff\nCapEff:	0000003fffffffff\nCapBnd:	0000003fffffffff\nCapAmb:	0000000000000000\nSeccomp:	0\nCpus_allowed:	f\nCpus_allowed_list:	0-3\nMems_allowed:	00000000,00000001\nMems_allowed_list:	0\nvoluntary_ctxt_switches:	2\nnonvoluntary_ctxt_switches:	0\n",
			"nvme"},
		{"Name: nvme\n", "nvme"},
		{"Name: fooBAR\n", "fooBAR"},
		{"Name: foo-bar\n", "foo-bar"},
		{"Name: foo_bar\n", "foo_bar"},
		{"Name: foo bar\n", "foo bar"},
		{"Name: foobar Umask: 0000", "foobar"},
		{"", ""},
	}

	for i, c := range cases {
		got := procName(c.in)

		if got != c.want {
			t.Errorf("%d: got: %q; want: %q", i, got, c.want)
		}
	}
}
