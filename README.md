<h1> Motivation: </h1>
hashcomparefiles is a robust CLI tool that computes file hashes to identify duplicate files regardless of filename, using SHA-256, and presents them to the user. This tool also makes it easy to selectively delete duplicate files, move them to trash, or output a list of all duplicate files with their filesize.
It runs in a multi-pass method, minimizing disk hits and improving efficiency. Initially, I developed this to help prune down a very messy, overly redundant NAS, where I had the same files haphazzardly thrown in a mess of folders, and needed a tool to help thin things out to help organize things.

### Notes: ###

* symlinks and empty files are ignored
* sub-folders are automatically walked and included
* compatible with piping in lists of folders/filenames for more customization
* The filesize is included for reference, and for the remote chance of hash collision.
- This tool was developed and tested in a Linux environment (Pop!_OS 24.04 LTS), but I would like it to eventually fully operate in Mac and Windows environments as well. In paricular, the trashing processes in Mac and Windows is incompatible, and the input when selectively removing files currently uses TTY (which likely won't work in Windows, but is currently untested).




## Quick Start: 

### Download compiled binaries manually, under "Releases" here: 
https://github.com/gainax2k1/hashcomparefiles/releases
-Then, from the command line, type:
```python
go install hashcomparefiles-(version)
```

### Or, using the Go CLI:
```python
go install github.com/gainax2k1/hashcomparefiles@latest
```

### Or, if you want to build your own:

- Install Go, if not already installed, instuctions here:
https://go.dev/doc/install

- Download repo
- In root folder of repo, run:
```python
go build
go install
```
# Usage:

```python
hashcomparefiles (filename/directory)
```
- returns hash value of a single file or through directory and sub-directories. Won't display lists of duplicate files without -v.

```python
hashcomparefiles -remove (directory)
```
- Processes all files, and goes through each group of hash-match files, allowing per-file deletion or trashing if desired. *Only fully supports FreeDesktop spec on primary drive in Linux based systems.* For non-Linux systems, or if running on an external mount in Linux, a folder is created in the working directory, and files are move into that, with corresponding .trashinfo files being created to record original file location. For this reason, it's highly recommended to run from the drive where files are stored. I would like to improve this in the future. 


```python
hashcomparefiles -log (directory/logfilename) ...
```
- Creates a log file in the given directory/logfilename, default is current working directory.

```python
hashcomparefiles -v (filename/directory)
```
- Verbose flag. This will output the final list of duplicates to terminal. By default (without flag), only progress will be output to terminal. Usefull if you want to log to file, but also want it output to screen.

```python
hashcomparefiles -min (integer) -max (integer) ...
```
- Set minimum and/or maximum filesizes (in bytes) to process.

```python
hashcomparefiles --help
```
- Shows list of available flags and descriptions

```python
cat (filename) | hashcomparefiles -(flag)
```
- Pipe in list of files and or folders to compare against each other. Flags maintain functionality.





# 🤝 Contributing
<h2>Submit a pull request</h2>

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.

# Examples:
<h2> Small run, no flags</h2>

```python
$ hashcomparefiles testdata/
Filecount after pass (1/3): 18
Filecount after pass (2/3): 12
Filecount after pass (3/3): 12
Groups of duplicates after shrink: 4


```

<h2>Full output from simple run with log output.</h2>

```python
hashcomparefiles -log testdata.log testdata/
Filecount after pass (1/3): 18
Filecount after pass (2/3): 12
Filecount after pass (3/3): 12
Groups of duplicates after shrink: 4
gainax2k1@pop-os-thinkpad:~/Documents/workspace/hashcomparefiles$ more testdata.log 
{"time":"2026-06-20T15:04:33.075265481-07:00","level":"INFO","msg":"pass_complete","pass":1,"count":18}
{"time":"2026-06-20T15:04:33.075749506-07:00","level":"INFO","msg":"pass_complete","pass":2,"count":12}
{"time":"2026-06-20T15:04:33.07581558-07:00","level":"INFO","msg":"pass_complete","pass":3,"count":12}
{"time":"2026-06-20T15:04:33.075825434-07:00","level":"INFO","msg":"pass_complete","pass":4,"count":4}
{"time":"2026-06-20T15:04:33.075841615-07:00","level":"INFO","msg":"duplicates_found","hash":"6f430d148a85e1475301f9bd44463cc8dc69bbc1a0e059eb7c7314734e8db6dd","foun
d":3,"files":[{"FilePath":"/home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileD.txt","FileSize":30},{"FilePath":"/h
ome/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileD.txt","FileSize":30},{"FilePath":"/home/gainax2k1/Documents/workspace/hashcomparef
iles/testdata/testFileDDup.txt","FileSize":30}]}
{"time":"2026-06-20T15:04:33.075931684-07:00","level":"INFO","msg":"duplicates_found","hash":"7368ac39295432a153b1532cacf30c1a4b55cc94c246d6cce820a42c06ff8c2f","foun
d":3,"files":[{"FilePath":"/home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileBDup.txt","FileSize":20},{"FilePath":"/home/gainax2k1/
Documents/workspace/hashcomparefiles/testdata/testFileB.txt","FileSize":20},{"FilePath":"/home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/
testFolderNested/testFileBDup.txt","FileSize":20}]}
{"time":"2026-06-20T15:04:33.075938841-07:00","level":"INFO","msg":"duplicates_found","hash":"a4978f74fe60dbc373e48f0486d767c8d866a8f94a45c661acf812e44d978a38","foun
d":3,"files":[{"FilePath":"/home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt","FileSize":22},{"FilePath":"/home/gainax2k1/
Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileCDup.txt","FileSize":22},{"FilePath":"/home/gainax2k1/Documents/workspace/hashco
mparefiles/testdata/testFileC.txt","FileSize":22}]}
{"time":"2026-06-20T15:04:33.075944135-07:00","level":"INFO","msg":"duplicates_found","hash":"c0f5efbef0fe98aa90619444250b1a5eb23158d6686f0b190838f3d544ec85b9","foun
d":3,"files":[{"FilePath":"/home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileADup.txt","FileSize":10},{"FilePath":
"/home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileADup.txt","FileSize":10},{"FilePath":"/home/gainax2k1/Documents/workspace/hashco
mparefiles/testdata/testFileA.txt","FileSize":10}]}
{"time":"2026-06-20T15:04:33.075949224-07:00","level":"INFO","msg":"scan_complete","total_duplicates":12}
gainax2k1@pop-os-thinkpad:~/Documents/workspace/hashcomparefiles$ 



```


<h2> Running with log on my home directory, note that the full run took 29 seconds.</h2>

```python
$ hashcomparefiles -log homeDir.log /home/gainax2k1/
 | Files processed: 695200 
Filecount after pass (1/3): 695260
Filecount after pass (2/3): 653967
Filecount after pass (3/3): 460684
Groups of duplicates after shrink: 130612

$  more homeDir.log 
{"time":"2026-06-20T15:05:54.457342887-07:00","level":"INFO","msg":"pass_complete","pass":1,"count":695260}
{"time":"2026-06-20T15:05:58.066783658-07:00","level":"INFO","msg":"pass_complete","pass":2,"count":653967}
{"time":"2026-06-20T15:06:23.065206981-07:00","level":"INFO","msg":"pass_complete","pass":3,"count":460684}
{"time":"2026-06-20T15:06:23.104763733-07:00","level":"INFO","msg":"pass_complete","pass":4,"count":130612}
{"time":"2026-06-20T15:06:23.104790161-07:00","level":"INFO","msg":"duplicates_found","hash":"5b3c161c054517a3cada248ab392c3ee4e05458137381e03b85c4ba5a1a58622","foun
d":2,"files":[{"FilePath":"/home/gainax2k1/.local/share/flatpak/runtime/org.kde.Platform.Locale/x86_64/6.9/8c9622d6e73f538cc116b94395f663791902f5296f63dca34263f18387
ac8a62/files/sl/share/sl/LC_MESSAGES/kio6_man.mo","FileSize":4130},{"FilePath":"/home/gainax2k1/.local/share/flatpak/repo/objects/ac/0406c4708e37dc55988e3b4c8a3991c5
d8168fcea80366b471c82cf46ed96d.file","FileSize":4130}]}

< cut for ammount of output >

{"time":"2026-06-20T15:06:23.68418657-07:00","level":"INFO","msg":"duplicates_found","hash":"cfdf828aa34498dd358b86e7711abed4a6302eb1d09b0ade74951d51bedcd8b3","found":2,"files":[{"FilePath":"/home/gainax2k1/.local/opt/go-v1.25.4/src/crypto/internal/fips140/aes/gcm/gcm_noasm.go","FileSize":572},{"FilePath":"/home/gainax2k1/.local/opt/go-bin-v1.25.4/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.0.linux-amd64/src/crypto/internal/fips140/aes/gcm/gcm_noasm.go","FileSize":572}]}
{"time":"2026-06-20T15:06:23.684189985-07:00","level":"INFO","msg":"scan_complete","total_duplicates":456353}
$

```

<h2> Running with min/max flags on home: </h2>

```python
$ hashcomparefiles -min 2048 -max 8128 -log minMaxHome.log /home/gainax2k1/
 | Files processed: 695200 
Filecount after pass (1/3): 695262
Filecount after pass (2/3): 148145
Filecount after pass (3/3): 105397
Groups of duplicates after shrink: 31862
gainax2k1@pop-os-thinkpad:~/Documents/workspace/hashcomparefiles$ more minMaxHome.log 
{"time":"2026-06-20T15:17:58.924086453-07:00","level":"INFO","msg":"pass_complete","pass":1,"count":695262}
{"time":"2026-06-20T15:17:59.792614463-07:00","level":"INFO","msg":"pass_complete","pass":2,"count":148145}
{"time":"2026-06-20T15:18:00.052912714-07:00","level":"INFO","msg":"pass_complete","pass":3,"count":105397}
{"time":"2026-06-20T15:18:00.058109563-07:00","level":"INFO","msg":"pass_complete","pass":4,"count":31862}
{"time":"2026-06-20T15:18:00.058142505-07:00","level":"INFO","msg":"duplicates_found","hash":"a427d1e7378772a69d44e27f733acf3069043dc1d94419c70981d58b84827dfc","foun
d":3,"files":[{"FilePath":"/home/gainax2k1/.local/share/flatpak/repo/objects/e7/c535e842eb5aec2c4c924e6552c7d06873fcc69f9e48f710e10887e1da8122.file","FileSize":7073}
,{"FilePath":"/home/gainax2k1/.local/share/flatpak/runtime/org.freedesktop.Platform.Locale/x86_64/25.08/40d64910991e3e290580e161c3719128cfe8227fbcd8a260e540c5239517c
942/files/fi/share/fi/LC_MESSAGES/iso_4217.mo","FileSize":7073},{"FilePath":"/home/gainax2k1/.local/share/flatpak/runtime/org.gnome.Platform.Locale/x86_64/48/8987d72
0c70ceca711f45c6caf28a479b3a805a25ab7387ccc151085909b4db2/files/fi/share/fi/LC_MESSAGES/iso_4217.mo","FileSize":7073}]}

< cut for ammount of output >

```


<hr>
<h2> "time -v" output run of home directory, uncached:</h2>

```
$ /usr/bin/time -v hashcomparefiles /home/gainax2k1/
 / Files processed: 696100 
Filecount after pass (1/3): 696142
Filecount after pass (2/3): 654841
Filecount after pass (3/3): 461121
Groups of duplicates after shrink: 130634
	Command being timed: "hashcomparefiles /home/gainax2k1/"
	User time (seconds): 76.53
	System time (seconds): 41.90
	Percent of CPU this job got: 343%
	Elapsed (wall clock) time (h:mm:ss or m:ss): 0:34.47
	Average shared text size (kbytes): 0
	Average unshared data size (kbytes): 0
	Average stack size (kbytes): 0
	Average total size (kbytes): 0
	Maximum resident set size (kbytes): 597200
	Average resident set size (kbytes): 0
	Major (requiring I/O) page faults: 0
	Minor (reclaiming a frame) page faults: 217270
	Voluntary context switches: 1574263
	Involuntary context switches: 36806
	Swaps: 0
	File system inputs: 74629280
	File system outputs: 205024
	Socket messages sent: 0
	Socket messages received: 0
	Signals delivered: 0
	Page size (bytes): 4096
	Exit status: 0


```

That should be everything! Any issues with the project, feel free to reach out to me. 
Thanks, and have a day. =^.^=
