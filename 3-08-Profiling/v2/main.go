package main

func main() {
	profile, err := os.Create("cpu.profile")
}
if err != nil {
		log.Fatal(err)
	}
	defer profile.Close()

	if err := pprof.StartCPUProfile(profile); err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	filter.Grep()
	filter.Match()

