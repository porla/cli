package ninja

var pool = []string{
	"archlinux-2022.10.01-x86_64.iso",
	"archlinux-2022.11.01-x86_64.iso",
	"archlinux-2022.12.01-x86_64.iso",
	"debian-11.6.0-amd64-netinst.iso",
	"debian-edu-11.6.0-amd64-netinst.iso",
	"debian-mac-11.6.0-amd64-netinst.iso",
	"ubuntu-20.04.5-desktop-amd64.iso",
	"ubuntu-20.04.5-live-server-amd64.iso",
	"ubuntu-22.04.1-desktop-amd64.iso",
	"ubuntu-22.04.1-live-server-amd64.iso",
}

func RandomLinuxTorrent(index int) string {
	return pool[index%len(pool)]
}
