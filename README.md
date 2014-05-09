GOSQM - SQM Toolchain
========================

gosqm is a toolchain (written in golang) for the sqm format used to represent mission files of the Arma Series published by Bohemia Interactive.
It provides both a low-level representation of the format to directly manipulate classes and attributes (gosqm/sqm) as well as a high-level representation to manipulate e.g. groups, units and vehicles.

Usage (Highlevel)
-----

See [mission.go](mission.go) for entities.

	f, _ := os.Open("mission.sqm")
	defer f.Close()
	dec := gosqm.NewDecoder(f)
	missionFile, err := dec.Decode()
	// Manipulate missionFile
	fo, _ := os.Create("mission.sqm.out")
	defer fo.Close()
	enc := gosqm.NewEncoder(buffer)
	err := enc.Encode(missionFile)

Usage (Lowlevel)
-----

See [sqm/parserentities.go](sqm/parserentities.go) for entities.

	b, _:= ioutil.ReadFile("mission.sqm")
	parser := sqm.MakeParser(bufstr)
	class, err := parser.Run()
	//Manipulate class
	fo, _ := os.Create("mission.sqm.out")
	enc := sqm.NewEncoder(fo)
	err = enc.Encode(class)

Stability
-----

	Alpha

Currently safely supports Arma2 mission files. Interfaces are subject to change.

Benchmarks
-----
A full decode and encode cycle in <6ms on a medium machine.

	BenchmarkFullDecode	         500	   4295337 ns/op
	BenchmarkFullEncode	        2000	    999859 ns/op
	BenchmarkMissionEncode	    5000	    508818 ns/op
	BenchmarkMissionDecode	   10000	    113710 ns/op
	BenchmarkFullEncodeDecode	 500	   5345543 ns/op

Issues
-----

There's a lot work left to make the high-level encoder and decoder intolerant to incorrect inputs. It should be fine if you know what you're doing and have a basic understanding of the format.

Contribution
-----

Feel free to make a pull request. For bigger changes create a issue first to discuss about it.


License
-----

See [LICENSE](LICENSE) file.
