package gpucard

import "testing"

var (
	iris   = Adapter{Name: "Intel(R) Iris(R) Xe Graphics", DedicatedMB: 128, Integrated: true, Hardware: true}
	apu    = Adapter{Name: "AMD Radeon(TM) Graphics", DedicatedMB: 8192, Integrated: true, Hardware: true}
	rtx    = Adapter{Name: "NVIDIA GeForce RTX 3070 Ti Laptop GPU", DedicatedMB: 8017, Hardware: true}
	gtx2   = Adapter{Name: "NVIDIA GeForce GTX 1050", DedicatedMB: 2048, Hardware: true}
	gtx4   = Adapter{Name: "NVIDIA GeForce GTX 1650", DedicatedMB: 3976, Hardware: true}
	exact4 = Adapter{Name: "AMD Radeon RX 6500 XT", DedicatedMB: 4096, Hardware: true}
	basic  = Adapter{Name: "Microsoft Basic Render Driver", DedicatedMB: 0, Hardware: false}
	soft   = Adapter{Name: "Software card", DedicatedMB: 16384, Hardware: false}
)

func TestTheCardRule(t *testing.T) {
	cases := []struct {
		name   string
		list   []Adapter
		fits   bool
		pick   string
		reason string
		mb     int
	}{
		{"integrated only", []Adapter{iris, basic}, false, "", ReasonIntegrated, 0},
		{"integrated with a large share", []Adapter{apu}, false, "", ReasonIntegrated, 0},
		{"2 GB", []Adapter{iris, gtx2}, false, "", ReasonSmall, 2048},
		{"4 GB", []Adapter{exact4}, true, exact4.Name, "", 4096},
		{"4 GB as reported by the driver", []Adapter{iris, gtx4}, true, gtx4.Name, "", 3976},
		{"two adapters", []Adapter{iris, rtx, basic}, true, rtx.Name, "", 8017},
		{"two adapters, the discrete one first", []Adapter{rtx, iris}, true, rtx.Name, "", 8017},
		{"two cards that fit", []Adapter{exact4, rtx}, true, rtx.Name, "", 8017},
		{"software only", []Adapter{basic, soft}, false, "", ReasonNone, 0},
		{"nothing", nil, false, "", ReasonNone, 0},
	}
	for _, c := range cases {
		got, why := Pick(c.list)
		if (why == "") != c.fits || got.Name != c.pick {
			t.Errorf("%s: picked %+v, why %q", c.name, got, why)
		}
		code, mb := Reason(c.list)
		if code != c.reason || mb != c.mb {
			t.Errorf("%s: reason %q %d", c.name, code, mb)
		}
	}
	if Fits(gtx2) || !Fits(exact4) || Fits(apu) || Fits(soft) {
		t.Fatal("Fits")
	}
}

const listOut = "Available devices:\r\n  Vulkan0: Intel(R) Iris(R) Xe Graphics (8022 MiB, 7365 MiB free)\r\n  Vulkan1: NVIDIA GeForce RTX 3070 Ti Laptop GPU (8017 MiB, 7240 MiB free)\r\n"

func TestTheDeviceIsPinnedByName(t *testing.T) {
	devs := ParseDevices(listOut)
	if len(devs) != 2 || devs[0].ID != "Vulkan0" || devs[1].Name != "NVIDIA GeForce RTX 3070 Ti Laptop GPU" || devs[1].TotalMB != 8017 || devs[1].FreeMB != 7240 {
		t.Fatalf("devices %+v", devs)
	}
	if d, ok := Match(devs, rtx.Name); !ok || d.ID != "Vulkan1" {
		t.Fatalf("rtx: %+v %v", d, ok)
	}
	if d, ok := Match(devs, "Intel Iris Xe Graphics"); !ok || d.ID != "Vulkan0" {
		t.Fatalf("iris without the marks: %+v %v", d, ok)
	}
	if d, ok := Match(devs, "NVIDIA GeForce RTX 3070 Ti"); !ok || d.ID != "Vulkan1" {
		t.Fatalf("a shorter name: %+v %v", d, ok)
	}
	if _, ok := Match(devs, "AMD Radeon RX 7600"); ok {
		t.Fatal("a card that is not there matched")
	}
	if _, ok := Match(devs, ""); ok {
		t.Fatal("an empty name matched")
	}
	two := ParseDevices("  Vulkan0: NVIDIA GeForce RTX 4090 (24564 MiB, 20000 MiB free)\n  Vulkan1: NVIDIA GeForce RTX 4090 D (24564 MiB, 24000 MiB free)\n")
	if d, ok := Match(two, "NVIDIA GeForce RTX 4090"); !ok || d.ID != "Vulkan0" {
		t.Fatalf("the exact name wins: %+v %v", d, ok)
	}
	if d, ok := Match(two, "RTX 4090"); ok {
		t.Fatalf("an ambiguous name matched %+v", d)
	}
	if got := ParseDevices("Available devices:\n"); len(got) != 0 {
		t.Fatalf("no devices: %+v", got)
	}
}
