# Unnatural

Unnatural is a an elf anomaly detector written in golang for linux systems.
It is capable of detecting malware infected binaries and disinfect those binaries to some extent

goals:
- detect and disinfect data segment infections
- detect and disinfect PT_NOTE infections
- Profiling view

ToDo:
- Detect Reverse text padding infections

Note:
- currently broken, cause im writing a parser from scratch since it is hard to work with debug/elf
