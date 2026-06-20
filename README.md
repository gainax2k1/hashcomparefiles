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
Filecount after pass (1/3): 14
Filecount after pass (2/3): 8
Filecount after pass (3/3): 8
Groups of duplicates after shrink: 3


```

<h2>Full output from simple run with log output.</h2>

```python
$ hashcomparefiles -log testdata.log testdata/
Filecount after pass (1/3): 14
Filecount after pass (2/3): 8
Filecount after pass (3/3): 8
Groups of duplicates after shrink: 3
$ more testdata.log 
2026/03/12 12:25:34 (Start)
2026/03/12 12:25:34 Filecount after pass (1/3): 14
2026/03/12 12:25:34 Filecount after pass (2/3): 8
2026/03/12 12:25:34 Filecount after pass (3/3): 8
2026/03/12 12:25:34 Groups of duplicates after shrink: 3
2026/03/12 12:25:34 Files with hash: 7368ac39295432a153b1532cacf30c1a4b55cc94c246d6cce820a42c06ff8c2f
2026/03/12 12:25:34  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileB.txt size: 20
2026/03/12 12:25:34  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileBDup.txt size: 20
2026/03/12 12:25:34  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileBDup.txt size: 20
2026/03/12 12:25:34  -- Duplicates: 3
2026/03/12 12:25:34 Files with hash: 6f430d148a85e1475301f9bd44463cc8dc69bbc1a0e059eb7c7314734e8db6dd
2026/03/12 12:25:34  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileDDup.txt size: 30
2026/03/12 12:25:34  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileD.txt size: 30
2026/03/12 12:25:34  -- Duplicates: 2
2026/03/12 12:25:34 Files with hash: a4978f74fe60dbc373e48f0486d767c8d866a8f94a45c661acf812e44d978a38
2026/03/12 12:25:34  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileC.txt size: 22
2026/03/12 12:25:34  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt size: 22
2026/03/12 12:25:34  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileCDup.txt size: 22
2026/03/12 12:25:34  -- Duplicates: 3
2026/03/12 12:25:34 (Done)
2026/03/12 12:25:48 (Start)
2026/03/12 12:25:48 Filecount after pass (1/3): 14
2026/03/12 12:25:48 Filecount after pass (2/3): 8
2026/03/12 12:25:48 Filecount after pass (3/3): 8
2026/03/12 12:25:48 Groups of duplicates after shrink: 3
2026/03/12 12:25:48 Files with hash: 7368ac39295432a153b1532cacf30c1a4b55cc94c246d6cce820a42c06ff8c2f
2026/03/12 12:25:48  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileB.txt size: 20
2026/03/12 12:25:48  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileBDup.txt size: 20
2026/03/12 12:25:48  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileBDup.txt size: 20
2026/03/12 12:25:48  -- Duplicates: 3
2026/03/12 12:25:48 Files with hash: a4978f74fe60dbc373e48f0486d767c8d866a8f94a45c661acf812e44d978a38
2026/03/12 12:25:48  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileC.txt size: 22
2026/03/12 12:25:48  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt size: 22
2026/03/12 12:25:48  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileCDup.txt size: 22
2026/03/12 12:25:48  -- Duplicates: 3
2026/03/12 12:25:48 Files with hash: 6f430d148a85e1475301f9bd44463cc8dc69bbc1a0e059eb7c7314734e8db6dd
2026/03/12 12:25:48  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileDDup.txt size: 30
2026/03/12 12:25:48  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileD.txt size: 30
2026/03/12 12:25:48  -- Duplicates: 2
2026/03/12 12:25:48 (Done)


```


<h2> Running with log on my home directory, note that the full run took 19 seconds.</h2>

```python
$ hashcomparefiles -log homeDir.log /home/gainax2k1/
 | Files processed: 359200 
Filecount after pass (1/3): 359233
Filecount after pass (2/3): 328556
Filecount after pass (3/3): 262734
Groups of duplicates after shrink: 81641

$ more homeDir.log 
2026/03/12 12:28:03 (Start)
2026/03/12 12:28:04 Filecount after pass (1/3): 359233
2026/03/12 12:28:06 Filecount after pass (2/3): 328556
2026/03/12 12:28:21 Filecount after pass (3/3): 262734
2026/03/12 12:28:21 Groups of duplicates after shrink: 81641
2026/03/12 12:28:21 Files with hash: 55a1de25329d50620a805b6cc39c800919f99fac14f9615c0bf579e56c6b7cc8
2026/03/12 12:28:21  - /home/gainax2k1/.local/share/flatpak/appstream/flathub/x86_64/8197362a174695bcc392c6595d0926780957528594840f54e56d42e0ca6240b
6/icons/128x128/com.jetbrains.PhpStorm.png size: 4147
2026/03/12 12:28:21  - /home/gainax2k1/.local/share/flatpak/repo/objects/a1/eeac0ac49f21c0663c73d43fa5db8ac203cfa9cee3884655dfc6cd583d1546.file size
: 4147
2026/03/12 12:28:21  -- Duplicates: 2

< cut for ammount of output >

2026/03/12 12:28:22 Files with hash: ed662514dbe6ee54801768a1a5a6fe1b825780fe2c276e5798298d808f106234
2026/03/12 12:28:22  - /home/gainax2k1/.local/share/Steam/steamapps/common/Proton - Experimental/files/share/wine/mono/wine-mono-10.4.1/lib/mono/4.6-api/Facades/System.Linq.Queryable.dll size: 5120
2026/03/12 12:28:22  - /home/gainax2k1/.local/share/Steam/steamapps/common/Proton - Experimental/files/share/wine/mono/wine-mono-10.4.1/lib/mono/4.6.1-api/Facades/System.Linq.Queryable.dll size: 5120
2026/03/12 12:28:22  - /home/gainax2k1/.local/share/Steam/steamapps/common/Proton - Experimental/files/share/wine/mono/wine-mono-10.4.1/lib/mono/4.6.2-api/Facades/System.Linq.Queryable.dll size: 5120
2026/03/12 12:28:22  -- Duplicates: 3
2026/03/12 12:28:22 Files with hash: cd63bfbbc16cb7b06b9a8dbc2fbf68122e9c768307b7e8649b4b717c9628e717
2026/03/12 12:28:22  - /home/gainax2k1/.local/share/flatpak/repo/objects/2c/aaaddfdba3a3229ff11bb5164866b614f6a184953c1f4dd4da6b0894aa8ac6.file size: 24472
2026/03/12 12:28:22  - /home/gainax2k1/.local/share/flatpak/runtime/org.gnome.Platform/x86_64/48/c7cdf5b885ebc70910207b6a4d62ffca3458036edc14d0972ca090b3797b52c3/files/lib/x86_64-linux-gnu/frei0r-1/multiply.so size: 24472
2026/03/12 12:28:22  -- Duplicates: 2
2026/03/12 12:28:22 (Done)
$
```

<h2> Running with min/max flags on home: </h2>

```python
$ hashcomparefiles -min 2048 -max 8128 -log minMaxHome.log /home/gainax2k1/
 / Files processed: 360100 
Filecount after pass (1/3): 360128
Filecount after pass (2/3): 71221
Filecount after pass (3/3): 60849
Groups of duplicates after shrink: 20555
$ hashcomparefiles -log fullHome.log /home/gainax2k1/
 / Files processed: 360100 
Filecount after pass (1/3): 360120
Filecount after pass (2/3): 329443
Filecount after pass (3/3): 263331
Groups of duplicates after shrink: 81669
$
```

<h2> Piping in list of files and folders, then selectively removing them: </h2>

```python
$ cat testFilesList.txt | hashcomparefiles -log piped.log  -remove
2026/03/12 12:42:24 (Start)
2026/03/12 12:42:24 Filecount after pass (1/3): 11
Filecount after pass (1/3): 11
2026/03/12 12:42:24 Filecount after pass (2/3): 9
Filecount after pass (2/3): 9
2026/03/12 12:42:24 Filecount after pass (3/3): 9
Filecount after pass (3/3): 9
2026/03/12 12:42:24 Groups of duplicates after shrink: 4
Groups of duplicates after shrink: 4
2026/03/12 12:42:24 Files with hash: 7368ac39295432a153b1532cacf30c1a4b55cc94c246d6cce820a42c06ff8c2f
2026/03/12 12:42:24  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileBDup.txt size: 20
2026/03/12 12:42:24  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileBDup.txt size: 20
2026/03/12 12:42:24  -- Duplicates: 2
Remove file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileBDup.txt?
 - (D)elete, (T)rash, (S)kip, (C)ontinue to next hash > s
Remove file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileBDup.txt?
 - (D)elete, (T)rash, (S)kip, (C)ontinue to next hash > t
2026/03/12 12:42:31 Deleted duplicate file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileBDup.txt
2026/03/12 12:42:31 Files with hash: a4978f74fe60dbc373e48f0486d767c8d866a8f94a45c661acf812e44d978a38
2026/03/12 12:42:31  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt size: 22
2026/03/12 12:42:31  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileCDup.txt size: 22
2026/03/12 12:42:31  -- Duplicates: 2
Remove file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt?
 - (D)elete, (T)rash, (S)kip, (C)ontinue to next hash > c
2026/03/12 12:42:35 Files with hash: 6f430d148a85e1475301f9bd44463cc8dc69bbc1a0e059eb7c7314734e8db6dd
2026/03/12 12:42:35  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileD.txt size: 30
2026/03/12 12:42:35  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileD.txt size: 30
2026/03/12 12:42:35  -- Duplicates: 2
Remove file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileD.txt?
 - (D)elete, (T)rash, (S)kip, (C)ontinue to next hash > c
2026/03/12 12:42:37 Files with hash: c0f5efbef0fe98aa90619444250b1a5eb23158d6686f0b190838f3d544ec85b9
2026/03/12 12:42:37  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileA.txt size: 10
2026/03/12 12:42:37  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileADup.txt size: 10
2026/03/12 12:42:37  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileADup.txt size: 10
2026/03/12 12:42:37  -- Duplicates: 3
Remove file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileA.txt?
 - (D)elete, (T)rash, (S)kip, (C)ontinue to next hash > s
Remove file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileADup.txt?
 - (D)elete, (T)rash, (S)kip, (C)ontinue to next hash > d
2026/03/12 12:42:43 Deleted duplicate file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileADup.txt
Remove file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileADup.txt?
 - (D)elete, (T)rash, (S)kip, (C)ontinue to next hash > d
2026/03/12 12:42:44 Deleted duplicate file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileADup.txt
2026/03/12 12:42:44 (Done)
$ more piped.log 
2026/03/12 12:38:32 (Start)
2026/03/12 12:38:32 Filecount after pass (1/3): 7
2026/03/12 12:38:32 Filecount after pass (2/3): 4
2026/03/12 12:38:32 Filecount after pass (3/3): 4
2026/03/12 12:38:32 Groups of duplicates after shrink: 2
2026/03/12 12:38:32 Files with hash: a4978f74fe60dbc373e48f0486d767c8d866a8f94a45c661acf812e44d978a38
2026/03/12 12:38:32  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt size: 22
2026/03/12 12:38:32  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileCDup.txt size: 22
2026/03/12 12:38:32  -- Duplicates: 2
2026/03/12 12:38:40 Deleted duplicate file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt
2026/03/12 12:38:41 Files with hash: 6f430d148a85e1475301f9bd44463cc8dc69bbc1a0e059eb7c7314734e8db6dd
2026/03/12 12:38:41  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileD.txt size: 30
2026/03/12 12:38:41  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileD.txt size: 30
2026/03/12 12:38:41  -- Duplicates: 2
2026/03/12 12:38:43 (Done)
2026/03/12 12:42:24 (Start)
2026/03/12 12:42:24 Filecount after pass (1/3): 11
2026/03/12 12:42:24 Filecount after pass (2/3): 9
2026/03/12 12:42:24 Filecount after pass (3/3): 9
2026/03/12 12:42:24 Groups of duplicates after shrink: 4
2026/03/12 12:42:24 Files with hash: 7368ac39295432a153b1532cacf30c1a4b55cc94c246d6cce820a42c06ff8c2f
2026/03/12 12:42:24  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileBDup.txt size: 20
2026/03/12 12:42:24  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileBDup.txt size: 20
2026/03/12 12:42:24  -- Duplicates: 2
2026/03/12 12:42:31 Deleted duplicate file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileBDu
p.txt
2026/03/12 12:42:31 Files with hash: a4978f74fe60dbc373e48f0486d767c8d866a8f94a45c661acf812e44d978a38
2026/03/12 12:42:31  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileCDup.txt size: 22
2026/03/12 12:42:31  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileCDup.txt size: 22
2026/03/12 12:42:31  -- Duplicates: 2
2026/03/12 12:42:35 Files with hash: 6f430d148a85e1475301f9bd44463cc8dc69bbc1a0e059eb7c7314734e8db6dd
2026/03/12 12:42:35  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileD.txt size: 30
2026/03/12 12:42:35  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileD.txt size: 30
2026/03/12 12:42:35  -- Duplicates: 2
2026/03/12 12:42:37 Files with hash: c0f5efbef0fe98aa90619444250b1a5eb23158d6686f0b190838f3d544ec85b9
2026/03/12 12:42:37  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testFileA.txt size: 10
2026/03/12 12:42:37  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileADup.txt size: 10
2026/03/12 12:42:37  - /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileADup.txt size: 10
2026/03/12 12:42:37  -- Duplicates: 3
2026/03/12 12:42:43 Deleted duplicate file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFileADup.txt
2026/03/12 12:42:44 Deleted duplicate file: /home/gainax2k1/Documents/workspace/hashcomparefiles/testdata/testSubFolder/testFolderNested/testFileADu
p.txt
2026/03/12 12:42:44 (Done)


```

<hr>
<h2> "time -v" output run of home directory, uncached:</h2>

```
$ /usr/bin/time -v hashcomparefiles /home/gainax2k1/
 | Files processed: 418400 
Filecount after pass (1/3): 418493
Filecount after pass (2/3): 386924
Filecount after pass (3/3): 303634
Groups of duplicates after shrink: 88734
	Command being timed: "hashcomparefiles /home/gainax2k1/"
	User time (seconds): 27.19
	System time (seconds): 15.88
	Percent of CPU this job got: 200%
	Elapsed (wall clock) time (h:mm:ss or m:ss): 0:21.48
	Average shared text size (kbytes): 0
	Average unshared data size (kbytes): 0
	Average stack size (kbytes): 0
	Average total size (kbytes): 0
	Maximum resident set size (kbytes): 300480
	Average resident set size (kbytes): 0
	Major (requiring I/O) page faults: 35
	Minor (reclaiming a frame) page faults: 103506
	Voluntary context switches: 1064335
	Involuntary context switches: 7981
	Swaps: 0
	File system inputs: 27459720
	File system outputs: 286544
	Socket messages sent: 0
	Socket messages received: 0
	Signals delivered: 0
	Page size (bytes): 4096
	Exit status: 0

```

That should be everything! Any issues with the project, feel free to reach out to me. 
Thanks, and have a day. =^.^=
