## HiveShadow

This is an implementation of CVE-2021–36934 written in GO (1.16.3)

### Usage 

```bash
Arguments:
         -h      Print this message
         -q      Quick wins - only scans for the first shadow copy
         -b      Brute force shadow copy number up to max depth (default 20)
         -d      Brute force max depth (default 20)
         -o      The output directory (make sure you can write here)

Example:
         .\GoHiveShadow.exe -b -d 20 -o C:\Windows\Temp

```

Original code: https://github.com/GossiTheDog/HiveNightmare